# Feed Ranking & Content Distribution

> LinkedIn feed is not a social feed — it rewards professionally relevant content, not just popular content.

**Type:** Build
**Company focus:** LinkedIn
**Learning goal:** Design LinkedIn's feed ranking system that surfaces professionally relevant content to 30M daily active users while managing content distribution fairness for creators.
**Prerequisites:** `16-application-backends/02-news-feed`, `26-meta-social-platform-design/06-ranking-ml-serving`
**Estimated time:** ~90 min
**Primary artifact:** feed ranking pipeline design + relevance signal taxonomy

## The Problem

Design LinkedIn's feed ranking system. The feed must:

1. **Surface professionally relevant content** — career-relevant posts rank above entertainment; your manager's job change ranks above a meme.
2. **Balance creator distribution** — ensure diverse creator voices appear; prevent one highly-engaging creator from dominating every member's feed.
3. **Suppress misinformation** — reduce organic reach of posts that spread too fast (anti-virality for potentially low-quality viral content).
4. **Personalize at scale** — 30M daily active users, each receiving a unique ranked feed from ~1B daily content events.
5. **Support A/B testing** — rank new ranking model versions against a baseline across 5% of members continuously.

## Clarify

- What content types are in scope? Posts, articles, job alerts, connection milestones, company updates, sponsored content?
- What is the target latency for feed load? (Target: <200ms p99 for ranked feed)
- What is the freshness requirement? (Members should see content posted within the last 4 hours as a priority)
- Is the feed personalized per member or segment-level?
- Are creator monetization features (newsletter, subscription) in scope?

## Requirements

### Functional Requirements

- Rank feed candidates for a member, returning the top-N posts in personalized order.
- Filter spam and low-quality content before ranking.
- Apply business rules: ad frequency (1 sponsored post per 5 organic), creator frequency cap (no single creator appearing more than 3× in the top 20).
- Serve Venice-computed features for ranking model inference.
- Support experiment framework: assign members to ranking model variants for A/B testing.

### Non-functional Requirements

- 30M daily active users; assume 50% open their feed at least once daily.
- ~1B feed update events per day (posts, likes, shares, connection milestones).
- Feed ranking latency: <200ms p99 from request to ranked result returned to client.
- Feed freshness: content from the last 4 hours weighted up regardless of engagement.
- Ranking model update: new model deployed without feed service restart; A/B test coverage selectable by percentage.

## Capacity Model

| Dimension | Estimate | Detail |
|-----------|----------|--------|
| Daily active users | 30M | feeds requested per day |
| Feed requests/sec | ~350/s avg, ~3,000/s peak (9–10 AM EST) | daily peak driven by commute window |
| Candidate posts per feed request | ~2,000 | from 1st + 2nd degree connections + followed companies |
| Fast ranker input | ~2,000 candidates | lightweight feature model |
| Slow ranker input | ~500 candidates | heavy ML model (Venice features) |
| Final feed size | ~20 items per page | paginated |
| Venice feature lookups per request | ~500 (member + item features) | precomputed offline features |

## Architecture

### LinkedIn Feed vs Social Feeds

| Dimension | LinkedIn | Facebook/Instagram/Twitter |
|-----------|----------|---------------------------|
| Primary ranking signal | Professional relevance (job title match, career stage alignment) | Engagement rate (likes, comments, shares) |
| Viral coefficient | Intentionally suppressed — anti-virality scoring reduces fast-spreading posts | Amplified — viral posts get boosted distribution |
| Content type priority | Career milestones, industry insights, company news | Entertainment, personal updates, trending topics |
| Creator economy | Organic reach algorithm, newsletter features | Paid boosts, sponsored, creator fund |
| Spam definition | Engagement bait, clickbait job posts, fake endorsements | Misinformation, coordinated inauthentic behavior |

### Multi-Stage Ranking Pipeline

