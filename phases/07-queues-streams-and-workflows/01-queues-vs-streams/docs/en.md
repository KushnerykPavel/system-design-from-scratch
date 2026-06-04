# Queues vs Streams vs Workflows

> Delivery shape is a product decision with operational consequences, not a vocabulary quiz.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Choose the right asynchronous primitive by matching business needs to fan-out, ordering, replay, state, and human-time semantics.
**Prerequisites:** `04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events`, `05-storage-indexing-and-access-patterns/02-access-pattern-first`
**Estimated time:** ~75 min
**Primary artifact:** messaging decision matrix

## The Problem

Senior candidates often say "use Kafka" or "put it on a queue" before stating what problem the async layer is solving. That sounds fluent, but it hides the important questions:

- is one consumer enough or do many teams need the same event?
- is replay a feature or a hazard?
- does the system need long-running business state?
- does the user need immediate confirmation or eventual progress tracking?

This lesson teaches the primitive-selection step before the deeper messaging lessons in the phase.

## Clarify

- Is this work one-to-one task execution, one-to-many event distribution, or multi-step business orchestration?
- Does the consumer need independent replay from an immutable history?
- Is strict per-key ordering required, or only best-effort ordering?
- Can the producer forget the work after handoff, or must it track long-running state and retries?

## Requirements

### Functional

- Pick an async primitive that matches the ownership and replay model.
- Explain how producers and consumers coordinate progress.
- Separate transient execution from durable business state when needed.

### Non-functional

- Keep latency expectations explicit for enqueue, processing, and completion.
- Bound operational complexity for replay, debugging, and schema evolution.
- Make failure handling visible instead of assuming the primitive solves it automatically.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Producer traffic | 150K events/tasks/s | enough scale that batching, partitions, and fan-out matter |
| Fan-out | 1 to 12 downstream consumers | drives queue versus stream choice |
| Replay window | hours to 30 days | influences retention and immutable log value |
| Long-running jobs | 200K active workflows | exposes when a plain queue is insufficient |
| Rough cost | broker storage + consumer compute + orchestration metadata | keeps primitive choice tied to operations |

## Architecture

Use these defaults:

- **Queue** when one logical consumer group should perform work and ack completion.
- **Stream** when multiple independent consumers need ordered durable history and replay.
- **Workflow engine** when the system must track multi-step progress, timers, compensation, and human-time waits.

```text
producer
  -> queue           for task dispatch
  -> stream          for event distribution and replay
  -> workflow engine for stateful orchestration
```

The primitive is not the whole design. You still need idempotency, ownership boundaries, and failure recovery.

## Data Model & APIs

Typical surfaces:

- Queue: `Enqueue(task)`, `Ack(task_id)`, `Retry(task_id)`
- Stream: `Append(topic, key, event)`, `Consume(group, partition, offset)`, `Commit(offset)`
- Workflow: `StartWorkflow(input)`, `Signal(id, event)`, `GetStatus(id)`, `Compensate(id)`

Core entities:

- task
- event
- workflow instance
- consumer group
- offset or execution checkpoint

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| using a queue where multiple downstream systems need replay | teams build side databases and custom taps | switch to log-based distribution or add durable event stream |
| using a stream for long-running business state | offset progress exists but business status is unclear | introduce workflow state store and explicit timers |
| producer assumes exactly-once from the broker alone | duplicate side effects appear downstream | make handlers idempotent and store dedupe keys where needed |
| orphaned work after consumer crash | lag grows and completion age breaches SLO | visibility timeout, retries, or resumable workflow tasks |

## Observability

- metric: enqueue rate, append rate, and fan-out per topic or queue
- metric: consumer lag, retry rate, and workflow step age
- metric: replay volume and oldest unprocessed item age
- log: task or event IDs with producer and consumer ownership
- trace: one user action across publish, consume, and completion
- SLO: handoff latency plus completion latency for the user-visible path

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| queue for single-owner execution | simple acknowledgement model | poor fit for multi-team replay | using a stream for all background work |
| stream for shared event history | durable fan-out and replay | higher operational complexity and storage cost | bespoke fan-out from the producer |
| workflow engine for stateful orchestration | timers, retries, compensation, and status become explicit | more metadata and platform overhead | encoding business state only in queue retries |

## Interview It

**Google framing:** "Design async processing for order fulfillment plus analytics consumers." The signal is whether you separate command execution from event distribution and resist collapsing both into one primitive.

**Cloudflare framing:** "Design policy updates and security event processing at global edge scale." The signal is whether you reason about fan-out, replay, and operational state instead of defaulting to a single broker pattern.

**Follow-ups:**
1. What changes if ten new teams want to consume the same payload six months later?
2. What if product needs a user-visible progress page for a multi-hour job?
3. What if one consumer must reprocess the last seven days from scratch?
4. What if one downstream system must never observe out-of-order events for the same tenant?
5. When do you intentionally combine queue, stream, and workflow in one design?

## Ship It

- `outputs/messaging-decision-matrix.md`

## Exercises

1. **Easy** — Choose between queue and stream for thumbnail generation plus audit logging.
2. **Medium** — Redesign a "Kafka for everything" proposal into clearer ownership boundaries.
3. **Hard** — Explain where a workflow engine becomes necessary in a cross-region payment dispute process.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — helpful baseline for structuring messaging trade-off discussions
- [Designing Data-Intensive Applications, Chapter 11](https://dataintensive.net/) — strong framing for streams, logs, and event processing
