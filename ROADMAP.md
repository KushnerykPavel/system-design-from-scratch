# Roadmap — System Design from Scratch

31 phases · 231 planned lessons · ~332 hours.

> Status legend: ✅ done · 🚧 in progress · ⬚ planned

## Phases

| #  | Phase | Hours |
|----|-------|-------|
| 00 | [Setup & Workflow](phases/00-setup-and-workflow) | ~8h |
| 01 | [Clarification & Scope Control](phases/01-clarification-and-scope) | ~8h |
| 02 | [Back-of-the-Envelope Estimation & Cost](phases/02-estimation-and-cost) | ~10h |
| 03 | [System Design Framework & Timing](phases/03-design-framework-and-timing) | ~8h |
| 04 | [APIs, Contracts & Schema Evolution](phases/04-apis-contracts-and-schema-evolution) | ~10h |
| 05 | [Storage, Indexing & Access Patterns](phases/05-storage-indexing-and-access-patterns) | ~10h |
| 06 | [Caching & Invalidation](phases/06-caching-and-invalidation) | ~10h |
| 07 | [Queues, Streams & Workflows](phases/07-queues-streams-and-workflows) | ~10h |
| 08 | [Consistency, Replication & Transactions](phases/08-consistency-replication-and-transactions) | ~12h |
| 09 | [Partitioning, Sharding & Rebalancing](phases/09-partitioning-sharding-and-rebalancing) | ~10h |
| 10 | [Reliability, Retries & Backpressure](phases/10-reliability-retries-and-backpressure) | ~10h |
| 11 | [Observability, SLOs & Incident Debugging](phases/11-observability-slos-and-debugging) | ~10h |
| 12 | [Security, Abuse & Multitenancy](phases/12-security-abuse-and-multitenancy) | ~10h |
| 13 | [Multi-Region, CDN & Edge Traffic](phases/13-multi-region-cdn-and-edge-traffic) | ~10h |
| 14 | [Rate Limiters, IDs & Consistent Hashing](phases/14-rate-limiters-ids-and-hashing) | ~12h |
| 15 | [KV Stores, Cache Clusters & Object Storage](phases/15-kv-cache-and-object-storage) | ~12h |
| 16 | [Application Backends](phases/16-application-backends) | ~12h |
| 17 | [Search, Crawl & Monitoring Systems](phases/17-search-crawl-and-monitoring-systems) | ~10h |
| 18 | [Messaging & Job Platforms](phases/18-messaging-and-job-platforms) | ~10h |
| 19 | [Payments, Wallets & Ordering Consistency](phases/19-payments-wallets-and-ordering) | ~10h |
| 20 | [Low-Latency, Location & Market Systems](phases/20-low-latency-location-and-market-systems) | ~12h |
| 21 | [Google Senior/Staff System Design](phases/21-google-senior-staff-system-design) | ~12h |
| 22 | [Cloudflare Edge & Platform Design](phases/22-cloudflare-edge-platform-design) | ~12h |
| 23 | [Mixed Mocks & Redesign Drills](phases/23-mixed-mocks-and-redesign-drills) | ~12h |
| 24 | [Netflix Streaming & Platform Design](phases/24-netflix-streaming-platform-design) | ~14h |
| 25 | [Common Interview Gaps](phases/25-common-interview-gaps) | ~10h |
| 26 | [Meta Social Platform Design](phases/26-meta-social-platform-design) | ~12h |
| 27 | [Amazon E-Commerce & AWS Systems Design](phases/27-amazon-ecommerce-aws-design) | ~12h |
| 28 | [Stripe Payments Infrastructure](phases/28-stripe-payments-infrastructure) | ~12h |
| 29 | [Uber Real-Time Platform Design](phases/29-uber-realtime-platform-design) | ~12h |
| 30 | [LinkedIn Professional Network Design](phases/30-linkedin-professional-network-design) | ~12h |

## Phase detail

### Phase 00 — Setup & Workflow (~8 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Repo Setup and Progress Tracking](phases/00-setup-and-workflow/01-repo-setup-and-progress) | ⬚ | ~45 min |
| 02 | [Architecture Note-Taking System](phases/00-setup-and-workflow/02-note-taking-system) | ⬚ | ~45 min |
| 03 | [Diagramming Fast Under Interview Pressure](phases/00-setup-and-workflow/03-fast-diagramming) | ⬚ | ~60 min |
| 04 | [Mistake Log and Feedback Loop](phases/00-setup-and-workflow/04-mistake-log) | ⬚ | ~45 min |
| 05 | [How to Use Capacity Sheets](phases/00-setup-and-workflow/05-capacity-sheets) | ⬚ | ~60 min |
| 06 | [Architecture Review Checklist](phases/00-setup-and-workflow/06-review-checklist) | ⬚ | ~60 min |
| 07 | [Mock Interview Workflow](phases/00-setup-and-workflow/07-mock-workflow) | ⬚ | ~60 min |
| 08 | [Design Debriefs and Iteration Cadence](phases/00-setup-and-workflow/08-design-debriefs) | ⬚ | ~45 min |

