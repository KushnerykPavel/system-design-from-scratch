# Ranking & ML Feature Serving at Feed Scale

> The feed ranking system is not an algorithm. It is a manufacturing line — thousands of candidate posts reduced to hundreds through successive quality filters, each operating under a strict time budget.

**Type:** Build  
**Company focus:** Meta  
**Learning goal:** Design Meta's feed ranking pipeline covering multi-stage retrieval and scoring, feature store architecture, model serving under 20ms latency budgets, feedback loops, and failure handling at 500M DAU scale.  
**Prerequisites:** `03-news-feed-fanout`, `05-media-pipeline`  
**Estimated time:** ~90 min  
**Primary artifact:** ranking pipeline design doc + feature store spec  

## The Problem

At 500M DAU, each active user opens their feed multiple times per day. The raw candidate pool for a single feed session is 10,000+ posts (from follows, groups, pages, recommendations). Meta must reduce this to ~500 visible items per session while maximizing engagement and relevance — and do it in under 200ms end to end.

Design the ranking system that powers News Feed and Reels ranking at 500M DAU.

## Clarify

- What is the end-to-end latency budget for ranking? (<200ms for the full pipeline; ranking stage has ~20ms)
- How many candidates enter the ranking pipeline per request? (~10,000 retrieved; scored to ~500 shown)
- What signals are available? (explicit: likes, shares; implicit: time spent, scrolling past; context: device, time of day)
- How stale can features be? (user short-term features: minutes; item engagement rate: hours; user long-term: daily)
- Are there regulatory constraints? (Yes — EU DSA requires explainability for recommendation decisions)
- How do we handle new users with no history (cold start)?

## Requirements

### Functional

- Retrieve ~10,000 candidates from follow graph and recommendation sources.
- Score candidates through a lightweight model and a heavy ranking model.
- Apply policy and diversity filters before returning ~500 items.
- Serve user and item features from a feature store at <5ms.
- Process engagement events to update features and retrain models.
- Provide audit logs for DSA explainability compliance.

### Non-functional

- Ranking pipeline latency: p99 <200ms end to end; ranking stage <20ms.
- Feature store read latency: p99 <5ms.
- Model freshness: short-term signals updated within minutes; models retrained daily.
- Scale: 500M DAU × ~5 feed opens/day = 2.5B ranking requests/day (~29,000 RPS).
- Availability: degradable — fall back to simpler heuristic ranking on model failure.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| DAU | 500M | Total ranking request volume |
| Feed opens per user per day | ~5 | Ranking RPS sizing |
| Ranking RPS | ~29,000 | Model serving cluster size |
| Candidates per request | ~10,000 | Lightweight scorer throughput |
| Items scored by heavy model | ~500 | Heavy model serving budget |
| Feature store reads per request | ~3 (user features, item features, context) | Feature store QPS |
| Feature store QPS | ~90,000 | Redis/Memcached cluster sizing |
| Engagement events per day | ~5B (clicks, likes, shares, scrolls) | Kafka throughput |

## Architecture

```text
[ranking request]
  feed load request (user_id, session_context)
  -> candidate retrieval service
     -> follow graph candidates (posts from followed accounts, last 72h)
     -> recommendation candidates (FAISS ANN from interest embeddings)
     -> union: ~10,000 candidates
  -> lightweight scorer
     -> batch feature fetch from feature store (Redis)
     -> fast model (GBT or small MLP, <1ms per item)
     -> prune to top ~500 by lightweight score
  -> heavy ranker
     -> full feature fetch (user + item + context)
     -> large ranking model (deep neural network, ~20ms total for 500 items)
     -> produce ranked list
  -> policy & diversity filter
     -> dedup (remove items user already saw)
     -> diversity (cap same creator, same topic)
     -> safety (remove policy-violating content)
  -> return top ~500 ordered items

[feature pipeline]
  engagement event (like, share, time-spent, scroll-past)
  -> Kafka topic: user-engagement-events
  -> stream processor (Flink): update short-term user features
  -> feature store (Redis): user short-term features updated in <2min
  -> batch pipeline (Spark daily): update long-term user features, item statistics
  -> model retraining (PyTorch on GPU cluster): daily full retrain
  -> model registry: versioned model artifacts
  -> model serving (TorchServe): A/B tested rollout
```

### Three-Stage Ranking

Multi-stage ranking trades off accuracy for latency at each stage:

