# Kafka Architecture — Stream Processing & Events

> LinkedIn built Kafka to unify all data movement — it became the backbone of modern event-driven architecture.

**Type:** Concept
**Company focus:** LinkedIn
**Learning goal:** Understand why LinkedIn built Kafka, how it works internally, and how LinkedIn uses it for the unified activity log, feed updates, notifications, and analytics.
**Prerequisites:** `07-queues-streams-and-workflows/01-queues-vs-streams`, `07-queues-streams-and-workflows/06-outbox-and-cdc`
**Estimated time:** ~75 min
**Primary artifact:** Kafka usage taxonomy at LinkedIn

## The Problem Kafka Was Built to Solve

In 2011, LinkedIn had dozens of data pipelines connecting services in a point-to-point mesh: the activity feed service pulled from the profile service, the recommendation service pulled from the activity service, the analytics pipeline pulled from each. Each new consumer required a new integration. N services × M data sources = O(N×M) pipelines.

The core insight from Jay Kreps, Neha Narkhede, and Jun Rao was: **treat every event as a log entry**. A log is append-only, ordered within a partition, and can be replayed from any offset. By publishing every event to a central log, each producer writes once and every consumer reads at its own pace, from its own position.

This is the **unified log** principle: the log is the source of truth; all state is a derivative of the log.

## Kafka Internals

### Topics and Partitions

A **topic** is a named stream of records. Every record has: a key, a value, a timestamp, and headers.

A topic is divided into **partitions**. A partition is an ordered, immutable log. Records within a partition are assigned monotonically increasing **offsets**.

```
Topic: member-activity
  Partition 0: [offset 0][offset 1][offset 2]...
  Partition 1: [offset 0][offset 1][offset 2]...
  Partition 2: [offset 0][offset 1][offset 2]...
```

Partitioning provides:
- **Parallelism:** multiple consumers can process different partitions concurrently.
- **Ordering:** within a partition, records are strictly ordered. Across partitions, ordering is not guaranteed.
- **Distribution:** partitions are distributed across broker nodes.

### Producers

A producer publishes records to a topic. The partition for each record is determined by:
- **Key hash:** if a key is provided, `partition = hash(key) % num_partitions`. Records with the same key always go to the same partition — guaranteeing ordered processing for a given entity (e.g., all events for member_id=123 go to the same partition).
- **Round-robin:** if no key, records are distributed evenly.

Producers can configure:
- `acks=all` — wait for all in-sync replicas (ISR) to acknowledge before confirming the write (durability guarantee)
- `acks=1` — wait for leader only (faster, less durable)
- `acks=0` — fire and forget

### Consumer Groups

A **consumer group** is a set of consumers that together consume a topic. Kafka assigns each partition to exactly one consumer within the group. This gives:
- **Parallelism:** up to N consumers can process N partitions simultaneously.
- **Exclusive consumption:** no two consumers in the same group process the same partition.

Multiple distinct consumer groups can consume the same topic independently, each maintaining its own offset.

```
Topic member-activity (3 partitions)
  Consumer Group: feed-updater
    consumer-A → partition 0, 1
    consumer-B → partition 2
  Consumer Group: pymk-recompute
    consumer-C → partition 0
    consumer-D → partition 1, 2
```

### Offset Management

Each consumer group commits its offset to the internal topic `__consumer_offsets`. This allows consumers to restart and resume from where they left off.

**Consumer lag** = `latest_offset − committed_offset`. Lag is the key operational metric: it measures how far behind a consumer is from the tip of the log.

### Retention

Kafka retains records either by:
- **Time-based:** default 7 days. Records older than the retention period are deleted.
- **Size-based:** when a partition exceeds a configured size limit, oldest segments are deleted.
- **Compact:** for changelog topics, only the latest record per key is retained (log compaction). Used for Kafka Streams state stores and Venice-style feature propagation.

### Replication

Each partition has one **leader** and zero or more **follower replicas**. The set of replicas that are caught up to the leader is the **In-Sync Replica set (ISR)**.

When `min.insync.replicas=2` and `acks=all`, a write succeeds only if at least 2 replicas (leader + 1 follower) have persisted the record. This prevents data loss if the leader fails immediately after the write.

If the leader fails, one of the ISR followers is elected as the new leader. Consumers and producers transparently reconnect.

## LinkedIn's Kafka at Scale

LinkedIn operates **hundreds of Kafka clusters** processing billions of messages per day across use cases including:

| Use Case | Topic | Producer | Consumer |
|----------|-------|----------|----------|
| Activity events | member-activity | Web/app servers | Feed ranking, PYMK recompute |
| Profile changes | member-profile-updates | Profile service | Elasticsearch sync, PYMK |
| Job applications | job-application-events | Jobs service | Recruiter notifications |
| Connection requests | connection-events | Graph service | PYMK invalidation, feed |
| Metrics | service-metrics | Every service | Pinot (real-time OLAP) |
| Notification triggers | notification-candidates | Samza jobs | Email/push service |