### Phase 01 — Clarification & Scope Control (~8 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Functional vs Non-Functional Requirements](phases/01-clarification-and-scope/01-functional-vs-non-functional) | ⬚ | ~60 min |
| 02 | [Clarifying Questions That Actually Change the Design](phases/01-clarification-and-scope/02-high-leverage-questions) | ⬚ | ~60 min |
| 03 | [How to Cut Scope Without Dodging the Prompt](phases/01-clarification-and-scope/03-scope-cuts) | ⬚ | ~60 min |
| 04 | [Assumption Logging Under Ambiguity](phases/01-clarification-and-scope/04-assumption-logging) | ⬚ | ~60 min |
| 05 | [Prioritizing Requirements with the Interviewer](phases/01-clarification-and-scope/05-prioritization) | ⬚ | ~60 min |
| 06 | [User Journeys and Workload Shape](phases/01-clarification-and-scope/06-workload-shape) | ⬚ | ~60 min |
| 07 | [Senior-Level Clarification Anti-Patterns](phases/01-clarification-and-scope/07-anti-patterns) | ⬚ | ~45 min |
| 08 | [Prompt Reframing Drill](phases/01-clarification-and-scope/08-prompt-reframing) | ⬚ | ~45 min |

### Phase 02 — Back-of-the-Envelope Estimation & Cost (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [QPS and Request Mix Estimation](phases/02-estimation-and-cost/01-qps-and-request-mix) | ⬚ | ~75 min |
| 02 | [Storage Growth and Retention Math](phases/02-estimation-and-cost/02-storage-growth) | ⬚ | ~75 min |
| 03 | [Bandwidth and Egress Cost](phases/02-estimation-and-cost/03-bandwidth-and-egress) | ⬚ | ~75 min |
| 04 | [Cache Hit Rate and Origin Load](phases/02-estimation-and-cost/04-cache-hit-rate) | ⬚ | ~75 min |
| 05 | [Peak Factors, Burstiness, and Queue Build-Up](phases/02-estimation-and-cost/05-burstiness) | ⬚ | ~75 min |
| 06 | [Rough Cost Modeling Without Getting Lost](phases/02-estimation-and-cost/06-cost-modeling) | ⬚ | ~75 min |
| 07 | [Bottleneck Math for CPU, Disk, and Network](phases/02-estimation-and-cost/07-bottleneck-math) | ⬚ | ~75 min |
| 08 | [Estimation Under Uncertain Inputs](phases/02-estimation-and-cost/08-uncertain-inputs) | ⬚ | ~60 min |

### Phase 03 — System Design Framework & Timing (~8 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Four-Step Interview Loop](phases/03-design-framework-and-timing/01-four-step-interview-loop) | ✅ | ~60 min |
| 02 | [How to Spend 45 Minutes Wisely](phases/03-design-framework-and-timing/02-time-boxing) | ✅ | ~60 min |
| 03 | [High-Level Diagram First, Deep Dive Second](phases/03-design-framework-and-timing/03-diagram-then-dive) | ✅ | ~60 min |
| 04 | [Choosing the Right Deep Dive](phases/03-design-framework-and-timing/04-deep-dive-selection) | ✅ | ~60 min |
| 05 | [How to Wrap Up Like a Senior Engineer](phases/03-design-framework-and-timing/05-wrap-up) | ✅ | ~45 min |
| 06 | [Constraint Change and Redesign Prompts](phases/03-design-framework-and-timing/06-redesign-prompts) | ✅ | ~45 min |
| 07 | [Common Interviewer Moves and How to Respond](phases/03-design-framework-and-timing/07-interviewer-moves) | ✅ | ~60 min |
| 08 | [Full-Loop Drill: Design a URL Shortener](phases/03-design-framework-and-timing/08-full-loop-drill) | ✅ | ~60 min |

### Phase 04 — APIs, Contracts & Schema Evolution (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [HTTP vs gRPC vs Event Interfaces](phases/04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events) | ✅ | ~75 min |
| 02 | [Idempotency Keys and Safe Retries](phases/04-apis-contracts-and-schema-evolution/02-idempotency-keys) | ✅ | ~75 min |
| 03 | [Pagination, Filtering, and Query Shape](phases/04-apis-contracts-and-schema-evolution/03-pagination-and-filtering) | ✅ | ~75 min |
| 04 | [API Versioning and Compatibility](phases/04-apis-contracts-and-schema-evolution/04-api-versioning) | ✅ | ~75 min |
| 05 | [Schema Evolution in Event-Driven Systems](phases/04-apis-contracts-and-schema-evolution/05-event-schema-evolution) | ✅ | ~75 min |
| 06 | [Ownership Boundaries and Contract Testing](phases/04-apis-contracts-and-schema-evolution/06-contract-testing) | ✅ | ~75 min |
| 07 | [Public API Safety Defaults](phases/04-apis-contracts-and-schema-evolution/07-api-safety-defaults) | ✅ | ~60 min |
| 08 | [Interface Design Drill](phases/04-apis-contracts-and-schema-evolution/08-interface-drill) | ✅ | ~60 min |

