# Ride-Sharing Dispatch System

> Matching a rider to a driver is a geospatial optimization problem masquerading as a simple lookup. Surge pricing, ETA accuracy, and supply rebalancing are where the interview actually lives.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Design the dispatch core of a ride-sharing platform: driver location ingestion, geospatial matching, ETA estimation, surge pricing signal generation, and supply rebalancing — not just "find the nearest driver."  
**Prerequisites:** `20-low-latency-location-and-market-systems/01-proximity-service`, `20-low-latency-location-and-market-systems/02-location-updates`, `07-queues-streams-and-workflows/05-workflow-engines`  
**Estimated time:** ~90 min  
**Primary artifact:** capacity sheet

## The Problem

Design the core dispatch system for a ride-sharing service operating in 100 cities. Riders request a trip, and the system must match them to an available nearby driver, estimate arrival time, charge appropriate pricing, and manage the ongoing trip state through pickup and dropoff.

The trap that catches mid-level candidates is treating dispatch as a database lookup: "find the driver with the minimum distance from the rider's coordinates." This ignores the core operational problems. Driver locations change every 4 seconds, producing enormous write throughput. The matching problem is not just minimum distance but minimum ETA, which requires live traffic conditions. Multiple competing riders may attempt to match the same driver simultaneously, requiring an atomic assignment protocol. And surge pricing is not a simple multiplier — it is a real-time demand-supply ratio signal that must respond within seconds to local supply shocks, yet must not oscillate so rapidly that drivers are confused.

## Clarify

- What is the matching objective: minimum ETA, or a weighted objective that includes driver acceptance rate, rider rating, and route efficiency?
- Should the system support ride pooling (matching multiple riders to one driver on overlapping routes), or single-passenger trips only?
- What is the required match latency: how long can a rider wait before receiving a driver assignment?

If the interviewer does not specify, assume single-passenger trips, ETA-minimizing match objective, and match latency target of under 5 seconds for 95% of requests in dense urban areas.

## Requirements

### Functional

- Accept driver location updates continuously from active driver apps.
- Accept ride requests from riders with pickup and dropoff locations.
- Match riders to available nearby drivers, preferring minimum ETA.
- Provide fare estimates before ride confirmation and charge final fare at completion.
- Compute and display surge pricing multipliers per geographic zone.
- Manage trip state machine: requested → matched → in-progress → completed or cancelled.

### Non-functional

- Driver location update ingestion: 500K active drivers × 1 update/4s = 125K writes/s.
- Match latency p95 under 5 seconds in dense cities; up to 30 seconds in low-density areas.
- Location data freshness: driver positions must be no more than 8 seconds stale when used for matching.
- High availability: matching service 99.99% uptime; a city-level outage should not affect other cities.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active drivers | 500K globally; 100K in peak city | determines location store write throughput and geospatial index size |
| Location updates | 125K writes/s (1 update per driver per 4s) | drives in-memory geospatial index update rate |
| Ride requests | 5K match requests/s at peak globally | concurrent match operations competing for the same driver pool |
| Geospatial index | 100K drivers per city × 2D coords = fits in 50 MB per city | enables per-city in-memory geohash index on dispatch nodes |
| Surge zones | 500 zones per city × 100 cities = 50K surge zones | surge state fits in a Redis cluster; must be updated every 30s |

## Architecture

```text
driver app
  -> location update (GPS coords, heading, timestamp)
  -> location ingestion service
  -> location store (Redis geospatial, per city)
  -> broadcast to local dispatch nodes via pubsub

rider app
  -> ride request (pickup, dropoff, desired ride type)
  -> matching service

matching service
  -> query location store: drivers within R km of pickup, sorted by ETA
  -> compute ETA using road network + live traffic (OSRM / Google Maps API)
  -> attempt atomic assignment (driver state: available -> pending)
  -> if driver declines or times out: retry with next candidate
  -> confirm assignment: create trip record, notify rider and driver

trip state machine
  -> states: requested, driver_assigned, pickup_en_route, arrived, in_progress, completed, cancelled
  -> transitions driven by driver app events and rider actions
  -> stored in relational DB (trips table) with state + timestamps

surge pricing
  -> demand signal: ride requests per zone per minute (windowed count)
  -> supply signal: available drivers per zone (from location store)
  -> surge ratio: demand / supply; multiplier from lookup table
  -> zone surge updated every 30s, propagated to rider app on next request
```

## Data Model & APIs

Core entities:

```text
Driver  { driver_id, current_lat, current_lng, status, vehicle_type, rating, city_id }
Rider   { rider_id, home_address, payment_method_id, rating }
Trip    { trip_id, rider_id, driver_id, status, pickup_loc, dropoff_loc,
          fare_estimate, fare_final, surge_multiplier, requested_at, matched_at,
          pickup_at, dropoff_at }
SurgeZone { zone_id, city_id, geohash_prefix, multiplier, updated_at }
```

