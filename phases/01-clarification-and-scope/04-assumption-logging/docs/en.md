# Assumption Logging Under Ambiguity

> Ambiguity is not the enemy. Untracked ambiguity is.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to record assumptions explicitly so you can keep moving under incomplete information while preserving honesty, adaptability, and interviewer trust.  
**Prerequisites:** `01-clarification-and-scope/02-high-leverage-questions`  
**Estimated time:** ~60 min  
**Primary artifact:** assumption log template + interview card  

## The Problem

Interview prompts are intentionally incomplete. If you wait for perfect answers, you stall. If you silently invent facts, your design becomes brittle and inconsistent.

Assumption logging is the bridge:

- ask the high-leverage question
- if the answer is vague or missing, state a reasonable assumption
- tie the assumption to the design
- name what changes if the assumption flips

This is especially important in senior interviews, because interviewers often evaluate how well you manage uncertainty, not just how quickly you draw boxes.

## Clarify

- Which unanswered question matters enough to record as an explicit assumption?
- Is the assumption about scale, consistency, geography, abuse, or product scope?
- What part of the design would change first if the assumption proves false?

## Requirements

### Functional

- Record assumptions in a compact and reviewable way.
- Link each assumption to a design decision or deferred branch.
- Update or replace assumptions when the interviewer gives new information.

### Non-functional

- Keep the assumption log short enough to use in real time.
- Preserve answer momentum while maintaining intellectual honesty.
- Avoid inconsistent assumptions across different parts of the design.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Explicit assumptions | 3 to 7 in a typical interview answer | enough to anchor the design without overwhelming the flow |
| Assumption categories | scale, consistency, geography, scope, abuse | keeps the log balanced and useful |
| Flip impact | 1 to 2 major design branches per assumption | helps prioritize which assumptions deserve airtime |
| Peak factor | highest when interviewer answers are vague | ambiguity management is most valuable early |
| Rough cost | low recording cost, high value in correctness | a short log avoids expensive design contradictions later |

## Architecture

An effective assumption log behaves like a control surface for the architecture:

1. record the assumption
2. build the design against it
3. mark the branch that changes if the assumption flips

Example:

- assumption: "I’ll assume the system starts single-region with disaster recovery, not active-active."
- immediate design effect: simpler write path and replication model
- flip consequence: if active-active is required, conflict resolution and regional failover become first-order concerns

The code artifact validates that assumptions have categories, rationale, and explicit impact notes.

## Data Model & APIs

Represent each assumption with:

- `statement`
- `category`
- `impact`
- `reversible`

Useful review prompts:

- Which assumption is most dangerous if wrong?
- Which assumption affects the write path?
- Which assumption is safe to change later versus expensive to reverse?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| silent assumptions | architecture contains facts never discussed aloud | state assumptions explicitly before using them |
| inconsistent assumptions | one part of the answer assumes one region, another assumes many | keep a compact written log and revisit it during follow-ups |
| too many assumptions | clarification turns into a catalog of guesses | record only design-relevant uncertainty |
| no flip analysis | assumptions are stated but not connected to redesign | always say what changes if the assumption breaks |

## Observability

- metric: number of explicit assumptions recorded before the high-level design is finalized
- metric: percentage of assumptions tied to a concrete design consequence
- metric: number of contradictions found during later follow-ups
- SLO: assumptions should reduce ambiguity without creating internal inconsistency

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit assumption log | preserves honesty and coherence | minor note-taking overhead | silently guessing |
| assumption + flip consequence | prepares for redesign follow-ups | slightly longer framing | assumptions with no adaptation path |
| short, categorized log | usable under pressure | omits some lower-value context | exhaustive documentation |

## Interview It

**Google framing:** "Design a distributed lock service. What assumptions would you record before you commit to consistency and failover choices?"

**Cloudflare framing:** "Design customer-configurable edge rules. Which assumptions must be explicit if propagation guarantees are underspecified?"

**Follow-ups:**
1. Which assumption is most expensive to reverse later?
2. Which assumption most changes your data model?
3. Which assumption would you test first with the interviewer?
4. What if the interviewer corrects one of your assumptions mid-design?
5. How do you avoid turning the assumption log into a transcript?

## Ship It

- `outputs/assumption-log-template.md`
- `outputs/interview-card-assumption-logging.md`

## Exercises

1. **Easy** — Write three assumptions for a URL shortener and note which one matters most.  
2. **Medium** — For a chat system, log assumptions about delivery guarantees, retention, and region count, then explain the flip impact.  
3. **Hard** — For a global API platform, create an assumption log that distinguishes reversible and expensive-to-reverse assumptions.  

## Further Reading

- [Google SRE workbook](https://sre.google/workbook/table-of-contents/) — helpful mindset for making operational assumptions explicit
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline on moving forward under reasonable assumptions