### Phase 05 — Storage, Indexing & Access Patterns (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Relational vs KV vs Document Stores](phases/05-storage-indexing-and-access-patterns/01-storage-models) | ✅ | ~75 min |
| 02 | [Primary Access Pattern First](phases/05-storage-indexing-and-access-patterns/02-access-pattern-first) | ✅ | ~75 min |
| 03 | [Indexes, Secondary Indexes, and Write Amplification](phases/05-storage-indexing-and-access-patterns/03-indexes) | ✅ | ~75 min |
| 04 | [Hot Rows, Cold Data, and Tiering](phases/05-storage-indexing-and-access-patterns/04-hot-and-cold-data) | ✅ | ~75 min |
| 05 | [Blob Metadata Separation](phases/05-storage-indexing-and-access-patterns/05-blob-metadata-separation) | ✅ | ~60 min |
| 06 | [Time-Series and Append-Heavy Workloads](phases/05-storage-indexing-and-access-patterns/06-time-series) | ✅ | ~75 min |
| 07 | [Retention, Deletion, and Compliance](phases/05-storage-indexing-and-access-patterns/07-retention-and-deletion) | ✅ | ~60 min |
| 08 | [Storage Deep-Dive Drill](phases/05-storage-indexing-and-access-patterns/08-storage-drill) | ✅ | ~60 min |

### Phase 06 — Caching & Invalidation (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Read-Through, Write-Through, and Write-Behind](phases/06-caching-and-invalidation/01-cache-patterns) | ✅ | ~75 min |
| 02 | [TTL, Explicit Invalidation, and Freshness Models](phases/06-caching-and-invalidation/02-freshness-models) | ✅ | ~75 min |
| 03 | [Eviction Policies and Working Set Shape](phases/06-caching-and-invalidation/03-eviction-policies) | ✅ | ~75 min |
| 04 | [Cache Stampedes and Request Coalescing](phases/06-caching-and-invalidation/04-cache-stampede) | ✅ | ~75 min |
| 05 | [Negative Caching and Error Caching](phases/06-caching-and-invalidation/05-negative-caching) | ✅ | ~60 min |
| 06 | [CDN, Browser, and Edge Cache Layers](phases/06-caching-and-invalidation/06-cache-layers) | ✅ | ~75 min |
| 07 | [Consistency Trade-offs in Cached Systems](phases/06-caching-and-invalidation/07-cache-consistency) | ✅ | ~60 min |
| 08 | [Caching Design Drill](phases/06-caching-and-invalidation/08-caching-drill) | ✅ | ~60 min |

### Phase 07 — Queues, Streams & Workflows (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Queues vs Streams vs Workflows](phases/07-queues-streams-and-workflows/01-queues-vs-streams) | ⬚ | ~75 min |
| 02 | [At-Most-Once, At-Least-Once, and Exactly-Once Claims](phases/07-queues-streams-and-workflows/02-delivery-semantics) | ⬚ | ~75 min |
| 03 | [Partitioning and Consumer Groups](phases/07-queues-streams-and-workflows/03-consumer-groups) | ⬚ | ~75 min |
| 04 | [Dead-Letter Queues and Replay](phases/07-queues-streams-and-workflows/04-dlq-and-replay) | ⬚ | ~60 min |
| 05 | [Workflow Engines and Long-Running State](phases/07-queues-streams-and-workflows/05-workflow-engines) | ⬚ | ~75 min |
| 06 | [Outbox and CDC Patterns](phases/07-queues-streams-and-workflows/06-outbox-and-cdc) | ⬚ | ~75 min |
| 07 | [Backpressure in Event Pipelines](phases/07-queues-streams-and-workflows/07-pipeline-backpressure) | ⬚ | ~60 min |
| 08 | [Messaging Architecture Drill](phases/07-queues-streams-and-workflows/08-messaging-drill) | ⬚ | ~60 min |

### Phase 08 — Consistency, Replication & Transactions (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Consistency Spectrum in Practice](phases/08-consistency-replication-and-transactions/01-consistency-spectrum) | ⬚ | ~75 min |
| 02 | [Leader-Follower Replication](phases/08-consistency-replication-and-transactions/02-leader-follower) | ⬚ | ~75 min |
| 03 | [Quorums, Read Repair, and Divergence](phases/08-consistency-replication-and-transactions/03-quorums) | ⬚ | ~75 min |
| 04 | [Replication Lag and Read Freshness](phases/08-consistency-replication-and-transactions/04-replication-lag) | ⬚ | ~75 min |
| 05 | [Transactions, Isolation, and Hotspotting](phases/08-consistency-replication-and-transactions/05-transactions) | ⬚ | ~75 min |
| 06 | [Sagas and Compensating Actions](phases/08-consistency-replication-and-transactions/06-sagas) | ⬚ | ~75 min |
| 07 | [Clock Skew, Ordering, and Time Assumptions](phases/08-consistency-replication-and-transactions/07-time-assumptions) | ⬚ | ~60 min |
| 08 | [Consistency Trade-off Drill](phases/08-consistency-replication-and-transactions/08-consistency-drill) | ⬚ | ~60 min |

