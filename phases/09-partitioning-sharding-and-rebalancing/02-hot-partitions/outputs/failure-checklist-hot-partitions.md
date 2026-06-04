---
name: hot-partitions-failure-checklist
phase: 09
lesson: 02
---

Check these before proposing a fix:

1. Is the heat caused by reads, writes, retries, locks, or background jobs?
2. Is the concentration by key, tenant, region, or time bucket?
3. Are neighboring tenants sharing the same failure domain?
4. Will more shards actually spread the hot owner, or only add idle capacity?
5. Which metric will confirm the mitigation worked within minutes?
