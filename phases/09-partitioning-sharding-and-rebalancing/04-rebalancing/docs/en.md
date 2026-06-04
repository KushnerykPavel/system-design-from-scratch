# Rebalancing Without Taking an Outage

> Moving data is easy until serving traffic still needs it.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Plan rebalance rollouts that protect the serving path by pacing moves, versioning ownership, and proving safety before switching traffic.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/03-placement`, `10-reliability-retries-and-backpressure/04-load-shedding`
**Estimated time:** ~75 min
**Primary artifact:** rebalance planner + failure checklist

## The Problem

Every sharded system eventually needs to rebalance:

- new nodes join
- old nodes drain
- one pool is too hot
- one region or AZ loses capacity

The mistake is to treat rebalancing as a background copy job. It is a serving-path change with correctness, capacity, and control-plane risk.

This lesson focuses on safe motion:

- move in bounded batches
- keep ownership versioned
- dual-serve or redirect carefully
- verify before final cutover

## Clarify

- Are we moving cache ownership, durable data, or both?
- Can reads be served from source and target during the move?
- What is the spare capacity budget for copy traffic?
- Is the main risk latency, stale routing, write loss, or replica under-protection?

If the interviewer is vague, assume durable partition moves under live read and write traffic with limited spare bandwidth.

## Requirements

### Functional

- Move ownership from source to target nodes safely.
- Keep reads and writes correct during transition.
- Support pause, rollback, and resumable progress.

### Non-functional

- Bound extra copy traffic and tail latency regression.
- Avoid mass ownership flips or control-plane confusion.
- Keep operator visibility high enough to catch silent divergence.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Data to move | 180 TB total candidate set | large enough that migration takes days, not minutes |
| Live write rate | 25K writes/s on moved ranges | requires change capture or dual-write strategy |
| Spare migration bandwidth | 12 Gbps cluster-wide | move pacing must respect serving traffic |
| Concurrent moves allowed | 4-8 ranges at once | controls blast radius |
| Rough cost | copy traffic + verification reads + warmup misses | rebalance has real temporary operating cost |

## Architecture

Safe rebalance often follows this sequence:

1. mark source ranges as `moving`
2. start background copy to target
3. stream deltas or dual-write while copy catches up
4. verify checksums, lag, and replica health
5. switch routing epoch
6. monitor, then retire source ownership

Key idea:

```text
control plane decides ownership version
serving path checks version
copy path stays rate-limited
rollback keeps source valid until target is proven
```

## Data Model & APIs

Useful metadata:

- `RangeMove(range_id, source, target, state, epoch, lag_bytes, verified_at)`
- `Ownership(range_id, epoch, owner_set)`

Useful interfaces:

- `StartMove(range_id, target)`
- `PauseMove(range_id)`
- `CutoverMove(range_id, expected_epoch)`
- `AbortMove(range_id)`

Senior-level detail:

- writes need either dual-write, change-log replay, or short freeze windows
- reads need a clear source of truth during cutover
- routing caches need epoch invalidation or short TTLs

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| too many moves saturate serving traffic | p99 latency and copy queue depth rise together | batch limits and rate caps |
| stale ownership cache sends writes to old node | stale-epoch errors or mismatch counters rise | versioned routing and redirect-with-epoch |
| target copy looks complete but missed recent writes | checksum mismatch or replication lag persists | CDC replay and pre-cutover verification |
| rollback path is missing | move pause becomes outage during incident | keep source valid until post-cutover health passes |

## Observability

- metric: bytes copied, lag bytes, and verification status per move
- metric: serving latency and error rate correlated with active move count
- metric: stale-epoch redirects and misrouted write attempts
- log: move state transitions with range ID, source, target, and epoch
- trace: requests crossing a move boundary, including redirects and dual-read fallback
- SLO: rebalances should complete without materially burning the serving path error budget

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| limited concurrent moves | smaller blast radius | slower completion | moving everything at once |
| dual-write or CDC catch-up | safer cutover under active writes | more implementation complexity | blind copy plus immediate flip |
| source retained through verification | strong rollback path | temporary extra capacity use | deleting source immediately after copy |

## Interview It

**Google framing:** "Your storage cluster needs rebalancing after capacity expansion." The signal is whether you treat ownership cutover as a correctness and rollout problem, not just a copy problem.

**Cloudflare framing:** "Rebalance customer state across a global fleet without control-plane instability." The signal is whether you discuss epochs, safe drains, and bounded move concurrency.

**Follow-ups:**
1. What if the target node falls behind on change replay?
2. What if routing caches are stale for 30 seconds?
3. How do you pause safely when serving latency spikes mid-move?
4. What if the rebalance is triggered by an AZ evacuation?
5. When is a short write freeze acceptable instead of dual-write complexity?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/failure-checklist-rebalancing.md`
- `outputs/interview-card-rebalancing.md`

## Exercises

1. **Easy** — List the minimum state machine for one range move.
2. **Medium** — Explain how you would cut over reads and writes separately during a live migration.
3. **Hard** — Design an AZ evacuation rebalance where capacity is already tight and rollback must stay available.

## Further Reading

- [Bigtable paper](https://research.google/pubs/pub27898/) — useful background for tablet movement and serving control
- [Google SRE books](https://sre.google/books/) — strong operational framing for safe rollouts and capacity-aware changes
