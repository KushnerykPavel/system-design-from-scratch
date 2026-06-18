# A/B Testing & Experimentation Platform

> At Netflix, if you cannot measure it, you cannot ship it. Every product decision goes through an experiment.

**Type:** Build  
**Company focus:** Netflix  
**Learning goal:** Design an experimentation platform at Netflix scale. Cover assignment service (consistent hashing by user_id), treatment bucketing, metrics collection pipeline, statistical significance, experiment isolation, holdout groups, and how Netflix runs 1,000+ simultaneous experiments.  
**Prerequisites:** `04-recommendation-engine`, `07-realtime-data-pipeline`  
**Estimated time:** ~75 min  
**Primary artifact:** experimentation platform design doc + assignment service spec  

## The Problem

Netflix runs over 1,000 simultaneous A/B experiments affecting every aspect of the product — recommendation algorithms, UI layouts, encoding quality, playback startup time, and more. Each experiment must:

- Assign each subscriber to exactly one treatment bucket per experiment.
- Keep assignments stable across sessions (same subscriber always sees the same variant).
- Prevent experiments from interfering with each other.
- Collect metrics for statistical analysis.
- Detect significance without requiring human intervention for routine experiments.

Design the platform that makes this possible at 270M subscriber scale.

## Clarify

- What is the assignment unit? (user_id, device_id, session_id — trade-off between consistency and granularity)
- How are experiment parameters stored? (experiment definition: name, variants, traffic allocation)
- What is the minimum detectable effect size? (drives required sample size and experiment duration)
- How do you handle experiments that must be mutually exclusive? (e.g., two UI experiments cannot overlap)
- How long do experiments run before a decision is required?
- Who can launch an experiment, and what review gates exist?

## Requirements

### Functional

- Assign each subscriber to a treatment bucket for each active experiment they are eligible for.
- Guarantee assignment stability: the same subscriber always receives the same treatment for a given experiment.
- Prevent cross-experiment contamination.
- Collect metric events associated with experiment assignments.
- Compute statistical significance metrics for experiment evaluation.
- Support holdout groups (subscribers excluded from all experiments).

### Non-functional

- Assignment latency: under 5ms per request (it is in the critical path of every page load).
- Scale: 270M subscribers, 1,000+ active experiments simultaneously.
- Consistency: assignments must be deterministic — no database lookup required at serving time.
- Auditability: assignment decisions must be reproducible from experiment config + user_id.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active experiments | 1,000+ | assignment computation must handle all experiments per request |
| Subscribers eligible per experiment | 1M–270M depending on targeting | drives bucket width and sample size calculation |
| Assignments per second | peak ~1M (concurrent homepage loads) | drives assignment service throughput requirement |
| Metric events per day | billions (clicks, plays, skips) | drives metrics pipeline sizing |
| Experiment duration | 2–4 weeks typical | drives data retention and significance computation cadence |

## Architecture

```text
[assignment path]
subscriber request
  -> assignment service
     -> load experiment definitions (in-memory, refreshed from config store)
     -> hash(user_id, experiment_id) -> bucket index
     -> look up treatment for bucket
     -> return {experiment_id: treatment_name} map
  -> caller uses treatment to decide behavior

[metrics path]
subscriber action (click, play, skip)
  -> event emitter (device or server)
  -> Kafka topic: experiment-events
  -> Flink job: join event with assignment (by user_id + experiment_id)
  -> Druid/Pinot: OLAP aggregation by experiment + treatment + metric
  -> significance calculator (runs daily or on demand)
  -> experiment dashboard
```

### Consistent Hashing Assignment

Assignments are computed deterministically without a database lookup:

```python
def assign(user_id: str, experiment_id: str, num_buckets: int = 10000) -> int:
    h = hash(f"{user_id}:{experiment_id}")
    return h % num_buckets

def get_treatment(user_id, experiment) -> str:
    bucket = assign(user_id, experiment.id)
    for treatment in experiment.treatments:
        if treatment.bucket_start <= bucket < treatment.bucket_end:
            return treatment.name
    return "control"
```

This guarantees:
- Same user always gets same bucket for the same experiment.
- Different experiments produce independent assignments (experiment_id in the hash key).
- No database needed at assignment time.

### Experiment Definition

```text
experiment_id: "rec-algo-v2"
status: active
start_date: 2026-01-01
treatments:
  - name: control
    bucket_start: 0
    bucket_end: 4000   # 40% of traffic
  - name: treatment_A
    bucket_start: 4000
    bucket_end: 6000   # 20% of traffic
  - name: treatment_B
    bucket_start: 6000
    bucket_end: 8000   # 20% of traffic
holdout_bucket_start: 9000
holdout_bucket_end: 10000  # 10% excluded from all experiments
targeting:
  - country: US
  - device_type: smart_tv
```

