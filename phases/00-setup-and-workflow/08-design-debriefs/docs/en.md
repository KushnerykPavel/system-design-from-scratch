# Design Debriefs and Iteration Cadence

> Improvement compounds when reflection has a schedule.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Turn each lesson or mock into a short debrief and a deliberate next-step cadence so practice compounds instead of fragmenting.
**Prerequisites:** `00-setup-and-workflow/04-mistake-log`, `00-setup-and-workflow/07-mock-workflow`
**Estimated time:** ~45 min
**Primary artifact:** debrief template + iteration cadence worksheet

## The Problem

Practice without a debrief creates motion without adaptation. Even when people do reflect, they often do it informally and inconsistently, so the next session starts from mood rather than evidence.

This lesson gives you a closing ritual: what changed, what remains weak, and what to do next.

## Clarify

- What should happen after every lesson versus after every mock?
- Which signals matter enough to influence the next week of practice?
- How long can the cadence be before weak areas start to drift out of focus?

## Requirements

### Functional

- Capture strengths, gaps, and the next deliberate drill.
- Schedule when the topic will be revisited.
- Connect debrief output back to the progress tracker and mistake log.

### Non-functional

- Debrief should take 5 to 10 minutes, not another full study block.
- Cadence should be realistic for part-time learners.
- The system must prevent favorite topics from crowding out weak areas.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Debriefs per week | several during active prep | workflow must stay lightweight |
| Follow-up horizon | 1 day to 2 weeks | enough to reinforce learning without overload |
| Active weak areas | 3 to 6 at a time | helps cap the next-action queue |
| Peak factor | heavy after mock-heavy weeks | cadence must protect recovery and focus |
| Rough cost | low | scheduling attention is the real resource |

## Architecture

Use a simple loop:

1. Finish the lesson or mock.
2. Record one strength and one highest-priority gap.
3. Link any recurring issue to the mistake log.
4. Schedule the next targeted drill.
5. Review cadence weekly and rebalance weak areas.

This creates a portfolio view of your prep rather than a pile of disconnected sessions.

## Data Model & APIs

Suggested debrief fields:

- `session_type`
- `prompt_or_lesson`
- `strength`
- `gap`
- `next_drill`
- `due_date`
- `linked_mistake_categories`

Useful review questions:

- What is overdue for retest?
- Which weak areas are improving?
- Which topics are being ignored because they feel uncomfortable?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| debrief skipped when tired | missing closeout after hard sessions | use a five-minute minimum debrief template |
| too many next actions | backlog grows and none get done | cap active follow-ups to a small set |
| weak areas never revisited | gaps repeat but due dates are absent | schedule a retest date in every debrief |
| cadence becomes guilt system | backlog only grows | prune aggressively and focus on the highest-leverage gaps |

## Observability

- metric: percentage of sessions with completed debriefs
- metric: average days until a logged weak area is retested
- metric: count of overdue drills
- SLO: every substantive mock should produce a scheduled next step before the day ends

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| mandatory short debrief | preserves learning signal | adds a little end-of-session effort | optional reflection |
| limited active follow-ups | improves completion | some low-priority gaps wait longer | chasing every issue immediately |
| weekly cadence review | keeps portfolio balanced | requires recurring discipline | ad hoc next-lesson choice |

## Interview It

**Google framing:** "How would you build a post-practice workflow that ensures each system design session influences the next one?"

**Cloudflare framing:** "How would you debrief edge and reliability-focused mocks so operational gaps are revisited on purpose rather than forgotten?"

**Follow-ups:**
1. How does the cadence change when interviews are two weeks away?
2. Which sessions deserve a full debrief versus a lightweight one?
3. How should you prioritize a new weak area against an overdue old one?
4. What if your debrief says the architecture was fine but communication failed?
5. How would you review improvement across a full phase?

## Ship It

- `outputs/design-debrief-template.md`
- `outputs/iteration-cadence-worksheet.md`

## Exercises

1. **Easy** — Write a debrief for a recent lesson and schedule one follow-up drill.
2. **Medium** — Design a weekly review ritual for three active weak areas.
3. **Hard** — Build a two-week cadence for mixed Google and Cloudflare mock prep.

## Further Reading

- [Atomic Habits](https://jamesclear.com/atomic-habits) — useful for designing low-friction review loops
- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — good operational examples of feedback loops and prioritized follow-through
