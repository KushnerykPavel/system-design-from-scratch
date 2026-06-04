---
lesson: 04-metadata-index-service
focus: balanced
---

## Clarify first
- Query shapes that matter to users versus operators
- Freshness expectations for listings and filters
- Whether enrichment updates can arrive asynchronously and out of order

## Must-size numbers
- Metadata write QPS
- Query QPS by shape
- Record size and index amplification
- Replay or backfill throughput budget

## Core design
- Canonical metadata owner
- Secondary projections tuned for actual query paths
- Replayable pipeline for rebuild and schema evolution

## Failure probes
- Lagging projection
- Stale overwrite from async enrichment
- Unsupported query combinations causing scans

## Trade-off summary
- Query power vs predictability
- Stronger write safety vs API complexity
- Faster evolution vs more backfill tooling
