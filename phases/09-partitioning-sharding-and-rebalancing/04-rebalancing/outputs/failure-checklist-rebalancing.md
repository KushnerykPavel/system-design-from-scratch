---
name: rebalancing-failure-checklist
phase: 09
lesson: 04
---

Before cutover, confirm:

1. Source and target agree on data checksums or lag bounds.
2. Routing has an epoch or versioned ownership check.
3. Copy bandwidth is capped relative to spare serving capacity.
4. Pause and rollback commands are real, tested, and fast.
5. Alerting links active move count to serving latency and errors.
