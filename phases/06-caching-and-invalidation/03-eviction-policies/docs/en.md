# Eviction Policies and Working Set Shape

> An eviction policy only looks smart when it matches the workload.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Compare eviction policies against real access traces and explain why working-set shape matters more than policy slogans.
**Prerequisites:** `02-estimation-and-cost/04-cache-hit-rate`, `05-storage-indexing-and-access-patterns/04-hot-and-cold-data`
**Estimated time:** ~75 min
**Primary artifact:** policy comparison simulator + trade-off matrix

## The Problem

Interview answers often name LRU, LFU, or FIFO as if one policy were universally best. In practice, eviction only makes sense relative to workload shape:

- short recency bursts
- long-term popularity
- scans that pollute the cache
- large hot sets that barely fit in memory

This lesson includes a small Go simulator so the learner can compare policies against an access trace and see how hit rate changes with capacity and skew.

## Clarify

- Is the workload dominated by bursty recency, repeated popularity, or wide scans?
- How large is the estimated working set compared with cache capacity?
- Are all objects roughly equal in size, or do large objects distort eviction value?
- Is protecting origin QPS more important than maximizing raw hit rate?

## Requirements

### Functional

- Compare at least three eviction policies on the same trace.
- Show how cache capacity changes hit rate.
- Explain which trace shapes help or hurt each policy.

### Non-functional

- Keep the tool small enough to understand during a lesson.
- Avoid pretending toy hit-rate results generalize without workload context.
- Connect simulator output back to real interview decisions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak reads | 150K req/s | enough to make a few points of hit rate valuable |
| Working set | 8M objects, 150K hot in an hour | defines whether memory is tight or generous |
| Cache capacity | 20K to 200K keys in examples | reveals policy sensitivity to size |
| Scan traffic | 5-15% on some workloads | tests pollution resistance |
| Rough cost | memory per key plus origin miss cost | ties policy choice to business trade-offs |

## Architecture

The teaching artifact simulates a fixed-capacity cache over an access trace:

1. feed accesses in order
2. record hit or miss
3. evict when capacity is full
4. report hit rate per policy

This is intentionally small. The goal is to make working-set reasoning visible, not to build a production cache.

## Data Model & APIs

The simulator treats accesses as string keys and supports:

- `Run(trace, capacity, policy)`
- `Compare(trace, capacities, policies)`

Policies included:

- `LRU`
- `LFU`
- `FIFO`

The most important interview point is not the exact simulator number. It is the explanation of why one policy wins on one trace and loses on another.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| policy chosen from folklore, not workload | weak explanation despite naming a famous policy | model the dominant access trace first |
| scan traffic destroys cache locality | sudden hit-rate collapse during wide sequential reads | segmented cache or admission control |
| LFU holds stale once-hot keys too long | old hot keys remain while new hot keys miss | aging, decay, or hybrid recency/frequency policy |
| memory budget is too small for working set | low hit rate across all policies | admit less, tier data, or increase cache where justified |

## Observability

- metric: hit rate segmented by endpoint and key class
- metric: eviction rate and churn
- metric: one-hit-wonder share or scan detection
- log: sampled large miss bursts and policy changes
- trace: cache participation on the highest-cost origin endpoints
- SLO: hit ratio target for the hottest read paths with origin protection as the operational goal

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| LRU | simple and good for recency-heavy traffic | vulnerable to scans | using only FIFO by default |
| LFU | strong for stable popularity | stale frequency can linger too long | assuming popularity never changes |
| simulator-based lesson | makes trade-offs concrete | still a simplification of production caches | memorizing policy definitions only |

## Interview It

**Google framing:** "Design a cache for a product feed where traffic is bursty and celebrity posts dominate briefly." The signal is whether you discuss recency, skew, and pollution rather than only naming LRU.

**Cloudflare framing:** "Design cache behavior for edge objects with hot keys and occasional broad scans." The signal is whether you reason about working set size and protecting the origin under adversarial patterns.

**Follow-ups:**
1. What if objects are not the same size?
2. What if a background job scans millions of cold keys?
3. When is hit rate the wrong main metric?
4. What if yesterday's hot keys should stop dominating today's cache?
5. Would you change the policy or the admission rule first?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/tradeoff-matrix-eviction-policies.md`
- `outputs/interview-card-eviction-policies.md`

## Exercises

1. **Easy** — Compare LRU and FIFO on a short bursty trace.
2. **Medium** — Explain why LFU can outperform LRU on stable popularity but still hurt after trend shifts.
3. **Hard** — Redesign the lesson for weighted objects where one large key displaces many smaller hot keys.

## Further Reading

- [Caching at Scale at Facebook](https://engineering.fb.com/2013/04/15/core-infra/scaling-memcache-at-facebook/) — practical examples of cache pressure and working-set thinking
- [System design notes](https://github.com/liquidslr/system-design-notes) — useful interview context for placing cache policy inside a broader system answer
