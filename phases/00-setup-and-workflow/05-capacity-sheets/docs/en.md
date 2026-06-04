# How to Use Capacity Sheets

> Rough numbers beat elegant guesses.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Use a compact capacity sheet to anchor system design conversations in throughput, storage, bandwidth, burstiness, and rough cost.
**Prerequisites:** `00-setup-and-workflow/01-repo-setup-and-progress`
**Estimated time:** ~60 min
**Primary artifact:** capacity sheet + estimator helper

## The Problem

Candidates often say they know sizing matters, then still jump straight to architecture. Capacity sheets are a forcing function: write down the few numbers that actually shape the design, then let those numbers constrain your choices.

This lesson is about speed, not spreadsheet perfection. You want interview-grade estimates that make topology and trade-offs credible.

## Clarify

- What is the dominant user action, and what percentage of traffic does it represent?
- What values can be estimated quickly from first principles: DAU, requests per user, object size, retention?
- Which number is most likely to change the topology if it is 10x larger than expected?

## Requirements

### Functional

- Estimate peak QPS, reads and writes per day, storage growth, and bandwidth.
- Capture a peak factor so average load does not hide burst risk.
- Preserve assumptions and units so the math can be reviewed quickly.

### Non-functional

- A first-pass sheet should take under five minutes.
- Math should be easy to audit by another person.
- The sheet should stay useful even when the inputs are uncertain.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak QPS | derived from DAU, actions, and burst factor | drives stateless fleet sizing and cache pressure |
| Storage growth | object size times writes times retention | shapes database and blob strategy |
| Bandwidth | bytes per request times QPS | drives edge, cache, and egress choices |
| Peak factor | average to peak multiplier | exposes overload and autoscaling risk |
| Rough cost | one or two dominant resource lines | keeps answers grounded in reality |

## Architecture

A good capacity sheet feeds directly into architecture:

- high read QPS suggests cache layers and read scaling
- large object size suggests blob separation and CDN use
- long retention changes storage tiering
- high burst factor raises admission control and queue questions

The helper tool in this lesson calculates a few common metrics from compact inputs so you can practice translating assumptions into usable numbers.

## Data Model & APIs

Core inputs:

- daily active users
- requests per active user
- peak factor
- average payload bytes
- write percentage
- retention days

Derived outputs:

- average QPS
- peak QPS
- daily writes
- daily storage bytes
- annual storage bytes
- peak bandwidth bytes per second

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hand-wavy numbers | no units or formulas recorded | require units and assumptions for each row |
| average load hides reality | no peak factor present | always include burst multiplier |
| one metric dominates everything | storage or bandwidth dwarfs the rest unnoticed | call out the biggest bottleneck explicitly |
| false precision | numbers quoted to too many digits | round aggressively and state uncertainty |

## Observability

- metric: time to produce first-pass capacity sheet
- metric: percentage of design notes with explicit assumptions and units
- metric: number of later architectural changes caused by corrected sizing
- SLO: every full design answer should produce rough sizing before the architecture deep dive

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| compact capacity sheet | fast and reusable | omits second-order effects initially | exhaustive spreadsheet |
| aggressive rounding | easier communication | less precise totals | exact-looking but fragile estimates |
| derive only a few key outputs | keeps focus on topology decisions | less detailed cost breakdown | broad forecasting model |

## Interview It

**Google framing:** "Show me the first numbers you would estimate for a photo-sharing backend and explain how they shape your architecture."

**Cloudflare framing:** "For a global API protection service, what sizing numbers would you compute before discussing edge caches, queues, or storage?"

**Follow-ups:**
1. What if DAU is unknown but MAU is given?
2. How do you model a heavy burst window like the start of a sports event?
3. Which assumptions matter most for a write-heavy system?
4. How do you express uncertainty without sounding unprepared?
5. What changes if only 1 percent of requests carry large payloads?

## Ship It

- `outputs/capacity-sheet-template.md`
- `outputs/observability-checklist-capacity-sheets.md`

## Exercises

1. **Easy** — Build a capacity sheet for a tiny internal URL shortener.
2. **Medium** — Recompute the sheet when payload size increases by 20x but request count stays flat.
3. **Hard** — Compare two systems with identical QPS but radically different retention and payload sizes.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — helpful reminder that performance behavior changes under scale and fan-out
- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) — useful baseline for rough sizing discipline
