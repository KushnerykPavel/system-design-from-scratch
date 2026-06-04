# Metrics Platform

> A metrics platform fails less often from raw ingest volume than from hidden cardinality explosions and unclear retention promises.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Design a metrics platform that treats cardinality, retention tiers, query isolation, and write-path durability as first-class design choices instead of implementation details.
**Prerequisites:** `05-storage-indexing-and-access-patterns/06-time-series`, `10-reliability-retries-and-backpressure/04-load-shedding`, `11-observability-slos-and-debugging/04-dashboards`
**Estimated time:** ~75 min
**Primary artifact:** retention-policy validator + observability checklist

## The Problem

Design a metrics platform that ingests time-series samples from many services, stores recent high-resolution data, supports dashboard and alert queries, and keeps costs under control as tenants add labels and new exporters.

This lesson matters because strong candidates know metrics systems are not "just append to storage." They talk about cardinality budgets, downsampling, write fan-in, query isolation, and why dashboards and alert evaluators should not share every bottleneck.

## Clarify

- Are we serving infrastructure metrics, product metrics, or both?
- What alert freshness is required for paging use cases?
- Are queries mostly short-range dashboards or long-range analytical scans?
- Do tenants control labels freely, or do we enforce schema and budgets?

If the interviewer leaves it open, assume infrastructure-heavy metrics, alert freshness in tens of seconds, dashboards on recent data, and stricter label governance than logs.

## Requirements

### Functional

- Ingest high-volume labeled time-series samples.
- Query recent and historical data with different resolution tiers.
- Power alert evaluation and dashboards.
- Enforce tenant or metric-family cardinality controls.
- Support downsampling and retention policies.

### Non-functional

- Keep ingest durable enough for operational use without making every write path too expensive.
- Prevent one tenant's cardinality explosion from harming everyone else.
- Isolate heavy ad hoc queries from alert evaluation.
- Make storage-cost trade-offs explicit through retention tiers.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Samples ingested | 50M samples/s peak | drives write path and batching design |
| Active series | 4B | cardinality and index structures dominate cost |
| Hot retention | 7 days raw | recent dashboards and alerts need full fidelity |
| Cold retention | 180 days downsampled | historical trends can trade fidelity for cost |
| Peak factor | 3x during incidents | telemetry spikes happen when the system is already stressed |

## Architecture

```text
agents / exporters
  -> ingest gateways
  -> validation + cardinality guardrails
  -> replicated write-ahead path
  -> hot TSDB shards
  -> downsampling jobs
  -> cold storage

dashboard queries / alert evaluators
  -> query routers
  -> recent-store readers + long-range readers
```

Design notes:

1. Put lightweight schema and cardinality checks near ingest before bad labels become permanent cost.
2. Split hot recent storage from colder historical tiers.
3. Keep alert evaluation on predictable query lanes, not mixed with every dashboard explorer query.
4. Downsampling is a product decision as much as a storage decision because it changes what users can infer historically.

## Data Model & APIs

Core dimensions:

```text
metric_name
label_set
timestamp
value
tenant_id
resolution_tier
retention_policy
```

Useful interfaces:

- `RemoteWrite(samples[])`
- `QueryRange(metric, labels, start, end, step)`
- `SetCardinalityBudget(metric_family, limit)`
- `PromoteRetentionPolicy(metric_family, raw_days, downsample_days)`

The answer is stronger when it explicitly distinguishes ingest validation from query-time interpretation.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| label cardinality explodes after a bad deploy | active-series growth and rejected-series metrics | budgets, allowlists, and per-tenant throttles |
| alert queries stall behind ad hoc dashboards | alert-evaluator lag and query-class saturation | isolated query pools and reserved capacity |
| hot shard fills during incident spike | shard ingestion lag and WAL backlog | shard splitting, backpressure, and selective sample dropping |
| downsampling jobs fall behind | downsample backlog age and cold-tier freshness drift | backlog prioritization and temporary raw-retention extension |

## Observability

- metric: samples accepted, rejected, and dropped by reason
- metric: active series count by tenant and metric family
- metric: alert evaluation lag and query latency by class
- metric: write-ahead backlog and hot shard disk pressure
- log: cardinality budget rejections with tenant, metric family, and label keys
- trace: remote-write ingest through validation, WAL, and shard commit
- SLO: alert evaluation reads are fresh enough to meet the paging target even during ingest spikes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| strict ingest budgets | protects shared platform | can frustrate teams during onboarding | unlimited labels with post-facto cleanup |
| hot/cold retention tiers | cost control with useful history | more query planning complexity | raw full-resolution retention forever |
| isolated alert query lanes | more predictable paging quality | some idle reserved capacity | fully shared query cluster |

## Interview It

**Google framing:** "Design an internal metrics platform for thousands of services." Expect follow-ups on label abuse, storage cost, and alerting correctness under incidents.

**Cloudflare framing:** "Design a multi-tenant metrics backend for edge and platform workloads." Expect questions on tenant isolation, regional ingest, and how alert evaluation stays trustworthy under load.

**Follow-ups:**
1. What happens when one team adds `request_id` as a label?
2. How would you serve global queries across regions?
3. What changes if SREs need one-year historical trends?
4. How do you degrade safely during an ingest storm?
5. When would you choose approximate or sampled metrics over exact series?

## Ship It

- `outputs/observability-checklist-metrics-platform.md`

## Exercises

1. **Easy** — Estimate the storage impact of keeping raw one-second samples for 180 days.
2. **Medium** — Design a guardrail to stop one tenant from exploding active-series count.
3. **Hard** — Redesign the platform so global alerts can evaluate across multiple regional ingest clusters.

## Further Reading

- [Prometheus storage](https://prometheus.io/docs/prometheus/latest/storage/) — helpful baseline for TSDB trade-offs
- [Site Reliability Engineering: Monitoring Distributed Systems](https://sre.google/sre-book/monitoring-distributed-systems/) — good framing for what metrics are actually for
