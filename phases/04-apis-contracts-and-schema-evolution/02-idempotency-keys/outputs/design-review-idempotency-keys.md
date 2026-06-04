---
lesson: 02-idempotency-keys
focus: balanced
---

## Design review prompts

- Where is the idempotency record written relative to the business commit?
- How do you detect same key, different payload?
- What is returned for a successful duplicate retry?
- What happens during regional failover or gateway retry storms?
- Which metrics tell you the mechanism is actually preventing duplicates?
