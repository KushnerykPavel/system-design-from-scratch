# Metadata Index Service for Large Blobs

> Blob bytes are not your query engine. Metadata is a product in its own right.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Design a metadata index service that makes large-blob platforms searchable, policy-aware, and operationally tractable without coupling every query to object storage.  
**Prerequisites:** `05-storage-indexing-and-access-patterns/03-indexes`, `05-storage-indexing-and-access-patterns/05-blob-metadata-separation`, `15-kv-cache-and-object-storage/03-object-storage`  
**Estimated time:** ~60 min  
**Primary artifact:** design review checklist  

## The Problem

An object store alone is rarely enough. Products want listing by tenant, filtering by type, recent uploads, moderation state, retention class, or billing tags. That requires a metadata index service with explicit ownership, schema evolution, and backfill tooling.

The hard part is not inventing indexes. It is defining which queries matter, how metadata is updated safely, and how stale or missing index rows are detected.

## Clarify

- Which metadata lookups are truly product-critical versus operator-only?
- Are queries mostly point lookups, prefix listings, time-window scans, or rich filtering?
- Can metadata become visible before background enrichment finishes?
- What is the correctness requirement when metadata and object state disagree temporarily?

If the prompt is ambiguous, assume tenant-scoped listings, time-sorted recent objects, filtering by type and state, and eventual enrichment from background jobs.

## Requirements

### Functional

- Serve point lookups and tenant-scoped listings for objects.
- Support filter dimensions such as type, status, storage class, or retention state.
- Ingest metadata updates from upload finalize and asynchronous enrichment jobs.

### Non-functional

- Keep listing latency predictable without scanning blob storage.
- Make index lag and backfill correctness measurable.
- Support schema growth without breaking older producers and consumers.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Metadata write rate | 120K updates/s peak | drives ingestion, index maintenance, and backpressure |
| Query rate | 80K req/s | index design must match dominant filters and sort orders |
| Record size | 1 to 3 KB | index amplification matters more than row count alone |
| Retention | 7 years for audit metadata, 90 days hot listing index | hot and cold paths should differ |
| Peak factor | 5x on batch imports | backfill and online ingest must coexist |

## Architecture

```text
upload finalize / enrichment jobs
  -> metadata event or write API
  -> canonical metadata store
  -> secondary listing indexes / search projections
  -> query API
```

Useful pattern:

1. One canonical metadata owner records state transitions.
2. Secondary indexes or projections are built for actual query shapes.
3. Backfill and replay are first-class because schemas and filters evolve.

## Data Model & APIs

Canonical metadata:

```text
object_id
tenant_id
bucket
object_state
content_type
created_at
retention_class
moderation_state
tags
```

Useful APIs:

- `GetObjectMetadata(object_id)`
- `ListObjects(tenant_id, filter, cursor)`
- `UpdateObjectState(object_id, expected_version, patch)`
- `ReplayIndex(from_offset)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| listing index lags canonical metadata badly | index lag metrics and stale result audits | replay pipeline and consumer lag alerts |
| enrichment pipeline overwrites newer state | version conflict counters | optimistic concurrency and monotonic state transitions |
| query surfaces unsupported filter combinations | slow scans and high tail latency | explicit query contracts and indexed access patterns |
| backfill floods hot partitions | ingest latency spikes during replay | throttled replay and isolated backfill workers |

## Observability

- metric: canonical write latency and conditional update conflicts
- metric: per-index lag from canonical source of truth
- metric: query latency by filter shape and sort order
- metric: backfill throughput and replay staleness
- log: rejected state transitions and schema-version mismatch events
- SLO: tenant listing remains fast while index lag stays within documented bounds

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| canonical store plus projections | simpler ownership and query tuning | eventual consistency between views | one over-generalized database for every access pattern |
| strict query contract | predictable performance | less ad hoc flexibility | arbitrary filters that degrade into scans |
| versioned metadata writes | safer concurrent updates | extra API complexity | blind overwrite from every async worker |

## Interview It

**Google framing:** "Design the searchable metadata plane for durable object storage." Expect questions about indexes, replay, and schema evolution.

**Cloudflare framing:** "Design a control-plane index for large assets distributed through the edge." Expect questions about listing latency, async enrichment, and safe rollouts of new metadata fields.

**Follow-ups:**
1. What if customers need cross-bucket search by tag?
2. How do you roll out a new filter without rewriting the whole store?
3. What if enrichment pipelines arrive out of order?
4. How do you keep audit records longer than hot listing indexes?
5. What changes when index lag becomes customer-visible?

## Ship It

- `outputs/design-review-metadata-index-service.md`
- `outputs/interview-card-metadata-index-service.md`

## Exercises

1. **Easy** — Choose the first three filter dimensions you would index and justify them.
2. **Medium** — Design a replay plan for rebuilding a corrupted secondary index.
3. **Hard** — Extend the metadata plane to support search by tags and retention policy without turning every query into a scan.

## Further Reading

- [The Log: What every software engineer should know about real-time data's unifying abstraction](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying) — strong framing for replayable index pipelines  
- [System design notes](https://github.com/liquidslr/system-design-notes) — useful canonical interview flow for explaining secondary indexes and query-oriented data design  
