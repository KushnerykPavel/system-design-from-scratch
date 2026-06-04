# QPS and Request Mix Estimation

> A sizing answer becomes useful the moment it changes the shape of the system.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Turn vague product usage into directional read/write QPS that can drive topology, cache, and storage choices.  
**Prerequisites:** `01-clarification-and-scope/06-workload-shape`, `01-clarification-and-scope/05-prioritization`  
**Estimated time:** ~75 min  
**Primary artifact:** interview card + capacity sheet  

## The Problem

Candidates often say "high scale" without translating it into request volume. That makes every later choice hand-wavy: stateless tier size, cache pressure, database write load, and regional fanout all depend on rough QPS and request mix.

This lesson teaches a repeatable way to estimate baseline, peak, read, and write traffic from incomplete inputs.

## Clarify

- Are we estimating global QPS, per-region QPS, or per-shard QPS?
- Is the product read-heavy, write-heavy, or dominated by background jobs?
- Do peak hours line up globally or only inside one region?
- Do retries, cache misses, and fanout count toward origin load or only user-visible requests?

## Requirements

### Functional

- Produce baseline QPS from users, sessions, or events.
- Split traffic into reads, writes, and asynchronous follow-up work.
- Surface peak factors explicitly instead of hiding them in one average number.

### Non-functional

- Estimates must be fast enough to say out loud in an interview.
- Numbers should be directionally correct, not finance-grade exact.
- The method must expose which assumptions most affect the architecture.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Daily active users | 5M | anchors the traffic model |
| Requests per active user per day | 24 | drives baseline request count |
| Baseline QPS | about 1,400 | sets average fleet load |
| Peak factor | 6x | drives autoscaling and queueing risk |
| Peak read/write QPS | 6,700 read / 1,700 write | changes cache and database planning |

## Architecture

Use a short ladder:

1. Pick one traffic anchor: DAU, API calls, uploads, or transactions.
2. Convert the daily or hourly number into average QPS.
3. Apply a peak factor separately.
4. Split by request mix: read, write, async, internal fanout.
5. Say which subsystem feels the peak first.

For example, a 5M DAU social product with 24 requests per user per day yields:

- 120M requests/day
- about 1,389 average QPS
- about 8,333 peak QPS at a 6x factor
- if traffic is 80% reads, then about 6,667 read QPS and 1,667 write QPS at peak

## Data Model & APIs

The code artifact models one traffic estimate:

```text
TrafficModel {
  DAU
  RequestsPerUserPerDay
  ReadRatio
  PeakFactor
}
```

Useful derived outputs:

- average QPS
- peak QPS
- peak read QPS
- peak write QPS

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| averages hide the real peak | traffic spikes exceed provisioned headroom | always state a peak factor separately |
| internal fanout ignored | downstream QPS is much higher than ingress QPS | multiply for feed fanout, retries, or cache misses |
| read/write split guessed badly | wrong storage or cache assumptions | say the split as an assumption and show sensitivity |
| global number used for one hot region | one region overloads before fleet average shows pain | convert global traffic into per-region peak |

## Observability

- metric: ingress QPS by route and region
- metric: read vs write ratio over time
- metric: retry amplification factor
- metric: cache miss QPS vs user-visible QPS
- SLO: estimation worksheet should identify the top 2 traffic-sensitive assumptions

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| start with one anchor metric | faster reasoning | less precision | collecting many weak inputs |
| use round numbers | easier communication | some arithmetic loss | pseudo-exact values that slow the answer |
| show separate peak factor | exposes risk honestly | adds another assumption | folding peak into one average number |

## Interview It

**Google framing:** "Design a notification platform. Start with sizing." The signal is whether you separate user-visible sends, retries, and background fanout.

**Cloudflare framing:** "Design an API gateway for a fast-growing SaaS product." The signal is whether you think in regional peak traffic instead of one global average.

**Follow-ups:**
1. What changes if 5% of users produce 60% of traffic?
2. What if reads are cacheable and writes are not?
3. What if one launch event creates a 20x burst for 15 minutes?
4. What if half of the traffic is from automated clients rather than humans?

## Ship It

- `outputs/interview-card-qps-and-request-mix.md`
- `outputs/capacity-sheet-qps-and-request-mix.md`

## Exercises

1. **Easy** — Estimate QPS for a photo upload service with 2M daily uploads.  
2. **Medium** — Rework the model when each write triggers 30 follower fanout events.  
3. **Hard** — Estimate separate regional peaks for a product with US evening concentration.  

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — good baseline on why early sizing matters  
- [Google SRE book](https://sre.google/books/) — useful context for capacity and overload thinking  
