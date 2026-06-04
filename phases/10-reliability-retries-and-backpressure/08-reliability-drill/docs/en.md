# Reliability Drill

> Reliability answers are only persuasive if they survive overload, ambiguity, and operator scrutiny together.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Integrate the phase by designing one system that needs timeouts, retries, idempotency, degradation, admission control, async backpressure, retry budgets, and bulkheads at the same time.
**Prerequisites:** `10-reliability-retries-and-backpressure/07-bulkheads`
**Estimated time:** ~60 min
**Primary artifact:** drill sheet + answer checker

## The Problem

Use this prompt:

"Design a global webhook delivery platform that accepts customer events, stores delivery intent durably, retries failed destinations, protects shared infrastructure during customer abuse or receiver outages, and gives operators clear reliability controls."

This drill is good because it forces the whole phase together:

- sync API reliability on ingest
- idempotent write semantics
- async retry and backpressure behavior
- tenant isolation and load shedding
- explicit observability and failure handling

## Clarify

- Are webhooks best-effort or contractually must-deliver?
- What is the maximum useful delay before a queued delivery attempt becomes stale?
- Can one customer's broken receiver be isolated without impacting others?
- What retry limits, ordering guarantees, and visibility do customers expect?

If the interviewer is vague, assume durable delivery intent, at-least-once semantics, strong per-tenant isolation, and a need to protect the shared platform from pathological receivers.

## Requirements

### Functional

- Accept webhook events durably and idempotently.
- Deliver with bounded retries and customer-visible status.
- Isolate abusive or failing destinations from healthy tenants.

### Non-functional

- Preserve platform health during receiver outages and customer retry bursts.
- Keep ingest latency low even when async delivery pipelines are stressed.
- Make backlog, retry behavior, and degraded modes observable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Ingest QPS | 150K events/s peak | front door must remain healthy under bursts |
| Delivery fanout | 1-20 endpoints per event | retries can multiply actual outbound attempt volume |
| Failure episode duplicate rate | 10% repeated submits | ingest path needs idempotency and replay safety |
| Maximum useful delivery lag | 5-30 minutes by plan tier | determines backlog policy and stale work handling |
| Rough cost | durable ingest + retry storage + per-tenant isolation | reliability features must remain economically credible |

## Architecture

A strong answer usually includes:

1. idempotent ingest API with durable intent record
2. async delivery queue with per-tenant pressure accounting
3. retry policy with budgets, jitter, and receiver-aware pause
4. load shedding and bulkheads to isolate bad tenants or endpoints
5. explicit degraded modes for status reads and optional analytics

```text
client
  -> ingest API
  -> idempotency + durable event log
  -> per-tenant delivery queues
  -> worker pools with retry budgets
  -> destination endpoints
  -> status / metrics / operator controls
```

## Data Model & APIs

Core entities:

- `WebhookIntent`
- `IdempotencyRecord`
- `DeliveryAttempt`
- `TenantQuota`
- `EndpointHealth`

Core APIs:

- `SubmitWebhook(tenant_id, idempotency_key, payload)`
- `GetDeliveryStatus(intent_id)`
- `PauseEndpoint(endpoint_id)`
- `ReplayFailedDeliveries(filter)`

The best answers keep ingest fast and safe even while delivery is degraded.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one receiver outage causes fleet-wide retry storm | attempt rate and queue age spike for that destination class | per-endpoint pause, retry budgets, and tenant isolation |
| duplicate client submit after ambiguous timeout creates duplicate intent | duplicate intent count or status mismatch rises | durable idempotency keys on ingest |
| backlog grows faster than workers can drain | queue age and expired work ratio rise steadily | publish throttling, stale-drop policy, or tiered quotas |
| status API depends on the same hot path as delivery workers | support visibility disappears during delivery incidents | separate read model or independently scaled status path |

## Observability

- metric: ingest QPS, accepted versus rejected, and idempotency hit ratio
- metric: queue age, attempt rate, success rate, and retry budget burn by tenant and endpoint class
- metric: shed rate, paused endpoints, and pool utilization across bulkheads
- metric: stale delivery share and replay backlog
- log: retry decisions, endpoint pauses, and idempotency conflicts
- trace: ingest request through enqueue and sampled delivery attempt path
- SLO: ingest availability and bounded delivery freshness for healthy tenants should survive isolated receiver incidents

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| durable ingest plus async delivery | fast front door and controllable retries | more state and eventual delivery | synchronous direct delivery on submit |
| per-tenant and per-endpoint isolation | healthy customers stay healthy | more queues and policy complexity | one global retry backlog |
| bounded retries with stale-drop policy | platform stays alive under bad receivers | some deliveries may age out | infinite retry with hidden backlog |

## Interview It

**Google framing:** "Design a large-scale webhook delivery system." Expect follow-ups on duplicate suppression, retry budgets, and how overload is contained.

**Cloudflare framing:** "Design a globally distributed event delivery platform that protects the shared fleet from bad customer endpoints." Expect pushback on backpressure, per-tenant isolation, and degradation strategy.

**Follow-ups:**
1. What if one tenant sends 100x more events than the median?
2. What if a major destination starts timing out globally?
3. How do you prevent operators from replaying a failure storm into the same outage?
4. What if some deliveries are legally required while others are best-effort?
5. Which metric tells you first that the system is becoming unreliable?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-reliability-drill.md`
- `outputs/reliability-drill-sheet.md`

## Exercises

1. **Easy** — Deliver a 10-minute version of the drill focused only on the highest-signal reliability choices.
2. **Medium** — Re-run the drill assuming premium tenants require much tighter delivery freshness than free-tier tenants.
3. **Hard** — Redesign the system when a downstream provider outage lasts two hours and customers can trigger manual replays.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline for structuring a full design answer
- [Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — strong reliability framing for overload and failure isolation
