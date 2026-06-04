# Real-Time Location Update Pipeline

> A location pipeline is not just an ingestion problem; it is a trust problem about how stale, duplicated, or reordered motion you are willing to serve.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a high-rate location update pipeline with dedupe, smoothing, bounded out-of-order handling, and latency-aware serving integration.
**Prerequisites:** `07-queues-streams-and-workflows/03-consumer-groups`, `10-reliability-retries-and-backpressure/05-async-backpressure`, `20-low-latency-location-and-market-systems/01-proximity-service`
**Estimated time:** ~75 min
**Primary artifact:** location-pipeline validator + operator checklist

## The Problem

Design the pipeline that receives GPS updates from mobile devices, vehicles, or IoT assets and makes those updates usable for real-time products. Producers are bursty, mobile networks are unreliable, and the serving system must decide when data is fresh enough to trust.

This lesson matters because many designs jump from "device sends coordinates" straight to "query nearby." Senior answers explain dedupe windows, jitter smoothing, sequence handling, and the serving contract for old or missing updates.

## Clarify

- How often do clients send updates under normal operation?
- Do devices send an incrementing sequence number or only timestamps?
- Is the primary goal live visibility, ETA estimation, geofencing, or analytics reuse?
- How stale can an accepted update be before it is excluded from live serving?

If no details are given, assume mobile devices send updates every 2 to 5 seconds while active, each update carries device time plus a monotonic client sequence, and live serving should ignore data older than 15 seconds.

## Requirements

### Functional

- Accept high-rate device location updates with authentication.
- Deduplicate retries and reject clearly stale or malformed events.
- Smooth noisy GPS jumps without hiding real movement.
- Publish accepted updates to both live-serving indexes and downstream analytics.
- Surface per-device freshness and delivery health.

### Non-functional

- Keep end-to-end publish-to-serve latency under 2 seconds for healthy paths.
- Bound duplicate amplification during flaky-network retries.
- Avoid one hot tenant or device cohort overwhelming the stream.
- Preserve enough ordering information to make last-write-wins credible.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active devices | 25M concurrently updating | shapes ingestion partitioning and auth fanout |
| Update rate | 8M events/s peak globally | stream durability and partition count matter immediately |
| Retry amplification | 3x during mobile-network incidents | dedupe storage is part of core capacity |
| Acceptable live staleness | 15 seconds | defines expiration and serving trust policy |
| Late-arrival window | 30 seconds max | controls how much reordering the pipeline absorbs |

## Architecture

```text
devices
  -> regional ingest gateway
  -> auth + schema validation
  -> dedupe / sequence guard
  -> stream partitions by device or region
  -> smoothing / normalization workers
  -> live index publisher + analytics sink
```

Design notes:

1. Separate acceptance, normalization, and serving publication so each stage has a clear responsibility.
2. Use sequence or version checks whenever possible; timestamps alone are weak in mobile environments.
3. Dedupe and stale-drop policies are product decisions, not only data-pipeline details.
4. Publish one trusted "latest live state" view and one fuller event stream for replay or analytics.

## Data Model & APIs

Core records:

```text
raw_update(device_id, seq, client_ts, server_ts, lat, lon, accuracy_m)
accepted_update(device_id, version, lat, lon, heading, speed, accepted_at)
device_freshness(device_id, age_seconds, state)
```

Useful interfaces:

- `POST /v1/location-events`
- `GetLatestLiveState(device_id)`
- `PublishAcceptedUpdate(device_id, version)`
- `MarkDeviceOffline(device_id, reason)`

Strong answers name the contract between the pipeline and nearby-serving systems: what "latest" means, how long it stays valid, and when it is excluded.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| retries flood the pipeline during weak connectivity | dedupe-hit rate and gateway retry histograms | bounded idempotency window and per-device backpressure |
| out-of-order updates move an object backward | stale-drop count and sequence-gap metrics | sequence guard with last accepted version |
| GPS jitter produces fake movement | speed outlier rate and path-jump metrics | smoothing, accuracy thresholds, and confidence flags |
| one region falls behind while others stay healthy | partition lag by region and age-to-serve metrics | isolate regional streams and fail data stale instead of silent |

## Observability

- metric: ingest QPS, auth reject rate, and schema reject rate
- metric: dedupe-hit rate, stale-drop rate, and out-of-order discard count
- metric: publish-to-serve latency and accepted-update age percentiles
- metric: smoothing corrections and impossible-speed detections
- log: rejected updates with device ID, reason class, and version
- trace: ingest gateway through stream, normalization, and live publication
- SLO: 99.9% of valid active-device updates become live-visible within the latency target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit dedupe and sequence guard | retry-safe and order-aware serving | more hot-path state | blind last-write-wins on timestamps only |
| separate live state from raw event history | fast serving and richer replay | dual data products to operate | only raw stream consumed directly by reads |
| smoothing and accuracy filters | less jitter in product behavior | risk of masking real movement | serving every raw GPS point unchanged |

## Interview It

**Google framing:** "Design the real-time location ingestion pipeline behind maps, nearby search, or ETA products." Expect follow-ups on ordering, late events, and how much staleness the product can tolerate.

**Cloudflare framing:** "Design a globally distributed ingest path for devices sending frequent updates." Expect questions on regional isolation, edge termination, and how to keep the path safe under retry storms.

**Follow-ups:**
1. What if devices cannot provide monotonic sequence numbers?
2. How do you keep one tenant from overwhelming shared partitions?
3. What changes if location updates also drive fraud or safety alerts?
4. How would you expose confidence or freshness to downstream products?
5. What if analytics wants the raw stream but product serving wants only smoothed state?

## Ship It

- `outputs/location-pipeline-operator-checklist.md`

## Exercises

1. **Easy** — Explain why server receive time alone is a weak ordering signal for mobile updates.
2. **Medium** — Design a dedupe key and expiry policy for retry-heavy devices.
3. **Hard** — Redesign the pipeline when network partitions cause minutes of delayed updates that arrive in bursts.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — useful for thinking through pipeline lag and graceful degradation
- [Kafka design](https://kafka.apache.org/documentation/) — helpful background for partitioned event ingestion and consumer lag
