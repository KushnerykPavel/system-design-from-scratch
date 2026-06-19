# Meta Full Mock Loop

> A 45-minute Meta system design interview is not a test of whether you know the answer. It is a test of how you think when you don't know the answer yet.

**Type:** Mock Interview  
**Company focus:** Meta  
**Learning goal:** Simulate a full 45-minute Meta-style system design interview. Practice the milestone cadence, Meta-specific technical vocabulary, common interviewer pivots, and the failure modes that separate strong hires from misses.  
**Prerequisites:** All lessons 01–06 in this phase  
**Estimated time:** ~90 min  
**Primary artifact:** completed mock design doc + self-assessment against rubric  

---

## Interview Prompt

**"Design Instagram's News Feed."**

*The interviewer gives you nothing more. The clarification round is yours to drive.*

---

## How to Use This Lesson

1. Set a 45-minute timer.
2. Work through the prompt independently using only a whiteboard or blank document.
3. After the timer ends, compare your design against the milestone map and failure mode list below.
4. Score yourself against the rubric.
5. Read the follow-up pivots and answer each one in writing (5 minutes each).

---

## Expected Strong-Hire Milestone Map

| Time mark | Expected progress | Failure mode if missed |
|-----------|-------------------|------------------------|
| **5 min** | Asked ≥3 clarifying questions covering: number of DAU, read/write ratio, latency target, whether the feed is ranked or chronological. Gave rough capacity numbers (500M DAU, ~5 opens/day, ~100 posts/write fanout). Scoped to one primary problem. | Jumped straight to boxes and arrows. Numbers never appeared. Asked one vague question ("what are the requirements?") and moved on. |
| **15 min** | Sketched high-level architecture: post write → fanout → feed store → feed read → ranking. Named the three hardest sub-problems: fanout at high-follower-count, ranking 10K candidates to 500, and serving at 500M DAU. Identified which sub-problem to go deep on. | Drew a generic "server + database + cache" diagram with no Meta-specific context. Did not identify hardest problems. |
| **25 min** | Drilled into one sub-system (fanout or ranking) with specifics: TAO for social graph, Scuba/Hive for engagement logs, Redis for feed cache, FAISS for candidate retrieval, three-stage ranking funnel. Discussed one failure mode with a specific mitigation. | Stayed at high level throughout. Mentioned "ML model" without explaining how it works or what features it uses. No failure modes. |
| **35 min** | Handled the interviewer's deliberate pivot without losing the thread. Adjusted design, showed awareness of the trade-offs the pivot introduces. | Got stuck on the pivot. Said "I haven't thought about that." Lost track of overall design while exploring the pivot. |
| **45 min** | Summarized: top 3 trade-off decisions made and why. Named the first metrics to instrument. Named one thing to do differently with 2× the time. | Ended mid-thought. Did not summarize. Did not reflect on trade-offs. |

---

## Clarifying Questions a Strong Candidate Asks First

1. "How many DAU are we designing for? I'll assume 500M unless you tell me otherwise."
2. "Is the feed ranked by relevance, or is it reverse-chronological? This changes the design significantly."
3. "What is the read/write ratio? Instagram skews heavily read — are we optimizing for read latency?"
4. "What is the latency target for loading the feed? Sub-200ms?"
5. "Are we handling celebrity accounts — users with millions of followers? That changes the fanout strategy entirely."

---

## Design Sketch: What to Cover

### Layer 1: Write Path — Post Creation and Fanout

```text
user creates post
  -> post service: stores post in TAO (social graph + content store)
  -> fanout service:
     if author has < 10,000 followers:  push fanout (write to each follower's feed cache)
     if author has ≥ 10,000 followers:  pull fanout (no write; followers fetch on read)
  -> for push fanout:
     -> read follower list from TAO (~100 followers for typical user)
     -> write post_id to each follower's Redis feed list (LPUSH, capped at 1000 entries)
  -> async: publish post event to Kafka for ranking signal ingestion
```

### Layer 2: Read Path — Feed Retrieval

```text
user opens feed
  -> feed service: fetch post IDs from Redis feed list (LRANGE 0 199)
  -> for celebrity posts (pull fanout): merge celebrity post IDs into feed list in real time
  -> fetch post content from TAO / Memcached for the merged ID list
  -> pass candidate list (~200–10,000 items) to ranking service
  -> ranking service: three-stage funnel (retrieval → lightweight score → heavy rank)
  -> return top ~50 items for first page
  -> client renders; next page fetches on scroll
```

### Layer 3: Ranking

```text
candidates (post IDs)
  -> lightweight scorer:
     -> feature fetch from Redis (user affinity, post engagement rate, recency)
     -> GBT model: score all candidates in <30ms
     -> prune to top ~500
  -> heavy ranker:
     -> full feature vector (100+ features) from feature store
     -> deep neural network inference on GPU: ~20ms for 500 items
  -> policy/diversity filter:
     -> dedup (already seen posts)
     -> creator capping (max 3 posts per creator per feed load)
     -> safety filters
  -> return top 50 for first page
```

