# Requirement Negotiation and Priority Framing

> Senior candidates do not accept every requirement at face value. They rank them, defend the ranking, and explain what the system optimizes first.

**Type:** Learn
**Company focus:** Google
**Learning goal:** Practice translating broad requirement lists into a priority order that makes the architecture coherent and the trade-offs explicit.
**Prerequisites:** `01-clarification-and-scope/05-prioritization`, `02-estimation-and-cost/01-qps-and-request-mix`, `10-reliability-retries-and-backpressure/04-load-shedding`
**Estimated time:** ~75 min
**Primary artifact:** priority framing matrix

## The Problem

Real interview prompts often contain more goals than any single design can optimize at once. Low latency, strong consistency, low cost, fast feature iteration, multi-region durability, tenant isolation, abuse controls, and simple operations rarely all point in the same direction.

The candidate's job is not only to list these goals. It is to order them and then explain what the system sacrifices.

## Clarify

- Which user-facing property matters most if two goals conflict?
- Which non-functional requirement is truly hard, not merely nice to have?
- Are we designing the first credible version or the long-term final platform?
- Which failures are acceptable and which are business-critical?

If the interviewer does not force a ranking, create one. Say what wins first, what comes second, and what is deferred.

## Requirements

### Functional

- Distinguish must-have from nice-to-have requirements.
- Turn conflicting goals into a visible ranking.
- Map the ranking to concrete architecture choices.
- Explain what gets slower, less fresh, or more expensive because of the ranking.
- Re-rank when the interviewer changes product or business constraints.

### Non-functional

- Avoid fake prioritization where every requirement is "critical."
- Avoid choosing priorities that never affect the architecture.
- Keep trade-offs honest enough that the design remains credible.
- Preserve room for phased rollout and future evolution.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Top-tier requirements | 2 to 3 | more than this usually blurs the design |
| Deferred requirements | 1 to 3 | shows scope control and sequencing |
| Major architecture decisions | 3 to 5 | each should tie back to a priority |
| Redesign triggers | medium | priorities often change under follow-up pressure |

## Architecture

A practical priority framing loop:

```text
collect requirements
  -> rank what matters most
  -> tie each top priority to a design decision
  -> name the cost of that decision
  -> defer or soften lower-priority goals
```

Examples:

1. If low read latency wins over strong freshness, add caching and accept controlled staleness.
2. If correctness wins over availability for writes, tighten the write path and admit slower failover.
3. If low operational complexity matters early, avoid premature multi-region active-active writes.

Priority framing is strongest when the interviewer can predict your next design choice from the requirement ranking.

## Data Model & APIs

Model requirements explicitly:

```text
requirement(name, priority, rationale, impacted_components)
decision(component, supports_requirement, cost)
```

Helpful verbal APIs:

- `RankRequirements()`
- `ExplainWinner(conflict)`
- `DeferRequirement(name)`
- `RebalanceAfterConstraintChange()`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| every requirement is labeled top priority | architecture never sharpens | force a conflict and pick a winner |
| ranking does not affect design | components look generic | tie each major component to one top requirement |
| lower-priority costs stay hidden | design sounds magically optimal | state what gets worse explicitly |
| interviewer changes business goal | answer becomes inconsistent | rerank requirements and point to the changed decisions |

## Observability

- metric: were top requirements named and ordered explicitly?
- metric: did at least one major decision clearly support a top-ranked goal?
- metric: were lower-priority sacrifices stated honestly?
- log: requirement ranking changes during follow-ups
- trace: collect -> rank -> decide -> trade off -> redesign
- SLO: produce an architecture whose priorities are visible from its main decisions

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit requirement ranking | coherent design | some goals are deferred | treating everything as equally urgent |
| phased optimization | lowers complexity early | some future needs wait | building the final platform on day one |
| honest sacrifice statements | credibility and trust | exposes imperfection | pretending there is no downside |

## Interview It

**Google framing:** "Design a system where latency, consistency, scale, and reliability all matter." The strong move is to explain which two matter most now and what that forces you to optimize for.

**Follow-ups:**
1. What if the business now values cost over latency?
2. What if writes become more important than reads?
3. What if compliance changes a deferred requirement into a top priority?
4. What if one region must now be sovereign and isolated?

## Ship It

- `outputs/tradeoff-matrix-priority-framing.md`

## Exercises

1. **Easy** - Rank requirements for a news feed system.
2. **Medium** - Explain why one ranking would favor caching and another would favor fresher origin reads.
3. **Hard** - Redesign a globally served API after cost becomes the top business constraint.

## Further Reading

- [Designing Data-Intensive Applications](https://dataintensive.net/) - strong background for how requirements drive storage and consistency choices
- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) - useful for thinking about operational cost of ambitious guarantees
