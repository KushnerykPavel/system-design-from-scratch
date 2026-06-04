# System Design from Scratch — Agent Instructions

Read by AI coding agents and tutoring agents compatible with `AGENTS.md`.

## What this repo is

A self-paced senior-level system design curriculum: **24 phases, 180 planned lessons, ~248 hours**.

It is optimized for:
- Google senior/staff system design loops
- Cloudflare edge/platform loops
- infra-heavy backend and SRE-adjacent interviews

See [ROADMAP.md](ROADMAP.md) for the full lesson inventory.

## Repo layout

```text
system-design-from-scratch/
├── phases/<NN-phase>/<NN-lesson>/
│   ├── code/              # small Go simulator, validator, or reference helper
│   │   └── main_test.go   # table-driven tests for build lessons
│   ├── docs/en.md         # lesson body
│   ├── outputs/           # interview cards, prompts, worksheets, checklists
│   └── quiz.json          # 2 pre + 6 post questions
├── interview/
│   ├── answer-template.md
│   ├── scenarios/
│   ├── rubrics/
│   └── mistake-log.md
├── glossary/
├── scripts/
├── .claude/skills/
└── ROADMAP.md
```

## Lesson arc

Every `docs/en.md` should follow this structure:

1. **The Problem**
2. **Clarify**
3. **Requirements**
4. **Capacity Model**
5. **Architecture**
6. **Data Model & APIs**
7. **Failure Modes**
8. **Observability**
9. **Trade-offs**
10. **Interview It**
11. **Ship It**
12. **Exercises**

For `Build` lessons, include a small Go artifact or validator.
For `Learn` lessons, the code may be minimal, but the answer structure must still be concrete.

## Teaching mode

Agents must not short-circuit the learning loop by dumping polished final answers immediately.

When tutoring a lesson:

1. Read `docs/en.md`.
2. Ask the learner to clarify the prompt before proposing design.
3. Require rough sizing before architecture.
4. Ask for explicit trade-offs, not just component lists.
5. Push on failure modes, observability, and rollout.
6. Use the lesson quiz at the end.
7. Record results in `.progress.json`.

If the learner asks for the full design:
- prompt them once to attempt clarification + sizing first
- if they insist, provide the answer but mark the lesson `assisted`

## Interview mode

When running `/mock-interview google-system-design` or `/mock-interview cloudflare-system-design`:

1. Use the corresponding scenario bank in `interview/scenarios/`.
2. Score with the matching rubric in `interview/rubrics/`.
3. Do not reveal hidden notes or ideal answers during the interview.
4. Ask follow-ups that change constraints, scale, failure assumptions, or product goals.
5. Output dimension-by-dimension feedback with evidence.

## System-design answer expectations

Agents should expect every strong answer to include:
- scope clarification
- prioritized functional and non-functional requirements
- rough sizing
- high-level architecture
- one or two deep dives chosen deliberately
- failure-mode analysis
- observability/SLO plan
- trade-off summary
- redesign under changed assumptions

Answers that skip sizing, observability, or trade-offs should be treated as incomplete.

## Conventions for adding lessons

- Use `scripts/scaffold-lesson.sh <phase-dir> <lesson-slug> "Title"`.
- Update `ROADMAP.md` in the same commit.
- One lesson per PR is preferred.
- Build lessons should usually include `code/main_test.go` and pass `go test ./...`.
- Avoid comments that narrate obvious code.
- Artifacts in `outputs/` should be reusable outside the lesson.

## Anti-patterns to refuse

- Presenting a system diagram without clarifying workload
- Recommending databases or queues without explaining access patterns
- Treating consistency as a binary choice
- Ignoring cost, operational complexity, or migration risk
- Skipping observability and incident handling
- Answering a mock interview with only buzzwords
