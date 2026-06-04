# Consistency Trade-offs in Cached Systems

> A cache is a consistency decision with a latency upside.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Explain the consistency envelope of a cached system in product terms, including what users may observe after writes, failures, and regional divergence.
**Prerequisites:** `01-cache-patterns`, `02-freshness-models`, `08-consistency-replication-and-transactions/01-consistency-spectrum`
**Estimated time:** ~60 min
**Primary artifact:** cached consistency checklist

## The Problem

Interview answers often claim "eventual consistency" without saying who experiences it, for how long, or under what failure. That is too vague. In cached systems, the real question is what a user can observe right after a write or invalidation event.

This lesson helps you describe:

- read-after-write behavior
- monotonic-read expectations
- cross-region freshness differences
- degraded-mode behavior when invalidation or replication is delayed

## Clarify

- Which user action must reflect immediately in a follow-up read?
- Is the main risk stale convenience data or stale safety-critical policy?
- Can users tolerate seeing old data as long as they never see time go backward?
- Are reads local to one region, or can the same user bounce between regions?

## Requirements

### Functional

- Serve low-latency reads from cache where possible.
- Define what freshness or monotonicity is promised after writes.
- Handle invalidation, replication lag, and regional divergence explicitly.

### Non-functional

- Make the consistency contract understandable to users and operators.
- Avoid accidental stronger guarantees than the system can truly provide.
- Keep degraded-mode behavior safe under partial failure.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read QPS | 250K req/s | large enough that cache is necessary |
| Write QPS | 8K req/s | enough to produce a steady invalidation stream |
| Cross-region share | 20% of sessions | exposes monotonic-read and freshness questions |
| Allowed stale window | 0-10 seconds by feature | determines whether local caches are acceptable |
| Rough cost | tighter invalidation, regional coordination, extra misses | makes stronger guarantees visibly expensive |

## Architecture

A strong cached consistency story names:

1. **Source of truth**
2. **Cache layers**
3. **Invalidation path**
4. **User-visible guarantee**

Examples:

- profile bio: bounded stale reads are acceptable
- price or entitlement: stale reads may be unacceptable immediately after write
- abuse or policy rule: fail-safe may matter more than hit ratio

Common patterns:

- session stickiness after write for monotonic reads
- write-through for hot read-after-write paths
- versioned reads where the client or server can require at least version `N`

## Data Model & APIs

Helpful metadata:

```text
record -> {
  value,
  version,
  committed_at,
  cache_fetched_at,
  region
}
```

Helpful interfaces:

- `Get(key, min_version)`
- `Invalidate(key, version)`
- `GetFresh(key)` for limited correctness-critical paths

If you promise more than TTL-based eventual consistency, explain the mechanism that actually enforces it.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| user reads stale data right after own write | read-after-write mismatch metric | sticky reads, write-through, or version-aware bypass |
| invalidation reaches some regions late | regional version skew | bounded TTL and version-based reconciliation |
| client alternates between regions and sees time go backward | monotonic-read violations | session affinity or client-carried version token |
| safety-critical rule served stale from cache | policy freshness alerts | shorter TTL, explicit push invalidation, or cache bypass on critical path |

## Observability

- metric: version skew between cache and source
- metric: read-after-write mismatch rate on critical flows
- metric: invalidation end-to-end latency by region
- log: stale-read incidents with cache age and version
- trace: write request linked to subsequent reads and invalidations
- SLO: freshness objective for critical entities paired with latency objective

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| bounded stale reads for noncritical content | cheaper and faster reads | users may briefly see old state | forcing every read to hit origin |
| version-aware bypass on critical reads | tighter correctness where it matters | extra metadata and selective miss cost | one uniform rule for every entity |
| session affinity after mutation | monotonic reads for one user | uneven load and routing complexity | global strict consistency everywhere |

## Interview It

**Google framing:** "Design cached reads for a collaborative product where edits should appear quickly to the editor but can lag slightly for others." The signal is whether you differentiate user classes and guarantees.

**Cloudflare framing:** "Design policy lookup caches where stale state can affect enforcement." The signal is whether you articulate fail-safe behavior, version skew, and propagation windows.

**Follow-ups:**
1. Which entities deserve a cache bypass after write?
2. What if users roam across regions within seconds?
3. What if invalidation is reliable but replication to the source follower lags?
4. What if stricter freshness doubles origin cost?
5. How would you explain the consistency envelope in one sentence to the interviewer?

## Ship It

- `outputs/cached-consistency-checklist.md`

## Exercises

1. **Easy** — Define the stale-read window you would accept for a profile page.
2. **Medium** — Design read-after-write behavior for inventory after a purchase.
3. **Hard** — Explain the consistency story for globally cached abuse rules and what happens under propagation lag.

## Further Reading

- [Designing Data-Intensive Applications](https://dataintensive.net/) — strong background on read-after-write, monotonic reads, and replication trade-offs
- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline for communicating trade-offs under interview pressure
