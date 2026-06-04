# Outbox and CDC Patterns

> The safest event is usually the one derived from a committed write, not the one published from a wishful code path.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design event publication around transactional truth so state changes and emitted events stay aligned under failure.
**Prerequisites:** `07-queues-streams-and-workflows/02-delivery-semantics`, `08-consistency-replication-and-transactions/01-consistency-spectrum`
**Estimated time:** ~75 min
**Primary artifact:** outbox decision guide + relay validator

## The Problem

A classic failure looks like this:

1. service writes the database row
2. service crashes before publishing the event
3. downstream systems never learn the change happened

Or the reverse:

1. service publishes the event
2. database transaction rolls back
3. downstream systems act on a state change that never committed

Outbox and CDC patterns exist to align durable state changes with emitted events without relying on distributed transactions across every dependency.

## Clarify

- What is the source of truth: relational DB, document store, event log, or something else?
- Does the service own the database transaction where the business change happens?
- Is ordering required per aggregate, tenant, or table row?
- What is the tolerated delay between commit and event visibility?

## Requirements

### Functional

- Ensure committed state changes eventually produce corresponding events.
- Preserve event ordering at the required boundary.
- Allow replay or recovery if the relay pipeline fails.

### Non-functional

- Avoid dual-write inconsistency under crashes and retries.
- Keep relay lag visible and bounded.
- Support schema evolution without breaking downstream consumers.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Transaction rate | 40K writes/s | outbox volume grows with primary writes |
| Event fan-out | 4 downstream systems | publication failures have broad impact |
| Relay lag target | under 5 seconds | determines polling or CDC tuning |
| Retention | 7 days hot, archive after | affects replay and audit cost |
| Rough cost | primary DB writes + relay service + broker storage | outbox safety is not free |

## Architecture

Two common patterns:

- **Transactional outbox**: application writes business row and outbox row in one DB transaction, then a relay publishes the outbox record.
- **CDC**: a log reader captures committed DB changes and emits events from the database change stream.

```text
API -> DB transaction
        -> business row
        -> outbox row
relay -> broker -> downstream consumers
```

Pick the variant that matches control over the schema, ordering needs, and operational maturity.

## Data Model & APIs

Outbox record fields:

- `event_id`
- `aggregate_id`
- `aggregate_version`
- `event_type`
- `payload`
- `created_at`
- `published_at`

Useful APIs:

- `WriteBusinessChange(tx, change)`
- `InsertOutbox(tx, event)`
- `Relay(batch_size)`
- `MarkPublished(event_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| business write commits but publish path crashes | outbox backlog age rises | relay republishes from durable outbox |
| relay publishes twice after uncertain ack | duplicate event IDs seen downstream | idempotent consumers and stable event IDs |
| outbox table grows until it hurts OLTP performance | table scan latency and storage usage rise | index carefully, partition, archive, and trim |
| CDC emits changes without business-friendly envelope | consumers couple to low-level DB details | transform into stable domain events before broad fan-out |

## Observability

- metric: outbox insert rate and unpublished backlog size
- metric: oldest unpublished record age
- metric: relay publish success rate and duplicate publish count
- log: event ID, aggregate ID, and relay outcome
- trace: transaction commit to broker publish latency
- SLO: committed critical changes become visible on the event bus within target lag

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| transactional outbox | explicit app-level event contract and strong alignment with writes | extra table, relay, and cleanup work | direct dual-write from request handler |
| CDC from database log | lower app-code touch and broad capture | more infra complexity and rougher event semantics | polling every table ad hoc |
| stable event IDs and versions | safer replay and dedupe | more envelope discipline | fire-and-forget opaque payloads |

## Interview It

**Google framing:** "Design change propagation from an orders database to billing, analytics, and notifications." The signal is whether you catch the dual-write hazard and design a safe relay path.

**Cloudflare framing:** "Design propagation of control-plane state changes to distributed consumers." The signal is whether you reason about ordering, lag, and replay in addition to raw event publication.

**Follow-ups:**
1. What if the outbox table becomes the hottest table in the database?
2. What if downstream consumers need only domain events, not row-level mutations?
3. What if relay lag breaches the product SLO during peak traffic?
4. What if the source database is managed by another team and you only get CDC access?
5. What if one consumer replays old events while others stay live?

## Ship It

- `outputs/outbox-decision-guide.md`
- `outputs/failure-checklist-outbox-cdc.md`

## Exercises

1. **Easy** — Explain why publishing directly after a DB write is a dual-write risk.
2. **Medium** — Compare transactional outbox and CDC for a service you do not fully control.
3. **Hard** — Redesign an overloaded outbox relay where backlog age is now a product incident.

## Further Reading

- [Transactional outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html) — canonical pattern framing
- [Debezium documentation](https://debezium.io/documentation/) — practical CDC model and operational considerations
