# Metrics That Actually Explain the System

> Good metrics do not merely describe the system. They narrow the search space.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Choose metrics that expose saturation, bottlenecks, and user harm instead of flooding dashboards with low-signal counters.
**Prerequisites:** `02-estimation-and-cost/07-bottleneck-math`, `05-storage-indexing-and-access-patterns/02-access-pattern-first`, `11-observability-slos-and-debugging/01-sli-slo-error-budget`
**Estimated time:** ~75 min
**Primary artifact:** metrics review checklist + interview card

## The Problem

Weak system design answers say "we'll add metrics for CPU, memory, and QPS." That is too generic. Senior answers pick metrics that answer:

- is the user unhappy?
- is the system overloaded, blocked, or broken?
- where is the bottleneck moving?
- which component should we inspect next?

The goal is not more metrics. The goal is fewer, better metrics with clear diagnostic value.

## Clarify

- Is the system latency-sensitive, throughput-heavy, or durability-centric?
- Is the main risk queue buildup, dependency failure, storage saturation, or skew?
- Which path matters most: ingest, read, write, or background processing?
- Do operators need per-tenant visibility or only fleet-level health?

If the interviewer is vague, assume a distributed API with synchronous request handling and an asynchronous worker pipeline behind it.

## Requirements

### Functional

- Pick a metric set that explains user impact, load, saturation, and bottlenecks.
- Tie service-level metrics to component-level diagnosis.
- Define how metrics differ across synchronous and asynchronous paths.

### Non-functional

- Keep metric cardinality and cost under control.
- Make dashboards explain incidents quickly for on-call engineers.
- Avoid metric sets that stay green while queues or latency collapse elsewhere.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Front-door traffic | 150K req/s peak | determines metric volume and aggregation needs |
| Worker throughput | 500K jobs/minute | async path needs lag and age metrics, not just counts |
| Distinct tenants | 50K active/day | per-tenant dimensions can explode series count |
| Tail latency target | p99 under 300 ms | requires histogram-style latency insight |
| Rough cost | metrics ingest, retention, and query cost | high-cardinality labels can become a platform problem |

## Architecture

A strong metric model spans layers:

1. **User outcome metrics** show success and latency.
2. **Workload metrics** show demand shape.
3. **Resource and saturation metrics** show capacity pressure.
4. **Pipeline and dependency metrics** show where work is getting stuck.

```text
request / job
  -> service
  -> dependency / queue / storage
  -> metrics with bounded labels
  -> dashboards
  -> incident triage
```

A practical default is:

- one outcome panel
- one throughput panel
- one saturation panel
- one dependency panel
- one backlog panel if async work exists

## Data Model & APIs

Common metric families:

- `requests_total`
- `request_duration_ms`
- `queue_age_seconds`
- `dependency_errors_total`
- `worker_utilization_ratio`
- `shed_requests_total`

Important label decisions:

- `route_class` instead of full path
- `tenant_tier` instead of `tenant_id` for default dashboards
- `dependency_name`
- `region`
- `status_class`

Useful interfaces:

- `RecordRequest(result, latency, routeClass, region)`
- `ObserveQueueAge(queue, ageSeconds)`
- `RecordDependencyCall(name, statusClass, latencyMs)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| dashboard tracks activity but not success | request count healthy while SLI burns | pair throughput with error and latency distributions |
| CPU looks fine while queue backlog explodes | worker utilization stays moderate but queue age rises | add wait-time and backlog metrics, not just host metrics |
| per-tenant labels blow up cost | metric store series count or query latency spikes | aggregate by tier or sampled tenant detail |
| dependency panel hides one bad downstream | overall service latency rises without local resource saturation | break out dependency latency and error rate by target |

## Observability

- metric: good/bad request ratio and latency percentiles on the user path
- metric: queue depth, queue age, retry rate, and stale work share for async systems
- metric: dependency latency/error by named downstream
- metric: saturation such as concurrent in-flight work, worker utilization, or connection pool exhaustion
- log: sampled abnormal decisions like shed, retry, fallback, or stale-read serve
- trace: slow-path examples to explain where time is spent
- SLO: metrics should support the primary SLI instead of competing with it

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| histogram-like latency metrics | reveals tail behavior | more storage and query cost | averages only |
| route-class labels instead of raw paths | keeps cardinality bounded | less exact per-endpoint detail | raw URL labels everywhere |
| queue age over queue depth alone | better signal for user harm | requires time-aware instrumentation | only backlog counts |

## Interview It

**Google framing:** "What metrics would you add to understand whether a serving stack is healthy?" The signal is whether you bridge user impact to bottlenecks instead of listing machine stats.

**Cloudflare framing:** "What metrics matter for a globally distributed edge service?" The signal is whether you discuss latency distributions, dependency isolation, and high-cardinality restraint.

**Follow-ups:**
1. When is p50 useful and when is it misleading?
2. Which metric tells you queueing pain before users notice it?
3. How would you expose tenant pain without melting the metrics backend?
4. What metrics change when the system is mostly asynchronous?
5. Which panel would you look at first during a latency incident?

## Ship It

- `outputs/metrics-review-checklist.md`
- `outputs/interview-card-metrics-that-matter.md`

## Exercises

1. **Easy** — Pick three metrics for a read-heavy API and justify why each narrows the debugging space.
2. **Medium** — Design a metrics set for a batch worker system where queue age matters more than request latency.
3. **Hard** — Redesign the metric labels when per-tenant visibility is needed but the metrics platform is already near cardinality limits.

## Further Reading

- [Monitoring Distributed Systems](https://sre.google/sre-book/monitoring-distributed-systems/) — strong mental model for useful metrics
- [USE Method](http://www.brendangregg.com/usemethod.html) — practical saturation-oriented debugging lens
