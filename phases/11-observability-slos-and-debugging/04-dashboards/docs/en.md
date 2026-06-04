# Dashboards and Cardinality Discipline

> A dashboard is a debugging interface, not a scrapbook of every metric you can export.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design dashboards that accelerate incident triage while keeping labels, queries, and storage economically sane.
**Prerequisites:** `06-caching-and-invalidation/06-cache-layers`, `09-partitioning-sharding-and-rebalancing/02-hot-partitions`, `11-observability-slos-and-debugging/02-metrics-that-matter`
**Estimated time:** ~75 min
**Primary artifact:** dashboard linter + cardinality checklist

## The Problem

Teams often ship dashboards with dozens of panels and uncontrolled labels:

- every path, tenant, and user ID becomes a label
- queries time out during incidents
- dashboards show data but not decision order
- the metrics platform becomes the next system that needs incident response

Senior answers emphasize dashboard purpose and label discipline together.

## Clarify

- Who is the primary consumer: on-call engineer, service owner, or executive reviewer?
- Is the dashboard for fast triage, capacity review, or deep forensics?
- Which dimensions genuinely need filtering: region, route class, tenant tier, dependency?
- How expensive is the metrics backend allowed to become?

If the interviewer is vague, assume one primary on-call dashboard for a global API service plus a smaller set of drill-down views.

## Requirements

### Functional

- Build a first-response dashboard with clear panel order.
- Select labels that support debugging without unbounded cardinality.
- Preserve drill-down paths for region, dependency, and route class.

### Non-functional

- Keep query latency reliable during incidents.
- Avoid series explosion from raw identifiers.
- Make the dashboard interpretable by someone new to the service.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Time series baseline | 300K active series | dashboard design must respect backend scale |
| Potential raw label explosion | millions of `tenant_id` or `path` values | shows why label discipline matters |
| Dashboard refresh | every 15-30 seconds during incident | expensive queries can self-sabotage triage |
| Global regions | 20+ | region filters are useful, but still bounded |
| Rough cost | metrics ingestion + long query times + operator confusion | bad dashboards cost both money and incident time |

## Architecture

A good on-call dashboard usually flows top to bottom:

1. SLO / user impact
2. request volume and latency
3. errors by route class and region
4. dependency health
5. saturation and queueing
6. deployment or config changes

```text
dashboard
  -> headline health
  -> narrowing panels
  -> safe filters
  -> links to logs / traces / runbooks
```

Cardinality discipline means:

- use bounded taxonomies
- pre-aggregate where possible
- reserve raw IDs for ad hoc exploration, not default dashboards

## Data Model & APIs

Safe label examples:

- `region`
- `service`
- `route_class`
- `tenant_tier`
- `dependency_name`
- `status_class`

Dangerous default labels:

- `user_id`
- `session_id`
- raw `path`
- raw `tenant_id`
- unbounded `error_message`

Useful interfaces:

- `ValidateLabelSet(metricName, labels)`
- `EstimateSeries(metricName, labelCardinality)`
- `RecommendDashboardPanels(serviceType)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| raw path labels create series explosion | metrics backend series count spikes | normalize into route classes |
| dashboard has too many equal-priority panels | operators bounce around during incident | arrange panels in triage order |
| default filter is too broad | one bad region is hidden in global averages | provide region and dependency drill-downs |
| every team builds a custom naming scheme | operators misread dashboards across services | standardize metric and panel conventions |

## Observability

- metric: dashboard query latency and failed dashboard queries
- metric: active series count by metric family and label source
- metric: top cardinality offenders and fastest-growing label sets
- log: linter or review results for newly proposed metrics
- trace: slow dashboard query traces if the observability backend supports them
- SLO: critical triage dashboards should load reliably enough during active incidents

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| bounded label taxonomies | predictable scale and faster queries | less exact per-entity detail on default dashboards | raw IDs on every metric |
| opinionated panel order | faster first response | less freedom for dashboard authors | free-form panel collage |
| pre-aggregated views | cheap and stable incident queries | some loss of raw flexibility | every dashboard query computed from raw high-cardinality data |

## Interview It

**Google framing:** "How would you build a dashboard for a serving system that is actually useful during incidents?" The signal is whether the dashboard structure reflects triage flow.

**Cloudflare framing:** "How would you expose health for a global edge product without blowing up metrics cardinality?" The signal is whether you can balance drill-down usefulness with fleet scale.

**Follow-ups:**
1. Which labels are safe by default and which are not?
2. How do you detect that your observability backend is becoming the bottleneck?
3. What would the first three panels be?
4. When do you allow raw tenant-level drill-down?
5. How would you review new metrics before rollout?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/cardinality-discipline-checklist.md`
- `outputs/interview-card-dashboards.md`

## Exercises

1. **Easy** — Sketch the first four panels of an on-call dashboard for an API service.
2. **Medium** — Replace a raw `path` label strategy with a bounded alternative.
3. **Hard** — Redesign the dashboard and labels when enterprise support requires tenant-specific drill-down but the metrics backend is already near cost limits.

## Further Reading

- [The RED Method](https://grafana.com/blog/2018/08/02/the-red-method-how-to-instrument-your-services/) — useful dashboard organization for request-driven systems
- [Cardinality is key](https://www.robustperception.io/cardinality-is-key/) — practical cardinality intuition
