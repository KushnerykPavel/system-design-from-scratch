# Uber Full Mock Loop
> A strong-hire answer is not about knowing every component — it is about making deliberate trade-offs and defending them.

**Type:** Mock Interview
**Company focus:** Uber
**Learning goal:** Practice a complete 45-minute system design interview for Uber's real-time ride-sharing platform, calibrate against strong-hire and weak-hire patterns, and identify personal weak spots.
**Prerequisites:** All lessons 01–06 in this phase.
**Estimated time:** 45 min interview + 30 min debrief
**Primary artifact:** self-scored interview rubric

---

## Interview Prompt

Choose one of the following prompts based on your preparation level:

**Broad prompt** (senior engineer, ~5 years exp):
> "Design Uber's real-time ride-sharing platform."

**Focused prompt** (staff engineer, system depth expected):
> "Design the Uber driver location tracking and dispatch notification system."

The focused prompt is harder because the interviewer expects you to go deeper on latency, partitioning, and failure modes in a narrower domain rather than surveying the whole system.

---

## Capacity Numbers to Memorize

| Metric | Value |
|---|---|
| Active drivers (peak) | 5 million |
| Active riders | 100 million registered |
| Daily trips | 25 million |
| Peak trips/second | ~290 trips/sec |
| Driver location update frequency | Every 4 seconds |
| Peak location events/second | **20 million/sec** (5M × 0.25Hz × safety factor) |
| Driver position staleness SLO | <1 second |
| Rider wait time SLO P50 | <3 minutes |

These numbers anchor every subsequent trade-off. State them in the first 5 minutes.

---

## Milestone Map (45 minutes)

### Minute 0–5: Requirements & Scope
- Clarify: global or single city? (answer: design for global, implement for one region)
- Clarify: all products or UberX only? (answer: UberX is fine for this session)
- State functional requirements: driver location ingestion, dispatch matching, trip lifecycle, surge pricing
- State non-functional requirements: 20M location events/sec, <1s position staleness, <3 min rider wait

### Minute 5–15: High-Level Architecture
- Draw the major components: mobile apps → gateway → Kafka → location workers → Redis → dispatch → trip service
- Identify three subsystems: location ingestion, dispatch matching, trip state machine
- Mention Cassandra for historical trails, Temporal/Cadence for trip workflow

### Minute 15–25: Location Ingestion Deep-Dive
- Kafka partitioned by driver_id (ordering), 10K partitions for 20M events/sec
- Redis TTL 30s for driver position (self-cleaning offline detection)
- Validation: accuracy filter (>50m drop), staleness filter (>5s drop), velocity filter (>300km/h drop)
- H3 cell update on cell boundary crossing → publishes to dispatch
- GPS spoofing detection strategy

### Minute 25–35: Dispatch & Surge Deep-Dive
- H3 geospatial indexing: res-9 for driver store, res-6 for surge zones
- Batch matching: 500ms window, Hungarian algorithm for optimal assignment
- Surge pipeline: Kafka consumer → zone aggregation → multiplier computation → Redis → push
- Upfront pricing: multiplier locked at request time, 2× requires explicit confirmation

### Minute 35–45: Failure Modes & Observability
- Driver app crash during trip: heartbeat TTL + dead trip detection
- State machine idempotency: same-state transitions are no-ops
- Billing double-charge prevention: idempotency key on payment processor
- Key metrics: Kafka consumer lag, Redis TTL freshness, rider wait P50/P95, ETA error rate

---

## Strong-Hire Outline

A strong-hire answer explicitly addresses:

1. **H3 geospatial indexing** — explains why H3 hex cells beat lat/lon grid or PostGIS for dispatch, mentions k-ring for proximity search, mentions polyfill for geofencing
2. **Kafka for location ingestion** — partitioned by driver_id for ordering, justifies 10K partitions with math, mentions consumer group isolation between location workers and surge workers
3. **Redis for driver store with TTL** — explains that TTL is the offline detection mechanism, not a background job; mentions pipeline writes for batching
4. **Batch matching with DISCO** — explains that greedy nearest driver is locally suboptimal; mentions 500ms batch window as the latency budget; mentions Hungarian algorithm by name or describes bipartite matching
5. **Trip state machine with Temporal/Cadence** — names durable execution as the recovery mechanism; explains idempotent transitions; explains billing idempotency key
6. **Surge pricing pipeline** — explains ratio computation, threshold tiers, multiplier lock at request time
7. **GPS spoofing detection** — velocity filter at minimum; bonus for statistical approaches

---

## Weak-Hire Anti-Patterns

Avoid these responses — they signal shallow preparation:

