# Recommendation Engine at Streaming Scale

> Recommendations are not about finding content you will like. They are about reducing the time it takes you to decide to watch something.

**Type:** Build  
**Company focus:** Netflix  
**Learning goal:** Design a Netflix-scale recommendation system. Cover collaborative filtering vs content-based approaches, two-tower model architecture, offline training pipelines, online serving with a feature store, the cold start problem, A/B testing for recommendations, and personalized homepage row ranking.  
**Prerequisites:** `04-recommendation-engine`, `06-ab-testing-platform`  
**Estimated time:** ~90 min  
**Primary artifact:** recommendation system design doc + feature store spec  

## The Problem

Netflix's entire product depends on showing each subscriber the right content at the right time. The homepage consists of rows, each representing a recommendation algorithm (Trending Now, Top 10, Because You Watched X, etc.). Every row must be populated with personalized content in under 200ms to avoid delaying page load.

Design the end-to-end recommendation system that powers these rows for 270M subscribers.

## Clarify

- What is the primary latency target for homepage recommendation serving? (<200ms p99)
- How many rows are on the homepage? (~20-40 rows per subscriber)
- What signals are available? (watch history, ratings, search queries, time-of-day, device type)
- How stale can recommendations be? (minutes to hours depending on row type)
- How do you handle new subscribers with no history (cold start)?
- How are recommendation algorithms selected and ranked per subscriber?

## Requirements

### Functional

- Serve personalized recommendations for each subscriber in under 200ms.
- Rank which row types to show each subscriber.
- Rank titles within each row.
- Handle cold start (new subscribers with no history).
- Support offline model training and online model serving.
- Enable A/B testing of recommendation algorithms.

### Non-functional

- Serving latency: p99 under 200ms per homepage request.
- Model freshness: recommendations updated at least every 24 hours; some signals updated in near-real-time.
- Scale: 270M subscribers, each with a personalized homepage.
- Availability: recommendation service must degrade gracefully (fallback to trending) if any component is slow.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Subscribers | 270M | drives pre-computed recommendation storage |
| Titles in catalog | ~36,000 | item embedding vector size |
| Rows per homepage | ~20 | total items to score per request |
| Items per row | 20–75 | scores to compute per row |
| Recommendations per subscriber | ~200 pre-computed candidates | balance freshness vs compute |
| Offline training cadence | daily full retrain + hourly incremental | drives Spark cluster sizing |

## Architecture

```text
[offline path]
  playback events + ratings + searches
  -> feature engineering (Spark)
  -> model training (TensorFlow/PyTorch on GPU cluster)
  -> model registry
  -> model deployment (online serving layer)

[online path]
  subscriber homepage request
  -> recommendation API (Zuul edge)
  -> row selector (which row types to show this subscriber)
  -> per-row ranker (retrieval + ranking)
     -> feature store (subscriber and item features)
     -> candidate retrieval (ANN index)
     -> scoring model (two-tower or GBT)
     -> re-ranker (diversity, freshness, business rules)
  -> response assembly
  -> fallback layer (EVCache pre-computed recommendations)
```

### Collaborative Filtering

Collaborative filtering predicts what a subscriber will watch based on what similar subscribers have watched:

```text
similarity(user_A, user_B) = cosine(embedding(A), embedding(B))
candidates_for_A = union(top_N_titles_of(most_similar_users_to_A))
```

Strengths: captures latent preferences without requiring content features.  
Weaknesses: cold start for new users and new titles; does not explain recommendations.

### Content-Based Filtering

Content-based filtering recommends titles similar in features (genre, cast, director, theme) to what a subscriber has previously enjoyed:

```text
similarity(title_A, title_B) = cosine(content_embedding(A), content_embedding(B))
candidates = top_N_titles similar_to(recently_watched_by(subscriber))
```

Strengths: works for new titles immediately; interpretable.  
Weaknesses: limited discovery — tends to recommend more of the same.

### Two-Tower Model Architecture

Netflix uses a two-tower (dual-encoder) model that learns separate embeddings for users and items:

```text
user_tower(subscriber_features) -> user_embedding (128-dim)
item_tower(title_features) -> item_embedding (128-dim)
score(user, item) = dot_product(user_embedding, item_embedding)
```

At serving time, the user embedding is computed once, then approximate nearest neighbor (ANN) search retrieves the top-K candidates from the item embedding space — typically in single-digit milliseconds.

### Feature Store

The feature store serves precomputed features for subscribers and titles to the online serving layer:

