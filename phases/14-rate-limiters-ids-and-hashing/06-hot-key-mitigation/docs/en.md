# Hot-Key Mitigation Strategies

> Average load does not melt systems. Concentrated load does.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Recognize when skewed keys dominate a distributed design, then choose mitigation patterns such as replication, request coalescing, isolation, and admission control without breaking correctness.  
**Prerequisites:** `06-caching-and-invalidation/04-cache-stampede`, `09-partitioning-sharding-and-rebalancing/02-hot-partitions`, `14-rate-limiters-ids-and-hashing/04-consistent-hashing`  
**Estimated time:** ~60 min  
**Primary artifact:** hot-key failure checklist + interview card  

## The Problem

Design a mitigation plan for a system where a tiny portion of keys carries a large portion of read or write traffic. Hashing and sharding distribute the average case, but a single celebrity user, API key, or configuration object can still overload one owner path.

## Clarify

- Is the hot key read-heavy, write-heavy, or both?
- Can the object be replicated safely, or does one owner need to serialize updates?
- Is the hotspot short-lived, predictable, or permanent?
- Would stale reads be acceptable during the spike?

If details are missing, assume read-heavy hotspots with occasional write bursts and strict protection of downstream storage.

## Requirements

### Functional

- Detect and mitigate skewed traffic concentrated on a few keys.
- Preserve correctness for writes and bounded staleness for reads.
- Avoid letting one hotspot degrade unrelated tenants.

### Non-functional

- Keep mitigation simple enough to activate during incidents.
- Bound the blast radius of one hot key.
- Preserve observability into why traffic shifted.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Normal traffic | 150K req/s | baseline looks safe |
| Hot-key spike | 40K req/s on one key | one shard can fail while fleet average stays fine |
| Cacheable share | 90% reads | replication and coalescing can pay off |
| Write amplification tolerance | low | some mitigations are read-only friendly |
| Rough cost | extra replicas + selective isolation | skew handling is cheaper than overprovisioning the whole fleet |

## Architecture

```text
detect hotspot
  -> classify read-heavy or write-heavy
  -> choose mitigation
     -> replicate or fanout cache
     -> request coalescing
     -> isolate to dedicated capacity
     -> tighter admission control
```

Useful strategy order:

1. Detect the skew quickly.
2. Separate read and write mitigation paths.
3. Protect the backing owner before global fleet latency degrades.
4. Keep the mitigation reversible.

## Data Model & APIs

Useful entities:

- `HotKeyProfile`
- `ReplicationPolicy`
- `AdmissionPolicy`
- `CoalescingGroup`

Useful APIs:

- `DetectHotKey(key)`
- `EnableMitigation(key, mode)`
- `ExplainMitigation(key)`
- `DisableMitigation(key)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| consistent hashing evenly places cold keys but one hot key still melts an owner | one-node saturation with normal fleet average | hot-key isolation or replication |
| replicated reads serve stale data too long | freshness breach metrics | bounded TTLs and explicit write invalidation |
| mitigation is activated globally instead of selectively | broad cost spike or noisy blast radius | target per key or key class only |
| write-heavy hotspot is treated like a cache problem | conflicts or serialized backend collapse | separate read replicas from write admission control |

## Observability

- metric: top keys by QPS and by backend cost
- metric: owner-node saturation versus fleet average
- metric: mitigation activation count and duration
- metric: coalescing savings or replica hit ratio
- log: sampled explanation of which mitigation was enabled and why

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| replicate hot read keys | spreads load quickly | stale-read and invalidation complexity | keep a single owner forever |
| request coalescing | collapses stampedes on misses | adds wait-path logic | allow every miss to hit backend |
| dedicated isolation for top tenants or keys | strong blast-radius control | higher reserved capacity cost | whole fleet absorbs pathological hotspots |

## Interview It

**Google framing:** "What do you do when a few keys dominate traffic?" Expect follow-ups on cost and correctness under mitigation.

**Cloudflare framing:** "How do you protect an edge-adjacent system from a global hot object?" Expect questions about fast activation, cache behavior, and tenant isolation.

**Follow-ups:**
1. What changes for write-heavy hotspots?
2. When is replication worse than admission control?
3. How do you know the hotspot is over?
4. How do you keep one tenant from consuming shared burst capacity?
5. What would you prebuild before the incident happens?

## Ship It

- `outputs/failure-checklist-hot-key-mitigation.md`
- `outputs/interview-card-hot-key-mitigation.md`

## Exercises

1. **Easy** — Pick a mitigation for a read-heavy celebrity profile page.
2. **Medium** — Redesign for a write-heavy shared counter.
3. **Hard** — Support a sudden global hot key while preserving fairness for other tenants.

## Further Reading

- [AWS Builders Library - Caching challenges and strategies](https://aws.amazon.com/builders-library/caching-challenges-and-strategies/) — practical hotspot and stampede lessons  
- [Google SRE - Handling overload](https://sre.google/sre-book/handling-overload/) — helpful when hot keys turn into overload events  
