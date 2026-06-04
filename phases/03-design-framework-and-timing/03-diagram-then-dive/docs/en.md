# High-Level Diagram First, Deep Dive Second

> The first diagram should orient the room, not exhaust the room.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Practice sequencing: present a shared mental model first, then go deep only after the main architecture and critical path are clear.  
**Prerequisites:** `03-design-framework-and-timing/02-time-boxing`  
**Estimated time:** ~60 min  
**Primary artifact:** diagram checklist + interview card  

## The Problem

Candidates often deep dive too early. They zoom into cache invalidation, replication, or queue semantics before the interviewer has even agreed on the system shape. That usually creates confusion, rework, and time loss.

This lesson teaches the discipline of earning the deep dive: first show the high-level flow, then choose one subsystem that deserves detailed treatment.

## Clarify

- What single end-to-end user flow should the first diagram explain?
- Which non-functional constraint most shapes the top-level architecture?
- Is there an obvious critical path the interviewer cares about?
- Are there any components we should explicitly defer from the first diagram?

If the interviewer does not steer, optimize the first picture for clarity, request flow, and storage boundaries.

## Requirements

### Functional

- Produce a high-level diagram that explains the main request path.
- Identify the subsystem most worth deep-diving.
- Transition cleanly from overview to detail.

### Non-functional

- Avoid overloading the first diagram with implementation detail.
- Keep the deep dive tied to the dominant system risk.
- Preserve interviewer alignment before investing in details.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| First-diagram time | 5-8 min | keeps overview concise |
| Components on first pass | 5-9 boxes | enough structure without clutter |
| Deep-dive branches | 1 primary branch | preserves focus |
| Main request paths | 1-2 | shows system shape cleanly |
| Rework cost | high if overview is missing | details may become irrelevant |

## Architecture

Think in two passes:

1. **Orientation pass**
   - clients
   - entry point
   - stateless tier
   - storage or async boundary
   - critical path
2. **Deep-dive pass**
   - choose the highest-risk or highest-leverage component
   - explain data flow, failure handling, and trade-offs

Bad sequence:

```text
cache eviction -> replica repair -> queue retry policy -> "wait, what system is this?"
```

Better sequence:

```text
user request flow -> main components -> critical path -> chosen subsystem detail
```

## Data Model & APIs

The code artifact checks whether a review is ready for deep dive:

- do we have a high-level diagram?
- is the critical path explicit?
- has the deep dive been chosen?

That mirrors the interview discipline you want to build.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| First diagram is too detailed | many low-level boxes but unclear request path | redraw with fewer components and clearer boundaries |
| Deep dive starts too early | topology still changing during detailed discussion | pause and confirm the system shape first |
| Critical path is ambiguous | interviewer cannot tell what matters most | narrate one canonical read or write flow |
| Deep dive is disconnected from constraints | detail does not map to main risk | restate why this subsystem matters now |

## Observability

- metric: time spent before the first high-level diagram is understandable
- metric: number of components in the first pass
- metric: whether the critical path was named before deep dive
- trace: interviewer follow-up sequence after the first diagram

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| simple first diagram | quick alignment | less immediate depth | drawing every subsystem upfront |
| delayed deep dive | reduces rework | can feel less technical at first | jumping into detail before buy-in |
| critical-path-based detail | ties discussion to risk | may defer interesting side systems | choosing a favorite subsystem at random |

## Interview It

**Google framing:** "Design a news feed." Start with the end-to-end flow and defer ranking or fanout detail until the interviewer agrees on the broad architecture.

**Cloudflare framing:** "Design edge request routing for an API gateway." Show request ingress, routing, policy enforcement, and origin interaction before diving into cache semantics or configuration propagation.

**Follow-ups:**
1. What belongs on the first diagram and what does not?
2. When should the first diagram include async pipelines?
3. How do you explain the critical path without turning it into a deep dive?
4. What if the interviewer asks for detail before the overview is complete?
5. How do you redraw quickly when a constraint changes?

## Ship It

- `outputs/diagram-review-checklist.md`
- `outputs/interview-card-diagram-then-dive.md`

## Exercises

1. **Easy** — Draw a first-pass diagram for a URL shortener using no more than seven boxes.  
2. **Medium** — Choose the right deep dive for a chat system after presenting the high-level shape.  
3. **Hard** — Re-sequence a messy architecture-first answer into a cleaner two-pass explanation.  

## Further Reading

- [The C4 model](https://c4model.com/) — useful reminder that abstraction level matters  
- [System design notes](https://github.com/liquidslr/system-design-notes) — classic sequence of overview before detail  