```
Feed Request
    ↓
[1. Candidate Generation]
  — fetch recent posts from 1st-degree connections (Kafka-fed cache)
  — fetch posts from followed companies and schools
  — fetch connection milestone events (X started a new role at Y)
  — total: ~2,000 candidates per member
    ↓
[2. Early-Stage Filter]
  — spam score threshold (remove SpamScore > 0.7)
  — blocklisted creators (member-level blocks, content policy blocks)
  — content policy filter (hate speech, explicit content)
  — survivors: ~1,500 candidates
    ↓
[3. Fast Ranker — lightweight model]
  — features: recency, engagement rate, creator_score (precomputed)
  — model: gradient boosted tree, <5ms inference
  — output: top-500 candidates with preliminary scores
    ↓
[4. Slow Ranker — heavy ML model]
  — fetch Venice features: member profile embedding, item content embedding, network proximity
  — model: deep neural network, ~50ms inference for 500 candidates
  — output: scored ranked list of 500 candidates
    ↓
[5. Diversity & Business Rules]
  — creator frequency cap: max 3× same creator in top 20
  — ad injection: 1 sponsored post per 5 organic (position 5, 10, 15, ...)
  — freshness boost: posts < 4h old get +0.1 score bonus
  — anti-virality: posts with viral_coefficient > 2.5 get −0.2 score penalty
    ↓
[6. Final Feed]
  — top-20 items for page 1
  — cached per member with 5-min TTL
```

### Venice Feature Store

Venice serves precomputed features to the slow ranker:

**Member features** (computed offline, updated daily):
- `job_title_embedding` — 64-dim vector representing member's current role
- `industry_category` — normalized industry (Finance, Engineering, Healthcare, ...)
- `engagement_history` — distribution of content types engaged with in last 30 days
- `network_influence_score` — normalized degree centrality in professional graph

**Item features** (computed on post creation + updated on engagement events, lag ~5 min):
- `content_category` — topic classifier output (Leadership, Tech, Marketing, ...)
- `creator_score` — creator's historical engagement rate normalized 0–1
- `spam_score` — spam classifier output 0–1
- `viral_coefficient` — posts_seen_to_shares ratio; high value = spreading unusually fast

**Cross features** (computed offline, lag ~1 hour):
- `member_item_relevance` — dot product of job_title_embedding × content_category_embedding
- `network_proximity_engagement` — fraction of member's 1st-degree connections who engaged with this post

### Relevance Signals Taxonomy

| Signal | Weight | Rationale |
|--------|--------|-----------|
| Dwell time (implicit) | High | Reading an article for 2+ minutes is a stronger engagement signal than a like |
| 1st-degree connection authored | High | Content from direct connections has the highest network proximity signal |
| Professional relevance match | High | job_title_embedding similarity to post content_category |
| Engagement by coworkers or same-industry peers | Medium | Social proof from professionally-aligned audience |
| Post freshness (<4h) | Medium | LinkedIn feed is semi-real-time; stale content should not dominate |
| Explicit positive engagement (like, comment, share) | Medium | Explicit signals, but weaker than dwell at LinkedIn |
| Creator score | Medium | Historically engaging creators likely to produce relevant content |
| Viral coefficient (penalty) | Negative | Posts spreading unusually fast may be engagement bait or misinformation |
| Spam score | Strong negative | Spam classifier output; high score → removed at filter stage |

### A/B Testing Framework

LinkedIn uses a member-level experiment assignment:

```
member_id → hash(member_id + experiment_id) % 100 → experiment bucket

Bucket 0–94: control (existing ranking model)
Bucket 95–99: treatment (new ranking model, 5% of members)
```

Metrics tracked per experiment bucket:
- 7-day feed CTR (click-through rate on posts)
- 7-day dwell time per session
- Connection request rate (proxy for feed-driven network growth)
- Creator reach (how many distinct creators appeared in feeds)
- Spam report rate

Experiment decision: if treatment CTR > control CTR by >2% with p-value < 0.05 after 7 days, promote new model to 100%.

### Anti-Virality Design

LinkedIn intentionally suppresses viral amplification more aggressively than other social networks:

**Viral coefficient** = (total views) / (organic shares in last 1 hour)

A viral_coefficient > 2.5 means a post is being reshared much faster than average. This could indicate:
- Genuine breakout professional insight (acceptable)
- Engagement bait ("Like if you agree...")
- Coordinated resharing campaigns
- Misinformation spreading

**Mechanism**: posts with viral_coefficient > 2.5 receive a −0.2 score penalty in the diversity rules stage. This does not remove the post; it deprioritizes it. A human content review flag is also set for posts with viral_coefficient > 5.0.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Slow ranker (ML model) latency spike | Feed latency p99 > 200ms; slow ranker timeout rate increases | Fallback to fast ranker output; serve fast-ranked top-20 without Venice features |
| Venice feature staleness | Ranking model using stale member or item features; item features lag > 15 min | Feature freshness SLA metrics per feature group; alert on item feature lag > 10 min |
| Spam filter false positive rate spike | Creator reach drops; high-quality posts incorrectly filtered | Shadow mode for new spam model versions; human review queue for disputed removals |
| Feed cache eviction storm | Feed latency spikes as cache misses hit Venice + slow ranker | TTL jitter (randomize TTL ±30s to spread expiry); pre-warm cache on feed service deploy |
| Creator frequency cap causes same-creator posts to disappear | Creator files complaint; organic reach metrics show sharp drop | Review cap threshold; A/B test cap at 2× vs 3× vs uncapped |
| A/B experiment contamination | Members in control group see treatment model behavior | Consistent hash on member_id + experiment_id; validate assignment is stable per member |

