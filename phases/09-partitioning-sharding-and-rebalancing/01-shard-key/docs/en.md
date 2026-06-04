# Choosing a Shard Key

> A shard key is a workload bet, not a schema afterthought.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Choose shard keys by tying access patterns, write skew, isolation, and migration risk to a concrete storage layout instead of defaulting to an ID column.
**Prerequisites:** `05-storage-indexing-and-access-patterns/02-access-pattern-first`, `08-consistency-replication-and-transactions/05-transactions`
**Estimated time:** ~75 min
**Primary artifact:** shard-key evaluator + trade-off matrix

## The Problem

Interview answers often say "we'll shard by user ID" before naming the hottest reads, write concurrency, tenant isolation needs, or future rebalancing plan.

That shortcut creates predictable failures:

- one tenant dominates a shard
- range scans become scatter-gather queries
- transactional boundaries cross shards too early
- migrations become full-table rewrites

This lesson gives you a practical frame: a shard key is only good if it makes the dominant workload, failure model, and growth path easier at the same time.

## Clarify

- What is the dominant read or write unit: user, tenant, region, object, or time bucket?
- Which requests must stay single-shard for latency or transactional reasons?
- Is the top operational risk hotspotting, uneven tenant size, or expensive cross-shard queries?
- Will the data model need future moves such as tenant splits, regional pinning, or archival tiering?

If the interviewer is vague, assume a multi-tenant product with strong per-tenant locality, uneven tenant sizes, and a need to support future tenant movement.

## Requirements

### Functional

- Choose a shard key that keeps the primary workload efficient.
- Explain how lookups, list queries, and writes route to the correct partition.
- Show what metadata or directory service is needed when the key is not self-locating.

### Non-functional

- Keep hot partitions and skew visible.
- Preserve room for rebalancing and tenant movement.
- Minimize unnecessary cross-shard fanout on the critical path.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Tenant count | 200K active tenants | enough variance to make average-based sizing misleading |
| Peak writes | 80K writes/s | exposes write hotspot risk for monotonic or narrow keys |
| Largest tenant share | 9% of traffic | one customer can dominate a shard if isolation is weak |
| Fanout-sensitive queries | 20% of reads | affects whether the key can support list and aggregation paths |
| Rough cost | shard count + directory metadata + migration overhead | a "simple" key is not free if it forces expensive future moves |

## Architecture

Pick the key from the workload backward:

1. Identify the unit that should stay local.
2. Decide whether routing is direct or requires metadata lookup.
3. Name which access patterns will stay single-shard.
4. Say what you will do about the patterns that do not.

Example for a multi-tenant issue tracker:

```text
request
  -> routing layer
  -> tenant directory
  -> shard owning tenant or tenant bucket
  -> local indexes for per-tenant lists
```

Good default patterns:

- shard by `tenant_id` when isolation and per-tenant locality dominate
- shard by `user_id` when user-owned write volume is dominant and tenant shape is flatter
- shard by `region + tenant_bucket` when legal placement and regional latency matter

Weak pattern:

- shard by random object ID when most useful queries are tenant- or user-scoped

## Data Model & APIs

Example entities:

- `TenantDirectory(tenant_id, shard_id, state, split_parent)`
- `Issue(issue_id, tenant_id, project_id, assignee_id, created_at, status)`
- `Project(project_id, tenant_id, ... )`

Useful routing interfaces:

- `ResolveTenant(tenant_id) -> shard_id`
- `CreateIssue(tenant_id, payload)`
- `ListIssues(tenant_id, filters, page_token)`

Design notes:

- keep the shard key present in primary records and secondary indexes
- use opaque pagination tokens that carry shard-local continuation state
- avoid APIs that require global scans unless they are explicitly async or analytical

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| key chosen for even distribution but not query locality | high scatter-gather rate on common reads | reframe around dominant access pattern, not just balance |
| one large tenant overwhelms a shard | per-shard CPU, queue depth, and tail latency skew sharply | introduce tenant splits, bucketed keys, or dedicated placement |
| routing metadata becomes stale during moves | miss or stale-owner errors after migration | dual-read routing state and versioned directory updates |
| time-based or sequential key creates write hotspot | shard write concentration and lock contention rise | hash or bucket the write key while preserving lookup path |

## Observability

- metric: requests per shard and per tenant percentile distribution
- metric: percentage of critical reads requiring cross-shard fanout
- metric: shard directory lookup latency and stale-owner retries
- log: routing decision, shard ID, and remap version on failures
- trace: directory lookup to shard execution for latency attribution
- SLO: dominant user flows stay single-shard and within latency targets for the 99th percentile tenant

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| shard by tenant | strong isolation and query locality | skew from large tenants | random object ID sharding |
| hashed tenant buckets | smoother distribution | extra routing indirection and harder tenant moves | raw tenant ID when tenant sizes vary wildly |
| regional prefix in key | legal placement and lower latency | more shards and more complex global views | fully global shard pool |

## Interview It

**Google framing:** "Design storage for a multi-tenant project management system." The signal is whether you choose the key from the access pattern and transaction boundary instead of from convenience.

**Cloudflare framing:** "Design configuration storage for millions of customer zones." The signal is whether you cover noisy customers, placement control, and migration safety.

**Follow-ups:**
1. What changes if one tenant becomes 100x larger than the median?
2. What if list-by-project is now hotter than list-by-tenant?
3. How would you support tenant moves with no write outage?
4. When is directory indirection worth it over self-locating hashing?
5. What does your key choice do to global reporting queries?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-shard-key.md`
- `outputs/tradeoff-matrix-shard-key.md`

## Exercises

1. **Easy** — Choose between `tenant_id`, `user_id`, and `project_id` for a bug tracker and defend one.
2. **Medium** — Redesign the key after one tenant starts driving 15% of all writes.
3. **Hard** — Explain a migration from tenant-based sharding to regional tenant buckets with no customer-visible outage.

## Further Reading

- [Spanner schema design best practices](https://cloud.google.com/spanner/docs/schema-design) — useful for locality and key-shape reasoning
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline framing for sizing and workload-first design
