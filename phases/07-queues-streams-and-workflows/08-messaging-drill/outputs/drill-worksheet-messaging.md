---
lesson: 07-messaging-drill
focus: balanced
---

# Messaging Drill Worksheet

## Clarify First

- Is this one-owner work, shared event history, or long-running orchestration?
- What must stay ordered together?
- What failure is worse: loss, duplication, or delay?
- What backlog age becomes a product incident?

## Must-Choose Decisions

- Primitive: queue, stream, workflow, or mix
- Delivery semantics for critical path
- Replay and DLQ policy
- Backpressure policy

## Close With

- Key failure mode
- Key observability signal
- Main trade-off
