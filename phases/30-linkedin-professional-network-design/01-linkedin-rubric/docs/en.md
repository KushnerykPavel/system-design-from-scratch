# LinkedIn Interview Rubric & Strong Signals

> Professional network design is graph design at scale — and LinkedIn invented Kafka to prove it.

**Type:** Concept
**Company focus:** LinkedIn
**Learning goal:** Understand what LinkedIn evaluates — graph-scale thinking, Kafka stream processing, professional relevance signals, and privacy-aware data design.
**Prerequisites:** `16-application-backends/02-news-feed`, `07-queues-streams-and-workflows/03-consumer-groups`
**Estimated time:** ~75 min
**Primary artifact:** LinkedIn rubric card + strong-hire signal checklist

## The Problem

You are preparing for a LinkedIn senior/staff system design interview. LinkedIn's bar is explicitly calibrated around professional graph scale, stream-processing depth (LinkedIn invented Kafka), and privacy-aware data design. Interviewers have built Espresso, Venice, and the PYMK pipeline. Generic distributed-systems answers are insufficient — you must demonstrate that you understand the operational reality of 950M members, 67M companies, 15M active job listings, ~1B feed updates per day, and a platform where a privacy violation directly harms professional reputations.

## Clarify

- Which product surface is in scope? Feed, PYMK (People You May Know), job matching, InMail, company pages, or the professional graph itself?
- Is the question a new-system design or scaling a known bottleneck?
- Does the interviewer want breadth (end-to-end architecture) or depth (one subsystem)?
- Are GDPR, data residency, and right-to-be-forgotten in scope?

## Requirements

### What LinkedIn Tests

LinkedIn interviewers assess several dimensions simultaneously:

1. **Graph thinking** — Can you reason about a graph with 950M nodes and 475B+ edges? BFS traversal for PYMK, 2nd-degree connection lookups, and graph partitioning challenges must appear early.
2. **Kafka familiarity** — LinkedIn invented Kafka. Deep stream-processing knowledge is expected: partitions, consumer groups, at-least-once semantics, lag monitoring, and Kafka Streams vs Samza.
3. **Feed ranking** — Can you distinguish professional relevance signals from general social engagement signals? LinkedIn's feed rewards career-relevant content, not just viral content.
4. **Privacy and compliance** — Do you identify GDPR right-to-be-forgotten, data residency requirements, and consent management proactively? Privacy is a hard constraint, not a feature.
5. **Job matching** — Can you design a search and recommendation system that matches 950M member profiles to 15M job listings with sub-second response times?
6. **Member data at scale** — Can you design for 950M member profiles with heterogeneous data (skills, experience, education, endorsements) across multiple storage systems?

### Non-functional Signals LinkedIn Cares About

- PYMK latency: 2nd-degree suggestions must return in under 100ms.
- Feed ranking throughput: ~1B feed updates per day, personalized per member.
- Kafka throughput: LinkedIn's Kafka clusters process billions of messages per day.
- Privacy propagation: GDPR deletion requests must be fulfilled within 30 days (regulatory), but deletions should propagate to search indexes within hours.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Members | 950M | drives graph storage, identity, and PYMK pipeline scale |
| Daily active users | ~30M | drives feed read QPS and ranking compute |
| Companies | 67M | drives company page indexing and job matching |
| Active job listings | 15M | drives job recommendation index size |
| Feed updates/day | ~1B | drives Kafka partition count and fanout pipeline capacity |
| Professional graph edges | ~475B (950M × 500 avg connections) | drives Espresso shard count and BFS traversal cost |
| InMail messages/day | ~100M | drives messaging pipeline and spam detection load |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- Designs PYMK with a real-time BFS without acknowledging that 250K traversals per request are too expensive at scale — precomputation is required.
- Mentions Kafka without explaining partitioning strategy, consumer group design, or at-least-once vs exactly-once trade-offs.
- Treats LinkedIn's feed like a social feed: optimizes for virality rather than professional relevance.
- No GDPR model: ignores right-to-be-forgotten, data residency, or consent management.
- "Use Elasticsearch" for job matching without discussing skills ontology, semantic matching, or recruiter vs member matching asymmetry.
- "Scale horizontally" as the answer to every bottleneck without naming the bottleneck precisely.

