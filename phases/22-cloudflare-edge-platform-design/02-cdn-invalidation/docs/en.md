# CDN Invalidation and Cache Propagation

> Invalidation is where CDN answers stop being theoretical. The hard part is not deleting bytes. It is deciding how much inconsistency, fanout, and control-plane risk the product can tolerate.

**Type:** Learn
**Company focus:** Cloudflare
**Learning goal:** Design cache invalidation and propagation strategies for globally distributed CDNs without hiding control-plane latency, purge fanout, or safety trade-offs.
**Prerequisites:** `06-caching-and-invalidation/02-freshness-models`, `13-multi-region-cdn-and-edge-traffic/03-cdn-layering`, `22-cloudflare-edge-platform-design/01-global-api-edge-gateway`
**Estimated time:** ~75 min
**Primary artifact:** purge strategy matrix + rollout checklist

## The Problem

Design CDN invalidation for a platform serving static assets, API-adjacent cached objects, and customer-managed content. Customers want fast purges, predictable propagation, and enough visibility to trust the platform during incidents.

This is a Cloudflare-style lesson because the interesting part is not just cache TTL. It is global purge fanout, cache-key design, customer blast radius, and control-plane safety under high churn.

## Clarify

- Are purges mostly single-object, prefix, tag, or full-zone purges?
- Is the main priority fastest visible purge, lowest control-plane cost, or strongest safety?
- What stale window is acceptable at the edge, shield, or browser layers?
- Can the product tolerate soft purge semantics before hard deletion lands everywhere?

## Requirements

### Functional

- Support object, prefix, and tag-based invalidation.
- Propagate purge intent to POPs and any intermediate cache layers.
- Expose purge status and failure visibility to operators and customers.
- Prevent one tenant from overwhelming the global purge pipeline.

### Non-functional

- Keep purge propagation operationally safe during bursts.
- Bound stale-serving windows explicitly.
- Avoid global control-plane meltdowns from high-cardinality purge patterns.
- Preserve explainability when a customer claims a purge did not work.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Cached objects | 10B+ keys globally | drives metadata scale and cache-key discipline |
| Purge rate | 50K requests/min bursty | shapes control-plane queueing and dedupe |
| Purge fanout | 100s of POPs per request | makes propagation cost visible |
| Common purge scope | mostly object and tag, rare full zone | affects indexing strategy |
| Acceptable stale window | seconds to low minutes depending on product | determines soft versus hard purge design |

## Architecture

```text
customer purge request
  -> authenticated purge API
  -> validation and quota checks
  -> purge control plane
  -> dedupe / batching / sequencing
  -> POP propagation channel
  -> local cache tombstone or version update
  -> lazy refill from origin on next miss
```

Key ideas:

1. Separate purge intent distribution from cache eviction execution.
2. Prefer versioned or tombstoned invalidation over expensive global scans.
3. Treat purge as a control-plane workload with quotas and rollout safety.
4. Make stale windows explicit by cache layer.

## Data Model & APIs

Useful records:

```text
purge_request(
  tenant_id,
  scope_type,
  scope_value,
  request_id,
  requested_at,
  priority,
  status
)
```

Helpful APIs:

- `POST /purges`
- `GET /purges/{request_id}`
- `ListPurgeProgress(tenant_id, scope)`
- `ExplainCacheState(key, pop)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| purge storm from one tenant | queue depth and per-tenant rate spikes | quotas, batching, and backpressure |
| POPs apply purge intents out of order | version skew or stale-hit anomalies | monotonic versioning and idempotent apply |
| tag purge index becomes too large | index growth and lookup latency rise | bounded tag cardinality and compaction rules |
| control plane is healthy but POP delivery lags | propagation lag metrics widen by POP | replayable streams and lag alarms |
| customers expect browser cache to purge too | support tickets after CDN purge success | document cache-layer semantics and cache-control defaults |

## Observability

- metric: purge propagation latency p50/p95/p99
- metric: stale-hit rate after purge by cache layer
- metric: per-tenant purge volume and throttling
- metric: POP apply lag and replay backlog
- log: purge request validation and rejection reasons
- trace: purge API -> control plane -> POP apply -> first post-purge miss

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| versioned invalidation | avoids scanning caches globally | requires key/version metadata discipline | delete-by-scan across POPs |
| tag purge support | great customer ergonomics | expensive metadata and fanout | object-only purge API |
| soft purge first, hard purge later | fast visible effect | more temporary complexity | waiting for full hard purge before any response |

## Interview It

**Google framing:** "Design a distributed cache invalidation system." Strong answers still discuss versioning, propagation lag, and observability, but may focus less on POP behavior.

**Cloudflare framing:** "Design global CDN purge." Strong answers must discuss control-plane fanout, per-tenant abuse control, cache-key discipline, and how operators explain stale content claims.

**Follow-ups:**
1. What changes if customers demand sub-second purge for a hot object?
2. What if full-zone purges are common during deployments?
3. What if some POPs are partitioned from the control plane?
4. What if purge cost must be surfaced as a billable dimension?

## Ship It

- `outputs/tradeoff-matrix-cdn-invalidation.md`
- `outputs/rollout-checklist-cdn-invalidation.md`

## Exercises

1. **Easy** — Compare object purge versus tag purge for cost and operator visibility.
2. **Medium** — Redesign for customers that need near-real-time purge of a product catalog.
3. **Hard** — Handle a regional control-plane lag incident while keeping stale-serving predictable.

## Further Reading

- [Cache invalidation](https://en.wikipedia.org/wiki/Cache_invalidation) — useful framing for why freshness is hard
- [Cloudflare engineering blog](https://blog.cloudflare.com/) — practical edge and purge system context
