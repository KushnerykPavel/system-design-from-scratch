# Constraint Change and Redesign Prompts

> The best redesign answers start by naming what assumptions broke.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Practice responding to interviewer constraint changes without restarting from zero or hand-waving over what really changes.  
**Prerequisites:** `03-design-framework-and-timing/05-wrap-up`  
**Estimated time:** ~45 min  
**Primary artifact:** redesign trigger sheet + interview card  

## The Problem

Senior system design interviews rarely end with the first architecture. The interviewer often changes scale, consistency, latency, cost, or trust assumptions to see whether your design has structure and whether you can evolve it intelligently.

Many candidates panic and either restart entirely or claim nothing important changes. Both lose signal.

## Clarify

- Which assumption changed: scale, latency, correctness, geography, cost, or abuse model?
- Is the original architecture still valid but stressed, or fundamentally mismatched?
- Which components are affected first?
- What can remain unchanged to preserve continuity?

If the change is ambiguous, restate the old assumption and the new one before proposing revisions.

## Requirements

### Functional

- Identify which prior assumptions broke.
- Explain what parts of the design stay stable and what must change.
- Redesign incrementally when possible.

### Non-functional

- Preserve composure under mid-interview change.
- Avoid full restarts unless the prompt truly becomes a new system.
- Tie the redesign to explicit trade-offs and migration risk.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Constraint changes | 1-3 per interview | common interviewer tactic |
| Redesign time | 5-10 min | must be fast and focused |
| First-order impacts | 2-4 components | keeps answer specific |
| Migration steps | 1-3 | demonstrates realism |
| Failure risks added | at least 1 new risk | shows operational awareness |

## Architecture

Use a four-part redesign response:

1. State the old assumption and the new assumption.
2. Identify the first components that break or become bottlenecks.
3. Propose the smallest architecture changes that address the new constraint.
4. Mention new trade-offs, migration steps, and observability changes.

Examples of high-leverage changes:

- single-region to multi-region
- 100K QPS to 10M QPS
- eventual consistency to strict correctness on writes
- low cost to low latency
- trusted internal users to hostile internet traffic

## Data Model & APIs

The code artifact scores redesign pressure from changed constraints. It is intentionally simple, but it teaches the right reflex: redesign starts from changed assumptions, not from generic component swapping.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Full restart on small change | answer discards stable parts unnecessarily | preserve unchanged components explicitly |
| No real redesign on large change | same architecture repeated despite broken assumptions | name the first bottleneck or correctness mismatch |
| Migration ignored | new topology described as instant cutover | mention staged rollout or dual-run plan |
| New risks unstated | redesign adds complexity without failure analysis | call out at least one new failure mode |

## Observability

- metric: time to identify which assumption changed
- metric: number of components whose behavior must change
- metric: whether migration and rollback were mentioned
- alert: redesign answer with no new trade-off or failure risk

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| incremental redesign | keeps continuity | may preserve legacy constraints | starting from a clean-sheet answer |
| explicit assumption delta | makes reasoning clear | costs some time upfront | jumping to component changes immediately |
| migration-aware answer | more realistic | less time for ideal end-state architecture | pretending cutover is free |

## Interview It

**Google framing:** "Your scheduler now needs strict job ordering across regions." The redesign must discuss what breaks, where coordination moves, and what throughput or latency cost follows.

**Cloudflare framing:** "Your edge configuration plane must now support global rollback within seconds." The redesign should focus on propagation topology, safety controls, blast radius, and stale-config detection.

**Follow-ups:**
1. When is a full restart actually justified?
2. How do you choose between a topology change and a tuning change?
3. What migration risks matter most in redesign discussions?
4. How do you preserve interview momentum after a big constraint shift?
5. Which observability additions usually accompany a redesign?

## Ship It

- `outputs/interview-card-redesign-prompts.md`
- `outputs/redesign-trigger-sheet.md`

## Exercises

1. **Easy** — Rework a URL shortener for a 10x QPS increase.  
2. **Medium** — Adapt a single-region notification system to multi-region delivery.  
3. **Hard** — Redesign a low-cost analytics pipeline for stronger correctness guarantees and explain the migration plan.  

## Further Reading

- [Google SRE workbook: disaster recovery](https://sre.google/workbook/disaster-recovery-in-sre/) — helpful thinking for changed constraints and recovery design  
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline framework before layering redesign maturity  
