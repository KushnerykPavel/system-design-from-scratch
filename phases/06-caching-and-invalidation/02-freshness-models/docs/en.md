# TTL, Explicit Invalidation, and Freshness Models

> Freshness is a product promise expressed as system behavior.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Pick a freshness model that matches user expectations and operational reality instead of defaulting to arbitrary TTLs.
**Prerequisites:** `01-cache-patterns`, `01-clarification-and-scope/05-prioritization`
**Estimated time:** ~75 min
**Primary artifact:** freshness decision checklist

## The Problem

Candidates often say "set a 5-minute TTL" as if TTL were a neutral default. It is not. TTL is a contract about how stale users may be, how much origin traffic you will pay, and how quickly updates become visible.

This lesson focuses on three common freshness mechanisms:

- time-based expiration
- explicit invalidation
- version-aware or conditional freshness

The strongest answer names the user-visible freshness promise before choosing the mechanism.

## Clarify

- Which user action most needs fresh data: create, edit, purchase, moderation, or configuration rollout?
- Is stale data merely annoying, or does it create safety, financial, or policy risk?
- Can the system tolerate bounded staleness if updates propagate quickly enough?
- Are invalidation events reliable and ordered, or should time-based expiry still exist as a safety net?

## Requirements

### Functional

- Serve hot reads with predictable latency.
- Refresh or invalidate entries when the source of truth changes.
- Bound stale windows for the most important user journeys.

### Non-functional

- Avoid cache invalidation mechanisms that are more fragile than the workload needs.
- Preserve debuggability when stale data is reported.
- Keep origin load within budget even when freshness becomes tighter.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read QPS | 220K req/s | drives cache value and miss sensitivity |
| Update rate | 4K writes/s | affects invalidation fanout and freshness pressure |
| Hot-object skew | top 0.1% keys get 35% of reads | makes invalidation efficiency more important |
| Freshness target | 500 ms to 5 min by path | determines whether TTL alone is enough |
| Rough cost | cache churn + invalidation transport + origin misses | ties freshness to real operating cost |

## Architecture

### TTL-first model

Entries expire after a fixed or jittered duration. This is the simplest model and often enough when:

- content changes slowly
- stale reads are acceptable for short windows
- invalidation events would be too expensive or unreliable

### Explicit invalidation model

Writers publish invalidation events or delete keys immediately after updating the source of truth. This fits when:

- freshness matters soon after writes
- key ownership is clear
- the invalidation path is reliable and observable

### Version-aware freshness

Responses carry versions, ETags, timestamps, or config generations. Readers can detect whether cached data is still current. This helps when:

- correctness depends on monotonic updates
- entries are reused across layers
- conditional revalidation is cheaper than full misses

Most real systems combine these. Explicit invalidation tightens freshness, while TTL limits damage if an invalidation event is lost.

## Data Model & APIs

Useful metadata:

```text
key -> {
  payload,
  version,
  expires_at,
  origin_updated_at
}
```

Useful operations:

- `Invalidate(key)`
- `InvalidatePrefix(prefix)` with great caution
- `GetIfVersionAtLeast(key, version)`
- `Revalidate(etag)`

In interviews, say clearly whether invalidation is best-effort, guaranteed, or replayable.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| TTL is too long for a hot mutable object | stale-read reports after updates | shorter TTL for mutable classes or explicit invalidation |
| invalidation event is dropped | cache version lags origin version | replayable event bus plus TTL backstop |
| invalidation fanout overwhelms cache nodes | invalidation queue depth and node CPU spike | batch invalidations, scope keys carefully, use version bumps |
| prefix invalidation blows away working set | hit ratio collapse and origin traffic surge | isolate namespaces and prefer versioned keys when possible |

## Observability

- metric: stale-read rate or version mismatch rate
- metric: invalidation publish-to-apply latency
- metric: hit ratio segmented by object mutability class
- log: invalidation failures, replay events, and version conflicts
- trace: update request through source write and cache invalidation path
- SLO: freshness for critical objects plus latency on the main read path

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| TTL for low-risk content | simple and cheap to operate | bounded but real staleness | per-write invalidation everywhere |
| explicit invalidation for hot mutable keys | tighter freshness right after writes | more moving parts and failure handling | only waiting for expiry |
| version-aware revalidation | strong debuggability and safer layering | more metadata and conditional logic | opaque values with no freshness context |

## Interview It

**Google framing:** "Design a product details service with heavy reads and frequent inventory updates." The signal is whether you distinguish stale descriptions from stale stock counts.

**Cloudflare framing:** "Design edge configuration distribution with local caches." The signal is whether you talk about propagation delay, version skew, and rollback safety.

**Follow-ups:**
1. Which fields can tolerate a longer TTL than others?
2. What if invalidation ordering is not guaranteed?
3. What if a multi-key object update must appear atomically?
4. What if invalidations cost more than misses at low write volume?
5. How do you prove whether a stale read came from the cache or source lag?

## Ship It

- `outputs/freshness-decision-checklist.md`

## Exercises

1. **Easy** — Pick TTL and invalidation behavior for a static marketing page cache.
2. **Medium** — Pick freshness rules for product price, description, and stock count separately.
3. **Hard** — Design a freshness strategy for firewall or routing policy rollout where stale data can create safety issues.

## Further Reading

- [HTTP Caching](https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching) — strong grounding for freshness and revalidation semantics
- [Caching best practices and max-age gotchas](https://web.dev/http-cache/) — practical cache-control behavior across layers
