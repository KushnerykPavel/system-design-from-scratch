# Live Streaming Platform (Twitch-style)

> Latency and quality are always in tension in live streaming. The real system design problem is ingest → transcode → edge distribution, not just "push to CDN."

**Type:** Learn  
**Company focus:** Cloudflare  
**Learning goal:** Design a live streaming platform covering RTMP ingest, adaptive bitrate transcoding, HLS/DASH edge distribution, viewer count at scale, and live chat — with explicit latency vs quality trade-offs for each layer.  
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/03-cdn-layering`, `07-queues-streams-and-workflows/01-queues-vs-streams`, `24-netflix-streaming-platform-design/02-video-streaming-pipeline`  
**Estimated time:** ~90 min  
**Primary artifact:** capacity sheet

## The Problem

Design a live streaming platform where streamers broadcast real-time video to potentially millions of concurrent viewers. Unlike on-demand video (Netflix), live streaming has fundamentally different constraints: content is produced and consumed simultaneously, viewers expect sub-30-second latency from streamer to viewer, and a stream must stay live even when one ingest node fails.

Mid-level candidates design a simple RTMP server that forwards to a CDN. Senior candidates recognize the four distinct layers: ingest (reliable receipt of the streamer's video from one source), transcoding (converting a single stream into multiple bitrate/resolution variants in near-real-time), distribution (delivering segments to millions of viewers at edge), and metadata (viewer count, chat, stream state). Each layer has different availability, latency, and scaling properties.

The chat-at-scale problem is also commonly underestimated. At 500K concurrent viewers, a single chat message sent by the streamer appears in every viewer's chat window within 1-2 seconds. This is a fan-out problem similar to social news feeds, but with much stricter latency requirements and much simpler content — and the scale can be sudden and unpredictable.

## Clarify

- What is the acceptable end-to-end latency from streamer to viewer? Ultra-low latency (sub-3s) or broadcast-style (15-30s HLS)? This fundamentally changes the protocol stack.
- Is viewer-to-viewer chat required, or only streamer-to-viewers broadcast? Viewer-to-viewer chat at 500K simultaneous users is much harder.
- What is the maximum expected concurrent viewer count per stream, and what is the total concurrent viewer count across all streams?

If the interviewer does not specify, assume HLS-based distribution with 10-30s latency, full viewer-to-viewer chat, and 1M concurrent viewers on peak streams with 50M global concurrent viewers across all streams.

## Requirements

### Functional

- Accept live video from streamers via RTMP or WebRTC.
- Transcode to multiple bitrate/resolution variants (1080p/4.5Mbps, 720p/2Mbps, 480p/1Mbps, 360p/400Kbps).
- Distribute HLS segments to viewers through a CDN with adaptive bitrate switching.
- Count concurrent viewers per stream and total platform viewers in near-real-time.
- Provide live chat with per-stream fan-out to all connected viewers.
- Detect and handle stream health issues (dropped frames, encoder failure, reconnect).

### Non-functional

- Ingest reliability: a streamer's connection to an ingest node should survive network blips with automatic reconnection.
- Transcode latency: add less than 5 seconds of latency in the transcoding pipeline.
- Viewer start time: first video segment loads within 3 seconds of joining a live stream.
- Chat delivery: messages appear within 2 seconds for 99% of recipients.
- Scale: handle 50M concurrent viewers platform-wide with per-stream peaks of 1M.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active streams | 500K concurrent live streams | sizes ingest fleet and transcode cluster |
| Ingest bandwidth | 500K streams × 6 Mbps avg = 3 Tbps ingest | ingest network capacity and cost |
| Transcode load | 500K streams × 4 renditions × 1 CPU-s/s = 2M CPU cores | largest infrastructure cost; GPU transcoding essential |
| Viewer egress | 50M viewers × 2 Mbps avg = 100 Tbps egress | must be fully served by edge CDN; no origin can serve this |
| Chat fan-out | peak stream: 1M viewers, 100 msgs/s = 100M deliveries/s | must use pub/sub fan-out per stream, not unicast |

## Architecture

```text
streamer encoder
  -> RTMP push to nearest ingest POP (anycast DNS)
  -> ingest service: validates stream key, accepts media
  -> segments video into 2-second GOP-aligned chunks
  -> forwards raw chunks to transcode workers

transcode worker cluster
  -> receives raw 2s chunks
  -> produces 4 renditions per chunk (GPU-accelerated)
  -> uploads each rendition segment to origin object store
  -> updates HLS manifest (m3u8) per stream

CDN distribution
  <- viewers pull HLS manifest from CDN edge
  <- viewers pull segments from CDN edge
  <- edge caches segments by stream/rendition/sequence
  <- origin pull-through for cache misses
  <- CDN edge nodes refresh manifests every 2s

viewer count
  -> viewer join/leave events -> Kafka
  -> stream processor (windowed count per stream, 10s window)
  -> viewer_counts Redis hash (stream_id -> count)
  -> streamer dashboard polls or subscribes via WebSocket

live chat
  -> viewer sends message -> chat API -> Kafka topic per stream
  -> chat fan-out worker reads Kafka, pushes to WebSocket gateway
  -> WebSocket gateway maintains persistent connections per viewer
  -> gateway delivers messages to viewer's browser/app
