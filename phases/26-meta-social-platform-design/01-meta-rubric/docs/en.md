# Meta Interview Rubric & Strong Signals

> Graph-scale thinking and privacy-aware design separate Meta offers from polite rejections.

**Type:** Concept
**Company focus:** Meta
**Learning goal:** Understand what Meta evaluates in system design — graph-scale thinking, real-time data pipelines, privacy-aware design, and ML integration. Distinguish strong-hire from generic answers.
**Prerequisites:** `16-application-backends/02-news-feed`, `09-partitioning-sharding-and-rebalancing/02-hot-partitions`
**Estimated time:** ~75 min
**Primary artifact:** rubric card + strong-hire signal checklist

## The Problem

You are preparing for a Meta senior/staff system design interview. Meta's bar is explicitly calibrated around social-graph scale, real-time data pipelines, and privacy-first engineering. Interviewers are engineers who built TAO, Haystack, and the news feed ranking stack. Generic distributed-systems answers are not enough — you need to demonstrate that you understand the operational reality of 3B users, 100B+ social graph edges, and a platform where a privacy bug can be front-page news.

## Clarify

- Which product surface is in scope? News feed, messaging, photos, groups, ads, or the social graph itself?
- Is the question about new-system design or scaling a known bottleneck?
- Does the interviewer want depth on one subsystem or an end-to-end sketch?
- Are privacy and regulatory constraints in scope?

## Requirements

### What Meta Tests

Meta interviewers assess several dimensions simultaneously:

1. **Social-graph scale** — Can you reason about a graph with 3B nodes and 100B+ edges? Numbers and graph-specific data structures must appear early.
2. **Real-time data pipelines** — Can you design a fanout system that delivers news feed updates to millions of users within seconds? Knowing push vs pull trade-offs is mandatory.
3. **Privacy-aware design** — Do you identify privacy constraints proactively? A design that ignores audience controls, data minimization, or right-to-deletion is a red flag at Meta.
4. **ML integration** — Can you hook your design into a ranking/recommendation pipeline? Feed ranking, ad ranking, and spam detection are first-class concerns.
5. **Failure-first thinking** — Do you name failure modes and mitigations before the interviewer prompts you? Meta's systems must degrade gracefully at global scale.

### Non-functional Signals Meta Cares About

- Low-latency graph reads: p99 < 10ms for social graph edge lookups.
- Fanout throughput: a post from a high-follower account must fan out to millions of feed caches within seconds.
- Privacy propagation: audience restriction changes must propagate before the next read, not eventually.
- Availability: feed reads and messaging must remain available during partial datacenter failures.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Monthly active users | 3B | drives graph storage, identity, and authentication scale |
| Daily active users | 500M | drives news feed write and read amplification |
| Social graph edges | 100B+ | drives TAO cache sizing and graph traversal cost |
| News feed writes/sec | ~5M (posts, likes, shares) at peak | drives fanout pipeline and Kafka partition count |
| Photos stored | 100B+ objects | drives Haystack/Everstore object storage capacity |
| Feed ranking inferences/sec | ~1M | drives ML inference fleet sizing |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- Designs a news feed with a single database and a CDN in front.
- Says "use Redis for caching" without specifying what is cached and with what eviction policy.
- Ignores fanout trade-offs: applies push-on-write uniformly to all users including celebrities.
- No privacy model, no mention of audience controls.
- "Scale horizontally" as the answer to every bottleneck question.
- No ML integration — treats ranking as a black box with no design implications.

### Strong-Hire Answer Pattern

- Opens with a capacity estimate and names the top two structural constraints (graph traversal cost and fanout amplification).
- Decomposes the system into services with clear ownership: graph serving, feed generation, ranking, notification, and privacy enforcement.
- Distinguishes push-on-write (for most users) from pull-on-read (for high-follower accounts, i.e., celebrities) — names the threshold explicitly.
- Integrates privacy at the data layer: audience controls enforced before writes propagate, not at read time.
- Names at least one ML integration point: where does the ranking model consume features, and how are feature updates propagated?
- Discusses trade-offs explicitly: "I chose X over Y because Z, at the cost of W."
- Uses Meta-specific vocabulary appropriately: TAO, Haystack, Scuba, ODS, Everstore, Iris, Prophet, Tupperware.

### Strong-Hire Milestone Map

| Time mark | Expected progress |
|-----------|-------------------|
| 5 min | Clarified scope, stated top 3 constraints, gave rough capacity numbers |
| 15 min | Sketched end-to-end architecture with named services, identified fanout strategy |
| 25 min | Drilled into the hardest subsystem (usually graph serving or feed fanout) with data flow and failure modes |
| 35 min | Addressed privacy propagation and at least one ML integration point |
| 45 min | Summarized trade-offs, named what you would instrument first |

## Meta-Specific Vocabulary

