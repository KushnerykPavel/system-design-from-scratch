---
lesson: 07-queues-vs-streams
focus: balanced
---

# Messaging Decision Matrix

| Need | Queue | Stream | Workflow Engine |
|------|-------|--------|-----------------|
| One logical worker owns execution | Strong fit | Possible but awkward | Overkill unless multi-step |
| Multiple independent consumers | Weak without fan-out add-ons | Strong fit | Only if business state must be tracked |
| Replay from durable history | Limited | Strong fit | Partial and workflow-specific |
| Long-running timers and waits | Weak | Weak | Strong fit |
| User-visible step status | Weak | Weak | Strong fit |
| Simple operational model | Strongest | Moderate | Most complex |

## Clarify First

- Who owns the work after publish?
- How many consumers need the same data?
- Is replay required for correctness, analytics, or backfill?
- Does the business process have explicit state transitions and timers?

## Common Mistakes

- Using a queue as a poor man's event log
- Using a stream as a poor man's workflow engine
- Claiming "exactly-once" without defining the side effect boundary
- Confusing transport choice with business ownership
