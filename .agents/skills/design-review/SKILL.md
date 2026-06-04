---
name: design-review
version: 1.0.0
description: Reviews a learner's architecture answer against the lesson's design criteria.
---

# Design Review

Activation:
- `/design-review <phase> <lesson>`

Review dimensions:
- clarification
- requirements
- sizing
- architecture
- deep_dive
- failure_modes
- observability
- trade_offs
- communication

Output:
- `session_type`: `design_review`
- strengths
- gaps
- one highest-leverage improvement
- per-dimension scores from 1 to 4
- evidence for each weak dimension
