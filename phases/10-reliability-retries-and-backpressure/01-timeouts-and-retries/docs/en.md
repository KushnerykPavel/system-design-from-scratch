# Timeouts, Retries, and Retry Storms

> A retry is a load multiplier, not a free reliability feature.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Choose timeout and retry policy from dependency behavior, fanout, and overload risk instead of repeating blanket rules like "retry three times."
**Prerequisites:** `02-estimation-and-cost/05-burstiness`, `03-design-framework-and-timing/01-four-step-interview-loop`, `07-queues-streams-and-workflows/07-pipeline-backpressure`
**Estimated time:** ~75 min
**Primary artifact:** retry policy evaluator + retry storm checklist

## The Problem

In interview answers, retries are often presented as a pure availability improvement:

- client times out
- client retries
- system becomes more reliable

Real systems behave differently. Bad retries can turn a small latency incident into a fleet-wide outage by multiplying in-flight work, queue depth, and tail latency.

This lesson trains a better instinct:

- set timeouts from deadlines and percentiles
- propagate deadlines across hops
- retry only when the operation and failure mode justify it
- cap retry amplification before it caps your service

## Clarify

- Is the call path single-hop or fanout-heavy?
- Are operations safe to retry, or can partial success create duplicate side effects?
- Is the main issue intermittent packet loss, dependency tail latency, or overload?
- Where does the deadline originate: user request, batch worker, or internal control plane?

If the interviewer is vague, assume a user-facing RPC path with 3-5 downstream calls, strict p99 latency goals, and occasional overload-driven timeouts.

## Requirements

### Functional

- Bound request latency with explicit deadlines and hop-level timeouts.
- Retry only safe operations and only on retryable failure classes.
- Avoid synchronized retries and uncontrolled request amplification.

### Non-functional

- Protect downstream services during overload instead of worsening it.
- Keep retry behavior explainable to operators and reviewers.
- Preserve enough availability gain to justify the extra logic.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Front-door QPS | 120K req/s | retry amplification gets expensive quickly |
| Downstream fanout | 4 RPCs on critical path | one client retry can trigger many backend attempts |
| Baseline p99 dependency latency | 180 ms | timeout choice must respect actual tails |
| Peak retry multiplier | up to 2 additional attempts | determines worst-case load expansion |
| Rough cost | extra attempts + queueing + wasted CPU | "availability" may be bought with expensive overload |

## Architecture

A strong answer treats retries as one part of deadline management:

```text
client deadline
  -> edge / API
     -> per-hop timeout derived from remaining budget
     -> jittered retry policy for safe failures
     -> dependency
```

Useful rules:

1. Set a request deadline at the caller.
2. Derive shorter hop timeouts from remaining budget.
3. Retry only transient failures such as connect resets or very short overload windows.
4. Add jitter and retry caps to avoid herd behavior.
5. Stop retrying when the request budget is already mostly spent.

## Data Model & APIs

Useful policy fields:

- `deadline_ms`
- `timeout_ms`
- `max_attempts`
- `backoff_strategy`
- `retryable_statuses`
- `idempotent`

Useful interfaces:

- `CallWithPolicy(request, policy)`
- `ShouldRetry(result, attempt, remaining_budget_ms)`
- `DeriveChildTimeout(parent_deadline, dependency_class)`

Senior-level detail:

- keep retry classification separate from business success
- make remaining budget part of the retry decision
- distinguish transport failures from application failures

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| timeout too short causes self-inflicted retries | timeout rate rises while downstream success on second attempt stays low | tune from real latency percentiles and deadline budget |
| every layer retries independently | attempt-per-request ratio spikes during incidents | central retry ownership and propagated deadlines |
| synchronized backoff creates retry waves | retries cluster on the same boundaries | add full jitter or decorrelated jitter |
| retries hit non-idempotent operations | duplicate writes or side effects appear after timeout | gate retries on idempotency guarantees |

## Observability

- metric: attempts per original request and retry amplification factor
- metric: timeout rate by dependency and by attempt number
- metric: success-on-retry rate versus total retry volume
- metric: remaining deadline budget at final failure
- log: retry decision with attempt number, failure class, and remaining budget
- trace: parent request and each child attempt correlated under one span tree
- SLO: retries should improve success rate without materially increasing overload during dependency incidents

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| derived per-hop deadlines | keeps total request bounded | requires careful budget propagation | one large static timeout everywhere |
| capped jittered retries | can recover transient failures | extra load and more complex debugging | unlimited immediate retries |
| retry only safe failure classes | reduces duplicate side effects | some transient app failures go un-retried | retry every non-200 response |

## Interview It

**Google framing:** "A user-facing service times out on a dependency during peak traffic. How would you make it more reliable?" The signal is whether you discuss deadline propagation and load amplification, not just backoff.

**Cloudflare framing:** "A global API path sees transient failures at the edge-to-origin hop." The signal is whether you protect origins from retry storms while keeping tail latency bounded.

**Follow-ups:**
1. What changes when the call path fans out to five dependencies?
2. Which failures should never be retried?
3. When do retries help less than admission control?
4. How do you stop nested services from multiplying retries?
5. What would you change at 10x traffic?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-timeouts-and-retries.md`
- `outputs/retry-storm-checklist.md`

## Exercises

1. **Easy** — Pick a timeout policy for a single dependency with p95 `40 ms` and p99 `180 ms`.
2. **Medium** — Explain how retries change when the service fans out to four downstream RPCs.
3. **Hard** — Redesign the policy when a dependency is overloaded and client retries are making the incident worse.

## Further Reading

- [Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — canonical reliability framing for retry storms and overload
- [The Tail at Scale](https://research.google/pubs/pub40801/) — strong background on latency tails and why naive retries are dangerous
