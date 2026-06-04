# Distributed Rate Limiter

> The hard part is not counting requests. It is counting them fast enough, fairly enough, and safely enough under skew and failure.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design a multi-node rate limiter that supports local fast paths, shared enforcement, and bounded inconsistency under failure.  
**Prerequisites:** `06-caching-and-invalidation/03-eviction-policies`, `09-partitioning-sharding-and-rebalancing/02-hot-partitions`, `10-reliability-retries-and-backpressure/04-load-shedding`  
**Estimated time:** ~90 min  
**Primary artifact:** trade-off matrix + failure checklist  

## The Problem

Design a rate limiter for an API gateway fleet. Any request can land on any edge node. Limits may be per API key, per customer, or per route pattern. Tail latency must stay low, but enforcement still needs to be meaningfully correct at high scale.

This is a canonical interview system because it combines data-path speed, distributed coordination, hot-key handling, and degraded-mode reasoning.

## Clarify

- Is this enforcing short-term rate, long-term quota, or both?
- What are the dimensions of the key: user, IP, API key, route, or tenant?
- What is the tolerated inconsistency window during partial failures?
- Should the system fail open or fail closed when coordination storage is unhealthy?

## Requirements

### Functional

- Per-key rate enforcement with configurable burst allowance.
- Shared policy propagation to all edge nodes.
- Optional multi-dimensional keys like `tenant + route + API key`.

### Non-functional

- p99 decision latency should stay in the low-millisecond or sub-millisecond range.
- The fleet should withstand skewed hot keys and partial backing-store outages.
- Operators must be able to explain reject decisions and tune policy safely.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak QPS | 1M across the fleet | determines coordination pressure and local cache value |
| Active keys | 10M/day, 500K hot in one hour | affects memory footprint and eviction policy |
| Decision latency | <1 ms local, <3 ms with coordination | constrains storage placement and fallback choices |
| Peak factor | 5x burst on hot tenants | drives burst policy and hotspot planning |
| Rough cost | shared store + local memory on each edge | forces trade-off between strictness and speed |

## Architecture

```text
client
  -> edge gateway
     -> local hot-key cache / token bucket
     -> consistent-hash owner or shared store
     -> allow / reject
```

Typical layers:

1. **Policy control plane** distributes limit rules.
2. **Local fast path** handles hot keys with bounded TTL.
3. **Shared enforcement layer** provides cross-node coordination.
4. **Metrics + audit path** explains why traffic was rejected.

## Data Model & APIs

Per-key state:

```text
key -> {
  tokens,
  burst,
  last_refill_ms,
  version
}
```

Useful APIs:
- `Check(key, cost, now)`
- `UpdatePolicy(policy_version, rules)`
- `Explain(key)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hot key overloads one shard | per-shard reject and latency spikes | local fast path + virtual nodes + emergency caps |
| shared store is slow | store latency metrics breach threshold | bounded fail-open or stricter local-only fallback |
| stale policy on some edges | policy version mismatch metrics | versioned rollout and expiry on old rules |
| duplicate retries inflate counts | mismatch between gateway retries and limiter counts | idempotency keys or retry-budget-aware counting |

## Observability

- metric: allows vs rejects by tenant, route, and policy class
- metric: shared-store p95/p99 latency
- metric: local cache hit ratio
- metric: policy version skew across edges
- log: sampled reject explanations with key class and rule ID
- trace: limiter decision attached to request path

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| local cache before shared store | very fast hot-path decisions | bounded inconsistency window | every request hits central store |
| token bucket | simple burst handling | only approximates some long-window limits | sliding log with higher memory cost |
| fail-open on store outage | preserves availability | weaker enforcement during incidents | fail-closed that can become self-inflicted outage |

## Interview It

**Google framing:** "Design a quota platform used by internal services." Expect questions about hot tenants, policy propagation, and when consistency actually matters.

**Cloudflare framing:** "Design edge rate limiting for global API traffic." Expect questions about POP-local latency, origin protection, abuse resistance, and key explosion.

**Follow-ups:**
1. What changes if customers need hourly and monthly quotas together?
2. What if one tenant is 40% of fleet traffic?
3. What if policy updates must propagate globally within 10 seconds?
4. What if some requests are retries and should not double count?
5. What changes if you must prove reject decisions to customers?

## Ship It

- `outputs/tradeoff-matrix-distributed-rate-limiter.md`
- `outputs/failure-checklist-distributed-rate-limiter.md`
- `outputs/interview-card-distributed-rate-limiter.md`

## Exercises

1. **Easy** — Add an explicit decision for whether policy should fail open or fail closed.  
2. **Medium** — Extend the design to support shared quota pools across multiple API keys.  
3. **Hard** — Redesign for multi-region limits where traffic can hit any region and still respect a global quota budget.  

## Further Reading

- [Counting things a lot of different things](https://blog.cloudflare.com/counting-things-a-lot-of-different-things/) — practical edge rate limiting trade-offs  
- [System design notes - rate limiter chapter](https://github.com/liquidslr/system-design-notes) — canonical interview framing  
