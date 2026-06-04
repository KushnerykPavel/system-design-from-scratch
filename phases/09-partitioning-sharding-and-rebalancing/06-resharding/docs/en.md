# Resharding and Data Migration Plans

> When the original shard shape no longer fits, the migration design becomes part of the product.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a resharding plan that changes key or partition boundaries safely, with explicit phases for backfill, validation, cutover, and rollback.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/04-rebalancing`, `09-partitioning-sharding-and-rebalancing/05-tenant-isolation`
**Estimated time:** ~75 min
**Primary artifact:** resharding planner + migration checklist

## The Problem

Rebalancing moves existing shards. Resharding changes the shard map itself:

- doubling shard count
- splitting large tenants
- changing from one key shape to another
- adding regional prefixes or bucketization

That makes it riskier than ordinary movement because clients, indexes, routing rules, and background jobs may all need to understand old and new layouts at the same time.

## Clarify

- Are we only splitting ranges, or also changing the shard key?
- Can old and new layouts coexist temporarily?
- Which callers are layout-aware and need migration compatibility?
- Is the biggest risk stale writes, dual-read cost, index rebuild time, or customer-visible latency?

If the prompt is vague, assume a live migration from `tenant_id` shards to `region + tenant_bucket` with active traffic and no full outage allowed.

## Requirements

### Functional

- Build the new shard map and backfill data into it.
- Keep writes correct while old and new layouts coexist.
- Cut over readers and writers incrementally with rollback support.

### Non-functional

- Avoid global stop-the-world rewrites.
- Bound duplicate storage and migration traffic.
- Keep customer-visible latency and correctness stable throughout the plan.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Existing shard count | 128 shards | large enough that whole-fleet coordination matters |
| New shard count | 256 shards | split plan must stay incremental |
| Data to duplicate during migration | 140 TB | dual-layout cost is material |
| Active writes during cutover | 30K writes/s | dual-write or change replay is unavoidable |
| Rough cost | temporary duplicate storage + index rebuild + operator time | resharding is both a technical and budget event |

## Architecture

A safe resharding plan often has four phases:

1. **Prepare**
   Create new shard map and routing compatibility.
2. **Backfill**
   Copy old data into new destinations.
3. **Dual phase**
   Dual-write or replay changes while validating parity.
4. **Cutover**
   Shift reads and writes by cohort, then retire the old layout.

The key distinction from simple rebalancing:

- clients may need to resolve two layouts
- secondary indexes may need rebuilding
- background jobs need idempotent handling across layouts

## Data Model & APIs

Useful migration metadata:

- `ShardMap(version, key_schema, state)`
- `MigrationCohort(id, tenant_range, read_cutover, write_cutover)`
- `ParityCheck(cohort_id, sample_rate, mismatch_count)`

Useful interfaces:

- `StartBackfill(cohort_id)`
- `EnableDualWrite(cohort_id)`
- `ShiftReadTraffic(cohort_id, percent)`
- `FinalizeShardMap(version)`

Senior-level notes:

- use cohorts to avoid all-at-once cutover
- if the key changes, old IDs may need translation or lookup indirection
- rebuild or validate secondary indexes explicitly; they do not stay correct by accident

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| readers and writers disagree on active shard map | parity errors or stale-layout misses rise | version pinning and compatibility routing |
| backfill completes but derived indexes are incomplete | query mismatches on new layout | explicit index rebuild checks before cutover |
| dual-write path diverges silently | mismatch counts or missing CDC offsets appear | idempotent writes and parity sampling |
| cutover cohort too large | latency regression and rollback complexity increase | smaller cohorts and percentage-based traffic shifting |

## Observability

- metric: backfill progress, dual-write lag, and parity mismatch rate by cohort
- metric: read and write traffic split between old and new layouts
- metric: shard-map version adoption across clients and services
- log: cohort transitions, cutover decisions, and rollback reasons
- trace: request path including old-layout lookup, new-layout lookup, and translation steps
- SLO: migration must preserve correctness while keeping latency regressions within a defined bound

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| dual-layout compatibility | safer incremental cutover | more code and temporary complexity | global big-bang switch |
| cohort-based migration | smaller blast radius | slower full completion | all tenants at once |
| parity validation before final cutover | stronger correctness confidence | extra read cost and tooling | trusting copy completion alone |

## Interview It

**Google framing:** "Your original sharding strategy no longer works. Migrate without a full outage." The signal is whether you produce a phased plan with compatibility, not just 'copy data and flip.'

**Cloudflare framing:** "Reshard a global control-plane dataset while configuration requests keep flowing." The signal is whether you talk about staged rollout, compatibility, and rollback under active change.

**Follow-ups:**
1. What if one client version cannot understand the new shard map yet?
2. How do you migrate secondary indexes that depend on the old key shape?
3. What if parity mismatches stay low but never reach zero?
4. When is a big-bang cutover actually acceptable?
5. How do you keep analytics and background jobs from double-counting during dual-write?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-resharding.md`
- `outputs/migration-checklist-resharding.md`

## Exercises

1. **Easy** — Sketch the state machine for one migration cohort.
2. **Medium** — Explain how you would cut over readers before writers, or vice versa.
3. **Hard** — Design a reshard from global tenant IDs to regional tenant buckets with secondary indexes and background jobs still active.

## Further Reading

- [Spanner schema design and migration docs](https://cloud.google.com/spanner/docs/schema-updates) — helpful for phased compatibility thinking
- [Google SRE books](https://sre.google/books/) — rollout discipline and correctness-focused change management
