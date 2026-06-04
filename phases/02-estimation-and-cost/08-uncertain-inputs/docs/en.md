# Estimation Under Uncertain Inputs

> Senior sizing does not wait for perfect numbers. It shows a range, names the swing assumptions, and keeps moving.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Handle incomplete or noisy interview inputs by building low/base/high estimates and showing which assumptions matter most.  
**Prerequisites:** `02-estimation-and-cost/01-qps-and-request-mix`, `03-design-framework-and-timing/01-four-step-interview-loop`  
**Estimated time:** ~60 min  
**Primary artifact:** uncertainty worksheet + interview card  

## The Problem

Interviewers often omit exact traffic, payload, or retention numbers on purpose. Weak candidates stall or pretend certainty. Strong candidates define a plausible range and explain how the architecture changes across that range.

This lesson teaches range-based estimation without losing decision-making clarity.

## Clarify

- Which input is most uncertain: users, payload size, peak factor, or retention?
- Which assumption would actually change the design if it moved 10x?
- Can we choose a base case and bound it with low and high scenarios?
- What decision must be made now despite uncertainty?

## Requirements

### Functional

- Build low, base, and high estimates quickly.
- Identify the top uncertainty driver.
- Explain which design choices are robust across the range.

### Non-functional

- Avoid analysis paralysis.
- Keep uncertainty explicit rather than hidden.
- Stay decisive enough to keep the interview moving.

## Capacity Model

| Dimension | Low | Base | High | Why it matters |
|-----------|-----|------|------|----------------|
| Peak QPS | 20K | 60K | 180K | controls stateless and cache sizing |
| Avg response size | 5 KB | 20 KB | 80 KB | changes bandwidth pressure |
| Retention | 7 d | 30 d | 180 d | changes storage shape |
| Monthly egress | bounded range | bounded range | bounded range | tests cost sensitivity |
| Design break point | around high case | | | reveals when architecture changes |

## Architecture

A simple workflow:

1. Pick the uncertain variable.
2. Create low/base/high values.
3. Compute rough impact for each.
4. Decide what architecture holds across the whole range.
5. Name the threshold where a different design is needed.

Example:

- base design may work from 20K to 60K QPS
- 180K QPS may require extra partitioning, better caching, or regionalization

## Data Model & APIs

The code artifact models one range:

```text
Range {
  Low
  Base
  High
}
```

Outputs:

- spread ratio
- whether uncertainty is narrow or wide

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate freezes waiting for exact data | answer stalls early | declare a base case and bounded range |
| false certainty | no assumptions called out | state the biggest uncertainty explicitly |
| no threshold thinking | same design claimed for any scale | identify the break point where architecture changes |
| too many variables ranged at once | discussion becomes unreadable | vary one or two major inputs only |

## Observability

- metric: actual vs forecasted workload range
- metric: assumption error by major sizing input
- metric: frequency of scaling events near predicted thresholds
- SLO: operating plans should define the next architecture breakpoint before it is reached

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| low/base/high ranges | honest under uncertainty | less neat than one number | pretending certainty |
| focus on one major swing factor | clearer discussion | less breadth | ranging every variable at once |
| state breakpoint explicitly | better redesign readiness | requires decisive judgment | vague "we'll scale later" answers |

## Interview It

**Google framing:** "Design a service, but assume product traffic could grow 10x next year." The signal is whether you can distinguish robust choices from choices that break at scale.

**Cloudflare framing:** "Design an edge service where traffic shape is uncertain across regions." The signal is whether you reason in ranges and thresholds instead of one fragile average.

**Follow-ups:**
1. What if the high case arrives on day one because of a launch?
2. What if payload size uncertainty matters more than QPS uncertainty?
3. What if only one region experiences the high case?
4. What if the low case makes a simpler architecture more rational for v1?

## Ship It

- `outputs/uncertainty-worksheet.md`
- `outputs/interview-card-uncertain-inputs.md`

## Exercises

1. **Easy** — Build low/base/high QPS for a new API with weak usage forecasts.  
2. **Medium** — Show which architecture stays stable across a 3x storage range.  
3. **Hard** — Decide where a single-region design should flip to multi-region under uncertainty.  

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline on rough sizing under ambiguity  
- [Google SRE book](https://sre.google/books/) — strong context on planning with imperfect information  
