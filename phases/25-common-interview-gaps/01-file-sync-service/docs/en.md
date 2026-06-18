# File Sync Service (Dropbox/Drive)

> The naive answer is "upload to S3." The senior answer explains delta sync, block deduplication, conflict resolution, and why the metadata service is the hardest part.

**Type:** Learn  
**Company focus:** Google  
**Learning goal:** Design a client-server file sync system that minimizes bandwidth via block-level deduplication, handles concurrent edits across devices, and resolves conflicts without data loss.  
**Prerequisites:** `05-storage-indexing-and-access-patterns/05-blob-metadata-separation`, `08-consistency-replication-and-transactions/06-sagas`, `15-kv-cache-and-object-storage/03-object-storage`  
**Estimated time:** ~90 min  
**Primary artifact:** capacity sheet

## The Problem

Design a file synchronization service like Dropbox or Google Drive. Users install a client on multiple devices and expect files to stay consistent across all of them without explicitly managing transfers.

Mid-level candidates immediately jump to "client uploads to S3, server notifies other devices." This misses almost every interesting problem: what happens when two devices edit the same file concurrently? How do you avoid re-uploading an unchanged 10 GB file when only one paragraph changed? What happens when a device is offline for a week and then reconnects with conflicting state? Senior candidates recognize that the metadata service is the true bottleneck: it must track file versions, block inventories, and device sync cursors consistently while the blob tier handles raw byte storage.

The block-level deduplication insight is especially powerful. Files are split into content-addressed blocks (typically 4 MB). Identical blocks — even across different files or users — are stored once. A file update transmits only the changed blocks. This reduces bandwidth by an order of magnitude in typical workloads and dramatically simplifies the conflict detection model since block hashes serve as canonical state identifiers.

## Clarify

- What is the target file size ceiling? (1 GB, 100 GB, or unlimited?) This affects block sizing, multipart upload strategy, and the metadata fan-out cost per file.
- Is real-time collaborative editing (like Google Docs) in scope, or is this last-write-wins file sync with conflict copies?
- What is the consistency model for cross-device visibility: eventual (minutes of lag acceptable) or near-real-time (< 10 seconds)?

If the interviewer is vague, assume file sizes up to 5 GB, no collaborative editing (conflict copies are acceptable), and sync latency target of under 30 seconds for changes up to 100 MB.

## Requirements

### Functional

- Upload and store files; propagate changes to all registered devices for the same user.
- Download the latest version of any file to any authenticated device.
- Detect and create conflict copies when two devices modify the same file concurrently.
- Support folder operations: create, rename, delete, move.
- Track version history and support file restoration to a previous version.

### Non-functional

- Minimize bandwidth by sending only changed blocks, not full file re-uploads.
- Deduplicate identical blocks across all users to reduce storage cost.
- Survive device disconnection gracefully; sync must complete correctly after reconnect.
- Metadata API p99 under 200 ms; upload and download throughput bounded by client network, not server.
- 99.99% file durability; 99.9% service availability.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active users | 100M users, 10M daily active | determines metadata DB sizing, notification fan-out |
| Files per user | 500 files avg, 2 TB storage avg | drives object storage capacity: ~200 PB total |
| Block size | 4 MB average blocks | file → block index ratio, dedup savings estimation |
| Metadata writes | 5M file change events/hour at peak | sizing metadata store and event bus |
| Dedup ratio | 30–60% block hit rate across users | cuts object storage cost by ~50% |

## Architecture

```text
client (desktop / mobile)
  -> chunk file into 4 MB content-addressed blocks
  -> compute SHA-256 per block
  -> diff local state vs last-known server state

client -> upload API
  -> check block presence (POST /blocks/check [block_hashes])
  -> upload only missing blocks (PUT /blocks/{hash})
  -> commit file version (POST /files/{file_id}/versions)
       - records block list, file metadata, device_id, timestamp
  -> metadata service writes version record
  -> event bus publishes file_changed event

notification service
  <- event bus
  -> long-poll or WebSocket per connected device
  -> delivers delta sync cursor to each registered device

other devices
  -> receive cursor update
  -> fetch changed blocks (GET /blocks/{hash} from CDN/cache)
  -> reconstruct file from block list
```

Conflict detection: when committing a new version, the metadata service checks whether the parent version the client declared matches the current head. If not, a conflict branch is created and the client receives both versions as a conflict copy.

## Data Model & APIs

Core entities:

```text
User           { user_id, quota_bytes, plan }
Device         { device_id, user_id, last_cursor, platform }
FileRecord     { file_id, user_id, path, current_version_id, deleted_at }
FileVersion    { version_id, file_id, block_list[], size_bytes, device_id, created_at, parent_version_id }
Block          { block_hash, size_bytes, storage_key, ref_count }
SyncCursor     { cursor_id, device_id, last_seen_event_seq }
```

