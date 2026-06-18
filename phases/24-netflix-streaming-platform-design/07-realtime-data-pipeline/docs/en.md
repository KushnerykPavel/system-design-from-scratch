# Real-Time Analytics Pipeline — Kafka, Flink, Druid

> Netflix generates billions of events per day. The challenge is not storing them — it is making them useful in seconds, not hours.

**Type:** Build  
**Company focus:** Netflix  
**Learning goal:** Design the streaming analytics platform. Cover event ingestion via Kafka, stream processing with Flink (windowed aggregations, sessionization), OLAP serving with Druid/Pinot, late data handling, exactly-once semantics, schema evolution, and the playback event → recommendation signal pipeline.  
**Prerequisites:** `06-ab-testing-platform`, `04-recommendation-engine`  
**Estimated time:** ~75 min  
**Primary artifact:** streaming analytics design doc + event schema spec  

## The Problem

Netflix devices emit billions of events per day: video play starts, pauses, seeks, buffering events, quality switches, clicks, search queries, and UI impressions. These events must be:

1. Ingested reliably at high throughput.
2. Processed in near-real-time for operational metrics (anomaly detection, QoE monitoring).
3. Joined with experiment assignments for A/B analysis.
4. Aggregated for OLAP queries by analysts and dashboards.
5. Converted into recommendation signals that update the personalization models.

Design this pipeline end to end.

## Clarify

- What is the latency target for real-time metrics? (seconds to minutes)
- What is the latency target for the batch path? (daily model retraining vs hourly incremental)
- What is the event schema? (structured JSON? Avro? Protobuf?)
- How long should raw events be retained? (Kafka TTL, cold storage retention)
- Who are the consumers of the analytics platform? (data scientists, dashboards, recommendation model trainers, ops on-call)
- What are the consistency requirements? (exactly-once? at-least-once? best-effort?)

## Requirements

### Functional

- Ingest billions of events per day from devices globally.
- Route events to the correct processing topology (real-time ops metrics, experiment metrics, recommendation signals).
- Perform windowed aggregations (per-minute, per-hour, per-day).
- Compute session-level metrics (session duration, total playback time, title completion rate).
- Serve aggregated metrics at low latency for dashboards and anomaly detection.
- Feed recommendation signal updates to the feature store.

### Non-functional

- Ingestion throughput: 1M+ events per second at peak.
- End-to-end latency (event to dashboard): under 60 seconds for operational metrics.
- Late data handling: events arriving up to 10 minutes late must be correctly incorporated.
- At-least-once delivery with idempotent processing (true exactly-once is a myth in distributed systems).
- Schema evolution: adding fields to an event schema must not break existing consumers.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Events per day | ~3 billion | Kafka partition and retention sizing |
| Peak events per second | ~1M | Kafka producer throughput and consumer lag tolerance |
| Event size (Avro) | ~500 bytes average | Kafka storage: ~1.5 TB/day raw before replication |
| Flink parallelism | 100–1,000 tasks depending on job | cluster sizing |
| Druid segment size | ~500 MB per hour per data source | OLAP storage and query planning |
| Recommendation signal update cadence | every 5–15 minutes | drives feature store write throughput |

## Architecture

```text
devices (smart TV, phone, browser)
  -> event collector (HTTP endpoint, validated and buffered)
  -> Kafka (partitioned by user_id, topic per event type)
     -> [real-time branch] Flink streaming job
        -> windowed aggregations (1-min, 5-min, 1-hour windows)
        -> sessionization (group events by session_id + inactivity gap)
        -> anomaly detection (QoE: buffer stall rate, playback start failures)
        -> Druid / Apache Pinot (OLAP, low-latency queries)
        -> Grafana / internal dashboards
     -> [experiment branch] Flink join job
        -> join event with assignment service (by user_id + experiment_id)
        -> experiment metrics by treatment group
        -> experiment dashboard
     -> [recommendation branch] Flink feature job
        -> extract implicit feedback (completed, skipped, paused early)
        -> update subscriber short-term features in feature store (EVCache)
        -> batch sink to S3 for daily model retraining
```

### Kafka Design

Each event type is a separate Kafka topic:

| Topic | Events | Partitions |
|-------|--------|------------|
| playback-start | play button tapped | 128 |
| playback-heartbeat | quality/buffer telemetry (every 30s) | 256 |
| playback-end | stop/completion | 128 |
| ui-interaction | clicks, searches, impressions | 256 |
| quality-event | buffer stall, bitrate switch | 128 |

Partition key: `user_id` (ensures all events for a user go to the same partition, enabling stateful processing without shuffle).

### Flink Stream Processing

**Windowed aggregations:**

```text
stream = kafka.source(playback-start)
  .keyBy(event -> event.title_id)
  .window(TumblingEventTimeWindows.of(Duration.ofMinutes(1)))
  .aggregate(PlaybackCountAggregator)
  -> sink to Druid
```

**Sessionization:**

```text
stream = kafka.source(all-events)
  .keyBy(event -> event.user_id)
  .window(SessionWindows.withGap(Duration.ofMinutes(30)))
  -> compute session_duration, titles_watched, completion_rate
  -> sink to feature store + S3
```

