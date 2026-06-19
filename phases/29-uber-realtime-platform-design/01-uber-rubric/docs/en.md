# Uber Interview Rubric & Strong Signals
> Generic real-time system answers pass the bar. Uber-specific geospatial and dispatch thinking gets the offer.

**Type:** Concept
**Company focus:** Uber
**Learning goal:** Understand what Uber evaluates in system design — real-time geospatial indexing, sub-second matching, location data pipelines, and consistency trade-offs in a marketplace system.
**Prerequisites:** `20-low-latency-location-and-market-systems/02-location-updates`, `09-partitioning-sharding-and-rebalancing/03-placement`
**Estimated time:** ~75 min
**Primary artifact:** Uber rubric card + strong-hire checklist

---

## Uber Scale

| Metric | Value |
|---|---|
| Daily trips | 25M+ |
| Active drivers (peak) | ~5M |
| Active riders | ~100M |
| Countries | 70+ |
| Driver GPS events/sec | ~20M (5M drivers × 4 Hz) |
| Matching decisions/sec | ~300 (25M trips / 86400s) |
| ETA computations/sec | ~10K |

These numbers explain the engineering choices. At 20M GPS events/sec, a naive database-per-event approach fails. The system must batch, compress, and index spatially.

---

## What Uber Tests in System Design

### 1. Geospatial Thinking

Weak answers describe "a database that stores lat/lon and queries by distance." Strong answers name the indexing structure.

Uber's answer is **H3** — a hierarchical hexagonal geospatial indexing system open-sourced by Uber. Key properties:

- **Resolution levels**: each level subdivides cells into ~7 children
  - Res 9: ~174m edge length (city block)
  - Res 7: ~1.2km edge length (neighborhood)
  - Res 5: ~8.5km edge length (district)
- **K-ring query**: all cells within k steps of a center cell, used for proximity search
- **Hexagonal geometry**: all 6 neighbors are equidistant from center (unlike squares where diagonal neighbors are ~41% farther)

Strong hire signal: "I'd index driver locations into H3 res-9 cells, keep a Redis set per cell, and answer proximity queries with a k-ring(1) lookup returning 7 cells."

### 2. Real-Time Data Pipelines

Drivers emit GPS updates at 4 Hz. The pipeline must:
- Ingest 20M events/sec
- Update geospatial index within 1–2 seconds
- Serve matching queries at <50ms

Uber's pipeline: driver app → **Kafka** topic per city → location consumer service → **Redis** cell store. Kafka partitioned by driver ID for ordering, consumed by matching workers.

Strong hire signal: naming Kafka for ingest, Redis for the cell store, and partitioning strategy.

### 3. Dispatch & Matching

The dispatch system is called **DISCO** (Dispatch Optimization). Two matching strategies:

| Strategy | Description | Trade-off |
|---|---|---|
| Greedy nearest | Match each rider to closest available driver immediately | Simple, suboptimal at scale |
| Batch matching | Collect requests over 500ms window, solve assignment problem | Optimal, adds latency |

Uber uses batch matching in high-density cities. The assignment problem is solved with a variant of the **Hungarian algorithm** (minimum cost bipartite matching).

Strong hire signal: knowing batch matching exists, naming the 500ms window, mentioning ETA as the primary match signal (not raw distance).

### 4. State Management

The **trip state machine** is a key design artifact:

```
DRIVER: OFFLINE → AVAILABLE → GOING_TO_PICKUP → ON_TRIP → AVAILABLE
TRIP:   CREATED → MATCHED → PICKUP → IN_PROGRESS → COMPLETED / CANCELLED
```

Strong hire signal: drawing this state machine, naming what happens on each transition (push notification, ETA update, price lock), and discussing recovery when a driver goes offline mid-trip.

### 5. Pricing: Surge Algorithms

Surge pricing balances supply and demand in real time:
- Input: rider requests/min vs. active drivers in each H3 cell
- Output: surge multiplier (1.0x–5.9x)
- Algorithm: demand_ratio = requests / available_drivers; surge = f(demand_ratio)
- Smoothing: exponential moving average to prevent oscillation

Strong hire signal: modeling surge per H3 cell, discussing dampening, mentioning driver repositioning incentives.

---

## Strong-Hire vs Weak-Hire Patterns

### Strong-Hire Signals

