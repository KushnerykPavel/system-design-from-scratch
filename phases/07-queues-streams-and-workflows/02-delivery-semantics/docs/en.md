# At-Most-Once, At-Least-Once, and Exactly-Once Claims

> Delivery guarantees matter only when you name the side effect boundary they actually cover.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Evaluate delivery guarantees honestly and design the deduplication, ack, and state boundaries that make them safe in practice.
**Prerequisites:** `07-queues-streams-and-workflows/01-queues-vs-streams`, `04-apis-contracts-and-schema-evolution/02-idempotency-keys`
**Estimated time:** ~75 min
**Primary artifact:** semantics review checklist + delivery planner

## The Problem

Messaging interviews are full of vague promises like "we will make it exactly once." Strong candidates slow that down and ask:

- exactly once between which components?
- does the broker guarantee it, or does the application logic?
- what happens if the consumer commits state but crashes before acknowledging?
- what happens if the producer retries after an uncertain publish?

This lesson forces you to translate marketing language into explicit failure windows and compensating design choices.

## Clarify

- What is the protected side effect: database write, email send, payment capture, or cache update?
- Can the consumer operation be made idempotent with a durable key?
- What should happen on uncertain delivery: retry, drop, or surface manual reconciliation?
- Is the primary risk duplication, loss, reordering, or user-visible delay?

## Requirements

### Functional

- Choose a realistic delivery model for the business operation.
- Explain producer retries, consumer acknowledgements, and duplicate suppression.
- State when the system can tolerate loss and when it cannot.

### Non-functional

- Keep throughput high without hiding failure windows.
- Bound storage and latency overhead from deduplication.
- Make correctness auditable when operators investigate duplication or loss.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Publish rate | 250K messages/s | drives broker throughput and duplicate pressure during incidents |
| Retry burst | 10x during downstream faults | exposes why semantics become harder exactly when failures happen |
| Dedup retention | 24 hours | sets state-store size for idempotent consumers |
| Consumer groups | 6 critical, 20 best-effort | different paths may need different guarantees |
| Rough cost | broker IO + dedupe store + reconciliation tooling | keeps guarantee choice grounded |

## Architecture

Practical default:

1. Producers retry publish with stable message IDs when they cannot confirm success.
2. Brokers provide at-least-once delivery to durable consumers.
3. Consumers make side effects idempotent with a dedupe key and durable commit record.
4. Ack happens only after the side effect boundary is durable.

```text
producer -> broker -> consumer -> side effect store
                     -> dedupe record
                     -> ack
```

Exactly-once claims are only credible if the side effect boundary and acknowledgement order are explained.

## Data Model & APIs

Useful fields:

- `message_id`
- `producer_id`
- `idempotency_key`
- `attempt`
- `effect_status`
- `processed_at`

Useful APIs:

- `Publish(message)`
- `Handle(message_id, payload)`
- `MarkProcessed(message_id, effect_ref)`
- `Ack(message_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| consumer writes side effect then crashes before ack | duplicate deliveries with same message ID | idempotent handler and durable processed record |
| producer times out after successful publish and retries | publish duplicates increase during network faults | stable producer message IDs and broker-side dedupe where supported |
| dedupe retention expires too early | old retries create second side effect | align dedupe TTL to realistic retry and replay windows |
| team claims exactly-once but only within broker transaction | duplicate emails or charges still occur | define application-side boundary and add reconciliation |

## Observability

- metric: duplicate delivery rate and duplicate suppression hit rate
- metric: ack latency and oldest unacked message age
- metric: processed-record conflicts or idempotency mismatches
- log: message ID, side effect ref, and ack outcome for sampled failures
- trace: publish to side effect to ack timeline
- SLO: no lost critical messages and duplicate side effects below explicit threshold

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| at-least-once plus idempotent consumer | realistic and robust under failure | dedupe state and handler complexity | pretending broker semantics alone solve side effects |
| at-most-once for low-value telemetry | lower latency and less state | possible loss | over-engineering every event path equally |
| exactly-once within narrow system boundary | simplifies some pipelines | higher coordination cost and still partial scope | claiming universal exactly-once |

## Interview It

**Google framing:** "Design asynchronous processing for sending receipts and updating an order ledger." The signal is whether you separate durable ledger writes from non-idempotent external effects like email.

**Cloudflare framing:** "Design event delivery for security signals where replay is acceptable but silent loss is not." The signal is whether you define acceptable duplication and auditability clearly.

**Follow-ups:**
1. What changes if one downstream effect is idempotent and another is not?
2. What if the broker supports transactions but the email provider does not?
3. What if compliance requires proving whether a message was ever processed?
4. What if replaying the topic from last week is now part of the incident response plan?
5. Where do you intentionally accept at-most-once semantics?

## Ship It

- `outputs/semantics-review-checklist.md`
- `outputs/interview-card-delivery-semantics.md`

## Exercises

1. **Easy** — Explain why "ack after processing" is not enough by itself.
2. **Medium** — Design delivery semantics for email, analytics, and billing events in the same product.
3. **Hard** — Redesign a system that marketed itself as exactly-once but still produced duplicate charges.

## Further Reading

- [Kafka semantics](https://kafka.apache.org/documentation/#semantics) — useful for narrowing what broker-level guarantees actually mean
- [Life beyond distributed transactions](https://queue.acm.org/detail.cfm?id=3025012) — strong framing for application-level correctness boundaries
