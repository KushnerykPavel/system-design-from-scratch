# Media Pipeline — Photos, Videos & CDN

> Meta does not store photos. It stores memories — at a scale that makes every other media company look like a weekend project.

**Type:** Build  
**Company focus:** Meta  
**Learning goal:** Design Meta's end-to-end media pipeline for photos and videos. Cover Haystack object storage, async video transcoding, CDN architecture for Reels, storage lifecycle tiering, Everstore, and failure handling across the entire upload-to-serve path.  
**Prerequisites:** `03-news-feed-fanout`, `04-messenger-realtime`  
**Estimated time:** ~90 min  
**Primary artifact:** media pipeline design doc + CDN strategy spec  

## The Problem

Meta handles 100 million photo uploads per day and serves 4 billion photos daily across Facebook and Instagram. Reels generates 140 billion video plays per day. Every piece of media must be:

1. Stored durably with no data loss.
2. Transcoded into multiple formats and resolutions.
3. Served globally with low latency via CDN.
4. Tiered to cold storage as it ages, without breaking links.

Design the media pipeline that handles photos and videos from upload to serve at this scale.

## Clarify

- Are we designing for photos, videos, or both? (Both — they share infrastructure but diverge in transcoding.)
- What is the latency target for a photo to appear after upload? (<5s for photo, <60s for video first preview)
- How long must media be stored? (Indefinitely unless deleted by user or policy)
- Is deduplication required? (Yes — content-addressable storage prevents storing the same photo twice)
- What is the acceptable read latency for serving a photo? (p99 <100ms from CDN edge)

## Requirements

### Functional

- Accept photo and video uploads from clients worldwide.
- Store media durably with 11-nines durability.
- Transcode videos into multiple formats and resolutions.
- Serve media via globally distributed CDN.
- Tier cold media to cheaper storage automatically.
- Detect and deduplicate identical content.

### Non-functional

- Upload throughput: 100M photos/day (~1,160 photos/second sustained).
- Storage: 4B+ photos in active storage; petabytes of video.
- Serve latency: p99 <100ms for photos from CDN edge; <500ms for video first segment.
- Durability: 11 nines (replicated across 3+ data centers).
- Availability: 99.99% for serve path; uploads may retry on failure.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Photo uploads | 100M/day = 1,160/s | Upload service sizing |
| Average photo size | 3 MB | Storage growth rate |
| Photo storage growth | ~300 TB/day | Capacity planning |
| Video uploads | ~5M/day | Transcoding worker sizing |
| Transcoded variants per video | ~15 (5 resolutions × 3 codecs) | Storage multiplier |
| Reels plays | 140B/day = 1.6M/s | CDN bandwidth sizing |
| CDN bandwidth for video | ~500 Gbps sustained | PoP sizing |
| Photos served | 4B/day = 46,000/s | CDN cache hit ratio target |

## Architecture

```text
[Photo upload path]
  client
  -> upload service (multipart, resumable, checksum validation)
  -> deduplication service (content hash lookup)
  -> Haystack / Everstore (write needle, return photo ID)
  -> CDN fill (async; photo ID resolves to CDN-cacheable URL)
  -> return photo URL to client

[Video upload path]
  client
  -> upload service (resumable chunks to blob store)
  -> original stored in cold object store
  -> transcoding job published to Kafka
  -> transcoding workers (parallel per variant)
     -> MP4 H.264, VP9, AV1 × 360p/720p/1080p/4K
  -> transcoded segments stored in Everstore
  -> manifest (DASH/HLS) generated and stored
  -> CDN fill triggered for popular content
  -> video URL returned to client

[Serve path]
  client request -> CDN PoP (cache hit: <10ms)
                 -> CDN miss -> Meta origin cache
                             -> Haystack / Everstore
```

### Haystack: Storing Small Files Efficiently

Traditional POSIX filesystems maintain per-file metadata (inode, directory entry, timestamps) which requires one disk seek per read even before the file data is read. At 4B photos, the inode table alone becomes the bottleneck.

Haystack solves this by packing many photos (needles) into a single large flat file (volume):

