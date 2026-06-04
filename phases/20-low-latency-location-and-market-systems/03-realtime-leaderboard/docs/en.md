# Real-Time Leaderboard

> Leaderboards look read-heavy on the surface, but the real interview question is whether your write path, ranking semantics, and anti-cheat policy stay credible under bursty score updates.

**Type:** Build
**Company focus:** Google
**Learning goal:** Design a real-time leaderboard with fast top-N reads, burst-tolerant score ingestion, and explicit trade-offs around exact ordering, tiered scopes, and anti-cheat delays.
**Prerequisites:** `05-storage-indexing-and-access-patterns/03-indexes`, `07-queues-streams-and-workflows/01-queues-vs-streams`, `16-application-backends/06-fanout-patterns`
**Estimated time:** ~75 min
**Primary artifact:** leaderboard validator + ranking trade-off sheet

## The Problem

Design a leaderboard that shows top players, teams, or creators in real time. Scores update continuously, users want their rank quickly, and the product may need global, regional, or friend-scoped views.

This lesson matters because weak answers say "sorted set." Senior answers explain whether scores are exact or eventually reflected, how updates are aggregated, and what happens when anti-cheat or tie-breaking rules delay visibility.

## Clarify

- Are we optimizing for global top-N, per-user around-me rank, or many segmented leaderboards?
- Do scores always increase, or can corrections and reversals happen?
- Must the ranking reflect events instantly, or is a small delay acceptable?
- Are suspicious scores withheld pending validation?

If left broad, assume a gaming or engagement leaderboard with global and regional scopes, high write bursts, top-100 reads plus around-me rank, and anti-cheat validation that may delay a subset of updates.

## Requirements

### Functional

- Accept score updates and maintain leaderboards by scope.
- Return top-N rankings quickly.
- Return a user's approximate or exact rank near their position.
- Support score corrections or reversals.
- Allow suspicious updates to be delayed or quarantined.

### Non-functional

- Keep p99 top-N reads under 40 ms.
- Bound write amplification when one user updates frequently.
- Preserve a clear tie-break policy.
- Avoid letting anti-cheat checks stall the whole leaderboard.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Score updates | 2.5M updates/s peak | drives ingestion aggregation and hot-key behavior |
| Top-N reads | 700K req/s peak | read path wants precomputed or in-memory structures |
| Segments | 500K active scoped leaderboards | per-scope storage cost matters |
| Around-me reads | 25% of total reads | exact rank retrieval can be more expensive than top-N |
| Correction rate | 0.5% delayed or reversed updates | exactness story must survive post-facto changes |

## Architecture

```text
score producers
  -> ingest API
  -> dedupe / anti-cheat gate
  -> score stream
  -> shard-local rank updaters
  -> in-memory sorted sets / segment indexes
  -> top-N cache + around-me rank service
```

Design notes:

1. Distinguish score acceptance from score publication if validation can delay visibility.
2. Keep top-N reads cheap by precomputing or caching ranked windows for hot scopes.
3. Be honest about approximate rank options for large long-tail leaderboards.
4. Treat corrections as part of the design, not as an edge case that magically never happens.

## Data Model & APIs

Core records:

```text
score_event(event_id, player_id, scope_id, delta, version, state)
published_score(player_id, scope_id, score, updated_at)
leaderboard_slice(scope_id, rank_start, rank_end, entries[])
```

Useful interfaces:

- `SubmitScore(scope_id, player_id, delta, version)`
- `GetTop(scope_id, limit)`
- `GetRank(scope_id, player_id, window)`
- `QuarantineScore(event_id, reason)`

Senior answers state whether around-me rank is exact, approximate, or cached with bounded staleness.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one player or team becomes a hot write key | per-entity write rate and shard skew | local aggregation, batching, and per-key throttles |
| anti-cheat pipeline lags and publication stalls | pending validation age and publish delay metrics | publish healthy traffic separately and isolate suspicious queue |
| rank corrections reorder the top unexpectedly | correction-count and rank-volatility metrics | explicit correction semantics and versioned recompute |
| segmented leaderboards explode memory | segment count, hot-scope set size, and cache churn | tiered storage and lazy materialization for cold scopes |

## Observability

- metric: score-ingest latency, publish latency, and validation queue age
- metric: top-N cache hit rate and around-me rank latency
- metric: correction rate and rank volatility for hot scopes
- metric: per-scope memory footprint and hot-key skew
- log: score quarantine, reversal, and manual leaderboard override actions
- trace: score ingest through validation, ranking update, and serving
- SLO: 99.9% of clean score updates become visible within the publication target while top-N reads stay inside latency target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| precomputed top-N slices | very fast hot reads | extra update work and cache churn | scanning full ranking structures on each read |
| separate validation from publication | isolates suspicious updates | more state transitions to explain | blocking all score updates on anti-cheat |
| approximate around-me rank for huge scopes | better latency and cost | less exact user position | exact rank recomputation on every request |

## Interview It

**Google framing:** "Design a real-time leaderboard for a global game or creator product." Expect follow-ups on rank semantics, hot keys, and how top-N differs from around-me reads.

**Cloudflare framing:** "Design low-latency leaderboard reads served globally." Expect pressure on regional caches, hot-scope isolation, and graceful degradation when write bursts spike.

**Follow-ups:**
1. What changes if scores can decrease because of fraud or moderation?
2. How do you support friends-only or regional leaderboards without precomputing everything?
3. What if anti-cheat needs 30 seconds before publication for premium tournaments?
4. How would you handle seasonal resets?
5. What changes if around-me rank matters more than global top-100?

## Ship It

- `outputs/tradeoff-matrix-realtime-leaderboard.md`

## Exercises

1. **Easy** — Explain why top-N reads and around-me rank can require different data paths.
2. **Medium** — Compare exact and approximate rank strategies for very large scopes.
3. **Hard** — Redesign the leaderboard when global tournaments generate 100x update spikes for short windows.

## Further Reading

- [Redis sorted sets](https://redis.io/docs/latest/develop/data-types/sorted-sets/) — useful intuition for ordered score structures
- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — relevant for hot-scope fanout and read latency