## Observability

- metric: feed ranking pipeline latency p50/p95/p99 broken down by stage (candidate generation, fast ranker, slow ranker, diversity rules)
- metric: Venice feature freshness per feature group (item features should be <10 min stale)
- metric: feed cache hit rate (high hit rate = lower ranking compute per member)
- metric: spam filter removal rate (sudden spike = new spam campaign or filter regression)
- metric: creator reach per feed session (diversity health indicator)
- metric: A/B experiment CTR and dwell time by bucket
- log: per-request ranking trace with top-5 scores and feature values (sampled at 0.1%)
- alert: slow ranker p99 > 150ms for 3 consecutive minutes (circuit breaker threshold)
- alert: Venice item feature lag > 10 minutes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Multi-stage ranking (fast + slow ranker) over single-stage ML | Fast ranker prunes 75% of candidates in <5ms, letting slow ranker focus on quality top-500 | Additional pipeline complexity; two model versions to maintain | Single-stage heavy ML on 2,000 candidates would take >200ms, exceeding feed latency SLA |
| Anti-virality scoring over pure engagement ranking | Suppresses misinformation and engagement bait; maintains content quality | Occasional suppression of genuinely viral quality content | Pure engagement ranking rewards optimizing for likes, driving engagement bait epidemic |
| Venice offline features over online feature computation | Microsecond feature serving; no compute at ranking time | Feature freshness lag (member features daily, item features ~5 min) | Online feature computation adds 100ms+ to ranking latency, breaking the 200ms SLA |
| 5-min feed cache TTL over no caching | Reduces slow ranker calls from 350/s to ~70/s (cache hit rate ~80%) | Members see content up to 5 min stale during active scroll | No cache: Venice + slow ranker at full QPS would require 5× infrastructure |

## Interview It

**LinkedIn framing:** "Design LinkedIn's feed ranking system." Strong answers decompose into: candidate generation (Kafka-fed cache), multi-stage ranking (fast + slow), Venice feature serving, business rules (ad injection, creator cap, anti-virality), and A/B testing framework.

**Follow-ups:**

1. The slow ranker ML model is down for 30 seconds. What does your feed serve?
2. A creator's post goes viral with 500K reshares in 1 hour. Walk me through what your system does.
3. How do you detect that a new ranking model is causing a 10% drop in creator reach before it ships to 100%?
4. A member complains that their feed shows only content from 3 creators. What metric would you check first?
5. How do you ensure the A/B test assignment is stable (the same member always gets the same model version)?

## Ship It

- `outputs/feed-ranking-pipeline-linkedin.md`
- `outputs/relevance-signal-taxonomy-linkedin.md`
- `outputs/ab-testing-feed-design.md`

## Exercises

1. **Easy** — List five signals that LinkedIn's feed ranking treats as positive and three signals that trigger a penalty or filter.
2. **Medium** — Design the anti-virality scoring mechanism: what is the viral coefficient formula, at what threshold does the penalty apply, and how does it interact with the multi-stage pipeline?
3. **Hard** — Write a capacity model for the full feed ranking pipeline: from Kafka feed candidate ingestion through Venice feature fetch to ranked result, with numbers at each stage.

## Further Reading

- [LinkedIn Feed Ranking Architecture](https://engineering.linkedin.com/blog/2022/feed-rankings-at-linkedin) — primary source
- [Venice: LinkedIn's Feature Store](https://engineering.linkedin.com/blog/2021/the-magic-of-venice) — feature store serving
- [LinkedIn's Anti-Virality Approach](https://engineering.linkedin.com/blog/2022/how-we-fight-misinformation) — content quality design
- [A/B Testing at LinkedIn](https://engineering.linkedin.com/blog/2020/a-b-testing-platform) — experiment framework
- [Pinot: Real-Time OLAP for Feed Analytics](https://engineering.linkedin.com/pinot) — analytics pipeline
