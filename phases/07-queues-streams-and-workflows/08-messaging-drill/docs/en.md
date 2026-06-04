# Messaging Architecture Drill

> The strongest messaging answer sounds like an ownership and failure model, not a broker shopping list.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Practice compressing primitive choice, delivery semantics, replay, ordering, and overload behavior into one coherent interview answer.
**Prerequisites:** `01-queues-vs-streams`, `02-delivery-semantics`, `03-consumer-groups`, `04-dlq-and-replay`, `06-outbox-and-cdc`, `07-pipeline-backpressure`
**Estimated time:** ~60 min
**Primary artifact:** drill worksheet + scoring rubric

## The Problem

This drill closes the phase. The learner gets an ambiguous system prompt and must produce a messaging design that sounds senior:

- choose the right primitive
- name the ordering scope
- define delivery semantics honestly
- explain replay and failure recovery
- show how overload becomes visible and controlled

The goal is not to mention every pattern. The goal is to make a small number of explicit, defensible choices.

## Clarify

- Is the system dispatching work, broadcasting state changes, or orchestrating long-running progress?
- Which consumer needs what ordering and replay behavior?
- What side effect boundary must be protected from duplicates or loss?
- What backlog or lag would become a product incident?

## Requirements

### Functional

- Pick an async primitive and justify it with ownership and replay needs.
- Define delivery semantics for the critical path.
- Handle poison messages, replay, and backpressure credibly.

### Non-functional

- Keep user-visible completion delay and consumer lag visible.
- Bound operational complexity rather than spraying patterns everywhere.
- Explain trade-offs between correctness, throughput, and simplicity.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak events | 400K events/s | enough to force partitioning and lag reasoning |
| Consumer groups | 5 critical, 12 best-effort | fan-out changes primitive choice |
| Replay window | 7 days | turns retention and event IDs into design constraints |
| Hot-key skew | top 1% keys cause 30% of traffic | exposes partitioning limits |
| Rough cost | broker storage + consumers + workflow state + ops tooling | makes trade-offs practical |

## Architecture

Recommended drill sequence:

1. Clarify whether this is queue, stream, workflow, or a combination.
2. Pick the ownership and ordering boundary.
3. State delivery semantics and dedupe plan for the critical path.
4. Add DLQ or replay only where it changes outcomes.
5. Close with backpressure and observability.

The answer should sound like a control story for work and state, not a product list.

## Data Model & APIs

A strong drill answer usually names:

- message or event ID
- ordering key
- consumer group or workflow status boundary
- replay or dead-letter handle

Example surfaces:

```text
Publish(event_id, key, payload)
Commit(group, partition, offset)
Replay(batch_id, rate_limit)
GetStatus(workflow_id)
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| wrong primitive chosen for the ownership model | compensating tools appear everywhere | restate queue vs stream vs workflow boundary |
| duplicates repeat external side effects | duplicate IDs and reconciliation incidents rise | idempotent handlers and stable event IDs |
| hot partition throttles the whole design | one partition lag diverges sharply | key redesign, isolation, or partition growth plan |
| backlog hides an incident until product impact | oldest-item age and completion delay breach target | alert on age and apply backpressure or shedding |

## Observability

- metric: producer rate, consumer lag, and backlog age
- metric: duplicate suppression hits and DLQ ingress rate
- metric: replay throughput and replay repeat-failure rate
- log: message ID, key, failure class, and replay decisions
- trace: one user action across publish, consume, and completion
- SLO: user-visible completion time plus critical consumer lag target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one clear ownership boundary | coherent design and easier reasoning | less room for vague breadth | mixing queue, stream, and workflow without purpose |
| explicit duplicate policy | safer retries and replay | more metadata and state | hand-waving exactly-once |
| bounded backlog policy | early overload visibility | may require rejection or degradation | infinite queue optimism |

## Interview It

**Google framing:** "Design asynchronous processing for order updates, billing, analytics, and user notifications." The signal is whether you separate command execution from event sharing and define correctness boundaries.

**Cloudflare framing:** "Design global processing for policy updates plus security telemetry." The signal is whether you reason about fan-out, lag, replay, and stateful rollout steps.

**Follow-ups:**
1. What if the interviewer changes the prompt from one consumer to ten independent consumers?
2. What if one external side effect cannot be retried safely?
3. What if replay traffic must coexist with live security events?
4. What if product suddenly needs visible progress on a multi-hour path?
5. Which deep dive would you pick if asked to go one level deeper?

## Ship It

- `outputs/drill-worksheet-messaging.md`
- `outputs/scoring-rubric-messaging.md`

## Exercises

1. **Easy** — Run the drill for image-processing jobs plus analytics side consumers.
2. **Medium** — Run the drill for a user-facing provisioning flow with callbacks and retries.
3. **Hard** — Run the drill for a mixed control-plane plus telemetry system where replay and hot-key skew both matter.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline for structuring the interview answer before deep messaging nuance
- [Kafka: The Definitive Guide](https://www.confluent.io/resources/kafka-the-definitive-guide/) — broader background for partitions, groups, and replay mental models