Key APIs:

- `POST /v1/blocks/check` — body: `{hashes: [string]}`, returns `{missing: [string]}` (avoids re-uploading known blocks)
- `PUT /v1/blocks/{hash}` — upload one block; idempotent by hash
- `POST /v1/files/{file_id}/versions` — body: `{block_list, parent_version_id, path, device_id}`; returns `{version_id, conflict?: true}`
- `GET /v1/files/{file_id}/versions/{version_id}/blocks` — returns ordered block list for download
- `GET /v1/sync/delta?cursor={cursor_id}` — returns file change events since cursor; used by all devices for incremental sync

Block storage uses content-addressing: the storage key is the block hash, enabling cross-user deduplication at the object layer. Ref-counting enables safe GC when no file version references a block.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Client goes offline mid-upload | upload session expires, block upload incomplete | resumable upload tokens with expiry; client retries from last committed block |
| Concurrent edit creates conflict | parent_version_id mismatch on commit | metadata service creates conflict copy branch; both devices receive conflict notification |
| Block storage node unavailable | object store health probe fails, upload error rate spikes | route uploads to healthy replica; reads served from CDN; block is replicated 3x |
| Metadata DB overload during reconnect surge | query latency p99 spike, connection pool exhaustion | rate-limit sync requests per user; stagger reconnect backoff; read-replicas for cursor queries |

## Observability

- metric: delta sync latency (time from file save on device A to availability on device B), p50/p95/p99
- metric: block deduplication ratio (blocks_skipped / blocks_checked) — signals storage savings and health of content-address check
- metric: conflict copy creation rate — rising rate indicates a multi-device edit pattern the UX should address
- log: every version commit with user_id, file_id, block_count, conflict flag, device_id, and duration
- trace: full upload path from client block check through metadata commit and notification dispatch
- SLO: 99.9% of sync delta queries resolve within 500 ms; 99.99% of committed blocks survive 30 days

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Content-addressed block storage with cross-user dedup | cuts storage by 30–60%; delta sync only sends changed blocks | block-level index per file adds metadata complexity; GC ref-counting is tricky | file-level versioning with full uploads each time — untenable at scale |
| Conflict copy instead of operational transform | simple to implement; no lock contention across devices | users see duplicate files; requires manual merge | full OT/CRDT — appropriate for collaborative editing, overkill for async file sync |
| Event-driven delta cursor for device notification | decouples server fanout from upload latency; scales notification separately | eventual consistency between metadata write and device notification | polling on a fixed interval — wastes bandwidth and adds unnecessary latency |

## Interview It

**Google framing:** "Design Drive for 1 billion users. How do you ensure a change on one device appears on another within 30 seconds even for large files?" Expect pushback on metadata consistency, deduplication strategy, and quota enforcement at scale.

**Cloudflare framing:** "Design the block storage and sync delivery layer for a file sync product. Where does the edge layer help, and where does it hurt?" Expect questions on CDN suitability for small metadata vs large block reads and cache invalidation for versioned content.

**Follow-ups:**
1. How would you enforce per-user storage quotas atomically without a distributed transaction on every upload?
2. What changes if you need block-level encryption per user (so deduplication across users is impossible)?
3. How would you handle a 100 GB file where only a 1 KB header changes?
4. How would you detect and prevent a client bug that causes a sync loop (file keeps being re-uploaded)?
5. What is your rollback plan if a metadata migration corrupts version history for 5% of users?

## Ship It

- `outputs/capacity-sheet-file-sync-service.md`

## Exercises

1. **Easy** — Estimate total object storage capacity for 500 M users with 1 TB each at a 40% dedup ratio. How many object storage nodes at 10 TB each?
2. **Medium** — Design the GC pipeline that safely deletes unreferenced blocks. How do you avoid a race between a new version commit and a concurrent GC pass?
3. **Hard** — Add end-to-end encryption where only the client holds the key. Which deduplication properties survive? What new problems does this introduce for quota enforcement and metadata search?

## Further Reading

- https://dropbox.tech/infrastructure/magic-pocket-infrastructure — Dropbox's block storage system explaining content-addressing and durability decisions at exabyte scale
- https://www.usenix.org/system/files/conference/fast16/fast16-papers-beaver.pdf — Facebook's Haystack paper on blob storage, relevant to small-block metadata pressure
- https://engineering.fb.com/2014/01/22/core-infra/scaling-the-facebook-data-warehouse-to-300-pb/ — context on deduplication economics and storage tiering decisions