```text
Volume = flat binary file, typically 100 GB
Needle = [magic | cookie | key | alt_key | flags | size | data | checksum | padding]

Lookup path:
  photo_id -> in-memory directory -> (volume_id, offset, size)
  one seek to offset, one read of size bytes
  zero filesystem metadata overhead
```

The directory (volume_id → physical machine, offset) fits entirely in RAM. Meta runs ~3 types of Haystack machines:

- **Logical machine**: stores needles, exports NFS-like interface.
- **Physical machine**: actual disk host with multiple logical volumes.
- **Cache machine**: in-memory LRU of recently accessed needles; absorbs read spikes on hot photos.

Key design properties:
- Write once, read many — needles are never modified, only soft-deleted via flag.
- Append-only writes eliminate random I/O during uploads.
- Compaction runs offline to reclaim space from deleted needles.

### Everstore: Content-Addressable Successor to Haystack

Everstore is Meta's newer object storage that replaces Haystack for new writes:

- **Content-addressed**: object key = SHA-256 of content. Identical photos stored once automatically.
- **Append-only**: writes are immutable once committed.
- **Replicated 3×**: across availability zones, synchronous write quorum of 2.
- **Tiered**: hot tier (SSD + RAM cache), warm tier (HDD), cold tier (erasure-coded, off-site).

```text
write(data):
  hash = SHA-256(data)
  if exists(hash): return existing_url  // deduplication
  store to hot tier, replicate to 2 more zones
  return content_url(hash)
```

### Video Transcoding Pipeline

Video transcoding is CPU-intensive and asynchronous. Meta runs it as a pipeline:

```text
upload complete
  -> publish TranscodeJob{video_id, s3_key, priority} to Kafka
  -> transcoding worker pool (auto-scaled):
     -> pull job from Kafka consumer group
     -> download original from cold store
     -> transcode to each variant (ffmpeg or custom encoder)
     -> upload variant to Everstore
     -> update video metadata (variants available, manifest URL)
  -> manifest generator: create MPEG-DASH + HLS manifests
  -> CDN proactive fill for predicted-popular videos
```

Transcoding variants for a typical Reels video:

| Codec | Resolutions | Notes |
|-------|-------------|-------|
| H.264 (MP4) | 360p, 720p, 1080p | Broadest device compatibility |
| VP9 (WebM) | 360p, 720p, 1080p | 30–50% smaller than H.264 |
| AV1 | 720p, 1080p, 4K | 50% smaller; higher encode cost; newest devices |

Transcoding is idempotent: if a worker crashes, the job is re-queued from Kafka and the output is overwritten.

### CDN Architecture

Meta operates its own CDN with Akamai as a fallback for regions with insufficient Meta infrastructure:

```text
client
  -> DNS anycast -> nearest Meta edge PoP
  -> cache hit: serve from edge SSD (< 10ms)
  -> cache miss:
       -> Meta regional hub (larger cache, L2)
       -> cache hit at hub: serve + fill edge
       -> cache miss at hub: fetch from origin (Haystack/Everstore)
       -> fill hub and edge
```

For Reels, Meta uses **Adaptive Bitrate (ABR)** streaming:

- **DASH (Dynamic Adaptive Streaming over HTTP)**: video split into 2–6 second segments; client selects quality per segment based on measured bandwidth.
- **HLS (HTTP Live Streaming)**: Apple device support.
- CDN serves manifests and segments as static cacheable assets.
- Popular Reels are proactively pushed to edge PoPs before the surge (prediction based on early engagement signals).

### Storage Lifecycle Tiering

Not all photos are accessed equally. The access distribution follows a power law: 90% of reads are on photos uploaded in the last 6 months.

```text
Age 0–30 days:  hot tier (SSD, high-IOPS Haystack/Everstore)
Age 31–180 days: warm tier (HDD-based Everstore)
Age 181 days+:  cold tier (erasure-coded, off-site, Glacier-equivalent)
```

