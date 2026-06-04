---
name: idempotency-failure-checklist
phase: 10
lesson: 02
---

- Duplicate submit after ambiguous timeout does not double-apply the write.
- Key reuse with different payload is rejected or surfaced clearly.
- Dedupe storage survives restart and failover.
- TTL covers realistic retry windows.
- Response replay behavior is defined if clients need stable results.
- Stuck `in_progress` intents are measurable and recoverable.
