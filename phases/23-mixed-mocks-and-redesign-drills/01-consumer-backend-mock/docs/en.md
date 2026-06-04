# Mixed Mock: Consumer Backend

> Consumer backend interviews are won by showing which user promise drives the design, not by drawing every familiar component.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Practice a timed mixed mock for consumer-style backend prompts where ranking, fanout, freshness, and product behavior all compete for answer time.  
**Prerequisites:** `16-application-backends/07-backend-product-drill`, `21-google-senior-staff-system-design/07-google-full-mock`, `22-cloudflare-edge-platform-design/07-cloudflare-full-mock`  
**Estimated time:** ~90 min  
**Primary artifact:** consumer backend mock scorecard validator + practice card  

## The Problem

Run a full mock on a consumer backend prompt that could plausibly show up in a senior interview loop:

- design a global photo-sharing feed
- design a high-scale group chat backend
- design a notification orchestration service
- design a collaborative whiteboard session backend

The challenge is not naming storage, caches, and queues. It is showing that you can pick the dominant user promise, size the hot path, and choose one deep dive that actually explains the hard part of the product.

## Clarify

- Is the product optimizing for freshness, delivery confidence, ranking quality, or low write cost?
- What is the main skew: celebrity fanout, reconnect bursts, write hotspots, or notification spikes?
- Which user-visible failure is least acceptable: message loss, stale timeline, duplicate notification, or session lag?
- Is the interviewer expecting a single-region answer first or multi-region from the start?

## Requirements

### Functional

- Bound the prompt around one concrete user journey.
- Produce rough traffic and storage sizing before choosing architecture.
- Choose one deep dive tied to the real product risk.
- Explain degraded behavior, observability, and redesign.

### Non-functional

- Keep the answer product-aware instead of generic.
- Make fanout, freshness, and correctness trade-offs explicit.
- Avoid hiding operational risk behind vague "eventual consistency" claims.
- Stay structured under follow-up pressure.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification window | 3 to 5 min | defines the product promise before architecture starts |
| Peak traffic | 500K to 3M active users in a hot hour | shapes fanout, cache pressure, and write amplification |
| Skew factor | top 0.1% producers create 10x to 100x more fanout | determines whether the design breaks under celebrity or group spikes |
| Latency target | sub-second user-visible updates for the hot path | forces prioritization of sync versus async work |
| Retention and storage | weeks to years depending on prompt | changes storage layout, compaction, and cost |

## Architecture

Use a repeatable loop:

```text
clarify user promise
  -> size read/write mix and skew
  -> propose high-level architecture
  -> choose one deep dive:
       fanout strategy
       consistency boundary
       session state
       ranking pipeline
  -> explain failure modes and observability
  -> redesign after a changed constraint
```

Strong answers usually separate:

1. user-facing write path
2. read-serving path
3. async distribution or enrichment path
4. state that must be strongly correct versus state that can lag

## Data Model & APIs

Useful answer skeleton:

```text
consumer_mock_answer(
  prompt,
  top_user_journey,
  workload_shape,
  write_path,
  read_path,
  deep_dive,
  degraded_mode,
  observability,
  redesign
)
```

Scorecard fields:

- prompt clarification quality
- sizing quality
- read/write asymmetry awareness
- deep-dive quality
- degraded-mode clarity
- observability quality
- redesign quality

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| answer treats every consumer product like feed plus cache plus queue | no product contract or user-visible guarantee appears | force one primary user promise and one unacceptable failure |
| fanout or skew is ignored | no hotspot discussion appears despite large producers or groups | estimate amplification and pick a skew-aware deep dive |
| correctness boundary is muddy | duplicate, stale, or missing state is discussed vaguely | name source of truth and tolerated anomaly explicitly |
| answer never reaches operational behavior | no degraded mode or metrics appear | reserve time for failure and observability before wrap-up |

## Observability

- metric: latency and freshness of the main user journey
- metric: fanout amplification or queue backlog on the dominant async path
- metric: duplicate, drop, or stale-read rate for the core product promise
- log: chosen assumptions, degraded-mode rules, and redesign triggers
- trace: user action through write path, async path, and read path
- SLO: define one measurable promise that a user would actually notice

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one concrete user journey | coherent architecture and metrics | less breadth in the initial answer | solving every product surface at once |
| one deep dive tied to the product bottleneck | high-signal reasoning | fewer peripheral details | many shallow dives |
| explicit tolerated anomaly | honest trade-off story | forces product judgment | pretending all paths are equally strict |

## Interview It

**Google framing:** Expect pressure on workload shape, system boundaries, and which one or two guarantees are truly worth paying for.

**Cloudflare framing:** Expect follow-ups on incident behavior, hotspot containment, and how regional or network conditions change user-visible freshness.

**Suggested prompts:**
1. Design a high-scale social feed.
2. Design group chat for millions of active rooms.
3. Design a notification orchestration backend.
4. Design a collaborative whiteboard backend.

## Ship It

- `outputs/interview-card-consumer-backend-mock.md`
- `outputs/skill-consumer-backend-mock.md`

## Exercises

1. **Easy** — Pick one prompt and write only the clarification and sizing sections.
2. **Medium** — For the same prompt, choose one deep dive and list three failure signals.
3. **Hard** — Re-answer the prompt after the interviewer changes the top goal from freshness to cost control.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — useful for turning product architectures into operational follow-ups  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — strong reference for state, streams, and consistency boundaries  
