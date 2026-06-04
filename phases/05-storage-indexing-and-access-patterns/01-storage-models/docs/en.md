# Relational vs KV vs Document Stores

> Pick the store that matches the access pattern and invariants, not the one with the best marketing page.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Explain when relational, key-value, or document storage is the right primary system of record and name the operational consequences of each choice.  
**Prerequisites:** `02-estimation-and-cost/02-storage-growth`, `04-apis-contracts-and-schema-evolution/03-pagination-and-filtering`  
**Estimated time:** ~75 min  
**Primary artifact:** trade-off matrix + interview card  

## The Problem

Interview answers often jump from "we need a database" to a product name without first naming access patterns, query flexibility, transactional needs, or ownership boundaries.

This lesson builds a sharper framing:

- use **relational stores** when multi-row invariants, strong constraints, and predictable joins dominate
- use **key-value stores** when lookup-by-key, low latency, and simple access paths dominate
- use **document stores** when object-shaped data evolves quickly and query flexibility matters more than rigid relational modeling

## Clarify

- What are the top three reads and writes by volume?
- Are there cross-entity invariants that must hold synchronously?
- Will most reads be primary-key lookups, bounded secondary queries, or ad hoc exploration?
- Does one team own the schema tightly, or will many producers evolve the shape over time?

## Requirements

### Functional

- Choose a primary storage model for a given workload.
- Explain what queries are first-class and which are intentionally expensive or unsupported.
- Identify when one system of record should be paired with a derived index or cache.

### Non-functional

- Keep latency, consistency, and operational complexity visible.
- Avoid accidental over-flexibility that creates unbounded queries.
- Make future growth and migration pressure explicit.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak reads | 120K req/s | differentiates key lookup paths from scan-heavy systems |
| Peak writes | 25K req/s | exposes write amplification and transaction cost |
| Data size | 80 TB primary over 3 years | affects compaction, tiering, and index budget |
| Peak factor | 4x during launches | stresses lock contention and hot partitions |
| Rough cost | primary store + replicas + indexes + backup | store choice is really a cost and operations choice too |

## Architecture

Think in terms of the system of record and derived views:

```text
clients
  -> service layer
     -> primary system of record
     -> optional derived search/index/cache path
```

A strong answer usually says:

1. which store owns correctness
2. which access paths are optimized directly
3. which queries move to derived systems later

## Data Model & APIs

Typical examples:

- **Relational:** `users`, `orders`, `payments`; foreign keys and transaction boundaries matter
- **KV:** `session_id -> blob`, `user_id -> profile cache`, `quota_key -> counters`
- **Document:** `project -> nested settings, permissions summary, feature flags`

Questions to answer explicitly:

- what is the primary key
- which secondary queries are allowed
- what update pattern is expected
- whether the shape is normalized, denormalized, or nested on purpose

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| relational store chosen for mostly key lookups at huge scale | rising tail latency and replica cost for simple reads | move hot lookup path to KV or cache while preserving source of truth |
| KV chosen for workflow with cross-row invariants | reconciliation work and correctness bugs | keep invariant-heavy writes in a transactional store |
| document store used without query discipline | slow secondary queries and index sprawl | cap query shapes and define approved indexes |
| team treats derived index as source of truth | stale or divergent behavior under lag | name one canonical system of record and measure lag explicitly |

## Observability

- metric: read and write volume by access pattern class
- metric: p95/p99 latency by primary key lookup, secondary query, and scan
- metric: index size growth relative to primary data size
- log: rejected or degraded query shapes
- trace: storage path chosen for critical requests
- SLO: primary user path stays within latency target without relying on accidental full scans

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| relational primary for invariant-heavy workflows | strong correctness and mature tooling | scaling joins and write hotspots can get expensive | storing everything in a schemaless document layer |
| KV for narrow lookup path | predictable low latency and operational simplicity | poor fit for flexible queries and joins | forcing simple lookup traffic through SQL joins |
| document model for evolving object shape | fast iteration and natural nested data | harder consistency rules and index discipline | over-normalizing early product schemas |

## Interview It

**Google framing:** "Design storage for a collaboration product with tasks, comments, and permissions." The signal is whether you separate transactional workflow data from derived read views.

**Cloudflare framing:** "Design storage for edge configuration objects and rollout metadata." The signal is whether you distinguish fast key lookups from control-plane auditability and version history.

**Follow-ups:**
1. What changes if one query now needs arbitrary filtering across many fields?
2. What if writes must atomically update two entities with business invariants?
3. What if one access path is 95% of total traffic?
4. What if operators need audit history for every mutation?
5. When should you introduce a second storage system instead of stretching the first one?

## Ship It

- `outputs/tradeoff-matrix-storage-models.md`
- `outputs/interview-card-storage-models.md`

## Exercises

1. **Easy** — Pick a storage model for sessions, feature flags, and billing invoices.  
2. **Medium** — Redesign a document-heavy prototype that now needs transactional order placement.  
3. **Hard** — Explain a mixed relational + KV + search architecture for a marketplace product without sounding overbuilt.  

## Further Reading

- [Amazon DynamoDB paper](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) — useful when discussing key-value trade-offs at scale  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — strong grounding for model and access-pattern trade-offs  
