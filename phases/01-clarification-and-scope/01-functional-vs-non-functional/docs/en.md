# Functional vs Non-Functional Requirements

> A system is not defined only by what it does, but by the promises it must keep while doing it.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to separate user-visible capabilities from quality constraints so the first five interview minutes produce a usable design frame instead of a requirement soup.  
**Prerequisites:** `00-setup-and-workflow/06-review-checklist`  
**Estimated time:** ~60 min  
**Primary artifact:** interview card + requirement sorting checklist  

## The Problem

Many interview answers go wrong before architecture starts because candidates mix different kinds of requirements together:

- "users can upload photos"  
- "p99 should be under 200 ms"  
- "data must survive region loss"  
- "we should support search later"

When those statements are not separated, trade-offs become blurry. The learner starts arguing about databases before deciding whether latency, durability, or time-to-market dominates the design.

This lesson gives you a simple split:

- **Functional requirements** describe what the system must enable.
- **Non-functional requirements** describe how well, how safely, how cheaply, or how reliably it must do it.

Senior candidates do not merely name both categories. They prioritize the non-functional ones because those usually drive architecture.

## Clarify

- What is the core user action or business workflow the system must support on day one?
- Which quality constraint would most change the architecture if tightened: latency, availability, consistency, durability, security, or cost?
- Which desirable features are explicitly out of scope for v1 so the design stays teachable?

## Requirements

### Functional

- Support a clear primary workflow such as create, read, update, publish, upload, search, or notify.
- Name one or two secondary workflows that matter, but are not the main driver.
- State at least one feature that is intentionally deferred.

### Non-functional

- Rank the important quality attributes instead of listing them alphabetically.
- Tie each top requirement to a design consequence.
- Be explicit when two constraints are in tension, such as low latency versus strong consistency.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak QPS | 10K to 100K for a common interview-sized consumer system | determines whether latency constraints become edge/cache problems |
| Storage growth | GB/day to TB/year depending on object size | separates metadata design from blob design |
| Bandwidth | request/response size x peak QPS | exposes CDN, compression, or fanout concerns |
| Peak factor | 3x to 10x over average | reveals whether availability and overload handling outrank steady-state optimization |
| Rough cost | low millions or less per year for a realistic v1 | forces trade-offs instead of unlimited infrastructure fantasy |

## Architecture

The requirement split changes the architecture conversation:

1. Start with the primary workflow.
2. Identify the one dominant quality constraint.
3. Choose an architecture that is good at that constraint.
4. Name what becomes weaker as a result.

Example:

- If the functional requirement is "users can upload and retrieve large media objects," the baseline architecture might be API tier + metadata store + object storage.
- If the dominant non-functional requirement is low global read latency, you likely add CDN and regional replication.
- If the dominant non-functional requirement is strict durability, you focus more on write acknowledgements, replication policy, and background repair than on cache placement.

The code artifact for this lesson is a small requirement classifier that checks whether a draft requirement set includes both categories and whether non-functional constraints were prioritized.

## Data Model & APIs

Represent each requirement with:

- `text`
- `kind` as `functional` or `non_functional`
- `priority`
- optional `driver` flag for the dominant non-functional requirement

Useful review prompts:

- Which requirement changes the storage choice?
- Which requirement changes the regional topology?
- Which requirement changes the failure policy?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| all requirements presented as one flat list | architecture discussion starts without clear priorities | separate functional from non-functional explicitly |
| too many "must-have" non-functionals | every design choice sounds equally urgent | force a top three ranking |
| features treated as quality goals | phrases like "support notifications" mixed with latency goals | classify each statement before designing |
| quality goals not tied to consequences | candidate says "highly available" but cannot explain how it changes design | ask what concrete architecture choice that requirement forces |

## Observability

- metric: number of requirements classified before architecture begins
- metric: count of explicitly ranked non-functional requirements
- metric: number of architecture choices justified by a named requirement
- SLO: every answer should identify one dominant non-functional driver before the first detailed component deep dive

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit requirement split | cleaner reasoning and better interviewer alignment | adds a small up-front structuring step | free-form brainstorming |
| ranking non-functional goals | exposes the real architecture driver | forces uncomfortable de-prioritization | treating all qualities as equally important |
| deferring secondary features | preserves design clarity | less product breadth in v1 | pretending the first answer covers every future requirement |

## Interview It

**Google framing:** "Design Google Drive sharing for large teams. Before architecture, tell me what the functional requirements are and which quality constraints dominate."

**Cloudflare framing:** "Design an API gateway product. Separate what customers need the product to do from the operational guarantees that shape the edge design."

**Follow-ups:**
1. Which non-functional requirement would most change your design if it tightened by 10x?
2. What requirement is important but intentionally not optimized in v1?
3. Which two requirements are in the most tension?
4. What if the interviewer says, "make it cheap" after you optimized for global latency?
5. Which requirement justifies the first deep dive?

## Ship It

- `outputs/interview-card-functional-vs-non-functional.md`
- `outputs/requirement-sorting-checklist.md`

## Exercises

1. **Easy** — Take a URL shortener prompt and sort eight mixed requirement statements into the two categories.  
2. **Medium** — For a chat system, rank five non-functional goals and justify the top two.  
3. **Hard** — Rework a design where durability, cost, and low global latency are all requested, and explain what you would deliberately under-optimize first.  

## Further Reading

- [Google SRE Book](https://sre.google/sre-book/table-of-contents/) — strong grounding for turning vague quality claims into operational consequences
- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline on clarifying scope before design
