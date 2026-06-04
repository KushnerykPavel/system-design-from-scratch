# Backpressure in Event Pipelines

> A queue absorbs bursts only until it becomes proof that the system is losing the race.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Design backpressure, admission control, and lag-management behavior for asynchronous pipelines before backlog becomes user-visible failure.
**Prerequisites:** `02-estimation-and-cost/05-burstiness`, `10-reliability-retries-and-backpressure/05-async-backpressure`
**Estimated time:** ~60 min
**Primary artifact:** backpressure response checklist

## The Problem

Async systems hide overload better than synchronous ones. That is useful for smoothing bursts, but dangerous because the pipeline can look "available" while lag quietly explodes.

This lesson focuses on recognizing when buffering stops helping and starts masking an incident.

## Clarify

- Is the workload loss-tolerant, delay-tolerant, or neither?
- Which backlog matters most: queue depth, oldest event age, or user-visible completion delay?
- Can producers slow down, shed, batch, or degrade when downstream falls behind?
- Which consumers are critical and which can fall behind safely?

## Requirements

### Functional

- Smooth short bursts without dropping critical work unnecessarily.
- Surface overload before user-visible latency or correctness degrades too far.
- Apply differentiated policies for critical versus best-effort traffic.

### Non-functional

- Protect downstream systems from unbounded retry or replay pressure.
- Keep backlog recovery operationally credible after incidents.
- Avoid turning one slow consumer into fleet-wide collapse.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Steady ingest | 120K events/s | sets baseline throughput |
| Peak burst | 8x for 5 minutes | enough to stress buffers and consumers |
| Consumer drain rate | 150K events/s healthy, 40K degraded | backlog math changes dramatically during faults |
| Critical traffic share | 20% of events | drives priority handling |
| Rough cost | queue storage + extra consumers + shed work | buffering strategy is an economic choice too |

## Architecture

Typical controls:

- bounded queues or partitions
- producer throttling or admission control
- priority lanes for critical work
- consumer autoscaling with lag-aware triggers
- load shedding for replay or low-value traffic

```text
producers -> admission control -> broker -> consumers -> downstreams
                     ^                |
                     |                -> lag signals
                     -> feedback loop
```

The key is to define what happens before the queue is effectively infinite.

## Data Model & APIs

Useful control-plane knobs:

- per-topic retention and max backlog age
- producer quotas
- priority classes
- replay rate limits

Useful APIs:

- `Publish(class, event)`
- `Throttle(producer_id, limit)`
- `DescribeLag(topic, group)`
- `PauseReplay(job_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| queue depth grows while no one pages | oldest-item age rises silently | alert on age and user-visible delay, not just count |
| autoscaling lags behind burst | backlog grows faster than new consumers help | pre-provision headroom and use rate limits upstream |
| one best-effort consumer starves critical traffic | critical completion SLO misses | isolate traffic classes or enforce priority lanes |
| replay after incident prevents recovery | live lag worsens during catch-up | cap replay rate and protect live path first |

## Observability

- metric: queue depth and oldest unprocessed age
- metric: ingest rate versus drain rate
- metric: publish throttling, shedding, and rejected work by class
- log: producer throttles and consumer pauses with reasons
- trace: one delayed item from publish to completion
- SLO: completion delay for critical traffic and maximum acceptable backlog age

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| bounded backlog age | overload becomes visible early | may require shedding or rejecting work | pretending deep queues mean safety |
| priority lanes | protects important work | less fair for best-effort traffic | one undifferentiated pipe |
| producer feedback loop | slows overload propagation | more coupling and control-plane work | letting producers publish blindly |

## Interview It

**Google framing:** "Design async processing for notifications where some events are critical and some are promotional." The signal is whether you prioritize by business value and define overload behavior.

**Cloudflare framing:** "Design protection for security or telemetry pipelines during sudden traffic spikes." The signal is whether you reason about lag age, replay pressure, and live-path protection.

**Follow-ups:**
1. What if consumers are already at max scale and lag is still growing?
2. What if producers cannot slow down because they are edge data planes?
3. What if replay traffic competes with fresh security events?
4. What if dropping data is forbidden but delay beyond 15 minutes is also unacceptable?
5. Which signal would you alert on first: depth, age, or user-visible delay?

## Ship It

- `outputs/backpressure-response-checklist.md`

## Exercises

1. **Easy** — Explain why queue depth alone can mislead.
2. **Medium** — Design separate overload policies for critical and best-effort events.
3. **Hard** — Redesign a pipeline where replay traffic caused a second incident after recovery.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/pub40801/) — strong mental model for percentile sensitivity and overload propagation
- [Google SRE book: Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — practical backpressure and protection concepts
