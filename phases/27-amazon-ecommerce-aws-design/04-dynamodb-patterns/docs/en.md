# DynamoDB Internals & NoSQL Access Patterns

> Model your access patterns first; the table schema follows — the opposite of relational thinking.

**Type:** Concept  
**Company focus:** Amazon  
**Learning goal:** Master DynamoDB single-table design, partition key selection, and when to use GSIs vs LSIs vs query patterns.  
**Prerequisites:** `05-storage-indexing-and-access-patterns/01-storage-models`, `09-partitioning-sharding-and-rebalancing/01-shard-key`  
**Estimated time:** ~75 min  
**Primary artifact:** DynamoDB access pattern worksheet  

## The Problem

DynamoDB is the backbone of Amazon's e-commerce platform — order records, inventory counts, session data, and product catalog attributes all live in DynamoDB tables. Engineers who treat it as "a NoSQL database where you just write JSON" produce systems that throttle under load, incur unexpected costs, and fail at Prime Day scale. The discipline is fundamentally different from relational modeling: you start with a list of access patterns, then choose keys that make every pattern a single O(1) lookup or a bounded O(k) range scan with no table scans, no multi-table joins.

## Clarify

- What are the concrete access patterns? List them before touching schema.
- What is the expected read/write ratio per pattern? Hot reads drive caching; hot writes drive partition key selection.
- Does ordering matter? LSIs and sort keys give range queries within a partition; GSIs give secondary projections across partitions.
- Is eventual consistency acceptable for reads? Strong consistency is available but doubles read cost.
- Assumption if no answer: high-read/low-write ratio, eventual consistency acceptable for catalog reads, strong consistency required for inventory writes.

## Requirements

### Functional
1. Retrieve a single entity by primary key in < 5 ms at p99.
2. List all related entities for an owner (e.g., all orders for a customer) in a single query.
3. Support secondary lookups (e.g., find all orders by status) without a full table scan.
4. Store multiple entity types in the same table using a single-table design pattern.

### Non-functional
1. Sustain Prime Day peak of 100K writes/second on inventory updates without hot-partition throttling.
2. GSI replication lag < 1 s acceptable for non-inventory reads.
3. Cost model on provisioned capacity with auto-scaling, switching to on-demand during unpredictable spikes.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Items per table | 350M product listings | drives partition count and key cardinality requirements |
| Peak read TPS | ~500K RCU/s on Prime Day | drives provisioned capacity and DAX caching decision |
| Peak write TPS | ~100K WCU/s on inventory | drives partition key uniformity requirements |
| Average item size | ~2 KB catalog item | drives RCU/WCU calculation (1 RCU = 4 KB strongly consistent) |
| GSIs per table | up to 20 | drives secondary access pattern budget |
| DynamoDB streams retention | 24 hours | drives CDC consumer lag tolerance |

## Architecture

### Dynamo Paper Origins

DynamoDB is Amazon's internal Dynamo paper (2007) productized. The paper introduced consistent hashing for partition placement, vector clocks for conflict resolution, and quorum-based reads/writes for tunable consistency. The production service adds a B-tree-based storage engine (not LSM) per shard, a replicated write-ahead log across three AZs, and a global control plane for partition rebalancing.

### Partition Key Design

The partition key (PK) determines which physical partition stores the item. DynamoDB hashes the PK uniformly across partitions — but only if the PK has high cardinality and uniform access frequency.

**Hot partition problem:** If 80% of your reads target `PK=PRODUCT#B07XYZ` (a viral listing), that one partition absorbs disproportionate traffic and throttles. The fix:
- **Shard the hot key** — append a random suffix `PRODUCT#B07XYZ#3` and read all shards then merge. Effective but adds application complexity.
- **Prefix with entity type** — `USER#<uuid>` and `ORDER#<uuid>` ensure no cross-entity key collision.
- **Avoid sequential keys** — `ORDER#00001`, `ORDER#00002` creates monotonically increasing keys that land on a single partition. Use UUID or hash-based IDs.

