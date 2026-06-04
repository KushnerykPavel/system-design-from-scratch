# Transactions, Isolation, and Hotspotting

> The best transaction answer is not "use ACID." It is knowing where atomicity matters, which anomalies are unacceptable, and what contention that choice creates.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Choose transaction boundaries and isolation levels deliberately while recognizing contention, lock cost, and hotspot risk.
**Prerequisites:** `02-leader-follower`, `05-storage-indexing-and-access-patterns/04-hot-and-cold-data`, `07-queues-streams-and-workflows/06-outbox-and-cdc`
**Estimated time:** ~75 min
**Primary artifact:** transaction review checklist

## The Problem

Candidates often reach for transactions as a correctness blanket. Interviewers want to know whether you understand the price:

- stronger isolation can reduce concurrency
- hot rows or counters can serialize the workload
- cross-entity transactions can drag many systems into one blast radius
- some workflows are better split into local transactions plus async compensation

This lesson helps you talk about transactions precisely rather than reverently.

## Clarify

- Which invariant must never be broken, even briefly?
- Is the workload mostly single-row, multi-row, or cross-service?
- What anomaly is unacceptable: lost update, dirty read, write skew, double spend?
- Where are the likely hotspots: account row, inventory item, quota bucket, or tenant counter?

## Requirements

### Functional

- Protect critical invariants with transactional boundaries where needed.
- Choose an isolation level that blocks the dangerous anomalies.
- Identify hotspots and redesign options before they become a scale wall.

### Non-functional

- Keep contention and latency visible.
- Avoid expanding distributed transactions across too many services.
- Preserve debuggability when retries and lock waits occur.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak write QPS | 12K req/s | enough that lock contention matters |
| Hot entities | top 0.5% keys drive 25% of writes | exposes serialization bottlenecks |
| Multi-row transactions | 15% of writes | enough to make isolation meaningful |
| P99 write target | <40 ms | constrains lock and retry behavior |
| Rough cost | lock waits, retries, coordination, reduced concurrency | shows stronger isolation is not free |

## Architecture

A senior transaction answer usually does four things:

1. Define the invariant.
2. Choose the smallest transactional boundary that protects it.
3. Name the isolation level or concurrency control reason.
4. Call out hotspots and redesign options early.

Common redesign moves:

- escrow or reservation counters
- per-entity serialization queues
- partitioning by account or tenant
- local transaction plus outbox instead of cross-service two-phase work

## Data Model & APIs

Helpful interfaces:

- `Begin()`
- `Transfer(from, to, amount, expected_version)`
- `Reserve(item_id, qty)`
- `Commit()` / `Rollback()`

Helpful metadata:

```text
row -> {
  id,
  version,
  owner_key,
  last_updated_at
}
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hotspot row serializes the workload | lock wait and abort rate rise on one key | shard, reserve, or pre-aggregate the hotspot |
| weak isolation permits dangerous anomaly | invariant violations appear despite success responses | tighten isolation or add version checks |
| transaction spans too many systems | latency and blast radius grow | shrink boundary and use async coordination |
| retries amplify contention | retry storm on lock conflicts | backoff, retry budgets, and admission control |

## Observability

- metric: lock wait time and transaction abort rate
- metric: hottest rows or keys by conflict count
- metric: invariant-violation or reconciliation incident count
- log: transaction retries with entity key, conflict cause, and isolation path
- trace: slow transaction spans with lock acquisition timing
- SLO: critical write success plus p99 latency and contention budget

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| stronger isolation on narrow critical path | protects invariants | lower concurrency and higher latency | blanket strong isolation everywhere |
| smaller transaction boundary | simpler ownership and better scale | more async coordination outside boundary | distributed transaction across every service |
| hotspot redesign | restores throughput | more modeling complexity | assuming the hot key will average out |

## Interview It

**Google framing:** "Design account transfers or inventory reservation under heavy concurrency." The signal is whether you identify the invariant and hotspot before naming an isolation level.

**Cloudflare framing:** "Design mutable control-plane state with tenant hotspots and strict correctness." The signal is whether you discuss ownership, contention, and where not to stretch transactions.

**Follow-ups:**
1. Which anomaly are you actually preventing?
2. What if one tenant generates 40% of all writes?
3. What if the system needs to update storage and publish an event?
4. When is optimistic concurrency enough?
5. When would you replace a transaction with a saga?

## Ship It

- `outputs/transaction-review-checklist.md`
- `outputs/hotspot-playbook-transactions.md`

## Exercises

1. **Easy** - Pick an isolation need for a simple inventory reservation flow.
2. **Medium** - Explain how a hot account row can throttle a payments-like system.
3. **Hard** - Redesign a cross-service checkout flow so only the truly critical invariant stays transactional.

## Further Reading

- [Designing Data-Intensive Applications](https://dataintensive.net/) - strong explanations of isolation and concurrency control trade-offs
- [System design notes](https://github.com/liquidslr/system-design-notes) - useful general interview structure before transaction nuance
