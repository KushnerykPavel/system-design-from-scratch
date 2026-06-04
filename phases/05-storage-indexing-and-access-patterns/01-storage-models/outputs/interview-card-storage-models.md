---
lesson: 01-storage-models
focus: balanced
---

## Clarify first
- What are the top reads and writes?
- Which invariants must hold synchronously?
- Which queries are flexible versus tightly bounded?

## Must-size numbers
- Peak read/write QPS
- Data growth and retention horizon
- Share of traffic that is simple key lookup versus secondary query

## Core design
- Name one canonical system of record
- Choose relational, KV, or document based on workload fit
- Push weakly matched queries into derived systems later

## Failure probes
- What breaks if the chosen store must support ad hoc filtering?
- What if cross-entity invariants appear later?
- What if the derived index lags behind the source of truth?

## Trade-off summary
- Relational: strong correctness, heavier operational scaling for joins and hotspots
- KV: low-latency lookup, weak fit for rich queries and invariants
- Document: flexible object shape, harder index and consistency discipline
