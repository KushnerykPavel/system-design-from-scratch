# Product Catalog & Search at E-Commerce Scale

> Customers who can't find it can't buy it — search is the most important revenue path.

**Type:** Build
**Company focus:** Amazon
**Learning goal:** Design a product catalog system and search index that handles 350M+ products, 300M+ customers, and Black Friday peak QPS.
**Prerequisites:** `05-storage-indexing-and-access-patterns/03-indexes`, `17-search-crawl-and-monitoring-systems/02-search-autocomplete`
**Estimated time:** ~90 min
**Primary artifact:** search index design + catalog schema

---

## The Problem

Amazon serves 350M+ active product listings to 300M+ customers globally. At any moment ~50,000 searches/sec hit the platform; on Black Friday that surges to ~500,000 searches/sec. The system must return relevant results in under 200ms at p99 while simultaneously showing accurate stock levels, prices that change every few seconds, and personalized rankings.

Two fundamentally different access patterns collide here:
- **Catalog reads** — structured, key-value, high fan-out (product page loads)
- **Search** — full-text + faceted, ranking-heavy, near-real-time index updates

---

## Scale Envelope

| Metric | Normal | Black Friday Peak |
|---|---|---|
| Active products | 350M | 350M |
| Customers | 300M | 300M |
| Searches/sec | ~50K | ~500K |
| Product page views/sec | ~200K | ~2M |
| Catalog write rate (new/updates) | ~10K/sec | ~10K/sec |
| Price/inventory updates | ~100K/sec | ~500K/sec |

Storage estimate: 350M products × ~2KB avg metadata = ~700GB catalog; search index is ~3–5× larger due to inverted index = ~3TB.

---

## Catalog Store Design

### Product Data — DynamoDB

Primary key: `product_id` (UUID). High read throughput, single-digit millisecond latency.

**Key attributes stored:**
- `title`, `description`, `brand`, `category_path` (e.g., `Electronics/Cameras/DSLR`)
- `image_urls[]`, `bullet_points[]`, `asin`
- `seller_id`, `prime_eligible`, `fulfillment_type`
- `avg_rating`, `review_count`
- `created_at`, `updated_at`

**Hot partition mitigation:** Viral products (e.g., a new iPhone launch) create hot partition risk. Mitigate via:
1. DynamoDB adaptive capacity (automatic, built-in)
2. Elasticache (Redis) caching layer in front — cache product data with TTL 1h
3. Write sharding with suffix randomization during extreme spikes

### Inventory & Pricing — Aurora (RDS)

Transactional: inventory levels and prices need ACID semantics for reservation. Aurora multi-AZ with read replicas for read scale.

Separate from catalog because:
- Update rate is 10–50× higher than catalog metadata
- Requires row-level locking for inventory reservation
- Price changes must not cascade to search index on every tick

### Price/Inventory Cache — Redis

TTL 5 seconds for prices and stock levels. Slightly stale is acceptable for browsing; reservation uses the Aurora source of truth.

---

## Search Architecture

### Index Technology

Amazon uses its proprietary **A9 search engine**. In a design interview, model it with **OpenSearch (Elasticsearch)**.

**Indexed fields:**
- `title` (highest weight), `description`, `bullet_points` (full-text, analyzed)
- `category_path` (keyword, for facets)
- `brand` (keyword)
- `price` (numeric range)
- `avg_rating`, `review_count` (numeric)
- `prime_eligible` (boolean)
- `sales_velocity_7d` (numeric, updated daily)
- `conversion_rate_30d` (numeric, ML-computed)

**Index sharding:** 350M products / ~50M docs per shard = 7 primary shards. Each shard replicated 2× for read throughput and fault tolerance.

### Search Ranking

Ranking is a multi-signal blend:

1. **Text relevance** — BM25 score from OpenSearch
2. **Sales velocity** — products selling faster rank higher (signals demand)
3. **Conversion rate** — historical click-to-purchase ratio
4. **Prime eligibility** — Prime customers prefer Prime results
5. **Average rating** — Wilson score lower bound (avoids bias for products with 1 review)
6. **ML re-ranking** — a lightweight model (LambdaMART or similar) combines signals and applies personalization for logged-in users

In interview: explain that the ML layer runs as a post-retrieval re-ranking step on the top-K (e.g., 200) candidates from OpenSearch, not on the full index.

### Faceted Search

Facets let users filter by category, brand, price range, Prime, rating. Pre-computing facet counts avoids per-query aggregation at scale.

