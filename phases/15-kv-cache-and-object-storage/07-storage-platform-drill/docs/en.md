# Storage Platform Drill

> The point of the drill is not to name every storage primitive. It is to choose which promises matter and defend them under changing constraints.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Practice stitching together KV, cache, object, metadata, durability, and lifecycle concepts into one coherent storage-platform answer under interview pressure.  
**Prerequisites:** `15-kv-cache-and-object-storage/01-distributed-kv-store`, `15-kv-cache-and-object-storage/03-object-storage`, `15-kv-cache-and-object-storage/05-durability-tiers`, `15-kv-cache-and-object-storage/06-compaction-and-lifecycle`  
**Estimated time:** ~60 min  
**Primary artifact:** full drill worksheet  

## The Problem

Design a storage platform for a product that stores user-uploaded assets, exposes fast metadata listing, serves frequently accessed derived assets with caching, and offers multiple durability classes. The interviewer cares less about vendor names and more about whether your design is internally consistent.

This drill is meant to force explicit prioritization:

- what is cached versus durably stored
- what is metadata versus blob data
- which data classes get which durability promise
- how repair and lifecycle work stay visible

## Clarify

- Which user flows dominate: upload, retrieval, listing, or policy actions?
- What durability tiers are product-visible versus internal?
- How fresh must metadata and derived assets be?
- Which failures are acceptable to degrade, and which are not?

If the interviewer stays broad, assume user-uploaded media with tenant-scoped listings, hot derived thumbnails, and archival retention for originals.

## Requirements

### Functional

- Upload and retrieve durable objects.
- List and filter metadata quickly.
- Cache hot derived assets.
- Support retention, deletion, and storage-class policies.

### Non-functional

- Keep upload and retrieval latency predictable under spikes.
- Bound metadata staleness and cache inconsistency.
- Restore degraded durability within documented time.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Upload QPS | 40K objects/s peak | drives ingest, multipart, and metadata finalize |
| Read mix | 95% cached thumbnail reads, 5% original object fetches | determines cache versus object-path focus |
| Metadata QPS | 70K req/s | listing plane cannot hide behind blob storage |
| Storage footprint | 15 PB hot and warm, 60 PB archival | durability classes and lifecycle costs dominate |
| Peak factor | 5 to 8x during campaign or backup windows | stress-tests cache warmup and maintenance backlogs |

## Architecture

```text
upload clients
  -> auth + upload session
  -> object store
  -> metadata finalize
  -> async derivation jobs
  -> cache warm path for hot derivatives

read clients
  -> metadata/list API
  -> derived asset cache
  -> object storage for originals

background control
  -> durability auditor
  -> repair workers
  -> lifecycle and retention engine
```

Suggested deep-dive choices:

1. upload finalize and metadata correctness
2. durability tiering and repair scheduling
3. cache stampede protection for derived assets

## Data Model & APIs

Key entities:

```text
object
object_version
metadata_record
derived_asset
durability_tier
lifecycle_rule
```

Useful APIs:

- `StartUpload`
- `CompleteUpload`
- `ListObjects`
- `GetDerivedAsset`
- `ChangeDurabilityTier`
- `DeleteObject`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| upload completes but metadata is missing | finalize gap and orphan metrics | staged finalize plus reconciliation |
| hot derivative cache flushes during traffic spike | hit-rate cliff and origin surge | request coalescing, prewarm, and load shedding |
| durability tier migration stalls | objects-in-transition age | resumable migration workflow with audit trail |
| lifecycle engine violates retention rule | delete-policy audit alarms | dry-run, rule versioning, and approval gates |

## Observability

- metric: end-to-end upload success and finalize delay
- metric: derived asset cache hit rate and origin offload
- metric: metadata freshness lag and listing latency
- metric: under-durable object age and repair backlog
- metric: lifecycle action backlog and destructive-action audit volume
- trace: upload to derivative-ready timeline

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| separate metadata and object planes | scalable query performance | more finalize complexity | direct blob scans for listings |
| dedicated derived-asset cache | protects origin and object store | invalidation and warmup complexity | serving all reads from object storage |
| multiple durability classes | cost aligned to value | more explanation and migration tooling | one expensive tier for everything |

## Interview It

**Google framing:** "Design a shared media storage platform used by multiple product teams." Expect pushback on index correctness, durability guarantees, and cost.

**Cloudflare framing:** "Design a globally distributed asset storage and delivery control plane." Expect follow-ups on cache locality, object metadata, and regional failure behavior.

**Follow-ups:**
1. What changes if legal hold becomes mandatory for a subset of tenants?
2. How do you keep metadata listings usable during a partial object-store outage?
3. What if derived assets are generated asynchronously and must be cached at the edge?
4. How do you migrate one tenant to a stronger durability tier with no downtime?
5. What changes when cost pressure forces colder storage classes much earlier?

## Ship It

- `outputs/storage-platform-drill-sheet.md`
- `outputs/scoring-rubric-storage-platform-drill.md`

## Exercises

1. **Easy** — Choose one deep dive and outline the clarifying questions you would ask first.
2. **Medium** — Redesign the drill for internal backup storage instead of user media.
3. **Hard** — Re-answer the drill after the interviewer cuts cost by 40% but tightens compliance.

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline interview flow for storage-centric prompts  
- [Amazon Dynamo paper](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) — helpful when reasoning about KV components inside broader storage systems  
