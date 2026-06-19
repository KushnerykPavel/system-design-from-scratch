# Event-Driven Architecture on AWS

> Decouple services with events — each service reacts to what happened, not what to do.

**Type:** Build  
**Company focus:** Amazon  
**Learning goal:** Design an event-driven system using SQS, SNS, EventBridge, Kinesis, and Step Functions for the right use case.  
**Prerequisites:** `07-queues-streams-and-workflows/01-queues-vs-streams`, `10-reliability-retries-and-backpressure/03-circuit-breakers`  
**Estimated time:** ~75 min  
**Primary artifact:** AWS messaging service decision matrix  

## The Problem

Amazon's e-commerce platform is a collection of independently deployable services that must react to each other without tight coupling. An order placement triggers inventory reservation, payment processing, fulfillment routing, notification delivery, analytics ingestion, and fraud scoring — all independently, all potentially in parallel, all with different reliability requirements. The wrong choice of messaging service creates cascading failures, ordering violations, duplicate processing, or runaway costs. The right choice is driven by a clear decision matrix built from use-case properties.

## Clarify

- Does message ordering matter? FIFO guarantees come at a throughput cost.
- Does more than one consumer need the same event? Fan-out changes the choice from queue to topic or event bus.
- How long must messages survive if consumers are offline? Retention requirements differ across services.
- Is this a high-throughput stream or a low-frequency command? Volume and consumer count shape the choice.
- Assumption if no answer: at-least-once delivery, eventual consistency acceptable, single consumer per message type unless told otherwise.

## Requirements

### Functional
1. Decouple order placement from downstream services so that each processes events independently.
2. Fan-out a single `order.placed` event to up to 10 downstream consumers without the producer knowing about them.
3. Replay events for a new analytics consumer onboarding without re-processing old orders in the fulfillment service.
4. Orchestrate a multi-step workflow (payment → reserve inventory → assign fulfillment → notify) with per-step retry and a durable audit trail.

### Non-functional
1. Order event delivery: at-least-once, processed within 30 seconds under normal load.
2. No message loss for payment events, even if the payment service is temporarily unavailable.
3. Analytics pipeline must be able to consume events from 7 days ago for backfill.
4. Step Functions workflow must complete the full order pipeline in < 5 minutes or escalate to ops.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Order events/day | ~35M | baseline SQS/Kinesis throughput sizing |
| Peak order events/min | ~100,000 | drives sharding strategy for FIFO queues |
| Consumers per event type | up to 10 (fan-out) | drives SNS vs EventBridge selection |
| Analytics retention needed | 7 days | drives Kinesis vs SQS for replay |
| Step Functions workflow steps | ~7 per order | drives cost: $0.025 per 1000 state transitions (standard) |
| DLQ inspection cadence | < 1 hour SLA | drives alarm configuration |

## Architecture

### SQS — Simple Queue Service

**What it is:** Point-to-point message queue. One message, one consumer group. At-least-once delivery. Messages persist until consumed or the retention period (default 4 days, max 14 days) expires.

**Visibility timeout:** When a consumer reads a message, SQS hides it from other consumers for the visibility timeout period. If the consumer does not delete the message before the timeout expires, it becomes visible again and is redelivered. Set visibility timeout to at least 6× the expected processing time to avoid spurious redelivery.

**Dead Letter Queue (DLQ):** After `maxReceiveCount` delivery attempts, SQS moves the message to a DLQ. DLQs are the first line of defense against poison pill messages — malformed events that cause the consumer to throw an exception on every attempt, blocking the queue indefinitely. Every production queue must have a DLQ with a CloudWatch alarm on `ApproximateNumberOfMessagesVisible > 0`.

**FIFO queues:** Guarantee strict ordering within a message group and exactly-once processing (deduplication window: 5 minutes). Hard ceiling: 300 TPS per queue, or 3,000 TPS with high-throughput mode. At Prime Day's 100K orders/min (~1,667 TPS), you need at least 6 FIFO queues sharded by order ID prefix.

**When to use SQS:**
- One consumer group per message type.
- Commands (do something), not events (something happened).
- No fan-out needed.
- Backpressure via queue depth: consumer scales based on `ApproximateNumberOfMessagesNotVisible`.

### SNS — Simple Notification Service

**What it is:** Pub/sub fan-out. One producer publishes to a topic; SNS pushes the message to all subscribers simultaneously. Subscribers can be SQS queues, Lambda functions, HTTP endpoints, email, or SMS. No message persistence — SNS does not store messages. If a subscriber is unavailable at delivery time, the message is lost (unless the subscriber is an SQS queue, which buffers it).

