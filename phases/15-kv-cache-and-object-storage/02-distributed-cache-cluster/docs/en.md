# Distributed Cache Cluster

> Caches look simple until miss storms, rebalancing, and skew reveal that most of the real work is operational.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design a distributed cache cluster that balances hit rate, memory efficiency, warmup behavior, and failure isolation under fleet-scale traffic.  
**Prerequisites:** `06-caching-and-invalidation/04-cache-stampede`, `09-partitioning-sharding-and-rebalancing/03-placement`, `14-rate-limiters-ids-and-hashing/04-consistent-hashing`  
**Estimated time:** ~75 min  
**Primary artifact:** topology checker + observability checklist  

## The Problem

Design a shared cache cluster used by many stateless services. The cluster should absorb read load from databases or APIs, but it must not become a hidden single point of failure or an amplifier during incident conditions.

Weak answers treat the cache as free performance. Strong answers explain working set shape, eviction policy, miss amplification, and what happens when nodes are added, lost, or cold.

## Clarify

- Is this cache protecting a database, an API, or an expensive computation?
- Are values reconstructable, or would cache loss create user-visible data loss?
- Is the workload dominated by a small hot set or a broad long tail?
- What staleness or TTL model is acceptable?

If no detail is given, assume reconstructable values, read-heavy traffic, hot-key skew, and a strong desire to avoid origin collapse during misses.

## Requirements

### Functional

- Serve low-latency reads for cached objects by key.
- Support TTL, explicit invalidation, and negative caching where appropriate.
- Rebalance smoothly when nodes are added or removed.

### Non-functional

- Keep miss storms from cascading into origin outages.
- Maintain predictable memory usage and eviction behavior.
- Make cache effectiveness visible enough that teams can tune it.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak read QPS | 2M req/s | cache is a hot path service, not a sidecar detail |
| Working set | 300 GB hot, 2 TB long tail | memory sizing and eviction policy dominate design |
| Object size | 0.5 to 32 KB | affects item overhead and shard packing |
| Hit-rate target | 92% overall, 99% for hottest 5% of keys | determines origin protection value |
| Node churn | weekly scale events + occasional failures | cold shard handling and consistent hashing matter |

## Architecture

```text
client service
  -> key hashing / shard selection
  -> cache node or replica pair
     -> local memory
     -> optional replication / persistence
  -> miss path
     -> request coalescing
     -> origin fetch
     -> populate cache
```

Important cluster behaviors:

1. Consistent hashing reduces key movement during scale changes.
2. Request coalescing protects origins during popular misses.
3. Replication or replica warming should match business pain, not habit.
4. Admission control can stop one noisy team from turning the cluster into a churn machine.

## Data Model & APIs

Cached entry:

```text
cache_key
value
ttl
inserted_at
size_bytes
version_token
```

Useful APIs:

- `Get(key)`
- `Set(key, value, ttl)`
- `Delete(key)`
- `Warm(keys)`
- `ExplainMiss(key)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| cache node loss triggers miss storm | miss rate and origin QPS spike together | request coalescing, circuit breaking, and warmup throttles |
| one tenant pollutes memory with low-value entries | eviction churn and low per-tenant hit rate | admission policies, tenant caps, or segmented pools |
| rebalancing moves too many hot keys | hit rate cliff after scale event | consistent hashing, virtual nodes, and staged warmup |
| negative cache entries outlive reality | stale not-found errors after data appears | shorter TTLs and versioned invalidation paths |

## Observability

- metric: hit rate, miss rate, and origin offload by service and tenant
- metric: eviction rate and item age distribution
- metric: request coalescing savings during miss storms
- metric: key movement and hit-rate drop after rebalancing
- log: sampled miss explanations with cause such as expired, evicted, or cold shard
- SLO: cache service latency plus protected-origin error budget consumption

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| shared cache cluster | high aggregate efficiency and easier warm hot set sharing | noisy-neighbor risk and operational coupling | per-service private caches only |
| consistent hashing | less data movement during scaling | imperfect balance under skew unless tuned | modulo hashing with massive key churn |
| request coalescing | origin protection on popular misses | coordination overhead for pending lookups | every miss independently hits origin |

## Interview It

**Google framing:** "Design a shared cache tier protecting high-QPS backend services." Expect questions on hit-rate economics, invalidation, and when teams should bypass the shared layer.

**Cloudflare framing:** "Design a cache cluster serving hot objects at edge or regional layers." Expect questions on locality, key movement, and degraded behavior when cache nodes churn.

**Follow-ups:**
1. How would you isolate a tenant with pathological large objects?
2. What changes if origin latency is highly variable?
3. How do you keep cold-starts from wiping out hit rate during autoscaling?
4. When is replication inside the cache worth it?
5. What changes if some values must never be stale for more than five seconds?

## Ship It

- `outputs/interview-card-distributed-cache-cluster.md`
- `outputs/observability-checklist-distributed-cache-cluster.md`
- `outputs/tradeoff-matrix-distributed-cache-cluster.md`

## Exercises

1. **Easy** — Choose an eviction policy for a skewed read-heavy workload and justify it.
2. **Medium** — Redesign the cluster to support large-object and small-object pools separately.
3. **Hard** — Explain how you would survive a regional cache flush without overloading the primary database.

## Further Reading

- [Scaling Memcache at Facebook](https://engineering.fb.com/2013/04/15/core-infra/scaling-memcache-at-facebook/) — strong practical lessons on cache clusters and operational pitfalls  
- [System design notes - cache chapter](https://github.com/liquidslr/system-design-notes) — canonical interview framing for cache tiers  