**Uniform distribution check:** Compute the standard deviation of request counts across partition keys. Divide by the mean. If `stddev / mean > 0.1`, investigate. DistributionScore formula: `score = 100 * max(0, 1 - stddev/mean)`, clamped to [0, 100]. A score of 100 means perfectly uniform; a score below 70 indicates a hot-partition risk.

### Single-Table Design (STD)

Store all entity types in one table using an overloaded PK/SK (partition key / sort key) pattern:

```
PK                    SK                      Attributes
USER#u-001            PROFILE                 name, email, tier
USER#u-001            ORDER#o-1001            status, total, created_at
USER#u-001            ORDER#o-1002            status, total, created_at
ORDER#o-1001          ITEM#p-501              qty, price
ORDER#o-1001          ITEM#p-502              qty, price
PRODUCT#p-501         METADATA                title, brand, category
```

This design supports:
- `GetItem(PK=USER#u-001, SK=PROFILE)` — fetch user profile.
- `Query(PK=USER#u-001, SK begins_with ORDER#)` — all orders for a user.
- `Query(PK=ORDER#o-1001, SK begins_with ITEM#)` — all items in an order.

**Entity adjacency list** is the pattern where relationships are stored as items alongside the entities they link. A `PRODUCT#p-501 / ORDER#o-1001` item captures the fact that order 1001 contains product 501, enabling bidirectional traversal without a join table.

**Trade-off:** Single-table design is harder to reason about for teams unfamiliar with it; it couples all entity lifecycle to one table; backups restore all entities together. Multi-table design is simpler to understand, supports independent scaling, and cleaner access controls — but requires application-side joins and more operational overhead.

### GSI vs LSI

| Property | LSI (Local Secondary Index) | GSI (Global Secondary Index) |
|----------|-----------------------------|-----------------------------|
| Scope | Same partition as base table | Entire table (different partition key allowed) |
| Consistency | Strongly consistent reads possible | Eventual consistency only |
| Creation | Must be defined at table creation | Can be added any time |
| Throughput | Shares table capacity | Independent capacity provisioned |
| Limit | 5 per table | 20 per table |
| Use case | Range queries within an owner partition | Cross-partition secondary lookups |

**GSI replication lag:** GSI writes are asynchronous. A write to the base table is eventually propagated to the GSI, typically within milliseconds but potentially seconds under high load. Never use a GSI for inventory counts that must be strongly consistent.

### Capacity Modes

**Provisioned:** You specify read capacity units (RCUs) and write capacity units (WCUs). Auto-scaling adjusts within bounds. Cheaper for predictable workloads. Risk: provisioning too low causes throttling; too high wastes money.

**On-demand:** No capacity planning required; you pay per request. 2x-3x more expensive at sustained load. Correct choice for unpredictable spikes (new product launch, marketing flash sales).

**Burst capacity:** DynamoDB retains unused capacity for up to 300 seconds. A partition can burst to 3,000 RCU/s or 1,000 WCU/s briefly. This masks small spikes but not sustained hot partitions.

**Adaptive capacity:** Automatically shifts capacity to hot partitions. Operates on a minutes timescale — does not help with sub-minute flash sale spikes. Do not rely on it instead of good key design.

### DAX (DynamoDB Accelerator)

In-memory write-through cache sitting in front of DynamoDB. Reads served from DAX at microsecond latency; cache misses fall through to DynamoDB. Write-through means DAX and DynamoDB are updated atomically on writes. Use for read-heavy catalog data where eventual consistency is acceptable. Do not use for inventory reservation — DAX may serve a stale count and allow oversell.

### DynamoDB Streams for CDC

