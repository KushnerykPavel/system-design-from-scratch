# News Feed / Timeline

> Feed systems are ranking and fanout systems disguised as simple post retrieval.

**Type:** Build  
**Company focus:** Google  
**Learning goal:** Design a home-timeline backend by reasoning about fanout strategy, celebrity skew, ranking freshness, and recovery when the precomputed feed gets behind.  
**Prerequisites:** `14-rate-limiters-ids-and-hashing/06-hot-key-mitigation`, `15-kv-cache-and-object-storage/01-distributed-kv-store`, `16-application-backends/01-url-shortener`  
**Estimated time:** ~90 min  
**Primary artifact:** feed strategy planner  

## The Problem

Design the backend for a social feed where users follow accounts and open a home timeline that mixes fresh content with relevance ranking. A naive answer says "store posts in a database and query followers." A strong answer distinguishes write amplification from read amplification and explains when to precompute timelines versus assemble them at read time.

This lesson is senior-level because the best design depends on follower-graph skew, ranking complexity, freshness goals, and what happens when celebrity accounts dominate the workload.

## Clarify

- Is the feed reverse-chronological, ranked, or a hybrid?
- Do we need hard freshness for new posts, or can ranking lag by seconds or minutes?
- How skewed is the follower graph, especially for celebrity or brand accounts?

If unspecified, assume a hybrid feed: fresh content should appear within seconds, ranking can refine within minutes, and a small number of accounts have very large fanout.

## Requirements

### Functional

- Users can publish posts and see them appear in follower feeds.
- Users can read a paginated home timeline.
- The system can blend ranking features, freshness, and policy filtering.
- Deleted or moderated posts should disappear from feeds quickly enough for product safety.

### Non-functional

- Keep timeline load p99 under 200 ms for ordinary users.
- Avoid exploding write cost for celebrity accounts.
- Support partial feed availability during ranking or fanout degradation.
- Make freshness lag measurable and explainable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Posts created | 120K posts/s peak | drives fanout work and queue sizing |
| Timeline reads | 18M requests/min peak | dominates storage access and cache design |
| Average follows | 400 per user, heavy long tail | makes one-size-fits-all fanout wrong |
| Celebrity fanout | top 0.1% accounts have 5M+ followers | forces mixed fanout strategy |
| Feed retention | 7 to 30 days hot timeline materialized | shapes storage and compaction choices |

## Architecture

```text
post write
  -> write API
  -> post store
  -> fanout classifier
     -> push fanout queue for normal accounts
     -> pull-on-read index for celebrity accounts
  -> ranking feature pipeline

timeline read
  -> feed service
  -> materialized timeline store
  -> on-read merge for celebrity or cold content
  -> ranking layer
  -> cache
```

The key design move is mixed fanout:

1. Fanout-on-write for ordinary accounts keeps read latency low.
2. Fanout-on-read for celebrity accounts prevents write explosions.
3. Read-time merge is the tax you pay for avoiding impossible write amplification.

## Data Model & APIs

Core entities:

```text
post(post_id, author_id, created_at, visibility, feature_refs)
follow(follower_id, followed_id)
timeline_entry(user_id, post_id, source, rank_hint, inserted_at)
```

APIs:

- `POST /v1/posts`
- `GET /v1/timeline?cursor=...`
- `POST /v1/follows`
- `POST /v1/posts/{id}/moderation-hide`

The timeline API should support cursors instead of offsets because inserts and ranking changes make offsets unstable.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| celebrity post floods fanout workers | queue depth by author cohort and worker saturation | classify to pull-on-read, cap synchronous fanout work |
| ranking feature pipeline lags | freshness gap and ranking feature age | serve chronological fallback with degraded ranking |
| moderation delete misses materialized entries | delete propagation lag and stale-content reports | tombstones plus background scrub of timeline stores |
| timeline cache stampedes on cold accounts | miss burst metrics and downstream saturation | request coalescing and partial prewarm |

## Observability

- metric: timeline read latency split by materialized-only versus merge-on-read
- metric: fanout queue depth by author cohort size
- metric: time-to-feed-visible for newly created posts
- metric: ranking feature age and stale-score fallback rate
- log: moderation delete propagation for sampled posts
- trace: post publish to first-feed-visibility path
- SLO: newly created posts from ordinary accounts become timeline-visible within the product freshness target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| mixed push and pull fanout | matches workload skew and caps celebrity write cost | more moving parts and read-time merge logic | pure push for everyone |
| cursor-based reads | stable pagination under inserts and rank updates | more stateful clients | offset pagination |
| degraded chronological fallback | preserves availability when ranking breaks | lower relevance quality | hard dependency on ranking for feed reads |

## Interview It

**Google framing:** "Design the home timeline for a social product with relevance ranking." Expect follow-ups on skew, cache strategy, and what timelines promise under ranking degradation.

**Cloudflare framing:** "Design a globally distributed feed read path with strong cache locality." Expect questions on edge caching limits, invalidation, and which computations stay regional.

**Follow-ups:**
1. What changes if every feed must be strictly chronological?
2. How do you hide a deleted or policy-violating post fast enough?
3. What if ranking models need features that arrive minutes late?
4. How do you keep one celebrity from consuming half the fleet?
5. What changes if the product adds ads or sponsored placement?

## Ship It

- `outputs/tradeoff-matrix-news-feed.md`

## Exercises

1. **Easy** — Choose a follower-count threshold for push versus pull fanout and justify it.
2. **Medium** — Design a fallback when the materialized timeline store is healthy but ranking is down.
3. **Hard** — Redesign the feed for a product where 60% of reads happen from one geographic region.

## Further Reading

- [Twitter Timelines at Scale](https://www.infoq.com/presentations/Twitter-Timeline-Scalability/) — classic fanout trade-offs and skew handling  
- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — useful for thinking about graceful degradation and freshness SLOs  