| Anti-Pattern | Why It's Wrong |
|---|---|
| "Store driver locations in a database and query with SELECT WHERE lat BETWEEN..." | O(N) full table scan on 5M rows every dispatch cycle — guaranteed SLA miss |
| "Use PostGIS radius query" | Doesn't explain why H3 is better; misses the fan-out / TTL advantages |
| "Poll driver locations every second from the app" | 5M drivers × 1s poll = 5M requests/sec to a REST API — no horizontal scaling story |
| No state machine for trips | Implies billing can fire on cancelled trips; no idempotency story |
| "Use a single Kafka topic with no partitioning strategy" | Cannot guarantee per-driver ordering; location workers get mixed driver streams |
| "Surge is just a percentage markup" | Misses the supply/demand ratio computation and the zone aggregation strategy |
| No mention of failure modes | System design without failure mode discussion is a weak answer at L5+ |

---

## Follow-Up Questions & Ideal Answers

### Q1: How would you handle a surge pricing calculation during a Redis cluster failover?
**Ideal answer**: Redis serves the last-known multiplier with a 120s TTL. During failover (typically <30s), stale multipliers are served from the remaining replicas. If a replica is unavailable, the system falls back to 1.0× (safe degradation). The impact is brief under-surge for affected zones, which is acceptable vs. blocking trip requests entirely.

### Q2: A concert ends and 50,000 people request Uber simultaneously. How does the system behave?
**Ideal answer**: The surge worker detects the demand spike within one 60s cycle and publishes elevated multipliers. Kafka consumer lag may spike briefly; the location worker auto-scales. The dispatch batch window may temporarily extend to absorb the backlog. Pre-positioned drivers (from the supply forecasting signal) are already in the zone, reducing match time. The surge multiplier begins attracting additional drivers. Peak is absorbed within 5–10 minutes.

### Q3: How do you prevent a driver from receiving duplicate trip offers?
**Ideal answer**: DISCO maintains a driver offer state in Redis (`driver:{id}:pending_offer`). Before sending an offer, it checks this key. The driver's acceptance or decline clears the key. The offer also has a 15-second timeout with a Redis TTL, after which the slot is reopened.

### Q4: What happens if two dispatch workers both try to assign a driver to different riders simultaneously?
**Ideal answer**: The driver offer state in Redis acts as a distributed lease. DISCO uses a Redis SET NX (set if not exists) with a 15s TTL as a compare-and-swap. Only the first worker to acquire the key wins the driver. The second worker detects the key already exists and removes that driver from its candidate pool.

### Q5: How would you design the Uber driver earnings settlement system?
**Ideal answer**: This is a separate bounded context from dispatch. Each trip produces a COMPLETED event with a fare amount. The earnings ledger service consumes these events and credits driver accounts. Weekly payouts batch these credits into a single ACH transfer. The key design constraint is that the ledger must be idempotent (tripID as idempotency key) and append-only for auditability. Redis is not appropriate here — this is a financial ledger requiring durable relational storage (PostgreSQL with a trips_earnings table).

---

## Scoring Rubric

| Dimension | Weight | What earns full marks |
|---|---|---|
| Geospatial design | 20 pts | H3 indexing, k-ring queries, polyfill for geofencing, explains why over alternatives |
| Location pipeline | 20 pts | Kafka partitioning strategy, Redis TTL design, validation filters, batching |
| Dispatch algorithm | 15 pts | Batch matching rationale, 500ms window trade-off, top-3 offer protocol |
| Trip state management | 15 pts | All states named, idempotent transitions, Temporal for durable recovery |
| Failure recovery | 15 pts | Dead trip detection, billing idempotency, Redis failover, state split resolution |
| Observability | 15 pts | Key metrics named with SLOs, Kafka lag monitoring, alert thresholds |

**Score bands:**
- 85–100: Strong Hire
- 70–84: Hire (some rough edges)
- 55–69: Weak Hire (needs coaching on one or two dimensions)
- < 55: No Hire (fundamental gaps)

---

## Self-Assessment Checklist

Before reviewing your answer, check off each item you explicitly mentioned:

- [ ] H3 hexagonal indexing for driver locations
- [ ] Kafka partitioned by driver_id with partition count math
- [ ] Redis TTL as offline detection mechanism
- [ ] Location event validation (accuracy, staleness, velocity)
- [ ] 500ms batch matching window with optimality justification
- [ ] Top-3 driver offer with 15s timeout
- [ ] Trip state machine with all states
- [ ] Idempotent state transitions
- [ ] Temporal/Cadence for durable trip workflow
- [ ] Surge zone aggregation (H3 res-6)
- [ ] Surge multiplier tiers and cap
- [ ] Upfront price lock at request time
- [ ] Dead trip detection (10 min heartbeat timeout)
- [ ] At least two failure modes with mitigations
- [ ] At least three key metrics with SLOs

Score 1 point per checked item. 12–15: strong hire territory. 8–11: hire. <8: study needed.
