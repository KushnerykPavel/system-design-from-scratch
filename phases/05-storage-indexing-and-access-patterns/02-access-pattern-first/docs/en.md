# Primary Access Pattern First

> Schema design gets easier when you admit which query pays the bill.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Turn a vague product prompt into an ordered list of access patterns, then use that ranking to justify storage layout, indexes, and denormalization choices.  
**Prerequisites:** `01-storage-models`, `04-apis-contracts-and-schema-evolution/03-pagination-and-filtering`  
**Estimated time:** ~75 min  
**Primary artifact:** design review prompt + capacity sheet  

## The Problem

Teams often say they are "designing the schema," but what they are really doing is choosing which reads and writes will be cheap. The mistake is trying to support every imagined query equally well.

This lesson trains you to:

- identify the primary access pattern first
- rank the next few patterns by business importance and traffic
- intentionally degrade or offload the rest

## Clarify

- Which user action happens most often and matters most to product success?
- Which requests are latency-sensitive versus batchable?
- Are reads mostly by ID, by tenant, by time range, or by compound filter?
- Which future query is tempting but should probably not shape the primary model yet?

## Requirements

### Functional

- Enumerate and rank the top access patterns for a workload.
- Choose schema and index layout based on the top-ranked paths.
- Separate online serving queries from exports, analytics, and backfills.

### Non-functional

- Keep the hot path simple and predictable.
- Avoid over-indexing for low-value rare queries.
- Make cost and operational complexity visible when supporting secondary patterns.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Primary read path | 90K req/s | the top query should shape storage layout |
| Secondary list path | 8K req/s | enough to matter, but not enough to dominate the model |
| Writes | 15K req/s | denormalization and index updates must stay affordable |
| Peak factor | 6x on one tenant cohort | stresses compound-key and partition design |
| Rough cost | primary table + 2-3 indexes + export pipeline | helps reject "optimize everything" thinking |

## Architecture

Recommended workflow:

1. List the user-visible reads and writes.
2. Rank them by traffic and business criticality.
3. Shape the primary data model for the top one or two.
4. Move weakly ranked queries to async jobs, search systems, or offline stores.

Example:

```text
write comment
read thread by thread_id
list recent threads by user_id
search comments by keyword
export all comments for compliance
```

Only the first three should usually shape the online serving model.

## Data Model & APIs

For a tenant-scoped feed product, a strong answer might say:

- primary key: `tenant_id + item_id`
- serving index: `tenant_id + created_at desc`
- export path: asynchronous job to object storage
- keyword search: derived search index, not the primary OLTP store

The key question is not "can the primary store do all of these?" It is "should it?"

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| low-value query shapes the schema | common path gets slower and more complex | re-rank patterns and move rare queries to derived systems |
| every stakeholder adds one more index | write cost and storage budget drift upward | require pattern ranking and explicit approval for new indexes |
| export or analytics queries hit OLTP directly | serving latency spikes during backfills | isolate heavy reads into replicas, warehouses, or async jobs |
| pattern ranking is not revisited after growth | old model no longer fits current traffic | review top access paths quarterly with real metrics |

## Observability

- metric: request volume and latency by access pattern label
- metric: write amplification per logical write
- metric: share of requests served by primary path versus fallback or export path
- log: rejected unbounded query attempts
- trace: which index or storage path was selected
- SLO: top-ranked path meets latency target under expected peak factor

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| optimize for top 1-2 patterns | keeps the hot path fast and understandable | weaker support for rare ad hoc queries | equal optimization for every stakeholder request |
| async exports for heavy reads | protects serving system from scans | higher user wait time for bulk operations | letting batch jobs share the serving path |
| derived search index | preserves OLTP simplicity | more moving parts and lag | forcing keyword search into the transactional store |

## Interview It

**Google framing:** "Design storage for a large task-tracking product." The signal is whether you identify the dominant query paths before naming tables or databases.

**Cloudflare framing:** "Design storage for customer configuration and fleet lookup APIs." The signal is whether you separate ultra-hot serving reads from audit, search, and export paths.

**Follow-ups:**
1. What if a rare admin query becomes a major product surface?
2. What if one tenant generates most of the list traffic?
3. What if exports must finish within five minutes?
4. How do you decide when to add a secondary index versus a new derived system?
5. Which access pattern would you deliberately not optimize in v1?

## Ship It

- `outputs/design-review-access-pattern-first.md`
- `outputs/capacity-sheet-access-pattern-first.md`

## Exercises

1. **Easy** — Rank the access patterns for a notification inbox.  
2. **Medium** — Redesign a schema that currently optimizes search before primary reads.  
3. **Hard** — Explain the online, search, and export paths for a multitenant activity feed.  

## Further Reading

- [Google Spanner schema design best practices](https://cloud.google.com/spanner/docs/schema-design) — useful reference for access-pattern-driven keys and indexes  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — strong framing for workload-first storage design  
