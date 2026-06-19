# Amazon Full Mock Loop

> 45 minutes to prove you are customer-obsessed, technically deep, and can own a system end to end.

**Type:** Mock Interview  
**Company focus:** Amazon  
**Learning goal:** Synthesize all Amazon-specific signals — Working Backwards, Leadership Principles, capacity modeling, AWS service selection, failure modes, and observability — into a coherent 45-minute interview answer.  
**Prerequisites:** All previous lessons in this phase (01–06).  
**Estimated time:** 45 min timed + 30 min debrief  
**Primary artifact:** Amazon mock interview scorecard  

## The Prompt

> "Design Amazon's product recommendation engine."

Take 30 seconds before speaking. Establish Working Backwards framing, then proceed.

## Working Backwards Framing

Before touching architecture, state the customer promise:

> "Customers see personalized recommendations that help them discover products they actually want — not generic bestsellers, but items that match their individual taste. This means recommendations must be relevant, fresh (reflecting the current session), and delivered fast enough not to interrupt the shopping experience."

This is your anchor. Every architectural decision you make must trace back to this promise.

**LP demonstrated:** Customer Obsession — you started from the customer experience, not the algorithm.

## Capacity Model

State numbers before drawing boxes. Interviewers at Amazon are engineers who know these numbers. Getting them wrong is a red flag.

| Dimension | Estimate | Derivation |
|-----------|----------|------------|
| Active customers | 300M | Amazon's public disclosure |
| Products in catalog | 350M | Amazon's public disclosure |
| Sessions per day | ~10M active shopping sessions | ~3% of MAU in session at any time |
| Recommendations per session | 50 | homepage + product page + cart widgets |
| Recommendation impressions/day | 500M | 10M sessions × 50 recommendations |
| Inference latency budget | < 100 ms | must not delay page render |
| Collaborative filtering model size | ~100 GB (user-item matrix, sparse) | 300M users × 350M items, < 0.01% non-zero entries |
| Real-time signal events/day | ~5B | clicks, views, add-to-cart, purchase |

## Architecture

### Three-Layer Recommendation Stack

**Layer 1: Offline collaborative filtering (batch, daily)**
- User-item interaction matrix built from purchase history, ratings, and browsing.
- Matrix factorization (ALS or SVD) produces user embeddings and item embeddings.
- Runs as a daily EMR/Spark job; outputs stored in S3 and loaded into ElastiCache (Redis) for fast lookup.
- Latency: embeddings are pre-computed; lookup is O(1) in Redis.

**Layer 2: Content-based filtering (product embeddings, weekly)**
- Product attributes (title, category, brand, description) encoded into dense embeddings via a transformer model.
- Similar products computed offline; stored as `product_id → [top-100 similar products]` in DynamoDB.
- Used for cold-start: new users with no history get recommendations based on currently viewed product.

**Layer 3: Real-time signals (session, streaming)**
- Lambda consumes Kinesis stream of click/view/add-to-cart events.
- Session state stored in ElastiCache (Redis) with a 30-minute TTL.
- Re-ranks the offline candidate list based on current session signal (if you viewed hiking boots, boost outdoor products).
- Kinesis → Lambda → Redis session update → re-ranking on next page load.

### Serving Path (Online, < 100 ms)

1. Page load triggers recommendation API call to the Recommendation Service (ECS on Fargate).
2. Recommendation Service fetches user embedding from Redis (cache hit: ~1 ms).
3. Fetches top-100 candidate products from DynamoDB (offline CF output).
4. Fetches current session signal from Redis (current browsing context).
5. Re-ranks candidates using a lightweight scoring function (no model inference in hot path).
6. Returns top-50 recommendations to the page.
7. Total latency budget: 5 ms Redis + 10 ms DynamoDB + 5 ms scoring = ~20 ms (well within 100 ms).

### LP Demonstration

| LP | Design Decision |
|----|----------------|
| Customer Obsession | Started with customer promise; real-time signals make recommendations reflect current intent, not just history |
| Dive Deep | Stated concrete numbers (300M users, 350M products, 500M impressions/day); explained offline vs online trade-off |
| Bias for Action | Proposed MVP (offline CF only) → iterate to real-time signals; don't design for perfection from day 1 |
| Frugality | Pre-computed embeddings in Redis avoid real-time ML inference (expensive GPU inference) on every page load |
| Ownership | Identified failure modes and mitigations before being asked; proposed concrete SLOs |

## Timeline: 5-Minute Milestones

| Time | Expected Progress |
|------|-------------------|
| 0–5 min | Working Backwards statement; customer promise; capacity numbers (10M sessions, 500M impressions) |
| 5–15 min | Three-layer architecture sketch; data stores named (Redis, DynamoDB, S3/EMR); offline vs real-time split |
| 15–25 min | Deep dive on serving path; latency budget breakdown; cold-start problem addressed |
| 25–35 min | Failure modes: stale recommendations, cold start for new users, Kinesis consumer lag, Redis cache eviction |
| 35–45 min | LP alignment summary; observability (SLOs, metrics); trade-offs named (offline CF vs neural CF; content-based vs collaborative) |

## Strong-Hire vs Weak-Hire Patterns

### Strong-Hire Signals
- Opens with Working Backwards before the first diagram.
- States capacity numbers unprompted and derives them from first principles.
- Distinguishes offline (quality, higher latency) from online (speed, lower quality) and explains the trade-off explicitly.
- Identifies the cold-start problem and proposes content-based filtering as mitigation.
- Names the top failure modes before the interviewer asks: stale recommendations during EMR job delay, Redis cache eviction under memory pressure, Kinesis consumer lag during traffic spike.
- Proposes concrete SLOs: recommendation API < 100 ms at p99, click-through rate (CTR) as the business metric.
- Mentions that CTR is the outcome metric but warns it can be gamed — A/B testing with revenue-per-session as the ground truth.

