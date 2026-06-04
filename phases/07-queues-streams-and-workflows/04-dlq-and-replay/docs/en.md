# Dead-Letter Queues and Replay

> A DLQ is not a graveyard. It is a controlled pause button that still needs ownership, triage, and safe re-entry.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Design dead-letter and replay flows that preserve debuggability and recovery without turning failed messages into a forgotten data pile.
**Prerequisites:** `07-queues-streams-and-workflows/02-delivery-semantics`, `07-queues-streams-and-workflows/03-consumer-groups`
**Estimated time:** ~60 min
**Primary artifact:** replay readiness checklist

## The Problem

Teams often add a dead-letter queue to "handle failures" without deciding:

- which failures deserve retry versus quarantine
- who owns triage
- how replay avoids repeating the same bug at scale
- whether replay preserves ordering and downstream safety

This lesson teaches that DLQ design is more about operational discipline than about adding another topic.

## Clarify

- What kinds of failures produce dead letters: malformed payload, transient dependency error, poison message, or schema drift?
- Is replay automatic, operator-driven, or product-triggered?
- Can replay reorder events in a way that breaks business invariants?
- How long must failed messages be retained for audit and recovery?

## Requirements

### Functional

- Capture unprocessable messages with context about why they failed.
- Support investigation, correction, and controlled replay.
- Prevent poison messages from blocking the entire pipeline indefinitely.

### Non-functional

- Keep DLQ growth visible so failure does not become silent backlog.
- Avoid replay storms that overload downstream systems.
- Preserve enough metadata to debug root cause quickly.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Normal failure rate | 0.01% of 200M messages/day | even small percentages create many dead letters at scale |
| Incident failure rate | 5% for 30 minutes | replay volume can dwarf normal traffic after recovery |
| Retention | 14 to 30 days | affects cost and auditability |
| Replay speed | 1x to 5x live traffic | determines need for throttling and isolation |
| Rough cost | extra storage + tooling + operator time | DLQ is an operational product, not a free add-on |

## Architecture

```text
consumer
  -> retry policy
  -> dead-letter topic / queue
  -> triage tooling
  -> controlled replay path
```

A good design separates:

- transient retry
- poison-message quarantine
- replay with guardrails

Do not use the DLQ as the default retry mechanism for ordinary transient faults.

## Data Model & APIs

Useful dead-letter record fields:

- original topic or queue
- partition and offset or task ID
- failure class
- stack or error code
- first failed at / last failed at
- replay attempt count

Useful APIs:

- `MoveToDLQ(record, reason)`
- `ListFailed(filter)`
- `Replay(batch, rate_limit)`
- `Drop(record_id, justification)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| DLQ fills silently for hours | DLQ age and depth alerts fire late or not at all | alert on oldest item age and growth slope |
| poison messages are replayed blindly | replay immediately fails again at scale | require classification or patch validation before replay |
| replay overwhelms downstream systems | consumer lag and error rate spike during recovery | throttle replay and isolate it from live traffic |
| replay breaks ordering-sensitive consumers | state divergence appears later | replay per key or partition with ordering-aware tooling |

## Observability

- metric: DLQ ingress rate by failure class
- metric: oldest DLQ age and total retained items
- metric: replay success rate, replay throughput, and repeat-failure rate
- log: original message metadata plus failure classification
- trace: one failed message from first error through replay outcome
- SLO: critical pipelines detect and surface dead letters within minutes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| structured failure classification | better triage and replay safety | more producer or consumer metadata work | one generic "failed" bucket |
| operator-controlled replay | reduces repeat incidents | slower recovery | automatic blind replay |
| separate live and replay capacity | protects serving path | extra resource overhead | replay on the same path with no rate control |

## Interview It

**Google framing:** "Design failure handling for an order event pipeline." The signal is whether you distinguish transient retries from poison-message quarantine and think about replay safety.

**Cloudflare framing:** "Design quarantine and replay for security telemetry consumers." The signal is whether you protect live analysis while preserving audit and recovery.

**Follow-ups:**
1. What if schema evolution caused the failures and all old messages now look malformed?
2. What if replay must preserve per-customer ordering?
3. What if the DLQ itself becomes too large to inspect manually?
4. What if replaying old events would send duplicate notifications?
5. What metric would wake you up before the product team notices missing downstream state?

## Ship It

- `outputs/replay-readiness-checklist.md`

## Exercises

1. **Easy** — Define what metadata must be stored with a dead-lettered record.
2. **Medium** — Explain when you should retry inline versus send to DLQ.
3. **Hard** — Redesign replay for a pipeline where per-account ordering must be preserved.

## Further Reading

- [Enterprise Integration Patterns: Dead Letter Channel](https://www.enterpriseintegrationpatterns.com/patterns/messaging/DeadLetterChannel.html) — classic framing for quarantine behavior
- [Building event-driven systems](https://www.confluent.io/event-driven-systems/) — useful discussion of replay and stream recovery trade-offs
