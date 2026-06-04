# Saga Design Review

## Workflow map

| Step | Local transaction owner | Idempotency key | Retry policy | Compensation or reconciliation |
|------|-------------------------|-----------------|--------------|-------------------------------|
|      |                         |                 |              |                               |
|      |                         |                 |              |                               |
|      |                         |                 |              |                               |

## Ask directly

- What state can users observe between steps?
- Which step is hardest to reverse?
- Where does workflow state live durably?
- How will operators resume or manually repair a stuck workflow?

## Review the hard parts

- irreversible side effects
- duplicate external calls
- long-running timeouts
- partial rollback
- progress visibility
