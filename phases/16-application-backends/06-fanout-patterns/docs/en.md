# Application Fanout Patterns Compared

> Push, pull, and hybrid fanout are workload decisions, not architecture ideology.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Compare push, pull, and mixed fanout patterns across feeds, chat, notifications, and collaboration workloads, using concrete workload signals instead of slogans.  
**Prerequisites:** `16-application-backends/02-news-feed`, `16-application-backends/03-chat-system`, `16-application-backends/04-notification-system`  
**Estimated time:** ~60 min  
**Primary artifact:** fanout decision matrix generator  

## The Problem

Many application backends boil down to deciding where work happens: at write time, at read time, or in a layered combination. Candidates often memorize "fanout on write for feeds" or "fanout on read for celebrities" without extracting the deeper rule.

This lesson builds a reusable mental model for choosing fanout strategy based on read skew, freshness targets, recipient counts, and tolerance for degraded partial views.

## Clarify

- Is the product dominated by writes, reads, or both?
- Are recipients known and bounded at write time?
- Can the system serve partial or stale results during downstream trouble?

If unspecified, assume the product has skew, peaks, and at least one high-value latency target that makes a uniform strategy suboptimal.

## Requirements

### Functional

- Select a fanout pattern for a given product shape.
- Explain when strategy should vary by cohort or object class.
- Identify how the chosen pattern affects caching, storage, and queueing.

### Non-functional

- Keep the explanation transferable across multiple backend products.
- Make hidden write amplification and read amplification visible.
- Tie strategy choices to concrete observability signals.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Recipient count per event | 1 to millions | defines fanout amplification |
| Read-to-write ratio | 0.1x to 100x | indicates where precomputation pays off |
| Freshness target | sub-second to minutes | changes what can be deferred |
| Skew | mild to celebrity-grade | determines need for hybrid patterns |
| Cost sensitivity | strict or relaxed | affects whether duplicate compute is acceptable |

## Architecture

Use this evaluation framework:

```text
event occurs
  -> classify audience size, urgency, skew, and read frequency
  -> choose:
     push fanout
     pull fanout
     mixed fanout
     two-stage fanout with buffering
  -> attach fallback and observability plan
```

Examples:

- Feeds: often mixed because authors vary widely in audience size.
- Chat: often push to online recipients, pull for offline replay.
- Notifications: push through orchestrators, but policy may suppress many deliveries.
- Collaboration: push live deltas, pull snapshots and replay on reconnect.

## Data Model & APIs

This lesson is comparative, so the reusable API is conceptual:

- `EvaluateFanout(workload_shape)`
- `ScoreAmplification(reads_per_write, recipients_per_event)`
- `RecommendFallback(dependency_failure_mode)`

The real senior signal is not API syntax. It is whether the candidate can choose the right work placement and explain the cost.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| chosen push pattern explodes under skew | write-path queue lag and per-object amplification | cohort split or hybrid strategy |
| chosen pull pattern overloads read path | read p99 and backend fan-in depth | partial materialization or caching |
| mixed strategy becomes impossible to reason about | decision drift across teams | explicit policy rules and ownership |
| no fallback for delayed computation | user-facing outages when secondary pipeline lags | degraded read mode and freshness indicators |

## Observability

- metric: write amplification per logical event
- metric: read amplification or merge depth per user request
- metric: freshness lag between source event and user-visible effect
- metric: queue lag or cache miss burst for each fanout cohort
- log: strategy decision for sampled workload classes
- trace: source event through fanout path to visible user outcome
- SLO: user-visible latency and freshness targets should be tied to the chosen work-placement strategy

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| push fanout | fast reads and predictable recipient retrieval | expensive writes under large audiences | pure pull for hot-read workloads |
| pull fanout | cheap writes and flexible ranking | more expensive reads | pure push for celebrity-scale audiences |
| mixed fanout | workload-aware efficiency | more policy complexity | single global strategy |

## Interview It

**Google framing:** "Compare backend work-placement strategies across several consumer products." Expect pushback if your answer stays at buzzwords instead of tying strategy to workload shape.

**Cloudflare framing:** "Explain when data should be pushed close to users versus assembled on demand." Expect questions on cache locality, edge fanout, and degraded behavior under dependency lag.

**Follow-ups:**
1. Which metrics would tell you the current fanout choice is wrong?
2. What if product semantics change from chronological to ranked?
3. How do you roll out a push-to-hybrid migration safely?
4. What if one region sees 80% of reads but only 20% of writes?
5. How do you keep different teams from choosing conflicting strategies?

## Ship It

- `outputs/interview-card-fanout-patterns.md`

## Exercises

1. **Easy** — Choose the default fanout mode for notifications and justify it.
2. **Medium** — Compare chat delivery for online versus offline devices.
3. **Hard** — Redesign a pure push news feed into a mixed strategy without dropping freshness SLOs.

## Further Reading

- [Twitter Timelines at Scale](https://www.infoq.com/presentations/Twitter-Timeline-Scalability/) — concrete fanout trade-offs in a famous application backend  
- [The Log](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying) — useful for reasoning about push, replay, and catch-up  
