# Open Connect CDN — Cache Hierarchy & POP Design

> Netflix built its own CDN not because Akamai could not deliver the bits, but because the economics and the control made it necessary at their scale.

**Type:** Build  
**Company focus:** Netflix  
**Learning goal:** Design a tiered CDN with ISP-embedded servers. Understand proactive vs reactive fill, cache eviction by popularity, ISP peering, and the economic and operational reasons Netflix built Open Connect instead of buying third-party CDN.  
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/`, `02-video-streaming-pipeline`  
**Estimated time:** ~90 min  
**Primary artifact:** CDN design doc + cache hierarchy spec  

## The Problem

Netflix needs to deliver hundreds of terabits per second to 270M subscribers worldwide. At this volume, buying CDN capacity from Akamai or CloudFront is cost-prohibitive. More importantly, ISPs have an incentive to participate because Netflix traffic can represent 30%+ of their peak downstream traffic — hosting a Netflix server saves them transit costs.

Design Open Connect: Netflix's own CDN embedded directly inside ISP networks.

## Clarify

- What is the cache hierarchy? (embedded ISP servers → regional Netflix POPs → Netflix cloud origin)
- Should fill be proactive (push popular content before it is requested) or reactive (pull on first cache miss)?
- How should popularity be determined for prefill decisions?
- What is the SLA for a cache miss from embedded server to regional origin?
- How should servers handle content that changes (re-encode, new DRM keys)?

## Requirements

### Functional

- Serve video segments and manifests from servers physically located inside ISP networks.
- Fill cache hierarchically: embedded server → regional Open Connect Appliance (OCA) → Netflix cloud origin.
- Prefill popular content proactively during off-peak hours.
- Evict unpopular content to make room for new or more popular content.
- Support content invalidation when a title is re-encoded or removed.

### Non-functional

- Cache hit ratio: target >99% for top-N popular titles.
- Segment delivery latency: p99 under 200ms from embedded server.
- Fill bandwidth: use off-peak hours (2–6 AM local) to minimize ISP transit cost impact.
- Storage per embedded server: 100–200 TB (commodity hardware).

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak downstream bandwidth | ~200 Tbps global | drives number of embedded servers required |
| Distinct titles actively watched | top ~1,000 titles drive ~80% of traffic | drives cache size requirements per server |
| Storage per OCA server | 100–200 TB | limits how much catalog can be proactively filled |
| Fill window | 4–6 hours nightly | drives fill bandwidth requirements |
| Cache miss latency to regional POP | ~20ms | acceptable for rare misses |

## Architecture

```text
subscriber device
  -> embedded OCA server (inside ISP, e.g. Comcast datacenter)
     [cache hit] -> serve directly
     [cache miss] -> regional Netflix POP
        [cache hit] -> fill embedded OCA + serve
        [cache miss] -> Netflix cloud origin (S3)
           -> fill regional POP + fill embedded OCA + serve
```

### Proactive Prefill

During nightly off-peak windows, Open Connect pushes popular content to embedded servers before subscribers request it:

```text
popularity_score(title, variant) = request_count_last_7d * recency_weight
top_K = sort(all_variants, by=popularity_score, desc)[:storage_budget]
for each variant in top_K not already on server:
    fetch from regional POP
    write to local storage
    update catalog index
```

Key design decisions:
- **Title ranking is ISP-local**: a regional ISP in Brazil gets different top-K than a server in Germany. Popularity is geographically segmented.
- **Variant selection**: not all variants of a top title fit. Prioritize the most-requested resolution+bitrate combination per region.
- **Prefill during low-traffic windows**: minimizes impact on subscriber experience and ISP transit costs.

### Reactive Fill

On a cache miss, the embedded server fetches from the nearest Netflix regional POP:

```text
miss_handler(request):
    chunk = fetch(regional_pop_url(request.segment_key))
    cache_locally(chunk, ttl=ttl_for_popularity(chunk.title_id))
    return chunk
