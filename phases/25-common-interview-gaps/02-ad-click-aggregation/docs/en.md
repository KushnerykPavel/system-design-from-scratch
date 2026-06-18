# Ad Click Aggregation Pipeline

> Counting clicks sounds trivial. Counting them correctly at 10 M/s, exactly once, in the presence of fraud, late arrivals, and partial failures is the actual interview.

**Type:** Learn  
**Company focus:** Google  
**Learning goal:** Design a streaming aggregation pipeline that handles at-least-once delivery, idempotent deduplication, tumbling and sliding aggregation windows, fraud filtering, and accurate late-data correction without reprocessing the full history.  
**Prerequisites:** `07-queues-streams-and-workflows/02-delivery-semantics`, `07-queues-streams-and-workflows/03-consumer-groups`, `17-search-crawl-and-monitoring-systems/03-metrics-platform`  
**Estimated time:** ~75 min  
**Primary artifact:** capacity sheet

## The Problem

Design the backend pipeline that receives ad click events from browsers and mobile apps, aggregates them by campaign and time window, and makes the results available to advertisers for billing and performance reporting within minutes.

The naive answer describes a Kafka topic and a database counter increment. That fails to explain several non-obvious problems that Google engineers deal with daily. First, click events arrive at-least-once because clients retry on network failure and CDN edge nodes buffer before forwarding. Without deduplication, you double-count and overbill advertisers. Second, late events arrive after their window has already been reported: a user clicked at 11:58 but the event arrived at 12:05 after the minute-0 window closed. A correct system must handle late data without re-running the entire pipeline. Third, fraudulent clicks from bots and click farms artificially inflate counts and must be filtered before billing aggregation, but fraud scoring is slow and cannot block ingestion.

Senior candidates segment the pipeline into ingestion, deduplication, fraud scoring, windowed aggregation, and reporting layers and articulate the guarantees each layer provides and violates.

## Clarify

- What is the aggregation granularity required: per-minute buckets, per-hour, or both? And what is the maximum acceptable reporting delay after a window closes?
- Is billing based on aggregated counts or on individual click records? This determines how strict the deduplication guarantee needs to be.
- How long after event time can late events still be accepted for correction, versus discarded as too old?

If the interviewer does not specify, assume per-minute and per-hour windows, billing based on aggregated counts, reporting delay under 5 minutes, and late-arrival acceptance window of 1 hour.

## Requirements

### Functional

- Ingest click events from millions of ad-serving endpoints at high throughput.
- Deduplicate events so each unique click is counted at most once.
- Apply fraud scoring and exclude fraudulent clicks from billing aggregates.
- Compute aggregated click counts by (campaign_id, ad_id, time_window).
- Serve aggregated results to the reporting API within 5 minutes of window close.
- Correct previously reported windows when late events arrive within the acceptance window.

### Non-functional

- Ingest throughput: 5–10 M click events per second at peak (major campaigns, holiday season).
- End-to-end pipeline latency: 95th percentile under 3 minutes from click to aggregate availability.
- Deduplication accuracy: less than 0.01% duplicate clicks reach billing aggregation.
- Late-event correction must not require reprocessing the full historical dataset.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Ingest rate | 5 M clicks/s peak, ~200 B clicks/day | sizes Kafka partition count and ingestion fleet |
| Dedup bloom filter | 200 B events × 20 bits = 500 GB/day bloom storage | determines whether in-memory bloom is feasible or must be Redis-sharded |
| Aggregation state | 10 M active (campaign, window) pairs × 100 bytes = 1 GB hot state | fits in-memory per Flink/Spark task slot with moderate parallelism |
| Storage for raw events | 200 B × 200 bytes = 40 TB/day compressed | drives retention vs replay cost tradeoff |
| Fraud scoring latency | 50–200 ms per batch | must be async to avoid blocking the ingestion critical path |

## Architecture

```text
client (browser / app)
  -> click event {click_id, ad_id, campaign_id, user_id, timestamp, geo, device}
  -> CDN edge (buffers, deduplicates in-flight via click_id Bloom filter)
  -> ingestion API (HTTP/2, batched write)

ingestion API
  -> Kafka topic: raw-clicks (partitioned by campaign_id)

stream processor (Flink or Kafka Streams)
  -> dedup stage: check click_id against Redis Bloom filter (TTL = 2 hours)
  -> fraud score stage (async): forward to fraud service, mark or drop
  -> windowed aggregation: tumbling 1-min windows keyed by (campaign_id, ad_id)
  -> emit window results to: aggregation store (ClickHouse / BigQuery)
  -> emit late corrections to: correction topic

late-event handler
  -> events arriving within acceptance window trigger delta patch to aggregation store
  -> events outside acceptance window are logged but excluded from billing

reporting API
  -> serves ad_click_counts(campaign_id, start_time, end_time)
  -> reads from aggregation store
  -> supports correction acknowledgement for advertiser billing reconciliation
```

