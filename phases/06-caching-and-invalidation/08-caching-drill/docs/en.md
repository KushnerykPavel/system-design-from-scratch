# Caching Design Drill

> The strongest caching answer sounds like a freshness contract, not a product list.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Practice compressing cache pattern, freshness, stampede prevention, and consistency reasoning into one coherent interview deep dive.
**Prerequisites:** `01-cache-patterns`, `02-freshness-models`, `04-cache-stampede`, `07-cache-consistency`
**Estimated time:** ~60 min
**Primary artifact:** drill worksheet + scoring rubric

## The Problem

This drill pulls together the whole caching phase. The learner must take an ambiguous prompt and produce a cache design that is fast, bounded in staleness, and operationally credible under failure.

The goal is not to mention every caching concept. The goal is to make a few clean, defensible choices:

- where to cache
- how entries become fresh or stale
- how hot keys behave on miss or expiry
- what users see after writes and failures

## Clarify

- Which objects are hottest and which are most mutable?
- Is freshness more important than offload, or vice versa?
- Which user journey most needs read-after-write correctness?
- What cache layers are allowed: application, shared cache, CDN, browser, edge?

## Requirements

### Functional

- Choose a cache pattern and freshness model for the main path.
- Explain invalidation or expiry for mutable data.
- Handle hot keys and miss storms safely.

### Non-functional

- Keep latency and origin offload visible.
- Make stale-read risk explicit rather than accidental.
- Show observability and degraded-mode thinking.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak reads | 300K req/s | large enough that caching is central, not optional |
| Writes | 6K req/s | enough to make freshness choices matter |
| Hot-key skew | 1% of keys cause 45% of reads | exposes stampede and eviction reasoning |
| Freshness target | sub-second to minutes depending on data class | forces prioritization instead of one TTL |
| Rough cost | memory, invalidation traffic, origin misses | keeps the answer practical |

## Architecture

Recommended drill sequence:

1. Clarify freshness and data classes.
2. Choose one primary cache path.
3. Add invalidation or revalidation only where it changes outcomes.
4. Explain stampede prevention and degraded mode.
5. Close with the user-visible consistency story.

The answer should sound like an intentional serving policy, not "cache everywhere."

## Data Model & APIs

A strong drill answer usually names:

- cache key composition
- metadata such as TTL or version
- invalidation trigger
- critical-path read API

Example:

```text
Get(key, min_version?)
Invalidate(key, version)
Warm(keys)
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hot entry expires and causes origin surge | hit ratio collapse with origin QPS spike | request coalescing, jitter, warming |
| user sees stale data after mutation | read-after-write mismatch | write-through or selective cache bypass |
| invalidation path drops events | version skew and stale-age alerts | replayable invalidation plus TTL safety net |
| cache layer leaks private data | tenant/context mismatch in served entries | key discipline and no-store on unsafe responses |

## Observability

- metric: hit ratio by layer and entity class
- metric: origin offload and miss amplification
- metric: invalidation latency and stale-read incidents
- log: cache-key decisions, bypasses, and rejected invalidations
- trace: one request across cache lookup, origin fetch, and refill
- SLO: latency objective plus freshness objective for the main user journey

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one explicit freshness contract | keeps answer coherent | less breadth across edge cases | many disconnected cache rules |
| selective strictness on critical reads | balances correctness and cost | requires entity classification | uniform strongest consistency |
| layered observability | makes debugging realistic | more instrumentation work | "we will monitor it" with no specifics |

## Interview It

**Google framing:** "Design caching for a collaborative notes feed." The signal is whether you separate immutable assets, feed fanout reads, and immediate read-after-write for the editor.

**Cloudflare framing:** "Design caching for global API responses with mutable policy data." The signal is whether you reason about edge layers, purge, stale policy, and safe bypasses.

**Follow-ups:**
1. What if the interviewer says freshness now matters more than cost?
2. What if one celebrity account causes 30x traffic on a single object?
3. What if users complain they sometimes see older state after moving regions?
4. What if the cache cluster fails and the origin is undersized for full traffic?
5. Which deep dive would you pick if the interviewer says "go one level deeper"?

## Ship It

- `outputs/drill-worksheet-caching.md`
- `outputs/scoring-rubric-caching.md`

## Exercises

1. **Easy** — Run the drill for a mostly static content feed.
2. **Medium** — Run the drill for product pages with rapidly changing inventory.
3. **Hard** — Run the drill for globally distributed policy reads where stale data can under-enforce security controls.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — helpful baseline for structuring the overall interview answer
- [Caching at Scale at Facebook](https://engineering.fb.com/2013/04/15/core-infra/scaling-memcache-at-facebook/) — practical examples to compare against your drill choices