### Phase 09 — Partitioning, Sharding & Rebalancing (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Choosing a Shard Key](phases/09-partitioning-sharding-and-rebalancing/01-shard-key) | ⬚ | ~75 min |
| 02 | [Hot Partitions and Skew](phases/09-partitioning-sharding-and-rebalancing/02-hot-partitions) | ⬚ | ~75 min |
| 03 | [Consistent Hashing and Placement](phases/09-partitioning-sharding-and-rebalancing/03-placement) | ⬚ | ~75 min |
| 04 | [Rebalancing Without Taking an Outage](phases/09-partitioning-sharding-and-rebalancing/04-rebalancing) | ⬚ | ~75 min |
| 05 | [Tenant Isolation and Noisy Neighbors](phases/09-partitioning-sharding-and-rebalancing/05-tenant-isolation) | ⬚ | ~60 min |
| 06 | [Resharding and Data Migration Plans](phases/09-partitioning-sharding-and-rebalancing/06-resharding) | ⬚ | ~75 min |
| 07 | [Cross-Shard Queries and Aggregation](phases/09-partitioning-sharding-and-rebalancing/07-cross-shard-queries) | ⬚ | ~60 min |
| 08 | [Sharding Drill](phases/09-partitioning-sharding-and-rebalancing/08-sharding-drill) | ⬚ | ~60 min |

### Phase 10 — Reliability, Retries & Backpressure (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Timeouts, Retries, and Retry Storms](phases/10-reliability-retries-and-backpressure/01-timeouts-and-retries) | ⬚ | ~75 min |
| 02 | [Idempotency Under Partial Failure](phases/10-reliability-retries-and-backpressure/02-idempotency-under-failure) | ⬚ | ~75 min |
| 03 | [Circuit Breakers and Graceful Degradation](phases/10-reliability-retries-and-backpressure/03-circuit-breakers) | ⬚ | ~75 min |
| 04 | [Admission Control and Load Shedding](phases/10-reliability-retries-and-backpressure/04-load-shedding) | ⬚ | ~75 min |
| 05 | [Backpressure Across Async Boundaries](phases/10-reliability-retries-and-backpressure/05-async-backpressure) | ⬚ | ~75 min |
| 06 | [Retry Budgets and Hedging](phases/10-reliability-retries-and-backpressure/06-retry-budgets) | ⬚ | ~60 min |
| 07 | [Bulkheads and Failure Isolation](phases/10-reliability-retries-and-backpressure/07-bulkheads) | ⬚ | ~60 min |
| 08 | [Reliability Drill](phases/10-reliability-retries-and-backpressure/08-reliability-drill) | ⬚ | ~60 min |

### Phase 11 — Observability, SLOs & Incident Debugging (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [SLIs, SLOs, and Error Budgets](phases/11-observability-slos-and-debugging/01-sli-slo-error-budget) | ⬚ | ~75 min |
| 02 | [Metrics That Actually Explain the System](phases/11-observability-slos-and-debugging/02-metrics-that-matter) | ⬚ | ~75 min |
| 03 | [Logs, Traces, and Correlation IDs](phases/11-observability-slos-and-debugging/03-logs-and-traces) | ⬚ | ~75 min |
| 04 | [Dashboards and Cardinality Discipline](phases/11-observability-slos-and-debugging/04-dashboards) | ⬚ | ~75 min |
| 05 | [Alert Design and Paging Quality](phases/11-observability-slos-and-debugging/05-alert-design) | ⬚ | ~60 min |
| 06 | [Runbooks and First-Response Workflow](phases/11-observability-slos-and-debugging/06-runbooks) | ⬚ | ~60 min |
| 07 | [Incident Debugging Narrative](phases/11-observability-slos-and-debugging/07-debugging-narrative) | ⬚ | ~60 min |
| 08 | [Observability Drill](phases/11-observability-slos-and-debugging/08-observability-drill) | ⬚ | ~60 min |

### Phase 12 — Security, Abuse & Multitenancy (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Authentication, Authorization, and Trust Boundaries](phases/12-security-abuse-and-multitenancy/01-auth-and-trust) | ⬚ | ~75 min |
| 02 | [Secrets, Key Management, and Rotation](phases/12-security-abuse-and-multitenancy/02-secrets-and-keys) | ⬚ | ~75 min |
| 03 | [Abuse Prevention and Rate Limiting Layers](phases/12-security-abuse-and-multitenancy/03-abuse-prevention) | ⬚ | ~75 min |
| 04 | [Tenant Isolation and Blast Radius](phases/12-security-abuse-and-multitenancy/04-tenant-isolation) | ⬚ | ~75 min |
| 05 | [Privacy, Retention, and Deletion Semantics](phases/12-security-abuse-and-multitenancy/05-privacy-and-deletion) | ⬚ | ~60 min |
| 06 | [Threat Modeling for Interview Design](phases/12-security-abuse-and-multitenancy/06-threat-modeling) | ⬚ | ~60 min |
| 07 | [Secure Defaults Drill](phases/12-security-abuse-and-multitenancy/07-secure-defaults-drill) | ⬚ | ~60 min |

