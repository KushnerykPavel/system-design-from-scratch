# Observability Checklist — Metrics Platform

- Accepted, rejected, and dropped sample counts are broken down by reason and tenant.
- Active-series growth is visible by tenant, metric family, and top label keys.
- Alert evaluation lag is tracked separately from dashboard query latency.
- Write-ahead backlog, shard pressure, and ingest replication health are always on the operator dashboard.
- Downsampling backlog age is monitored so cold-tier freshness drift is visible.
- Cardinality rejections are logged with enough context for fast user follow-up.
- Query classes are isolated so exploratory reads cannot silently break paging quality.
