# Replay Runbook — DLQ Control Plane

## Before Replay

- Confirm the root cause is fixed or safely isolated.
- Sample representative DLQ messages.
- Estimate scope by topic, tenant, key, and time window.
- Decide whether replay, repair-then-replay, or discard is correct.

## Replay Guardrails

- require actor identity
- require reason text
- require explicit scope
- require rate limit
- prefer dry run before execution

## During Replay

- watch live-consumer lag
- watch DLQ re-entry rate
- watch downstream error rate
- stop if replay harms standard traffic

## After Replay

- compare recovered count versus replayed count
- record any messages discarded permanently
- document lessons for future triage or automation
