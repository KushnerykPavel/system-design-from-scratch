# Multi-Region Reliability on AWS

> Design for failure of an entire AWS region — it happens, and recovery must be automatic.

**Type:** Build  
**Company focus:** Amazon  
**Learning goal:** Design an active-active multi-region architecture on AWS with Route53, Global Accelerator, DynamoDB Global Tables, and Aurora Global Database.  
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive`, `10-reliability-retries-and-backpressure/07-bulkheads`  
**Estimated time:** ~75 min  
**Primary artifact:** multi-region architecture checklist  

## The Problem

AWS regions fail. In 2011, 2012, 2017, and 2021, AWS regions experienced multi-hour outages affecting significant portions of their services. Amazon's own e-commerce platform cannot afford to be unavailable for hours — each minute of downtime during peak costs millions of dollars and erodes customer trust. The solution is multi-region architecture, but it introduces hard problems: data consistency across regions, conflict resolution on concurrent writes, and traffic routing during partial failures. Getting these wrong produces split-brain, stale reads, or recovery procedures that are slower than a single-region rebuild.

## Clarify

- Is this active-active (traffic served from multiple regions simultaneously) or active-passive (traffic served from one, failover to another)?
- What are the RTO (Recovery Time Objective) and RPO (Recovery Point Objective) requirements? These determine the DR tier.
- Which services require strong consistency across regions vs eventual consistency?
- Does the application use stateful sessions? Session affinity complicates active-active routing.
- Assumption if no answer: Amazon-scale e-commerce, active-active for reads, near-zero RTO for customer-facing paths, < 1 minute RPO for order data.

## Requirements

### Functional
1. Serve traffic from at least two regions simultaneously to minimize latency for geographically distributed customers.
2. Automatically failover DNS and traffic routing within 60 seconds of a region health check failure.
3. Replicate order and inventory data across regions with < 1 second lag under normal conditions.
4. Support global secondary reads from any region with eventual consistency.

### Non-functional
1. RTO < 1 minute for customer-facing paths (checkout, order lookup).
2. RPO < 1 minute for order data (no more than 1 minute of writes lost on region failure).
3. 99.99% availability target, requiring < 52 minutes of downtime per year.
4. Stale reads acceptable for product catalog; not acceptable for inventory reservation.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active regions | 2–3 (US-East-1, EU-West-1, AP-Northeast-1) | drives replication topology design |
| DynamoDB Global Tables replication lag | ~1 second typical | sets RPO floor for DynamoDB-backed data |
| Aurora Global Database replication lag | < 1 second typical | sets RPO floor for relational data |
| Route53 health check interval | 10 s (fast health checks) | sets failover detection floor |
| Global Accelerator failover time | ~30–60 s | sets practical RTO for anycast routing |
| S3 CRR (Cross-Region Replication) lag | seconds to minutes | sets RPO floor for object storage |

## Architecture

### AWS Regions and Availability Zones

Each AWS region contains at least three Availability Zones (AZs). AZs are physically separated data centers within 60–100 km of each other, connected by redundant low-latency links. A region-level failure affects all AZs in that region simultaneously — AZ redundancy alone does not protect against regional events such as large-scale power failures, network backbone issues, or software control-plane outages.

### Active-Active vs Active-Passive vs Hot-Warm

| Strategy | Description | RTO | RPO | Cost |
|----------|-------------|-----|-----|------|
| Active-active | All regions serve live traffic; any region can take full load | near-zero | near-zero | very high |
| Active-passive | Primary region serves all traffic; passive region is a warm standby | 1–5 min | < 1 min | high |
| Hot-warm | Primary region live; warm region pre-scaled but not serving traffic | 5–15 min | minutes | medium |
| Pilot light | Minimal infrastructure running in DR region; scale up on failover | 30–60 min | minutes | low |
| Backup/restore | No running DR infrastructure; restore from S3 backup | hours | hours | very low |

Amazon's customer-facing paths (checkout, order lookup) use **active-active** across regions. Internal analytical pipelines use **warm standby** or **pilot light**.

### Route53: DNS-Based Traffic Routing

Route53 supports multiple routing policies that enable multi-region traffic management:

**Latency-based routing:** Routes each request to the region with the lowest measured latency for the user's location. A customer in Tokyo gets routed to `ap-northeast-1`; a customer in New York gets `us-east-1`. This is the default choice for active-active global deployments.

**Failover routing:** Primary record + secondary record. Route53 monitors the primary endpoint with a health check. On failure, traffic shifts to the secondary automatically. Used for active-passive pairs.

**Health checks:** Route53 probes an endpoint every 10–30 seconds (fast health checks at 10 s). A health check fails after 3 consecutive failures. With fast health checks, detection takes ~30 seconds + propagation time. Health check false positives from transient network conditions cause unnecessary failovers — implement health check endpoints that validate application-layer health, not just TCP connectivity.

### Global Accelerator

AWS Global Accelerator provides two static anycast IP addresses that route traffic to the nearest healthy AWS endpoint via the AWS global network backbone — bypassing the public internet.

- **Anycast routing:** A single IP address is advertised from multiple AWS edge locations globally. A client's request is routed to the nearest edge, then travels over the AWS backbone to the target region. Reduces jitter and latency compared to public internet routing.
- **Health checks and failover:** Global Accelerator continuously monitors endpoint health. On failure, it shifts traffic to the next healthy endpoint in approximately 30–60 seconds — faster than DNS TTL-based failover.
- **Use case:** Customer-facing APIs (checkout, product search) where < 60 s failover is required. Global Accelerator provides faster failover than Route53 because it does not depend on DNS TTL expiry.

### DynamoDB Global Tables

DynamoDB Global Tables replicate a table across multiple AWS regions with **multi-master writes**. Any region can serve both reads and writes.

- **Replication:** Writes are propagated to all replica regions via DynamoDB Streams. Replication lag is typically < 1 second but can grow to seconds during a cross-region network event.
- **Conflict resolution:** **Last-writer-wins (LWW)** based on timestamp. If two regions write to the same item simultaneously (split-brain scenario), the write with the higher timestamp wins and the other is discarded. This is an eventual consistency guarantee, not a strong one.
- **Strong consistency:** Only available for reads in the local region. A read from `eu-west-1` will not reflect a write just completed in `us-east-1` until replication propagates.
- **Use case:** Order status, user sessions, inventory soft counts — data where eventual consistency is acceptable for reads and conflict resolution by LWW is acceptable for writes.
- **Not suitable for:** Inventory hard reservation (requires strong consistency and conditional writes within one region), financial ledger entries (LWW can silently discard a debit), or anything requiring serializable isolation.

### Aurora Global Database

Aurora Global Database extends a single Aurora cluster across up to 6 regions. One region is the **primary** (read/write); up to 5 regions are **secondary** (read-only replicas).

- **Replication lag:** < 1 second typical for secondary regions.
- **RPO:** < 1 minute. On primary region failure, at most 1 minute of writes may be lost.
- **Failover (promote):** Promoting a secondary to primary is a **manual operation** that takes approximately 1 minute. This is not automatic — an operator or automated runbook must trigger the promotion. Automatic failover is a planned feature but not generally available at this writing.
- **RTO:** ~1 minute after promotion is initiated, plus application DNS update time.
- **Use case:** Relational data requiring SQL queries: seller accounts, payment records, order history for reporting. Secondary regions serve read traffic (product reviews, order history display) without affecting the primary.
- **Trade-off vs DynamoDB Global Tables:** Aurora is relational and supports ad-hoc SQL but requires manual failover promotion and has a single write region. DynamoDB Global Tables is multi-master but requires access-pattern-first schema design and LWW conflict resolution.

### S3 Cross-Region Replication (CRR)

S3 CRR asynchronously replicates objects from a source bucket to one or more destination buckets in different regions. Replication lag is typically seconds to minutes. S3 objects are eventually consistent after replication — reads in the destination region may serve a stale version immediately after a write in the source region.

- **Use case:** Product images, static assets, data lake exports, backup archives.
- **Not suitable for:** Mutable state that requires consistent reads across regions.

### Disaster Recovery Tiers

DR tier selection is driven by RTO and RPO requirements:

| DR Tier | RTO | RPO | Cost multiplier |
|---------|-----|-----|----------------|
| BACKUP_RESTORE | > 60 min | > 60 min | 1× (baseline) |
| PILOT_LIGHT | 15–60 min | 5–60 min | 2× |
| WARM_STANDBY | 1–15 min | 1–5 min | 4× |
| MULTI_ACTIVE | < 1 min | < 1 min | 8× |

**Tier selection thresholds used in this lesson:**
- `MULTI_ACTIVE`: RTO < 1 min AND RPO < 1 min
- `WARM_STANDBY`: RTO ≤ 15 min AND RPO ≤ 5 min
- `PILOT_LIGHT`: RTO ≤ 60 min AND RPO ≤ 60 min
- `BACKUP_RESTORE`: RTO > 60 min OR RPO > 60 min

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Split-brain on DynamoDB Global Tables | Two regions accept conflicting writes during network partition | Accept LWW semantics; avoid Global Tables for financial writes; use single-region strong consistency for reservations |
| Replication lag causing stale reads across regions | Secondary region serves stale inventory count; customer oversells | Route inventory reservation reads to primary region only; use strong consistency within local region |
| Route53 health check false positive | Transient network blip triggers failover; customers routed to warm region with cold caches | Use application-layer health checks; require 3 consecutive failures before routing change; implement jitter on retry |
| Aurora failover: promotion delay | Primary region fails; secondary promotion takes 2–3 minutes during operator response lag | Automate promotion via AWS Lambda + CloudWatch alarm; pre-warm application connection pools in secondary |
| Global Accelerator failover to wrong region | Traffic routed to region with stale data after DynamoDB replication lag | Implement read-your-writes consistency by routing reads to the last-written region for session-scoped operations |

## Observability

- metric: Route53 `HealthCheckStatus` per endpoint — alert if 0 (failed) for > 30 s
- metric: DynamoDB Global Tables `ReplicationLatency` per region pair — alert if > 5,000 ms
- metric: Aurora Global Database `AuroraGlobalDBReplicationLag` — alert if > 2,000 ms
- metric: Global Accelerator `NewFlowCount` and `ProcessedBytesIn` per endpoint group — detect traffic shift on failover
- log: application-level cross-region read staleness events (read timestamp vs write timestamp delta)
- runbook: documented Aurora secondary promotion procedure with target RTO of < 3 min from alarm to promotion complete

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| DynamoDB Global Tables over Aurora Global for order status | Multi-master writes from any region; no manual failover | LWW conflict resolution; no strong cross-region consistency | Aurora Global requires manual failover and single write region |
| Aurora Global for payment records | SQL, ACID, relational integrity; < 1s replication lag to read replicas | Single write region; manual failover; ~1 min RPO | DynamoDB for payment records risks LWW discarding writes silently |
| Global Accelerator over Route53-only failover | 30–60 s failover without DNS TTL expiry delay; anycast reduces latency | ~$0.025/hour per accelerator plus data transfer cost | Route53 failover depends on DNS TTL (60 s min) and client cache behavior |
| Active-active for customer-facing paths | Near-zero RTO; serves users from nearest region; full load tolerance | Application must handle eventual consistency and LWW conflicts; much higher cost | Active-passive simpler but has 1–5 min RTO and cold cache on failover |

## Interview It

**Amazon framing:** "Design Amazon.com to survive a complete us-east-1 outage during Prime Day." LP tie-in: Ownership (you own the failover runbook end to end), Dive Deep (explain DynamoDB LWW conflict resolution and when it is not safe), Bias for Action (multi-active is expensive — justify the cost with Prime Day revenue at risk).

**Follow-ups:**
1. us-east-1 just failed. Walk me through exactly what happens in your architecture in the first 5 minutes.
2. DynamoDB Global Tables uses last-writer-wins. Give me a concrete scenario where this is dangerous for Amazon's order system.
3. Why does Aurora Global Database require manual failover? What is the AWS-recommended automation for this?
4. A customer places an order in `us-east-1`. 500 ms later, us-east-1 fails. The order write was acknowledged. Is that order visible in `eu-west-1`? Why or why not?
5. What is the difference in cost between warm-standby and multi-active for a service processing 10,000 requests/second?

## Ship It

- `outputs/multi-region-architecture-checklist.md`

## Exercises

1. **Easy** — List the five AWS disaster recovery tiers and their RTO/RPO ranges. For each, name one Amazon service that would be appropriate for that tier.
2. **Medium** — Design the Route53 + Global Accelerator routing setup for Amazon's checkout service across `us-east-1` and `eu-west-1`. Specify routing policy, health check interval, and failover thresholds.
3. **Hard** — A split-brain event occurs: both `us-east-1` and `eu-west-1` accept a write to the same DynamoDB Global Tables item within the same 500 ms window. Trace exactly what happens with timestamps. Which write survives? What is lost? How would you detect this in production?

## Further Reading

- [AWS Well-Architected — Reliability Pillar](https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/welcome.html)
- [DynamoDB Global Tables](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GlobalTables.html)
- [Aurora Global Database](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/aurora-global-database.html)
- [AWS Global Accelerator Developer Guide](https://docs.aws.amazon.com/global-accelerator/latest/dg/what-is-global-accelerator.html)
- [Werner Vogels: Eventually Consistent](https://dl.acm.org/doi/10.1145/1435417.1435432) — foundational paper on consistency models for distributed systems