### Phase 13 — Multi-Region, CDN & Edge Traffic (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Active-Active vs Active-Passive](phases/13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive) | ⬚ | ~75 min |
| 02 | [Regional Failover and Recovery Objectives](phases/13-multi-region-cdn-and-edge-traffic/02-failover-and-rto) | ⬚ | ~75 min |
| 03 | [CDN Layering and Cache Hierarchies](phases/13-multi-region-cdn-and-edge-traffic/03-cdn-layering) | ⬚ | ~75 min |
| 04 | [Traffic Steering, Anycast, and Regional Routing](phases/13-multi-region-cdn-and-edge-traffic/04-traffic-steering) | ⬚ | ~75 min |
| 05 | [Edge Compute and Data Gravity](phases/13-multi-region-cdn-and-edge-traffic/05-edge-compute) | ⬚ | ~60 min |
| 06 | [Geo-Distributed Consistency Trade-offs](phases/13-multi-region-cdn-and-edge-traffic/06-geo-consistency) | ⬚ | ~60 min |
| 07 | [Global Traffic Drill](phases/13-multi-region-cdn-and-edge-traffic/07-global-traffic-drill) | ⬚ | ~60 min |

### Phase 14 — Rate Limiters, IDs & Consistent Hashing (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Token Bucket vs Sliding Window vs Leaky Bucket](phases/14-rate-limiters-ids-and-hashing/01-rate-limiter-primitives) | ✅ | ~75 min |
| 02 | [Distributed Rate Limiter](phases/14-rate-limiters-ids-and-hashing/02-distributed-rate-limiter) | ✅ | ~90 min |
| 03 | [Unique ID Generator: Snowflake, DB, and Random IDs](phases/14-rate-limiters-ids-and-hashing/03-unique-id-generator) | ✅ | ~75 min |
| 04 | [Consistent Hashing and Ring Rebalancing](phases/14-rate-limiters-ids-and-hashing/04-consistent-hashing) | ✅ | ~75 min |
| 05 | [Service Discovery and Placement Decisions](phases/14-rate-limiters-ids-and-hashing/05-service-discovery-placement) | ✅ | ~60 min |
| 06 | [Hot-Key Mitigation Strategies](phases/14-rate-limiters-ids-and-hashing/06-hot-key-mitigation) | ✅ | ~60 min |
| 07 | [Control-Plane vs Data-Plane Trade-offs Drill](phases/14-rate-limiters-ids-and-hashing/07-control-vs-data-plane-drill) | ✅ | ~60 min |

### Phase 15 — KV Stores, Cache Clusters & Object Storage (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Distributed Key-Value Store](phases/15-kv-cache-and-object-storage/01-distributed-kv-store) | ✅ | ~90 min |
| 02 | [Distributed Cache Cluster](phases/15-kv-cache-and-object-storage/02-distributed-cache-cluster) | ✅ | ~75 min |
| 03 | [Object Storage and Blob Metadata](phases/15-kv-cache-and-object-storage/03-object-storage) | ✅ | ~90 min |
| 04 | [Metadata Index Service for Large Blobs](phases/15-kv-cache-and-object-storage/04-metadata-index-service) | ✅ | ~60 min |
| 05 | [Durability Tiers and Data Repair](phases/15-kv-cache-and-object-storage/05-durability-tiers) | ✅ | ~60 min |
| 06 | [Compaction, GC, and Lifecycle Policies](phases/15-kv-cache-and-object-storage/06-compaction-and-lifecycle) | ✅ | ~60 min |
| 07 | [Storage Platform Drill](phases/15-kv-cache-and-object-storage/07-storage-platform-drill) | ✅ | ~60 min |

### Phase 16 — Application Backends (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [URL Shortener](phases/16-application-backends/01-url-shortener) | ✅ | ~75 min |
| 02 | [News Feed / Timeline](phases/16-application-backends/02-news-feed) | ✅ | ~90 min |
| 03 | [Chat System](phases/16-application-backends/03-chat-system) | ✅ | ~90 min |
| 04 | [Notification System](phases/16-application-backends/04-notification-system) | ✅ | ~75 min |
| 05 | [Collaborative Document / Presence Backend](phases/16-application-backends/05-collaboration-backend) | ✅ | ~75 min |
| 06 | [Application Fanout Patterns Compared](phases/16-application-backends/06-fanout-patterns) | ✅ | ~60 min |
| 07 | [Backend Product Drill](phases/16-application-backends/07-backend-product-drill) | ✅ | ~60 min |

### Phase 17 — Search, Crawl & Monitoring Systems (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Web Crawler](phases/17-search-crawl-and-monitoring-systems/01-web-crawler) | ✅ | ~75 min |
| 02 | [Search Autocomplete](phases/17-search-crawl-and-monitoring-systems/02-search-autocomplete) | ✅ | ~75 min |
| 03 | [Metrics Platform](phases/17-search-crawl-and-monitoring-systems/03-metrics-platform) | ✅ | ~75 min |
| 04 | [Log Aggregation Pipeline](phases/17-search-crawl-and-monitoring-systems/04-log-pipeline) | ✅ | ~75 min |
| 05 | [Alert Routing and On-Call Signal Quality](phases/17-search-crawl-and-monitoring-systems/05-alert-routing) | ✅ | ~60 min |
| 06 | [Index Freshness and Ranking Updates](phases/17-search-crawl-and-monitoring-systems/06-index-freshness) | ✅ | ~60 min |
| 07 | [Search and Monitoring Drill](phases/17-search-crawl-and-monitoring-systems/07-search-monitoring-drill) | ✅ | ~60 min |

