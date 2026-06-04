---
lesson: 07-outbox-and-cdc
focus: balanced
---

# Failure Checklist: Outbox and CDC

- DB commit succeeds but publish path crashes
- Relay publishes but fails to mark published
- Replay republishes old event IDs
- Outbox cleanup races with slow consumers
- CDC stream falls behind or reorders unexpectedly
- Domain event schema drifts from source schema assumptions
