# Consistency Trade-offs for Serving and Storage

> Google-style consistency discussions are rarely about saying "strong" or "eventual." They are about choosing where correctness is strict, where staleness is acceptable, and how the system behaves when those boundaries are stressed.

**Type:** Build
**Company focus:** Google
**Learning goal:** Practice choosing consistency models deliberately for senior-level prompts that mix serving latency, storage correctness, and multi-region realities.
**Prerequisites:** `08-consistency-replication-and-transactions/01-consistency-spectrum`, `08-consistency-replication-and-transactions/05-transactions`, `13-multi-region-cdn-and-edge-traffic/06-geo-consistency`
**Estimated time:** ~75 min
**Primary artifact:** consistency drill validator + decision checklist

## The Problem

Many Google interview prompts become interesting when the answer must decide where freshness, ordering, and correctness actually matter. A configuration service, document store, ad-serving platform, social graph, and metadata service all make different choices.

This lesson exists to stop consistency talk from becoming vague. You should be able to say what is strongly consistent, what is eventually consistent, what user-visible anomaly is acceptable, and what operational cost follows.

## Clarify

- Which operation must be correct immediately: writes, reads, ordering, uniqueness, or billing?
- Is the dominant path interactive serving, background processing, or both?
- Is the system regional, multi-region, or globally writable?
- Which anomaly would be unacceptable to the product or business?

If the prompt is broad, choose one authoritative write path and one read path, then reason about the weakest acceptable freshness model for the read side.

## Requirements

### Functional

- Define the authoritative write path and source of truth.
- Choose a consistency model for each major read or write flow.
- Explain user-visible anomalies that the design permits or rejects.
- Show how the system degrades under lag, partition, or failover.
- Redesign when stronger consistency or lower latency becomes the new priority.

### Non-functional

- Avoid "global strong consistency everywhere" as a default answer.
- Avoid hiding anomalies behind generic replication language.
- Keep latency, availability, and operational complexity visible.
- Make regional and multi-region trade-offs explicit.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read/write split | 20:1 to 100:1 in many serving systems | often determines whether stale replicas are acceptable |
| Cross-region RTT | 80 to 200 ms | bounds the cost of synchronous global coordination |
| Allowed read staleness | sub-second to minutes | shapes replica and cache design |
| Recovery window | seconds to tens of minutes | affects failover and repair choices |

## Architecture

Start with a consistency map:

```text
write path -> source of truth -> replication path -> read path
```

Then classify:

1. Which writes require synchronous acknowledgment from the authority?
2. Which reads can tolerate stale replicas, caches, or async propagation?
3. Which background repairs, reconciliation jobs, or version checks contain divergence?
4. Which failure mode breaks the product promise first?

A strong answer often separates:

- correctness-sensitive metadata from high-scale serving reads
- write authority from read fanout
- user-visible freshness from back-office reconciliation

## Data Model & APIs

Consistency decisions should attach to concrete records:

```text
record(key, version, region, committed_at)
read_policy(flow, max_staleness, fallback)
write_policy(flow, quorum, idempotency_boundary)
```

Useful interfaces:

- `Write(record)`
- `ReadStrong(key)`
- `ReadBoundedStale(key, max_age)`
- `Failover(region)`
- `Reconcile(version_gap)`

If you cannot name the record whose correctness matters most, the consistency discussion is probably still too abstract.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| stale replicas serve surprising reads | freshness lag and version-skew metrics | bounded-staleness policy, read fencing, or strong-read escape hatch |
| cross-region synchronous writes miss latency target | write latency by region pair | localize authority or narrow strong-consistency scope |
| failover promotes lagging state | replica lag and divergence counters | promotion gates, last-committed markers, replay and repair |
| candidate promises correctness without anomaly model | no explicit stale-read or lost-update story | force description of one tolerated anomaly and one forbidden anomaly |

## Observability

- metric: replica lag and version skew
- metric: write latency by quorum scope and region pair
- metric: stale-read age or bounded-staleness compliance
- metric: reconciliation backlog and repair success rate
- log: failover decisions, divergence detection, and fenced reads
- trace: write commit through replica visibility and client read
- SLO: the chosen consistency guarantee is measurable, not just stated

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| strong write authority for critical metadata | protects correctness-sensitive updates | higher write latency and reduced fault tolerance | multi-master for all writes |
| stale or cached read path for scale | lower read latency and better fanout | freshness anomalies | strong reads everywhere |
| bounded-staleness model | explicit user-visible contract | extra enforcement and monitoring | undefined replica freshness |

## Interview It

**Google framing:** "Design a globally used serving or metadata system." Expect pressure on which flows need stronger correctness, what anomalies are acceptable, and how cross-region coordination affects latency.

**Follow-ups:**
1. What if write correctness becomes more important than latency?
2. What if the product adds globally visible writes?
3. What if caches are returning data older than the product can tolerate?
4. What if a region fails during a version rollout?
5. What changes at 10x reads but only 2x writes?

## Ship It

- `outputs/consistency-checklist-google-serving-storage.md`

## Exercises

1. **Easy** - Give one example of a read that can tolerate bounded staleness.
2. **Medium** - Compare a single-writer regional authority with multi-region write coordination.
3. **Hard** - Redesign a metadata service when clients now require globally visible updates within two seconds.

## Further Reading

- [Spanner paper](https://research.google/pubs/pub39966/) - useful reference for globally consistent data trade-offs
- [Designing Data-Intensive Applications](https://dataintensive.net/) - strong grounding for replication, anomalies, and consistency models
