# Index Freshness and Ranking Updates

> Search quality is not just about the ranking model; it is about how quickly the serving index can reflect the world without destabilizing itself.

**Type:** Build
**Company focus:** Google
**Learning goal:** Explain how index freshness, document updates, ranking-signal refresh, and rollout safety interact in search systems with large serving fleets.
**Prerequisites:** `08-consistency-replication-and-transactions/04-replication-lag`, `13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive`, `17-search-crawl-and-monitoring-systems/01-web-crawler`
**Estimated time:** ~60 min
**Primary artifact:** update-plan validator + design review prompt

## The Problem

Design the update path for a search index so new or changed documents become searchable quickly while ranking updates and schema changes roll out safely. The system should handle continuous updates, deletes, and large backfills.

This lesson matters because many search answers hand-wave "the crawler updates the index." Senior answers separate document ingestion, segment publishing, ranking-feature refresh, delete handling, canary rollout, and rollback strategy.

## Clarify

- Is document freshness more important than ranking perfection?
- Are updates append-only, full rewrites, or mixed with deletes?
- Do ranking changes require model rollout separate from index data rollout?
- What freshness SLO applies to high-priority content?

If the prompt stays general, assume document freshness in minutes, ranking signals refreshing on a separate cadence, support for deletes, and staged rollout across serving shards.

## Requirements

### Functional

- Publish newly crawled or updated documents into serving indexes.
- Remove or tombstone deleted documents.
- Refresh ranking signals without full reindex every time.
- Support backfills and schema-compatible reprocessing.
- Roll out changes safely with canaries and rollback.

### Non-functional

- Keep update lag within target for high-value content.
- Avoid globally corrupting serving results during rollout.
- Bound resource impact of reindexing and backfills.
- Preserve search availability while updates continue.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Document updates | 20M docs/min peak | drives segment build and publish strategy |
| Delete events | 2M docs/min peak | delete handling cannot require full rewrites |
| Serving shards | 25K | rollout safety and config distribution matter |
| High-priority freshness SLO | under 10 minutes | hot-path publishing must be selective |
| Backfill size | 30B docs | reindex tooling must coexist with live traffic |

## Architecture

```text
crawl / content updates
  -> document transform
  -> segment builders
  -> validation + canary compare
  -> snapshot or segment publisher
  -> serving shard rollout

ranking features / models
  -> feature refresh pipeline
  -> serving feature stores
```

Design notes:

1. Separate document visibility from ranking perfection so you can trade freshness against quality safely.
2. Use staged publish and canary diffing before broad rollout.
3. Represent deletes with tombstones or fast invalidation paths rather than waiting for full reindex.
4. Treat schema evolution as a rollout problem, not only an indexing problem.

## Data Model & APIs

Core fields:

```text
doc_id
content_version
segment_id
publish_state
ranking_feature_version
delete_tombstone
schema_version
```

Useful interfaces:

- `PublishSegment(segment_id, shard_set, schema_version)`
- `ApplyDelete(doc_id, tombstone_version)`
- `PromoteRankingFeatures(feature_version)`
- `StartBackfill(range, transform_version)`

This topic becomes much clearer when the answer names which version boundaries exist in the system.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| a bad segment rollout corrupts results | canary diff failures and result quality checks | blue/green publish and fast rollback |
| deletes lag far behind inserts | delete backlog and stale-hit audits | tombstone fast path and prioritized delete processing |
| ranking features mismatch document snapshot | version skew metrics and serving compare tests | explicit feature-version pinning |
| backfill starves live freshness | live publish lag during backfill windows | isolated backfill resources and throttled catch-up |

## Observability

- metric: document publish lag by priority tier
- metric: delete backlog age and stale-hit rate
- metric: shard rollout success and rollback count
- metric: feature-version skew between publishers and serving shards
- log: publish transitions with segment ID, schema version, and canary verdict
- trace: update ingestion through transform, segment build, and shard promotion
- SLO: high-priority document updates are visible in search within the freshness target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| segment or snapshot publishing | stable serving behavior | bounded freshness between publishes | mutating serving index directly per document |
| tombstone delete path | faster removal from results | extra compaction and cleanup work | waiting for next full rebuild |
| separate ranking refresh cadence | safer model rollout | temporary quality lag behind content | tying every ranking change to full index publish |

## Interview It

**Google framing:** "Design how a search system updates results after pages change." Expect questions on freshness SLOs, delete semantics, and safe ranking rollout.

**Cloudflare framing:** "Design index propagation for a globally distributed content-discovery product." Expect follow-ups on rollout safety, regional skew, and how stale shards are detected.

**Follow-ups:**
1. What changes if breaking news pages must appear in search within one minute?
2. How do you delete unsafe content from results quickly?
3. How would you backfill a new extracted field across the whole corpus?
4. What if ranking model updates are riskier than document content updates?
5. How do you compare two index versions safely before full promotion?

## Ship It

- `outputs/design-review-index-freshness.md`

## Exercises

1. **Easy** — Explain why deletes usually need a faster path than full document reindex.
2. **Medium** — Design a canary plan for a new ranking-feature version.
3. **Hard** — Redesign the update pipeline for a product where freshness beats ranking quality for a subset of content.

## Further Reading

- [Managing Gigabytes: Inverted Files and Search Engines](https://nlp.stanford.edu/IR-book/) — classic search indexing background
- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — useful for thinking about rollout safety and serving skew
