---
name: sharding-drill-sheet
phase: 09
lesson: 08
---

Prompt:

Design the storage and serving architecture for a multi-tenant feature flag platform with millions of tenants, uneven enterprise customers, and a requirement to move large tenants safely over time.

Checklist:

1. Clarify workload locality and tenant skew.
2. Choose a shard key and defend it.
3. Explain hotspot and noisy-neighbor controls.
4. Describe placement, rebalancing, and resharding paths.
5. Bound or derive cross-shard admin queries.
6. Close with observability and trade-offs.
