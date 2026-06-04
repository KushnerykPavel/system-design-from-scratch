# Storage Growth and Retention Math

> Storage mistakes are expensive because they accumulate quietly before they fail loudly.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Estimate raw, replicated, and retained storage growth so data model and lifecycle choices stay realistic.  
**Prerequisites:** `02-estimation-and-cost/01-qps-and-request-mix`, `05-storage-indexing-and-access-patterns/02-access-pattern-first`  
**Estimated time:** ~75 min  
**Primary artifact:** retention worksheet + interview card  

## The Problem

Engineers often quote total data size without breaking out ingestion rate, index overhead, replication, and retention policy. That leads to weak answers about storage engines, tiering, and deletion semantics.

This lesson teaches a compact way to turn event volume into daily, yearly, and replicated footprint.

## Clarify

- Are we storing raw objects, metadata, logs, or derived indexes?
- What is the retention requirement for hot, warm, and cold data?
- Is the data mutable, append-only, or periodically compacted?
- Does the number need to include replicas, backups, and indexes?

## Requirements

### Functional

- Estimate daily ingest size.
- Estimate retained footprint across time windows.
- Separate primary data from replicas, indexes, and backups.

### Non-functional

- The model should be simple enough for interview use.
- It should highlight the cost of long retention and over-replication.
- It should show when tiered storage becomes necessary.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Events per day | 400M | anchors ingest volume |
| Bytes per event | 1200 B | controls raw footprint |
| Daily raw ingest | about 447 GB | sets primary storage growth |
| Replication factor | 3x | changes durable footprint materially |
| 30-day retained size | about 39 TB before indexes | informs tiering and cost |

## Architecture

A useful decomposition is:

1. Raw ingest per day.
2. Replicated durable size.
3. Index or metadata overhead.
4. Retention window by tier.
5. Backup or compliance copy.

Example:

- 400M events/day
- 1.2 KB each
- about 447 GB/day raw
- about 1.34 TB/day at 3x replication
- about 40 TB for 30 days before secondary indexes and backups

## Data Model & APIs

The code artifact models a simple storage plan:

```text
StoragePlan {
  EventsPerDay
  BytesPerEvent
  ReplicationFactor
  RetentionDays
  IndexOverheadRatio
}
```

Outputs:

- daily raw GB
- daily durable GB
- retained total GB

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| retention ignored | cluster fills far faster than forecast | state retention in days for each tier |
| indexes forgotten | actual disk use is much higher than raw payload math | add index overhead explicitly |
| hot and cold data mixed | expensive primary disks store stale data | introduce storage tiers and lifecycle moves |
| backups omitted | recovery plan lacks capacity budget | count backups as separate copies |

## Observability

- metric: daily bytes ingested by data class
- metric: retained bytes by storage tier
- metric: index-to-primary size ratio
- metric: deletion lag against retention policy
- SLO: retained footprint should stay within forecasted growth envelope

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| separate raw data from index overhead | more honest sizing | slightly more math | pretending payload size is total size |
| retention by tier | reduces cost | adds operational lifecycle complexity | one forever-hot storage class |
| rounded storage units | faster interview flow | some unit drift | exact byte-level accounting |

## Interview It

**Google framing:** "Design an event logging pipeline with one-year retention." The signal is whether you distinguish hot queryable data from archive.

**Cloudflare framing:** "Design request-log storage for a global edge network." The signal is whether you think about sheer ingest rate, compression, and cheap retention tiers.

**Follow-ups:**
1. What changes if only 7 days need fast query and the rest can be cold?
2. What if an index adds 40% overhead?
3. What if compliance requires undeletable backup copies for 90 days?
4. What if one tenant produces half the daily ingest?

## Ship It

- `outputs/retention-worksheet-storage-growth.md`
- `outputs/interview-card-storage-growth.md`

## Exercises

1. **Easy** — Estimate retained size for 100M photos with 2 MB average size.  
2. **Medium** — Add a warm tier with 30-day hot retention and 365-day archive.  
3. **Hard** — Explain how log compaction changes the estimate for mutable key histories.  

## Further Reading

- [Google SRE book](https://sre.google/books/) — strong background on capacity planning  
- [System design notes](https://github.com/liquidslr/system-design-notes) — helpful interview framing for rough sizing  
