# Search and Monitoring Drill

> Senior performance in interviews comes from choosing one good design path under pressure, not from naming every component you know.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Practice full-loop system design across search and monitoring prompts by forcing explicit clarification, sizing, deep-dive selection, failure analysis, and redesign.
**Prerequisites:** `03-design-framework-and-timing/08-full-loop-drill`, `17-search-crawl-and-monitoring-systems/03-metrics-platform`, `17-search-crawl-and-monitoring-systems/06-index-freshness`
**Estimated time:** ~60 min
**Primary artifact:** drill scorecard validator + practice skill sheet

## The Problem

Run a timed drill that combines the phase themes. You might be asked to design a website monitoring platform, a search freshness pipeline, or an operational telemetry backend, then defend the design under changing constraints.

This lesson exists to convert theory into interview rhythm. The learner should feel the pressure to clarify first, do rough math, pick one deep dive, and speak concretely about failure modes and observability.

## Clarify

- Is the prompt more search-heavy, monitoring-heavy, or mixed?
- What is the most important axis: freshness, latency, coverage, cost, or signal quality?
- Which parts are in scope for the first 20 minutes versus follow-up rounds?
- Are we optimizing for internal platform operators, end users, or both?

If the interviewer stays ambiguous, make your assumptions explicit and proceed instead of stalling.

## Requirements

### Functional

- Clarify the prompt and state assumptions.
- Produce rough sizing before architecture.
- Present a high-level design and one deliberate deep dive.
- Explain failure modes, observability, and rollout or migration.
- Redesign after one changed constraint.

### Non-functional

- Stay organized under time pressure.
- Avoid sprawling component catalogs without justification.
- Surface the most important trade-offs early.
- Keep the answer adaptable when the interviewer changes assumptions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification time | 3 to 5 minutes | sets scope before premature architecture |
| Sizing time | 5 to 7 minutes | anchors storage, QPS, and freshness claims |
| Initial design time | 10 to 12 minutes | enough to get buy-in before deep dive |
| Deep dive time | 8 to 10 minutes | proves senior judgment in choosing complexity |
| Redesign time | 5 minutes | tests adaptability under changed constraints |

## Architecture

Use the four-step rhythm:

```text
clarify scope
  -> size the workload
  -> propose high-level architecture
  -> deep dive one critical area
  -> close with trade-offs, failure modes, and redesign
```

Recommended deep dives for this phase:

- crawler frontier and politeness
- autocomplete freshness and ranking
- metrics cardinality and retention
- alert routing quality
- index update rollout safety

## Data Model & APIs

For the drill, the important API is your answer structure:

- clarification assumptions
- capacity sheet
- core component boundaries
- one deep-dive interface
- observability and SLO plan
- redesign trigger

If you cannot explain one key data model cleanly, your design is probably still too vague.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate skips sizing and hand-waves scale | interviewer notices ungrounded claims | force a quick rough estimate before architecture |
| answer fans out into too many deep dives | time runs out and no part is defended well | pick one high-value deep dive deliberately |
| failure modes are generic | no metrics or alerts tied to them | require one detection metric per major failure |
| redesign answer is superficial | changed constraint has no architectural consequence | state which component or policy actually changes |

## Observability

- metric: did the answer name concrete latency, freshness, or retention targets?
- metric: did each major failure mode include a detection signal?
- metric: was at least one trade-off expressed as benefit and cost, not only preference?
- log: note where the candidate made or revised assumptions
- trace: track the flow of the interview answer from scope to redesign
- SLO: complete a coherent first-pass design with one deep dive inside the allotted interview time

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one chosen deep dive | depth where it matters | less coverage elsewhere | shallow discussion of every component |
| rough sizing before architecture | grounded design decisions | spends a few minutes up front | intuitive but unmeasured capacity claims |
| explicit assumptions | keeps momentum under ambiguity | may need revision later | waiting for perfect prompt precision |

## Interview It

**Google framing:** "Design a metrics platform for internal services, then explain what breaks when one team adds extreme label cardinality." Expect pushback on guardrails, query isolation, and operator usability.

**Cloudflare framing:** "Design website monitoring with global probes and alert routing." Expect follow-ups on probe placement, noisy checks, routing, and how the system behaves during wide regional incidents.

**Follow-ups:**
1. What if cost becomes more important than freshness?
2. What if the user asks for a 10x increase in query or ingest rate?
3. What if one region must operate independently during control-plane loss?
4. What if policy or privacy constraints change the storage plan?
5. What if the interviewer asks you to cut scope and ship phase one quickly?

## Ship It

- `outputs/skill-search-monitoring-drill.md`

## Exercises

1. **Easy** — Give a three-minute clarification and sizing opening for "design website monitoring."
2. **Medium** — Choose one deep dive for "design autocomplete" and justify why it is the right one.
3. **Hard** — Redesign a metrics platform answer after the interviewer introduces multi-tenant abuse and a lower cost target.

## Further Reading

- [System Design Interview - An Insider's Guide](https://github.com/liquidslr/system-design-notes) — useful reminder of the four-step interview flow
- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — strong source of operational follow-up pressure
