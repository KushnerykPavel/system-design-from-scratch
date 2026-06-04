---
lesson: 02-idempotency-under-failure
focus: balanced
---

## Design review prompts

- What is the exact scope of the idempotency key?
- Can the same key arrive with a different payload, and what happens then?
- Is the key reserved before the side effect starts?
- How long must duplicate suppression last in realistic incidents?
- Does the system return a stable prior result or only a generic duplicate response?
- What happens to `in_progress` records after crashes?