| Signal | Example answer |
|---|---|
| Geospatial indexing | "H3 res-9 cells in Redis, k-ring proximity" |
| Location pipeline | "Kafka partitioned by driver ID, 20M events/sec" |
| Dispatch decomposition | "DISCO batch matching, 500ms window, Hungarian-variant" |
| Surge pricing design | "Per-cell demand ratio, EMA smoothing, 1-min recalculation" |
| Trip state machine | "Named all states plus recovery paths" |
| Failure recovery | "Driver offline: re-dispatch, Kafka replay on worker crash" |
| Observability | "P99 ETA error, match rate, driver utilization per city" |
| Uber vocabulary | "H3, DISCO, Cadence, Ringpop, TChannel, M3, Docstore" |

### Weak-Hire Signals

| Anti-pattern | Why it fails |
|---|---|
| "Store lat/lon in Postgres with PostGIS" | Does not scale to 20M events/sec without specialized pipeline |
| "Use Redis GEO commands" | Redis GEORADIUS is O(N+log M) — breaks at 5M drivers per region |
| "Real-time means low latency" | Vague; doesn't name Kafka, partitioning, or consumer lag budget |
| "Nearest driver wins" | Ignores batch optimization and ETA accuracy |
| "Use a message queue for location" | SQS/RabbitMQ not designed for 20M/sec stream ingest |

---

## Interview Milestone Map

| Time | Milestone | Strong-hire action |
|---|---|---|
| 0–5 min | Clarify scope | Ask: city vs global, driver count, trip volume, SLA |
| 5–15 min | Capacity estimate | 20M GPS events/sec, Redis ops/sec, Kafka throughput |
| 15–30 min | Core components | Location pipeline, geospatial index, matching engine |
| 30–45 min | Deep dive | H3 cell design OR batch matching OR surge pricing |
| 45–55 min | State machine | Trip + driver states, transitions, recovery |
| 55–65 min | Failure modes | Redis shard down, Kafka lag, driver offline mid-trip |
| 65–75 min | Observability | Key metrics, alerting thresholds, on-call runbook |

---

## Uber-Specific Vocabulary

| Term | What it is |
|---|---|
| H3 | Hierarchical hexagonal geospatial index (open-sourced by Uber, 2018) |
| DISCO | Dispatch Optimization — Uber's matching engine |
| Cadence / Temporal | Workflow orchestration system (Cadence open-sourced by Uber, became Temporal) |
| Ringpop | Consistent hashing ring for in-process cluster membership |
| TChannel | Uber's binary RPC protocol (multiplexed, request pipelining) |
| M3 | Uber's distributed metrics platform (Prometheus-compatible) |
| Docstore | Uber's document store built on MySQL with sharding layer |
| Peloton | Uber's resource manager for Mesos/Kubernetes workloads |

Dropping these names signals familiarity with Uber's engineering culture. Explain at least two in depth.

---

## Failure Modes

| Failure | Impact | Mitigation |
|---|---|---|
| Redis cell shard failure | Drivers in N cells invisible to dispatch | Replica failover, brief dispatch gap accepted |
| Kafka consumer lag | Location data stale, ETA errors | Auto-scale consumers, lag alert at 5s |
| Matching worker crash | In-flight batch lost | Kafka offset not committed, re-consumed on restart |
| Driver app offline | Trip orphaned | 30s heartbeat timeout → re-dispatch, driver penalized |
| ETA model degraded | Poor match quality | Fall back to map-only routing, alert on ETA error >20% |
| Surge model oscillation | Driver repositioning chaos | EMA dampening, floor multiplier at 1.0x |

---

## Observability

| Metric | SLO |
|---|---|
| Rider wait time P50 | <3 min |
| Rider wait time P95 | <8 min |
| ETA error (predicted vs actual) | <15% |
| Location update latency P99 | <2s |
| Match rate (requests matched / total) | >98% |
| Driver utilization per city | >70% |
| Kafka consumer lag | <5s |

---

## Trade-offs Table

| Decision | Option A | Option B | Uber's choice |
|---|---|---|---|
| Geospatial index | Redis GEO | H3 cell sets | H3 — predictable O(1) cell lookup |
| Matching strategy | Greedy nearest | Batch Hungarian | Batch in dense cities, greedy elsewhere |
| Location ingest | REST API per update | Kafka stream | Kafka — throughput and replay |
| Trip state | In-memory only | Durable workflow | Cadence/Temporal — recovery |
| Surge precision | Per-city | Per-H3-cell | Per-H3-cell — granular incentives |
| Driver comm | Polling | Push (long poll / WebSocket) | Push — reduced latency |

---

## Summary

Uber interviews test geospatial intuition (H3), pipeline thinking (Kafka → Redis), dispatch optimization (batch matching, ETA), and operational fluency (state machines, failure recovery). The vocabulary signals cultural fit. The depth signals engineering maturity.

A strong hire names the data structures, quantifies the scale, and draws the state machine without prompting.
