# Schema Evolution in Event-Driven Systems

> Producers change once; consumers upgrade on many different clocks.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Evolve event schemas without breaking consumers, replay pipelines, or long-lived data contracts across loosely coupled systems.  
**Prerequisites:** `04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events`, `04-apis-contracts-and-schema-evolution/04-api-versioning`  
**Estimated time:** ~75 min  
**Primary artifact:** schema checklist + migration playbook  

## The Problem

Event systems make fanout easy, but they make compatibility harder because consumers do not all upgrade together. Some consumers are real time, some replay old history, and some are forgotten until they break.

This lesson focuses on safe evolution patterns:

- additive fields
- default handling
- explicit schema registration
- dual-publish or upcasters when meaning changes

## Clarify

- Are events append-only facts or mutable snapshots?
- How many consumer teams exist, and do they all upgrade quickly?
- Will old events be replayed months later?
- Is schema validation enforced centrally or only by convention?

## Requirements

### Functional

- Allow producers to add useful fields safely.
- Keep old consumers working during migration.
- Support replay and backfill without ambiguous event meaning.

### Non-functional

- Prevent silent consumer breakage.
- Keep event contracts discoverable and reviewable.
- Bound the operational cost of multi-version consumption.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Event throughput | 250K events/s | schema validation overhead must stay cheap |
| Consumer count | 20-100 services | many upgrade cadences exist at once |
| Replay window | 180 days | old schema shapes remain operationally relevant |
| Schema changes | weekly additive, rare breaking | suggests compatibility-first process |
| Rough cost | registry + validation + dual-publish during migrations | event evolution has real control-plane cost |

## Architecture

Recommended pattern:

1. Register schemas and review changes.
2. Prefer additive evolution with defaults.
3. Keep old fields until consumer migration completes.
4. Use an upcaster, translator, or new event type for semantic breaks.

```text
producer -> schema validation -> event bus -> consumers
                           -> registry / compatibility checks
```

The core operational insight: replay must be part of the design. If a new consumer cannot read six-month-old events, the contract is weaker than it looks.

## Data Model & APIs

Event contract concerns:

- unique event type and version strategy
- required vs optional fields
- field deprecation policy
- semantic meaning over time

Example:

```text
order.created.v1
order.created.v2
```

Sometimes a new event name is cleaner than stretching one event shape across incompatible meanings.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| producer adds required field unexpectedly | consumer deserialization errors spike | enforce compatibility checks before publish |
| event meaning changes in place | replay outputs diverge from original behavior | create new version or event type |
| abandoned consumer breaks silently | dead-letter spikes after deploy | consumer inventory and staged rollout |
| backfill publishes mixed schema meanings | analytics inconsistencies | explicit migration window and version tagging |

## Observability

- metric: consumer deserialization failure rate by schema version
- metric: dead-letter count by event type and version
- metric: producer publish count by version
- log: schema compatibility check results at deploy time
- trace: event version through producer, bus, and consumer stages
- SLO: supported consumers continue processing new events without breakage during additive evolution

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| additive schema evolution | safer for slow consumers | old fields linger longer | frequent breaking rewrites |
| registry-enforced compatibility | catches mistakes early | more process and tooling | conventions only |
| new event type for semantic change | clearer meaning | more migration work | mutating semantics in place |

## Interview It

**Google framing:** "Design an event contract strategy for an order pipeline with many consumers." The signal is whether you think about replay and consumer lag, not just protobuf syntax.

**Cloudflare framing:** "Evolve edge telemetry events without breaking downstream analytics and security pipelines." The signal is whether you treat long-tail consumers and backfill as first-class constraints.

**Follow-ups:**
1. What if one consumer upgrades quarterly while others deploy daily?
2. What if a field must become required in the future?
3. What if old events must be replayed into a brand new consumer?
4. When is dual-publish worth the cost?
5. How do you migrate semantics, not just schema shape?

## Ship It

- `outputs/schema-checklist-event-schema-evolution.md`
- `outputs/migration-playbook-event-schema-evolution.md`

## Exercises

1. **Easy** — Decide whether four example event changes are additive or breaking.  
2. **Medium** — Write a migration plan for renaming a field used by many consumers.  
3. **Hard** — Redesign an event family where the same event name is used for two different business meanings.  

## Further Reading

- [CloudEvents specification](https://cloudevents.io/) — useful for consistent event envelope thinking  
- [Protocol Buffers compatibility guidance](https://protobuf.dev/programming-guides/proto3/) — practical field-evolution rules  