### Weak-Hire Signals
- Opens with "I'd use a machine learning model" without naming the customer outcome.
- No capacity numbers; architecture is a generic ML system diagram.
- Treats collaborative filtering and content-based filtering as interchangeable without explaining the trade-off.
- No cold-start mitigation.
- No failure modes discussed until explicitly prompted.
- Uses LP terms decoratively: "I'd be customer-obsessed about the recommendations" without connecting it to a specific design decision.
- Does not mention latency budget or that serving path must return in < 100 ms.

## Scoring Rubric

Total: 100 points. Hire threshold: 70 points. Strong-hire threshold: 85 points.

| Dimension | Weight | What earns full points |
|-----------|--------|------------------------|
| LP Framing | 15 | Working Backwards statement in first 2 minutes; at least 2 specific LPs demonstrated with design decisions, not decoration |
| Working Backwards | 15 | Customer promise stated before architecture; design decisions traced back to customer outcome throughout the session |
| Capacity Model | 10 | Correct order-of-magnitude for sessions, impressions, and data sizes; derived from first principles, not memorized |
| Architecture | 20 | Three-layer stack (offline CF + content-based + real-time); correct data stores (Redis, DynamoDB, S3); serving path < 100 ms |
| Deep Dive | 20 | Cold-start problem identified and mitigated; latency budget broken down per component; offline model update cadence addressed |
| Failure Modes | 10 | At least 3 failure modes with concrete mitigations: stale recommendations, cold start, Kinesis lag |
| Trade-offs | 5 | At least 2 explicit trade-offs: offline CF vs neural CF; content-based vs collaborative; pre-compute vs on-demand inference |
| Observability | 5 | SLOs stated (API latency, CTR); named metrics (cache hit rate, Kinesis iterator age); business metric vs technical metric distinction |

## Common Follow-Up Questions

**Q1: How does your recommendation engine handle a brand-new customer who has never purchased on Amazon?**

Ideal answer: Cold-start problem. Three mitigations in priority order: (1) content-based filtering on the currently viewed product — if they're looking at a tent, recommend camping gear; (2) geographic and demographic signals (location, device type) to infer likely interests; (3) global bestsellers in the browsed category as a fallback. The offline CF model has nothing to work with for new users, so it must be bypassed entirely until sufficient interaction data is collected (typically 5+ events).

**Q2: Your real-time Kinesis consumer is 10 minutes behind during a Prime Day spike. What is the customer impact and what do you do?**

Ideal answer: Recommendations revert to the offline collaborative filtering layer — still personalized but not reflecting the current session. Customer sees "you might also like" based on history rather than current browsing context. Mitigation: (1) scale Kinesis shards before Prime Day; (2) Lambda enhanced fan-out to avoid shared-throughput bottleneck; (3) alert on `GetRecords.IteratorAgeMilliseconds` > 60,000 ms; (4) graceful degradation — offline CF is better than no recommendations.

**Q3: How would you A/B test a new recommendation algorithm without harming revenue?**

Ideal answer: Holdout experiment with traffic splitting at the user level (consistent assignment via user ID hash). Primary metric: revenue per session (not CTR — CTR can be gamed by showing cheap items). Secondary: click-through rate, add-to-cart rate, return rate. Guardrail metric: page load time (new algorithm must not breach p99 latency SLO). Minimum detectable effect: ~1% revenue-per-session improvement to be worth the complexity cost. Run for at least one full week to capture weekly seasonality.

**Q4: How do you prevent recommendations from becoming a filter bubble — showing users only what they already know?**

Ideal answer: Explicit exploration budget: 10–20% of recommendation slots are filled with serendipitous items outside the user's known interest graph (epsilon-greedy exploration). Diversity constraint on the candidate set: no more than 3 items from the same brand or category in the top-10. Long-term novelty metric tracked separately from CTR. This demonstrates awareness that optimizing a single metric (CTR) can degrade the customer experience in ways that only surface over months.

**Q5: The offline collaborative filtering job takes 6 hours on EMR. A product goes viral at 2 PM. When do recommendations reflect it?**

Ideal answer: Real-time signals (Layer 3) reflect it within seconds — if enough users are clicking the viral product in their sessions, the session-based re-ranking will boost it. But the offline CF model won't include it until the next daily batch. Mitigation: trending product detection as a separate lightweight pipeline — identify items with >5× normal click velocity over a 30-minute window, inject them into the candidate set as a "trending now" slot, bypassing the CF model entirely for that position.

## Ship It

- `outputs/amazon-mock-scorecard.json`

## Further Reading

- [Amazon Builder's Library: Challenges with distributed systems](https://aws.amazon.com/builders-library/challenges-with-distributed-systems/)
- [Amazon's Working Backwards Process](https://www.allthingsdistributed.com/2006/11/working_backwards.html)
- [Amazon Leadership Principles](https://www.amazon.jobs/content/en/our-workplace/leadership-principles)
- [Matrix Factorization Techniques for Recommender Systems (Netflix paper)](https://datajobs.com/data-science-repo/Recommender-Systems-%5BNetflix%5D.pdf)
- [Real-time Machine Learning at Scale (AWS re:Invent)](https://www.youtube.com/watch?v=DmsXRiPjt3Y)
