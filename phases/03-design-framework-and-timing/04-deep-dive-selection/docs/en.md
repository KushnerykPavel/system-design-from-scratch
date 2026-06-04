# Choosing the Right Deep Dive

> Depth is only impressive when it is aimed at the system's real risk.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn how to choose one or two deep dives that maximize interview signal instead of defaulting to whichever subsystem feels most familiar.  
**Prerequisites:** `03-design-framework-and-timing/03-diagram-then-dive`  
**Estimated time:** ~60 min  
**Primary artifact:** deep-dive scorecard + interview card  

## The Problem

Strong candidates do not deep dive randomly. They choose the area where scale, correctness, operational complexity, or company-specific expectations are most concentrated.

Weak candidates often deep dive a subsystem because they know how to talk about it, not because it is the bottleneck, failure hotspot, or trade-off center of the design.

## Clarify

- What is the hardest part of this system under the stated constraints?
- Which subsystem would most likely fail first at 10x scale?
- Where is the most controversial trade-off: latency, consistency, cost, or abuse prevention?
- Does the company context suggest an obvious deep-dive target?

If no area dominates, choose the component with the largest blast radius or hardest correctness story.

## Requirements

### Functional

- Rank candidate deep dives.
- Explain why the chosen deep dive matters to the overall system.
- Tie the detail level to explicit risks or constraints.

### Non-functional

- Avoid choosing detail targets based on personal familiarity alone.
- Balance business relevance and technical complexity.
- Adapt the deep dive if the interviewer changes priorities.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Candidate subsystems | 3-5 | enough options to compare |
| Deep dives selected | 1 primary, maybe 1 secondary | keeps focus |
| Risk dimensions | scale, correctness, ops, novelty | drives selection quality |
| Decision time | under 2 min | avoids analysis paralysis |
| Blast radius | directional estimate | identifies what matters most |

## Architecture

A good deep-dive selector asks:

1. Where can the system fail in the most expensive way?
2. Which subsystem absorbs the key requirement tension?
3. Which component changes the most under 10x growth?
4. Which area reveals senior judgment, not just memorized patterns?

Useful heuristics:

- choose write path over read path when correctness dominates
- choose routing or cache layers when latency dominates
- choose rollout, recovery, or isolation when operational risk dominates
- choose config propagation and origin protection for Cloudflare-shaped prompts

## Data Model & APIs

The code artifact ranks deep-dive candidates by a simple composite score:

- `Risk`
- `Scale`
- `Novelty`
- `Dependency`

It is a toy model, but it reflects a strong interview habit: select depth systematically.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Familiarity bias | candidate chooses easiest subsystem to explain | ask which component has the hardest failure story |
| Detail disconnected from prompt | deep dive does not affect main constraints | restate top requirement and re-rank |
| Too many deep dives | conversation fragments across topics | commit to one primary deep dive |
| No deep dive at all | answer stays generic | pick the component with the highest combined risk and scale pressure |

## Observability

- metric: number of candidate subsystems considered explicitly
- metric: whether the chosen deep dive maps to the top non-functional requirement
- metric: time spent justifying the choice
- log: rejected deep-dive candidates and reasons

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| risk-based deep-dive choice | maximizes signal | may skip interesting but lower-value areas | choosing the easiest subsystem |
| one primary deep dive | preserves coherence | less breadth | covering three topics shallowly |
| company-context-aware selection | fits interviewer expectations | can bias away from product breadth | ignoring company shape entirely |

## Interview It

**Google framing:** "Design a task scheduling system." The best deep dive might be fairness, queue partitioning, or retry semantics depending on the workload and SLOs you clarified.

**Cloudflare framing:** "Design a global firewall rule distribution system." The best deep dive is likely configuration propagation, consistency windows, rollback, or origin safety rather than just the API shape.

**Follow-ups:**
1. How do you compare two equally plausible deep dives?
2. When should you pick correctness over scale?
3. How do you justify a deep dive in one sentence?
4. What if the interviewer wants a different area than the one you chose?
5. How does a good deep-dive choice change under a 10x traffic jump?

## Ship It

- `outputs/deep-dive-scorecard.md`
- `outputs/interview-card-deep-dive-selection.md`

## Exercises

1. **Easy** — Rank three deep-dive options for a URL shortener.  
2. **Medium** — Choose the right deep dive for a globally distributed chat system.  
3. **Hard** — Re-rank your deep dive after the interviewer shifts the dominant constraint from latency to correctness.  

## Further Reading

- [Google SRE workbook](https://sre.google/workbook/table-of-contents/) — useful framing for failure-focused deep dives  
- [System design notes](https://github.com/liquidslr/system-design-notes) — helpful baseline for choosing what to explore after the high-level design  
