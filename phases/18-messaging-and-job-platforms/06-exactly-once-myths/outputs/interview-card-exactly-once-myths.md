# Interview Card — Exactly-Once Myths

## Strong Opening

"I would avoid claiming end-to-end exactly-once without defining the boundary. The practical goal is usually exactly-once state transition inside a bounded store, with idempotent handling for retries and side effects beyond that boundary."

## Must Cover

- where the guarantee starts and ends
- how duplicates are detected or tolerated
- how state update and progress commit relate
- what replay does
- what residual risk remains

## Common Mistakes

- saying "exactly once" without scope
- ignoring replay and long-delay duplicates
- forgetting side effects like webhooks, email, or billing
- omitting observability for duplicate attempts
