# Mistake Log and Feedback Loop

> The fastest way to improve is to stop relearning the same lesson.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Build a mistake log that captures repeat failure patterns and turns them into deliberate practice loops.
**Prerequisites:** `00-setup-and-workflow/01-repo-setup-and-progress`, `00-setup-and-workflow/02-note-taking-system`
**Estimated time:** ~45 min
**Primary artifact:** mistake log template

## The Problem

Most candidates remember painful mocks but fail to encode them into a system. The result is familiar: every week they rediscover that they skipped sizing, ignored trade-offs, or went too deep on the wrong component.

This lesson turns feedback into a stable dataset rather than emotional residue.

## Clarify

- Which mistakes are one-off misses versus recurring habits?
- Do you want to classify mistakes by interview stage, technical concept, or communication pattern?
- What evidence would prove a logged mistake is actually improving over time?

## Requirements

### Functional

- Record the mistake, evidence, root cause, and corrective drill.
- Tag mistakes so repeated themes become visible.
- Link a mistake to the lesson, mock, or artifact where it appeared.

### Non-functional

- Logging must be quick enough to finish while the memory is fresh.
- Categories should stay stable across many sessions.
- The feedback loop should emphasize behavior change, not self-criticism.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Mistakes logged | dozens over a full course run | needs searchable categories |
| Update frequency | after lessons and mocks | encourages simple markdown or JSON records |
| Retention | keep all records | trend visibility matters more than storage |
| Peak factor | clustered around mock interviews | template must be fast to fill |
| Rough cost | negligible | optimize human attention instead |

## Architecture

Use a simple loop:

1. Capture evidence from the session.
2. Classify the mistake.
3. Name the root cause.
4. Assign a corrective action.
5. Revisit the issue after the next similar prompt.

Good mistake categories:

- clarification
- sizing
- architecture decomposition
- deep-dive choice
- consistency reasoning
- failure mode coverage
- observability
- trade-off communication

## Data Model & APIs

Recommended fields:

- `date`
- `prompt`
- `lesson_or_mock`
- `category`
- `symptom`
- `root_cause`
- `corrective_drill`
- `retested`
- `status`

Queries you should be able to answer:

- What are my top three recurring categories?
- Which mistakes have not been retested?
- Which failure modes show up only in company-specific mocks?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| log becomes a wall of shame | entries describe frustration but no corrective action | require a concrete next drill for each entry |
| categories sprawl | every entry invents a new label | keep a controlled category list |
| feedback never closes | mistakes are logged but not retested | track retest date and status |
| shallow root causes | entries stop at surface symptoms | force a `why did this happen` field |

## Observability

- metric: recurring mistakes by category over rolling four weeks
- metric: percentage of mistakes with a completed corrective drill
- metric: time from logging a mistake to retesting it
- SLO: every mock should produce either zero findings or explicit logged findings within 24 hours

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| structured categories | easier trend analysis | slightly less expressive | free-form journaling |
| mandatory corrective action | improves follow-through | more effort per entry | passive note collection |
| keep historical mistakes | preserves trend visibility | old records may feel noisy | deleting resolved issues |

## Interview It

**Google framing:** "How would you design a practice feedback loop so repeated system design mistakes actually disappear over time?"

**Cloudflare framing:** "How would you track operational and communication mistakes across multiple infrastructure mock interviews so they become actionable?"

**Follow-ups:**
1. What if the same mistake appears across unrelated prompts?
2. How should the system treat a mistake that is fixed in solo practice but returns in live mocks?
3. What categories would you use for edge-specific prompts?
4. How do you keep the log honest without making it exhausting to maintain?
5. What is the minimum data you would still capture after a bad session?

## Ship It

- `outputs/mistake-log-template.md`
- `outputs/design-review-feedback-loop.md`

## Exercises

1. **Easy** — Write two example mistakes from a recent design session and classify them.
2. **Medium** — Define a retest rule for sizing mistakes versus communication mistakes.
3. **Hard** — Design a weekly review that selects the next three drills from the mistake log automatically.

## Further Reading

- [Thinking in Bets](https://www.penguinrandomhouse.com/books/547787/thinking-in-bets-by-annie-duke/) — useful framing for feedback under uncertainty
- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — strong examples of using post-incident learning loops productively
