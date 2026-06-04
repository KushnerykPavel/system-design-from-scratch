# Admission Control and Load Shedding

> When capacity is gone, honesty beats queueing theater.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design admission control that protects latency and critical work under overload instead of letting queues grow until everything is slow.
**Prerequisites:** `02-estimation-and-cost/07-bottleneck-math`, `07-queues-streams-and-workflows/07-pipeline-backpressure`, `10-reliability-retries-and-backpressure/03-circuit-breakers`
**Estimated time:** ~75 min
**Primary artifact:** admission policy evaluator + overload checklist

## The Problem

Overload usually starts with a good intention:

- "let's queue a little more"
- "let's keep accepting traffic"
- "let's hope the spike passes"

That often ends with:

- queues that are already too old to be useful
- CPU pegged and context switching heavily
- retries amplifying the spike
- high-priority traffic drowning with everything else

Admission control is about deciding which work the system is still willing to accept.

## Clarify

- Is overload CPU-bound, IO-bound, downstream-bound, or memory-bound?
- Which requests are critical and which are optional or deferrable?
- Can some work be queued asynchronously, or must it be rejected immediately?
- Is fairness more important than raw throughput during overload?

If the interviewer is vague, assume an API service with mixed critical and optional traffic, bursty peaks, and a requirement to preserve low latency for high-priority requests.

## Requirements

### Functional

- Enforce concurrency or queue limits before the service collapses.
- Prefer critical traffic when capacity is scarce.
- Return fast, explicit overload responses for rejected work.

### Non-functional

- Protect latency SLOs for the most important flows.
- Avoid unbounded memory growth from queues.
- Make overload policy understandable and tunable in production.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak incoming QPS | 250K req/s | overload can arrive faster than autoscaling |
| Safe concurrency | 18K in-flight | admission limits should reflect where latency breaks |
| Burst factor | 6x over baseline for 30 seconds | queueing strategy must survive spikes, not averages |
| Critical traffic share | 20% of requests | drives priority lane design |
| Rough cost | rejected work + reserved capacity + policy complexity | protecting the core path requires leaving some work out |

## Architecture

Strong admission control usually combines:

1. a concurrency cap
2. bounded queues
3. priority classes
4. fast reject behavior

```text
request
  -> classifier
  -> priority lane
  -> concurrency / token gate
  -> bounded queue or immediate reject
  -> worker pool
```

Useful pattern:

- accept critical requests first
- keep queue age bounded, not just queue length
- reject optional work early when downstreams are unhealthy

## Data Model & APIs

Useful policy fields:

- `max_inflight`
- `max_queue_depth`
- `max_queue_age_ms`
- `priority_class`
- `shed_when_dependency_unhealthy`

Useful interfaces:

- `Admit(request_class, now)`
- `RejectReason()`
- `UpdatePolicy(capacity_snapshot)`

Senior-level detail:

- queue age is often a better signal than raw length
- separate admission by request class when one flow must survive
- tie shedding to downstream saturation, not only local CPU

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| queue absorbs too much overload and hides failure until latency explodes | queue age and p99 latency rise before rejects start | strict queue-age and concurrency caps |
| high-priority traffic shares one pool with best-effort work | critical latency degrades during spikes | dedicated priority lanes or reserved tokens |
| service accepts work that downstream cannot handle | local success looks fine but dependency errors climb | couple admission to downstream health and budgets |
| overload responses trigger client retry storms | reject rate and retry amplification rise together | clear retry hints, backoff headers, and caller budgets |

## Observability

- metric: admit rate versus reject rate by priority class
- metric: queue depth, queue age, and in-flight concurrency
- metric: downstream saturation signals used by the admission policy
- metric: success rate of critical traffic during overload
- log: rejection reason, priority class, and policy version
- trace: tags for admitted, queued, and shed requests
- SLO: critical-path requests should stay within latency target during bounded overload, even if optional traffic is shed

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| fast reject under overload | protects latency and system health | some user-visible failures | queue everything and hope |
| reserved capacity for critical traffic | preserves core business flow | lower total utilization in calm periods | one shared pool |
| queue-age limit | prevents stale work from consuming capacity | more aggressive shedding | long queues with uncertain usefulness |

## Interview It

**Google framing:** "Your service stays up during bursts but latency becomes terrible. How would you fix overload behavior?" The signal is whether you reject work deliberately rather than hand-wave autoscaling.

**Cloudflare framing:** "A global edge service must keep security checks and core routing healthy during a traffic surge." The signal is whether you can classify traffic and preserve critical work.

**Follow-ups:**
1. Which traffic should be dropped first?
2. What if clients blindly retry every `429`?
3. How do you keep fairness across tenants during a burst?
4. What if the bottleneck is a downstream database, not local CPU?
5. When is queueing better than shedding?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/failure-checklist-load-shedding.md`
- `outputs/tradeoff-matrix-load-shedding.md`

## Exercises

1. **Easy** — Pick one signal you would use to stop queue growth before latency collapses.
2. **Medium** — Design separate admission behavior for premium and best-effort traffic.
3. **Hard** — Explain how you coordinate load shedding across a fleet when the shared database is the real bottleneck.

## Further Reading

- [Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — strong overload and admission-control framing
- [Google SRE books](https://sre.google/books/) — useful background on protecting SLOs under load
