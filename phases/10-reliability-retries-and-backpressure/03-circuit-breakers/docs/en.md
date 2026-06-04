# Circuit Breakers and Graceful Degradation

> Failing fast is only useful if the user experience degrades on purpose instead of by accident.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Combine circuit breakers, fallback behavior, and user-visible degradation so a sick dependency stops burning the rest of the system.
**Prerequisites:** `06-caching-and-invalidation/02-freshness-models`, `08-consistency-replication-and-transactions/04-replication-lag`, `10-reliability-retries-and-backpressure/01-timeouts-and-retries`
**Estimated time:** ~75 min
**Primary artifact:** breaker policy evaluator + graceful degradation matrix

## The Problem

When one dependency becomes slow or error-prone, many interview answers stop at "add a circuit breaker."

That is incomplete. A breaker only changes whether we keep calling the dependency. The real design work is:

- what traffic should be cut off
- what degraded answer can still be served
- how to probe recovery safely
- how to keep the fallback path from sharing the same failure

## Clarify

- Is the dependency optional, critical, or critical only for some request classes?
- What degraded behavior is acceptable: stale data, partial response, queue for later, or explicit error?
- How quickly do we need recovery probes once the dependency improves?
- Can some callers bypass the breaker for high-priority flows?

If the interviewer is vague, assume a user-facing service with one optional enrichment dependency and one critical core dependency, both of which can fail independently.

## Requirements

### Functional

- Detect dependency distress and stop sending unlimited traffic into failure.
- Serve a degraded but intentional experience when possible.
- Probe recovery without unleashing a full retry wave.

### Non-functional

- Keep overload from cascading into healthy components.
- Make degradation mode visible to operators and product owners.
- Bound the risk that the fallback path is stale or misleading.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Critical path QPS | 90K req/s | determines how fast a bad dependency can poison callers |
| Optional dependency share | 65% of requests use it | degradation can preserve most value if the fallback is credible |
| Recovery probe budget | <1% of normal volume | half-open behavior must stay bounded |
| Fallback freshness target | 30-120 seconds stale allowed | constrains cache or static fallback design |
| Rough cost | fallback storage + probe traffic + complexity | graceful degradation is not free, but outages are worse |

## Architecture

A good breaker design separates three ideas:

1. **Trip policy** based on recent failures, timeouts, or saturation.
2. **Fallback policy** based on feature criticality.
3. **Recovery probe policy** that limits reopened traffic.

```text
request
  -> dependency classifier
  -> breaker state machine
  -> primary dependency OR fallback path
  -> sampled recovery probes when half-open
```

## Data Model & APIs

Useful policy fields:

- `trip_error_rate`
- `trip_timeout_rate`
- `evaluation_window`
- `half_open_max_probes`
- `fallback_mode`
- `priority_class`

Useful interfaces:

- `Allow(primary, request_class)`
- `Fallback(request_class)`
- `RecordOutcome(dependency, result)`

Senior-level detail:

- fallback can be stale cache, partial response, async completion, or explicit shed
- breaker state should usually be per dependency class, not global across everything
- recovery probes should be small and cancelable

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| breaker trips too late and overload spreads | error rate, queue depth, and tail latency climb before state changes | combine saturation and error signals, not only hard failures |
| fallback path shares same dependency or datastore | degraded mode fails with the primary | isolate fallback dependencies and precompute data where possible |
| half-open allows too much traffic back at once | second outage wave during partial recovery | cap probe concurrency and ramp slowly |
| degraded response confuses users or downstreams | support tickets and semantic mismatch errors rise | define product contract for stale or partial behavior explicitly |

## Observability

- metric: breaker open rate, half-open probe count, and fallback hit ratio
- metric: user-visible success rate in normal and degraded modes
- metric: dependency error rate, timeout rate, and queue depth driving trip decisions
- log: breaker state transitions with reason, threshold, and dependency class
- trace: tags for normal path versus fallback path
- SLO: degradation should preserve a defined minimum service level for non-critical features during dependency incidents

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| per-dependency breaker | smaller blast radius | more policy tuning | one global breaker for all calls |
| stale fallback data | preserves availability | weaker freshness and potential user confusion | hard fail every request |
| cautious half-open probing | avoids re-triggering outage | slower recovery confirmation | instant full reopen |

## Interview It

**Google framing:** "One downstream service becomes flaky and slow. How do you prevent a cascade?" The signal is whether you pair the breaker with user-facing degradation and safe recovery.

**Cloudflare framing:** "An edge feature depends on a control-plane lookup that is intermittently unavailable." The signal is whether you can keep the data plane serving with bounded staleness.

**Follow-ups:**
1. Which features should degrade versus fail closed?
2. What if the fallback cache is stale for two minutes?
3. What if only premium traffic should get recovery probes first?
4. How do you avoid a false trip during a brief blip?
5. What changes if the dependency is billing or auth and cannot degrade safely?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-circuit-breakers.md`
- `outputs/tradeoff-matrix-circuit-breakers.md`

## Exercises

1. **Easy** — Define one acceptable degraded response for an optional profile-enrichment dependency.
2. **Medium** — Explain the half-open probe policy for a dependency recovering from overload.
3. **Hard** — Redesign the breaker strategy when the dependency is auth and incorrect fallback is worse than temporary unavailability.

## Further Reading

- [Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — strong foundation for load shedding and isolation patterns
- [Release It!](https://pragprog.com/titles/mnee2/release-it-second-edition/) — classic operational framing for breakers and stability patterns