```text
Stage 1: Candidate Retrieval (latency budget ~10ms)
  input:  follow graph + interest graph
  method: graph traversal + FAISS ANN on interest embeddings
  output: ~10,000 candidates
  model:  no ML model — heuristic graph traversal + ANN index

Stage 2: Lightweight Scoring (latency budget ~30ms for all 10,000)
  input:  10,000 candidates
  method: fast GBT model with 50 features per candidate
  output: top ~500 candidates sorted by lightweight score
  feature fetch: batched Redis lookup, ~2ms for 10,000 items

Stage 3: Heavy Ranking (latency budget ~20ms for 500 items)
  input:  500 candidates
  method: large deep neural network (100+ features, cross-feature interactions)
  output: final ranked list
  feature fetch: full user + item + context features from feature store
```

The critical insight: 99.5% of candidates are eliminated before the expensive model sees them.

### Feature Store Architecture

Features are categorized by update frequency and latency requirement:

| Feature type | Examples | Source | Update cadence | Store |
|--------------|----------|--------|----------------|-------|
| User long-term | genre affinity, avg session length, top 10 interests | Spark batch | Daily | Redis (hot) |
| User short-term | last 5 posts interacted with, recent search queries | Flink stream | Minutes | Redis (hot) |
| Item engagement rate | likes/1000 views in last 1h, share rate, completion rate | Flink stream | Minutes | Redis (hot) |
| Creator score | historical engagement rate, account standing | Spark batch | Daily | Redis (hot) |
| Context features | time of day, device type, network quality | Request-time computed | Per-request | Not stored — computed inline |

Online feature store: Redis Cluster with consistent hashing. Feature vectors serialized as MessagePack for compact encoding. Batch TTL: 25 hours (refreshed before expiry by daily pipeline).

Meta's equivalent of LinkedIn's Venice/Voldemort: an internal key-value store purpose-built for feature serving with:
- Sub-millisecond read latency (SSD-backed, in-memory for hot keys)
- High write throughput for streaming feature updates
- Atomic multi-key reads (get all features for one request in one call)

### Model Serving

The heavy ranking model runs under strict latency constraints:

```text
model serving stack:
  TorchServe (or internal equivalent) with batching
  -> GPU inference server
  -> 500 items × 100 features → batch score in ~15ms on T4 GPU
  -> return 500 scores

latency breakdown per ranking request:
  feature fetch (Redis): ~5ms
  model inference (GPU): ~15ms
  total ranking stage:   ~20ms
```

FAISS (Facebook AI Similarity Search) for embedding-based candidate retrieval:

```text
user interest embedding (128-dim) computed daily from watch/engage history
FAISS IVF index over 1B+ post embeddings
ANN query: top-K nearest posts to user embedding
query time: <5ms for top-1000 candidates from 1B items
```

### Feedback Loop and Model Retraining

```text
engagement events
  -> Kafka (topic: feed-engagement, 5B events/day)
  -> Flink (stream): update per-item engagement rates, per-user short-term features
  -> Redis feature store: updated within 2 minutes of event
  -> Spark (batch daily): aggregate long-term features, compute training labels
  -> GPU training cluster: full model retrain (PyTorch DDP, multi-GPU)
  -> model registry: new model artifact versioned and tested
  -> shadow serving: new model runs in shadow alongside production for 24h
  -> gradual rollout: A/B test new model on 5% of traffic → 100%
```

**Filter bubble risk:** Positive engagement signals amplify what users already engage with. Mitigations:
1. **Exploration bonus**: inject a random sample of non-personalized posts (5%) to break loops.
2. **Diversity constraints**: cap same creator/topic at stage 4 filter.
3. **Offline evaluation**: measure recommendation diversity and filter bubble metrics before deploying new models.

### Explainability (EU DSA Compliance)

Every ranking decision must be auditable. Meta logs:

