# Ride Dispatch & Matching Engine
> The best match is not always the nearest driver — it's the one who minimizes total system waiting time.

**Type:** Build
**Company focus:** Uber
**Learning goal:** Design Uber's dispatch system that matches 25M+ daily trips with drivers, optimizing for rider wait time, driver utilization, and ETA accuracy.
**Prerequisites:** `29-uber-realtime-platform-design/02-geospatial-indexing`, `20-low-latency-location-and-market-systems/04-stock-exchange`
**Estimated time:** ~90 min
**Primary artifact:** dispatch matching algorithm design + supply/demand flow

---

## What DISCO Does

DISCO (Dispatch Optimization) is Uber's proprietary dispatch system. Its core job is:

> Given a set of rider requests and a set of available drivers, compute the optimal assignment that minimizes total system wait time.

Key inputs:
- **Rider request**: location, trip type (UberX, Black, XL), timestamp
- **Driver supply**: location (H3 cell), rating, vehicle type, surge zone
- **ETA model**: predicted time-to-pickup for each driver-rider pair

Key outputs:
- **Match offer**: sent to top-N drivers simultaneously
- **Trip creation**: once a driver accepts

At 25M trips/day, DISCO handles ~300 new trip requests/second and ~10,000 ETA computations/second.

---

## Matching Strategies

### 1. Greedy Nearest Driver

For each rider request:
1. Find all drivers within k-ring of rider's H3 cell
2. Compute ETA to each driver
3. Assign the driver with the lowest ETA

**Pros**: simple, low latency (<50ms per match)
**Cons**: locally optimal but globally suboptimal — two riders 100m apart may both want the same driver while two drivers nearby go unmatched

### 2. Batch Matching (Uber's approach in dense cities)

Collect all rider requests and available drivers over a **500ms time window**, then solve the assignment as a bipartite matching problem:

```
Riders:  R1, R2, R3
Drivers: D1, D2, D3

Cost matrix (ETA in seconds):
        D1   D2   D3
R1 [    90   45  120 ]
R2 [   180   60   30 ]
R3 [    30  150   90 ]

Greedy assigns: R1→D2 (45s), R2→D3 (30s), R3→D1 (30s) → total = 105s
Hungarian finds: R1→D2 (45s), R2→D3 (30s), R3→D1 (30s) → same in this case
But in practice with asymmetric distributions, batch matching saves 10-20% total wait time.
```

The **Hungarian algorithm** (or auction algorithm) finds the minimum-cost bipartite matching in O(n³) — feasible for batches of up to ~50 riders/drivers per 500ms window in a city shard.

### 3. Match Score Components

ETA alone is insufficient. The match score incorporates:

| Component | Weight | Rationale |
|---|---|---|
| ETA (predicted pickup time) | Primary | Rider wait time is the primary SLA |
| Driver rating | Secondary | Protect rider experience |
| Trip type compatibility | Hard filter | UberBlack driver cannot take UberX request |
| Surge zone alignment | Bonus | Prefer drivers already in high-demand zone |
| Driver heading | Minor | Bonus if driver is already moving toward rider |

`match_score = ETA_seconds × (1 / driver_rating_normalized) × type_penalty`

Lower score is better. type_penalty = ∞ for incompatible types (hard filter).

---

## ETA Computation

ETA is the single most important input to matching. Uber uses two layers:

### Map-Based ETA
- Base: road network + turn restrictions + traffic speed segments
- Uber maintains its own map layer updated from telemetry
- Output: routing distance + estimated drive time at current speed

### ML-Predicted ETA
- Inputs: base map ETA, time of day, day of week, weather, historical pickup times at this location
- Output: corrected ETA (typically 5-15% more accurate than map-only)
- Fallback: if ML model is degraded (>20% error rate on recent trips), fall back to map-only ETA

ETA accuracy is measured as: `|predicted_eta - actual_arrival_time| / actual_arrival_time`

SLO: P50 ETA error < 10%, P95 < 25%.

---

## Driver Offer Protocol

After DISCO computes the optimal match for a batch:

1. Send offer simultaneously to **top-3 drivers** (ranked by match score) for each rider
2. First driver to **accept** within 15 seconds wins the trip
3. If all 3 decline (or timeout): recompute with next-best candidates, penalize repeat decliners
4. Driver acceptance triggers: trip state CREATED→MATCHED, push notification to rider

Why top-3? Hedges against the best driver being momentarily distracted. Avoids the single-offer waterfall which adds 15s × N latency per declined offer.

### Driver State Machine

