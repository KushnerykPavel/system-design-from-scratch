# Real-Time GPS Ingestion & Driver Tracking
> 5 million drivers sending location every 4 seconds is 20 million events per second at peak.

**Type:** Build
**Company focus:** Uber
**Learning goal:** Design a location ingestion pipeline that processes 20M GPS events/sec, maintains last-known position per driver with <1s staleness, and feeds dispatch and surge pricing systems.
**Prerequisites:** `07-queues-streams-and-workflows/03-consumer-groups`, `15-kv-cache-and-object-storage/01-distributed-kv-store`
**Estimated time:** ~75 min
**Primary artifact:** location ingestion architecture diagram

---

## Scale Context

At peak, Uber operates approximately 5 million active drivers globally. Each driver's mobile SDK sends a GPS update every 4 seconds:

```
5,000,000 drivers × (1 event / 4 seconds) = 1,250,000 events/second baseline
```

With geographic clustering (rush hour in 50 cities simultaneously), the peak rate reaches **~20 million events/second**. This is not a database workload — it is a streaming ingestion problem.

---

## Location Event Structure

Each GPS update contains:

```
LocationEvent {
    driver_id:        "driver_abc123"    // globally unique
    lat:              40.7128            // WGS84 decimal degrees
    lon:              -74.0060
    bearing:          270.0              // degrees from north (0–360)
    speed:            35.5               // km/h
    accuracy_meters:  8.0               // GPS accuracy estimate
    timestamp:        1719000000000      // Unix milliseconds (device clock)
}
```

`accuracy_meters` is provided by the device OS. It represents the radius of the 68% confidence circle around the reported position. Lower is better.

---

## Ingestion Pipeline

```
Mobile SDK (4Hz)
       │
       ▼ HTTPS/gRPC
Uber Gateway (load balanced, ~500 nodes)
       │
       ▼ Kafka producer (keyed by driver_id)
Kafka: driver-locations topic
  · ~10,000 partitions
  · partitioned by driver_id (ordering per driver)
  · retention: 24 hours
       │
       ├──────────────────────────────────────┐
       ▼                                      ▼
Location Workers (consumer group)    Surge Workers (consumer group)
  · validate event                     · aggregate by H3 zone
  · decode H3 cell                     · compute multiplier
  · write to Redis (driver position)
  · publish H3 cell change events
       │
       ▼
Redis Cluster (driver position store)
  · key: driver:{id}:location → LocationEvent JSON
  · TTL: 30 seconds
  · key: driver:{id}:cell → H3 cell ID
```

### Why 10,000 Kafka Partitions?

With 20M events/sec and roughly 2000 events/sec per partition throughput budget, the math yields:

```
20,000,000 / 2,000 = 10,000 partitions
```

Partitioning by `driver_id` ensures all events from one driver land in the same partition, giving the location worker a consistent per-driver ordered stream without cross-partition coordination.

---

## Location Workers

Each location worker:

1. **Reads** a batch of events from its Kafka partition(s)
2. **Validates** each event (accuracy, recency, velocity)
3. **Decodes** the lat/lon to an H3 cell ID
4. **Detects** cell transitions (driver moved to a new cell)
5. **Batches** 100ms of writes, then flushes to Redis via pipeline
6. **Publishes** H3 cell change events to `driver-cell-changes` Kafka topic
7. **Commits** the Kafka offset after successful Redis write

The 100ms batch window reduces per-event Redis round-trips by ~100×. At 2000 events/sec per worker, this batches ~200 events per pipeline flush.

---

## Validation Rules

Location workers apply three filters:

### 1. Accuracy Filter
```
if event.AccuracyMeters > 50 → drop event
```
GPS readings with > 50m accuracy are too imprecise for dispatch matching. This typically occurs in urban canyons (tall buildings blocking satellite view) or when the driver has just started the app.

### 2. Staleness Filter
```
if now - event.Timestamp > 5 seconds → drop event
```
The driver SDK timestamps events using the device clock. Events older than 5 seconds indicate network delay or a reconnection replay from a crash. Accepting them would cause position "jumps" backward in time.

### 3. Velocity (Teleportation) Filter
```
if implied_speed_from_last_position > 300 km/h → drop event
```
A driver cannot move from lat/lon A to lat/lon B faster than ~300 km/h (the threshold for an impossible vehicle speed). Events that imply faster movement are either GPS glitches, device clock errors, or spoofing attempts.

