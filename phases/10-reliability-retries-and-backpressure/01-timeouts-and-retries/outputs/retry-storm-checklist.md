---
name: retry-storm-checklist
phase: 10
lesson: 01
---

1. Deadlines are propagated end to end.
2. Only one layer owns most retry behavior.
3. Retryable failure classes are explicit.
4. Backoff includes jitter.
5. Attempt count is capped by remaining budget.
6. Attempt-per-request ratio is on dashboards.
7. Overload responses do not accidentally trigger blind immediate retries.