**Late data handling:**

Flink uses event time (timestamp embedded in the event) with a configurable watermark:

```text
watermark = max_seen_event_time - 10 minutes
events arriving after watermark is past their window: sent to side output (late data sink)
late data sink: reprocessed hourly to correct OLAP segments
```

### Exactly-Once Semantics in Practice

True exactly-once across Kafka source + Flink + Druid sink requires:
1. Kafka consumer offsets committed only after successful processing.
2. Flink checkpointing with two-phase commit to the sink.
3. Idempotent sink writes (deduplication key).

Netflix's approach: **at-least-once delivery + idempotent writes**. Processing the same event twice produces the same aggregated result because aggregation functions are designed to be idempotent (count distinct, replace-not-accumulate for state).

### Schema Evolution

Events are serialized in Avro with a schema registry (Confluent Schema Registry):

| Evolution type | Safe? | Notes |
|---------------|-------|-------|
| Add optional field with default | Yes | Old consumers ignore it; new consumers use it |
| Remove optional field | Yes | New consumers use default; old consumers still see it |
| Add required field | No | Breaks old producers that do not populate it |
| Rename field | No | Breaks all consumers |
| Change field type | No | Type incompatibility |

All schema changes must be backward-compatible (new writers, old readers) and forward-compatible (old writers, new readers) to allow rolling deploys.

### Playback Events → Recommendation Signal

```text
playback-end event (user_id, title_id, completion_percent, duration_watched_seconds)
  -> Flink filter: completion_percent > 80% = strong positive signal
  -> Flink filter: completion_percent < 15% = negative signal (did not engage)
  -> update feature: subscriber_short_term_preferences[user_id][genre] += signal_weight
  -> write to EVCache (TTL: 7 days)
  -> batch to S3 for daily training run
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Kafka consumer lag grows | Lag metric exceeds alert threshold | Scale Flink parallelism; prioritize high-value event types; drop or compress low-priority events |
| Flink job fails mid-window | Job restarts from last checkpoint | Checkpointing every 30 seconds limits reprocessing to 30s of events |
| Druid ingestion falls behind | Query results stale; lag metric alert | Druid can scale ingestion horizontally; degrade to batch queries during catch-up |
| Event schema breaks backward compatibility | Consumer parse errors spike | Schema registry rejects incompatible schema versions; deploy gates require schema validation |
| Late events cause incorrect session metrics | Session completion rate diverges from batch | Late data side output reprocessed in next hourly batch; metrics marked as preliminary |

## Observability

- metric: Kafka consumer lag by topic, consumer group, and partition
- metric: Flink job throughput and checkpoint duration
- metric: event-to-dashboard latency (event timestamp vs dashboard render time)
- metric: late event rate by topic and window
- metric: feature store write latency and error rate
- log: schema registry validation failures with schema version and error
- alert: consumer lag > 5 minutes for any high-priority topic

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Separate Kafka topics per event type | Independent scaling, retention, and consumer groups per event type | More topics to manage | Single topic with event type field is simpler but prevents per-type tuning |
| Event-time processing with watermarks | Correct handling of out-of-order events; windows are accurate | Late events require a side output path | Processing-time windows are simpler but incorrect when devices have clock skew or network delays |
| At-least-once + idempotent writes | Simpler than true exactly-once; resilient to failures | Must design aggregation functions to be idempotent | True exactly-once adds significant complexity and latency |
| Avro + schema registry | Schema evolution safety; compact binary encoding | Schema registry is a critical dependency | JSON is simpler but 3-5x larger and has no schema enforcement |

## Interview It

**Netflix framing:** "Design Netflix's real-time analytics pipeline." Strong answers cover Kafka partitioning strategy, Flink processing topology, late data handling, and how events flow into the recommendation signal pipeline. Weak answers describe only "store events in S3 and query with Spark" (batch, not streaming).

**Follow-ups:**
1. How do you handle a Kafka consumer that falls 10 minutes behind?
2. What happens to a Flink window if the job crashes 30 seconds before the window closes?
3. How do you evolve the playback event schema to add a new field without breaking downstream consumers?
4. How would you design the pipeline to detect a sudden spike in buffer stall events across a region?
5. How do recommendation signals get from a playback event to a subscriber's next homepage recommendation?

## Ship It

- `outputs/design-doc-realtime-data-pipeline.md`
- `outputs/event-schema-spec.md`
- `outputs/interview-card-realtime-data-pipeline.md`

## Exercises

1. **Easy** — Size the Kafka storage requirement for one week of raw events (3B events/day, 500 bytes/event, 3x replication).  
2. **Medium** — Design the Flink job that detects a regional spike in buffer stall rate and alerts within 60 seconds.  
3. **Hard** — Design the late data reprocessing pipeline that corrects OLAP segment data for events arriving up to 10 minutes late.  

## Further Reading

- [Netflix Keystone data pipeline](https://netflixtechblog.com/keystone-real-time-stream-processing-platform-a3ee651812a)  
- [Apache Flink documentation](https://nightlies.apache.org/flink/flink-docs-stable/)  
- [Apache Druid architecture](https://druid.apache.org/docs/latest/design/architecture.html)  
