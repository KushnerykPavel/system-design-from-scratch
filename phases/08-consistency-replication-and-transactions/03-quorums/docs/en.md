# Quorums, Read Repair, and Divergence

> Quorums reduce uncertainty; they do not erase it. You still need to talk about stale reads, conflicting versions, and repair work.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Reason about quorum read and write choices, overlap guarantees, and the operational consequences of divergence repair.
**Prerequisites:** `01-consistency-spectrum`, `02-leader-follower`, `05-storage-indexing-and-access-patterns/02-access-pattern-first`
**Estimated time:** ~75 min
**Primary artifact:** quorum planning sheet

## The Problem

Quorum systems are attractive because they can keep serving through some failures while offering a stronger story than blind eventual consistency. But interview answers often stop at `R + W > N` and miss the operational reality:

- reads may still be stale
- conflicting versions still need reconciliation
- repair traffic can become a hidden cost
- tail latency changes with quorum size

This lesson helps you explain those trade-offs honestly.

## Clarify

- Is availability during replica failure more important than freshest possible reads?
- Are conflicts rare enough for last-write-wins, or do they need application reconciliation?
- Is the data read-heavy, write-heavy, or balanced?
- How much repair traffic is acceptable in the steady state?

## Requirements

### Functional

- Support reads and writes through a replica set of size `N`.
- Define read quorum `R` and write quorum `W` intentionally.
- Detect and repair divergent replicas after failures.

### Non-functional

- Keep latency within reason for the chosen quorum sizes.
- Make divergence visible rather than silently optimistic.
- Bound repair work and hotspot impact.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Replica set size | 3 or 5 | changes failure tolerance and cost |
| Read/write mix | 80/20 | influences whether to bias cheaper reads or writes |
| P99 target | <25 ms | constrains quorum fan-out |
| Conflict rate | <0.1% normal, much higher during incidents | determines reconciliation strategy |
| Rough cost | extra replicas, quorum fan-out, repair traffic | exposes the price of higher availability |

## Architecture

A useful interview structure:

1. State `N`, `R`, and `W`.
2. Explain whether `R + W > N` is required for the critical path.
3. Name the versioning strategy, such as timestamps, vector clocks, or application version checks.
4. Explain read repair and anti-entropy without pretending they are free.

## Data Model & APIs

Helpful per-record metadata:

```text
record -> {
  key,
  value,
  version,
  timestamp,
  tombstone
}
```

Useful interfaces:

- `Put(key, value, write_quorum)`
- `Get(key, read_quorum)`
- `Repair(key, expected_version)`
- `MerkleSync(range)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| read returns stale value despite quorum choice | version mismatch across replicas | stronger `R`, read repair, or bounded stale contract |
| write succeeds on too few replicas | write acks below intended `W` | reject success or downgrade guarantee explicitly |
| repair traffic overloads the cluster | background sync queue and network saturation rise | throttle repair and isolate hot ranges |
| conflict resolution loses intent | reconciliation incidents appear | domain-aware merge or tighter serialization for critical entities |

## Observability

- metric: read and write quorum latency by path
- metric: divergent replica count and repair backlog
- metric: stale-read detection on sampled quorum reads
- log: conflicting version sets and chosen resolution
- trace: quorum fan-out and slow replica contributors
- SLO: freshness or conflict-rate objective paired with read/write latency target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| overlapping read and write quorums | stronger odds of fresh reads | more latency and lower partial-failure availability | tiny quorums everywhere |
| background anti-entropy | eventual convergence | network and disk overhead | relying only on foreground reads to repair |
| simple conflict resolution | lower complexity | can hide lost intent | pretending conflicts never happen |

## Interview It

**Google framing:** "Design a metadata store that stays up through node loss but still gives a credible freshness story." The signal is whether you go beyond quorum math into divergence handling.

**Cloudflare framing:** "Design a distributed edge data plane store with replica loss and repair." The signal is whether you talk about latency, conflict handling, and repair cost under global traffic.

**Follow-ups:**
1. What changes when the data is 95% reads?
2. What if writes are safety-critical and conflicts are unacceptable?
3. What if one replica is always slower but still healthy?
4. How do tombstones change deletion semantics?
5. When would you reject quorums and pick a leadered system instead?

## Ship It

- `outputs/quorum-planning-sheet.md`
- `outputs/read-repair-checklist.md`

## Exercises

1. **Easy** - Compare `N=3, R=1, W=3` against `N=3, R=2, W=2`.
2. **Medium** - Explain how read repair helps and what it costs.
3. **Hard** - Design quorum choices for edge policy reads where stale data is risky but availability still matters.

## Further Reading

- [Amazon DynamoDB papers and background](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) - useful historical grounding for quorum and repair trade-offs
- [Designing Data-Intensive Applications](https://dataintensive.net/) - practical explanation of quorum caveats and convergence work