### Phase 18 — Messaging & Job Platforms (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Distributed Message Queue](phases/18-messaging-and-job-platforms/01-distributed-message-queue) | ⬚ | ~90 min |
| 02 | [Workflow Engine for Long-Running Jobs](phases/18-messaging-and-job-platforms/02-workflow-engine) | ⬚ | ~75 min |
| 03 | [Job Scheduler and Retry Platform](phases/18-messaging-and-job-platforms/03-job-scheduler) | ⬚ | ~75 min |
| 04 | [Pub/Sub Fanout and Subscription Isolation](phases/18-messaging-and-job-platforms/04-pubsub-fanout) | ⬚ | ~75 min |
| 05 | [Dead-Letter and Replay Control Plane](phases/18-messaging-and-job-platforms/05-dlq-control-plane) | ⬚ | ~60 min |
| 06 | [Exactly-Once Myths in Interview Design](phases/18-messaging-and-job-platforms/06-exactly-once-myths) | ⬚ | ~60 min |
| 07 | [Messaging Platform Drill](phases/18-messaging-and-job-platforms/07-messaging-platform-drill) | ⬚ | ~60 min |

### Phase 19 — Payments, Wallets & Ordering Consistency (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Payment Ledger](phases/19-payments-wallets-and-ordering/01-payment-ledger) | ⬚ | ~90 min |
| 02 | [Digital Wallet with Holds and Settlements](phases/19-payments-wallets-and-ordering/02-digital-wallet) | ⬚ | ~75 min |
| 03 | [Inventory Reservation System](phases/19-payments-wallets-and-ordering/03-inventory-reservation) | ⬚ | ~75 min |
| 04 | [Order State Machine and Recovery](phases/19-payments-wallets-and-ordering/04-order-state-machine) | ⬚ | ~75 min |
| 05 | [Audit, Compliance, and Data Retention Constraints](phases/19-payments-wallets-and-ordering/05-audit-and-compliance) | ⬚ | ~60 min |
| 06 | [Fraud and Risk Hooks Without Blocking the Core Path](phases/19-payments-wallets-and-ordering/06-fraud-hooks) | ⬚ | ~60 min |
| 07 | [Financial Consistency Drill](phases/19-payments-wallets-and-ordering/07-financial-consistency-drill) | ⬚ | ~60 min |

### Phase 20 — Low-Latency, Location & Market Systems (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Proximity Service / Nearby Search](phases/20-low-latency-location-and-market-systems/01-proximity-service) | ⬚ | ~75 min |
| 02 | [Real-Time Location Update Pipeline](phases/20-low-latency-location-and-market-systems/02-location-updates) | ⬚ | ~75 min |
| 03 | [Real-Time Leaderboard](phases/20-low-latency-location-and-market-systems/03-realtime-leaderboard) | ⬚ | ~75 min |
| 04 | [Stock Exchange Matching Engine](phases/20-low-latency-location-and-market-systems/04-stock-exchange) | ⬚ | ~90 min |
| 05 | [Market Data Fanout and Subscriber Tiers](phases/20-low-latency-location-and-market-systems/05-market-data-fanout) | ⬚ | ~60 min |
| 06 | [Map Tiles and Read-Mostly Geo Data](phases/20-low-latency-location-and-market-systems/06-map-tiles) | ⬚ | ~60 min |
| 07 | [Low-Latency Systems Drill](phases/20-low-latency-location-and-market-systems/07-low-latency-drill) | ⬚ | ~60 min |

### Phase 21 — Google Senior/Staff System Design (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Google Rubric and Strong Signals](phases/21-google-senior-staff-system-design/01-google-rubric) | ⬚ | ~75 min |
| 02 | [Handling Ambiguous Prompts Like a Senior Engineer](phases/21-google-senior-staff-system-design/02-ambiguity-handling) | ⬚ | ~75 min |
| 03 | [Requirement Negotiation and Priority Framing](phases/21-google-senior-staff-system-design/03-priority-framing) | ⬚ | ~75 min |
| 04 | [Consistency Trade-offs for Serving and Storage](phases/21-google-senior-staff-system-design/04-google-consistency-drills) | ⬚ | ~75 min |
| 05 | [Communicating Risks, Rollouts, and Fallbacks](phases/21-google-senior-staff-system-design/05-risk-communication) | ⬚ | ~60 min |
| 06 | [Staff-Level Deep-Dive Selection](phases/21-google-senior-staff-system-design/06-staff-deep-dive) | ⬚ | ~60 min |
| 07 | [Google Full Mock Loop](phases/21-google-senior-staff-system-design/07-google-full-mock) | ⬚ | ~90 min |

### Phase 22 — Cloudflare Edge & Platform Design (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Global API Edge Gateway](phases/22-cloudflare-edge-platform-design/01-global-api-edge-gateway) | 🚧 | ~90 min |
| 02 | [CDN Invalidation and Cache Propagation](phases/22-cloudflare-edge-platform-design/02-cdn-invalidation) | ⬚ | ~75 min |
| 03 | [Origin Protection and Health-Based Failover](phases/22-cloudflare-edge-platform-design/03-origin-protection) | ⬚ | ~75 min |
| 04 | [Traffic Steering Across POPs and Regions](phases/22-cloudflare-edge-platform-design/04-traffic-steering) | ⬚ | ~75 min |
| 05 | [Bot Mitigation and Abuse Control Planes](phases/22-cloudflare-edge-platform-design/05-bot-mitigation) | ⬚ | ~75 min |
| 06 | [Cost, Latency, and Edge Cache Trade-offs](phases/22-cloudflare-edge-platform-design/06-cost-latency-tradeoffs) | ⬚ | ~60 min |
| 07 | [Cloudflare Full Mock Loop](phases/22-cloudflare-edge-platform-design/07-cloudflare-full-mock) | ⬚ | ~90 min |

