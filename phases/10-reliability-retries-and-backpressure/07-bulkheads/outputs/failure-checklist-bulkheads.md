---
name: bulkhead-failure-checklist
phase: 10
lesson: 07
---

- Critical and optional work do not share every pool.
- One large tenant cannot monopolize shared capacity silently.
- Healthy cells can remain healthy during local failure.
- Fallback paths do not depend on the same broken bottleneck.
- Capacity borrowing is bounded and revocable.
- Blast-radius boundaries are visible in dashboards and runbooks.
