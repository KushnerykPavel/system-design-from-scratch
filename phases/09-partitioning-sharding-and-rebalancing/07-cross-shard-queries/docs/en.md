# Cross-Shard Queries and Aggregation

> If every query touches every shard, you did not really partition the problem.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Decide when scatter-gather is acceptable, when to precompute or index differently, and how to explain the latency and correctness costs of cross-shard reads.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/01-shard-key`, `07-queues-streams-and-workflows/06-outbox-and-cdc`
**Estimated time:** ~60 min
**Primary artifact:** trade-off matrix + interview card

## The Problem

Sharding helps scale writes and ownership, but it makes global reads harder. The danger is pretending every query can stay as simple after sharding as it was before.

Cross-shard queries show up as:

- global search or filter views
- tenant-spanning analytics
- sorted feeds across many owners
- counts, joins, or aggregates across partitions

Senior answers do not ban cross-shard work. They classify it:

- synchronous and latency-sensitive
- asynchronous and precomputed
- approximate
- restricted by product scope

## Clarify

- Is the query on the critical user path, or can it be async?
- Must results be exact, fresh, or globally ordered?
- Is the query common enough to deserve a dedicated index or materialized view?
- Can the product scope be narrowed to one tenant, region, time window, or cohort?

If the prompt is vague, assume tenant-local requests must be fast, while cross-tenant analytics can be eventual or precomputed.

## Requirements

### Functional

- Support the important global or multi-shard read patterns.
- Explain whether reads are direct scatter-gather, precomputed, or served from a derived system.
- Preserve clear semantics for ordering, freshness, and pagination.

### Non-functional

- Keep critical-path latency under control.
- Avoid fanout explosions as shard count grows.
- Bound cost on storage, compute, and network when derived views are introduced.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Shard count | 256 shards | direct fanout cost becomes material |
| Tenant-local read QPS | 220K req/s | should stay local and cheap |
| Cross-shard read QPS | 2K req/s | maybe okay for bounded scatter-gather |
| Heavy analytical requests | 50 req/s but expensive | often better in a derived system |
| Rough cost | fanout RPCs + merge CPU + derived index storage | query support choices have real operating cost |

## Architecture

Use a decision tree:

1. Can the product scope stay shard-local?
2. If not, can the answer be precomputed or materialized?
3. If it must fan out, can the fanout be bounded and merged quickly?
4. If even that is too costly, move the query to search, analytics, or batch systems.

Common patterns:

- shard-local list views for core UX
- materialized global counters or rollups via CDC
- search cluster for flexible multi-shard filtering
- asynchronous exports for large analytical queries

## Data Model & APIs

Useful structures:

- `ShardLocalIndex(...)`
- `GlobalRollup(metric_key, time_bucket, value)`
- `SearchDocument(...)`

Useful interfaces:

- `ListTenantItems(tenant_id, page_token)`
- `QueryGlobalDashboard(time_window, filters)`
- `ExportCrossTenantReport(job_spec)`

Important semantics to state:

- global pagination often needs merge tokens, not simple offsets
- exact global ordering is expensive under fanout
- derived systems usually trade freshness for latency and cost

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| synchronous scatter-gather grows with shard count | p99 latency rises with topology size | precompute or bound fanout by scope |
| global counters lag silently | dashboard drift or parity checks fail | expose freshness lag and rebuild workflows |
| flexible filters overload primary storage | cross-shard query CPU and network spike | move to search or analytics index |
| pagination across shards becomes unstable | duplicates or gaps appear between pages | merge-aware continuation tokens and stable ordering keys |

## Observability

- metric: fanout width, merge latency, and cross-shard query rate
- metric: freshness lag on materialized views or rollups
- metric: query cost by class, including remote RPC count and bytes merged
- log: query plan choice, bounded shard set, and fallback path
- trace: coordinator plus shard sub-requests with merge timing
- SLO: critical user-path queries should not depend on unbounded cross-shard fanout

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| bounded scatter-gather | simpler than extra systems for low-rate queries | latency and coordination cost | forcing everything into one shard |
| materialized views | fast reads on common global queries | lag, pipeline complexity, and duplicate storage | live fanout for every dashboard load |
| search or analytics backend | flexible global queries | more systems and sync complexity | loading primary shards with exploratory filters |

## Interview It

**Google framing:** "Your product now needs a global admin dashboard after adopting tenant sharding." The signal is whether you distinguish serving-path reads from analytical reads.

**Cloudflare framing:** "Expose customer-wide and fleet-wide views without melting the control plane." The signal is whether you scope global queries and use derived systems deliberately.

**Follow-ups:**
1. What if the dashboard needs exact counts within 2 seconds?
2. How do you paginate a globally ordered feed across many shards?
3. When is approximate aggregation acceptable?
4. What if cross-shard queries are rare today but growing quickly?
5. How do you explain freshness lag to product stakeholders?

## Ship It

- `outputs/interview-card-cross-shard-queries.md`
- `outputs/tradeoff-matrix-cross-shard-queries.md`

## Exercises

1. **Easy** — Classify three read patterns as shard-local, scatter-gather, or derived.
2. **Medium** — Design a materialized view for a fleet-wide operational dashboard.
3. **Hard** — Explain a global, sorted feed across 256 shards with pagination and freshness constraints.

## Further Reading

- [Streaming Systems](https://www.oreilly.com/library/view/streaming-systems/9781491983867/) — strong background for derived views and aggregation pipelines
- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline for deciding when to simplify scope versus adding new systems
