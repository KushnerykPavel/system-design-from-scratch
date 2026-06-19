# Surge Pricing & Dynamic Pricing Pipeline
> Surge is not gouging — it is the signal that repositions supply to where demand needs it.

**Type:** Build
**Company focus:** Uber
**Learning goal:** Design Uber's surge pricing system that computes per-zone multipliers in near-real-time based on supply/demand imbalance.
**Prerequisites:** `20-low-latency-location-and-market-systems/05-market-data-fanout`, `06-caching-and-invalidation/02-freshness-models`
**Estimated time:** ~75 min
**Primary artifact:** surge multiplier computation pipeline

---

## What Is Surge Pricing?

Surge pricing is Uber's mechanism for balancing supply and demand in real time. When rider demand outpaces available driver supply in a geographic zone, the price multiplier rises. Higher prices accomplish two things simultaneously:

1. **Demand-side**: some riders choose to wait, reducing request volume
2. **Supply-side**: higher earnings attract idle or distant drivers into the zone

The critical insight is that surge is a **market signal**, not a profit lever. The goal is equilibrium, not extraction.

---

## Zone Model: H3 Aggregation

Driver locations and trip requests are indexed at H3 resolution 9 (neighborhood level). For surge computation, Uber aggregates res-9 cells into larger surge zones using H3 resolution 6 (city district level). Each res-6 cell covers roughly 36 km² — large enough to capture meaningful supply/demand dynamics, small enough to be geographically useful.

**Why not compute surge at res-9?**
- Individual cells can be empty by chance even in active areas
- Aggregation smooths out momentary gaps in driver coverage
- Fewer zones = fewer Kafka partitions + lower fan-out to rider apps

---

## Multiplier Thresholds

The surge multiplier maps the demand/supply ratio to a price tier:

| Demand/Supply Ratio | Multiplier |
|---|---|
| 0.0 – 1.0 | 1.0× (no surge) |
| 1.0 – 1.25 | 1.2× |
| 1.25 – 1.5 | 1.5× |
| 1.5 – 2.0 | 2.0× |
| 2.0 – 4.0 | 4.0× |
| > 4.0 | 8.0× (cap) |

The cap at 8× is both regulatory (many cities ban higher multipliers) and brand protection — an uncapped multiplier during a disaster event would be catastrophic PR.

---

## Pipeline Architecture

```
Kafka: trip-requests ──────────────────────────────────┐
                                                        ▼
Kafka: driver-locations ──────────────────▶  Surge Worker
                                             · aggregate by zone (60s window)
                                             · count open requests per zone
                                             · count available drivers per zone
                                             · compute ratio → multiplier
                                             · publish SurgeZone events
                                                        │
                               ┌────────────────────────┤
                               ▼                        ▼
                     Redis (zone → multiplier)    Kafka: surge-zones
                               │                        │
                     Surge API (poll 30s)      Rider WebSocket push
```

### Kafka Topics

| Topic | Producer | Consumer | Key |
|---|---|---|---|
| `trip-requests` | Trip service | Surge worker | zone_id |
| `driver-locations` | Location service | Surge worker | driver_id |
| `surge-zones` | Surge worker | Rider app, Driver app | zone_id |

Partitioning by `zone_id` ensures that all events for a given zone land in the same partition, giving the surge worker consistent ordering without cross-partition coordination.

### Surge Worker Internals

Every 60 seconds per zone:

1. Count `OpenRequests`: trip requests received in the last 30 seconds that haven't been matched
2. Count `AvailableDrivers`: drivers in the zone with `AVAILABLE` status in Redis
3. Compute ratio: `demand_ratio = max(OpenRequests, 1) / max(AvailableDrivers, 1)`
4. Map ratio to multiplier using threshold table
5. Write `zone_id → multiplier` to Redis with TTL 120s
6. Publish `SurgeZoneEvent` to Kafka `surge-zones` topic