```

Reactive fill handles the long tail of the catalog that cannot fit in the proactive fill budget.

### Cache Eviction

Each embedded server runs a popularity-weighted eviction policy:

| Eviction dimension | Rule |
|--------------------|------|
| Access frequency | least-frequently-accessed segments evicted first |
| Recency | content not accessed in N days is candidate for eviction |
| Title-level granularity | evict all variants of a title together to preserve storage alignment |
| Free space threshold | trigger eviction when free space drops below 10% |

### ISP Peering vs Transit

| Path | Cost | Latency |
|------|------|---------|
| Embedded server → subscriber | near-zero (ISP internal) | <5ms |
| Regional POP → subscriber (peering) | low | 10–30ms |
| Netflix cloud origin → ISP (transit) | high | 30–100ms |

Netflix negotiates free colocation of OCA hardware with ISPs in exchange for reducing the ISP's transit costs. The ISP benefits because the Netflix server handles traffic that would otherwise cross expensive upstream transit links.

## Data Model & APIs

OCA catalog entry:
```text
segment_key, title_id, variant_id, segment_index, local_path, cached_at, last_accessed_at, access_count
```

Popularity signal (sent from Netflix cloud to each OCA):
```text
{title_id, variant_key, popularity_rank, isp_id, updated_at}
```

Cache miss event (sent to Netflix monitoring):
```text
{oca_id, segment_key, miss_reason, fill_source, fill_latency_ms, timestamp}
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| OCA server disk failure | Health check fails, no heartbeat | Route traffic to regional POP; alert Netflix ops; schedule replacement hardware |
| Regional POP unreachable | OCA fill requests timeout | Fall back to second-nearest POP; ultimately fall back to cloud origin |
| Prefill fills wrong content for a region | Cache hit rate drops; wrong-region requests increase | Popularity signals are ISP-tagged; verify geo-segmentation in fill scheduler |
| Content is pulled (DMCA or editorial) | Invalidation message sent to all OCAs | OCA deletes segments by title_id within N minutes; manifests stop referencing removed titles |
| Fill bandwidth saturates ISP uplink | Fill throughput metric hits cap | Rate-limit fill during peak hours; defer low-priority fills to next nightly window |

## Observability

- metric: cache hit ratio by OCA, ISP, region, and title
- metric: fill latency and fill bandwidth consumed by window
- metric: eviction rate by reason (popularity-based vs age-based)
- metric: miss escalation rate (OCA miss → POP miss → origin hit)
- metric: OCA health (disk utilization, CPU, network throughput)
- log: prefill decisions with popularity score, variant chosen, and bytes written
- alert: cache hit ratio drops below 95% for any ISP cluster

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Proactive prefill over pure reactive | Hit rate for top titles approaches 100%; no subscriber waits for cold cache | Fill scheduling complexity; wasted fill if popularity shifts suddenly | Pure reactive simpler but cold cache causes origin load spikes at launch |
| ISP-local popularity ranking | Accurate regional demand; better hit rates per region | Popularity computation must be distributed per ISP | Global ranking would fill wrong content for regional ISPs |
| Title-level eviction granularity | Keeps variant sets coherent; avoids partial title in cache | Eviction granularity is coarser, may waste space on low-quality variants of popular titles | Per-segment eviction gives finer control but fragments title storage |
| Own CDN over Akamai | Cost control at extreme scale; ISP peering economics; faster iteration | Massive capex, hardware ops team, ISP relationship management | Third-party CDN viable at 10x lower scale but cost prohibitive for Netflix volume |

## Interview It

**Netflix framing:** "Why did Netflix build Open Connect instead of using a third-party CDN?" Strong answers discuss economics (transit cost savings), control (proactive prefill, custom eviction), and operational depth (ISP relationships, hardware choices). Weak answers say "latency" without explaining the underlying economics.

**Follow-ups:**
1. How does Netflix decide how much storage to allocate per ISP?
2. What happens during a major content release (e.g., a popular series finale) when demand spikes 10x?
3. How do you handle a content invalidation for a title across thousands of OCA servers globally?
4. What is the fallback path if an entire region's OCA fleet goes dark?
5. How do you measure and improve cache hit ratio over time?

## Ship It

- `outputs/design-doc-open-connect-cdn.md`
- `outputs/cache-hierarchy-spec.md`
- `outputs/interview-card-open-connect-cdn.md`

## Exercises

1. **Easy** — Calculate the fill bandwidth required to prefill the top 500 titles (100 GB average, 4-hour fill window) per embedded server.  
2. **Medium** — Design the popularity signal pipeline: how does Netflix compute per-ISP title popularity from playback events?  
3. **Hard** — Design the content invalidation system that removes a pulled title from thousands of OCA servers globally within 15 minutes.  

## Further Reading

- [Netflix Open Connect overview](https://openconnect.netflix.com/en/)  
- [Netflix ISP speed index](https://ispspeedindex.netflix.com/)  
- [Peering vs transit economics](https://blog.cloudflare.com/the-internet-is-a-system-of-agreements/)  