Key APIs:

- `PUT /v1/drivers/{driver_id}/location` — body: `{lat, lng, heading, speed, timestamp}`; idempotent by timestamp
- `POST /v1/rides` — body: `{pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, ride_type}`; returns `{ride_id, eta_seconds, fare_estimate, surge_multiplier}`
- `POST /v1/rides/{ride_id}/confirm` — rider confirms the fare estimate and initiates matching
- `GET /v1/rides/{ride_id}` — returns trip status, driver location, ETA to pickup
- `GET /v1/surge?lat=X&lng=Y` — returns current surge multiplier for the given location

Driver assignment is atomic: a Redis SETNX or compare-and-swap on `driver:{driver_id}:status` transitions from `available` to `pending:{ride_id}`. Only the process that wins the CAS may send the dispatch request to the driver.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Driver declines or does not respond to dispatch | no acceptance within 15s, match timeout alert | retry matching with next-best driver from candidate list; limit retries to 3 before returning no match |
| Location store partition loses driver positions for a city | stale driver count alert; match latency spike | fall back to last-known position with staleness flag; trigger re-sync from driver apps |
| Surge zone computation service crashes | surge multiplier stuck at last value; no update within 2 × update interval | serve last-known surge value; alert on-call; surge computation is non-critical for trips in progress |
| Two match processes assign same driver to two rides | duplicate trip_id on same driver_id in DB; driver app reports conflict | atomic CAS on driver status prevents this; conflict detection in driver app triggers re-dispatch for loser |

## Observability

- metric: match latency p50/p95/p99 per city — the primary user-visible SLO
- metric: driver acceptance rate per city — falling rate signals poor match quality or driver pool issues
- metric: driver location data age histogram — detects stale position data used in matching
- metric: surge zone update latency — how long after demand shifts before the zone reflects it
- log: every match attempt with rider_id, candidate drivers, ETA estimates, final assigned driver, and total match duration
- SLO: p95 of matches in tier-1 cities complete within 5 seconds; driver position freshness > 95% under 8 seconds

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| In-memory geospatial index per city on dispatch nodes | sub-millisecond geospatial lookup; no DB round-trip for candidate set | location store must be updated with every driver ping; consistency across nodes requires pubsub | PostgreSQL PostGIS on every match — works at small scale, too slow at 125K updates/s and 5K matches/s |
| Atomic CAS for driver assignment instead of reservation queue | simple, fast; no extra coordination service needed | a few milliseconds of lock contention when multiple riders compete for the same driver | pessimistic locking in RDBMS — would create a transaction bottleneck at the match layer |
| Zone-level surge pricing updated every 30s | smooth enough for rider UX; prevents oscillation | 30-second lag means a sudden supply shock takes up to 30s to reflect in pricing shown to riders | per-request dynamic pricing — accurate but confusing to riders who see price change between estimate and confirm |

## Interview It

**Google framing:** "Design the dispatch backend for a ride-hailing service in 50 cities. Focus on matching accuracy and handling driver supply shocks." Expect pushback on ETA accuracy under live traffic and how surge pricing avoids oscillation.

**Cloudflare framing:** "How would you architect the location ingestion and geospatial query layer to run at the edge, closer to drivers?" Expect questions on edge state consistency, geospatial index synchronization across POPs, and what happens when edge nodes partition.

**Follow-ups:**
1. How would you add ride pooling where two riders with overlapping routes share one car?
2. What happens to in-flight trips when the matching service node handling a city crashes?
3. How would you implement driver supply rebalancing that suggests drivers move to high-demand zones?
4. If you had to reduce match latency to under 1 second, what would you change?
5. How would you detect and suppress a driver location spoofing attack?

## Ship It

- `outputs/capacity-sheet-ride-sharing-dispatch.md`

## Exercises

1. **Easy** — Calculate how many geospatial index updates per second the location store must handle for 500K active drivers each sending one GPS ping every 4 seconds.
2. **Medium** — Design the retry protocol when the first driver declines: who holds the list of candidate drivers, how long is each candidate given to respond, and when do you give up and tell the rider no match was found?
3. **Hard** — Design a supply rebalancing system that issues directional suggestions to idle drivers. How do you prevent all idle drivers from moving to the same zone simultaneously?

## Further Reading

- https://eng.uber.com/uber-speed-scale/ — Uber's engineering blog on dispatch system evolution from monolith to microservices, including geospatial indexing choices
- https://s3.amazonaws.com/systemsandpapers/papers/amazon.pdf — Amazon Dynamo paper; useful mental model for the consistency vs availability tradeoff in driver state management
- https://engineering.lyft.com/how-lyft-thinks-about-location-systems-bd6f8e617b69 — Lyft's approach to location system design and the tradeoffs between accuracy, freshness, and cost