The 60-second window is a deliberate trade-off: more frequent updates would increase CPU and Kafka throughput, but the 60s lag is imperceptible to riders browsing the app.

---

## Upfront Pricing

Uber quotes an exact price at the moment of trip request — the **upfront price**. This means:

1. Rider sees "This trip will cost $14.50" before booking
2. The surge multiplier is **locked** at request time
3. If traffic is worse than predicted, Uber absorbs the difference
4. If traffic is better, the rider pays the upfront price (no discount)

**Why lock the multiplier?** Rider experience. A multiplier that changes between the "confirm" and "charged" moment destroys trust.

**Surge zone boundary race condition**: if the rider's zone and the driver's zone have different multipliers at request time, the rider's zone multiplier governs. This is stored on the trip document at creation time.

---

## Price Transparency & Confirmation

When the surge multiplier is ≥ 2.0×, Uber requires **explicit rider confirmation**:

```
┌─────────────────────────────────┐
│  🔴 High Demand — 2.3× surge   │
│                                 │
│  Rides are in high demand.      │
│  This trip costs $24.80         │
│  (normal price: $10.78)         │
│                                 │
│  [Confirm]        [Cancel]      │
└─────────────────────────────────┘
```

This pattern mirrors Stripe's strong customer authentication: the UX friction is intentional. It reduces accidental high-price trip confirmations and associated support tickets.

---

## Driver Repositioning Incentive

High-surge zones push a **heat map notification** to nearby offline or idle drivers:

- Triggered when a zone reaches ≥ 1.5× multiplier
- Shows potential earnings uplift (e.g., "+$3 per trip in your area")
- Displayed on surge heat map overlay in driver app

This is the mechanism by which surge pricing actually increases supply — without it, surge would be pure demand destruction.

---

## Regulatory Constraints

Several cities ban surge pricing during declared emergencies (natural disasters, public health events):

| City/Region | Policy |
|---|---|
| New York City | Price gouging law caps surge during declared emergencies |
| California | AB 1383 bans price increases > 10% during emergencies |
| EU | Platform-to-Business regulation limits dynamic pricing transparency |

Uber implements this as a geo-fenced override: when an emergency flag is set for a region (manually by an on-call engineer or via automated ingestion of government feeds), the surge worker caps the multiplier at 1.0× for all zones in that region.

---

## Failure Modes

| Failure | Impact | Mitigation |
|---|---|---|
| Surge worker crash | Stale multipliers | Redis TTL (120s) serves last-known-good; rider sees "loading" after TTL expires |
| Demand spike faster than 60s update cycle | Multiplier underestimates surge | Adaptive polling: worker detects rapid request rate increase and re-computes at 10s intervals |
| Zone boundary mismatch between apps | Driver and rider see different multipliers | Trip document stores multiplier at creation; authoritative for billing |
| Redis cluster partition | Zone multipliers unavailable | Fall back to default 1.0× (safe degradation); alert on-call |
| Kafka consumer lag | Surge lags reality | Monitor consumer group lag; alert at > 120s |

---

## Observability

| Metric | SLO | Alert |
|---|---|---|
| Surge computation latency P99 | <5s | >15s |
| Redis multiplier freshness | <120s | >180s (TTL expired) |
| Rider confirmation rate at 2×+ | >80% | <60% (UX problem) |
| Surge lift (% trips above 1×) | Informational | — |
| Worker Kafka consumer lag | <60s | >120s |

---

## Trade-offs Summary

| Decision | Alternative | Uber's choice | Rationale |
|---|---|---|---|
| 60s update cycle | 10s cycle | 60s | Saves 6× compute; lag imperceptible |
| Zone aggregation at res-6 | Per-cell (res-9) | Res-6 aggregation | Smooths noise; fewer zones |
| Cap at 8× | No cap | 8× cap | Regulatory + brand risk |
| Upfront price lock | Dynamic in-trip price | Lock at request | Trust + support ticket reduction |
