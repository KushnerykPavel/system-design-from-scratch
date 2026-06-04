# Cache Hit Rate and Origin Load

> Cache hit rate is not a vanity metric when it directly decides whether your origin survives peak traffic.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Translate cache hit rate into origin QPS, capacity savings, and failure exposure instead of treating it as a decorative percentage.  
**Prerequisites:** `02-estimation-and-cost/01-qps-and-request-mix`, `06-caching-and-invalidation/01-cache-patterns`  
**Estimated time:** ~75 min  
**Primary artifact:** origin-load worksheet + interview card  

## The Problem

Candidates often say "add a cache" without quantifying what it buys. A cache with the wrong hit rate, object mix, or invalidation policy may barely help the origin under real peak traffic.

This lesson turns hit ratio into concrete origin QPS and explains how to reason about miss pain.

## Clarify

- Is the cache in front of reads only, or can it absorb write-adjacent work too?
- Is the workload dominated by a hot working set or a long tail?
- Are misses expensive because of origin latency, cost, or rate limits?
- Does the cache protect one origin tier or multiple downstream dependencies?

## Requirements

### Functional

- Estimate origin QPS from frontend QPS and cache hit ratio.
- Compare multiple hit-rate scenarios quickly.
- Show how hit-rate drops affect cost and overload risk.

### Non-functional

- Make miss cost visible, not just hit percentage.
- Keep the math simple enough for live discussion.
- Expose whether a cache is reducing average load or true peak pain.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Frontend read QPS | 90K | cache demand |
| Cache hit rate | 92% | determines miss pressure |
| Origin miss QPS | 7.2K | actual origin read load |
| Miss penalty | 40 ms extra | shapes latency savings |
| Hit-rate drop scenario | 92% to 75% | tests degradation risk |

## Architecture

Reason in three numbers:

1. Frontend read QPS.
2. Cache hit rate.
3. Miss penalty in QPS, latency, or dollars.

At 90K read QPS:

- 92% hit means 7.2K origin QPS
- 75% hit means 22.5K origin QPS

That drop is not "17 percentage points worse." It is more than 3x the origin load.

## Data Model & APIs

The code artifact models:

```text
CacheModel {
  ReadQPS
  HitRatio
  MissLatencyMillis
}
```

Outputs:

- origin QPS
- cache-served QPS
- weighted average latency contribution

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hit rate quoted globally only | hot keys or cold objects hide true misses | break down by route and object class |
| invalidation event drops hit rate | origin overloads during rollout | pre-warm or stage invalidations |
| cache capacity too small for working set | hit ratio decays under peak | size memory to the actual hot set |
| miss cost ignored | hit-rate drop seems harmless on paper | quantify origin QPS and latency amplification |

## Observability

- metric: hit ratio by route and content class
- metric: origin QPS attributable to cache misses
- metric: miss latency distribution
- metric: cache eviction rate
- SLO: origin headroom remains safe even during expected hit-rate dips

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| focus on origin QPS saved | clearer capacity impact | less elegant than quoting one ratio | celebrating hit rate without context |
| optimize hottest objects first | largest capacity return | leaves long tail less protected | uniform treatment for all objects |
| push higher hit rate aggressively | lowers cost and origin risk | harder invalidation and warmer complexity | accepting more misses |

## Interview It

**Google framing:** "Design a product catalog service with heavy read traffic." The signal is whether you quantify what a cache actually saves the datastore.

**Cloudflare framing:** "Protect origin for dynamic-but-cacheable API responses." The signal is whether you connect hit ratio to origin shielding under bursts.

**Follow-ups:**
1. What if 1% of objects receive 80% of requests?
2. What if misses trigger expensive database joins?
3. What if freshness requirements cap TTL at 5 seconds?
4. What if a deployment invalidates most keys at once?

## Ship It

- `outputs/origin-load-worksheet-cache-hit-rate.md`
- `outputs/interview-card-cache-hit-rate.md`

## Exercises

1. **Easy** — Estimate origin QPS for a 98% hit static asset service.  
2. **Medium** — Compare two caches: 85% hit with cheap misses vs 95% hit with expensive invalidation.  
3. **Hard** — Explain why weighted object classes matter more than a single fleet-wide hit number.  

## Further Reading

- [Cloudflare blog](https://blog.cloudflare.com/) — strong practical material on caching and origin protection  
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview context  
