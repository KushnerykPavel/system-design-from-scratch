# Log Aggregation Pipeline

> Logs are the easiest telemetry to emit and the easiest telemetry to drown in if you do not design the pipeline around cost, redaction, and backpressure.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a log pipeline that cleanly separates ingestion, enrichment, indexing, archival, and query isolation while handling PII, burstiness, and replay.
**Prerequisites:** `05-storage-indexing-and-access-patterns/07-retention-and-deletion`, `07-queues-streams-and-workflows/04-dlq-and-replay`, `11-observability-slos-and-debugging/03-logs-and-traces`
**Estimated time:** ~75 min
**Primary artifact:** pipeline-policy validator + failure checklist

## The Problem

Design a centralized log aggregation system for many services. It must collect structured and unstructured events, support search on recent logs, archive older data cheaply, and avoid turning incident spikes into a second outage.

The senior-level signal is whether you distinguish logs from metrics: weaker schemas, heavier payloads, stronger privacy risk, and a more painful balance between indexing everything versus storing cheaply and querying later.

## Clarify

- Are logs mainly for debugging, compliance, security, or customer analytics?
- How long must searchable logs stay hot before moving to archive?
- Do we require replay into downstream processors after parser bugs or policy fixes?
- What privacy or redaction constraints apply before persistence?

If the interviewer leaves it broad, assume operational and security use cases, hot search for a few days, cold archive for months, and mandatory redaction before indexable storage.

## Requirements

### Functional

- Collect logs from many services and regions.
- Parse and enrich logs with service, tenant, and trace metadata.
- Support recent full-text or fielded search.
- Archive older logs cheaply with replay support.
- Enforce redaction and retention policies.

### Non-functional

- Survive incident bursts without losing all operational visibility.
- Prevent one noisy producer from overwhelming shared indexing.
- Keep hot search responsive while archiving large volumes.
- Make privacy-sensitive data handling auditable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Events ingested | 15M events/s peak | shapes queueing, batching, and regional fan-in |
| Average event size | 1.5 KB | storage and egress costs grow quickly |
| Hot searchable retention | 3 days | expensive indexing tier must stay bounded |
| Cold archive retention | 180 days | cost and replay workflow dominate operations |
| Peak factor | 8x during incidents | backpressure strategy matters when the platform is most needed |

## Architecture

```text
agents
  -> regional collectors
  -> parse + enrich + redact
  -> durable event bus
  -> hot indexers
  -> searchable hot store
  -> archive writers
  -> replay / reprocess jobs
```

Design notes:

1. Redact as early as possible so unsafe payloads do not leak into every downstream system.
2. Use a durable bus between collection and indexing so hot search can fall behind briefly without dropping all data.
3. Separate searchable hot storage from cold archive because query shape and cost profile differ sharply.
4. Decide which logs are droppable during backpressure and which must be retained for security or compliance.

## Data Model & APIs

Representative fields:

```text
timestamp
service
region
severity
trace_id
tenant_id
event_body
redaction_state
retention_class
```

Useful interfaces:

- `IngestLogs(batch[])`
- `SearchLogs(query, start, end, cursor)`
- `ReplayLogs(source_range, parser_version)`
- `SetRetentionPolicy(service, hot_days, archive_days)`

The data model should make retention class and redaction state explicit rather than implied.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| parser deployment breaks structured extraction | parse-failure rate and schema drift alerts | versioned parsers and replay from durable bus |
| hot indexing tier saturates during incident storm | indexing lag and query latency blowup | backpressure, tiered sampling for low-value logs, archive-first buffering |
| PII reaches hot search index | redaction-miss audits and canary detectors | early redaction, quarantine, and reprocess-delete workflow |
| one tenant floods shared storage | tenant ingest skew and quota breaches | per-tenant quotas, routing isolation, and deferred indexing |

## Observability

- metric: ingest accepted, delayed, sampled, and dropped by class
- metric: parser failure rate and top schema drift sources
- metric: hot index lag, search latency, and archive write backlog
- metric: replay job age and bytes pending by parser version
- log: redaction actions, quarantine decisions, and tenant quota enforcement
- trace: collector to parse to bus to index pipeline for sampled batches
- SLO: hot incident logs from protected classes are searchable within the documented freshness target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| hot search plus cold archive | balances usability and cost | dual storage and replay complexity | indexing everything forever |
| early redaction | lower privacy blast radius | some transforms become irreversible | storing raw logs everywhere then cleaning later |
| differentiated backpressure | preserves critical logs during storms | policy complexity and fairness debates | equal treatment for all log classes |

## Interview It

**Google framing:** "Design internal centralized logging for services and SREs." Expect questions on retention tiers, search cost, and what happens during incident spikes.

**Cloudflare framing:** "Design a globally distributed log pipeline for edge and platform events." Expect pressure on regional ingestion, privacy handling, and replay after parsing mistakes.

**Follow-ups:**
1. Which logs are safe to sample and which are not?
2. How do you replay old logs after changing parsing rules?
3. What changes if security teams need near-real-time detections from the same stream?
4. How do you prove PII redaction happened before indexing?
5. How would you isolate a noisy tenant without hiding their logs completely?

## Ship It

- `outputs/failure-checklist-log-pipeline.md`

## Exercises

1. **Easy** — Split three example log classes into hot search, archive-only, and droppable-under-pressure.
2. **Medium** — Design a replay flow after a parser bug misclassified an important field.
3. **Hard** — Redesign the pipeline for strict regional data-boundary requirements.

## Further Reading

- [OpenTelemetry logs data model](https://opentelemetry.io/docs/specs/otel/logs/data-model/) — useful for schema and correlation thinking
- [Site Reliability Engineering: Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — relevant to telemetry backpressure during incidents
