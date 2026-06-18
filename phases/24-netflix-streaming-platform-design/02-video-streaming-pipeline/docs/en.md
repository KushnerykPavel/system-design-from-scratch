# Video Streaming Pipeline — Encoding, Chunking, ABR

> A video file is not a streamable asset. It becomes one only after a pipeline of encoding, packaging, and manifest generation runs on it.

**Type:** Build  
**Company focus:** Netflix  
**Learning goal:** Design the video ingestion → transcoding → packaging → delivery pipeline. Understand multi-bitrate encoding ladders, HLS/DASH segmentation, manifest generation, chunk storage, and adaptive bitrate algorithms at the client.  
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/`, `03-open-connect-cdn`  
**Estimated time:** ~90 min  
**Primary artifact:** pipeline design doc + encoding ladder spec  

## The Problem

A studio delivers a 4K RAW or camera original file to Netflix. Before a single subscriber can press play, that file must be encoded into dozens of quality variants at multiple resolutions, segmented into small chunks for adaptive streaming, packaged with encryption, and stored durably in object storage. Any subscriber anywhere must then receive the right bitrate for their current network conditions, with seamless switching between bitrates as conditions change.

Design this pipeline end to end.

## Clarify

- What is the target device matrix? (phones, smart TVs, browsers, set-top boxes — each may support different codecs and container formats)
- What are the latency requirements between content delivery and subscriber availability? (hours? minutes?)
- What encoding formats are required? (H.264, H.265/HEVC, AV1 — different complexity/quality trade-offs)
- Is the content live or video-on-demand? (live streaming has very different latency and segment-size constraints)
- What is the DRM requirement? (Widevine, FairPlay, PlayReady)

## Requirements

### Functional

- Ingest source video files from studio delivery (S3-compatible object store or direct upload).
- Transcode into multiple resolutions (240p through 4K) and multiple codecs.
- Segment each encoded variant into fixed-duration chunks for HTTP streaming.
- Generate HLS (`.m3u8`) and DASH (`.mpd`) manifests describing all available variants.
- Store chunks and manifests durably in object storage.
- Serve chunks and manifests from the CDN.
- Select and switch bitrates at the client based on current network throughput.

### Non-functional

- Processing latency from ingest to availability: under 4 hours for catalog titles, under 30 seconds for live.
- Durability of encoded assets: 11 nines (replicated across regions).
- Parallel encoding: thousands of titles simultaneously without head-of-line blocking.
- Chunk delivery latency: p99 under 200ms from CDN edge.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Catalog size | ~36,000 titles | total storage footprint for encoded variants |
| Variants per title | ~120 (codecs × resolutions × bitrates) | encoding compute demand |
| Chunk size | 2–10 seconds of video | affects seek latency and buffer efficiency |
| Chunks per title | ~6,000 per variant at 2-hour runtime, 4-sec chunks | storage and manifest complexity |
| Peak delivery | hundreds of Tbps | CDN and origin sizing |
| Storage per title | ~100 GB across all variants | object store capacity planning |

## Architecture

```text
studio delivery
  -> ingest service (validates, checksums, stores source)
  -> job scheduler (splits work into per-segment jobs)
     -> encoding workers (parallel, codec-specific)
        -> packager (chunks + DRM encryption + manifest generation)
           -> object store (S3-like, multi-region replication)
              -> CDN (Open Connect + tiered cache)
                 -> client player (ABR algorithm)
```

### Encoding Ladder

An encoding ladder defines the set of (resolution, bitrate, codec) combinations produced for every title. Netflix's published ladder includes:

| Profile | Resolution | Bitrate (H.264) |
|---------|------------|-----------------|
| 0 | 320x240 | 235 kbps |
| 1 | 384x288 | 375 kbps |
| 2 | 512x384 | 560 kbps |
| 3 | 512x384 | 750 kbps |
| 4 | 640x480 | 1050 kbps |
| 5 | 720x480 | 1750 kbps |
| 6 | 1280x720 | 2350 kbps |
| 7 | 1280x720 | 3000 kbps |
| 8 | 1920x1080 | 4300 kbps |
| 9 | 1920x1080 | 5800 kbps |

Netflix also generates per-title encoding where simpler content (animation, talking heads) uses lower bitrates than complex action scenes for the same quality — reducing storage and bandwidth without perceptible quality loss.

### HLS Segmentation

HLS works by splitting each variant stream into small `.ts` or `.fmp4` chunks referenced by a `.m3u8` playlist:

```text
#EXTM3U
#EXT-X-VERSION:6
#EXT-X-TARGETDURATION:4
#EXT-X-SEGMENT-SEQUENCE:0
#EXTINF:4.0,
seg-0000.ts
#EXTINF:4.0,
seg-0001.ts
...
```

A master manifest then references all variant playlists:

```text
#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=235000,RESOLUTION=320x240
/v1/{title_id}/h264/low/manifest.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=5800000,RESOLUTION=1920x1080
/v1/{title_id}/h264/high/manifest.m3u8
```

### Segment Size Trade-offs

| Segment duration | Benefit | Cost |
|-----------------|---------|------|
| 2 seconds | Fast bitrate switch, low buffer for live | More HTTP requests, higher manifest parsing overhead |
| 4 seconds | Good balance for VOD | Moderate seek latency |
| 10 seconds | High throughput, fewer HTTP requests | Slow ABR response, poor seek experience |

Netflix uses ~4-second segments for VOD and ~2-second segments for live content.

### Adaptive Bitrate Algorithm (Client Side)

The client player monitors available bandwidth and selects the appropriate variant:

```text
bandwidth estimate = EWMA(recent_chunk_download_time)
if bandwidth_estimate > variant_bitrate * 1.5:
    switch up
