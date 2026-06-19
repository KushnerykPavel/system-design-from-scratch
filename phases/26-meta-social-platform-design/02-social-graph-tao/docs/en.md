# Social Graph — TAO Cache & Graph Serving

> You can't query MySQL 50 billion times a second. TAO is why Meta can.

**Type:** Deep Dive
**Company focus:** Meta
**Learning goal:** Understand how Meta serves the social graph at 3B-user scale using the TAO two-level cache. Know the object/association data model, the read/write path, and the failure modes interviewers probe.
**Prerequisites:** `06-caching-and-invalidation/01-cache-aside-and-write-through`, `08-consistency-replication-and-transactions/02-read-your-writes`
**Estimated time:** ~90 min
**Primary artifact:** TAO architecture diagram + capacity model

## The Problem

Facebook's social graph has 3B user nodes and 100B+ edges (friendships, likes, follows, group memberships). Every news feed load, every "mutual friends" query, and every privacy check requires traversing this graph. The naive solution — reading from MySQL on every request — would require tens of billions of queries per second. TAO solves this by placing a purpose-built, graph-aware two-level cache in front of MySQL.

## Clarify

- What read patterns dominate? Mostly `(user, friends)` and `(user, likes on post)` — association lists sorted by recency.
- What consistency level is acceptable? Eventual consistency is fine for social graph reads (a stale friend list for a few seconds is acceptable). Privacy changes are the exception.
- How are writes handled? Writes go through leaders and propagate to followers; read-your-own-writes is not guaranteed by default.

## Requirements

### Functional

- Look up a user object (profile, metadata) by ID.
- Query associations for a given (id1, atype) pair, returning associated id2s sorted by time descending.
- Write new associations (friendship, like, follow) with single-writer guarantee per TAO leader.

### Non-functional

- Read latency: p99 < 10ms for graph lookups serving the news feed.
- Read throughput: ~50B read QPS at peak across the social graph.
- Write throughput: ~1M write QPS (new posts, likes, follows, unfriends).
- Availability: TAO followers can serve stale data when TAO leader is degraded.

## Capacity Model

| Dimension | Estimate | Notes |
|-----------|----------|-------|
| User nodes | 3B | Each object ~1KB → ~3TB of object data |
| Social graph edges | 100B+ | Each assoc ~100B → ~10TB of association data |
| Avg friends per user | 500 | Skewed: celebrities have millions of followers |
| Peak graph read QPS | ~50B | Spread across TAO follower fleet |
| Peak write QPS | ~1M | Spread across TAO leader fleet → MySQL shards |
| TAO follower cache hit rate | ~99% | Most social graph reads hit hot working set |
| MySQL shard count | ~1,000+ | Sharded by object/user ID |

## Architecture

### Data Model

TAO stores two primitive types:

**Objects:** Entities with a unique ID.
```
object(id uint64, otype string, data json)
```
Examples: user (otype=user), post (otype=post), photo (otype=photo).

**Associations:** Directed edges between objects.
```
assoc(id1 uint64, atype string, id2 uint64, time int64, data json)
```
Examples: (alice, friend, bob), (alice, liked, post_123), (alice, member_of, group_456).

Association queries are always by `(id1, atype)`, returning results sorted by `time DESC`. This access pattern drives the entire TAO cache design.

### Two-Level Cache Architecture

```
App servers
    │
    ▼
TAO Followers  (many, close to app servers, regional)
    │  cache miss
    ▼
TAO Leaders    (one per region, sharded by id1)
    │  cache miss
    ▼
MySQL Shards   (persistent store, sharded by object ID)
```

**Read path:**
1. App server queries TAO follower for `(id1, atype)`.
2. Follower returns cached result if present (hot path, sub-millisecond).
3. On follower miss: follower forwards to TAO leader.
4. On leader miss: leader fetches from MySQL shard and caches result.
5. Leader returns to follower; follower caches result and returns to app server.

**Write path:**
1. App server sends write to TAO leader (sharded by id1).
2. Leader writes to MySQL shard synchronously.
3. Leader invalidates/updates its own cache.
4. Leader sends invalidation messages to TAO followers.
5. Followers invalidate their caches; next read triggers a refill from leader.

### Consistency Model

TAO provides **eventual consistency** by design. Follower caches may serve stale association lists for a brief window after a write. This is explicitly acceptable for social graph reads: a stale friend list or like count for a few seconds does not harm the product.

