# Bandwidth and Egress Cost

> Throughput is a performance concern first and a billing concern immediately after.

**Type:** Build  
**Company focus:** Cloudflare  
**Learning goal:** Estimate bandwidth load and egress cost so cache, compression, and regional serving choices stay grounded.  
**Prerequisites:** `02-estimation-and-cost/01-qps-and-request-mix`, `06-caching-and-invalidation/06-cache-layers`  
**Estimated time:** ~75 min  
**Primary artifact:** trade-off matrix + egress worksheet  

## The Problem

Many interview answers notice request count but miss response size. A service with modest QPS can still be dominated by bandwidth if objects are large, cache misses are frequent, or cross-region replication is noisy.

This lesson focuses on rough throughput math and how it changes edge placement and cost posture.

## Clarify

- Are responses tiny metadata reads or large media payloads?
- Does the number need to include client egress, origin egress, or both?
- What fraction of requests are cache hits at the edge or CDN?
- Is traffic concentrated in one geography or spread globally?

## Requirements

### Functional

- Estimate peak bandwidth from QPS and average payload size.
- Estimate origin bandwidth after cache hit reduction.
- Translate transferred data into rough monthly egress cost.

### Non-functional

- Keep the arithmetic quick enough for interview use.
- Make cache effectiveness visible in the cost model.
- Highlight when egress dominates infra cost.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak QPS | 120K | traffic anchor |
| Average response size | 220 KB | drives throughput |
| Edge cache hit rate | 85% | changes origin load dramatically |
| Peak origin bandwidth | about 4 GB/s | affects origin links and cost |
| Rough monthly egress | petabyte scale | can dominate operating budget |

## Architecture

Simple flow:

1. Estimate total peak bytes per second.
2. Reduce origin load by cache hit rate.
3. Convert to monthly transferred data.
4. Apply rough per-GB egress pricing.
5. Ask whether compression, cache, or edge compute can reduce it.

Example:

- 120K QPS
- 220 KB average response
- about 25.8 GB/s total edge-serving throughput
- at 85% cache hit, origin sees about 3.9 GB/s

## Data Model & APIs

The code artifact models:

```text
BandwidthModel {
  PeakQPS
  ResponseKB
  CacheHitRatio
  CostPerGB
}
```

Useful outputs:

- total served GB/s
- origin GB/s
- monthly origin egress GB
- monthly origin egress cost

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| response size underestimated | real bandwidth exceeds link budget | size by percentile payload classes |
| cache hit too optimistic | origin cost explodes | show sensitivity around hit ratio |
| cross-region traffic ignored | inter-region bill surprises team | count replication and cache fill traffic separately |
| compression assumed but not achieved | expected savings never appear | measure payload compression ratio directly |

## Observability

- metric: bytes served by route, object class, and region
- metric: edge hit ratio and origin miss bytes
- metric: compressed vs uncompressed payload ratio
- metric: egress cost by service or tenant
- SLO: origin bandwidth stays below planned headroom during peak events

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| invest in higher cache hit rate | lowers origin bandwidth and cost | more invalidation and cache-management complexity | accepting repeated origin fetches |
| compress aggressively | saves egress | adds CPU and latency overhead | shipping raw payloads |
| push content to edge | reduces central bottlenecks | adds cache consistency concerns | serving everything from origin |

## Interview It

**Google framing:** "Design an image delivery pipeline for a social app." The signal is whether you connect object size and cacheability to both cost and topology.

**Cloudflare framing:** "Design global asset delivery with origin protection." The signal is whether you separate edge throughput from origin egress and think about cache fills.

**Follow-ups:**
1. What if large objects are only 10% of QPS but 90% of bytes?
2. What if egress pricing differs sharply by region?
3. What if cache hit rate drops from 85% to 60% during a rollout?
4. What if compression helps text but not already-compressed media?

## Ship It

- `outputs/tradeoff-matrix-bandwidth-and-egress.md`
- `outputs/egress-worksheet-bandwidth-and-egress.md`

## Exercises

1. **Easy** — Estimate origin bandwidth for a 95% cache-hit static asset service.  
2. **Medium** — Split traffic into text API payloads and large image payloads.  
3. **Hard** — Explain how regional pricing and cache-fill traffic change the cost model.  

## Further Reading

- [Cloudflare blog](https://blog.cloudflare.com/) — many practical discussions about bandwidth, caching, and origin protection  
- [System design notes](https://github.com/liquidslr/system-design-notes) — helpful interview framing  
