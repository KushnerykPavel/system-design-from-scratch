# Google Rubric and Strong Signals

> Google-style senior interviews reward answers that stay structured under ambiguity and keep earning trust all the way through the deep dive.

**Type:** Build
**Company focus:** Google
**Learning goal:** Internalize the evaluation signals that separate a senior-level Google design answer from a component dump or premature deep dive.
**Prerequisites:** `01-clarification-and-scope/05-prioritization`, `03-design-framework-and-timing/01-four-step-interview-loop`, `11-observability-slos-and-debugging/01-sli-slo-error-budget`
**Estimated time:** ~75 min
**Primary artifact:** rubric scorecard validator + interview signal card

## The Problem

You are not only designing a system. You are showing judgment. In Google senior and staff loops, interviewers are usually scoring whether you clarify scope well, estimate realistically, choose a coherent architecture, go deep in the right place, and communicate trade-offs like an owner.

Weak candidates think the rubric is "name enough distributed systems parts." Strong candidates realize the rubric is really about decision quality under time pressure.

## Clarify

- Is this a senior loop or a staff loop?
- Is the interviewer looking for product-serving scale, infrastructure depth, or reliability reasoning?
- Is the system mostly read-heavy, write-heavy, or correctness-sensitive?
- Which part of the answer should earn trust first: scope control, sizing, or operational depth?

If the interviewer is vague, assume a senior/staff-serving prompt with one primary workload, one primary bottleneck, and one deep dive that should be chosen deliberately rather than reactively.

## Requirements

### Functional

- Open by narrowing scope and defining the primary user journey.
- Turn vague requirements into prioritized functional and non-functional goals.
- Do rough sizing before committing to topology.
- Present a high-level design with one explicit deep-dive candidate.
- Close with failure handling, observability, and realistic follow-ups.

### Non-functional

- Sound organized without sounding scripted.
- Avoid magical guarantees or infrastructure cargo culting.
- Expose trade-offs clearly enough that the interviewer can trust later claims.
- Make redesign under changed constraints feel natural.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification time | 3 to 5 min | early structure is a scoring signal by itself |
| Sizing time | 4 to 6 min | proves the architecture is workload-aware |
| High-level design time | 8 to 12 min | creates shared context before depth |
| Deep-dive time | 8 to 10 min | where judgment and systems fluency become visible |
| Wrap-up time | 3 to 5 min | reveals whether you think about risk and evolution |

## Architecture

Treat the rubric like a pipeline:

```text
scope and assumptions
  -> requirements and sizing
  -> high-level architecture
  -> deliberate deep dive
  -> failure and observability story
  -> redesign and trade-offs
```

Signals that usually score well:

1. You choose one main workload shape and keep returning to it.
2. You tie major components to concrete requirements rather than listing them generically.
3. You name what will probably fail first and how you would detect it.
4. You admit uncertainty and state the assumption instead of bluffing.

Signals that usually score poorly:

1. Starting with a giant diagram before clarifying scale or product scope.
2. Claiming exact correctness, zero downtime, or "just shard it" without boundaries.
3. Spending half the interview on one subsystem the prompt did not require.
4. Treating observability and rollout as optional add-ons.

## Data Model & APIs

For rubric-driven practice, the "data model" is the answer structure itself:

```text
prompt_scope(
  users,
  main actions,
  scale,
  latency target,
  durability target
)

answer_plan(
  top_requirements,
  sizing_assumptions,
  chosen_architecture,
  deep_dive_area,
  top_risks,
  observability_plan
)
```

Helpful mental APIs:

- `Clarify(prompt) -> scope`
- `Estimate(scope) -> capacity model`
- `Propose(scope, capacity) -> high-level design`
- `DeepDive(component) -> guarantees and trade-offs`
- `Redesign(change) -> updated architecture`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate starts architecture before scope | no explicit workload or priorities in first minutes | force a one-sentence scope statement first |
| sizing is skipped entirely | component choices are ungrounded | estimate order of magnitude traffic and storage before topology |
| deep dive is accidental | conversation drifts into random subsystem detail | choose the deep dive explicitly and say why |
| answer ends before risk discussion | no failure or observability narrative appears | reserve wrap-up time and use a short checklist |

## Observability

- metric: did the answer define one primary workload and one primary non-functional constraint?
- metric: did the answer include at least one quantitative estimate?
- metric: did every major component have a reason to exist?
- metric: were failure detection signals named, not only mitigations?
- log: which assumptions changed during the discussion
- trace: clarify -> size -> design -> deep dive -> redesign
- SLO: produce a coherent Google-style system design answer that stays structured from opening to wrap-up

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| lead with scope and sizing | builds trust early | less time for decorative detail | diagram-first answer |
| one deliberate deep dive | shows judgment and depth | not every subsystem gets airtime | shallow coverage of everything |
| explicit risks and observability | operational credibility | forces hard choices into the open | idealized happy-path design |

## Interview It

**Google framing:** "Design a large-scale system and explain the key trade-offs." Interviewers often reward clear structure, good prioritization, and willingness to reason about bottlenecks and degraded behavior instead of trying to sound encyclopedic.

**Strong signals to demonstrate:**
1. Clarifying questions that change the design.
2. Rough sizing that actually informs the topology.
3. A calm high-level design before the deep dive.
4. Honest trade-offs around consistency, latency, cost, and complexity.
5. A redesign answer that changes something real.

## Ship It

- `outputs/interview-card-google-rubric.md`
- `outputs/design-review-google-strong-signals.md`

## Exercises

1. **Easy** - Write a two-minute opening for "design a metrics ingestion system."
2. **Medium** - Explain how you would choose a deep dive for "design Google Photos."
3. **Hard** - Critique an answer that has perfect component coverage but no sizing, no risks, and no observability.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) - useful for the operational instincts interviewers often probe
- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) - a classic reference for tail latency thinking in large distributed systems
