# Amazon Interview Rubric & Strong Signals

> Generic system design answers get you to the phone screen. Amazon Leadership Principle fluency gets you an offer.

**Type:** Concept  
**Company focus:** Amazon  
**Learning goal:** Understand what Amazon actually evaluates in a system design loop — LP demonstration, working backwards from the customer, two-pizza team service ownership, and AWS-native technology selection.  
**Prerequisites:** `10-reliability-retries-and-backpressure/03-circuit-breakers`, `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`  
**Estimated time:** ~75 min  
**Primary artifact:** rubric card + strong-hire signal checklist  

## The Problem

You are preparing for an Amazon senior/staff system design interview. Amazon's bar is high and deeply tied to culture. Interviewers are engineers who built actual Amazon systems and who are simultaneously evaluating technical depth and Leadership Principle alignment. A technically correct answer that ignores the customer perspective, avoids trade-off discussion, or treats AWS as an afterthought will not clear the bar.

Amazon operates at a scale that forces specific decisions: 300M+ active customers, ~1.5M third-party sellers on Marketplace, Prime Day peaks reaching ~100,000 orders per minute, and AWS with 200+ services underpinning both Amazon.com and the broader cloud business. Understanding these numbers and which LP drives which architectural choice is the difference between strong-hire and borderline.

## Clarify

- Which Amazon domain is the question about? E-commerce (order management, catalog, recommendations, fulfillment), AWS infrastructure (DynamoDB internals, S3 consistency, Lambda cold starts), or platform reliability (Route53 health checks, multi-region failover)?
- Is the focus on a customer-facing system (apply Customer Obsession, Bias for Action) or an internal platform (apply Operational Excellence, Frugality)?
- Does the interviewer want depth-first (one subsystem deeply) or breadth-first (end-to-end sketch including fulfillment center integration)?
- Assumption if no answer: customer-facing, senior-level, 45 minutes, breadth-first with one deep dive.

## Requirements

### What Amazon Tests

Amazon interviewers evaluate on two simultaneous tracks:

**Track 1: Technical Depth**
1. **Scale intuition** — Can you size the Amazon problem correctly? The catalog has 350M+ product listings. The fulfillment network spans 175+ fulfillment centers. Prime Video streams to 200M+ Prime members. Numbers that are off by an order of magnitude are a red flag.
2. **Service decomposition** — Can you apply the two-pizza team model? Each service should have clear ownership, an independent data store, and a well-defined API surface. Monolith answers miss Amazon's operating model entirely.
3. **AWS technology selection** — Can you choose the right AWS service for the job and justify the choice? DynamoDB vs Aurora vs ElastiCache each have specific sweet spots. Saying "use a database" is not an answer at Amazon.
4. **Failure-first thinking** — Do you design for the failure before the success path? Amazon runs Chaos Engineering exercises, game days, and pre-mortem reviews as standard practice.
5. **Working backwards** — Does your design start with the customer experience and work backwards to the service contract? The PR/FAQ format (Press Release + FAQ) is Amazon's standard for new system design.

**Track 2: Leadership Principle Alignment**
1. **Customer Obsession** — Are you explicitly connecting each technical decision to a customer outcome?
2. **Dive Deep** — Are you going beyond high-level buzzwords to specific implementation choices with justification?
3. **Bias for Action** — Are you proposing a concrete design even with incomplete information, then iterating?
4. **Frugality** — Are you questioning whether every component is necessary and choosing cost-efficient AWS services?
5. **Ownership** — Do you treat the system as something you would own end-to-end, including operations?

### Non-functional Signals Amazon Cares About

- Single-digit millisecond reads for catalog and session data (DynamoDB sweet spot).
- 99.99% availability for checkout and payment paths — one minute of downtime on Prime Day costs millions.
- Eventual consistency is acceptable for catalog updates; strong consistency required for inventory reservation.
- Cost-at-scale matters: a 10x more expensive solution is a design smell unless justified.
- Idempotency for all write operations — distributed systems duplicate messages; handlers must be safe to replay.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active customers | 300M+ | drives account, session, and recommendation scale |
| Product listings | 350M+ | drives catalog storage, search indexing, and CDN asset scale |
| Third-party sellers | ~1.5M | drives seller onboarding, pricing update, and fraud detection throughput |
| Orders per day (normal) | ~35M | drives order management, fulfillment, and payment processing sizing |
| Orders per minute (Prime Day peak) | ~100,000 | drives peak capacity planning for checkout and inventory |
| AWS services | 200+ | drives service mesh, IAM, and observability complexity |
| Fulfillment centers | 175+ globally | drives inventory distribution, routing, and SLA modeling |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- Draws a web server, a database, and a load balancer.
- Says "store products in RDS" without discussing access patterns.
- Does not mention DynamoDB, SQS, SNS, Kinesis, or Step Functions.
- Ignores idempotency — "just write to the database when the order comes in."
- No mention of fulfillment center integration, seller APIs, or third-party logistics.
- Uses LP terms as decoration ("I'd be customer-obsessed about this") rather than design drivers.

### Strong-Hire Answer Pattern

- Starts with a working-backwards statement: "The customer experience is X; what must be true technically to deliver that?"
- Decomposes into two-pizza-team-sized services: Catalog, Search, Cart, Order, Inventory, Fulfillment, Notifications, each with own data store.
- Selects AWS services by access pattern: DynamoDB for single-item reads, Aurora for reporting, ElastiCache (Redis) for sessions and hot inventory counts, SQS for decoupled order events.
- Identifies the top failure modes: inventory oversell, payment duplicate charge, fulfillment routing failure.
- Closes with observability: CloudWatch metrics, X-Ray traces, named SLOs.

