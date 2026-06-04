# Exactly-Once Myths in Interview Design

> "Exactly once" is usually a boundary definition plus careful side-effect handling, not a magic end-to-end property.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Learn to answer exactly-once questions precisely by defining scope, naming residual risk, and translating the claim into idempotency, dedupe, and transactional boundaries.
**Prerequisites:** `04-apis-contracts-and-schema-evolution/02-idempotency-keys`, `07-queues-streams-and-workflows/02-delivery-semantics`, `08-consistency-replication-and-transactions/06-sagas`
**Estimated time:** ~60 min
**Primary artifact:** exactly-once claim validator + interview card

## The Problem

Interviewers often ask whether a messaging or job platform can deliver exactly once. The real challenge is not to say "yes" or "no" quickly, but to define the boundary: exactly once into what store, under which failures, and whether external side effects are included.

This lesson matters because sloppy exactly-once claims are a strong senior-level negative signal. Strong candidates narrow the claim, explain where duplicates can still surface, and propose practical mitigations such as idempotency keys, dedupe windows, transactional outboxes, or fenced consumers.

## Clarify

- Are we talking about broker delivery, consumer state update, or end-user-visible side effects?
- What failures are in scope: crashes, network ambiguity, retries, replays, or external callback timeouts?
- Can the consumer write to a transactional store with the dedupe key?
- Are irreversible side effects like email, payments, or webhooks part of the claim?

If the interviewer stays broad, assume the safe answer is "exactly-once within a bounded storage or processing boundary, with idempotent side effects beyond that boundary."

## Requirements

### Functional

- Define the exact scope of the guarantee.
- Explain how duplicate delivery attempts are detected or tolerated.
- Show how committed state and message progress stay consistent.
- Address replay and consumer restart behavior.
- Name what is not fully guaranteed.

### Non-functional

- Avoid overclaiming correctness that depends on perfect networks.
- Keep the design understandable enough to defend under follow-up questions.
- Preserve operational visibility into duplicate or ambiguous processing.
- Make the cost of stronger semantics explicit.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Normal duplicate attempt rate | 0.01% to 0.1% | even low rates matter for side effects |
| Replay volume after incidents | 10M to 1B events | dedupe storage and windows must survive burst recovery |
| Dedup retention window | hours to days | short windows lower cost but miss long-delayed retries |
| Cross-system side effects | many downstream APIs | end-to-end guarantees get weaker quickly |
| Peak factor | 100x duplicate pressure during partial outages | the hard case is ambiguity, not the happy path |

## Architecture

```text
durable message
  -> consumer reads
  -> idempotency / dedupe check
  -> transactional state update
  -> ack or commit progress
  -> optional external side effect with fence or outbox
```

Design notes:

1. Define whether the exactly-once boundary ends at the consumer's database write or includes downstream effects.
2. Prefer idempotent consumers even if the broker claims strong semantics.
3. Separate "duplicate delivery happened" from "duplicate user-visible effect happened."
4. Use outbox or transactional write patterns when consumer state and downstream publication must align.

## Data Model & APIs

Core records:

```text
message_id
idempotency_key
consumer_group
processing_state
dedupe_expiry
side_effect_token
commit_cursor
```

Useful interfaces:

- `Process(message_id, payload, idempotency_key)`
- `UpsertDedupeKey(key, state)`
- `CommitCursor(group, partition, offset)`
- `PublishOutboxBatch(batch_id)`
- `RecordSideEffectAttempt(token, status)`

If you cannot say where the dedupe key lives and how it is written atomically, the exactly-once answer is probably still fuzzy.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| consumer crashes after side effect but before ack | duplicate-attempt audit and dedupe-hit rate | idempotent side effect token or outbox |
| dedupe window expires too early | duplicate-after-window count | longer retention or business-level dedupe |
| replay bypasses dedupe checks | replay duplicate spike and downstream anomalies | force replay through normal idempotency path |
| end-to-end claim hides external uncertainty | incident review shows duplicate user action | narrow the guarantee and state residual risk clearly |

## Observability

- metric: dedupe-hit rate and duplicate-attempt rate
- metric: outbox lag or side-effect confirmation lag
- metric: replayed events that were suppressed by idempotency
- metric: ambiguous processing outcomes requiring manual review
- log: dedupe decisions with key, boundary, and result reason
- trace: message read through dedupe, state write, and side-effect handling
- SLO: duplicate user-visible effects stay below the agreed business threshold for the protected operation

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| bounded exactly-once claim | accurate and defensible | less marketing-friendly | vague end-to-end promise |
| longer dedupe retention | catches delayed duplicates | higher storage and lookup cost | tiny best-effort dedupe window |
| outbox or transactional coupling | stronger state-to-event consistency | more implementation complexity | separate best-effort writes |

## Interview It

**Google framing:** "Can your queue provide exactly-once processing?" A strong answer narrows the claim, describes the storage boundary, and explains what remains at-least-once operationally.

**Cloudflare framing:** "How do you prevent duplicate side effects in a shared async platform?" Expect follow-ups on retries, replay, and what happens when downstream systems are outside your trust boundary.

**Follow-ups:**
1. What changes if the consumer sends email or charges a card?
2. How do you handle a replay one week later?
3. When is at-least-once with idempotency the better answer?
4. How would you explain this to an interviewer in two minutes?
5. What if the dedupe store is temporarily unavailable?

## Ship It

- `outputs/interview-card-exactly-once-myths.md`

## Exercises

1. **Easy** — State a precise exactly-once claim for "write order status to a database."
2. **Medium** — Extend that answer to include webhook delivery and explain what weakens.
3. **Hard** — Redesign the consumer when replay storms happen after a week-long outage and the dedupe window is too short.

## Further Reading

- [Kafka exactly-once semantics](https://www.confluent.io/blog/simplified-robust-exactly-one-semantics-in-kafka-2-5/) — useful example of bounded semantics
- [Idempotent Consumer pattern](https://microservices.io/post/microservices/patterns/2020/10/16/idempotent-consumer.html) — practical framing for duplicate tolerance
