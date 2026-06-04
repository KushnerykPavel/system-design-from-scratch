# Consistency Trade-off Drill

> The strongest answer sounds like a guarantee map with failure handling, not a lecture on database theory.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Practice turning ambiguous requirements into a scoped consistency contract, replication choice, transaction boundary, and observability plan.
**Prerequisites:** `01-consistency-spectrum`, `02-leader-follower`, `03-quorums`, `05-transactions`, `06-sagas`, `07-time-assumptions`
**Estimated time:** ~60 min
**Primary artifact:** drill worksheet + scoring rubric

## The Problem

This phase drill combines the main questions senior candidates must answer under pressure:

- what guarantee does each critical flow need
- which replication or coordination model provides it
- where should transactions stop
- what degraded mode is acceptable under lag, partition, or partial failure

The goal is not to prove knowledge of every theorem. The goal is to make a small number of explicit, defensible choices.

## Clarify

- Which user journey breaks if data is stale or reordered?
- Which invariants require atomic protection, and which can reconcile later?
- Is the system leadered, quorum-based, or hybrid for different entities?
- What failure is more dangerous here: stale reads, lost writes, duplicate side effects, or blocked availability?

## Requirements

### Functional

- Define consistency guarantees by entity or flow.
- Pick a replication and transaction model intentionally.
- Explain compensation or replay where atomicity ends.

### Non-functional

- Keep guarantees honest under failover and lag.
- Bound latency and coordination cost.
- Make freshness, divergence, and contention visible.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read QPS | 200K req/s | enough to pressure freshness-routing decisions |
| Write QPS | 10K req/s | enough to make replication and transactions visible |
| Cross-service workflows | 8% of writes | enough to force saga discussion |
| Hot keys | top 1% create 20% of writes | exposes contention and routing risks |
| Rough cost | replicas, coordination latency, retries, workflow state | keeps trade-offs concrete |

## Architecture

Recommended drill sequence:

1. Clarify the highest-risk user and business failures.
2. Assign each critical entity a consistency target.
3. Pick the simplest replication model that supports those targets.
4. Keep transactions narrow and name where sagas begin.
5. Close with lag handling, observability, and redesign options.

Strong answer pattern:

- critical account or policy state: tighter consistency
- convenience or analytics reads: bounded stale
- multi-service side effects: local transaction plus saga
- time-based logic: fencing or versions where correctness depends on order

## Data Model & APIs

A strong drill answer often names:

- version or sequence metadata
- read path with `min_version` or `max_staleness`
- transactional boundary or idempotency key
- workflow or compensation status handle

Example surfaces:

```text
Read(id, min_version)
Write(id, expected_version)
Reserve(resource_id, qty)
StartWorkflow(workflow_id)
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| stale data reaches critical path | freshness mismatch or version skew rises | leader read, stronger quorum, or version-aware bypass |
| transaction hotspot throttles system | lock wait and abort concentration on one key | redesign hotspot or shrink boundary |
| compensation path is under-specified | stuck workflow age and manual queue grow | explicit state machine and operator recovery |
| clock or failover assumptions are too weak | fence violations or read-after-write misses appear | add fencing, buffers, or safer ordering metadata |

## Observability

- metric: freshness age and version skew for critical entities
- metric: commit latency, abort rate, and hotspot conflict concentration
- metric: workflow completion, compensation, and oldest in-progress age
- log: replication, promotion, and workflow recovery decisions
- trace: one user action across write, replication, transaction, and async completion
- SLO: critical correctness target paired with latency and recovery objectives

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| selective strong guarantees | correctness where it matters | more metadata and coordination | one uniform guarantee for every path |
| narrow transactions plus sagas | smaller blast radius | visible intermediate state and compensation work | wide distributed transactions |
| freshness-aware routing | cheaper scaling for noncritical reads | more operational logic | replica reads without explicit lag policy |

## Interview It

**Google framing:** "Design order state, inventory, billing visibility, and notifications for a global commerce workflow." The signal is whether you separate invariants and choose different consistency tools for each.

**Cloudflare framing:** "Design globally propagated policy changes, tenant metadata, and rollout status." The signal is whether you reason about read freshness, ownership, compensation, and skew under failure.

**Follow-ups:**
1. Which deep dive would you pick first and why?
2. What if one tenant becomes a write hotspot?
3. What if strict consistency doubles tail latency for the hot read path?
4. What if rollback is impossible for one external side effect?
5. How does the design change at 10x write rate?

## Ship It

- `outputs/drill-worksheet-consistency.md`
- `outputs/scoring-rubric-consistency.md`

## Exercises

1. **Easy** - Run the drill for profile settings plus analytics.
2. **Medium** - Run the drill for inventory reservation plus payment visibility.
3. **Hard** - Run the drill for global policy rollout with tenant hotspots, lag, and partial rollback.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) - useful baseline interview flow before applying replication and transaction nuance
- [Designing Data-Intensive Applications](https://dataintensive.net/) - background for consistency, replication, and transaction trade-offs
