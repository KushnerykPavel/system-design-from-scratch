# Token Bucket vs Sliding Window vs Leaky Bucket

> A senior answer does not just pick a rate-limiting algorithm. It explains what failure or fairness property that algorithm is buying.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Choose the right limiter primitive for burst handling, memory cost, fairness, and operator simplicity instead of treating all "rate limiting" as the same mechanism.  
**Prerequisites:** `02-estimation-and-cost/05-burstiness`, `06-caching-and-invalidation/03-eviction-policies`, `10-reliability-retries-and-backpressure/04-load-shedding`  
**Estimated time:** ~75 min  
**Primary artifact:** primitive trade-off matrix  

## The Problem

Design the primitive behind a request limiter for an API or shared platform. You need to decide whether the system should absorb bursts, smooth output, approximate rolling windows, or enforce stricter fairness under skew.

Interview candidates often say "use a rate limiter" and move on. Stronger answers explain why token bucket, sliding window, and leaky bucket produce different system behavior under bursty traffic and partial overload.

## Clarify

- Is the business trying to cap abuse, protect a backend, or shape traffic smoothly?
- Should short spikes be tolerated if the long-term average is acceptable?
- Is fairness across adjacent seconds more important than implementation simplicity?
- How much memory and per-request work can the hot path afford?

If no more detail is given, assume public API traffic with bursty clients, low-latency enforcement, and moderate tolerance for approximation.

## Requirements

### Functional

- Enforce a configurable per-key rate.
- Support burst behavior that matches product expectations.
- Provide predictable reject decisions at high QPS.

### Non-functional

- Keep per-request decision cost low.
- Avoid excessive per-key memory amplification.
- Make operator tuning understandable during incidents.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak QPS | 300K req/s | shapes hot-path efficiency requirements |
| Active keys | 5M/day, 200K hot/hour | drives state size and eviction policy |
| Allowed burst | 5x steady rate for 1 to 3 seconds | changes whether smoothing or bucket semantics fit |
| Decision budget | sub-millisecond | rules out expensive per-request logs for the common path |
| Rough cost | in-memory state + optional shared coordination | primitive choice affects both compute and memory |

## Architecture

```text
request
  -> identify key
  -> evaluate limiter primitive
  -> optionally update shared state
  -> allow / reject
```

Primitive intuition:

1. **Token bucket** permits burst accumulation and cheap decisions.
2. **Sliding window** approximates fairness over a rolling interval.
3. **Leaky bucket** smooths egress and protects a downstream worker pool.

The architecture question is not just "what algorithm?" It is also where state lives, how often it is refreshed, and whether exactness is worth the cost.

## Data Model & APIs

Common state:

```text
key -> {
  policy,
  current_state,
  last_update_ms
}
```

Useful APIs:

- `Check(key, now, cost)`
- `Refill(key, now)`
- `UpdatePolicy(key_class, rule)`
- `ExplainDecision(key, now)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| token bucket allows too much short-term spikiness | backend latency rises despite compliant average rate | lower burst or add downstream concurrency limits |
| sliding window state grows too large | memory pressure on hot nodes | use approximations or coarser buckets |
| leaky bucket is misused for abuse control | queue delay rises before rejects happen | separate shaping from hard enforcement |
| team argues about algorithm labels instead of workload shape | vague discussion with no concrete constraints | anchor on bursts, fairness, and memory budget |

## Observability

- metric: allows and rejects by key class and limiter primitive
- metric: per-key state count and eviction rate
- metric: burst absorption versus smoothed throughput
- log: sampled reject explanation including algorithm and remaining budget
- trace: attach limiter decision to downstream latency spikes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| token bucket for public APIs | simple and burst friendly | adjacent-second fairness is approximate | exact timestamp log for every request |
| sliding window for stricter fairness | better rolling-window behavior | more state and math on hot path | pure token bucket when fairness matters more |
| leaky bucket for worker protection | smooths output rate toward fragile backends | queueing can hide overload | immediate reject-only policy |

## Interview It

**Google framing:** "Choose the limiter primitive for a shared quota platform." Expect pressure on approximation, memory growth, and what changes when limits span multiple windows.

**Cloudflare framing:** "Choose the right edge enforcement primitive." Expect questions about bursty clients, latency at the POP, and whether smoothing belongs in the data path.

**Follow-ups:**
1. What changes if customers care about fairness inside every rolling second?
2. When is a leaky bucket the wrong answer for abuse prevention?
3. How would you combine a token bucket with a monthly quota?
4. What if one key becomes 20% of total traffic?
5. Which primitive is easiest to explain during an incident review?

## Ship It

- `outputs/tradeoff-matrix-rate-limiter-primitives.md`
- `outputs/interview-card-rate-limiter-primitives.md`

## Exercises

1. **Easy** — Pick a primitive for a login endpoint and justify it.
2. **Medium** — Redesign for a backend that cannot tolerate bursts at all.
3. **Hard** — Support both per-second burst limits and per-hour tenant quotas in one answer.

## Further Reading

- [Cloudflare - How we built rate limiting capable of scaling to millions of domains](https://blog.cloudflare.com/counting-things-a-lot-of-different-things/) — practical trade-offs around distributed enforcement  
- [System design notes](https://github.com/liquidslr/system-design-notes) — canonical interview framing for limiter primitives  
