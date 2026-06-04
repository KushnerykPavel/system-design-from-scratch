# Retry Budgets and Hedging

> Tail-latency tricks are only safe when they spend from a budget instead of from wishful thinking.

**Type:** Build
**Company focus:** Google
**Learning goal:** Use retry budgets and hedged requests to improve tail latency without letting speculative work become an invisible overload tax.
**Prerequisites:** `02-estimation-and-cost/04-cache-hit-rate`, `10-reliability-retries-and-backpressure/01-timeouts-and-retries`, `10-reliability-retries-and-backpressure/04-load-shedding`
**Estimated time:** ~60 min
**Primary artifact:** retry budget evaluator + hedging trade-off matrix

## The Problem

Retries and hedged requests can rescue tail latency, but they spend extra capacity. If the system cannot answer "how much speculative work are we willing to create?" then every latency incident risks turning into an overload incident.

This lesson focuses on two senior-level upgrades:

- give retries and hedges explicit budgets
- cancel speculative work aggressively when one attempt wins

## Clarify

- Are we optimizing for success rate, tail latency, or both?
- Which request classes justify hedging, and which are too expensive?
- Is the backend replicated enough that hedging hits independent failure domains?
- What fraction of total traffic can be spent on speculative attempts safely?

If the interviewer is vague, assume a latency-sensitive read-heavy service where p99 matters, replicas are independent enough to hedge occasionally, and capacity margins are finite.

## Requirements

### Functional

- Cap extra requests created by retries and hedges.
- Improve tail latency on selected request classes.
- Cancel losing speculative work quickly when another attempt succeeds.

### Non-functional

- Keep speculative load from overwhelming healthy dependencies.
- Make budget spend observable by caller, endpoint, and dependency class.
- Avoid masking chronic latency problems with endless extra attempts.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Baseline read QPS | 180K req/s | even a small hedge rate creates large absolute extra load |
| Allowed extra-attempt budget | 5-8% of baseline | keeps latency tricks from consuming all spare headroom |
| p99 tail gap | median `25 ms`, p99 `220 ms` | enough spread to justify selective hedging |
| Replica independence | 2-3 replicas in separate failure domains | hedging only helps if failures are not fully correlated |
| Rough cost | speculative CPU + network + canceled work | tail improvement is purchased with real capacity |

## Architecture

A healthy design separates normal traffic from budgeted extra work:

```text
caller
  -> budget manager
  -> primary attempt
  -> delayed hedge only if request is slow and budget remains
  -> cancel losers immediately on success
```

Good defaults:

- budget as a percentage of primary traffic
- hedge only reads or other safe idempotent work
- trigger hedge after a percentile threshold, not immediately
- disable hedging automatically when budget burn is too high

## Data Model & APIs

Useful policy fields:

- `retry_budget_ratio`
- `hedge_budget_ratio`
- `hedge_after_ms`
- `max_extra_attempts`
- `cancel_on_first_success`

Useful interfaces:

- `CanSpendBudget(class)`
- `ShouldHedge(latency_so_far, budget_remaining)`
- `RecordExtraAttempt(class, type)`

Senior-level detail:

- separate retry budgets from hedge budgets when semantics differ
- hedge only where alternate replicas are meaningfully independent
- cancellation needs to be real, not "we ignore the result but still do the work"

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hedge traffic doubles load during a broad incident | hedge rate and backend saturation rise together | suspend hedging when budget burn or saturation exceeds threshold |
| retries and hedges share one vague policy | operators cannot tell which mechanism caused overload | separate accounting and controls |
| losing requests are not canceled | backend still does duplicate work after the winner returns | strong cancellation semantics and short hedge timeouts |
| hedging masks chronic latency regressions | p99 improves temporarily while budget burn stays elevated | track budget burn alongside latency and fix root cause |

## Observability

- metric: retry budget burn and hedge budget burn over time
- metric: hedge win rate and success-on-retry rate
- metric: extra attempts per original request by endpoint
- metric: canceled-loser latency and cancellation effectiveness
- log: budget-denied speculative attempts and policy state changes
- trace: original, retried, and hedged spans under one root request
- SLO: tail-latency improvement should stay within a bounded speculative-load budget

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit retry and hedge budgets | bounded extra load | more policy complexity | unlimited "best effort" retries |
| percentile-delayed hedging | helps real tail outliers | later trigger may miss some recoverable cases | immediate duplicate requests |
| aggressive loser cancellation | lowers wasted backend work | requires cancellation support end to end | fire-and-forget speculative attempts |

## Interview It

**Google framing:** "Your read path has a bad p99 even though median latency is excellent. How would you improve it safely?" The signal is whether you can talk about hedging with budgets instead of unbounded duplicate requests.

**Cloudflare framing:** "A globally distributed read service needs tighter tail latency without overloading origins." The signal is whether speculative work stays bounded and dependency-aware.

**Follow-ups:**
1. Which requests should never be hedged?
2. How large should the retry budget be?
3. What if replicas are correlated and hedging does not help?
4. How do you measure whether hedging is actually worth the cost?
5. When should the system disable hedging automatically?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-retry-budgets.md`
- `outputs/tradeoff-matrix-retry-budgets.md`

## Exercises

1. **Easy** — Pick a safe first hedge threshold for a read path with median `20 ms` and p99 `200 ms`.
2. **Medium** — Explain how you would account separately for retries and hedges in dashboards.
3. **Hard** — Redesign the policy when broad regional degradation makes replicas highly correlated.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/pub40801/) — foundational paper for hedged requests and tail mitigation
- [Google SRE books](https://sre.google/books/) — practical reliability framing for budgeted operational behavior