### Mutual Exclusion

Experiments that test the same UI surface or interact with the same feature must be mutually exclusive to avoid confounding results:

```text
exclusion_group: "homepage-ui"
experiments_in_group:
  - homepage-layout-v3
  - thumbnail-size-test
  - row-order-experiment
```

Within a group, each subscriber is assigned to at most one experiment. Assignment uses a group-level hash to select the experiment, then a within-experiment hash for the treatment bucket.

### Holdout Groups

A holdout group is a set of subscribers excluded from all (or a category of) experiments. They provide a baseline to measure the cumulative effect of all launched features over a period:

```text
holdout_bucket: [9000, 10000)  # global holdout: 10% of all users
```

Comparing holdout to the rest of the fleet after 6 months reveals the aggregate impact of shipped features.

### Statistical Significance

For each experiment, the platform computes:
- Mean metric value per treatment group.
- Standard error of the difference.
- p-value (t-test or Z-test for proportions).
- Confidence interval for the effect size.

Netflix uses a two-sided test with p < 0.05 as the default threshold, with additional checks for multiple comparisons (Bonferroni correction for simultaneous experiments measuring many metrics).

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Assignment service is down | Homepage loads fail assignment lookup | Fall back to control group for all subscribers; never fail page load due to assignment failure |
| Experiment config store is stale | Assignment produces wrong buckets | Assignment service holds in-memory snapshot with version check; refuse to serve if snapshot is too old |
| Metric event dropped before joining with assignment | Experiment metrics become incomplete | Kafka at-least-once delivery; idempotent deduplication in Flink join |
| Two experiments contaminate each other | Holdout group metrics shift | Enforce mutual exclusion at experiment creation time; holdout group provides contamination signal |
| Experiment runs too long (peeking problem) | Team declares significance prematurely | Platform enforces minimum run duration and sample size requirements before surfacing significance |

## Observability

- metric: assignment service latency at p50/p95/p99
- metric: assignments per second by experiment
- metric: metric event lag (time from action to joining with assignment)
- metric: experiment coverage (% of eligible subscribers assigned)
- metric: holdout group metric delta vs treated population
- log: experiment config version loaded by assignment service
- alert: assignment service latency exceeds 10ms

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Deterministic hash assignment (no DB) | Sub-millisecond assignment, infinite scale | Cannot easily reassign users mid-experiment | DB lookup allows dynamic reassignment but adds latency and DB dependency in the hot path |
| Mutual exclusion groups | Clean separation of interacting experiments | Reduces available traffic for each experiment in the group | Allowing overlapping experiments is simpler but produces confounded results |
| Global holdout group | Measures cumulative effect of all shipped features | 10% of traffic permanently excluded from experiments | No holdout means you can never measure long-term experiment accumulation effect |
| Two-sided tests | Guards against both positive and negative effects | Less statistical power per experiment at fixed N | One-sided tests are more powerful but miss unexpected negative effects |

## Interview It

**Netflix framing:** "Design Netflix's A/B testing platform." Strong answers cover deterministic assignment, experiment isolation, metric collection and joining, and statistical significance. Weak answers describe only a bucketing mechanism without discussing metric pipelines, contamination, or holdouts.

**Follow-ups:**
1. How do you prevent the same subscriber from being in two conflicting UI experiments simultaneously?
2. A new experiment is launched at 10% traffic. How do you verify the assignment is correct?
3. What is the peeking problem, and how does your platform prevent it?
4. How does the holdout group help you measure the effect of a whole year's worth of features?
5. How would you design the system to support experiments at the device_id level (before login)?

## Ship It

- `outputs/design-doc-ab-testing-platform.md`
- `outputs/assignment-service-spec.md`
- `outputs/interview-card-ab-testing-platform.md`

## Exercises

1. **Easy** — Compute the bucket assignment for user_id="abc123", experiment_id="rec-algo-v2" using a 10,000-bucket space. Which treatment does the subscriber receive?  
2. **Medium** — Design the Flink job that joins playback events with experiment assignments to produce per-treatment metric aggregations.  
3. **Hard** — Design the experiment creation review gate: what checks must pass before an experiment goes live (traffic allocation, mutual exclusion, minimum sample size, metric selection)?  

## Further Reading

- [Netflix experimentation platform blog](https://netflixtechblog.com/its-all-a-bout-testing-the-netflix-experimentation-platform-4e1ca458c15b)  
- [Trustworthy online controlled experiments (Kohavi)](https://www.amazon.com/Trustworthy-Online-Controlled-Experiments-Practical/dp/1108724264)  
- [Statsig experimentation platform](https://statsig.com/blog)  
