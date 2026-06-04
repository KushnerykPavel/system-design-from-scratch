# Geo-Distributed Consistency Trade-offs

> Distance turns consistency into a latency budget conversation long before it becomes a theorem conversation.

**Type:** Learn
**Company focus:** Google
**Learning goal:** Pick geo-distributed consistency semantics that match user expectations, write geography, and failure tolerance instead of defaulting to vague "eventual consistency" answers.
**Prerequisites:** `08-consistency-replication-and-transactions/01-consistency-spectrum`, `08-consistency-replication-and-transactions/03-quorums`, `13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive`
**Estimated time:** ~60 min
**Primary artifact:** consistency decision matrix

## The Problem

A global service spans regions, but not every operation needs the same consistency semantics. Some reads can lag, some writes can be home-region only, and some user-visible workflows break if ordering is weak.

This lesson trains the senior move: choose semantics per workflow and explain the cost in latency, availability, and operational repair.

## Clarify

- Which user actions require fresh reads immediately after writes?
- Are writes naturally home-region scoped, or truly global?
- What inconsistency is visible but acceptable, and what is unacceptable?
- How should the system behave during inter-region partition or lag?

If the interviewer is vague, separate metadata reads, user-generated writes, and control-plane changes instead of giving one consistency answer for the whole system.

## Requirements

### Functional

- Support a mix of reads and writes across multiple regions.
- Preserve stronger semantics where user workflows demand them.
- Make degraded behavior explicit during partition or lag.

### Non-functional

- Keep global write latency within a realistic budget.
- Avoid pretending strong semantics are free across oceans.
- Ensure repair and reconciliation are operationally manageable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Inter-region RTT | 70 to 180 ms | bounds strong cross-region write latency |
| Write classes | 3 to 5 | different workflows deserve different semantics |
| Partition tolerance window | minutes during incidents | determines degraded-mode behavior |
| Repair backlog | thousands to millions of mutations | affects reconciliation tooling |
| Rough cost | replicated storage + coordination + repair ops | stronger global semantics cost both latency and money |

## Architecture

```text
client
  -> nearest region
     -> read path
        -> local replica or cache
     -> write path
        -> home region or quorum path
     -> async replication / repair
```

Common patterns:

1. **Read local, write home region**
   - simple ownership
   - good for user-scoped data
2. **Per-entity home region**
   - avoids multi-writer conflict for most objects
3. **Selective quorum writes**
   - used only for workflows that justify higher latency
4. **Asynchronous replicated reads**
   - acceptable when freshness can lag explicitly

## Data Model & APIs

Core entities:

- `ConsistencyClass`
- `HomeRegion`
- `VersionVector` or `LogicalVersion`
- `RepairEvent`
- `ReadPolicy`

Useful APIs:

- `Write(entity, consistency_class)`
- `Read(entity, freshness_bound)`
- `ResolveConflict(entity_versions)`
- `ExplainFreshness(entity)`

Important modeling question: which APIs promise read-your-writes, monotonic reads, or only bounded staleness?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| local reads are too stale for user workflow | freshness SLA breach or user-visible mismatch | route critical reads to home region or tighten bounded staleness |
| multi-writer conflict exceeds repair capacity | conflict queue growth | constrain ownership and reduce global concurrent writes |
| partition handling is undefined | inconsistent operator action during incident | define fail-closed, read-only, or stale-read policy per workflow |
| control-plane metadata lags globally | version skew and rollout mismatch | use stronger replication semantics for critical control data |

## Observability

- metric: replication lag and bounded-staleness SLA breach rate
- metric: conflict volume, repair backlog, and repair age
- metric: share of reads served locally versus escalated to home region
- log: write path chosen, consistency class, and conflict outcome
- trace: cross-region hops on critical reads and writes
- SLO: freshness promises should be stated per workflow, not only globally

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| per-workflow consistency classes | matches cost to user need | more design complexity | one vague consistency mode for everything |
| home-region writes | simpler correctness | some users pay higher write latency | arbitrary multi-writer global writes |
| bounded staleness for reads | lower latency and higher availability | occasional stale views | global synchronous reads everywhere |

## Interview It

**Google framing:** "How would you handle consistency across regions?" Expect questions on what the user actually notices and which semantics deserve the cost.

**Cloudflare framing:** "What consistency can you afford globally at edge scale?" Expect pressure on bounded staleness, control-plane data, and partitions.

**Follow-ups:**
1. Which workflow needs read-your-writes and which does not?
2. What happens during an inter-region partition?
3. How do you explain stale reads to product partners?
4. What data deserves stronger replication than the rest?
5. How does the design change at 10x write volume?

## Ship It

- `outputs/tradeoff-matrix-geo-consistency.md`
- `outputs/interview-card-geo-consistency.md`

## Exercises

1. **Easy** — Assign consistency classes to profile reads, comment writes, and billing updates.
2. **Medium** — Redesign a collaborative document backend for bounded staleness reads and home-region writes.
3. **Hard** — Add a globally replicated control plane with stricter semantics than user-content data.

## Further Reading

- [Designing Data-Intensive Applications - Replication chapter summary](https://www.oreilly.com/library/view/designing-data-intensive-applications/9781491903063/) — classic framing for replication trade-offs
- [Google Spanner paper](https://research.google/pubs/pub39966/) — strong reference point for what globally consistent writes really cost and enable