---

## Redis Position Store

Each driver's current position is stored in Redis as:

```
SET driver:{driver_id}:location <LocationEvent JSON> EX 30
SET driver:{driver_id}:cell <H3CellID> EX 30
```

The 30-second TTL ensures:
- Drivers who go offline are automatically removed from dispatch candidate queries
- No manual cleanup process is needed
- Dispatch service can check `EXISTS driver:{id}:location` to detect staleness

### Staleness Detection

The dispatch service considers a driver **UNREACHABLE** when:
```
TTL(driver:{id}:location) == -2  (key expired or never existed)
```

This fires approximately 30 seconds after the last valid GPS event. For drivers on trips, the trip is kept `IN_PROGRESS` even if the driver goes unreachable — the trip state machine handles recovery separately.

---

## Cassandra: Historical Location Trail

While Redis stores the **current** position (low latency, short TTL), Cassandra stores the **historical trail** for:
- Trip replay (reconstruct the exact route for earnings calculation)
- Incident investigation (where was the driver at time T?)
- Dispute resolution (was the driver actually at the pickup location?)

Schema:
```
CREATE TABLE driver_location_history (
    driver_id   TEXT,
    ts          TIMESTAMP,
    lat         DOUBLE,
    lon         DOUBLE,
    speed       DOUBLE,
    PRIMARY KEY (driver_id, ts)
) WITH CLUSTERING ORDER BY (ts DESC)
  AND default_time_to_live = 7776000;  -- 90 days
```

Cassandra's wide-row model makes per-driver time range queries efficient. The 90-day TTL covers regulatory retention requirements.

---

## Geofencing

Location events also trigger geofence evaluations:

- **Airport zones**: entering an airport geofence switches driver to airport queue rules
- **Surge zones**: entering/exiting a surge zone triggers surge assignment refresh
- **Restricted areas**: entering a no-pickup zone triggers driver notification

Geofence containment is O(1) using H3 polyfill: the geofence polygon is pre-converted to a set of H3 cells. Containment = Redis `SISMEMBER geofence:{name}:cells {driver_cell}`.

---

## GPS Spoofing Detection

Drivers occasionally use GPS spoofing apps to fake their location (e.g., appearing in a surge zone without physically going there). Detection layers:

| Signal | Detection Method |
|---|---|
| Velocity impossible (> 300 km/h) | Teleportation filter |
| Perfect grid movement | Statistical: no natural GPS jitter |
| Device sensor mismatch | Accelerometer says stationary, GPS says moving |
| Multiple drivers, same coordinates | Clustering detection (shared spoofing service) |

Detected spoofers are flagged for manual review, not immediately deactivated, to avoid false positives on legitimate GPS glitches.

---

## Failure Modes

| Failure | Impact | Mitigation |
|---|---|---|
| Kafka consumer lag spike at peak hour | Location data delayed → stale dispatch | Auto-scale workers; monitor consumer group lag; alert at >10s |
| Redis cluster partition | Driver positions unavailable in affected shard | Replica promotion; dispatch marks affected drivers UNREACHABLE |
| GPS spoofing | Drivers earn surge rates without being in zone | Velocity filter + statistical detection; manual review queue |
| Device clock drift | Events arrive with wrong timestamp | Use server-side arrival time as fallback for staleness check |
| Network batch replay on reconnect | Old events flood pipeline | Staleness filter (> 5s) discards replayed historical events |

---

## Observability

| Metric | SLO | Alert |
|---|---|---|
| Kafka consumer lag P99 | <2s | >10s |
| Redis position freshness P99 | <1s staleness | >3s |
| Event validation drop rate | <5% | >15% (accuracy or clock issue) |
| Teleportation detection rate | <0.1% | >1% (spoofing spike) |
| Location worker throughput | >2000 events/sec/worker | <1000 (worker degraded) |

---

## Trade-offs Summary

| Decision | Alternative | Chosen | Rationale |
|---|---|---|---|
| Kafka partitioned by driver_id | Partitioned by zone | driver_id | Ordering per driver prevents position jumps |
| Redis TTL 30s | Explicit DELETE on driver offline | TTL | Self-cleaning; no coordination needed |
| 100ms batch write | Per-event write | Batching | 100× fewer Redis round-trips |
| Cassandra for history | PostgreSQL | Cassandra | Wide-row model efficient for per-driver time ranges at scale |
