---
name: load-shedding-checklist
phase: 10
lesson: 04
---

- Safe concurrency and queue-age limits are explicit.
- Critical traffic has a protected path.
- Reject reasons are visible to operators and callers.
- Downstream saturation can influence local admission.
- Client retry behavior for `429` or overload responses is considered.
- Queue growth cannot hide overload for long periods.
