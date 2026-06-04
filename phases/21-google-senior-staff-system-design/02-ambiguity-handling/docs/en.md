# Handling Ambiguous Prompts Like a Senior Engineer

> Ambiguity is not a trap to escape quickly. It is often the fastest path to showing senior judgment.

**Type:** Learn
**Company focus:** Google
**Learning goal:** Learn how to convert vague prompts into a bounded system design problem without sounding defensive, rigid, or overly academic.
**Prerequisites:** `01-clarification-and-scope/02-high-leverage-questions`, `01-clarification-and-scope/04-assumption-logging`, `03-design-framework-and-timing/07-interviewer-moves`
**Estimated time:** ~75 min
**Primary artifact:** ambiguity-handling checklist

## The Problem

Google-style prompts are often deliberately under-specified. "Design Google Calendar." "Design a file sync system." "Design a recommendations service." The goal is not to guess the hidden exact problem. The goal is to scope responsibly, expose the design-changing uncertainties, and move forward decisively once enough is known.

Senior answers do not ask endless questions. They ask the few questions that determine the shape of the system.

## Clarify

- Who are the primary users and what is their main action?
- What scale or growth matters enough to affect architecture?
- Which guarantee matters most: latency, freshness, correctness, availability, or cost?
- What can reasonably be deferred or out of scope?

When the interviewer answers loosely, state an assumption, explain why it is reasonable, and continue. Do not wait for perfect information.

## Requirements

### Functional

- Extract the main product workflow from the vague prompt.
- Turn ambiguity into a prioritized requirement set.
- Log assumptions explicitly and revisit them when needed.
- Bound the scope enough to size and design confidently.
- Keep room for the interviewer to redirect the design later.

### Non-functional

- Avoid interrogating the interviewer with a long questionnaire.
- Avoid over-scoping the problem to sound thorough.
- Make uncertainty visible without sounding stuck.
- Preserve momentum while still being rigorous.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification questions | 3 to 6 | too few leaves hidden risk, too many burns time |
| Assumptions stated aloud | 2 to 4 | enough to unlock design without overfitting |
| Scope cuts | 1 to 3 | keeps the answer bounded and credible |
| Redirection frequency | medium | the design should stay adaptable |

## Architecture

Use an ambiguity-handling loop:

```text
extract core user journey
  -> ask 1 to 3 design-changing questions
  -> state assumptions for unresolved areas
  -> prioritize requirements
  -> design against the chosen shape
```

Helpful question categories:

1. Workload shape: read-heavy, write-heavy, bursty, or correctness-sensitive?
2. Product expectations: interactive UX, batch analytics, or hybrid?
3. Trust boundary: public clients, internal services, or both?
4. Evolution pressure: global scale, multi-region, or regulatory constraints?

Strong candidates stop clarifying once the system can be sized and decomposed. They do not keep asking once the answer is already moving.

## Data Model & APIs

Useful structure for assumptions:

```text
assumption(
  topic,
  chosen_value,
  reason,
  possible_design_impact
)
```

Example API for your own answer:

- `ClarifyUsers()`
- `ClarifyScale()`
- `ClarifyPrimaryConstraint()`
- `SetOutOfScope()`
- `ProceedWithAssumptions()`

If your assumptions are hidden, the interviewer cannot tell whether your decisions were thoughtful or accidental.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate asks many low-impact questions | time passes without design movement | ask only questions that change data model, topology, or guarantees |
| candidate never states assumptions | architecture depends on unstated guesses | summarize assumptions before sizing |
| candidate over-expands scope | answer becomes broad but thin | pick one core workflow and defer secondary features |
| interviewer changes a key assumption later | design feels brittle | tie decisions to assumptions so redesign becomes straightforward |

## Observability

- metric: number of clarification questions that changed the design materially
- metric: whether scope cuts were stated before architecture
- metric: whether one primary non-functional goal was chosen
- log: assumptions and later corrections
- trace: prompt -> questions -> assumptions -> scoped design
- SLO: convert an ambiguous prompt into a designable problem in under five minutes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit assumptions | keeps momentum and honesty | may need redesign later | waiting for perfect certainty |
| narrow core workflow | deeper, more coherent answer | less breadth | trying to solve every feature equally |
| few high-impact questions | preserves interview time | some ambiguity remains | exhaustive discovery session |

## Interview It

**Google framing:** "Design a broad product or platform." Interviewers often want to see whether you can turn ambiguity into a manageable design problem without becoming either passive or domineering.

**Follow-ups:**
1. What if you picked the wrong assumption about scale?
2. What if the interviewer says, "make it global" halfway through?
3. What if two requirements conflict and both sound important?
4. What if you are unsure which deep dive matters most?

## Ship It

- `outputs/ambiguity-checklist-google-system-design.md`

## Exercises

1. **Easy** - Ask three design-changing questions for "design a calendar system."
2. **Medium** - Convert "design an internal analytics platform" into a bounded first-pass scope.
3. **Hard** - Redesign your assumptions after the interviewer changes the prompt from internal-only to global external traffic.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) - useful reminder of the four-step interview flow and early scoping rhythm
- [Google SRE book](https://sre.google/sre-book/table-of-contents/) - helpful for thinking about how assumptions affect operational guarantees