### Layer 4: Storage

```text
TAO:          social graph + post metadata + user profiles (sharded MySQL + memcache)
Redis:        per-user feed list (post IDs), feature store (user + item features)
Everstore:    photos and videos (content-addressed object store)
Hive:         offline logs (engagement events, training data)
Kafka:        engagement events, fanout tasks, transcoding jobs
FAISS index:  post embeddings for ANN retrieval (updated daily)
```

### Layer 5: Reliability

```text
fanout failures:  async queue; retry with exponential backoff; DLQ after 3 failures
ranking failures: circuit breaker → fall back to lightweight scorer → fall back to recency sort
Redis failure:    fall back to pull-based feed from TAO (higher latency, but always available)
region failure:   active-active multi-region; TAO replicates cross-region; traffic reroutes via DNS
```

---

## Failure Modes to Demonstrate

Name at least three of these without prompting:

| Failure | What strong candidates say |
|---------|---------------------------|
| Celebrity posts not appearing in follower feeds | "We use pull fanout for high-follower accounts. On feed load, we merge celebrity post IDs in real time by fetching the celebrity's recent posts from TAO. The merge happens server-side before ranking." |
| Redis feed cache is unavailable | "We fall back to a direct TAO query for recent posts from followed accounts. Latency increases from <50ms to ~200ms, but feed availability is maintained." |
| Heavy ranking model is too slow | "Circuit breaker trips. Feed is served using lightweight scorer output. Users see a less personalized but still relevant feed within SLA." |
| Fanout queue backs up during traffic spike | "Fanout is async. The queue absorbs the spike. Feed reads for affected followers may see slightly stale data for a few seconds, which is acceptable." |
| New viral post not yet in rankings | "Engagement rate feature is updated by Flink within minutes. The ranking model incorporates it at the feature level. New viral posts naturally rise in rankings as their engagement signals propagate." |

---

## Interviewer Pivots (Follow-Up Challenges)

Practice answering each in under 5 minutes:

**Pivot 1: Stories (ephemeral content)**  
"Instagram also has Stories — content that disappears after 24 hours. How does the design change?"

*What to cover:* Stories have a TTL — use Redis with 24-hour expiry or a separate Ephemeral Post Service. Viewing Stories is a write (marks as seen), unlike feed which is read-only. Story fanout is simpler: push to all followers' story tray immediately (Stories have fewer items than feed). CDN caching must respect TTL to prevent serving expired Stories.

**Pivot 2: Ads in the feed**  
"Meta needs to insert ads into the feed at a ratio of 1 in 5 posts. How does this integrate?"

*What to cover:* Ads are not stored in the user's feed list — they are fetched from the Ad Serving System separately and injected at feed assembly time. The ad server returns a ranked list of ads based on targeting. The feed service merges organic posts and ads at the final assembly step, respecting the 1-in-5 ratio. Ads have different ranking signals (bid price, predicted CTR, ad quality score) from organic content.

**Pivot 3: Feed for a new user with no followers**  
"A new user just signed up and follows nobody. What do they see?"

*What to cover:* Cold start: onboarding flow asks for interests and suggests accounts to follow. Until they follow enough people, the feed is populated with: (1) trending content in their region, (2) content from suggested accounts in their interest category, (3) Reels from accounts with high engagement. The ranking model falls back to global popularity + interest-category affinity derived from onboarding.

**Pivot 4: Feed for a user who follows 10,000 accounts**  
"What if a power user follows 10,000 accounts? How does the read path change?"

*What to cover:* Pull fanout exclusively for all 10,000 followed accounts is expensive at read time (~100ms just to fetch post IDs). Optimization: pre-compute a merged feed for power users during off-peak periods and cache it with a 5-minute TTL. Alternatively, sample from the follow graph — don't fetch all 10,000 account feeds every time, instead use the ranking model to pre-select the 100 most relevant accounts to sample on each load.

**Pivot 5: Regulatory — right to explanation**  
"The EU DSA requires Meta to explain to each user why each post appears in their feed. How do you add this?"

*What to cover:* Log SHAP-style feature attributions for every ranking decision (sampled or full for EU users). Store attributions in a queryable audit log. Expose an endpoint: GET /feed/{post_id}/explanation → returns top-3 feature contributions in human-readable form. User-facing UI shows: "Shown because you follow [creator], and this post has high engagement in your network." Audit log retained for 90 days per DSA requirements.

---

## Scoring Rubric