### Phase 23 — Mixed Mocks & Redesign Drills (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Mixed Mock: Consumer Backend](phases/23-mixed-mocks-and-redesign-drills/01-consumer-backend-mock) | ⬚ | ~90 min |
| 02 | [Mixed Mock: Infra Platform](phases/23-mixed-mocks-and-redesign-drills/02-infra-platform-mock) | ⬚ | ~90 min |
| 03 | [Redesign After a 10x Scale Change](phases/23-mixed-mocks-and-redesign-drills/03-ten-x-redesign) | ⬚ | ~75 min |
| 04 | [Redesign After a Major Failure Mode](phases/23-mixed-mocks-and-redesign-drills/04-failure-redesign) | ⬚ | ~75 min |
| 05 | [Lightning Capacity Rounds](phases/23-mixed-mocks-and-redesign-drills/05-lightning-capacity-rounds) | ⬚ | ~60 min |
| 06 | [Lightning Trade-off Rounds](phases/23-mixed-mocks-and-redesign-drills/06-lightning-tradeoffs) | ⬚ | ~60 min |
| 07 | [Final Capstone Review Panel](phases/23-mixed-mocks-and-redesign-drills/07-final-capstone-review) | ⬚ | ~90 min |

### Phase 24 — Netflix Streaming & Platform Design (~14 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Netflix Interview Rubric & Strong Signals](phases/24-netflix-streaming-platform-design/01-netflix-rubric) | ⬚ | ~75 min |
| 02 | [Video Streaming Pipeline — Encoding, Chunking, ABR](phases/24-netflix-streaming-platform-design/02-video-streaming-pipeline) | ⬚ | ~90 min |
| 03 | [Open Connect CDN — Cache Hierarchy & POP Design](phases/24-netflix-streaming-platform-design/03-open-connect-cdn) | ⬚ | ~90 min |
| 04 | [Recommendation Engine at Streaming Scale](phases/24-netflix-streaming-platform-design/04-recommendation-engine) | ⬚ | ~90 min |
| 05 | [Chaos Engineering — Resilience by Design](phases/24-netflix-streaming-platform-design/05-chaos-engineering) | ⬚ | ~75 min |
| 06 | [A/B Testing & Experimentation Platform](phases/24-netflix-streaming-platform-design/06-ab-testing-platform) | ⬚ | ~75 min |
| 07 | [Real-Time Analytics Pipeline — Kafka, Flink, Druid](phases/24-netflix-streaming-platform-design/07-realtime-data-pipeline) | ⬚ | ~75 min |
| 08 | [Netflix Full Mock Loop](phases/24-netflix-streaming-platform-design/08-netflix-full-mock) | ⬚ | ~90 min |

### Phase 25 — Common Interview Gaps (~10 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [File Sync Service (Dropbox/Drive)](phases/25-common-interview-gaps/01-file-sync-service) | ⬚ | ~90 min |
| 02 | [Ad Click Aggregation Pipeline](phases/25-common-interview-gaps/02-ad-click-aggregation) | ⬚ | ~75 min |
| 03 | [Ride-Sharing Dispatch System](phases/25-common-interview-gaps/03-ride-sharing-dispatch) | ⬚ | ~90 min |
| 04 | [Live Streaming Platform (Twitch-style)](phases/25-common-interview-gaps/04-live-streaming) | ⬚ | ~90 min |
| 05 | [Hotel/Flight Booking & Seat Reservation](phases/25-common-interview-gaps/05-booking-and-reservation) | ⬚ | ~75 min |
| 06 | [Code Deployment Pipeline](phases/25-common-interview-gaps/06-code-deployment-pipeline) | ⬚ | ~75 min |
| 07 | [Feature Flag Service](phases/25-common-interview-gaps/07-feature-flag-service) | ⬚ | ~60 min |
| 08 | [Gap Scenarios Drill](phases/25-common-interview-gaps/08-gap-scenarios-drill) | ⬚ | ~90 min |

### Phase 26 — Meta Social Platform Design (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Meta Interview Rubric & Strong Signals](phases/26-meta-social-platform-design/01-meta-rubric) | ⬚ | ~75 min |
| 02 | [Social Graph — TAO Cache & Graph Serving](phases/26-meta-social-platform-design/02-social-graph-tao) | ⬚ | ~90 min |
| 03 | [News Feed Fanout at 3 Billion Scale](phases/26-meta-social-platform-design/03-news-feed-fanout) | ⬚ | ~90 min |
| 04 | [Real-Time Messaging — WhatsApp & Messenger](phases/26-meta-social-platform-design/04-real-time-messaging) | ⬚ | ~90 min |
| 05 | [Media Pipeline — Photos, Videos & CDN](phases/26-meta-social-platform-design/05-media-pipeline) | ⬚ | ~75 min |
| 06 | [Ranking & ML Feature Serving at Feed Scale](phases/26-meta-social-platform-design/06-ranking-ml-serving) | ⬚ | ~75 min |
| 07 | [Meta Full Mock Loop](phases/26-meta-social-platform-design/07-meta-full-mock) | ⬚ | ~90 min |