## Data Model & APIs

Core event schema (Avro):

```text
ClickEvent {
  click_id: string       // globally unique, client-generated UUID
  ad_id: string
  campaign_id: string
  user_agent_hash: bytes
  ip_hash: bytes
  geo_country: string
  event_time: timestamp  // client clock; used for windowing
  ingest_time: timestamp // server clock; used for late detection
}

AggregateRecord {
  campaign_id: string
  ad_id: string
  window_start: timestamp
  window_end: timestamp
  click_count: int64
  unique_users: int64
  fraud_filtered_count: int64
  correction_version: int
}
```

Key APIs:

- `POST /v1/clicks` — ingest endpoint; accepts batch of up to 500 events; returns `{accepted: N, duplicate_detected: M}`
- `GET /v1/reports/clicks?campaign_id=X&start=T1&end=T2&granularity=1m` — aggregated click counts per window
- `GET /v1/reports/corrections?campaign_id=X&since=T` — returns correction log for advertiser reconciliation

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Kafka consumer lag grows during ingestion spike | consumer lag metric crossing alert threshold | auto-scale stream processor replicas; apply backpressure to ingestion API |
| Bloom filter false positive deduplicates a real click | dedup rejection rate relative to ingest rate anomaly | tune bloom false-positive rate to < 0.001%; log rejected click_ids for audit |
| Fraud service is slow or unavailable | fraud scoring latency p99 alert; async queue depth grows | mark events as fraud-pending, emit without fraud label, backfill asynchronously when service recovers |
| Stream processor crashes mid-window | checkpoint failure alert; window output delayed | Flink checkpoints allow restart from last committed offset; window state is recovered from durable state backend |

## Observability

- metric: click ingest rate (events/s), dedup rejection rate (%), and fraud filter rate (%) — measured per campaign and globally
- metric: aggregation pipeline lag (time between window close and aggregate available in reporting store)
- metric: late-event rate by arrival delay bucket (0–60s, 60s–10m, 10m–1h, >1h)
- log: every dedup rejection with click_id, campaign_id, and detection stage (edge vs stream processor)
- trace: end-to-end latency from ingest API receipt to aggregate store write for a sampled click batch
- SLO: 99% of closed windows have final aggregate available within 5 minutes; late corrections applied within 10 minutes of detection

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Bloom filter deduplication (probabilistic) | O(1) lookup, low memory, works at 5M events/s | false positives reject a tiny fraction of valid clicks | exact dedup in a key-value store — too slow at 5M/s write rate and too expensive for 200B daily events |
| Async fraud scoring, out of band from billing aggregation | fraud service outage does not block billing | some fraudulent clicks may reach billing temporarily | synchronous fraud gate on ingestion — adds 200ms+ to every click, unacceptable at scale |
| Late-event correction with delta patches | avoids reprocessing full history; correction is incremental | correction logic is complex; advertisers must reconcile delta | ignore late events — legally and contractually unacceptable for billing accuracy |

## Interview It

**Google framing:** "Design the click aggregation system for Google Ads. It processes billions of clicks per day and must be accurate enough for billing." Expect pushback on deduplication guarantees, fraud filtering coupling, and late event handling.

**Cloudflare framing:** "Design an edge-integrated click counting system where the first dedup pass happens at the POP before data reaches the core pipeline." Expect questions on edge state management, POP-to-core consistency, and what happens when a POP is partitioned.

**Follow-ups:**
1. How would you handle a 10x traffic spike during a major product launch without pre-provisioning capacity?
2. If a stream processor bug double-counted all clicks for a 2-hour window, how would you detect and correct it?
3. How would you add a real-time fraud detection model that needs 100ms to score each click?
4. Advertisers demand click-level audit logs for legal disputes. How does that change the pipeline?
5. How would you evolve the schema if you need to add a new event field required for a new bidding model?

## Ship It

- `outputs/capacity-sheet-ad-click-aggregation.md`

## Exercises

1. **Easy** — Calculate the Bloom filter memory required to deduplicate 200 billion daily click_ids at a 0.001% false positive rate.
2. **Medium** — Design the correction pipeline: when a late click arrives 45 minutes after its window closed, trace the full data path from ingestion to updated aggregate in the reporting store.
3. **Hard** — A bot farm generates 50M fake clicks with valid-looking click_ids before the fraud model catches them. Design the retroactive correction process that ensures billing is accurate and auditable.

## Further Reading

- https://engineering.fb.com/2018/11/29/web/facebook-redesigns-pixel/ — Facebook's approach to deduplication and event reliability at advertising scale
- https://flink.apache.org/2015/12/04/introducing-stream-windows-in-apache-flink/ — Flink windowing internals, essential for understanding tumbling vs sliding vs session windows
- https://dataintensive.net — Chapter 11 of "Designing Data-Intensive Applications" covers stream processing, windowing, and exactly-once semantics in depth
