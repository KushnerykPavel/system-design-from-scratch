# SLIs, SLOs, and Error Budgets

> A target is only useful if it changes what operators do next.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Define user-meaningful SLIs, set realistic SLOs, and reason about error budget burn instead of reciting uptime percentages mechanically.
**Prerequisites:** `02-estimation-and-cost/01-qps-and-request-mix`, `03-design-framework-and-timing/05-wrap-up`, `10-reliability-retries-and-backpressure/01-timeouts-and-retries`
**Estimated time:** ~75 min
**Primary artifact:** SLO budget evaluator + interview card

## The Problem

Interview answers often say "we'll monitor uptime" and stop there. That is weak because:

- uptime alone hides whether users can complete important actions
- availability without latency can reward a slow but technically "up" system
- a target without a budget policy does not guide rollout or incident trade-offs

Senior answers define:

- what users actually experience
- how success is measured
- what failure budget remains
- what decisions change when the budget burns too fast

## Clarify

- Which user journey matters most: read path, write path, or end-to-end workflow?
- Is the system consumer-facing, internal platform, or edge infrastructure?
- What level of badness matters more: outright failures or extreme latency?
- Does the business prefer tighter targets for premium tenants or one shared objective?

If the interviewer is vague, assume a user-facing API with a critical read path, a smaller but important write path, and explicit latency plus success targets.

## Requirements

### Functional

- Define one or two SLIs tied to a meaningful user action.
- Express an SLO target over a clear window.
- Compute remaining error budget and burn rate.

### Non-functional

- Keep the objective simple enough for operators and interviewers to reason about quickly.
- Avoid vanity metrics that look healthy while users are unhappy.
- Make the SLO actionable for rollout, paging, and prioritization decisions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read traffic | 80K req/s peak | large volume makes tiny failure percentages operationally significant |
| Write traffic | 8K req/s peak | lower volume may justify separate objectives |
| Monthly request count | ~200B reads | converts percentages into concrete allowed failures |
| Candidate target | 99.9% success, 99% under 250 ms | shows how strict targets change budget size |
| Rough cost | tighter objectives require more redundancy, testing, and staffing | SLO choice is an economic decision too |

## Architecture

An SLO system has two layers:

1. **Measurement layer** computes SLIs from trusted telemetry.
2. **Decision layer** turns budget burn into actions such as paging, launch freezes, or reliability work.

```text
user request
  -> service path
  -> metrics / logs / traces
  -> SLI computation
  -> rolling window objective
  -> burn-rate policy
  -> operator action
```

Important modeling choices:

- measure at the user-facing boundary when possible
- separate availability and latency if both matter
- define exclusions explicitly, not implicitly
- keep one "headline" SLO and a few supporting indicators

## Data Model & APIs

Useful entities:

- `SLIDefinition`
- `SLOObjective`
- `BudgetWindow`
- `BurnAlertPolicy`

Useful fields:

- `name`
- `numerator_events`
- `denominator_events`
- `latency_threshold_ms`
- `target_ratio`
- `window_days`
- `alert_burn_multiple`

Useful interfaces:

- `ComputeBudget(total_events, bad_events, target_ratio)`
- `EvaluateBurn(consumed_budget_ratio, elapsed_window_ratio)`
- `RecommendAction(burn_rate)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| SLI tracks proxy health instead of user success | SLI stays green while support tickets or trace failures rise | move measurement closer to completed user action |
| target is unrealistically strict | budget is exhausted constantly with no room for launches | reset objective to match user need and operating maturity |
| denominator excludes hard cases | suspiciously healthy SLI during known incidents | explicitly define inclusion rules and review exclusions |
| one blended SLO hides premium tenant pain | global SLO looks fine while important customers degrade | add segmentation by plan, region, or request class |

## Observability

- metric: good events, total events, and bad event ratio per SLI
- metric: rolling error budget remaining and multi-window burn rate
- metric: latency distribution for the user path behind the SLI
- log: budget-policy actions such as freeze, rollback, or paging escalation
- trace: representative failing requests to explain why the SLI is burning
- SLO: one primary availability or latency objective tied to a named user journey

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| user-journey SLI | aligns with actual experience | harder to instrument end-to-end | host uptime as primary objective |
| separate read and write objectives | more honest about distinct paths | more policy and dashboard complexity | one blended objective for all traffic |
| explicit error budget policy | makes trade-offs operationally real | forces uncomfortable prioritization | target with no action thresholds |

## Interview It

**Google framing:** "How would you define SLOs for a storage or serving system?" The signal is whether you tie objectives to user impact and error budget policy rather than just percentage targets.

**Cloudflare framing:** "How would you define reliability objectives for an edge API product?" The signal is whether you account for latency, partial regional issues, and the cost of globally tight objectives.

**Follow-ups:**
1. Should internal services have the same SLO design as customer-facing APIs?
2. When would you split one SLO into per-region or per-tier objectives?
3. What do you do if the error budget is burning but uptime still looks good?
4. How strict can an SLO be before it becomes organizational theater?
5. What changes at 10x traffic?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-sli-slo-error-budget.md`
- `outputs/error-budget-checklist.md`

## Exercises

1. **Easy** — Define a primary SLI for a read-heavy API and explain why you rejected pure host uptime.
2. **Medium** — Compute the monthly error budget for `99.95%` success at `25B` requests per month.
3. **Hard** — Redesign the objective when premium enterprise tenants need stricter guarantees than the free tier.

## Further Reading

- [Service Level Objectives](https://sre.google/workbook/service-level-objectives/) — strong framing for practical SLO design
- [Embracing Risk](https://sre.google/sre-book/embracing-risk/) — the canonical error-budget mindset