### Strong-Hire Milestone Map

| Time mark | Expected progress |
|-----------|-------------------|
| 5 min | Clarified scope, stated customer experience goal, gave capacity numbers |
| 15 min | Sketched service decomposition with AWS service mapping |
| 25 min | Deep-dived the hardest subsystem (usually inventory reservation or order state machine) |
| 35 min | Addressed interviewer pivot — likely failure mode or alternative technology |
| 45 min | Summarized LP alignment, trade-offs, and what you would instrument first |

## Amazon/AWS-Specific Vocabulary

| Term | What it signals |
|------|-----------------|
| DynamoDB single-table design | Understanding of access pattern-first NoSQL modeling |
| SQS FIFO + deduplication ID | Idempotent, ordered message processing for order events |
| Step Functions | Durable workflow orchestration for multi-step order processing |
| Aurora Multi-AZ | Managed relational HA for transactional data (payments, seller accounts) |
| ElastiCache (Redis) | Sub-millisecond hot data: sessions, cart, inventory counts |
| Kinesis Data Streams | High-throughput event streaming for clickstream, order events |
| CloudFront + S3 | Static asset delivery: product images, HTML at global edge |
| Route53 latency/health routing | Multi-region traffic steering with health-check-based failover |
| Lambda + EventBridge | Event-driven serverless for catalog updates, notification fanout |
| Working Backwards / PR/FAQ | Amazon's mechanism for customer-first requirement definition |

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Inventory oversell during Prime Day | Orders created > inventory count | Pessimistic lock via DynamoDB conditional writes; inventory reservation step before payment |
| SQS message processed twice (duplicate charge) | Payment provider duplicate detection alert | Idempotency key on all charge requests; SQS deduplication ID for FIFO queues |
| Fulfillment routing service down | Order stuck in PAYMENT_CONFIRMED for > 5 min | Step Functions catch + retry; fallback to manual fulfillment queue; SLA alert to ops |
| DynamoDB hot partition on product listing | Throttling errors on item reads | Caching via ElastiCache; adaptive capacity in DynamoDB; request-level retry with jitter |
| Third-party seller pricing storm | Sudden burst of price update writes | SQS buffering of seller price updates; rate limit per seller; async index rebuild |

## Observability

- metric: order placement success rate by region and device type
- metric: inventory reservation latency at p50/p95/p99
- metric: SQS queue depth for order-events (alert if > 10K messages)
- metric: DynamoDB consumed capacity units vs provisioned (throttle rate)
- log: idempotency key collision events with order ID and deduplication outcome
- trace: customer checkout path from cart submit to order confirmation (X-Ray)
- SLO: checkout end-to-end < 300ms at p99; order confirmation delivery < 2s

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| DynamoDB over Aurora for order reads | Single-digit ms reads, infinite horizontal scale | No ad-hoc SQL joins; access patterns must be pre-modeled | Aurora is easier to query but cannot sustain Prime Day read throughput without massive read replicas |
| SQS FIFO for order events | Guaranteed ordering and deduplication per order ID | Max 300 TPS per queue (requires sharding by order ID prefix) | Standard SQS higher throughput but no ordering or deduplication |
| Step Functions for fulfillment | Durable state, built-in retry, audit trail | Cost per state transition; latency per step | Lambda chain without state management loses durability on failure |
| EventBridge for cross-service events | Decoupled producers/consumers; schema registry | Event delivery guarantee is at-least-once; consumers must be idempotent | SNS simpler but no content-based filtering or schema enforcement |

## Interview It

**Amazon framing:** "Design Amazon's order management system for Prime Day." LP tie-in: Customer Obsession (the customer must get an order confirmation within seconds), Dive Deep (how does inventory reservation work under 100K orders/min peak), Ownership (who owns the order state machine end to end?).

**Follow-ups:**
1. What happens to an order that is placed but the payment service is temporarily unavailable?
2. How would you prevent a third-party seller from accidentally being oversold to 10,000 customers during a flash sale?
3. Walk me through how an order event flows from checkout to the fulfillment center routing decision.
4. How does your design handle a customer clicking "Place Order" twice within 500ms due to a double-tap?
5. What is your rollout strategy for changing the order state machine schema without downtime?

## Ship It

- `outputs/rubric-card-amazon.md`

## Exercises

1. **Easy** — List five Amazon/AWS-specific terms and explain the LP or trade-off each one represents.
2. **Medium** — Sketch the working-backwards PR/FAQ for Amazon's checkout service: what is the customer promise, and what must be technically true to deliver it?
3. **Hard** — Write a capacity model for Amazon's inventory reservation system on Prime Day: from customer click to DynamoDB write to SQS event, including failure rates and retry budgets.

## Further Reading

- [Amazon Builder's Library](https://aws.amazon.com/builders-library/) — primary source for idempotency, retry, and distributed systems patterns used internally at Amazon
- [DynamoDB Developer Guide — Best Practices](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html) — access pattern design for single-table modeling
- [Werner Vogels' All Things Distributed](https://www.allthingsdistributed.com/) — CTO blog covering Amazon's distributed systems philosophy
- [Amazon Leadership Principles](https://www.amazon.jobs/content/en/our-workplace/leadership-principles) — the 16 LPs that shape every Amazon design decision
