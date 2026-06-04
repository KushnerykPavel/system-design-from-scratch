---
name: check-understanding
version: 1.0.0
description: Phase quiz for System Design from Scratch. Tests whether the learner can explain the phase's trade-offs and failure modes, not just recite terms.
---

# Check Understanding

Quiz the learner on one phase using `quiz.json` files from the phase lessons.

Prioritize questions that probe:
- trade-offs
- failure modes
- observability
- sizing impact

Adaptive behavior:
- prefer lessons that are due or overdue for review over untouched lessons
- if the learner has repeated weak dimensions, bias question selection toward them
- if the learner has a recurring `mistake_tag`, bias one question toward the matching corrective drill
- if the learner recently improved on a dimension, include one confirmation question before moving on

After grading, write a structured feedback payload to `.progress.json` with:
- `session_type`: `check_understanding`
- one short summary
- strengths
- gaps
- one highest-leverage improvement
- 1-4 scores for the shared dimensions:
  - `clarification`
  - `requirements`
  - `sizing`
  - `architecture`
  - `deep_dive`
  - `failure_modes`
  - `observability`
  - `trade_offs`
  - `communication`

Then reschedule the lesson review:
- low scores or `assisted` outcomes should come back in about `1` day
- medium scores should come back in about `3-7` days
- high scores should come back in about `7-14` days
- repeated misses should reset the interval and increment lapses

## Verdicts

- `8/8`: mastered
- `6-7`: solid
- `4-5`: shaky
- `0-3`: revisit
