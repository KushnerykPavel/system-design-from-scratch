---
lesson: 02-idempotency-keys
focus: balanced
---

## Duplicate-safety checks

- Does the client send a stable idempotency key across retries?
- Is the request hash stored and compared?
- Can the service replay the original outcome?
- What happens if the business write commits before the response returns?
- How are pending records repaired or expired?

## Review questions

- What TTL matches real retry behavior?
- Which writes truly need idempotency and which do not?
- Is key scope per user, per tenant, or global?
