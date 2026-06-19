# News Feed Fanout at 3 Billion Scale

> Push to active followers, pull for celebrities — hybrid wins at 3B scale.

**Type:** Build
**Company focus:** Meta
**Learning goal:** Design a hybrid push/pull news feed that handles celebrity accounts with 100M followers without a fanout storm.
**Prerequisites:** `16-application-backends/02-news-feed`, `06-caching-and-invalidation/04-cache-stampede`
**Estimated time:** ~90 min
**Primary artifact:** fanout strategy decision matrix

## The Problem

Facebook has 3 billion monthly active users, roughly 500 million daily active users, each with an average of 200 friends. If every user posts 10 times per day, that is 5 billion posts per day and — in a pure push model — 5 billion × 200 = 1 trillion fanout writes per day. Now add celebrities: Cristiano Ronaldo has over 100 million followers. A single post triggers 100 million writes. A naive push implementation would exhaust your write capacity in minutes.

Design a news feed system that delivers relevant, ranked content to each user within seconds of a post, without collapsing under celebrity traffic.

## Clarify

- What is the follower threshold that separates "normal" from "celebrity" accounts?
- Does feed ranking happen at write time (precomputed score) or read time (online inference)?
- What freshness SLA is required? (seconds vs minutes vs best effort)
- How important is consistency — can a user miss a post from a famous account for up to 60 seconds?
- Must GDPR delete or content moderation reach all materialized feeds promptly?

## Requirements

### Functional

- Each user sees a ranked news feed when they open the app.
- A new post from a followed account appears in the follower's feed within 30 seconds.
- Feed supports infinite scroll via cursor pagination.
- Deleting a post must propagate to visible feeds within 60 seconds.

### Non-functional

- Feed read latency: p99 under 100ms (Memcached-backed).
- Fanout write throughput: sustain 1 trillion fanout operations per day without a celebrity post causing a spike that degrades reads.
- Availability: 99.99% for feed reads; graceful degradation to chronological order if ranking is unavailable.
- Storage: feed store holds last 500 items per user; older items are re-hydrated from the post store on scroll.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| DAU | 500M | number of feed caches to maintain |
| Avg friends per user | 200 | fanout multiplier for push path |
| Posts per user per day | 10 | total posts: 5B/day |
| Fanout ops per day (pure push) | 1T | reveals write amplification problem |
| Feed items cached per user | 500 | Memcached entry size per user |
| Celebrity follower ceiling | 100M | worst-case push cost per post |
| Facebook celebrity threshold | ~50K followers | crossover point for hybrid routing |

## Architecture

```text
post writer
  -> post store (MySQL + TAO graph)
  -> fanout decision service
       |-- normal account (<50K followers)
       |     -> Kafka topic: fanout-jobs
       |          -> feed workers (parallel)
       |               -> Memcached: prepend post_id to user's feed list
       |
       `-- celebrity account (>=50K followers)
             -> celebrity post store (indexed by author_id + timestamp)

feed read path
  -> load user's Memcached feed list (push-populated items)
  -> identify followed celebrities from TAO graph (cached ~5 min)
  -> fetch recent posts from celebrity post store for each celebrity
  -> merge + rank: ML ranker assigns relevance score per item
  -> return top-N items sorted by score
```

### Fanout Strategy Decision Matrix

| Account type | Follower count | Write path | Read path |
|---|---|---|---|
| Normal user | < 50K | Push to each follower's MC feed list | Read directly from MC |
| Semi-celebrity | 50K–500K | Push to online followers only; skip offline | Read MC; merge missing items on first open |
| Full celebrity | > 500K | No push fanout | Pull from celebrity post store at read time |

Facebook reportedly uses ~50,000 followers as the push/pull threshold. Twitter/X uses a similar threshold (~100K followers).

### Feed Store Design

Each user's feed in Memcached is a sorted list of `(post_id, score, timestamp)` tuples. The feed worker prepends new items atomically using MC's `cas` (check-and-set) operation to avoid lost updates under concurrent writes.

```
key:   feed:{user_id}
value: [(post_id_1, score_1, ts_1), (post_id_2, score_2, ts_2), ...]
TTL:   7 days (evict inactive users)
max:   500 items; older items dropped on overflow
```

### Ranking: From EdgeRank to ML

Facebook's original EdgeRank formula weighted items by `affinity_score × edge_weight × time_decay`. The current ML-based ranker (deployed around 2011 and continuously evolved since) treats ranking as a prediction problem:

```
score(post, viewer) = P(engagement | features)

Features:
  - affinity_score: how often viewer interacts with author
  - post_type: text, photo, video, link
  - recency: exponential decay with half-life ~6h
  - social proof: likes/comments from viewer's friends
  - predicted view time (for video content)
```

At read time, the ranker scores each candidate item and returns the top-K. The model is served from a low-latency inference cluster (p99 under 10ms).

### Celebrity Post Store

For celebrity posts, a separate store (think: a wide-column store like HBase or a dedicated Cassandra table) holds:

```
partition key: author_id
sort key:      post_id DESC (newest first)
columns:       content_id, timestamp, type, engagement_snapshot
```

At read time, the feed service queries each celebrity the viewer follows. To bound latency, the viewer's celebrity follow list is cached in TAO and the maximum number of celebrity queries per feed load is capped at 50. Queries run in parallel.

### Fanout Worker Pipeline

```text
Kafka topic: fanout-jobs
  partition key: post_id % N (spread across partitions)

