---
name: mock-interview-cloudflare-system-design
version: 1.0.0
description: Cloudflare-style edge/platform system design mock interview.
---

# Mock Interview — Cloudflare System Design

Use scenarios from `interview/scenarios/cloudflare-system-design.md`.
Score using `interview/rubrics/cloudflare-edge-platform.md`.

Emphasize:
- traffic realism
- cache semantics
- origin protection
- abuse handling
- latency and egress trade-offs

Persist feedback using the shared structured payload:
- `session_type`: `mock_interview`
- strengths
- gaps
- one highest-leverage improvement
- dimension scores for `clarification`, `requirements`, `sizing`, `architecture`, `deep_dive`, `failure_modes`, `observability`, `trade_offs`, `communication`

Map edge-specific rubric findings into the shared dimensions so the learner can compare mock results with lesson and review sessions.
