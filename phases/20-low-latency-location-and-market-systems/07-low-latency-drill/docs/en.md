# Low-Latency Systems Drill

> The point of this drill is not to say "optimize latency." It is to defend which path must stay hot, which path may lag, and which guarantee you refuse to compromise.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Practice a full interview loop on low-latency systems by forcing explicit path decomposition, freshness budgets, hotspot handling, and degraded-mode choices.
**Prerequisites:** `20-low-latency-location-and-market-systems/01-proximity-service`, `20-low-latency-location-and-market-systems/04-stock-exchange`, `20-low-latency-location-and-market-systems/06-map-tiles`
**Estimated time:** ~60 min
**Primary artifact:** drill worksheet + scoring rubric

## The Problem

Run a timed drill on prompts such as "design nearby search," "design a real-time leaderboard," "design a market-data feed," or "design a matching engine." Your goal is to make latency, isolation, and fallback decisions visible instead of hiding behind buzzwords.

This lesson exists to turn the entire phase into interview reflex. You should identify the hot path, state what can be async, size the skew, and explain what happens when freshness or fanout falls behind.

## Clarify

- What is the one user-visible latency promise that matters most?
- Which writes or updates feed that promise?
- Is exactness or freshness the harder constraint?
- Which deep dive would reveal the real engineering judgment?

If the prompt stays broad, pick one hot path, one correctness boundary, and one degraded mode. Answers that try to optimize everything at once usually become vague.

## Requirements

### Functional

- Clarify the primary user-visible low-latency path.
- Size steady-state load and skew or burst risk.
- Separate synchronous and asynchronous responsibilities.
- Explain failure handling and degraded behavior explicitly.
- Redesign after a changed requirement without hiding trade-offs.

### Non-functional

- Keep the answer concrete under time pressure.
- Avoid magical global ordering or perfect freshness claims.
- Show hotspot and noisy-neighbor thinking.
- Make observability part of the low-latency story.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification time | 3 to 5 minutes | defines the hot path before architecture sprawl starts |
| Sizing time | 5 to 7 minutes | skew and candidate/fanout growth must be quantified |
| Initial design time | 10 to 12 minutes | enough to establish sync versus async boundaries |
| Deep dive time | 8 to 10 minutes | where latency or correctness trade-offs become visible |
| Redesign time | 5 minutes | proves the design is grounded in real constraints |

## Architecture

Use the same loop every time:

```text
clarify the hot path and latency/freshness target
  -> size steady-state plus hotspot or fanout skew
  -> define sync versus async boundaries
  -> deep dive one pressure point
  -> close with degraded mode, observability, and redesign
```

Good deep-dive choices for this phase:

- hot-cell mitigation for proximity search
- stale-update policy for location ingestion
- top-N versus around-me rank path
- sequencing and ack boundaries in matching
- subscriber isolation in market-data fanout
- version rollout and cache invalidation for map tiles

## Data Model & APIs

For the drill, your answer structure is the main interface:

- hot path statement
- freshness or exactness boundary
- one core state record or partitioning scheme
- one detection metric for each major failure
- one degraded mode the product can survive

If you cannot say what becomes asynchronous, the design is probably still too hand-wavy.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate optimizes everything equally | no clear hot path in first minutes | force one primary latency promise |
| skew is ignored | no hotspot estimate or partition plan | add one explicit worst-case density or fanout estimate |
| async boundaries are fuzzy | no clear ack or publish contract | state what must happen before success is returned |
| degraded mode is missing | failures end with "add retries and replicas" | name what becomes slower, stale, or unavailable |

## Observability

- metric: did the answer define the hot path and target clearly?
- metric: did the sizing include skew, retries, or fanout amplification?
- metric: did each failure mode name a detection signal?
- log: note where the design allowed staleness, delay, or approximation
- trace: clarify -> size -> architecture -> deep dive -> redesign
- SLO: produce a coherent low-latency design answer inside interview time without hiding the key trade-offs

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one explicit hot path | focused credible design | less breadth elsewhere | vague optimize-everything answer |
| one deliberate deep dive | reveals senior judgment | less time for side topics | shallow coverage of all components |
| honest degraded mode | operational credibility | fewer flashy guarantees | magical always-fresh always-fast story |

## Interview It

**Google framing:** "Design a low-latency user-facing or infrastructure-heavy system and explain which path stays hot under scale." Expect follow-ups on skew, partitioning, and exactness versus freshness.

**Cloudflare framing:** "Design a globally distributed low-latency service with strong isolation and graceful degradation." Expect pressure on edge behavior, noisy neighbors, and control-plane safety.

**Follow-ups:**
1. What if traffic grows 10x but only in one region or one symbol?
2. What if cost pressure forces more approximation or more cache dependence?
3. What if the user-visible latency target stays fixed but write freshness worsens?
4. What if the control plane is healthy but the data plane is saturated?
5. What is the first guarantee you would relax, and why?

## Ship It

- `outputs/low-latency-drill-sheet.md`
- `outputs/scoring-rubric-low-latency-drill.md`

## Exercises

1. **Easy** — Give a three-minute opening for "design nearby search."
2. **Medium** — Pick the best deep dive for "design a market-data feed" and justify it.
3. **Hard** — Re-answer "design a matching engine" after the interviewer adds a much hotter top symbol and stricter durability expectations.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — strong framing for low-latency design pressure
- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) — useful reminder of the four-step interview flow