| Term | What it signals |
|------|-----------------|
| TAO | Two-level graph cache (objects + associations) — shows graph-serving depth |
| Haystack | Photo object storage optimized for small random reads — shows storage architecture literacy |
| Scuba | In-memory time-series analytics store — shows observability depth |
| ODS (Operational Data Store) | Time-series metrics aggregation pipeline — shows monitoring architecture |
| Everstore | Distributed object store successor to Haystack — shows awareness of storage evolution |
| Iris | Meta's internal messaging platform backbone — shows messaging pipeline familiarity |
| Prophet | Open-source time-series forecasting library (from Meta) — shows ML/forecasting awareness |
| Tupperware | Meta's container orchestration platform (pre-dates Kubernetes ubiquity at Meta) — shows infra depth |
| EdgeRank / Feed ranking | Learned ranking model for news feed — shows ML integration awareness |
| Multifeed | Architecture where each user's feed is a materialized ranked list — shows feed design depth |

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| TAO follower diverges from leader | Stale read rate spikes; consistency checks fail | TAO refill from MySQL; read-repair on follower divergence |
| Celebrity post causes fanout storm | Write QPS spike on fan-out workers; queue lag grows | Switch celebrity accounts to pull-on-read above follower threshold (e.g., 10K) |
| Privacy change fails to propagate | Audit log gap; post visible to wrong audience | Synchronous write to privacy store before fan-out begins; block fan-out on privacy confirmation |
| Ranking model serves stale features | CTR drops suddenly; A/B holdout diverges | Feature store heartbeat checks; fallback to heuristic ranking |
| Graph shard hotspot | Latency spike on specific TAO shard | Consistent hashing with virtual nodes; read replicas for hot nodes |
| Feed cache eviction storm | Cache hit rate drops; origin QPS spikes | Pre-warm caches on deploy; use TTL jitter to prevent synchronized expiry |

## Observability

- metric: TAO cache hit rate by object type (user, post, association)
- metric: feed fanout latency p50/p95/p99 from write to cache population
- metric: privacy propagation lag (time from audience change to consistent read)
- metric: ranking inference latency at p99 by model version
- log: fan-out failures with affected user count and retry state
- trace: user request from post creation through fan-out to feed read
- alert: TAO follower divergence > 0.1% of reads within 5-minute window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Push-on-write for normal users, pull-on-read for celebrities | Low feed read latency for the common case | Complexity of hybrid fanout logic; need to track follower count thresholds | Pure pull-on-read is simple but too slow for accounts with few followers |
| TAO two-level cache over direct DB reads | Graph reads at <10ms p99 without hitting MySQL on every request | Cache invalidation complexity; eventual consistency for non-privacy-critical graph data | Direct MySQL reads would not scale to 50B+ read QPS on the social graph |
| Privacy enforcement at write time | Privacy bugs surface early; no audit gap | Higher write latency; fan-out blocked on privacy check | Enforcing at read time is simpler but creates windows where wrong audience sees content |
| Separate ranking service from feed retrieval | Ranking model can be updated and A/B tested independently | Extra network hop; need feature store synchronization | Inline ranking couples model changes to feed service deploys |

## Interview It

**Meta framing:** "Design Facebook's news feed." Strong answers decompose into at least: graph serving, fan-out pipeline, feed cache, ranking service, and privacy enforcement. Weak answers stay at the CDN + database level.

**Follow-ups:**
1. What happens to the news feed when the ranking service is down for 30 seconds?
2. A user with 50M followers posts for the first time in a year. Walk me through what happens in your system.
3. How do you enforce a privacy change (post becomes friends-only) so it is consistent on the next read?
4. How does your design support running ML ranking experiments without affecting the baseline feed?
5. What is your strategy for a graph traversal that requires 2nd-degree friends-of-friends at scale?

## Ship It

- `outputs/rubric-card-meta.md`
- `outputs/strong-hire-checklist-meta.md`
- `outputs/interview-card-meta-rubric.md`

## Exercises

1. **Easy** — List eight Meta-specific terms and explain the design trade-off each one represents.
2. **Medium** — Sketch a fanout architecture that handles both a normal user (200 friends) and a celebrity (50M followers) posting at the same time.
3. **Hard** — Write a capacity model for Meta's news feed pipeline: from post creation through fan-out to ranked feed delivery, with numbers at each stage.

## Further Reading

- [TAO: Facebook's Distributed Data Store for the Social Graph](https://www.usenix.org/conference/atc13/technical-sessions/presentation/bronson) — primary source
- [Finding a needle in Haystack: Facebook's photo storage](https://www.usenix.org/legacy/event/osdi10/tech/full_papers/Beaver.pdf) — Haystack architecture
- [Meta Engineering Blog](https://engineering.fb.com/) — primary source for Meta-specific architecture decisions
- [Scuba: Diving into Data at Facebook](https://research.facebook.com/publications/scuba-diving-into-data-at-facebook/) — observability stack
