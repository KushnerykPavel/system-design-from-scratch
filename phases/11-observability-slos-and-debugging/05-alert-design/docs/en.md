# Alert Design and Paging Quality

> Alerts should wake people because action is required, not because telemetry happens to exist.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Design alerts that detect real user harm, reduce noise, and point responders toward the next action.
**Prerequisites:** `10-reliability-retries-and-backpressure/03-circuit-breakers`, `10-reliability-retries-and-backpressure/04-load-shedding`, `11-observability-slos-and-debugging/01-sli-slo-error-budget`
**Estimated time:** ~60 min
**Primary artifact:** paging-quality checklist + interview card

## The Problem

Bad paging systems burn people out:

- low-signal warnings page overnight
- symptoms and causes both alert separately
- every threshold breach is treated the same
- alerts tell you nothing about user impact or first action

Senior interview answers separate:

- page-worthy user harm
- ticket-worthy capacity drift
- informational signals for dashboards

## Clarify

- Which alerts should wake humans immediately versus create tickets?
- Is this for a user-facing API, data pipeline, or internal control plane?
- Is the system allowed brief self-healing before paging?
- What runbook or mitigation exists once the page fires?

If the interviewer is vague, assume a 24/7 user-facing service with an error-budget-based paging policy and a supporting set of ticket-level alerts.

## Requirements

### Functional

- Page on fast, credible indicators of meaningful user impact.
- Separate urgent pages from non-urgent trend alerts.
- Include enough context to start mitigation quickly.

### Non-functional

- Minimize noisy pages and duplicate alerts.
- Keep alert logic understandable and reviewable.
- Ensure alerts remain useful under partial outages and telemetry gaps.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Service traffic | 120K req/s peak | small percentage shifts can still be severe |
| Operator pool | 5-8 engineers sharing on-call | noisy policies have real organizational cost |
| Acceptable page rate | low single digits per week for one service | forces strong selectivity |
| Burn windows | short and long windows together | catches fast outages and slow burns |
| Rough cost | alert fatigue, MTTR, and missed incidents | paging quality is an operations design problem |

## Architecture

A mature alert stack usually has:

1. **SLO or symptom alerts** for user harm.
2. **Cause-oriented alerts** for routing once symptoms exist.
3. **Ticket alerts** for hygiene, capacity, or trend degradation.

```text
telemetry
  -> symptom evaluation
  -> severity routing
  -> page / ticket / info
  -> runbook link and ownership
```

Practical rules:

- page on symptoms first
- deduplicate related causes
- use multi-window burn rates for SLOs
- add suppression when the upstream or region-wide event is already known

## Data Model & APIs

Useful entities:

- `AlertPolicy`
- `Severity`
- `SuppressionRule`
- `RunbookLink`

Useful fields:

- `signal_name`
- `threshold`
- `window`
- `severity`
- `owner`
- `action_hint`

Useful interfaces:

- `EvaluateAlert(signalState)`
- `RouteSeverity(policyResult)`
- `SuppressIfCovered(parentIncident)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| pages on causes instead of user symptoms | many pages during one incident with unclear user impact | make SLO and user-path alerts primary |
| thresholds too sensitive | pages spike during normal traffic shifts | test against historical data and expected noise |
| alert fires without action path | responders escalate manually with no first step | require runbook link and initial mitigation note |
| telemetry outage triggers false page storm | multiple services alert simultaneously on missing data | distinguish no-data cases from genuine system failure |

## Observability

- metric: alert volume, page volume, and duplicate alert ratio
- metric: acknowledged-to-actionable ratio and noisy alert suppression counts
- metric: MTTA and MTTR by alert policy
- log: alert evaluations, suppressions, and escalations
- trace: not usually primary, but useful for symptom-to-cause drill-down after a page
- SLO: paging system quality can itself be tracked by noise rate and missed-incident rate

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| SLO burn-rate paging | pages on real user harm | requires decent SLO instrumentation | paging on every infrastructure threshold |
| separate page/ticket/info channels | protects operator attention | more policy design upfront | one severity for everything |
| suppression and dedupe | less alert storm behavior | risk of hiding distinct incidents if done badly | naive page-per-signal routing |

## Interview It

**Google framing:** "How would you design alerts for a large serving system?" The signal is whether you can distinguish signal quality from simple threshold setting.

**Cloudflare framing:** "How would you page on problems in a global edge platform without waking engineers for every noisy regional blip?" The signal is whether you think about symptom-first paging and suppression.

**Follow-ups:**
1. What should never page by itself?
2. How do you combine short and long burn windows?
3. What if the metrics backend itself is degraded?
4. How do you stop one dependency incident from paging every downstream team?
5. How would you review alert quality after a quarter?

## Ship It

- `outputs/paging-quality-checklist.md`
- `outputs/interview-card-alert-design.md`

## Exercises

1. **Easy** — Convert a CPU-threshold page into a more user-impact-driven alert.
2. **Medium** — Design short-window and long-window burn alerts for an API SLO.
3. **Hard** — Redesign the paging strategy when one upstream dependency can take down many customer-facing services at once.

## Further Reading

- [Practical Alerting from Time-Series Data](https://sre.google/workbook/alerting-on-slos/) — concrete SLO-based paging guidance
- [Being On-Call](https://sre.google/sre-book/being-on-call/) — strong operational framing for alert quality
