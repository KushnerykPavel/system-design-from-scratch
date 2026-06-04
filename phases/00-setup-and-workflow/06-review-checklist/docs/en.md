# Architecture Review Checklist

> Review answers for missing reasoning, not missing buzzwords.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Use a short, high-signal checklist to review whether a design answer covered the dimensions senior interviewers actually care about.
**Prerequisites:** `00-setup-and-workflow/03-fast-diagramming`, `00-setup-and-workflow/05-capacity-sheets`
**Estimated time:** ~60 min
**Primary artifact:** design review checklist

## The Problem

Without a review standard, feedback becomes vague: "good answer," "needs more detail," or "should think about scale." That kind of feedback does not help. A senior-level review checklist should reveal what was missing, why it mattered, and what to improve next.

## Clarify

- Is the checklist for self-review, peer review, or interviewer-style scoring?
- Which missing sections most often explain weak answers: scope, sizing, failure modes, or trade-offs?
- How much review depth is appropriate after a 20-minute drill versus a full 45-minute mock?

## Requirements

### Functional

- Review clarification, requirements, sizing, architecture, deep dives, failure modes, and wrap-up.
- Make it easy to record evidence rather than vibes.
- Produce one or two concrete next actions.

### Non-functional

- Full review should take under ten minutes.
- The checklist must work across product and infrastructure prompts.
- Signal should favor completeness and prioritization over stylistic preferences.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Review items | 8 to 12 | enough coverage without overwhelming the reviewer |
| Reviews per week | several during active prep | favors concise prompts and stable categories |
| Evidence notes | 1 to 3 lines per item | keeps reviews actionable |
| Peak factor | highest after mocks | review flow must stay fast while memory is fresh |
| Rough cost | low | time is the dominant resource |

## Architecture

Organize the checklist by answer stages:

1. Clarify and scope
2. Requirements and priorities
3. Capacity model
4. High-level architecture
5. Deep dives
6. Failure modes and reliability
7. Observability
8. Trade-offs and redesign

Each stage should ask two things:

- Was it covered?
- Was it covered well enough to change the design quality?

## Data Model & APIs

Recommended checklist fields:

- item name
- observed evidence
- score such as `missed`, `partial`, `strong`
- suggested next drill

Minimal output:

- one major strength
- one critical gap
- one next-practice action

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| review becomes generic praise | no evidence lines attached | require a short evidence note per item |
| checklist becomes too long | reviewers skip half the items | cap to the highest-signal dimensions |
| feedback is only negative | strengths are omitted entirely | capture one concrete strength too |
| reviewer nitpicks style | comments focus on box colors or wording | keep criteria tied to design substance |

## Observability

- metric: distribution of recurring weak checklist categories
- metric: average review completion time
- metric: percentage of reviews that include a next action
- SLO: every mock review should end with one priority fix rather than a long undifferentiated list

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| short checklist | high completion rate | less nuance on rare cases | detailed rubric for every lesson |
| evidence-first review | feedback becomes actionable | a bit slower than gut-feel scoring | numeric score only |
| same checklist across prompts | better longitudinal comparison | less prompt-specific nuance | entirely custom feedback per session |

## Interview It

**Google framing:** "How would you review a candidate's system design answer so the feedback reflects senior-level expectations rather than style preferences?"

**Cloudflare framing:** "How would you review an edge-platform design answer where operational gaps may matter more than one component choice?"

**Follow-ups:**
1. What evidence proves the candidate actually prioritized requirements?
2. How would the checklist change for a 20-minute design drill?
3. What if the answer was technically correct but never reached failure modes?
4. How should the checklist handle a strong redesign follow-up after a weak initial architecture?
5. What is the smallest useful summary you can hand back to the learner?

## Ship It

- `outputs/design-review-checklist.md`
- `outputs/interview-card-review-checklist.md`

## Exercises

1. **Easy** — Review one old answer using the checklist and identify the strongest missing section.
2. **Medium** — Adapt the checklist for a peer reviewer who only has ten minutes.
3. **Hard** — Design a rubric mapping from checklist outcomes to recommended next lessons.

## Further Reading

- [Google SRE hiring and interview advice](https://sre.google/) — useful context for the style of reasoning senior interviews reward
- [The Manager's Path](https://www.oreilly.com/library/view/the-managers-path/9781491973882/) — useful for evidence-based feedback habits