### Strong-Hire Answer Pattern

- Opens with a capacity estimate and names the top two structural constraints (BFS cost at 2nd-degree scale and Kafka-based feed pipeline throughput).
- Decomposes the system into services with clear ownership: graph store, PYMK pipeline, feed generation, ranking, job matching, InMail, and privacy enforcement.
- Explicitly distinguishes real-time BFS (acceptable for 1st-degree lookups) from precomputed batch candidates (required for 2nd-degree PYMK).
- Names LinkedIn's own technologies appropriately: Espresso (NoSQL store), Venice (feature store), Kafka (unified log), Pinot (real-time OLAP), Samza (stream processing), Azkaban (workflow scheduler).
- Integrates privacy at the data layer: GDPR right-to-be-forgotten flows through a dedicated deletion pipeline that purges Espresso, Kafka compacted topics, Pinot segments, and search indexes.
- Uses professional-relevance vocabulary in feed design: dwell time over 1st-degree author content, career-relevant topics, network proximity signal.
- Discusses trade-offs explicitly: "I chose batch PYMK precomputation over real-time BFS because 2nd-degree traversal at 250K candidates per query at 30M DAU would require 7.5T operations/day — Spark is the right tool."

### Strong-Hire Milestone Map

| Time mark | Expected progress |
|-----------|-------------------|
| 5 min | Clarified scope, stated top 3 constraints, gave rough capacity numbers |
| 15 min | Sketched end-to-end architecture with named services, identified PYMK precomputation strategy |
| 25 min | Drilled into the hardest subsystem (usually graph traversal or Kafka pipeline) with data flow and failure modes |
| 35 min | Addressed GDPR deletion pipeline and at least one feed ranking integration point |
| 45 min | Summarized trade-offs, named what you would instrument first (Kafka consumer lag, PYMK hit rate) |

## LinkedIn-Specific Vocabulary

| Term | What it signals |
|------|-----------------|
| Espresso | LinkedIn's home-grown NoSQL document store (horizontally scalable, built on MySQL) — shows storage depth |
| Venice | Feature store serving precomputed offline features for real-time ranking — shows ML serving architecture literacy |
| Kafka | Unified log backbone LinkedIn invented in 2011 — deep knowledge expected: partitions, compaction, consumer groups |
| Pinot | Real-time OLAP analytics system for low-latency aggregation — shows analytics pipeline awareness |
| Samza | Stream processing framework built on Kafka — LinkedIn's alternative to Apache Flink |
| Azkaban | Workflow scheduler for Hadoop/Spark batch jobs — shows awareness of offline pipeline orchestration |
| Brooklin | Data streaming service for change-data-capture (CDC) across heterogeneous data stores |
| Wherehows | Data lineage and metadata management platform — shows data governance depth |
| Dali | Remote virtual file system abstracting Hadoop HDFS and other storage — shows data platform depth |
| PYMK | People You May Know — LinkedIn's graph-based connection recommendation engine |

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| PYMK BFS timeout on highly-connected member | p99 latency spike in PYMK API; timeout errors | Cap BFS at depth-2 with cardinality limit (e.g., max 10K 2nd-degree candidates); precompute for top-k connected members |
| Hot Espresso shard for celebrity member | QPS spike on single shard; latency degradation | Replicate read-only copy of high-degree member adjacency list across multiple shards |
| Kafka consumer lag grows during traffic spike | Consumer lag metric exceeds threshold alert | Increase partition count (requires consumer restart); add consumer instances to group |
| Stale PYMK from batch precomputation lag | PYMK suggestions include already-connected members | Add real-time 1st-degree filter at query time; reduce Spark job interval |
| GDPR deletion not propagating to all stores | Compliance audit finds member data in Pinot segments or search index | Deletion pipeline with per-store confirmation and audit log; SLA alert on deletion propagation lag |
| Venice feature store staleness | Ranking model using stale member or item features | Feature freshness metrics per feature group; fallback to heuristic ranking on freshness SLA breach |

