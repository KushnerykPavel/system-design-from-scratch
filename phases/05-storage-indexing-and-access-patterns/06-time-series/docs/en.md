# Time-Series and Append-Heavy Workloads

> Appends are easy until retention, compaction, and query windows arrive.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design storage for append-heavy and time-series data with explicit choices around partitioning by time, rollups, retention, and query-window constraints.  
**Prerequisites:** `02-access-pattern-first`, `04-hot-and-cold-data`  
**Estimated time:** ~75 min  
**Primary artifact:** capacity sheet + trade-off matrix  

## The Problem

Time-series workloads look simple because writes are mostly appends. The complexity appears later in partition sizing, recent-window query latency, rollups, downsampling, and retention.

This lesson trains you to reason about:

- append-heavy ingest
- bounded time-window reads
- retention and compaction
- hot recent partitions versus cold historical windows

## Clarify

- Are writes immutable appends, or can points be updated later?
- What are the most common read windows: 5 minutes, 1 hour, 7 days, 1 year?
- Do users need raw points, rollups, or both?
- What is the acceptable lag for aggregation and downsampling?

## Requirements

### Functional

- Ingest high-volume append events.
- Query bounded time ranges efficiently.
- Support rollups, retention, and backfill without destroying the serving path.

### Non-functional

- Keep recent-window queries fast.
- Control cost for long-term historical retention.
- Avoid one partition or shard becoming the write bottleneck.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Ingest rate | 2M points/s | drives append path and partition sizing |
| Cardinality | 50M active series/day | affects index strategy and memory pressure |
| Common query window | 15 min to 24 hr | shapes storage layout and cacheability |
| Retention | raw 7 days, rollups 13 months | determines compaction and cost plan |
| Rough cost | hot recent storage + rollup store + archive | append-heavy systems often need multiple retention tiers |

## Architecture

Typical shape:

```text
producers
  -> ingest gateway
  -> append log / write buffer
  -> recent hot partitions
  -> rollup / compaction jobs
  -> long-term tier
```

The key discipline is to keep queries bounded by tenant, metric, and time range.

## Data Model & APIs

Common key shape:

- `tenant_id + metric_id + time_bucket + timestamp`

Useful APIs:

- `WritePoints(batch)`
- `QueryRange(metric, start, end, step)`
- `CreateRollup(metric, window)`

A strong answer explicitly names:

- time partition size
- late-arriving data policy
- raw versus aggregated retention

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| unbounded query windows overwhelm storage | query timeout and scanned-bytes metrics spike | enforce max lookback and route large jobs async |
| recent partition becomes too hot | shard latency and queue depth rise | hash within recent buckets or pre-split heavy tenants |
| rollup lag grows until dashboards are stale | freshness metrics breach threshold | separate ingest from rollup compute and scale compaction workers |
| retention deletes raw data needed for audit | restore or query gaps surface later | define retention classes and protect critical streams separately |

## Observability

- metric: ingest throughput, queue depth, and reject rate
- metric: query scanned bytes per requested window
- metric: rollup lag and freshness by aggregation tier
- metric: partition skew for recent buckets
- log: rejected query windows and late data handling decisions
- SLO: recent dashboard queries meet latency target while ingest remains durable under peak write rate

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| time-bucketed partitions | efficient recent reads and retention | awkward cross-bucket queries | one giant append table |
| raw plus rollup tiers | cheap long-range queries and lower storage cost | compaction complexity and lag | serving year-long queries from raw points |
| bounded query windows | protects cluster from runaway scans | less analyst flexibility on serving path | unlimited ad hoc queries on the OLTP/TSDB path |

## Interview It

**Google framing:** "Design metrics storage for service dashboards." The signal is whether you discuss bounded queries, rollups, and retention instead of only append throughput.

**Cloudflare framing:** "Design high-volume edge telemetry ingestion." The signal is whether you separate ingest durability from recent reads and historical cost control.

**Follow-ups:**
1. What if one metric family is 30% of all writes?
2. What if users now want 13 months of hourly rollups?
3. How do you handle late-arriving points?
4. What if raw retention must increase during an incident?
5. When does a warehouse or offline analytics path become necessary?

## Ship It

- `outputs/capacity-sheet-time-series.md`
- `outputs/tradeoff-matrix-time-series.md`

## Exercises

1. **Easy** — Choose a key and partitioning scheme for host metrics.  
2. **Medium** — Explain raw versus rollup retention for API latency dashboards.  
3. **Hard** — Redesign a time-series system after one tenant begins dominating recent ingest.  

## Further Reading

- [The Datadog time series database](https://www.datadoghq.com/blog/engineering/timeseries-metric-storage-at-scale/) — practical discussion of scale and storage choices  
- [Prometheus storage docs](https://prometheus.io/docs/prometheus/latest/storage/) — helpful framing for append-heavy recent storage behavior  
