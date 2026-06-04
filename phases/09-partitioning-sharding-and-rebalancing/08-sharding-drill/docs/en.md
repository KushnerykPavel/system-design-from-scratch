# Sharding Drill

> A sharding answer is only strong if it survives skew, movement, and follow-ups.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Integrate the phase by designing a sharded system end to end and checking whether the answer covers key choice, skew, placement, movement, cross-shard queries, and rollout safety.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/07-cross-shard-queries`
**Estimated time:** ~60 min
**Primary artifact:** drill sheet + answer checker

## The Problem

This drill uses one prompt to force the whole phase together:

"Design the storage and serving architecture for a multi-tenant feature flag platform with millions of tenants, heavy read traffic, uneven enterprise customers, and a requirement to migrate large tenants safely over time."

The point is not to name every database feature. The point is to show an interviewer that you can:

- choose a shard key from the workload
- protect against hotspotting and noisy neighbors
- explain placement and movement
- limit cross-shard pain
- ship the migration story safely

## Clarify

- Are reads mostly tenant-local configuration fetches, or are there many fleet-wide admin views too?
- How strict is freshness for config reads and control-plane writes?
- How uneven are the largest tenants, and do premium tiers justify dedicated placement?
- Do configuration writes need zero-downtime tenant movement?

If the interviewer stays vague, assume tenant-local reads dominate, a few large enterprise tenants create skew, and admin dashboards can tolerate eventual consistency.

## Requirements

### Functional

- Store and serve feature flag configuration per tenant.
- Support safe tenant moves, tenant splits, and growing shard count.
- Provide cross-tenant operational views without melting the primary path.

### Non-functional

- Low-latency tenant-local reads.
- Strong protection against noisy neighbors in the control plane.
- Safe migration and rollback under live customer traffic.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read QPS | 300K req/s | tenant-local reads must stay cheap |
| Write QPS | 15K req/s | control-plane updates need consistent routing |
| Largest tenant skew | 40x median | default placement will eventually fail |
| Global admin query QPS | 500 req/s | some cross-shard work is acceptable, but must be scoped |
| Rough cost | shard pool + dedicated heavy-tenant tier + migration tooling | the design must be operationally believable |

## Architecture

A strong drill answer typically includes:

1. shard by tenant or tenant bucket, not random config ID
2. tenant directory for flexible placement and moves
3. shard classes or dedicated pools for outsized tenants
4. bounded rebalance and reshard workflows with epochs
5. derived global views for admin dashboards

Useful sketch:

```text
client / control API
  -> auth + tenant identity
  -> directory lookup
  -> shard pool or dedicated tenant pool
  -> local indexes for tenant reads
  -> CDC to global dashboard / analytics views
```

## Data Model & APIs

Core entities:

- `TenantDirectory`
- `FlagConfig`
- `TenantClass`
- `MigrationCohort`
- `GlobalRollup`

Core APIs:

- `GetFlags(tenant_id, environment)`
- `UpdateFlag(tenant_id, flag_id, payload)`
- `MoveTenant(tenant_id, target_pool)`
- `QueryAdminDashboard(filters)`

The best answers keep write and read routing compatible during tenant moves and avoid pretending global dashboards can stay free.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one enterprise tenant overwhelms a shared shard | per-tenant usage dominates shard latency | dedicated pool or tenant split |
| stale routing sends writes to old owner during move | stale-epoch errors and parity mismatches | versioned directory and redirect logic |
| global dashboards query primaries directly | fanout and shard CPU rise with admin traffic | derived rollups or search-style index |
| migration copy steals serving headroom | p99 latency spikes during active moves | bounded concurrency and rate-limited backfill |

## Observability

- metric: per-tenant and per-shard traffic concentration
- metric: directory lookup latency and stale-epoch redirects
- metric: active move count, parity mismatch rate, and backfill lag
- metric: fanout width and freshness lag for global admin views
- SLO: tenant-local reads and config writes stay within target latency even during controlled migrations

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| tenant directory with flexible placement | easier moves and dedicated pools | extra routing hop | fully self-locating tenant hash with rigid placement |
| shard classes for large tenants | stronger fairness | more fleet management complexity | one pool for all customer sizes |
| derived admin views | protect serving path | freshness lag and pipeline cost | global scatter-gather on every admin request |

## Interview It

**Google framing:** "Design a global multi-tenant configuration store." Expect follow-ups on shard key, tenant movement, cross-shard admin queries, and migration safety.

**Cloudflare framing:** "Design customer configuration storage for an edge control plane." Expect pushback on noisy neighbors, regional placement, and zero-downtime movement.

**Follow-ups:**
1. What changes if large tenants demand region pinning?
2. How would you split one tenant that outgrows its current placement?
3. What if admin dashboards now need near-real-time counts?
4. How would you explain rollback during a partially completed migration?
5. Which parts of the system become more expensive at 10x tenant count?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-sharding-drill.md`
- `outputs/sharding-drill-sheet.md`

## Exercises

1. **Easy** — Deliver a 10-minute version of the drill with only the highest-signal decisions.
2. **Medium** — Re-run the drill assuming large tenants require dedicated regional placement.
3. **Hard** — Redesign the drill so cross-tenant admin analytics must be available within 2 seconds globally.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline for structuring a full interview answer
- [Google SRE books](https://sre.google/books/) — strong operational framing for rollout, reliability, and migration discipline
