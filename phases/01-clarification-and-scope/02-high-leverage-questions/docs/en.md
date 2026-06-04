# Clarifying Questions That Actually Change the Design

> The best clarification question is the one that eliminates an entire wrong architecture branch.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to ask a small number of high-leverage questions that change storage, consistency, routing, security, or cost, instead of burning time on trivia.  
**Prerequisites:** `01-clarification-and-scope/01-functional-vs-non-functional`  
**Estimated time:** ~60 min  
**Primary artifact:** question ladder + design-change cheat sheet  

## The Problem

Weak clarification looks active but produces no signal:

- "Should the button be blue or green?"
- "Do users like notifications?"
- "Can I assume a modern tech stack?"

Strong clarification changes the shape of the solution:

- Is the system read-heavy or write-heavy?
- Are writes allowed to be asynchronous?
- Is cross-region failover required?
- Are requests tenant-isolated or globally pooled?

The goal is not to ask many questions. The goal is to ask the few that collapse uncertainty quickly.

## Clarify

- What is the primary workload shape: read-heavy, write-heavy, bursty, or fanout-heavy?
- Which consistency or freshness promise matters for the core user journey?
- Is the system single-region to start, or must it survive regional failure?

## Requirements

### Functional

- Reach a designable v1 scope in under three to five clarification questions.
- Ask questions that map directly to architecture pivots.
- State assumptions when the interviewer does not answer.

### Non-functional

- Minimize clarification time while maximizing architectural signal.
- Avoid low-impact product trivia unless it truly affects capacity or abuse patterns.
- Preserve enough ambiguity to keep momentum when necessary.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification budget | 3 to 5 high-leverage questions | too many questions can stall the answer |
| Architecture branches avoided | ideally 2 or more | good questions eliminate bad paths early |
| Major pivots touched | workload, consistency, scope, failure, cost | these are the categories that usually matter |
| Peak factor | highest in the first 5 interview minutes | front-loaded value matters most here |
| Rough cost | low direct cost, high opportunity cost | bad questions consume time better spent on sizing and design |

## Architecture

A high-leverage question usually targets one of these pivots:

1. **Workload shape** — changes cache, partitioning, and async choices.
2. **Freshness or consistency** — changes write path, replication, and read policy.
3. **Failure scope** — changes regional topology and degradation plan.
4. **Tenant or abuse model** — changes isolation, quotas, and safety controls.
5. **Time horizon** — changes how much extensibility to design for.

The code artifact for this lesson scores draft questions by whether they map to one of these pivots and whether they are likely to influence the design.

## Data Model & APIs

Represent a draft question as:

- `text`
- `pivot` such as `workload`, `consistency`, `scope`, `failure`, `cost`, `security`
- `changes_design` as a boolean

Useful review API:

- `ScoreQuestions([]Question)` to flag weak or redundant clarifiers

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| too many low-impact questions | many questions asked, few design consequences named | tie each question to a possible architecture pivot |
| clarification becomes product brainstorming | discussion drifts away from system behavior | ask about scale, consistency, failures, and scope |
| candidate waits for complete information | architecture never begins | state assumptions and move forward |
| redundant questions | several questions probe the same pivot with no new signal | use a question ladder and stop once the branch is resolved |

## Observability

- metric: number of clarifying questions asked before sizing begins
- metric: percentage of questions that map to a major architecture pivot
- metric: number of assumptions made after unanswered questions
- SLO: at least two of the first three clarification questions should materially constrain the design

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| ask fewer but higher-signal questions | more time for sizing and architecture | some details remain ambiguous | exhaustive discovery |
| move on after unanswered questions | preserves momentum | some assumptions may later be wrong | waiting for perfect clarity |
| bias toward pivot questions | architecture becomes defensible faster | less product-color detail early | trivia-first questioning |

## Interview It

**Google framing:** "Design a global metrics ingestion system. What are the first questions you ask, and why do they matter?"

**Cloudflare framing:** "Design DDoS protection for an API product. Which clarification questions change the edge data path versus the control plane?"

**Follow-ups:**
1. Which of your first questions most changes the storage model?
2. What question most changes the regional design?
3. Which question would you skip if the interviewer is in a hurry?
4. What if the interviewer gives vague answers like 'pretty high scale'?
5. Which question best reveals whether caching is central or optional?

## Ship It

- `outputs/question-ladder.md`
- `outputs/design-change-cheat-sheet.md`

## Exercises

1. **Easy** — For a URL shortener, list three questions that would materially change the design.  
2. **Medium** — For a chat system, pick the best four clarifiers from a list of ten candidate questions.  
3. **Hard** — For an edge abuse prevention system, explain how different answers to the same clarifier produce different architectures.  

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — helpful baseline on early clarification and sizing
- [Google SRE workbook](https://sre.google/workbook/table-of-contents/) — good operational instincts for asking failure-oriented questions
