# Lesson Template — System Design from Scratch

Copy this structure for each new lesson.

## Folder structure

```text
NN-lesson-name/
├── code/
│   ├── main.go                 # small simulator, validator, or helper
│   └── main_test.go            # table-driven tests for build lessons
├── docs/
│   └── en.md
├── quiz.json
└── outputs/
    ├── design-review-*.md
    ├── interview-card-*.md
    ├── capacity-sheet-*.md
    ├── tradeoff-matrix-*.md
    ├── failure-checklist-*.md
    └── observability-checklist-*.md
```

## Documentation format (`docs/en.md`)

```markdown
# [Lesson Title]

> [One-line motto — the idea the learner should remember]

**Type:** Learn | Build
**Company focus:** Google | Cloudflare | Balanced
**Learning goal:** [one sentence]
**Prerequisites:** [prior lesson slugs]
**Estimated time:** ~[N] min
**Primary artifact:** [capacity sheet / review prompt / simulator / checklist]

## The Problem

[Concrete prompt and why this lesson matters.]

## Clarify

- [top clarification question]
- [top clarification question]
- [assumptions if interviewer does not answer]

## Requirements

### Functional

- [requirement]
- [requirement]

### Non-functional

- [latency / availability / consistency / compliance / cost]

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| QPS | | |
| Storage | | |
| Bandwidth | | |
| Peak factor | | |
| Rough cost | | |

## Architecture

[ASCII diagram or bullet narrative.]

## Data Model & APIs

[Key entities, ownership boundaries, API surface, idempotency/versioning notes.]

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
|         |           |            |

## Observability

- [metric]
- [log]
- [trace]
- [SLO / alert]

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
|          |         |      |                      |

## Interview It

**Google framing:** [prompt + likely pushback]

**Cloudflare framing:** [prompt + likely pushback]

**Follow-ups:**
1. [change scale]
2. [change consistency]
3. [change latency/cost]
4. [introduce failure]
5. [migration or rollout]

## Ship It

[List the reusable artifact(s) created in `outputs/`.]

## Exercises

1. **Easy** — [scoping or sizing variation]
2. **Medium** — [deep-dive or trade-off variation]
3. **Hard** — [redesign under new constraints]

## Further Reading

- [link] — [why it matters]
```

## Quiz format (`quiz.json`)

Use the locked schema:

```json
{
  "questions": [
    {
      "stage": "pre",
      "question": "Question text.",
      "options": ["A", "B", "C", "D"],
      "correct": 1,
      "explanation": "Why this is correct."
    }
  ]
}
```

Rules:
- 2 `pre` questions + 6 `post` questions
- options should have roughly similar length
- at least one question should probe trade-offs
- at least one question should probe failure handling or observability

## Go code conventions

- `code/main.go` should compile cleanly
- build lessons should usually have tests in `code/main_test.go`
- prefer small, sharp tools over big frameworks
- no comments unless the *why* is non-obvious

## Artifact suggestions

### Interview card

```markdown
---
lesson: NN-lesson-slug
focus: google|cloudflare|balanced
---

## Clarify first
- [...]

## Must-size numbers
- [...]

## Core design
- [...]

## Failure probes
- [...]

## Trade-off summary
- [...]
```

### Capacity sheet

```markdown
---
lesson: NN-lesson-slug
---

| Metric | Value | Notes |
|--------|-------|-------|
| QPS | | |
| Reads/day | | |
| Writes/day | | |
| Storage/year | | |
```
```