elif bandwidth_estimate < current_variant_bitrate * 0.8:
    switch down
```

Key design decisions:
- **Buffer-based vs bandwidth-based**: Pure bandwidth estimation reacts to short spikes; buffer occupancy gives smoother behavior.
- **Conservative switching up, aggressive switching down**: A bad quality drop is more disruptive than a delayed quality upgrade.
- **Pre-buffer threshold**: Do not start playback until at least N seconds of the lowest quality are buffered.

## Data Model & APIs

Encoding job:
```text
job_id, title_id, source_s3_key, codec, resolution, bitrate, status, created_at, completed_at
```

Manifest record:
```text
title_id, codec, format (hls|dash), s3_manifest_key, variant_count, created_at
```

Segment record:
```text
title_id, variant_id, segment_index, duration_ms, s3_chunk_key, encrypted (bool)
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Encoding worker crashes mid-job | Job heartbeat expires | Reschedule job to another worker; idempotent segment output |
| Source file is corrupted | Checksum mismatch at ingest | Reject at ingest; alert studio delivery; never enqueue corrupt source |
| Manifest points to missing segment | Client gets 404 on segment fetch | Segment existence validation before manifest is published |
| DRM key service is unavailable | Packaging fails | Queue packager, retry with backoff; do not publish unencrypted segments |
| CDN serves stale manifest after re-encode | Client gets wrong segment list | Use versioned manifest URLs; CDN cache bust on publish |

## Observability

- metric: encoding job queue depth by codec and priority
- metric: encoding worker utilization and error rate by job type
- metric: time from ingest to first-segment availability
- metric: client-side ABR switch rate (up vs down) per session
- metric: client buffer stall events per 1000 plays
- log: per-job encoding start, completion, and failure with codec and resolution
- trace: ingest → encode → package → manifest-publish latency

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Per-title encoding ladder | Lower bitrate for same quality on simple content | More complex encoding pipeline; longer processing per title | Fixed ladder simpler but wastes bandwidth on simple scenes |
| 4-second segments for VOD | Good ABR responsiveness, reasonable request rate | Seek requires fetching partial segment at boundaries | Longer segments reduce CDN request count but hurt seek and ABR |
| Parallel encoding workers | Thousands of titles encoded simultaneously | Worker scheduling and result coordination complexity | Sequential encoding cannot meet catalog-size requirements |
| DRM encryption at packaging time | Decouples encoding from content protection decisions | Packaging step adds latency; key service is a critical dependency | Encrypting at encoding time ties codec and DRM concerns together |

## Interview It

**Netflix framing:** "Walk me through how a video file gets from a studio to a subscriber's screen." Strong answers cover the full pipeline: ingest, transcoding, packaging, CDN, and client ABR. Weak answers jump straight to "store in S3 and serve from CDN."

**Follow-ups:**
1. How does the client decide to switch from 1080p to 720p mid-stream?
2. What happens if a single encoding worker is 10x slower than average?
3. How would you handle a codec transition (e.g., adding AV1 support) without re-encoding the entire catalog?
4. How do you prevent a corrupt encode from reaching subscribers?
5. What changes for live streaming vs VOD?

## Ship It

- `outputs/design-doc-video-streaming-pipeline.md`
- `outputs/encoding-ladder-spec.md`
- `outputs/interview-card-video-streaming-pipeline.md`

## Exercises

1. **Easy** — Draw the manifest hierarchy for a title with 3 codecs × 5 resolutions. How many files are generated?  
2. **Medium** — Design the job scheduler that parallelizes encoding of one 2-hour title across 100 workers.  
3. **Hard** — Extend the design to support live streaming with sub-5-second end-to-end latency (low-latency HLS or CMAF).  

## Further Reading

- [Netflix per-title encoding](https://netflixtechblog.com/per-title-encode-optimization-7e99442b62a2)  
- [HLS specification (Apple)](https://developer.apple.com/streaming/)  
- [DASH Industry Forum](https://dashif.org/)  
