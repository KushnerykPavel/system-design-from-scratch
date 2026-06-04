---
lesson: 07-delivery-semantics
focus: balanced
---

# Interview Card: Delivery Semantics

## Clarify First

- What failure matters more: loss or duplication?
- What is the side effect boundary?
- Can handlers be idempotent?
- How long must replay or retries remain safe?

## Strong Default

- At-least-once delivery
- Idempotent consumers
- Stable message IDs on retry
- Explicit reconciliation for external effects

## Push Back On

- "Exactly-once everywhere"
- Ack before durable side effect
- No dedupe retention policy
- No observability for duplicates or loss
