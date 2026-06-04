# Observability Drill

> Observability is complete only when it supports targets, diagnosis, action, and redesign together.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Integrate Phase 11 by designing the observability model for a realistic distributed system, including SLOs, metrics, traces, dashboards, alerts, runbooks, and incident narrative.
**Prerequisites:** `11-observability-slos-and-debugging/07-debugging-narrative`
**Estimated time:** ~60 min
**Primary artifact:** observability drill rubric + answer checker

## The Problem

Use this prompt:

"Design the observability and incident-response model for a globally distributed API gateway that enforces auth, rate limits requests, fans out to multiple internal services, and serves both free-tier and enterprise traffic."

This drill is strong because it forces the whole phase together:

- user-facing SLIs and SLOs
- fleet and dependency metrics
- edge-friendly tracing and logging decisions
- dashboard hierarchy and cardinality control
- page-worthy alerts and first-response workflow
- debugging narrative under regional or tenant-scoped incidents

## Clarify

- Which user journeys matter most: request acceptance, low latency, or policy correctness?
- Do enterprise customers need stricter objectives or just better visibility?
- Are edge POPs allowed to degrade independently without global paging?
- What is the biggest operational risk: dependency failures, noisy tenants, or observability cost explosion?

If the interviewer is vague, assume latency-sensitive global traffic, multi-tenant behavior, and a need for both regional triage and fleet-wide policy decisions.

## Requirements

### Functional

- Define user-meaningful SLIs and SLOs for the gateway path.
- Specify core metrics, logs, traces, dashboards, and alerts.
- Provide a first-response incident workflow and debugging narrative.

### Non-functional

- Keep telemetry affordable at global scale.
- Preserve per-region and per-tier visibility without unbounded cardinality.
- Ensure observability still helps when one dependency or region is degraded.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Global QPS | 1M req/s peak | telemetry volume and sampling policy matter immediately |
| Regions / POPs | 30+ | requires regional scoping and incident isolation |
| Internal fanout | 3-6 dependencies on hot path | traces and dependency metrics must narrow causality |
| Customer tiers | free and enterprise | may justify segmented SLOs and dashboards |
| Rough cost | metrics, logs, traces, paging, and human time | observability design is a platform trade-off |

## Architecture

A strong answer usually contains:

1. one primary request-success SLI and one latency SLI
2. metrics for volume, errors, dependencies, queueing, and saturation
3. structured logs and sampled traces with request correlation
4. a triage dashboard with bounded dimensions
5. SLO-burn and dependency-aware paging rules
6. runbook and debugging sequence for regional versus global issues

```text
client
  -> edge gateway
  -> auth / rate limit / routing
  -> internal services
  -> telemetry pipeline
  -> SLI / dashboard / alert / runbook / investigation
```

## Data Model & APIs

Core entities:

- `SLIConfig`
- `AlertPolicy`
- `DashboardPanel`
- `TraceSamplePolicy`
- `RunbookStep`

Useful APIs:

- `EvaluateGatewayBudget(windowState)`
- `ValidateMetricLabels(metricDefinition)`
- `PropagateTraceContext(request)`
- `RecommendEscalation(scope, severity)`

The best answers connect these pieces rather than describing them independently.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| global averages hide one bad region | enterprise complaints or regional SLI burn while fleet metric looks fine | add regional and tier drill-down with bounded labels |
| logs and traces are too expensive at peak | observability backend cost or ingest lag spikes | structured sampling, redaction, and error-focused retention |
| alerting pages every team during one dependency outage | duplicated pages and confused ownership | symptom-first routing and suppression |
| dashboards become unusable during incidents | query latency rises or panels time out | pre-aggregation and label discipline |

## Observability

- metric: request success, latency percentiles, shed rate, dependency latency, and auth/rate-limit decision rates
- metric: region-level and tier-level budget burn with bounded labels
- log: structured request decisions, dependency error classes, and mitigation actions
- trace: sampled end-to-end request path with auth, rate limit, and downstream spans
- SLO: user-path success and latency should remain explainable and actionable across tiers and regions

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| segmented enterprise visibility | clearer support and incident targeting | more telemetry policy complexity | one fleet-wide view only |
| error-biased tracing | good debugging value per dollar | less complete healthy-traffic coverage | tracing everything |
| symptom-first alerting | better paging quality | requires well-defined SLOs | pure infrastructure-threshold paging |

## Interview It

**Google framing:** "Design the observability model for a serving platform." Expect follow-ups on SLI choice, burn-rate paging, and how telemetry narrows incidents.

**Cloudflare framing:** "Design observability for a global edge gateway." Expect pushback on cardinality, regional visibility, and how to debug partial edge failures.

**Follow-ups:**
1. What is your primary SLI and why is it user-meaningful?
2. Which labels do you allow and forbid on default metrics?
3. What pages immediately and what becomes a ticket?
4. How do you debug one-region degradation differently from a global outage?
5. What changes at 10x traffic?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/observability-drill-rubric.md`
- `outputs/interview-card-observability-drill.md`

## Exercises

1. **Easy** — Give a 10-minute answer focused only on the highest-signal observability choices.
2. **Medium** — Redesign the drill assuming enterprise tenants get stricter latency guarantees than free-tier traffic.
3. **Hard** — Redesign the observability model when the metrics backend itself becomes a scaling bottleneck during incidents.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline for structuring a full design answer
- [The Site Reliability Workbook](https://sre.google/workbook/table-of-contents/) — practical material on SLOs, alerting, and incident response
