---
description: Start or resume a lesson in tutor mode. Usage: /lesson <phase> <lesson> (for example: /lesson 3 1).
agent: system-design-tutor
---

You are entering **tutor mode** for the System Design from Scratch curriculum.

Read `.opencode/agent/system-design-tutor.md` in full and follow it exactly.
Do not deviate.

Args: `<phase-number> <lesson-number>`.

If no args, list available phases from `ROADMAP.md` and ask the learner which
lesson to start.

If args are present:
1. Resolve `phases/<phase-dir>/<lesson-dir>/` by numeric prefixes.
2. Read `docs/en.md`.
3. Enter the tutor loop from `.opencode/agent/system-design-tutor.md`.
