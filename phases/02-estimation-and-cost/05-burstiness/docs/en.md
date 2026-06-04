# Peak Factors, Burstiness, and Queue Build-Up

> The average request rate tells you what the system does all day. Burst math tells you whether it survives the next five minutes.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Model arrival bursts, service capacity, and backlog growth so overload discussions stay concrete.  
**Prerequisites:** `02-estimation-and-cost/01-qps-and-request-mix`, `07-queues-streams-and-workflows/01-queues-vs-streams`  
**Estimated time:** ~75 min  
**Primary artifact:** backlog worksheet + failure checklist  

## The Problem

Many candidates size systems using average rate and miss the real question: what happens when arrivals temporarily exceed service capacity? Burstiness decides queue growth, latency blowups, retry storms, and when backpressure becomes mandatory.

This lesson turns burst factors into backlog and recovery time.

## Clarify

- Is the burst user-driven, cron-driven, or retry-driven?
- How long does the burst last: seconds, minutes, or hours?
- Can the system buffer the extra work, or must it reject immediately?
- Is processing parallelizable, or are there ordered partitions that cap throughput?

## Requirements

### Functional

- Estimate backlog created during an overload window.
- Estimate drain time once the burst ends.
- Separate stable average load from unstable peak arrival.

### Non-functional

- Make queue growth visible before architecture choices.
- Expose when load shedding is better than buffering.
- Keep the math simple enough for live use.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Normal arrival rate | 25K msg/s | steady-state load |
| Burst arrival rate | 60K msg/s for 10 min | overload window |
| Service rate | 40K msg/s | maximum sustainable processing |
| Backlog growth | 20K msg/s | determines queue depth |
| Drain time | 20 min after burst | affects latency and recovery |

## Architecture

Compute:

1. Arrival rate during burst.
2. Service rate during burst.
3. Backlog growth = arrival - service.
4. Backlog size = growth x burst duration.
5. Drain time = backlog / post-burst spare capacity.

Example:

- 60K arrivals/s
- 40K processed/s
- 20K/s backlog growth
- over 10 minutes, backlog grows by 12M items
- if post-burst spare capacity is 10K/s, recovery takes about 20 minutes

## Data Model & APIs

The code artifact models:

```text
QueueModel {
  ArrivalRate
  ServiceRate
  BurstSeconds
  RecoveryServiceRate
}
```

Outputs:

- backlog growth rate
- backlog items
- drain seconds

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| average rate used as design limit | queues explode during peak events | size for burst windows, not just daily mean |
| ordered partition caps throughput | one partition backlogs while fleet looks fine | model per-partition throughput |
| retries amplify burst | backlog grows faster than original demand | add retry budgets and client backoff |
| queue absorbs too much | long stale work becomes useless | drop, expire, or degrade low-value work |

## Observability

- metric: queue depth and age percentiles
- metric: arrival rate vs service rate
- metric: rejected or expired work during overload
- metric: recovery time after burst
- SLO: backlog should drain within the agreed recovery window after expected bursts

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| buffer short bursts | preserves availability | higher tail latency and queue memory | immediate rejection |
| shed excess load | protects core path | drops work | pretending all bursts are bufferable |
| overprovision for peak | simplest mental model | high idle cost | elastic or degraded-mode handling |

## Interview It

**Google framing:** "Design a job ingestion pipeline for analytics events." The signal is whether you reason about queue age and recovery, not just steady-state throughput.

**Cloudflare framing:** "Design a DDoS-protected pipeline for edge logs." The signal is whether you distinguish absorbable bursts from traffic that must be dropped or sampled.

**Follow-ups:**
1. What if the burst lasts longer than the queue retention window?
2. What if only one shard is hot?
3. What if retries double effective arrival rate?
4. What if downstream work loses value after 5 minutes?

## Ship It

- `outputs/backlog-worksheet-burstiness.md`
- `outputs/failure-checklist-burstiness.md`

## Exercises

1. **Easy** — Compute backlog for a 2-minute burst with 20% spare capacity.  
2. **Medium** — Compare buffering vs dropping for user notifications.  
3. **Hard** — Rework the estimate when one hot partition limits throughput.  

## Further Reading

- [Google SRE book](https://sre.google/books/) — strong grounding on overload and queueing trade-offs  
- [System design notes](https://github.com/liquidslr/system-design-notes) — useful interview context  
