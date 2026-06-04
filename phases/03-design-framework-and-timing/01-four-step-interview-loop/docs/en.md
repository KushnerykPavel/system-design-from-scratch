# Four-Step Interview Loop

> Senior answers do not start with boxes. They start with scope, numbers, and intent.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Internalize a repeatable 45-minute system design loop that starts with clarification, forces sizing, and leaves room for deep dive and wrap-up.  
**Prerequisites:** `01-clarification-and-scope/05-prioritization`, `02-estimation-and-cost/01-qps-and-request-mix`  
**Estimated time:** ~60 min  
**Primary artifact:** interview card + design review prompt  

## The Problem

Many candidates know system components but still perform poorly because they spend too much time drawing architecture before agreeing on the problem. They skip clarification, skip sizing, and then run out of time before trade-offs or failure handling.

This lesson gives you a default interview operating system: a 4-step loop inspired by classic system design prep frameworks, but upgraded for senior-level interviews where trade-offs, reliability, and communication matter as much as component choice.

## Clarify

- Who are the users and what is the single most important user journey?
- What scale are we designing for in v1 and at 10x?
- Which non-functional requirements dominate the design: latency, availability, consistency, abuse resistance, or cost?
- What can we explicitly defer out of scope?

If the interviewer is vague, say your assumptions out loud and continue.

## Requirements

### Functional

- Define the core user flow.
- State one or two explicitly out-of-scope features.

### Non-functional

- Prioritize latency, availability, consistency, durability, cost, and security in order.
- Name the one requirement that is most likely to change the design.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak QPS | rough order of magnitude | affects stateless tier, cache pressure, and queue sizing |
| Storage growth | daily + yearly | affects storage model and retention strategy |
| Bandwidth | response size x peak QPS | affects egress and edge/cache decisions |
| Peak factor | normal vs burst | affects autoscaling and overload handling |
| Rough cost | one sentence estimate | keeps the design realistic |

## Architecture

Use this pacing:

1. **Clarify and prioritize**  
2. **Size quickly**  
3. **Propose high-level design and get buy-in**  
4. **Deep dive into the critical area**  
5. **Wrap up with risks, trade-offs, and next steps**

Suggested 45-minute split:

```text
0-7   clarify + priorities
7-12  rough sizing
12-24 high-level architecture
24-38 deep dive on 1-2 critical areas
38-45 failure modes, observability, trade-offs, wrap-up
```

## Data Model & APIs

This lesson's code artifact is a tiny planner that validates whether your interview plan fits into a 45-minute loop. It is intentionally small, because the main value here is the process discipline.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Candidate rushes to architecture | no clarified assumptions recorded | stop and restate the prompt first |
| Candidate skips sizing | no QPS/storage numbers before diagram | force at least directional estimates |
| Candidate deep dives on the wrong thing | deep dive does not connect to the main constraint | pause and choose the critical path explicitly |
| Candidate runs out of time | wrap-up never happens | reserve the last 7 minutes for failure modes and trade-offs |

## Observability

- metric: minutes spent per interview phase
- metric: number of explicit assumptions recorded
- metric: number of trade-offs stated
- SLO: every answer should reach failure modes and wrap-up before time expires

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| fixed 45-minute loop | keeps answers structured | can feel rigid | free-form discussion that often drifts |
| mandatory sizing before design | keeps design grounded | adds early pressure | architecture-first answers that become hand-wavy |
| one deliberate deep dive | shows prioritization | less breadth in one round | shallow commentary on many components |

## Interview It

**Google framing:** "Design a quota platform. Talk me through your process." The signal is not only the final design, but how you clarify the prompt and decide where to go deep.

**Cloudflare framing:** "Design a global API gateway." The signal is whether you preserve time for edge-specific failure modes, cache semantics, and origin protection rather than spending the entire interview on request routing.

**Follow-ups:**
1. What if the interviewer interrupts earlier than expected?
2. What if the numbers are unclear and you must assume a range?
3. What if the system has multiple equally critical deep dives?
4. What if the interviewer changes scale halfway through?
5. How do you recover when you realize your first design choice was weak?

## Ship It

- `outputs/interview-card-four-step-loop.md`
- `outputs/design-review-four-step-loop.md`

## Exercises

1. **Easy** — Time-box a 20-minute answer for a URL shortener using the same stages.  
2. **Medium** — Adapt the loop for a prompt where compliance is the dominant constraint.  
3. **Hard** — Run the loop on a low-latency trading system and justify a different deep-dive choice.  

## Further Reading

- [System design framework chapter in liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes/tree/main/03.%20System%20Design%20Framework) — baseline interview flow  
- [Google system design interviews](https://sre.google/) — good context for reliability and communication expectations  