**SNS → SQS fan-out pattern:** The canonical pattern for durable fan-out. Producer → SNS topic → N SQS queues (one per consumer service). Each consumer reads from its own queue at its own pace. This is how Amazon broadcasts `order.placed` to fulfillment, notifications, fraud, and analytics simultaneously with backpressure isolation.

**No content-based filtering:** SNS delivers all messages to all subscribers. Consumers must discard irrelevant events. EventBridge solves this.

**When to use SNS:**
- Fan-out to multiple consumers of the same event.
- No replay needed.
- No content-based filtering needed.
- Push delivery to Lambda or HTTP required.

### EventBridge — Event Bus

**What it is:** Serverless event bus with content-based filtering rules, schema registry, and cross-account delivery. Producers send events to a bus; EventBridge evaluates rules and routes matching events to targets.

**Content-based filtering:** A rule like `{ "source": ["order-service"], "detail-type": ["order.placed"], "detail": { "status": ["CONFIRMED"] } }` delivers only confirmed orders to the fulfillment service — without the fulfillment service receiving and discarding PENDING or CANCELLED events. This decouples consumers from irrelevant event types.

**Schema registry:** EventBridge can discover event schemas automatically and generate code bindings. Reduces contract drift between producer and consumer.

**Cross-account delivery:** EventBridge can route events to buses in other AWS accounts, enabling platform-level event sharing without VPC peering or custom API layers.

**When to use EventBridge:**
- Multiple consumers, each interested in a subset of event types.
- Cross-account event distribution.
- Schema evolution and discovery matter.
- Low-to-medium throughput (soft limit: 10,000 events/s per bus by default).

### Kinesis Data Streams

**What it is:** Ordered, sharded, persistent stream. Records within a shard are strictly ordered by sequence number. Multiple independent consumer groups can read the same stream simultaneously (enhanced fan-out gives each consumer a dedicated 2 MB/s throughput pipe). Default retention: 24 hours; extendable to 365 days.

**Shard throughput:** 1 MB/s writes or 1,000 records/s per shard. Scale by adding shards. At 100K orders/min from Prime Day, you need at least 2 shards for writes plus additional shards for each consumer.

**Consumer groups via enhanced fan-out:** Unlike SQS where consuming destroys the message, Kinesis lets an analytics service and a fraud service each read the full stream independently at their own sequence positions.

**Kinesis Firehose:** Managed delivery of Kinesis records to S3, Redshift, or OpenSearch. No consumer group control — Firehose reads from the stream and batches to destination. Use Firehose for one-way ETL pipelines, not for application-level consumers that need to track position.

**When to use Kinesis:**
- Ordered events required.
- Multiple independent consumers of the same stream.
- Replay or long retention needed (hours to days).
- High-throughput ingestion (clickstream, order events at scale).

### Step Functions — Workflow Orchestration

**What it is:** Serverless state machine orchestrator. Each step in the workflow is an activity (Lambda invocation, ECS task, DynamoDB write, SQS send, human approval). State and retry logic are stored durably by Step Functions — not in your application code.

**Standard vs Express workflows:**
| Property | Standard | Express |
|----------|----------|---------|
| Max duration | 1 year | 5 minutes |
| Execution semantics | Exactly-once per state | At-least-once |
| Audit trail | Full history in console | CloudWatch Logs only |
| Cost | $0.025 / 1,000 transitions | $1.00 / 1M transitions |
| Use case | Order pipeline, human approvals | High-volume, short-lived |

**When to use Step Functions:**
- Multi-step workflow where failure of step N must not re-execute steps 1..N-1.
- Human approval gates (order fraud review, seller verification).
- Long-running workflows where a single Lambda timeout (15 min) is insufficient.
- Audit trail required for compliance.

### Decision Matrix

| Requirement | Service |
|-------------|---------|
| Single consumer, reliable queue, backpressure | SQS Standard |
| Single consumer, ordered, deduplication | SQS FIFO |
| Fan-out to N consumers, push, no persistence | SNS |
| Fan-out with content-based routing, schema | EventBridge |
| Ordered stream, multiple consumer groups, replay | Kinesis Data Streams |
| ETL to S3/Redshift, no consumer control | Kinesis Firehose |
| Multi-step durable workflow, audit trail | Step Functions Standard |
| High-volume short-lived orchestration | Step Functions Express |

### Lambda Integration Considerations

