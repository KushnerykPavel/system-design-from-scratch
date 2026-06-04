---
name: load-shedding-tradeoff-matrix
phase: 10
lesson: 04
---

| Choice | Benefit | Cost | Best fit |
|--------|---------|------|----------|
| fast reject | protects latency | visible dropped work | low-latency critical APIs |
| bounded queue | absorbs short spikes | stale work if too loose | small transient bursts |
| reserved priority capacity | preserves core flows | lower average utilization | mixed-priority traffic |
| fairness by tenant | limits abuse blast radius | more policy state | multi-tenant shared services |
