# Read-Through, Write-Through, and Write-Behind

> A cache pattern is a correctness choice first and a latency choice second.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Choose the right cache integration pattern for a workload by reasoning from source of truth, write path risk, and freshness tolerance.
**Prerequisites:** `04-apis-contracts-and-schema-evolution/02-idempotency-keys`, `05-storage-indexing-and-access-patterns/02-access-pattern-first`
**Estimated time:** ~75 min
**Primary artifact:** cache pattern decision matrix

## The Problem

Many interview answers say "add Redis" without explaining how reads and writes interact with it. That skips the real question: who owns correctness, who updates the cache, and what failure do users see when the cache and source of truth disagree?

This lesson covers the three patterns that appear constantly in senior-level design discussions:

- read-through
- write-through
- write-behind

The goal is not to memorize definitions. The goal is to match each pattern to workload shape and operational risk.

## Clarify

- Is the cache only accelerating reads, or is it also part of the write acknowledgement path?
- What is the source of truth: database, object store, search index, or a derived view?
- Can users tolerate stale reads, delayed writes, or dropped async updates?
- Is the workload dominated by read amplification, write amplification, or burst fanout after a write?

## Requirements

### Functional

- Support a hot read path with lower latency than the primary store.
- Keep writes durable in the system of record.
- Make cache misses and repopulation behavior explicit.

### Non-functional

- Avoid hidden data-loss modes on the write path.
- Keep failure behavior explainable during incidents.
- Prefer patterns that fit the dominant access pattern instead of mixing all of them blindly.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak reads | 180K req/s | high enough that cache placement materially changes origin load |
| Peak writes | 12K req/s | enough to make synchronous cache updates a real cost |
| Hot key fanout | 20K reads in 30 seconds after an update | exposes invalidation and stale-read pressure |
| Freshness target | 1-30 seconds depending on path | determines whether async cache population is acceptable |
| Rough cost | cache memory + write amplification + miss traffic | keeps the design grounded beyond latency alone |

## Architecture

### Read-through

The application checks the cache first. On miss, it reads from the source of truth, returns the response, and fills the cache.

Use it when:

- read traffic is much larger than write traffic
- misses are acceptable but should be amortized
- the application already owns the data access path

Risk:

- stampedes on hot misses
- inconsistent serialization across callers if fill logic is duplicated

### Write-through

The application writes to the source of truth and cache in the same logical path, usually before acknowledging success to the caller.

Use it when:

- read-after-write freshness matters
- the cache entry is cheap to derive at write time
- write latency can afford the extra step

Risk:

- higher write latency
- partial failure handling becomes more complex

### Write-behind

The application writes to the cache or durable queue first and pushes to the backing store asynchronously.

Use it when:

- writes arrive at very high rate
- batching or smoothing storage traffic matters
- the system tolerates delayed persistence and replay-driven recovery

Risk:

- async loss or duplication if the write-behind path is not durable
- harder correctness story in interviews unless the invariants are carefully bounded

## Data Model & APIs

Useful cache entry fields:

```text
key -> {
  value,
  version,
  fetched_at,
  ttl_seconds
}
```

Representative APIs:

- `Get(key)` for read-through
- `Put(key, value, version)` for explicit population or write-through
- `EnqueueWrite(record)` for write-behind pipelines

The most important design statement is which store is authoritative. The cache should not quietly become a second database unless the lesson explicitly calls for it.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| cache miss storm after expiry | origin QPS spike and drop in cache hit rate | request coalescing, jittered TTLs, warmup for hot keys |
| cache updated but DB write fails | mismatch between cache version and durable version | write ordering rules plus compensating invalidation |
| async write-behind worker falls behind | queue lag and age-of-oldest-write alarms | durable queue, backpressure, replay, degraded write mode |
| stale object survives after source update | version skew or user-visible stale-read reports | versioned keys, explicit invalidation, bounded TTL |

## Observability

- metric: cache hit ratio by endpoint and key class
- metric: write latency split by source-of-truth write and cache update
- metric: write-behind queue lag and retry rate
- log: sampled cache fill, invalidate, and async write failure events
- trace: cache lookup, source read, and cache population on one span tree
- SLO: hot read path latency with a paired freshness target for user-visible data

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| read-through for most serving paths | simple mental model and good miss amortization | miss penalty still hits origin | always reading origin directly |
| write-through for strict read-after-write views | tighter freshness after mutation | higher write latency and more partial-failure logic | TTL-only freshness |
| write-behind only on explicitly async domains | absorbs bursty writes and enables batching | delayed durability semantics to explain | pretending every write path can be async |

## Interview It

**Google framing:** "Design profile serving for a consumer app where reads are much hotter than writes." The signal is whether you can explain when write-through is worth the latency tax and when read-through is enough.

**Cloudflare framing:** "Design a configuration cache in front of control-plane storage." The signal is whether you separate fast serving from authoritative config state and reason about stale policy risk.

**Follow-ups:**
1. What if users must see their own writes immediately?
2. What if one update causes millions of reads in the next minute?
3. What if the cache cluster is healthy but the backing store is degraded?
4. What if the cache is cheaper to update than to invalidate?
5. Which pattern changes most if the system becomes multi-region?

## Ship It

- `outputs/cache-pattern-decision-matrix.md`

## Exercises

1. **Easy** — Choose a cache pattern for product catalog reads where updates happen every few minutes.
2. **Medium** — Choose a pattern for quota configuration where stale reads can temporarily under-enforce policy.
3. **Hard** — Explain why write-behind is dangerous for payments but acceptable for some analytics ingestion paths.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline for placing caching inside a full interview answer
- [Caching at Scale at Facebook](https://engineering.fb.com/2013/04/15/core-infra/scaling-memcache-at-facebook/) — practical examples of cache architecture and consistency pressure