```text
per ranked item:
  - top-3 feature contributions (SHAP-style attribution)
  - which stage produced the final score
  - whether content safety rules were applied
  - timestamp and model version

user-facing explanation (on request):
  "Shown because: you recently liked posts from [creator], 
   and posts about [topic] are trending in your network."
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Heavy ranking model latency spike | p99 latency > 25ms (SLA breach) | Circuit breaker: fall back to lightweight scorer result as final ranking; degrade quality gracefully |
| Feature store unavailable | Redis health check failure or read error rate > 1% | Use stale features from local in-process cache (30s TTL); if no cache, fall back to heuristic scoring (recency + creator score) |
| Kafka consumer lag grows (feature staleness) | Consumer lag > 10 min | Alert oncall; temporarily increase Flink parallelism; features degrade gracefully (older features still valid, just less fresh) |
| Model bias introduced by feedback loop | Diversity metric drops below threshold in A/B test | Rollback model; increase exploration bonus; review training label construction |
| ANN index stale (rebuilt daily) | New posts not appearing in candidates | Keep previous index active during rebuild; atomic swap on completion |

## Observability

- metric: ranking pipeline p50/p95/p99 latency (per stage)
- metric: feature store read latency and cache hit ratio
- metric: model inference latency (p99 on GPU cluster)
- metric: Kafka consumer lag on engagement event topics
- metric: feed engagement rate (CTR, completion rate, share rate) per model version
- metric: recommendation diversity score (unique creators, topics in served feed)
- metric: exploration injection rate (% non-personalized posts served)
- log: per-request ranking decisions with feature attributions (sampled 0.1%)
- trace: full ranking pipeline from candidate retrieval to final sort
- alert: ranking p99 > 200ms; feature store p99 > 10ms; model serving error rate > 0.1%

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Three-stage funnel | 99.5% of candidates eliminated cheaply; expensive model only runs on top-500 | Complex pipeline to maintain; each stage introduces potential quality loss | Single heavy model for all 10,000 candidates: latency budget blown (10,000 × 0.1ms = 1 second) |
| Online features in Redis vs offline in Hive only | Low-latency feature serving; real-time personalization (user's last 5 interactions visible immediately) | Additional operational complexity of Redis cluster; feature synchronization | Hive-only (batch): simpler but feature staleness would be 24h, missing session context |
| Daily full retrain vs continuous training | Simpler infrastructure; predictable compute budget | 24-hour model staleness; new viral posts not well-modeled until next retrain | Continuous training: fresher models but much harder to validate before deploying; training/serving skew risk |
| Exploration bonus injection | Breaks filter bubble; enables discovery; data diversity for training | Small reduction in short-term engagement metrics | No exploration: maximizes local engagement but degrades long-term retention as feed becomes monotonous |
| FAISS for embedding retrieval | Sub-millisecond ANN queries over 1B+ posts | Index rebuild required on model update; approximate (not exact) | Exact nearest neighbor: O(N) per query, too slow at 1B items |

## Interview It

**Meta framing:** "Design Meta's News Feed ranking system." Strong answers cover the three-stage funnel, feature store with online vs offline split, model serving latency budget, feedback loop with filter bubble mitigation, and failure degradation. Weak answers describe a single ML model without addressing candidate retrieval, feature serving latency, or what happens when the model is down.

**Follow-ups:**

1. A new user just joined Meta. They have no watch history and no engagement history. How does the ranking system handle their first session (cold start)?
2. Your heavy ranking model just deployed and diversity metrics dropped 30% overnight. Walk me through how you detect and respond to this.
3. How would you A/B test two different ranking model architectures simultaneously on 10% of traffic each without the experiments contaminating each other?
4. The EU DSA requires that users can understand why a specific post was ranked first. How does your design support this?
5. How would you reduce ranking pipeline end-to-end latency from 200ms to 50ms?

## Ship It

- `outputs/design-doc-ranking-pipeline.md`
- `outputs/feature-store-spec.md`
- `outputs/interview-card-ranking-ml-serving.md`

## Exercises

1. **Easy** — List the three stages of the ranking funnel, the number of candidates at each stage, and the latency budget for each stage.  
2. **Medium** — Design the Kafka consumer that reads engagement events and updates the Redis feature store. What guarantees does it need (at-least-once vs exactly-once), and why?  
3. **Hard** — Design the cold-start ranking strategy for a user who joined 3 minutes ago. What features can you use, and what fallback ranking do you apply until enough engagement history accumulates?  

## Further Reading

- [Meta's paper on deep learning for ranking (DLRM)](https://arxiv.org/abs/1906.00091)  
- [FAISS — A Library for Efficient Similarity Search](https://engineering.fb.com/2017/03/29/data-infrastructure/faiss-a-library-for-efficient-similarity-search/)  
- [EU Digital Services Act — Article 27 (Recommender systems transparency)](https://digital-strategy.ec.europa.eu/en/policies/digital-services-act-package)  