**Kafka Cruise Control** is LinkedIn's open-source cluster balancing tool that automatically rebalances partition assignments across brokers to equalize disk usage and network traffic.

## Samza: LinkedIn's Stream Processing Framework

**Apache Samza** was built at LinkedIn before Flink and Spark Structured Streaming matured. Its key design principle: **stateful stream processing with local state**.

Instead of storing state in a remote database, Samza keeps state in a **local RocksDB instance** on the same machine as the processing task. State changes are backed by a changelog Kafka topic, so state can be rebuilt by replaying the changelog.

This gives:
- Sub-millisecond state lookups (no network hop)
- Durable state (RocksDB backed by Kafka changelog)
- Horizontal scalability (each partition = one stateful task)

**Samza vs. Flink vs. Spark Streaming:**

| Feature | Samza | Flink | Spark Structured Streaming |
|---------|-------|-------|---------------------------|
| Origin | LinkedIn, 2013 | Berlin, 2014 | Databricks, 2016 |
| State backend | Local RocksDB + Kafka changelog | RocksDB or heap + checkpoints | Spark RDDs |
| Latency | Low (record-at-a-time) | Low (record-at-a-time) | Micro-batch (100ms+) |
| Deployment | YARN or Kubernetes | Kubernetes | Kubernetes / Spark |
| Adoption at LinkedIn | Core platform | Growing | Limited |

## Brooklin: Change Data Capture at LinkedIn

**Brooklin** is LinkedIn's data movement platform. Its primary use case: streaming changes from Espresso (LinkedIn's NoSQL store) to Kafka topics, which are then consumed by Elasticsearch for index updates, Venice for feature propagation, and the data lake.

Without Brooklin, every service that needed to react to database changes would need to poll the database or maintain a manual dual-write. Brooklin provides CDC as a platform primitive.

## Consumer Lag: The Critical Operational Metric

Consumer lag is the most important Kafka operational metric. High lag means:
- **Notifications are delayed** — if the notification consumer is lagging, members receive job alerts late.
- **Feed ranking is stale** — if the feed candidate consumer is lagging, new content doesn't appear in feeds.
- **Data pipeline SLA breaches** — if the analytics consumer is lagging, dashboards show stale data.

**Backpressure:** when consumers process records slower than producers publish them, lag grows. The correct response is to scale out consumers (add more instances) or optimize the consumer processing logic, not to throttle the producer.

**Lag alarm:** alert when `total_lag > threshold` for a sustained period (e.g., lag > 100K records for >5 minutes).

## Failure Modes

| Mode | Cause | Mitigation |
|------|-------|------------|
| Broker failure | Machine crash, disk failure | ISR replication + `min.insync.replicas=2`; automatic leader election |
| Consumer group rebalance storm | Many consumers join/leave simultaneously | Incremental cooperative rebalancing (Kafka 2.4+); sticky partition assignment |
| Disk full on broker | Topic retention not tuned; unexpected traffic spike | Capacity planning alarms at 75% disk; automatic size-based retention enforcement |
| Schema registry unavailable | If using Avro/Protobuf, schema registry failure blocks producers/consumers | Schema registry HA (multiple instances); cached schema local to producer/consumer |
| Producer timeout under load | Producer buffer full; broker slow | Tune `buffer.memory`, `batch.size`, `linger.ms`; monitor producer metrics |
| Log compaction lag | Compaction falls behind for changelog topics | Monitor `kafka.log.Log:type=LogFlushStats`; tune `log.cleaner.min.cleanable.ratio` |

## Interview Trade-offs to Discuss

- **Partitioning by member_id vs. random:** Partitioning by member_id ensures all events for a member are processed in order by the same consumer (needed for PYMK computation). Random partitioning gives better load distribution but loses per-member ordering.
- **Consumer group isolation:** Using separate consumer groups for feed ranking and PYMK recompute means each system can scale and fail independently. A single shared consumer group would couple their throughput.
- **Retention policy:** 7-day retention is enough for most event-driven pipelines. Longer retention enables event replay for debugging and reprocessing after consumer bugs. Cost grows linearly with retention.
- **Exactly-once vs. at-least-once:** Kafka supports exactly-once semantics (EOS) via idempotent producers + transactional APIs. At-most-once risks data loss. At-least-once is the most common production choice — consumers handle idempotency via deduplication keys.
- **Samza vs. Flink at LinkedIn today:** Flink has taken over new stateful streaming workloads at LinkedIn due to better community support and richer SQL interface. Samza remains critical for existing pipelines.