**Implementation:**
- OpenSearch aggregations for real-time counts (expensive at 500K QPS)
- For top categories: pre-computed facet counts stored in Redis, refreshed every 5 minutes
- Facet cache key: `facet:{query_hash}:{filter_combination}`
- On cache miss: delegate to OpenSearch aggregation; async warm the cache

### Autocomplete

- Prefix trie stored in Redis (ZSET per prefix, scored by popularity)
- Top-10 suggestions per prefix
- Personalized for logged-in users: merge global suggestions with user's search history
- p99 latency target: <20ms

**Data pipeline:** Offline batch job (daily) computes top queries from search logs → loads into Redis. Near-real-time trending queries (last 1h) patched on top via a streaming job (Kinesis → Lambda → Redis).

---

## Product Page Assembly

A product page assembles data from multiple sources:

```
Browser → CDN (CloudFront)
        → Product Service (reads DynamoDB + Elasticache)
        → Price/Inventory Service (reads Redis, falls back to Aurora)
        → Review Service (reads separate DynamoDB table)
        → Recommendation Service (ML model)
```

**Caching strategy:**
- Static product content (title, images, description): TTL 1h in Elasticache + CDN
- Price and stock: TTL 5s in Redis (near-real-time)
- Reviews: TTL 10min

**Cache stampede prevention:** Probabilistic early expiration (PER) — re-fetch slightly before TTL expires with probability proportional to remaining TTL. Prevents thundering herd when a cached item expires.

---

## Index Freshness Pipeline

Catalog changes (new products, price updates) must reach the search index quickly:

```
Seller Portal → Catalog Service → DynamoDB
                                → SNS/SQS → Index Worker → OpenSearch
```

- Normal catalog updates: eventual consistency acceptable; index lag ~30s
- Price changes: NOT propagated to search index on every tick (too expensive); price filter in search reads Redis cache
- New product launch: expedited queue with higher priority; target index lag <5min

**Reindex risk:** Full reindex of 350M products takes hours. Mitigate by using blue/green indexes (build new index in background, cut over atomically via index alias swap).

---

## Failure Modes

### 1. Search Index Reindex During Traffic
- **Risk:** OpenSearch cluster undergoing shard rebalance or reindex serves degraded latency or timeouts
- **Mitigation:** Blue/green index aliases; read from old index until new is ready; circuit breaker to fall back to DynamoDB keyword scan (degraded but functional)

### 2. DynamoDB Hot Partition on Viral Product
- **Risk:** Single product_id gets millions of reads/sec (new iPhone launch event)
- **Mitigation:** Elasticache in front with TTL 1h; DAX (DynamoDB Accelerator) for microsecond cache; if needed, pre-warm cache before anticipated launch

### 3. Cache Stampede on Black Friday Launch
- **Risk:** Redis TTL expiry at midnight causes all servers to simultaneously query OpenSearch/DynamoDB
- **Mitigation:** Jittered TTLs (TTL = base ± random(0, base×0.1)); probabilistic early expiration; pre-warming cache before event

---

## Trade-offs

| Decision | Chosen | Rejected | Reason |
|---|---|---|---|
| Catalog DB | DynamoDB | PostgreSQL | DynamoDB scales horizontally to 350M products without schema migrations; PostgreSQL would require sharding |
| Search | OpenSearch (A9 proxy) | DynamoDB full scan | Full-text ranking impossible without inverted index |
| Price freshness | 5s Redis TTL | Strong consistency from Aurora | Aurora can't serve 500K reads/sec; 5s stale is acceptable for browsing |
| Facets | Pre-computed in Redis | On-the-fly aggregation | OpenSearch aggregation at 500K QPS is prohibitively expensive |
| Index update | Async via SQS | Synchronous on write | Synchronous indexing would double write latency and couple availability |

---

## Follow-up Questions

1. **Cold-start for new products:** A product with zero sales history has no conversion rate or velocity signal. How do you rank it? (Answer: boost new products temporarily, use category-level signals as prior, A/B test)

2. **Seller fraud:** A seller inflates ratings or uses fake purchases to boost conversion rate. How do you detect and mitigate this in ranking signals?

3. **Multi-language search:** Amazon operates in 20+ countries. How does your search index handle Japanese kanji, Arabic right-to-left, and German compound words? (Answer: per-locale analyzers in OpenSearch, locale-aware tokenization)

4. **Image search:** Customers want to search by uploading a photo. How would you add visual search without changing the core text search architecture?

5. **Personalized ranking at scale:** How do you serve personalized rankings to 300M different users without computing a unique ranking per user per query?