**Exception — privacy changes:** Audience restriction changes (post becomes friends-only) are handled outside the standard TAO eventual-consistency path. Privacy checks are enforced synchronously before fan-out begins.

### Sharding Strategy

- Objects and associations are sharded by the first object ID (id1 for associations).
- All associations for a given user (id1) land on the same MySQL shard and the same TAO leader shard.
- This ensures association list queries are single-shard and the leader cache is cohesive per user.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| TAO follower diverges from leader | Stale read rate spikes; association count drift on consistency checks | Periodic cache refill from leader; read-repair on detected divergence |
| TAO leader is unreachable | Follower cache miss rate spikes; latency p99 rises | Followers can serve stale data up to a configured TTL; writes queue until leader recovers |
| MySQL shard hotspot | Leader QPS on specific shard spikes; latency rises | Read replicas for hot shards; query result caching at leader with longer TTL |
| Association count drift (cache vs DB mismatch) | Periodic consistency audits fail | Background reconciliation job; leader refetch on drift detected |
| Celebrity node fan-out | Write QPS to leader for high-follower id1 spikes | Separate write path for high-follower accounts; pull-on-read for association lists above threshold |

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Eventual consistency for social graph reads | TAO followers can serve stale data without blocking on leader, enabling massive read throughput | Brief inconsistency window after writes | Strong consistency would require all reads to hit the leader, eliminating the follower tier's value |
| Association queries always by (id1, atype) | Single-shard query, cache-friendly access pattern | Cannot efficiently query by id2 (e.g., "who liked this post" requires separate association type) | Bi-directional indexes double write cost and cache complexity |
| Invalidation via message to followers vs. TTL expiry | Near-real-time cache freshness after writes | Invalidation fan-out cost; messages can be lost, requiring TTL as a safety net | Pure TTL would allow longer staleness windows and more MySQL load |
| Separate object and association caches | Different TTL and eviction policies per primitive type | More cache management complexity | Unified cache would conflate hot objects (users) with hot associations (friend lists) |

## Observability

- metric: TAO follower cache hit rate by association type (atype)
- metric: TAO leader to MySQL query rate (cache miss rate at the leader)
- metric: TAO write propagation latency (leader write to follower invalidation)
- metric: MySQL shard QPS and p99 latency per shard
- log: follower refill events with id1, atype, and staleness duration
- trace: full read path from app server through follower, leader, to MySQL
- alert: follower divergence > 0.1% of reads within a 5-minute window

## Interview It

**Meta framing:** "How does Facebook serve the social graph at scale?" A strong answer names TAO by structure (two-level cache, objects + associations, eventual consistency) and names the access pattern that drives the design: `(id1, atype)` sorted by time. A weak answer says "use Redis in front of MySQL" without explaining why the association query pattern makes a generic cache insufficient.

**Follow-ups:**
1. What if a celebrity account has 100M followers? How does your design handle association list queries for that node?
2. How do you handle a privacy change (unfriend) that must be consistent on the next read?
3. A TAO leader shard is slow due to a hot MySQL shard beneath it. What is your mitigation?
4. How does TAO handle the association count for "number of likes on a post" — is the count stored separately?
5. What would you change if read-your-own-writes was a hard requirement?

## Ship It

- `outputs/tao-architecture-diagram.md`
- `outputs/capacity-model-tao.md`
- `outputs/tao-data-model-reference.md`

## Exercises

1. **Easy** — Draw the read path for a query: "who are Alice's friends?" Name every hop from app server to MySQL shard.
2. **Medium** — Design the write path for an "unfriend" operation. Identify where consistency is critical and where eventual consistency is acceptable.
3. **Hard** — A TAO follower cache has a 99% hit rate with 1B cached associations. Calculate the memory required and the MySQL QPS that the 1% miss rate generates at 50B total read QPS.

## Further Reading

- [TAO: Facebook's Distributed Data Store for the Social Graph (USENIX ATC 2013)](https://www.usenix.org/conference/atc13/technical-sessions/presentation/bronson) — definitive primary source
- [Scaling Memcache at Facebook (NSDI 2013)](https://www.usenix.org/conference/nsdi13/technical-sessions/presentation/nishtala) — caching architecture context
- [Meta Engineering Blog — Graph API](https://engineering.fb.com/) — updated architectural context
