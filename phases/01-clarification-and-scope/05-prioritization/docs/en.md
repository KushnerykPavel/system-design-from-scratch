# Prioritizing Requirements with the Interviewer

> Good design starts when the interviewer agrees what matters most.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to negotiate and state requirement priority order so the design is optimized for the right thing instead of trying to win every dimension at once.  
**Prerequisites:** `01-clarification-and-scope/01-functional-vs-non-functional`, `01-clarification-and-scope/04-assumption-logging`  
**Estimated time:** ~60 min  
**Primary artifact:** prioritization matrix + negotiation prompts  

## The Problem

Interview prompts often imply multiple goals at once:

- low latency
- low cost
- high availability
- strict consistency
- rapid launch
- abuse resistance

The candidate who tries to optimize all of them equally usually produces a vague, expensive, and internally inconsistent design. Senior signal comes from surfacing the tension and choosing a priority order with the interviewer.

This lesson focuses on saying things like:

- "I’ll optimize for low read latency over write simplicity."
- "I’ll prioritize availability over strict global consistency."
- "Given the product stage, I’ll trade some operational sophistication for faster delivery."

## Clarify

- Which requirement is the primary optimization target?
- Which requirement is important but allowed to be second-order in v1?
- Which pair of goals is most in tension for this prompt?

## Requirements

### Functional

- State the core workflow and the goal ordering that supports it.
- Confirm at least one requirement the design will not optimize first.
- Revisit priorities if the interviewer adds new constraints.

### Non-functional

- Make trade-off ordering explicit and easy to defend.
- Avoid false precision when the interviewer does not provide hard targets.
- Keep prioritization stable enough to guide all later decisions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Priority set size | top 3 constraints is usually enough | too many "top priorities" weakens decision quality |
| Major tensions | usually 1 or 2 pairs dominate | identifies where trade-offs must be verbalized |
| Reprioritization events | 1 to 3 in a follow-up-heavy interview | tests adaptability without losing coherence |
| Peak factor | highest when the interviewer adds new constraints late | priority order helps absorb changes |
| Rough cost | strong optimization in one area often increases cost elsewhere | prioritization is inseparable from realistic economics |

## Architecture

Prioritization should visibly change the architecture:

- prioritize **latency** and you may add cache layers, edge presence, and read replicas
- prioritize **consistency** and you may simplify regional topology or tighten write acknowledgement rules
- prioritize **cost** and you may reduce replication, narrow hot-path caching, or prefer simpler services
- prioritize **availability** and you may accept weaker write guarantees or more operational redundancy

The code artifact in this lesson validates a ranked requirement list and flags cases where too many goals are declared equally dominant.

## Data Model & APIs

Represent the priority set with:

- `requirement`
- `rank`
- `rationale`

Useful review prompts:

- Which architecture choice would flip if rank one and rank two were swapped?
- Which requirement did you intentionally not maximize?
- Which priority order best matches the business stage?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| everything marked top priority | no architecture choice feels constrained | force a rank ordering |
| priorities stated but ignored later | design choices do not reflect the ranking | tie major components back to ranked goals |
| no tension identified | trade-offs sound generic or absent | name at least one conflicting pair explicitly |
| reprioritization handled poorly | follow-up constraints cause narrative collapse | restate the new ordering before redesigning |

## Observability

- metric: number of explicitly ranked non-functional requirements
- metric: percentage of major architecture decisions tied to a ranked goal
- metric: number of times the ranking is updated after follow-up constraints
- SLO: the top priority should be referenced in the first architecture justification

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit ranking | makes trade-offs defensible | forces you to under-optimize something | equal-weight wish list |
| top-three cap | keeps discussion sharp | less room for nuance up front | ranking every possible concern |
| negotiated priority order | creates interviewer alignment | may expose disagreement early | unilateral design assumptions |

## Interview It

**Google framing:** "Design an internal build artifact service. Which requirements would you prioritize first, and how would that change the architecture?"

**Cloudflare framing:** "Design a customer-facing edge firewall product. How do latency, availability, explainability, and abuse resistance rank against one another?"

**Follow-ups:**
1. Which design choice most clearly reflects your top priority?
2. Which priority would you swap if the product were enterprise instead of consumer?
3. What if compliance becomes more important than latency halfway through?
4. Which requirement is important but deliberately not first-order?
5. How do you avoid sounding arbitrary when ranking goals?

## Ship It

- `outputs/prioritization-matrix.md`
- `outputs/negotiation-prompts.md`

## Exercises

1. **Easy** — Rank latency, cost, and durability for a personal photo backup product.  
2. **Medium** — Reprioritize a chat system when the interviewer changes the prompt from consumer to enterprise.  
3. **Hard** — Defend a priority order for an edge security product where availability, explainability, and strict enforcement all matter.  

## Further Reading

- [Google SRE Book](https://sre.google/sre-book/table-of-contents/) — useful perspective on availability, reliability, and operational trade-offs
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview framing before you add explicit prioritization discipline
