# Order State Machine and Recovery

> The order is not one status field; it is a recovery strategy encoded as a state machine.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design an order orchestration state machine that coordinates payment, inventory, fulfillment, and recovery without hiding ambiguity behind a single "failed" status.
**Prerequisites:** `07-queues-streams-and-workflows/05-workflow-engines`, `18-messaging-and-job-platforms/02-workflow-engine`, `19-payments-wallets-and-ordering/03-inventory-reservation`
**Estimated time:** ~75 min
**Primary artifact:** state-machine validator + recovery runbook

## The Problem

Design the order state machine for checkout, payment authorization, inventory reservation, fulfillment request, cancellation, and refund or compensation. The system must recover from retries, partial success, and downstream ambiguity.

This lesson matters because many system-design answers pretend the order is either "pending" or "done." Senior answers model which subsystem owns each transition, which states are terminal, and how operators recover when callbacks arrive out of order.

## Clarify

- Is the order flow synchronous checkout only, or does it include long-running fulfillment?
- Which steps are reversible, and which require compensating actions instead?
- Is "payment accepted" equivalent to "order accepted," or can inventory still fail afterward?
- What operator actions are allowed when automated recovery stalls?

If the interviewer stays broad, assume checkout creates an order, then inventory reserve and payment authorize happen with retries, fulfillment is asynchronous, and ambiguous partial failures require explicit recovery states instead of silent rollback assumptions.

## Requirements

### Functional

- Track explicit order states and allowed transitions.
- Coordinate payment, inventory, and fulfillment side effects safely.
- Support idempotent transition requests and repeated callbacks.
- Recover from stuck or ambiguous states with retries or operator actions.
- Preserve a full event history for debugging and customer support.

### Non-functional

- Keep the workflow understandable under failure.
- Avoid hidden double-execution of external side effects.
- Make terminal versus recoverable states explicit.
- Let operators inspect and resume stuck orders safely.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Order creations | 90K/s peak | drives workflow state writes and event history volume |
| Average side effects | 4 to 7 per order | determines retry fanout and ambiguity surface |
| Callback delay tail | up to hours for some providers | requires durable intermediate states |
| Stuck-order rate target | <0.01% | operational tooling matters even for rare cases |
| Peak factor | 5x during promotions | recovery queues must scale with happy-path traffic |

## Architecture

```text
checkout API
  -> order workflow service
  -> state machine + event history
  -> payment, inventory, fulfillment integrations
  -> timeout / retry scheduler
  -> operator recovery console
```

Design notes:

1. Keep the workflow state separate from downstream resource truth; the order references payment and reservation objects rather than pretending it owns them.
2. Use explicit intermediate states such as `payment_pending`, `inventory_reserved`, `awaiting_compensation`, or `manual_review`.
3. Treat repeated callbacks as normal, and make transition handling idempotent.
4. Model ambiguity explicitly when a timeout leaves the external outcome unknown.

## Data Model & APIs

Core entities:

```text
order(order_id, customer_id, state, version, created_at)
order_event(event_id, order_id, type, source, payload_ref, created_at)
order_step(order_id, step_name, attempt, status, next_retry_at)
recovery_task(task_id, order_id, reason_code, state, owner)
```

Useful interfaces:

- `CreateOrder(cart_ref, idempotency_key)`
- `AdvanceOrder(order_id, transition, source_ref, idempotency_key)`
- `GetOrderTimeline(order_id)`
- `RetryRecoverableStep(order_id, step_name)`
- `ResolveRecoveryTask(task_id, action)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| callback arrives twice or out of order | duplicate transition attempts and invalid-transition counts | idempotent transition handling with version checks |
| payment succeeded but inventory failed | compensation-state counts and mismatch alerts | explicit compensation flow for release or refund |
| order gets stuck in pending forever | state-age SLOs and retry scheduler backlog | timeout states, escalation queues, and operator tooling |
| hidden side effects happen on retry | repeated external reference use | outbox or dedupe keys per downstream integration |

## Observability

- metric: orders by state, state age, and invalid transition rate
- metric: compensation workflow rate and stuck-order backlog
- metric: callback delay percentiles by provider and step
- log: manual recovery actions, force-cancel, and force-complete events
- trace: checkout -> payment -> inventory -> fulfillment -> terminal state
- SLO: 99.9% of orders reach a valid terminal or recoverable intermediate state within the target workflow window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit workflow states | debuggable recovery and operator clarity | more design effort and state explosion risk | one generic pending/failed status |
| durable event history | excellent replay and debugging | higher storage and indexing cost | only storing latest order row |
| compensation instead of distributed transaction | practical cross-service recovery | temporary inconsistency windows | pretending global rollback is available |

## Interview It

**Google framing:** "Design the order orchestration layer for checkout, payment, and fulfillment." Expect pushback on ambiguous timeouts, compensation, and how operators recover stuck orders.

**Cloudflare framing:** "Design a provisioning order flow for paid infrastructure resources." Expect pressure on long-running workflows, partial provisioning, and control-plane recovery semantics.

**Follow-ups:**
1. What if fulfillment can take days and supports partial shipment?
2. What if payment callbacks can arrive after cancellation?
3. How do you model manual-review states without freezing the entire workflow?
4. What if one downstream service has much weaker reliability than the rest?
5. How would you migrate from cron-based status repair to explicit workflow history?

## Ship It

- `outputs/replay-runbook-order-state-machine.md`

## Exercises

1. **Easy** — Draw a minimal valid order state machine from create to success or cancellation.
2. **Medium** — Add an ambiguity state for payment timeout where final processor outcome is unknown.
3. **Hard** — Redesign the workflow for partial shipment and split payment capture.

## Further Reading

- [Temporal workflow concepts](https://docs.temporal.io/workflows) — helpful for long-running orchestration thinking
- [Saga pattern](https://microservices.io/patterns/data/saga.html) — good grounding for compensation-based workflows
