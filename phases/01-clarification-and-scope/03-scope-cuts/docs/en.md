# How to Cut Scope Without Dodging the Prompt

> Scope cuts are not an escape hatch. They are how you preserve a coherent answer under time pressure.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to narrow a design honestly by trimming breadth while preserving the core job of the system and the interviewer’s intent.  
**Prerequisites:** `01-clarification-and-scope/02-high-leverage-questions`  
**Estimated time:** ~60 min  
**Primary artifact:** scope-cut playbook + interview card  

## The Problem

Senior interviews reward focused answers, but many candidates cut scope badly:

- they remove the hardest requirement without naming the trade-off
- they defer the core user journey instead of a secondary feature
- they say "for simplicity" so often that the prompt stops meaning anything

Good scope cuts do three things:

1. keep the main user journey intact
2. remove complexity that is not essential to the first architecture pass
3. explain the consequence of the cut and what changes later

This is the difference between disciplined framing and dodging the problem.

## Clarify

- What is the irreducible core workflow that must remain in v1?
- Which complexity source is secondary: multi-region, advanced search, analytics, rich permissions, or long-tail product features?
- Which scope cut best preserves the architectural lesson while keeping the answer finishable in 45 minutes?

## Requirements

### Functional

- Preserve the core user journey and success criterion.
- Defer secondary workflows explicitly and honestly.
- State how a deferred feature would alter the design later.

### Non-functional

- Keep the narrowed scope credible for an interview setting.
- Reduce architecture breadth enough to leave room for depth and trade-offs.
- Avoid cuts that erase the main scaling or reliability challenge.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Core workflows kept | 1 to 2 | too many v1 workflows create architectural sprawl |
| Deferred features | 2 to 4 common extras | forces the answer to prioritize |
| Architecture branches removed | ideally 1 to 3 | healthy cuts simplify without hollowing out the prompt |
| Peak factor | highest when prompt is very broad | scope control prevents early overload |
| Rough cost | lower complexity now, migration cost later | every cut shifts future work into redesign or rollout planning |

## Architecture

A disciplined scope cut should target one of these complexity buckets:

- **geography** — start single-region before active-active global design
- **consistency envelope** — start with bounded eventual consistency before global strict semantics
- **product breadth** — support the main workflow before adding analytics, search, or admin surfaces
- **operational features** — mention but defer advanced self-serve tuning, complex reporting, or full automation

Example:

- Prompt: "Design Dropbox."
- Good cut: "For v1, I’ll focus on upload, download, and metadata sync for personal accounts, not collaborative editing."
- Weak cut: "For simplicity, I’ll ignore storage and synchronization."

The code artifact checks whether a proposed scope cut preserves the core workflow and whether the deferred items were explained instead of silently dropped.

## Data Model & APIs

Represent a scoped answer with:

- `core_workflows`
- `deferred_features`
- `reason`
- `preserves_prompt_intent`

Useful review prompts:

- Did the cut remove difficulty or merely rename it?
- Would the interviewer still recognize the original problem?
- Did you explain what changes when the deferred feature returns?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| cut removes the core problem | the remaining design no longer answers the prompt | preserve the main user journey first |
| every hard part deferred | answer becomes shallow and unrealistic | keep one meaningful challenge in scope |
| no consequence named | deferred feature vanishes without redesign discussion | say what topology, state, or API changes later |
| "for simplicity" used as a shield | repeated vague deferrals with no rationale | tie each cut to time-boxing or learning value |

## Observability

- metric: number of explicit scope cuts made before architecture detail expands
- metric: whether each cut has a stated reason and future design consequence
- metric: percentage of core workflows still represented in the final architecture
- SLO: scope cuts should reduce complexity without changing the identity of the problem

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| narrow to one core workflow | preserves clarity and depth | less breadth in the initial answer | attempt full feature parity up front |
| defer multi-region or advanced extras | keeps v1 teachable | later redesign work is required | instantly designing the final global system |
| explain cuts explicitly | builds interviewer trust | spends a little time on framing | silent deferrals that feel evasive |

## Interview It

**Google framing:** "Design Google Photos sharing. What would you intentionally defer in a 45-minute interview answer, and why?"

**Cloudflare framing:** "Design an edge configuration platform. Which parts of the product would you keep in v1 so the control plane and propagation model stay discussable?"

**Follow-ups:**
1. Which cut preserves the hardest architectural lesson?
2. Which cut would be suspicious because it removes the core challenge?
3. What changes when the deferred feature comes back at 10x scale?
4. How do you scope-cut without sounding like you are avoiding operational detail?
5. Which cut is safest when the interviewer has not yet clarified scale?

## Ship It

- `outputs/scope-cut-playbook.md`
- `outputs/interview-card-scope-cuts.md`

## Exercises

1. **Easy** — Propose two honest scope cuts for a URL shortener and explain which one is better.  
2. **Medium** — Narrow a feed system prompt so the answer still includes a meaningful scaling challenge.  
3. **Hard** — Reframe a global file-sync prompt into a credible v1 while preserving the later path to multi-region expansion.  

## Further Reading

- [C4 model](https://c4model.com/) — helpful reminder that abstraction and scoping are deliberate choices
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline framing for narrowing scope before deeper design
