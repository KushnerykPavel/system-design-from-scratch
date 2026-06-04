# Diagramming Fast Under Interview Pressure

> The diagram is a communication tool, not a museum exhibit.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Learn to produce a fast high-level diagram that exposes request flow, state boundaries, and scaling choices without drowning the interview in boxes.
**Prerequisites:** `00-setup-and-workflow/02-note-taking-system`
**Estimated time:** ~60 min
**Primary artifact:** diagram stencil + diagram budget checker

## The Problem

Candidates often waste precious minutes polishing diagrams instead of using them to guide the conversation. Senior-level diagrams should answer the interviewer’s next question before they ask it: where does the request go, where is state stored, where can the system fail, and what part deserves a deeper dive?

This lesson teaches a box budget. If a component does not change a trade-off, it probably does not belong in the first-pass diagram.

## Clarify

- What is the single request path or data flow that the interviewer must understand first?
- Which two or three stateful boundaries actually shape the design?
- What information is better spoken than drawn in the first five minutes?

## Requirements

### Functional

- Show the main clients, routing layer, compute path, and durable state.
- Preserve one visible place to annotate scale, cache, or queue decisions.
- Leave room for at least one follow-up deep dive.

### Non-functional

- Initial diagram should fit in roughly 2 to 4 minutes.
- The number of top-level boxes should stay low enough to explain quickly.
- Visual structure must survive interviewer interruptions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| First-pass boxes | 5 to 8 | too many boxes slow explanation and hide priorities |
| Primary arrows | 4 to 10 | enough to show flow without turning into wiring |
| Annotation slots | 3 to 5 | reserve space for latency, QPS, and risk notes |
| Peak factor | highest during the first 10 interview minutes | must optimize for speed and clarity |
| Rough cost | low | the real cost is time, not tooling |

## Architecture

Use a layered stencil:

1. Clients or producers
2. Entry point such as gateway, load balancer, or ingestion tier
3. Stateless service tier
4. Async path if essential
5. Durable state

Then annotate only the information that changes the conversation:

- read-heavy versus write-heavy
- cache or queue presence
- regional boundary
- failure isolation boundary

The code artifact in this lesson validates whether a proposed first-pass diagram stays within a reasonable complexity budget.

## Data Model & APIs

Represent a diagram as:

- components with `name`, `kind`, and `critical`
- edges with `from`, `to`, and `purpose`
- optional annotations such as `qps`, `latency`, `consistency`, `cache`

The first-pass API should answer:

- what are the top-level boxes?
- which arrows are required to tell the primary request story?
- which component is likely to be the deep dive?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| diagram too detailed | more time spent drawing than explaining | enforce a first-pass box budget |
| diagram hides stateful boundaries | caches, queues, and databases all blurred together | label durable and transient state explicitly |
| arrows are ambiguous | interviewer asks where writes or reads actually go | name arrows by purpose, not just direction |
| deep dive starts before alignment | candidate zooms into internals immediately | get buy-in on high-level topology first |

## Observability

- metric: number of boxes on the initial diagram
- metric: time to first complete explanation of request flow
- metric: number of interviewer clarification questions caused by diagram ambiguity
- SLO: first-pass diagram should be explainable in under four minutes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| small box budget | clearer communication | some details deferred verbally | full topology up front |
| layered layout | request path is obvious | some cross-cutting concerns are omitted at first | free-form whiteboard placement |
| verbal annotation of nuance | faster momentum | less persistent detail in the drawing | fully annotated diagram |

## Interview It

**Google framing:** "Show me the first diagram you would draw for a global file metadata service. What do you intentionally leave out?"

**Cloudflare framing:** "Sketch the first-pass architecture for a global API gateway. How do you show edge routing and origin protection without over-drawing?"

**Follow-ups:**
1. What if the interviewer wants the write path before you finish the initial diagram?
2. How do you represent multi-region concerns without drawing every region?
3. When should a queue appear in the first-pass view?
4. What if the system has two equally important request paths?
5. How would you redraw the diagram at 10x scale while preserving continuity?

## Ship It

- `outputs/diagram-stencil.md`
- `outputs/interview-card-fast-diagramming.md`

## Exercises

1. **Easy** — Sketch a first-pass diagram for a URL shortener using at most six boxes.
2. **Medium** — Redraw a chat system diagram to make state boundaries more explicit.
3. **Hard** — Produce two versions of a CDN diagram: one for the opening five minutes and one for the deep dive.

## Further Reading

- [C4 model](https://c4model.com/) — useful discipline for separating high-level and detailed views
- [Google SRE Book](https://sre.google/sre-book/table-of-contents/) — strong reminder that communication must preserve operational boundaries