## Observability

- metric: Kafka consumer lag per consumer group per topic partition
- metric: PYMK API latency p50/p95/p99 and BFS traversal depth distribution
- metric: feed ranking inference latency by model version
- metric: Espresso read/write latency per keyspace
- metric: GDPR deletion pipeline propagation lag per store (Espresso, Kafka, Pinot, search)
- log: PYMK candidate generation with graph traversal depth, dedup count, and ranking result
- trace: member feed request from candidate generation through Venice feature fetch to ranked result
- alert: Kafka consumer lag > 10 minutes on feed fanout consumer group
- alert: GDPR deletion propagation lag > 4 hours on any store

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Batch PYMK precomputation (Spark daily) over real-time BFS | Avoids 250K 2nd-degree lookups per query at 30M DAU; cache-friendly | Suggestions can be up to 24h stale; already-connected members may appear | Real-time BFS is too expensive: 250K × 500 = 125M edge reads per query at scale |
| Espresso (adjacency list) over native graph DB | Horizontal scaling with consistent hashing; operationally simpler | No native graph query language; BFS must be implemented in application code | Graph DB (Neo4j) cannot scale to 950M nodes with LinkedIn's SLA requirements |
| Kafka as unified log over direct service calls | Decouples producers from consumers; enables replay, CDC, audit trail | Additional operational complexity; at-least-once requires idempotent consumers | Direct RPC would couple feed, PYMK, and analytics pipelines to the same release cycle |
| Venice offline feature store over online feature computation | Pre-computed features served in microseconds; ranking model stable | Feature staleness: offline features updated in batch; online events (likes, views) lag by minutes | Computing all features at ranking time would add 100ms+ to feed latency |

## Interview It

**LinkedIn framing:** "Design LinkedIn's People You May Know." Strong answers decompose into at least: graph store (Espresso adjacency list), batch precomputation pipeline (Spark + Azkaban), online scoring (Venice features + ranking model), and a real-time 1st-degree exclusion filter.

**Follow-ups:**

1. A member with 30K connections joins your system. Walk me through the PYMK pipeline for them.
2. A member files a GDPR right-to-be-forgotten request. What does your deletion pipeline look like?
3. How would you detect that your Kafka consumer group for feed fanout is falling behind?
4. Your PYMK batch job runs once per day. How do you ensure a member who connects in the morning doesn't see the new connection in PYMK suggestions in the afternoon?
5. How do you handle data residency requirements (EU member data must stay in EU datacenters)?

## Ship It

- `outputs/rubric-card-linkedin.md`
- `outputs/strong-hire-checklist-linkedin.md`
- `outputs/interview-card-linkedin-rubric.md`

## Exercises

1. **Easy** — List the eight LinkedIn-specific technologies and explain the design problem each one solves.
2. **Medium** — Sketch a PYMK pipeline architecture that handles both a normal member (500 connections) and a super-connector (30K connections) in a single consistent design.
3. **Hard** — Write a capacity model for the LinkedIn feed pipeline: from member post creation through Kafka fanout to ranked feed delivery with Venice features, with numbers at each stage.

## Further Reading

- [LinkedIn Engineering Blog — Kafka](https://engineering.linkedin.com/kafka) — primary source for Kafka's origin at LinkedIn
- [Espresso: LinkedIn's Distributed Document Store](https://engineering.linkedin.com/espresso/espresso-linkedins-distributed-document-store) — storage architecture
- [Venice: LinkedIn's Feature Store](https://engineering.linkedin.com/blog/2021/the-magic-of-venice) — feature store design
- [Pinot: Realtime Distributed OLAP Datastore](https://engineering.linkedin.com/blog/2019/engineering-smart-infrastructure-for-mobile-at-linkedin) — analytics pipeline
- [PYMK at LinkedIn Scale](https://engineering.linkedin.com/blog/2021/optimizing-people-you-may-know) — PYMK pipeline deep dive