Streams capture item-level changes (INSERT, MODIFY, REMOVE) in order within a shard, with 24-hour retention. Lambda functions subscribe to streams for change-data-capture (CDC) use cases: rebuilding search indexes, updating materialized views, and feeding Kinesis for analytics. Streams are a DynamoDB-native alternative to Debezium/Kafka CDC used in relational systems.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Hot partition throttling on viral product | `ThrottledRequests` spike on table CloudWatch metric | Write sharding with random suffix; ElastiCache caching layer for reads |
| GSI replication lag on cross-entity query | Stale results from GSI vs base table discrepancy | Use base table query when strong consistency required; add lag monitoring via stream position |
| Capacity spike on flash sale | Provisioned WCU exhausted within seconds | Switch to on-demand mode before event; implement exponential backoff with jitter in client |
| Item size exceeds 400 KB limit | `ValidationException` on large documents | Store large blobs in S3; store S3 key in DynamoDB item |
| Conditional write collision under high contention | `ConditionalCheckFailedException` bursts | Implement optimistic locking retry loop with bounded attempts; alert on high collision rate |

## Observability

- metric: `ThrottledRequests` per table and per index — alert if non-zero for > 60 s
- metric: `SuccessfulRequestLatency` at p99 per operation type (GetItem, Query, PutItem)
- metric: GSI consumed capacity vs provisioned ratio
- metric: DynamoDB Streams `IteratorAgeMilliseconds` — alert if > 5000 ms (consumer falling behind)
- log: conditional write failures with partition key and retry count
- trace: DynamoDB calls in X-Ray service map, correlated with upstream Lambda or ECS service

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Single-table design | All access patterns served by one round-trip; consistent performance | Schema harder to grok; teams must learn STD before contributing | Multi-table is simpler to read but requires application-side joins at query time |
| GSI for secondary lookups | Decouples secondary access pattern from base table partition key | Eventual consistency; extra write cost; replication lag | Application-side scan is O(N) and unacceptable at 350M items |
| On-demand capacity for flash sales | No capacity planning; absorbs traffic spike instantly | 2-3x cost vs provisioned at steady state | Pre-scaling provisioned capacity requires accurate forecast and misses late-breaking spikes |
| Partition key sharding for hot items | Eliminates single-partition hot spot | Scatter-gather reads across all shards; application complexity | Adaptive capacity is too slow to react to sub-minute spikes |

## Interview It

**Amazon framing:** "Design DynamoDB's access pattern for Amazon's order management system." LP tie-in: Dive Deep (explain single-table design trade-offs), Frugality (provisioned vs on-demand cost modeling), Ownership (you choose the key design and own any throttling that results).

**Follow-ups:**
1. How would you model the access pattern for "get all orders for a customer sorted by recency" using single-table design?
2. A product goes viral during Prime Day. Your GSI on `category_id` is throttling. What do you do right now vs what do you change for next time?
3. When would you reject single-table design and use multiple tables instead?
4. How does DynamoDB Streams differ from Kinesis, and when do you use each?
5. What is adaptive capacity, and why can't you rely on it instead of good partition key design?

## Ship It

- `outputs/dynamodb-access-pattern-worksheet.md`

## Exercises

1. **Easy** — List three bad partition key choices for an e-commerce order table and explain why each creates a hot partition.
2. **Medium** — Design the PK/SK schema for a single-table design that supports: fetch order by ID, list orders by customer, list items in an order, and look up all orders in SHIPPED status.
3. **Hard** — Calculate the RCU cost of fetching all 50 items in a large order using `Query` with eventual consistency vs 50 individual `GetItem` calls with strong consistency.

## Further Reading

- [DynamoDB Best Practices — AWS Docs](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html)
- [Amazon DynamoDB Design Patterns (re:Invent 2019)](https://www.youtube.com/watch?v=HaEPXoXVf2k) — Rick Houlihan's canonical single-table design talk
- [The DynamoDB Book — Alex DeBrie](https://www.dynamodbbook.com/) — comprehensive single-table design reference
- [Dynamo: Amazon's Highly Available Key-Value Store (paper)](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) — original Dynamo paper
