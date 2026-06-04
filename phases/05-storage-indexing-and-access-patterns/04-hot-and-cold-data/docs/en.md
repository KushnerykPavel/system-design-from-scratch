# Hot Rows, Cold Data, and Tiering

> Most storage pain comes from a tiny fraction of data being touched constantly.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Recognize hotspots, separate hot from cold data, and explain when tiering is the right move for latency, cost, and operational safety.  
**Prerequisites:** `02-access-pattern-first`, `03-indexes`  
**Estimated time:** ~75 min  
**Primary artifact:** failure checklist + interview card  

## The Problem

Data is rarely accessed uniformly. A handful of tenants, users, keys, or recent objects often receive most of the traffic, while the majority of data sits cold but still consumes storage and operational budget.

This lesson teaches you to distinguish:

- **hot rows or keys** that dominate read or write pressure
- **warm data** that still matters for interactive requests
- **cold data** that should move to cheaper or slower tiers

## Clarify

- Is the hotspot caused by one tenant, one object, recent time windows, or one write-heavy counter?
- Does hotness matter more for reads, writes, or both?
- What latency target applies to recent versus historical data?
- Can older data tolerate archive retrieval times or asynchronous restore steps?

## Requirements

### Functional

- Identify where the hot working set lives.
- Separate serving paths for hot and cold data when justified.
- Explain how data moves between tiers safely.

### Non-functional

- Keep the hot path low latency under skew.
- Reduce storage cost for low-value cold data.
- Avoid migrations that silently break compliance or restore expectations.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Working set | 3% of rows serve 85% of reads | proves hot/cold separation value |
| Historical data | 500 TB retained for 2 years | tiering materially changes cost |
| Write hotspot | top 0.1% of keys see 40% of writes | exposes lock contention and skew |
| Peak factor | 8x during incidents or launches | stresses hotspot mitigation |
| Rough cost | SSD hot tier + HDD or object cold tier | cost difference motivates lifecycle policies |

## Architecture

Common pattern:

```text
request
  -> hot serving store / cache / recent partition
  -> warm replica or bounded historical store
  -> archive tier for rare restores
```

The design is incomplete unless you explain:

1. how data qualifies as hot or cold
2. when it moves
3. how it is restored if needed

## Data Model & APIs

Useful patterns:

- recent data partitioned by time for fast current reads
- archived blobs stored separately with lightweight metadata pointers
- hot counters split or aggregated to avoid one write hotspot

Possible APIs:

- `GetRecentEvents(user_id, cursor)`
- `RequestArchiveExport(range)`
- `RestoreArchivedObject(object_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one tenant or row dominates traffic | per-key skew metrics and lock waits spike | shard the hotspot, cache aggressively, or precompute summaries |
| cold tier restore path never tested | restore jobs fail during audit or incident | run scheduled restore drills with time objectives |
| tiering policy moves data too early | user-facing latency jumps on recent lookups | delay demotion or keep a larger warm tier |
| archived data loses metadata linkage | object exists but cannot be found or restored | keep durable metadata catalog with integrity checks |

## Observability

- metric: read and write skew by tenant, key class, or time window
- metric: hot tier hit rate and cold tier retrieval count
- metric: archive restore success rate and restore latency
- log: lifecycle movement decisions and restore failures
- trace: which tier served the request
- SLO: recent data path meets latency target while archive restores meet a documented recovery target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit hot and cold tiers | better latency and lower cost | more lifecycle logic and restore complexity | one uniform premium tier for everything |
| recent-time partitioning | fast common reads and simpler retention | awkward cross-window queries | one giant unbounded table |
| archive restore workflow | much cheaper long-term storage | slower access during rare historical reads | keeping years of data in the serving store |

## Interview It

**Google framing:** "Design storage for activity history where recent reads matter most." The signal is whether you model the working set instead of assuming all data is equal.

**Cloudflare framing:** "Design a control-plane audit history with rare long-term retrieval." The signal is whether you pair cheap retention with a credible restore and compliance story.

**Follow-ups:**
1. What if the top customer is 50% of all traffic?
2. What if historical queries suddenly become a paid feature?
3. How do you prove archived data is actually restorable?
4. What if data temperature changes after an incident or product launch?
5. When is caching enough, and when do you need true storage tiering?

## Ship It

- `outputs/failure-checklist-hot-and-cold-data.md`
- `outputs/interview-card-hot-and-cold-data.md`

## Exercises

1. **Easy** — Define hot, warm, and cold tiers for a photo-sharing product.  
2. **Medium** — Redesign a log store where old data is too expensive but occasionally needed within 24 hours.  
3. **Hard** — Explain how to handle one customer creating most of the write pressure in a multitenant system.  

## Further Reading

- [Bigtable: A Distributed Storage System for Structured Data](https://research.google/pubs/pub27898/) — helpful for hotspot and tablet locality discussions  
- [Amazon S3 Storage Classes](https://aws.amazon.com/s3/storage-classes/) — concrete examples of tiering trade-offs and restore behavior  