| Dimension | Points | Strong hire | Hire | No hire |
|-----------|--------|-------------|------|---------|
| Clarification & scope | 10 | ≥3 targeted questions with numbers; scoped to one problem | 1–2 questions; partial numbers | No questions; jumped to architecture |
| Capacity estimation | 10 | Estimated DAU, RPS, storage growth, fanout ratio | Estimated 2–3 dimensions | No estimates provided |
| Architecture | 20 | All major layers named with correct data flow; Meta-specific components (TAO, Everstore, FAISS) | 3 of 4 layers; some Meta context | Generic diagram only; no Meta context |
| Deep-dive | 20 | Went deep on 1–2 subsystems with data model, failure mode, and specific component choices | Deep on 1 subsystem; no failure modes | Stayed high-level throughout |
| Failure modes | 15 | Named ≥3 failure modes with specific mitigations and fallback chains | Named 1–2 failure modes | No failure modes mentioned |
| Trade-offs | 15 | Named ≥2 trade-offs with rejected alternatives and quantified reasoning | Named 1 trade-off | No trade-offs mentioned |
| Observability | 10 | Named specific metrics (feed load latency, fanout lag, ranking p99, cache hit ratio) + alert thresholds | Named 2–3 generic metrics | No observability |

**Maximum score: 100 points**  
- 85–100: Strong hire  
- 65–84: Hire  
- 45–64: Mixed — needs another round  
- 0–44: No hire  

---

## Common Meta Interviewer Follow-Ups and Ideal Responses

**"How does TAO differ from a standard relational database?"**  
TAO is a distributed key-value and graph store optimized for read-heavy social graph queries. Unlike MySQL, it supports graph traversal (get all friends of a user, get all likes on a post) natively. It uses a two-tier cache (leader + follower) on top of MySQL, and achieves billions of reads per day with single-millisecond latency through Memcached integration.

**"Why is push fanout problematic for celebrities?"**  
A celebrity with 10M followers creates a fanout write amplification of 10M writes per post. At peak posting time, this creates a write storm that can overwhelm the Redis cluster. Pull fanout avoids this: the celebrity's post is not pushed; instead, followers fetch it at read time by querying the celebrity's post list. The trade-off is higher read latency on feed load.

**"What happens to the feed if Meta's main data center loses power?"**  
Active-active multi-region: each region is a full replica of the social graph (TAO), feed cache (Redis), and media store (Everstore). DNS health checks detect the region failure within 60 seconds and reroute traffic. Feed reads from the failover region may lag by 30–60 seconds due to replication lag, but the feed always loads.

**"How would you reduce feed load time from 200ms to 50ms?"**  
Three levers: (1) Pre-rank the feed offline and cache the result — user opens feed, it's already computed. (2) Reduce feature fetch latency by co-locating the feature store with the ranking service. (3) Reduce the heavy model inference time by using a quantized model or reducing the number of candidates it scores from 500 to 100 using the lightweight scorer more aggressively.

---

## Post-Mock Debrief Checklist

After your 45-minute mock, review each item:

- [ ] Did you clarify at least 3 requirements before drawing anything?
- [ ] Did you give capacity estimates (DAU, RPS, storage) within the first 10 minutes?
- [ ] Did you name Meta-specific components: TAO, Everstore, FAISS, Redis feed list, Kafka?
- [ ] Did you explain push vs pull fanout and when each is appropriate?
- [ ] Did you cover the three-stage ranking funnel?
- [ ] Did you name at least 3 failure modes with specific mitigations?
- [ ] Did you articulate at least 2 trade-offs with rejected alternatives?
- [ ] Did you handle the interviewer's pivot without losing the design thread?
- [ ] Did you name at least 3 metrics to instrument?
- [ ] Did you close with a 60-second summary of your top trade-off decisions?

---

## Weak-Hire Anti-Patterns to Avoid

1. **Generic "S3 + CDN + database" diagram.** This signals no Meta preparation. Use: TAO, Everstore, FAISS, Redis feed list, Kafka fanout queue.
2. **Ignoring the celebrity problem.** Every Meta interviewer will ask about it. If you don't mention push vs pull fanout unprompted, it signals you haven't studied the problem deeply.
3. **Saying "we can scale horizontally" without explaining how.** How do you shard TAO? How do you partition the Kafka fanout queue? What is the sharding key?
4. **No failure modes.** Meta's core value is reliability. Interviewers who live through incidents every month want to know you think about what breaks.
5. **Staying in the happy path.** Spending 40 minutes on the write path and never discussing ranking or read performance.
6. **No numbers.** "We'll have a lot of users" is not an answer. "500M DAU × 5 opens/day × 50 items/open = 125B feed items served/day = ~1.4M RPS" is.
7. **Getting derailed by the pivot.** The pivot is intentional. Have a mental model flexible enough to adapt in 5 minutes.

---

## Ship It

- `outputs/mock-design-doc-meta-news-feed.md`
- `outputs/self-assessment-rubric-meta-mock.md`
- `outputs/pivot-answers-meta-mock.md`

## Further Reading

- [TAO: Facebook's Distributed Data Store for the Social Graph (ATC 2013)](https://www.usenix.org/system/files/conference/atc13/atc13-bronson.pdf)  
- [Meta Engineering Blog — News Feed](https://engineering.fb.com/2021/09/13/production-engineering/disruptive-news-feed/)  
- [DLRM: A Deep Learning Recommendation Model for Personalization and Recommendation Systems](https://arxiv.org/abs/1906.00091)  
