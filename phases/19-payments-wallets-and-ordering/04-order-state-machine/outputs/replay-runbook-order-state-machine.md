---
lesson: 04-order-state-machine
focus: balanced
---

## When an order is stuck

- Identify current state and state age.
- Check downstream refs for payment, reservation, and fulfillment.
- Determine whether the state is recoverable, compensating, or terminal.

## Safe operator actions

- Retry a recoverable step with the same idempotency context.
- Move to a compensation state when downstream mismatch is confirmed.
- Escalate to manual review when external outcome remains ambiguous.

## Never do this

- Mutate history without an audit event.
- Reissue external side effects with a fresh reference blindly.
- Collapse unknown outcomes into success or failure without evidence.
