# Distributed Message Queue

> A queue is not a magical buffer; it is an explicit contract about ordering, retry, ownership, and backlog pain.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a distributed message queue that makes partitioning, delivery semantics, backlog recovery, and consumer isolation explicit instead of hand-waving them as broker internals.
**Prerequisites:** `07-queues-streams-and-workflows/02-delivery-semantics`, `09-partitioning-sharding-and-rebalancing/01-shard-key`, `10-reliability-retries-and-backpressure/05-async-backpressure`
**Estimated time:** ~90 min
**Primary artifact:** queue-plan validator + design review sheet

## The Problem

Design a distributed message queue used by internal services for asynchronous jobs, fan-in ingestion, and decoupled workflows. Producers publish messages durably, consumers pull or receive deliveries, and operators need replay, redelivery, and backlog visibility without pretending the system gives universal exactly-once delivery.

This lesson matters because interview answers often stop at "use Kafka" or "use SQS." Senior answers explain what is ordered, who owns offsets or acknowledgements, how poison messages are isolated, and how slow consumers stop being a platform-wide incident.

## Clarify

- Do consumers need per-key ordering, total ordering, or only eventual processing?
- Is the product closer to a task queue, an event log, or both with different access patterns?
- Who owns retry policy: the broker, the consumer, or a higher-level workflow engine?
- What replay guarantees and retention windows are required?

If the interviewer stays broad, assume durable at-least-once delivery, per-partition ordering, consumer-group style scaling, and retention long enough for replay and operational debugging.

## Requirements

### Functional

- Accept producer writes durably before acknowledging success.
- Partition messages for scale while preserving ordering within a partition.
- Allow consumers to track progress and recover after crashes.
- Support retry, visibility timeout or lease semantics, and poison-message isolation.
- Expose replay or seek controls within retention windows.

### Non-functional

- Scale producer throughput independently from consumer speed.
- Prevent one hot topic or slow consumer group from destabilizing the cluster.
- Make backlog age and delivery lag observable.
- Recover cleanly from broker, partition-leader, or consumer failure.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Producer writes | 15M messages/s peak | drives partition count, batching, and replication fanout |
| Average message size | 4 KB | shapes disk, network, and page-cache pressure |
| Retention window | 72 hours hot, 14 days cold | determines storage and replay posture |
| Consumer groups | 8K active | affects offset storage and fanout overhead |
| Peak factor | 5x during incident replays | replay safety matters as much as steady-state ingest |

## Architecture

```text
producers
  -> broker ingress
  -> partition router
  -> replicated append log
  -> retention / compaction policy

consumer groups
  -> partition assignment
  -> delivery lease or fetch cursor
  -> ack / offset commit
  -> retry and poison-message path
```

Design notes:

1. Partition by key when ordering matters, and say what happens when no stable key exists.
2. Separate durable append from delivery acknowledgement so consumer retries stay explicit.
3. Keep replay a control-plane feature with quotas because it can look like a self-inflicted DDoS.
4. Treat backlog age as a first-class signal; queue depth alone hides whether the system is catching up.

## Data Model & APIs

Core records:

```text
topic
partition_id
message_id
ordering_key
enqueue_time
delivery_attempt
visibility_deadline
consumer_group
committed_offset
```

Useful interfaces:

- `CreateTopic(name, partitions, retention_policy)`
- `Publish(topic, ordering_key, payload, headers)`
- `Fetch(topic, consumer_group, max_batch, lease_seconds)`
- `Ack(topic, consumer_group, partition, offset)`
- `Seek(topic, consumer_group, partition, offset)`

Strong answers explicitly separate producer durability, consumer delivery, and replay administration.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hot partition receives most writes | per-partition throughput and queue-age skew | better keying, partition expansion, or key bucketing |
| consumer crashes after delivery | expired lease count and redelivery attempts | visibility timeout plus idempotent consumers |
| poison message blocks progress | repeated attempt count on same offset | DLQ or side queue after bounded retries |
| replay overwhelms live traffic | replay job rate and live-lag regression | replay quotas, isolated lanes, and admission control |

## Observability

- metric: publish latency and replication-ack latency by topic tier
- metric: backlog age percentile by topic and consumer group
- metric: redelivery count, expired leases, and poison-message rate
- metric: partition skew for bytes, writes, and lag
- log: replay, seek, and retention-policy changes with actor identity
- trace: publish to durable append to consumer ack for sampled messages
- SLO: 99% of standard-priority messages become visible to healthy consumers within the queue target under normal load

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| per-partition ordering | simple reasoning for keyed workloads | hot-key and rebalancing pain | global total ordering |
| at-least-once delivery | practical durability and recovery | duplicate processing risk | pretending exactly-once is free |
| replay as control-plane action | safer operational behavior | slower emergency catch-up | unlimited self-serve replay |

## Interview It

**Google framing:** "Design a queue for asynchronous internal jobs used by many services." Expect follow-ups on partitioning, offset ownership, replay safety, and consumer lag diagnosis.

**Cloudflare framing:** "Design an event backbone for edge-generated async workloads." Expect pressure on tenant isolation, regional durability, and protecting the shared platform from replay storms.

**Follow-ups:**
1. What changes if consumers need FIFO within tenant but not across tenants?
2. How do you expand partition count without breaking ordering assumptions?
3. When would you choose push delivery over pull delivery?
4. How would you isolate a subscription that is always behind?
5. What if compliance requires immutable retention for selected topics?

## Ship It

- `outputs/design-review-distributed-message-queue.md`

## Exercises

1. **Easy** — Pick a partition key for email-send tasks and explain its ordering consequence.
2. **Medium** — Redesign the queue for a workload with large replay bursts after downstream outages.
3. **Hard** — Compare queue behavior when one consumer group needs strict ordering but another only needs fast fanout.

## Further Reading

- [Kafka design](https://kafka.apache.org/documentation/#design) — helpful grounding for log-based messaging trade-offs
- [Amazon SQS visibility timeout](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html) — useful for delivery-lease thinking
