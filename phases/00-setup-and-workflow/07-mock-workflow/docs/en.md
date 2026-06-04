# Mock Interview Workflow

> A mock is only useful if the workflow surfaces real signal.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Run mocks with enough structure to simulate interview pressure, preserve evidence, and produce actionable feedback instead of loose conversation.
**Prerequisites:** `00-setup-and-workflow/03-fast-diagramming`, `00-setup-and-workflow/06-review-checklist`
**Estimated time:** ~60 min
**Primary artifact:** mock workflow checklist + session planner

## The Problem

Unstructured mocks often fail in predictable ways: the interviewer helps too much, timing drifts, feedback arrives late, or nobody records the actual weak point. That means an hour of practice can produce less signal than a disciplined 25-minute drill.

This lesson gives you a repeatable workflow for solo, paired, and interviewer-led practice.

## Clarify

- Is the mock focused on process, technical depth, or company-specific style?
- Who controls time and constraint changes during the session?
- What evidence will be captured: notes, score, mistake log entries, or recording timestamps?

## Requirements

### Functional

- Define pre-brief, live mock, feedback, and debrief stages.
- Reserve time for follow-up constraint changes.
- Produce concrete artifacts after the session.

### Non-functional

- Workflow must stay usable for both 20-minute drills and 45-minute mocks.
- The interviewer role should be consistent enough to compare sessions.
- Feedback quality should not depend on perfect memory.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Mock duration | 20 to 45 minutes live | drives agenda budgeting |
| Feedback duration | 10 to 20 minutes | must be protected, not squeezed out |
| Artifacts per mock | 2 to 4 notes or score objects | small enough for low-friction logging |
| Peak factor | highest near interview loops | workflow should scale without extra coordination cost |
| Rough cost | low monetary cost, high attention cost | structure is the main leverage point |

## Architecture

Use a four-stage workflow:

1. **Pre-brief** — choose prompt, focus area, and timing.
2. **Live design** — run the mock with realistic interruptions and follow-ups.
3. **Feedback** — review with evidence and the checklist.
4. **Debrief** — log mistakes, next drills, and next lesson.

The code artifact validates whether a proposed mock agenda fits the intended duration and includes the required stages.

## Data Model & APIs

Session fields:

- `mode` such as `solo`, `peer`, `interviewer`
- `duration_minutes`
- `prompt`
- `company_focus`
- `stages`
- `review_completed`
- `mistakes_logged`

Useful API questions:

- does this agenda fit the session budget?
- does it include a feedback stage?
- how much time is reserved for changed constraints?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| mock turns into collaborative design | interviewer feeds missing ideas too early | define intervention rules before starting |
| timing drifts | no wrap-up or follow-up stage reached | enforce stage budget upfront |
| feedback is vague | summary lacks evidence or next action | require checklist-backed findings |
| nothing changes afterward | no mistake log or next lesson update | make debrief mandatory |

## Observability

- metric: percentage of mocks that reach both feedback and debrief
- metric: time spent per stage versus plan
- metric: count of concrete next actions per mock
- SLO: every mock should end with at least one logged finding and one scheduled follow-up

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| fixed stage agenda | comparable practice sessions | less spontaneity | purely free-form mock |
| evidence-based feedback | higher improvement quality | requires a little more discipline | vibes-only conversation |
| mandatory debrief | closes the learning loop | adds a final step when tired | stopping after verbal feedback |

## Interview It

**Google framing:** "How would you design your own mock interview process so it tests prioritization, sizing, and adaptability rather than just component recall?"

**Cloudflare framing:** "How would you run infrastructure-heavy mocks that preserve time for operational follow-ups and changed assumptions?"

**Follow-ups:**
1. What changes for a solo mock where you play both roles?
2. How should a peer interviewer decide when to interrupt?
3. What if the mock must target Cloudflare edge topics specifically?
4. How do you preserve evidence if the mock is done on a whiteboard only?
5. Which stage would you compress first in a 20-minute drill?

## Ship It

- `outputs/mock-workflow-checklist.md`
- `outputs/interview-card-mock-workflow.md`

## Exercises

1. **Easy** — Plan a 20-minute mock with one redesign follow-up.
2. **Medium** — Create different agendas for a Google-style and Cloudflare-style mock.
3. **Hard** — Design a rotating three-week mock cadence for two partners with different weak areas.

## Further Reading

- [Google Interview Warmup resources](https://www.google.com/about/careers/applications/interview-tips/) — useful framing around structured interview behavior
- [The SRE Workbook](https://sre.google/workbook/table-of-contents/) — strong example of disciplined operational review loops