```
OFFLINE ────────────────────────────────────────────
           ↑ app exit                               |
           |                                        |
AVAILABLE ←──── trip complete ←── ON_TRIP          |
    |                                  ↑            |
    | offer accepted                   |            |
    ↓                                  |            |
GOING_TO_PICKUP ──── rider in car ─────            |
    |                                              |
    | app crash / GPS loss (30s timeout) ──────────
    ↓
    (DISPATCH removes from pool, re-dispatches trip)
```

### Trip State Machine

```
CREATED → MATCHED → DRIVER_EN_ROUTE → PICKUP → IN_PROGRESS → COMPLETED
                                                             ↘ CANCELLED
```

Durable state managed by Cadence/Temporal workflow. Each transition emits an event to:
- Rider app (push notification + ETA update)
- Driver app (navigation instructions)
- Pricing service (trip start/end for billing)
- Analytics pipeline (via Kafka)

---

## Supply Forecasting

DISCO doesn't only react to requests — it predicts supply/demand imbalances to pre-position drivers.

**Per-H3-cell supply forecast (res 7 = neighborhood level)**:
- Input: current driver count in cell, historical trip rates, time of day, events
- Horizon: 5/10/15 minutes
- Output: predicted driver surplus/deficit

If a concert ends in 10 minutes in cell X, DISCO increases surge pricing in X now to attract drivers before the surge hits. This is the **pre-positioning incentive** system.

---

## Market Segmentation

Each trip type runs a separate matching pool:

| Pool | Vehicles | Priority |
|---|---|---|
| UberX | Standard 4-door | Highest volume |
| UberXL | 6-seat vehicles | Medium |
| UberBlack | Licensed livery | Premium SLA |
| UberPool / UberShare | Standard + route-matching | Complex (multi-rider) |
| Eats / Freight | Commercial vehicles | Separate dispatch system |

Pools are independent: a driver registered only for UberBlack never appears in UberX dispatch candidates. This prevents complexity and maintains service-level promises.

---

## System Architecture

```
                    ┌──────────────────────────┐
Driver GPS (4Hz) ──>│  Kafka: driver-location  │
                    └──────────┬───────────────┘
                               │
                    ┌──────────▼───────────────┐
                    │  Location Consumer       │
                    │  (H3 cell update)        │
                    └──────────┬───────────────┘
                               │
                    ┌──────────▼───────────────┐
                    │  Redis Cell Store        │
                    │  h3cell:{id} → {drivers} │
                    └──────────┬───────────────┘
                               │
Rider Request ─────────────────▼
                    ┌──────────────────────────┐
                    │  Dispatch Service        │
                    │  (DISCO batch matcher)   │
                    │  · k-ring cell lookup    │
                    │  · ETA computation       │
                    │  · Hungarian assignment  │
                    │  · 500ms batch window    │
                    └──────────┬───────────────┘
                               │
                    ┌──────────▼───────────────┐
                    │  Offer Distribution      │
                    │  · push to top-3 drivers │
                    │  · 15s acceptance window │
                    └──────────┬───────────────┘
                               │
                    ┌──────────▼───────────────┐
                    │  Trip Workflow (Cadence) │
                    │  · durable state machine │
                    │  · state change events   │
                    └──────────────────────────┘
```

---

## Failure Modes

| Failure | Impact | Mitigation |
|---|---|---|
| Matching worker crash mid-batch | In-flight batch lost | Kafka offset not committed; re-consumed on restart |
| Driver goes offline after match | Trip orphaned | 30s heartbeat timeout → re-dispatch |
| ETA model degraded | Poor match quality | Fall back to map-only routing; alert on P50 error >20% |
| Top-3 drivers all decline | Rider waits | Broaden k-ring to k=2, re-offer to next candidates |
| Batch window too long (>500ms) | Match latency grows | Auto-scale dispatch workers; shed load by reducing batch window |
| Redis cell shard failure | Drivers invisible | Replica failover; dispatch gap for affected cells |

---

## Observability

| Metric | SLO | Alert threshold |
|---|---|---|
| Rider wait time P50 | <3 min | >4 min |
| Rider wait time P95 | <8 min | >12 min |
| ETA error P50 | <10% | >20% |
| Match latency P99 | <500ms | >1s |
| Offer acceptance rate | >85% | <70% |
| Re-dispatch rate | <5% | >10% |
| Kafka consumer lag | <5s | >15s |

---

## Trade-offs Summary

| Decision | Greedy | Batch (Hungarian) | Uber's choice |
|---|---|---|---|
| Latency | <50ms | 500ms | Batch in dense cities, greedy in sparse |
| Optimality | Local | Global | Batch — 10-20% better total wait time |
| Complexity | Low | High | Justified at density |
| Failure blast radius | One match | Entire batch | Batch re-consumed from Kafka |

The 500ms batch window is a deliberate latency budget trade-off: riders don't perceive 500ms of extra wait, but the system gains meaningful optimality improvement.
