---
lesson: 06-time-series
---

| Metric | Value | Notes |
|--------|-------|-------|
| Ingest | 2M points/s | append-heavy hot path |
| Active series | 50M/day | cardinality pressure |
| Query window | 15 min to 24 hr | bounded interactive reads |
| Raw retention | 7 days | recent detail only |
| Rollup retention | 13 months | cheaper long-range queries |
