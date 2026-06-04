# Unique ID Generator: Snowflake, DB, and Random IDs

> ID generation is a design problem about ordering, coordination, and failure recovery, not just uniqueness.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Compare database sequences, Snowflake-style IDs, and random IDs, then validate whether a proposed ID plan actually matches ordering, throughput, and operational constraints.  
**Prerequisites:** `02-estimation-and-cost/01-qps-and-request-mix`, `08-consistency-replication-and-transactions/07-time-assumptions`, `09-partitioning-sharding-and-rebalancing/01-shard-key`  
**Estimated time:** ~75 min  
**Primary artifact:** ID design review checklist + validator  

## The Problem

Design an ID strategy for a high-scale service. Product teams want globally unique identifiers, but they also often want ordering, locality, debugging ergonomics, and safe generation during failures.

This lesson exists because "just use UUIDs" and "just use auto-increment" both ignore workload shape and operational constraints.

## Clarify

- Do IDs need to be globally unique, sortable, or only unique inside a shard?
- Is the system write-heavy enough that a central sequence becomes a bottleneck?
- Are clients allowed to generate IDs offline or only trusted servers?
- Is clock skew acceptable, or would non-monotonic IDs break downstream assumptions?

If details are missing, assume a write-heavy backend with multiple regions, operator need for rough temporal ordering, and no requirement for perfectly gapless sequences.

## Requirements

### Functional

- Generate unique IDs at high write throughput.
- Support a model for shard ownership or generator identity.
- Avoid turning ID creation into the system bottleneck.

### Non-functional

- Bound or explain ordering anomalies.
- Keep operational recovery understandable after node restart or clock problems.
- Avoid leaking more predictability than the product can tolerate.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak writes | 250K inserts/s globally | central generators may bottleneck |
| Regions | 3 active regions | coordination and ordering semantics get harder |
| Node count | 200 writers | generator ID space and drift handling matter |
| Restart frequency | dozens per week | generator recovery must not create collisions |
| Rough cost | local generation + worker coordination | ordering guarantees are the expensive part |

## Architecture

```text
writer
  -> choose ID strategy
     -> db sequence or allocator
     -> snowflake-style generator
     -> random ID generator
  -> persist entity with generated ID
```

Rule of thumb:

1. **DB sequence** works well when throughput is modest and strict central ordering matters.
2. **Snowflake-style IDs** balance locality, scale, and rough time order.
3. **Random IDs** remove coordination but sacrifice ordering and can complicate storage locality.

## Data Model & APIs

Useful state:

```text
generator -> {
  strategy,
  region_id,
  worker_id,
  last_timestamp_ms,
  sequence
}
```

Useful APIs:

- `Generate(now_ms)`
- `ValidatePlan(plan)`
- `ExplainCollisionRisk(plan)`
- `ReserveWorkerID(node)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| clock moves backward on Snowflake node | timestamp regression alert | block generation briefly or move to a safe sequence reserve |
| central sequence becomes hot or unavailable | elevated insert latency and allocator saturation | cache ranges or move to distributed generation |
| random IDs destroy storage locality | write amplification and fragmented index behavior | separate public IDs from storage keys |
| worker IDs collide after restart | duplicate generator ownership | explicit lease or registration of worker identity |

## Observability

- metric: generated IDs per strategy and per region
- metric: timestamp-regression or generator-stall count
- metric: allocator latency for sequence-backed IDs
- log: sampled generator decisions with region and worker identity
- audit: generator lease history for collision investigations

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Snowflake-style IDs | scalable and roughly sortable | clock discipline and worker coordination | central DB sequence for all writes |
| random public IDs | no hot central allocator and hard to guess | no order and poorer locality | exposing internal sequence IDs |
| decoupled storage key and public ID | better indexing and safer external surface | extra mapping complexity | one ID serves every purpose |

## Interview It

**Google framing:** "Design an ID service for a write-heavy product." Expect questions about ordering requirements, replays, and why downstream systems assume too much from IDs.

**Cloudflare framing:** "Generate identifiers safely across many machines and regions." Expect questions about clock issues, worker identity, and whether control-plane registration is required.

**Follow-ups:**
1. What if IDs must sort roughly by creation time?
2. What if clock skew is common in one region?
3. Should the same ID be used as the public resource ID and the storage key?
4. What changes if clients need offline generation?
5. How would you migrate from auto-increment IDs?

## Ship It

- `outputs/design-review-unique-id-generator.md`
- `outputs/tradeoff-matrix-unique-id-generator.md`

## Exercises

1. **Easy** — Pick an ID strategy for an internal admin tool with low write volume.
2. **Medium** — Redesign for a multiregion write path that needs rough time ordering.
3. **Hard** — Support public IDs that should not reveal write volume while preserving storage locality.

## Further Reading

- [Twitter Snowflake notes](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake) — classic framing for time-ordered distributed IDs  
- [UUID revision draft overview](https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-14.html) — useful when discussing modern sortable/random UUID trade-offs  
