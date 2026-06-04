# Indexes, Secondary Indexes, and Write Amplification

> Every index is a promise to make one read cheaper by making many writes more expensive.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Explain how primary and secondary indexes change read latency, write cost, storage growth, and failure behavior in an interview design.  
**Prerequisites:** `02-access-pattern-first`, `02-estimation-and-cost/07-bottleneck-math`  
**Estimated time:** ~75 min  
**Primary artifact:** trade-off matrix + observability checklist  

## The Problem

Indexes are often treated as a free answer to "how will we query this?" In real systems, every extra index increases storage usage, write path work, backfill time, and repair complexity.

This lesson helps you speak precisely about:

- clustered or primary indexes
- secondary indexes
- covering indexes
- write amplification and index backfills

## Clarify

- Which read path truly needs low latency, and which can tolerate async or batch behavior?
- How frequently do writes happen relative to reads?
- Are secondary queries bounded and predictable, or open-ended?
- What is the tolerance for stale or eventually built derived indexes?

## Requirements

### Functional

- Choose which queries deserve direct index support.
- Explain primary key and secondary index layout.
- Estimate the write amplification of maintaining the chosen indexes.

### Non-functional

- Keep write latency and storage growth within budget.
- Avoid index sprawl caused by low-value filter combinations.
- Make index build and repair risk visible during rollout.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Logical writes | 40K writes/s | index count multiplies write work |
| Indexed reads | 150K reads/s | justifies some secondary structures |
| Average row size | 1.5 KB | affects covering index decisions and storage blowup |
| Peak factor | 3x during import windows | stresses compaction and backfill scheduling |
| Rough cost | base table + replicas + N indexes | shows why "add an index" is not free |

## Architecture

Think of the write path explicitly:

```text
write request
  -> primary row update
  -> secondary index updates
  -> replication / WAL / compaction
```

The senior move is to say:

1. which query gets first-class index support
2. which query is good enough with a scan over a bounded partition
3. which query belongs in a separate search or analytics system

## Data Model & APIs

Example for a feed item table:

- primary key: `tenant_id + item_id`
- list index: `tenant_id + created_at desc`
- moderation index: `tenant_id + status + created_at`

Then answer:

- how many index entries each write updates
- which indexes are sparse or partial
- how new indexes are rolled out or backfilled safely

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| too many indexes slow writes unexpectedly | logical write volume stable but commit latency climbs | prune low-value indexes and move rare reads elsewhere |
| unbounded composite index growth | storage and compaction cost drift upward | require access-pattern review before new compound indexes |
| backfill competes with serving traffic | p99 latency spikes during index build | throttle backfill and use shadow validation before cutover |
| stale secondary index produces confusing reads | mismatch checks and repair metrics fire | expose freshness lag and keep correctness in primary record |

## Observability

- metric: write latency by number of maintained indexes
- metric: index storage size and growth rate
- metric: index hit ratio by query class
- metric: backfill progress and lag
- log: query planner fallback or rejected unindexed query
- SLO: indexed reads stay fast without pushing write latency outside budget

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| targeted secondary indexes | fast support for important reads | more write work and storage | scanning everything on the hot path |
| partial or sparse indexes | lower index cost for selective queries | extra complexity in query behavior | indexing every row equally |
| external search for flexible queries | protects OLTP write path | more systems and eventual consistency | piling many flexible indexes into primary store |

## Interview It

**Google framing:** "Design storage for a large issue tracker with many list views." The signal is whether you reason about which filters deserve index support and what that costs writes.

**Cloudflare framing:** "Design an indexable metadata layer for customer configuration state." The signal is whether you keep write amplification and rollout safety visible.

**Follow-ups:**
1. What changes if writes jump 10x while read volume stays flat?
2. What if a new admin filter is needed only once a day?
3. How do you add a new secondary index without hurting production?
4. What if a secondary index can lag by 30 seconds?
5. When is a covering index worth the extra storage?

## Ship It

- `outputs/tradeoff-matrix-indexes.md`
- `outputs/observability-checklist-indexes.md`

## Exercises

1. **Easy** — Choose one primary key and two secondary indexes for a notification table.  
2. **Medium** — Estimate the write cost difference between one index and four indexes.  
3. **Hard** — Explain why a search backend is better than adding another six compound indexes.  

## Further Reading

- [Use The Index, Luke](https://use-the-index-luke.com/) — practical grounding for index behavior and trade-offs  
- [Cloud Spanner secondary indexes](https://cloud.google.com/spanner/docs/secondary-indexes) — good reference for cost and backfill discussions  