| Feature type | Examples | Update cadence |
|--------------|----------|----------------|
| Subscriber long-term features | genre affinity, platform preferences | daily |
| Subscriber short-term features | last 5 titles watched, last search query | near-real-time (minutes) |
| Title features | genre, cast, runtime, popularity score | daily |
| Context features | time of day, device type, country | request-time |

Netflix's EVCache (Memcached-based) serves feature lookup at <1ms.

### Cold Start Problem

New subscribers have no watch history. Strategies:

1. **Onboarding questionnaire**: ask for genre preferences and surface content accordingly.
2. **Trending fallback**: show global or regional trending content as default.
3. **Transfer learning from demographics**: subscriber's country and signup context can seed an initial embedding.
4. **Fast-path personalization**: after the first 3–5 views, short-term signals update recommendations in near-real-time.

### Homepage Row Ranking

Not all rows are shown to all subscribers. A meta-ranker decides which row types to display and in what order:

```text
row_score(subscriber, row_type) = model(subscriber_features, row_type_affinity, session_context)
show top-K rows with highest score
within each row: rank titles by per-row ranker
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Recommendation API is slow | p99 latency exceeds threshold | Serve pre-computed recommendations from EVCache; fall back to trending |
| Feature store unavailable | Feature lookup timeout | Use stale features from local cache; degrade to context-only scoring |
| ANN index is stale (model redeployment lag) | Recommendation diversity drops | Keep previous model hot until new index is warmed |
| A/B experiment contaminates control group | Control group metrics shift | Assignment service enforces strict bucketing; holdout group monitors contamination |
| Cold start subscriber sees irrelevant content | High skip rate for new users | Accelerate short-term signal update cadence for first-session subscribers |

## Observability

- metric: homepage recommendation API latency at p50/p95/p99
- metric: cache hit rate for pre-computed recommendations (EVCache)
- metric: recommendation acceptance rate (did subscriber click / start watching)
- metric: diversity score per row (unique genres, studios)
- metric: cold start duration (time from signup to first personalized recommendation)
- metric: model staleness (age of currently deployed model)
- log: fallback activations (trending served instead of personalized) with reason
- trace: recommendation request from row selector through ANN retrieval to response

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Two-tower model over matrix factorization | ANN retrieval is faster at serving; user and item embeddings can be updated independently | Training is more complex; ANN index must be rebuilt on model update | Matrix factorization simpler but full user-item matrix does not scale to 270M users |
| Pre-computed recommendations in EVCache | Sub-millisecond serving for most requests | Recommendations go stale between compute runs | Pure on-demand scoring would miss 200ms SLA at scale |
| Daily full retrain + hourly incremental | Balance between freshness and compute cost | New content may not appear in recommendations for up to 24 hours | Continuous training would improve freshness but requires very expensive streaming training infrastructure |
| Trending fallback for cold start | Always shows something compelling to new subscribers | Generic; does not reflect individual interests | No-recommendation fallback (empty homepage) is unacceptable |

## Interview It

**Netflix framing:** "Design Netflix's recommendation system." Strong answers cover offline training pipeline, online serving with feature store and ANN retrieval, cold start handling, and fallback chains. Weak answers describe only collaborative filtering without discussing serving latency, feature freshness, or experimentation.

**Follow-ups:**
1. How would you handle a popular new title that is not yet in the trained model?
2. What happens to recommendations if the feature store is down for 5 minutes?
3. How do you ensure that A/B testing a new recommendation algorithm does not pollute the control group?
4. How would you measure whether a recommendation change actually improved subscriber satisfaction?
5. How would you reduce recommendation latency from 200ms to 50ms?

## Ship It

- `outputs/design-doc-recommendation-engine.md`
- `outputs/feature-store-spec.md`
- `outputs/interview-card-recommendation-engine.md`

## Exercises

1. **Easy** — List five subscriber features and five title features you would include in the feature store. What is the freshness requirement for each?  
2. **Medium** — Design the offline training pipeline that produces two-tower embeddings using Spark for feature engineering and a GPU cluster for model training.  
3. **Hard** — Design the real-time signal pipeline that updates short-term subscriber features within 2 minutes of a playback event.  

## Further Reading

- [Netflix recommendation blog series](https://netflixtechblog.com/system-architectures-for-personalization-and-recommendation-e081aa94b5d8)  
- [Two-tower models (Google)](https://research.google/pubs/pub48951/)  
- [Feature stores for ML (Feast)](https://feast.dev/)  