```

## Data Model & APIs

Core entities:

```text
Stream  { stream_id, user_id, stream_key, status, title, category,
          ingest_server, started_at, viewer_count, peak_viewer_count }
Segment { stream_id, rendition, sequence_num, duration_ms, storage_key, created_at }
Manifest{ stream_id, rendition, m3u8_content, updated_at }
ChatMessage { message_id, stream_id, user_id, username, content, sent_at }
```

Key APIs:

- `POST /v1/streams` — create stream, returns stream_key and RTMP ingest URL
- `GET /v1/streams/{stream_id}/hls/master.m3u8` — HLS master manifest (served via CDN)
- `GET /v1/streams/{stream_id}/hls/{rendition}/live.m3u8` — per-rendition manifest
- `GET /v1/streams/{stream_id}/viewer-count` — current viewer count
- WebSocket `/v1/chat/{stream_id}` — bidirectional; publish and subscribe to chat messages

HLS manifests include only the last N segments (sliding window); segment URLs point to CDN paths. The sequence number in each segment name enables efficient CDN caching without cache busting.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Streamer's connection to ingest node drops | ingest service detects TCP FIN or timeout within 2s | hold the HLS manifest with last known segments; streamer reconnect resumes from next sequence number |
| Transcode worker crashes mid-stream | segment missing from manifest, manifest staleness alert | redundant transcode workers; lost segment causes a brief playback stall but stream continues |
| CDN edge node is overloaded during peak | segment fetch error rate rising, CDN health alert | CDN load balancer redirects viewers to alternative edge POP; request surge handled by horizontal scale |
| Chat fan-out lag exceeds SLO | chat message delivery latency p99 alert | shed lower-priority metadata (emotes, badges) during load; increase WebSocket gateway capacity |

## Observability

- metric: viewer playback start latency p95 — first segment load time from stream join
- metric: transcode pipeline latency — time from raw chunk arrival to all rendition segments available
- metric: CDN cache hit rate per stream rendition — low rate means origin is getting hammered
- metric: chat fan-out latency p99 per stream — key SLO for popular streams
- log: every ingest connect/disconnect event with stream_id, ingest_server, duration, and disconnect reason
- SLO: 99.5% of viewers receive their first video segment within 3 seconds; chat messages delivered within 2 seconds for 99% of recipients

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| HLS segmentation (2s segments, 10-30s live edge) | CDN-cacheable; works with all clients; scales to 100 Tbps via edge | 10-30s latency is unacceptable for interactive use cases (watch parties, esports) | WebRTC relay — sub-1s latency but cannot scale to 1M viewers per stream without complex selective forwarding unit (SFU) mesh |
| Kafka fan-out for chat delivery via WebSocket gateways | decouples message production from delivery; survives gateway restarts | adds 100-500ms latency in the fan-out path | direct DB fan-out (poll DB per viewer) — does not scale past 10K concurrent chatters per stream |
| GPU-accelerated transcoding | cost-effective at scale; real-time for all 500K streams | GPU transcoding infrastructure is complex to operate; spot instance availability risk | CPU-only transcoding — too slow for real-time at scale; would require 10x the CPU core count |

## Interview It

**Google framing:** "Design a live streaming product that can scale to the Super Bowl — 20M concurrent viewers on one stream." Expect pushback on ingest reliability, transcode redundancy, and CDN burst handling.

**Cloudflare framing:** "How would you use edge infrastructure to reduce live streaming latency from 30s to 5s without sacrificing scalability?" Expect questions on LL-HLS, edge manifest generation, and the consistency model for segment availability.

**Follow-ups:**
1. How would you implement stream clips — a viewer captures and shares the last 30 seconds of a stream?
2. How would you support 100K concurrent streamers, each with a 4K 60fps source feed?
3. If a streamer's audio and video go out of sync by 2 seconds, how do you detect and correct it?
4. How would you implement channel subscriptions that notify 2M subscribers when a streamer goes live?
5. What changes if the platform must comply with COPPA and restrict chat for viewers under 13?

## Ship It

- `outputs/capacity-sheet-live-streaming.md`

## Exercises

1. **Easy** — Calculate the total CDN egress bandwidth for 50M concurrent viewers each watching at an average of 2 Mbps.
2. **Medium** — Design the stream reconnect protocol: a streamer's encoder crashes and reconnects 10 seconds later. How does the ingest service, transcode pipeline, and HLS manifest handle the gap?
3. **Hard** — Add Low-Latency HLS (LL-HLS) to reduce viewer latency from 10s to 2s. What changes in the segment generation, manifest format, CDN configuration, and player?

## Further Reading

- https://blog.twitch.tv/en/2022/04/26/ingesting-live-video-streams-at-global-scale/ — Twitch's ingest architecture at scale, including their SRT migration from RTMP
- https://developer.apple.com/documentation/http-live-streaming — Apple's HLS spec, the foundation for understanding segment-based live streaming delivery
- https://webrtc.org/getting-started/overview — WebRTC primer for understanding the ultra-low-latency alternative and why it doesn't scale to 1M viewers
