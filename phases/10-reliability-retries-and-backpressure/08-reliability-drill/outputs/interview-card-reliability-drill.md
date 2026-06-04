---
lesson: 08-reliability-drill
focus: balanced
---

## Must cover

- timeout and retry ownership
- idempotent ingest or write behavior
- bounded async backlog and retry policy
- admission control under overload
- isolation boundaries and blast radius
- operator-visible metrics and controls

## Weak answer smells

- queues used as infinite shock absorbers
- retries discussed with no budgets
- duplicate submits ignored
- no story for one bad tenant or one bad receiver

## Best follow-up pivots

- one dependency times out globally
- one tenant sends 100x more traffic
- backlog lasts long enough that stale work matters