Lifecycle policy runs nightly: scans for objects past their age threshold, moves them to the next tier, updates the directory to point to the new location. Reads on cold-tier objects require a re-warm step (~seconds latency).

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Transcoding job failure | Worker reports failure; Kafka offset not committed | Retry up to 3× with exponential backoff; dead-letter queue after 3 failures; alert oncall |
| CDN cache miss storm on viral video | Origin traffic spike > 10× baseline | Pre-fill edge caches via proactive CDN push when engagement prediction fires; origin auto-scales with circuit breaker |
| Storage node failure mid-upload | Upload service timeout or write quorum failure | Client retries with resumable upload token; partial write discarded; new write to healthy nodes |
| Transcoding DLQ backlog grows | DLQ depth alert | Oncall reviews DLQ; common causes: corrupt input, unsupported codec; mark video as processing-failed and notify user |
| Cold tier retrieval latency spike | p99 read latency > 5s on cold-tier objects | Cache frequently-re-accessed cold objects back to warm tier after first re-access |

## Observability

- metric: upload success rate (photos and videos) by region
- metric: transcoding job queue depth and consumer lag per priority bucket
- metric: CDN cache hit ratio at edge and hub tiers
- metric: CDN origin fetch rate (proxy for cache miss rate)
- metric: storage tier distribution (% of objects in hot/warm/cold)
- metric: photo serve latency at p50/p95/p99 from CDN edge
- log: transcoding failures with input video metadata and error code
- trace: upload request from client through to CDN fill confirmation
- alert: DLQ depth > 1000 jobs; CDN cache hit ratio < 85%

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Content-addressed storage (Everstore) | Automatic deduplication; simpler replication (same hash = same bytes) | Hash computation on write; slightly more complex lookup (hash → URL) | Path-based storage simpler to implement but duplicates identical photos |
| Eager transcoding (all variants at upload) | All resolutions immediately available | Higher storage cost; higher CPU cost per upload | Lazy transcoding saves cost but first viewer waits; unacceptable for Reels |
| Own CDN vs pure Akamai | Lower per-request cost at Meta's scale; tighter integration with origin | High CapEx for PoP build-out | Pure Akamai cheaper at small scale but Meta's volume makes own CDN 10× cheaper per request |
| Append-only Haystack volumes | No random I/O on write; simple consistency model | Wasted space from deleted needles until compaction | Mutable storage would require complex locking and metadata updates at 100M writes/day |
| Async transcoding via Kafka | Decouples upload latency from transcoding latency; survives worker crashes | Video not immediately available in all resolutions | Sync transcoding blocks upload response; unacceptable at scale |

## Interview It

**Meta framing:** "Design Instagram's photo and video storage system." Strong answers cover Haystack or Everstore internals (needles/volumes or content addressing), async transcoding pipeline with failure handling, CDN with ABR for video, and storage tiering. Weak answers describe generic S3 + CDN without addressing the metadata overhead problem, deduplication, or transcoding variants.

**Follow-ups:**

1. Why does Haystack avoid filesystem metadata overhead, and what data structure does it use to replace it?
2. A Reels video goes viral 10 minutes after upload. How does the CDN handle the sudden traffic surge without the origin being overwhelmed?
3. A user uploads a photo that is identical to one already stored. How does Everstore detect this, and what does the write path look like?
4. Why does Meta transcode videos to AV1 when it costs significantly more compute than H.264?
5. How would you design the storage tier migration process so that a photo moving from warm to cold tier never returns a 404 during the migration?

## Ship It

- `outputs/design-doc-media-pipeline.md`
- `outputs/cdn-strategy-spec.md`
- `outputs/interview-card-media-pipeline.md`

## Exercises

1. **Easy** — List the Haystack needle fields. Which field enables soft-delete without removing data from disk?  
2. **Medium** — Design the transcoding job schema on Kafka. What fields are needed to make jobs idempotent and retriable?  
3. **Hard** — Design the CDN proactive fill system that pushes Reels videos to edge PoPs before they go viral. What signals would you use to predict which videos to push?  

## Further Reading

- [Haystack: Finding a Needle in Facebook's Haystack (OSDI 2010)](https://www.usenix.org/legacy/event/osdi10/tech/full_papers/Beaver.pdf)  
- [Meta Engineering Blog — video infrastructure](https://engineering.fb.com/2021/01/08/video-engineering/video-processing-pipeline/)  
- [MPEG-DASH specification and ABR overview](https://dashif.org/docs/DASH-IF-IOP-v4.0.pdf)  