Consumer group: feed-workers (auto-scaled)
  for each job:
    1. Fetch follower list from TAO in pages of 1000
    2. Filter: skip offline users (last_active > 48h)
    3. Write to MC: PREPEND post_id to feed:{follower_id}
    4. Trim feed to max 500 items
    5. Ack Kafka offset
```

Skipping offline users reduces fanout by 30-50% for normal accounts and avoids filling the MC feeds of users who will never read them. When a dormant user reactivates, their feed is cold-built by pulling the most recent posts from followed accounts.

## Data Model

```sql
-- Post store (simplified)
post(post_id BIGINT PK, author_id BIGINT, content TEXT, 
     created_at TIMESTAMP, deleted_at TIMESTAMP)

-- Feed item (materialized in MC, not a DB table in hot path)
-- feed_item(user_id, post_id, score FLOAT, timestamp)

-- Celebrity follows (cached from TAO)
celebrity_follow(follower_id BIGINT, celebrity_id BIGINT, 
                 followed_at TIMESTAMP)

-- Celebrity post index (Cassandra)
celebrity_post(author_id BIGINT, post_id BIGINT DESC, 
               content_id TEXT, ts TIMESTAMP)
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Fanout worker lag spike | Kafka consumer lag metric > threshold | Auto-scale workers; shed load by skipping offline users aggressively |
| MC eviction during celebrity viral post | Cache miss rate spikes on feed reads | Pull directly from celebrity post store; rate-limit MC writes during storm |
| ML ranker serving latency spike | p99 ranker latency > 50ms | Fall back to recency-sorted feed; degrade gracefully rather than block reads |
| Post deleted but still in MC feeds | Moderation alert; user complaint | Tombstone record blocks deleted post_ids at read time; async cleanup |
| GDPR delete request | Delete request received | Propagate tombstone immediately; async removal from all MC entries and post stores |
| TAO graph partition | Celebrity follow list unavailable | Use stale cached celebrity list with TTL extension; fail open with empty celebrity section |

## Observability

- metric: fanout_worker_lag_seconds (Kafka consumer lag)
- metric: feed_read_latency_ms p50/p99 (end-to-end)
- metric: mc_miss_rate per feed type (push vs celebrity pull)
- metric: fanout_ops_per_second by account tier (normal/semi/celebrity)
- metric: ranker_inference_latency_ms p99
- metric: posts_missing_from_feed_rate (freshness signal)
- log: per-post fanout start/complete/skipped with follower count and tier
- trace: post_created -> fanout_job_enqueued -> mc_write_complete latency

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Hybrid push/pull at 50K threshold | No fanout storm for celebrities; fast reads for normals | Extra complexity; two code paths; threshold tuning needed | Pure push collapses on celebrity posts; pure pull is too slow for normals |
| Skip offline users in fanout | Reduces fanout by 30-50%; saves MC writes | Dormant users see stale feed on reactivation; cold-build cost on first open | Fanout to everyone is wasteful and causes lag spikes |
| Prepend to MC feed list | O(1) write; constant-time feed prepend | MC CAS contention under high fan-in (viral post to active user) | Rebuild from DB on every read is too slow at scale |
| Pull celebrities at read time | No write amplification for celebrity posts | Latency at read time; requires parallel fan-in from celebrity store | Push to 100M followers per post is catastrophically expensive |
| ML ranking at read time | Personalized; adapts to signals not available at write time | Inference latency in hot read path; model freshness requirements | Pre-ranked feed loses personalization signals known only at view time |

## Interview It

**Meta framing:** "Design Facebook's news feed for 3 billion users." Strong answers identify the celebrity fanout problem early, propose a hybrid threshold, explain the MC feed store, and cover ranking degradation as a failure mode. Weak answers propose pure push without mentioning write amplification.

**Follow-ups:**

1. A celebrity post goes viral (10M reactions in 1 hour). How does your system respond without melting the fanout workers?
2. A user requests GDPR deletion. How do you propagate the removal from all materialized feeds within 72 hours?
3. Your ML ranker degrades and returns scores of 0 for all posts. What does the user see?
4. How would you implement "close friends" feed — a subset of the main feed shown in a separate tab?
5. A new user has no feed history. How do you cold-start their feed?

## Ship It

- `outputs/fanout-strategy-decision-matrix.md`
- `outputs/design-doc-news-feed-fanout.md`
- `outputs/interview-card-news-feed-fanout.md`

## Exercises

1. **Easy** — Calculate the total Memcached storage required if each feed entry is 24 bytes (post_id + score + timestamp) and you store 500 entries for all 500M DAU.
2. **Medium** — Design the cold-start feed for a user who just joined and has no posts from their follows yet.
3. **Hard** — Extend the fanout system to support "stories" (24-hour ephemeral content) alongside the regular feed, sharing the same worker infrastructure.

## Further Reading

- [Facebook TAO: Facebook's Distributed Data Store for the Social Graph](https://www.usenix.org/conference/atc13/technical-sessions/presentation/bronson)
- [Feeding Frenzy: Selectively Materializing Users' Event Feeds (ACM, 2010)](https://dl.acm.org/doi/10.1145/1807167.1807207)
- [How News Feed Works (Meta Engineering)](https://engineering.fb.com/2021/01/26/ml-applications/news-feed/)