### Phase 27 — Amazon E-Commerce & AWS Systems Design (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Amazon Interview Rubric & Strong Signals](phases/27-amazon-ecommerce-aws-design/01-amazon-rubric) | ⬚ | ~75 min |
| 02 | [Product Catalog & Search at E-Commerce Scale](phases/27-amazon-ecommerce-aws-design/02-product-catalog-search) | ⬚ | ~90 min |
| 03 | [Order Management & Fulfillment Pipeline](phases/27-amazon-ecommerce-aws-design/03-order-fulfillment) | ⬚ | ~90 min |
| 04 | [DynamoDB Internals & NoSQL Access Patterns](phases/27-amazon-ecommerce-aws-design/04-dynamodb-patterns) | ⬚ | ~75 min |
| 05 | [Event-Driven Architecture on AWS](phases/27-amazon-ecommerce-aws-design/05-event-driven-aws) | ⬚ | ~75 min |
| 06 | [Multi-Region Reliability on AWS](phases/27-amazon-ecommerce-aws-design/06-aws-multi-region) | ⬚ | ~75 min |
| 07 | [Amazon Full Mock Loop](phases/27-amazon-ecommerce-aws-design/07-amazon-full-mock) | ⬚ | ~90 min |

### Phase 28 — Stripe Payments Infrastructure (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Stripe Interview Rubric & Strong Signals](phases/28-stripe-payments-infrastructure/01-stripe-rubric) | ⬚ | ~75 min |
| 02 | [Payment Processing & Idempotency](phases/28-stripe-payments-infrastructure/02-payment-processing) | ⬚ | ~90 min |
| 03 | [Ledger Design & Reconciliation](phases/28-stripe-payments-infrastructure/03-ledger-design) | ⬚ | ~90 min |
| 04 | [Fraud Detection Pipeline](phases/28-stripe-payments-infrastructure/04-fraud-detection) | ⬚ | ~75 min |
| 05 | [Global Money Movement & Multi-Currency](phases/28-stripe-payments-infrastructure/05-global-money-movement) | ⬚ | ~75 min |
| 06 | [Developer API Design — Webhooks & Versioning](phases/28-stripe-payments-infrastructure/06-developer-api-design) | ⬚ | ~75 min |
| 07 | [Stripe Full Mock Loop](phases/28-stripe-payments-infrastructure/07-stripe-full-mock) | ⬚ | ~90 min |

### Phase 29 — Uber Real-Time Platform Design (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [Uber Interview Rubric & Strong Signals](phases/29-uber-realtime-platform-design/01-uber-rubric) | ⬚ | ~75 min |
| 02 | [Real-Time Geospatial — H3 Indexing & Proximity Search](phases/29-uber-realtime-platform-design/02-geospatial-indexing) | ⬚ | ~90 min |
| 03 | [Ride Dispatch & Matching Engine](phases/29-uber-realtime-platform-design/03-dispatch-matching) | ⬚ | ~90 min |
| 04 | [Surge Pricing & Dynamic Pricing Pipeline](phases/29-uber-realtime-platform-design/04-surge-pricing) | ⬚ | ~75 min |
| 05 | [Real-Time GPS Ingestion & Driver Tracking](phases/29-uber-realtime-platform-design/05-location-tracking) | ⬚ | ~75 min |
| 06 | [Trip State Machine & Recovery](phases/29-uber-realtime-platform-design/06-trip-state-machine) | ⬚ | ~75 min |
| 07 | [Uber Full Mock Loop](phases/29-uber-realtime-platform-design/07-uber-full-mock) | ⬚ | ~90 min |

### Phase 30 — LinkedIn Professional Network Design (~12 hours)

| #  | Lesson | Status | Time |
|----|--------|--------|------|
| 01 | [LinkedIn Interview Rubric & Strong Signals](phases/30-linkedin-professional-network-design/01-linkedin-rubric) | ⬚ | ~75 min |
| 02 | [Professional Graph — Connections & Degree Traversal](phases/30-linkedin-professional-network-design/02-professional-graph) | ⬚ | ~90 min |
| 03 | [Feed Ranking & Content Distribution](phases/30-linkedin-professional-network-design/03-feed-ranking) | ⬚ | ~90 min |
| 04 | [Job Matching — Search & Skills Graph](phases/30-linkedin-professional-network-design/04-job-matching) | ⬚ | ~75 min |
| 05 | [InMail & Messaging Platform](phases/30-linkedin-professional-network-design/05-inmail-messaging) | ⬚ | ~75 min |
| 06 | [Kafka Architecture — Stream Processing & Events](phases/30-linkedin-professional-network-design/06-kafka-stream-infra) | ⬚ | ~75 min |
| 07 | [LinkedIn Full Mock Loop](phases/30-linkedin-professional-network-design/07-linkedin-full-mock) | ⬚ | ~90 min |