Lambda integrates with SQS, Kinesis, DynamoDB Streams, and SNS as event sources. Key failure modes:
- **SQS poison pill:** A malformed message causes Lambda to throw on every invocation. SQS retries until `maxReceiveCount`, then moves to DLQ. Without a DLQ, a poison pill blocks the queue and stops all processing.
- **Lambda concurrency exhaustion:** If SQS delivers 10,000 messages simultaneously and Lambda hits its concurrency limit, excess invocations are throttled. SQS messages remain in the queue and are retried after the visibility timeout — but the queue depth grows.
- **Cold starts:** Lambda containers spin up on first invocation after idle period (~100–500 ms for Go). For latency-sensitive paths, use Provisioned Concurrency or ensure event volume keeps containers warm.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| SQS poison pill message | DLQ `ApproximateNumberOfMessagesVisible` > 0 | DLQ on every queue; CloudWatch alarm; dead-letter replay tooling |
| Lambda concurrency exhaustion | `ConcurrentExecutions` at account limit; SQS queue depth growing | Increase concurrency limit; add concurrency reservation per function; use SQS batch size tuning |
| Step Functions timeout on long workflow | Execution status `TIMED_OUT` | Set `HeartbeatSeconds` on long activities; use .waitForTaskToken for async steps |
| SNS subscriber unavailable | Message lost (SNS has no persistence) | Always pair SNS with SQS subscriber for durability; do not subscribe Lambda directly for critical events |
| Kinesis shard hot spot | One shard at throughput limit while others are idle | Use a partition key with high cardinality; use `EnhancedMonitoring` to detect per-shard metrics |

## Observability

- metric: SQS `ApproximateAgeOfOldestMessage` — alert if > 5 min (consumer falling behind)
- metric: SQS DLQ `ApproximateNumberOfMessagesVisible` — alert if > 0
- metric: Kinesis `GetRecords.IteratorAgeMilliseconds` — alert if > 60,000 ms (1 minute behind)
- metric: Step Functions `ExecutionsFailed` and `ExecutionsTimedOut` per state machine
- metric: Lambda `ConcurrentExecutions` approaching account limit
- trace: X-Ray traces across Lambda → SQS → Lambda chains using correlation IDs propagated in message attributes
- log: DLQ message contents with original error context for triage

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| SQS FIFO over Kinesis for per-order ordering | Simpler consumer; built-in deduplication; no shard management | 300 TPS ceiling requires sharding by order prefix at Prime Day scale | Kinesis gives ordering but requires consumer group management and shard scaling |
| EventBridge over SNS for cross-service events | Content-based filtering; schema registry; cross-account | Lower throughput ceiling; slightly higher per-event cost | SNS is simpler but delivers all events to all subscribers, wasting consumer compute |
| Step Functions Standard over Lambda chain | Durable state; per-step retry; audit trail; no lost-progress on failure | $0.025/1000 transitions; adds latency between steps | Lambda chain loses progress on failure; no native retry per step; no audit trail |
| SNS → SQS fan-out over direct SNS → Lambda | SQS buffers events if consumer Lambda is unavailable; decouples consumer scaling | Two hops instead of one; slightly higher cost | Direct SNS → Lambda loses events if Lambda is throttled at delivery time |

## Interview It

**Amazon framing:** "Design the event-driven order pipeline for Amazon's checkout service." LP tie-in: Ownership (you own end-to-end reliability, including DLQ), Frugality (Step Functions Standard vs Express cost comparison), Bias for Action (choose SNS→SQS fan-out as MVP, evolve to EventBridge when filtering is needed).

**Follow-ups:**
1. A bug is deployed that causes the inventory service's Lambda to throw an exception on every message. Walk me through what happens in your SQS-based design and how you recover.
2. You need to add a new analytics consumer that needs to read the last 3 days of order events. What service do you use and why?
3. Why is SNS not sufficient for durable fan-out to the fulfillment service?
4. The Step Functions order workflow is timing out at the "assign fulfillment center" step during a regional outage. How do you handle this without losing the order?
5. What is the difference between EventBridge content-based filtering and SNS filter policies?

## Ship It

- `outputs/aws-messaging-decision-matrix.md`

## Exercises

1. **Easy** — Draw the SNS → SQS fan-out pattern for `order.placed` going to four downstream consumers. Label each component and its delivery guarantee.
2. **Medium** — Design the SQS sharding strategy for FIFO queues at 100K orders/min. How many queues? What is the sharding key? How does a consumer find the right queue for a given order ID?
3. **Hard** — Write the Step Functions state machine definition (as a state chart, not code) for the order fulfillment pipeline: payment → inventory reservation → fulfillment assignment → notification → done. Include catch/retry clauses and a DLQ fallback for irrecoverable failures.

## Further Reading

- [Amazon SQS Developer Guide](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/welcome.html)
- [Amazon EventBridge User Guide](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-what-is.html)
- [AWS Step Functions Developer Guide](https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html)
- [Kinesis Data Streams Developer Guide](https://docs.aws.amazon.com/streams/latest/dev/introduction.html)
- [The Amazon Builder's Library: Avoiding fallback in distributed systems](https://aws.amazon.com/builders-library/avoiding-fallback-in-distributed-systems/)
