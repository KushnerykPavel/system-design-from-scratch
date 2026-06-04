# Backpressure Across Async Boundaries

> A queue is not a pressure release valve unless the producer can feel the pressure too.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design async pipelines where producers, brokers, and consumers coordinate around bounded backlog instead of treating queues as infinite shock absorbers.
**Prerequisites:** `07-queues-streams-and-workflows/01-queues-vs-streams`, `07-queues-streams-and-workflows/07-pipeline-backpressure`, `10-reliability-retries-and-backpressure/04-load-shedding`
**Estimated time:** ~75 min
**Primary artifact:** pipeline pressure planner + observability checklist

## The Problem

Async systems often look safe because producers are decoupled from consumers. Then one consumer tier slows down and:

- backlog grows silently
- retries re-enqueue more work
- old messages become useless but still expensive
- operators only notice when the queue is huge

Backpressure is the missing design contract. The question is not only where work waits. It is who slows down, who gets rejected, and when work stops being worth doing.

## Clarify

- Is the workload lossless, lossy, or priority-based?
- How stale can queued work become before it loses value?
- Can producers be slowed or rejected, or do they only know how to enqueue?
- Is the real bottleneck the broker, consumers, or a downstream dependency behind consumers?

If the interviewer is vague, assume a high-volume async pipeline where producer spikes are common, some messages expire in minutes, and consumer capacity is the main bottleneck.

## Requirements

### Functional

- Keep backlog bounded and visible.
- Slow or reject producers when consumers fall behind.
- Distinguish durable high-value work from droppable low-value work.

### Non-functional

- Avoid memory and storage blowups in the broker.
- Prevent stale work from crowding out fresh critical work.
- Make bottlenecks attributable to a specific stage.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Producer ingress | 400K events/s peak | enough to bury consumers quickly |
| Consumer capacity | 280K events/s steady | backlog grows unless pressure reaches producers |
| Message useful lifetime | 2-15 minutes depending on class | determines whether queue age matters more than depth |
| Retry inflation | 1.2-1.5x under failures | async retries can quietly multiply work |
| Rough cost | broker storage + replay + delayed value | old backlog is often negative value, not just delayed work |

## Architecture

A mature async answer includes pressure signals across boundaries:

```text
producer
  -> enqueue gate / credits
  -> broker with bounded partitions
  -> consumer pool
  -> downstream dependency
  -> lag and queue-age feedback to producers
```

Useful patterns:

- quota or credit-based producer admission
- priority queues or separate topics
- message TTL and stale-drop policy
- consumer auto-scaling tied to useful-lag, not only depth

## Data Model & APIs

Useful metadata:

- `message_class`
- `enqueue_time`
- `expires_at`
- `delivery_attempt`
- `priority`

Useful interfaces:

- `Publish(event, class)`
- `CanAccept(class, pressure_snapshot)`
- `Ack / Nack / DropExpired`

Senior-level detail:

- backlog age usually tells you more than raw count
- low-priority topics may need independent drop policies
- retries should re-enter the same pressure accounting, not bypass it

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| broker stores unbounded backlog of low-value work | queue age and expired-message share rise steadily | TTL, stale-drop policies, and class-based quotas |
| producers never see consumer distress | backlog grows until storage or latency collapses | feedback channel, credits, or publish rejection |
| retry workers amplify a slow downstream | attempt rate rises while completion rate stays flat | retry budget, exponential backoff, and dependency-aware pause |
| one noisy topic starves all others | partition lag skew and per-class queue age diverge sharply | isolation by topic, quota, or dedicated consumer pools |

## Observability

- metric: queue depth and queue age by topic and priority class
- metric: ingress rate, consume rate, and backlog growth delta
- metric: expired or dropped message ratio
- metric: downstream latency seen by consumers and retry inflation
- log: publish rejects, stale drops, and pressure transitions
- trace: enqueue to completion latency for sampled messages
- SLO: high-priority messages should stay within freshness target even during bursty producer spikes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| bounded queues with publish rejection | protects broker and consumers | some work is lost or deferred | infinite buffering |
| queue-age-based dropping | preserves fresh value | older work may never execute | FIFO loyalty for useless stale work |
| separate classes or topics | better fairness and isolation | more operational surface area | one shared backlog for everything |

## Interview It

**Google framing:** "Design an async processing pipeline that survives spikes and slow consumers." The signal is whether pressure propagates upstream instead of hiding in the broker.

**Cloudflare framing:** "Design a high-volume event pipeline where edge nodes enqueue more work than central processors can always absorb." The signal is whether you separate droppable versus must-process traffic.

**Follow-ups:**
1. What if backlog keeps growing but CPU is low?
2. Which messages can be dropped, and when?
3. How do producers learn that consumers are overwhelmed?
4. What if the downstream database is slow only for one message class?
5. How does the design change at 10x event volume?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/observability-checklist-async-backpressure.md`
- `outputs/tradeoff-matrix-async-backpressure.md`

## Exercises

1. **Easy** — Pick one metric that is stronger than queue depth for deciding whether a backlog is dangerous.
2. **Medium** — Explain how you would propagate pressure back to producers in a multi-tenant queueing system.
3. **Hard** — Redesign the pipeline when some traffic is legally required to be retained while other traffic can be dropped aggressively.

## Further Reading

- [The Log: What every software engineer should know about real-time data's unifying abstraction](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying) — useful mental model for queue and stream pipelines
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview framing for queues and workload trade-offs
