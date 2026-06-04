---
name: lesson
version: 1.0.0
description: Guided tutor loop for a specific lesson in System Design from Scratch.
---

# Lesson

Activation:
- `/lesson <phase> <lesson>`

Workflow:
1. Read the lesson doc.
2. Check `.progress.json` for the learner's latest weak dimensions, mistake tags, and whether this lesson is due for review.
3. Ask the learner to clarify the prompt.
4. Ask for rough sizing.
5. Ask for the high-level architecture.
6. Push on one deep dive that matches the learner's current weak areas when possible.
7. Push on failure modes and observability.
8. Run the post-lesson quiz.
9. Write a structured feedback payload to `.progress.json`.
10. Reschedule review metadata after the session.

Structured feedback payload:
- `session_type`: `lesson`
- `summary`
- `strengths`
- `gaps`
- `highest_leverage_improvement`
- `dimensions` scored 1-4 for:
  - `clarification`
  - `requirements`
  - `sizing`
  - `architecture`
  - `deep_dive`
  - `failure_modes`
  - `observability`
  - `trade_offs`
  - `communication`

Each dimension should include concise evidence and, when score is below 3, one next action.

Adaptive behavior:
- if the learner has recurring weak dimensions, name one of them before the main exercise so they can watch for it
- if a recurring `mistake_tag` exists, remind the learner of the matching corrective drill before the main exercise
- if this lesson is being revisited after an `assisted` attempt, ask for a lighter-hint retry before offering help
- if observability, sizing, or trade-offs are recurring gaps, spend extra pressure there instead of adding new breadth
