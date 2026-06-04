# Search Autocomplete

> The user experiences autocomplete as a latency feature, but the system is really a freshness, memory, and ranking trade-off machine.

**Type:** Build
**Company focus:** Google
**Learning goal:** Design an autocomplete system that separates candidate generation from ranking, manages memory-heavy prefix indexes, and explains freshness and personalization costs clearly.
**Prerequisites:** `05-storage-indexing-and-access-patterns/03-indexes`, `06-caching-and-invalidation/07-cache-consistency`, `16-application-backends/02-news-feed`
**Estimated time:** ~75 min
**Primary artifact:** autocomplete-config validator + trade-off matrix

## The Problem

Design search autocomplete for a large query workload. As users type prefixes, the system should return low-latency suggestions based on popularity, freshness, language, and possibly personalization.

This lesson matters because many answers stop at "use a trie." Senior answers go further: memory layout, update pipelines, abusive query suppression, multi-language segmentation, ranking freshness, and fallback behavior during partial outages.

## Clarify

- Is autocomplete based only on historical query popularity, or also on documents and trending events?
- How fresh do updates need to be after a new trend appears?
- Do we personalize results per user, per region, or not at all?
- Are adult, abusive, or policy-blocked queries filtered from suggestions?

If details are left open, assume query-based autocomplete with region-aware popularity, trend updates within a few minutes, light personalization, and policy filtering before serving.

## Requirements

### Functional

- Return top suggestions for each typed prefix.
- Support language or region-aware popularity.
- Filter blocked or low-quality suggestions.
- Update rankings as query popularity shifts.
- Fall back gracefully when fresh signals are delayed.

### Non-functional

- Keep p99 suggestion latency under 50 ms.
- Bound memory growth for prefix indexes.
- Avoid serving dangerous or policy-blocked suggestions.
- Allow popularity refresh without rebuilding the entire index synchronously.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Query request rate | 1.2M req/s peak | latency path must stay memory-local |
| Typed keystrokes per session | 8 average | amplifies backend load per user action |
| Distinct daily queries | 500M | affects candidate corpus and long-tail storage |
| Suggestion payload | top 8 results | limits network and ranking work |
| Freshness target | under 5 minutes for trends | shapes update ingestion and cache invalidation |

## Architecture

```text
clients
  -> edge cache
  -> autocomplete serving tier
  -> in-memory prefix index
  -> ranking features cache

query logs + trend signals
  -> aggregation pipeline
  -> candidate builder
  -> index snapshot publisher
  -> staged rollout to serving shards
```

Design notes:

1. Generate candidates cheaply from prefix indexes, then rank a small set with fresher features.
2. Keep a stable snapshot for serving and roll forward snapshots gradually.
3. Separate policy filtering from popularity scoring so unsafe suggestions never enter the fast path.
4. Plan for a stale-but-safe fallback if trend updates fail.

## Data Model & APIs

Core structures:

```text
prefix
candidate_query
region
language
base_popularity
trend_score
policy_state
last_updated_at
```

Useful interfaces:

- `GET /v1/autocomplete?prefix=dis&region=pl&lang=en`
- `PublishSnapshot(snapshot_id, shard, created_at)`
- `BlockSuggestion(query, reason, policy_version)`

Autocomplete answers are often stronger when they name the data ownership boundary between logging, ranking, and serving.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| trend pipeline stalls | snapshot age and feature lag alerts | serve last good snapshot and degrade freshness messaging internally |
| one prefix family becomes extremely hot | prefix hot-key metrics and shard CPU saturation | edge caching, prefix replication, and hot-prefix pinning |
| policy service fails open | blocked suggestion count drop and policy check health | fail closed for uncertain suggestions and keep last good blocklist |
| full snapshot rollout corrupts one shard | canary mismatch and serving result diffing | blue/green snapshots and staged rollout rollback |

## Observability

- metric: latency by prefix length, region, and cache hit level
- metric: snapshot age and publish success by shard
- metric: top hot prefixes as share of request volume
- metric: blocked suggestion serve attempts and policy decision latency
- log: selected suggestions with snapshot ID, region, and fallback reason
- trace: serving request through candidate generation and ranking feature lookup
- SLO: 99.9% of autocomplete requests return within latency target using an approved snapshot

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| snapshot-based serving index | predictable low latency | freshness is bounded by publish cadence | online mutation on every keystroke |
| two-stage candidate then ranking | better ranking quality with bounded compute | extra pipeline complexity | ranking the full corpus per request |
| stale-safe fallback snapshot | preserves availability | trend freshness can lag | hard dependency on real-time feature pipeline |

## Interview It

**Google framing:** "Design autocomplete for a global search product." Expect follow-ups on hot prefixes, language segmentation, and how trend freshness affects ranking.

**Cloudflare framing:** "Design low-latency prefix suggestions served near users." Expect questions on edge caching, global rollout, and safe degradation when origin ranking signals lag.

**Follow-ups:**
1. What changes if suggestions must include document titles, not just prior queries?
2. How do you handle brand-new trending phrases with little historical data?
3. How do you support per-user personalization without exploding memory?
4. What if one region requires stricter policy filtering than others?
5. How would you test ranking changes before global rollout?

## Ship It

- `outputs/tradeoff-matrix-search-autocomplete.md`

## Exercises

1. **Easy** — Estimate the extra request load caused by per-keystroke autocomplete.
2. **Medium** — Compare snapshot publishing against per-record online updates.
3. **Hard** — Redesign the system so one set of suggestions is personalized while another remains policy-reviewed and globally shared.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — good intuition for keystroke-driven latency sensitivity
- [Google Search Central: Search features and ranking systems](https://developers.google.com/search/docs/fundamentals/ranking-systems-guide) — useful framing for freshness and ranking signals
