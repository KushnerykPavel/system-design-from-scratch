# User Journeys and Workload Shape

> The architecture follows the traffic shape, not the feature list.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to turn the core user journey into an explicit workload shape so you can size and design for reads, writes, bursts, fanout, and hot paths instead of vague "high scale."  
**Prerequisites:** `01-clarification-and-scope/02-high-leverage-questions`, `02-estimation-and-cost/01-qps-and-request-mix`  
**Estimated time:** ~60 min  
**Primary artifact:** workload worksheet + interview card  

## The Problem

Many candidates say "design for high scale" without defining what kind of scale:

- mostly reads or mostly writes?
- bursty or steady?
- tiny payloads or large objects?
- one request triggers one write, or one request triggers fanout to millions?

The same product label can hide very different systems. A chat product, a feed product, and a metrics product each have distinct workload shapes even when their user counts are similar.

This lesson trains the habit of converting user journeys into traffic geometry before architecture starts.

## Clarify

- What is the primary user journey and what operations does it create on the backend?
- Is the system dominated by reads, writes, scans, fanout, or asynchronous pipelines?
- Where are the bursts, hotspots, or amplification effects likely to appear?

## Requirements

### Functional

- Map the main user journey to concrete backend operations.
- Identify at least one hot path and one secondary path.
- Distinguish interactive operations from background work.

### Non-functional

- Quantify the rough request mix and amplification factors.
- Highlight burstiness and hotspot risk, not just average load.
- Use the workload shape to justify early architecture choices.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read/write mix | e.g. 90/10, 50/50, or 10/90 | drives cache utility and write-path complexity |
| Fanout factor | 1x to 1M+ depending on product | changes queueing, delivery, and storage design |
| Payload size | bytes to MBs | separates metadata paths from blob paths |
| Peak factor | 3x to 20x or higher | determines burst handling and admission control |
| Rough cost | storage, egress, and compute all depend on path shape | stops the design from optimizing the wrong resource |

## Architecture

Translate journey into workload shape with a simple template:

1. user action
2. backend operations created
3. synchronous versus asynchronous split
4. amplification or fanout
5. hotspot or skew risk

Examples:

- posting to a social feed may be one write followed by massive fanout
- object upload may be modest request rate but very large bandwidth and storage
- metrics ingest may be append-heavy, bursty, and aggregation-centric

The code artifact summarizes a workload profile and flags when the stated user journey does not include any hot path or burst characterization.

## Data Model & APIs

Represent a workload with:

- `journey`
- `read_qps`
- `write_qps`
- `fanout`
- `burst_factor`
- `hot_key_risk`

Useful review prompts:

- Which path dominates CPU, storage, or network?
- Which user action creates asynchronous follow-up work?
- What is the worst hotspot if one tenant goes viral?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| average-only thinking | no mention of peaks, bursts, or hot tenants | estimate peak factor and skew explicitly |
| journey not translated into operations | architecture appears before request mix is stated | list backend reads, writes, fanout, and async tasks first |
| interactive and async paths blurred together | latency goals become confusing | separate user-facing and background paths |
| hotspot risk ignored | later design cannot explain load concentration | ask where traffic or keys cluster unnaturally |

## Observability

- metric: explicit read/write mix recorded before architecture begins
- metric: fanout or amplification factor noted for core workflow
- metric: number of hotspots or burst sources identified
- SLO: the primary workload shape should be describable in two or three sentences before drawing the first architecture box

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| journey-first workload model | architecture reflects real traffic | requires a little up-front translation work | generic "high scale" language |
| highlight peaks and skew | exposes operational risk early | rough numbers may be approximate | average-only sizing |
| separate sync and async paths | cleaner latency reasoning | more concepts to track | one blended path that hides amplification |

## Interview It

**Google framing:** "Design YouTube comments. What workload shape matters more: average QPS or viral hotspot behavior?"

**Cloudflare framing:** "Design a global logs ingest product. Walk from customer journey to request mix, bursts, and where backpressure shows up."

**Follow-ups:**
1. Which path is hot and user-facing versus hot and background?
2. What part of the workload most changes the storage choice?
3. Which amplification factor is easy to forget but expensive in production?
4. How would a celebrity or large tenant distort the average case?
5. What if the system becomes write-heavy during incident conditions?

## Ship It

- `outputs/workload-worksheet.md`
- `outputs/interview-card-workload-shape.md`

## Exercises

1. **Easy** — Turn a URL shortener prompt into a read/write mix and burst estimate.  
2. **Medium** — Compare the workload shape of chat delivery versus feed fanout.  
3. **Hard** — Model a global metrics ingest service and explain which path dominates cost and reliability risk.  

## Further Reading

- [Google SRE workbook](https://sre.google/workbook/table-of-contents/) — useful framing for operational load and failure behavior
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline examples of workload-driven capacity estimation
