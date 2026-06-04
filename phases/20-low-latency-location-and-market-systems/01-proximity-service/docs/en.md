# Proximity Service / Nearby Search

> Nearby search is rarely blocked by one clever data structure; it is blocked by how often you must refresh moving locations without destroying latency.

**Type:** Build
**Company focus:** Google
**Learning goal:** Design a proximity service that keeps nearby queries fast while balancing moving-object freshness, geo-partition skew, and fallback behavior.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/01-shard-key`, `14-rate-limiters-ids-and-hashing/04-consistent-hashing`, `17-search-crawl-and-monitoring-systems/02-search-autocomplete`
**Estimated time:** ~75 min
**Primary artifact:** proximity-plan validator + design review sheet

## The Problem

Design a service that returns nearby drivers, stores, scooters, or users within a radius. Reads must feel immediate, but the write path is continuous because locations keep moving.

This lesson matters because weak answers stop at "use geohash." Strong answers explain how moving writers, hot downtown cells, freshness budgets, and widening-radius fallbacks shape the design.

## Clarify

- Are objects mostly static places, continuously moving devices, or both?
- What freshness target matters: seconds, tens of seconds, or minutes?
- Is ranking purely distance-based, or do availability and business rules also matter?
- Do we need exact geometry, or is approximate candidate retrieval enough before reranking?

If the interviewer leaves this open, assume continuously moving providers, user queries within 1 to 20 km, freshness under 5 seconds, and a two-stage flow: approximate geo lookup followed by exact filtering and ranking.

## Requirements

### Functional

- Ingest continuous location updates for moving objects.
- Return nearby candidates inside a requested radius.
- Support widening-radius fallback when nearby supply is sparse.
- Filter by availability, category, or region-specific constraints.
- Prefer exact distance checks after fast geo candidate retrieval.

### Non-functional

- Keep p99 nearby-query latency under 80 ms.
- Bound stale-location risk during update bursts or partial pipeline lag.
- Avoid hotspots in dense urban cells.
- Degrade safely when fresh updates are delayed.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Nearby queries | 350K req/s peak | drives read-path caching and shard fanout budget |
| Moving objects | 80M active objects globally | shapes partition count and update index size |
| Update rate | 4 location writes/object/min active average | write amplification can dominate cost |
| Query radius | 500 m default, 10 km fallback max | affects candidate explosion and cell fanout |
| Density skew | 50x between suburbs and event hotspots | forces hotspot mitigation in downtown cells |

## Architecture

```text
mobile clients
  -> edge/API gateway
  -> nearby query service
  -> geo-index shards (cell-based)
  -> availability cache
  -> exact distance reranker

device updates
  -> ingest gateway
  -> update dedupe / smoothing
  -> geo-index mutation stream
  -> shard-local hot cell replication
```

Design notes:

1. Use grid or geohash-like cells for candidate retrieval, then compute exact distance on a bounded candidate set.
2. Separate geo presence from business availability so frequent location writes do not require rewriting all ranking metadata.
3. Plan for density skew by splitting or replicating hot cells instead of assuming one fixed partitioning is enough forever.
4. Make widening-radius logic explicit because product behavior under low supply is part of the system design.

## Data Model & APIs

Core entities:

```text
moving_object(object_id, lat, lon, heading, speed, updated_at)
geo_cell(cell_id, shard_id, object_ids[])
availability_state(object_id, state, attributes, expires_at)
nearby_query(request_id, center, radius_m, filters, max_results)
```

Useful interfaces:

- `POST /v1/location:update`
- `GET /v1/nearby?lat=...&lon=...&radius_m=...&type=driver`
- `ExpandRadius(request_id, previous_radius_m, next_radius_m)`
- `MarkAvailability(object_id, state)`

Senior answers call out the difference between approximate candidate selection and exact distance validation. That keeps the latency story honest.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hot city cell overloads one shard | per-cell QPS, shard CPU, and candidate fanout metrics | split hot cells, replicate popular cells, cache top availability sets |
| stale location data after ingest lag | update age histogram and stale-result complaints | TTL old locations, surface freshness budget, fall back to last-good snapshot carefully |
| candidate explosion for wide-radius rural queries | candidates scanned per request and rerank latency | widening-radius caps, hierarchical cells, and result-count-aware early exit |
| duplicate or out-of-order updates move an object backward | sequence-gap metrics and update discard counts | monotonic update versions and bounded dedupe window |

## Observability

- metric: query latency by radius, region, and candidate count
- metric: location update age, dedupe discard rate, and per-cell write skew
- metric: candidates scanned versus results returned
- metric: hot-cell split frequency and replicated-cell read share
- log: widening-radius decisions with supply count and freshness age
- trace: client request through geo lookup, availability filter, and exact rerank
- SLO: 99.9% of nearby queries return compliant results within latency target using data younger than the freshness budget

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| cell-based approximate lookup plus exact rerank | fast bounded search | two-stage serving complexity | exact geometry scan on all objects |
| separate availability cache from geo index | reduces write amplification | extra consistency boundary | rewriting the full geo record for every business-state change |
| hot-cell replication and split strategy | protects dense regions | more control-plane logic | one static partitioning forever |

## Interview It

**Google framing:** "Design nearby search for drivers, stores, or friends." Expect follow-ups on moving updates, dense urban hotspots, and how you handle sparse areas without scanning the world.

**Cloudflare framing:** "Design geo lookup served close to users with regional caches and failure isolation." Expect pressure on edge cacheability, data freshness, and regional failover.

**Follow-ups:**
1. What changes if many objects move every second instead of every 15 seconds?
2. How do you prevent one stadium event from melting the hottest cells?
3. What if availability changes more often than location?
4. How would you support polygon search for delivery zones, not only radius search?
5. What changes at 10x read load but only 2x write load?

## Ship It

- `outputs/design-review-proximity-service.md`

## Exercises

1. **Easy** — Explain why exact distance checks usually happen after candidate retrieval.
2. **Medium** — Compare geohash-style cells with an R-tree or S2-style approach for this workload.
3. **Hard** — Redesign the service when dense-event hotspots create 100x skew for 20 minutes at a time.

## Further Reading

- [S2 Geometry](https://s2geometry.io/) — helpful for reasoning about geo partitioning and coverage
- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — useful when nearby lookups fan out across multiple candidate cells
