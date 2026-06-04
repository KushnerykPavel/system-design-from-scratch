---
lesson: 02-digital-wallet
focus: balanced
---

## Hold lifecycle checks

- Can you create, settle, release, and expire holds safely?
- Are partial settlements modeled explicitly?
- Are retries scoped to the specific hold action?

## Failure probes

- What happens if the order service times out after hold creation?
- How are aged holds surfaced to operators?
- How is overspend prevented during concurrent debits?

## Observability prompts

- Track hold age percentiles.
- Track insufficient-funds rejection rate.
- Track manual overrides and forced releases.
