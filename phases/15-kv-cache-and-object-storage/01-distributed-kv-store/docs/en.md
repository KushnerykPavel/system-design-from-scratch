# Distributed Key-Value Store

> A strong storage answer begins with workload shape and failure semantics, not with "use Dynamo-like replication" as a reflex.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design a distributed KV store around access patterns, replica placement, read and write guarantees, and repair behavior under skew and failure.  
**Prerequisites:** `05-storage-indexing-and-access-patterns/01-storage-models`, `08-consistency-replication-and-transactions/03-quorums`, `09-partitioning-sharding-and-rebalancing/04-rebalancing`  
**Estimated time:** ~90 min  
**Primary artifact:** topology validator + trade-off matrix  

## The Problem

Design a distributed key-value store for application metadata, feature flags, session-like state, or configuration objects. The store must serve very high read volume, predictable write behavior, and reasonable durability even when nodes or racks fail.

Interview answers often jump straight to replication buzzwords. Senior answers explain what the keys represent, how large values are, whether reads require freshness, and how the system repairs itself after drift.

## Clarify

- Are values small and immutable-ish, or can they be large and frequently overwritten?
- Is the primary goal low-latency reads, high write throughput, or strong durability?
- What consistency contract do callers actually need: monotonic reads, read-your-writes, or eventual?
- Can clients tolerate stale reads during repair or region failover?

If the interviewer leaves it open, assume small values under 16 KB, read-heavy traffic, moderate write rates, and region-local low-latency service with background repair.

## Requirements

### Functional

- Get and put small values by key.
- Replicate data across failure domains.
- Survive node loss without losing acknowledged durable writes.
- Support range-lite operational APIs such as scan-by-prefix only if the shard key allows it.

### Non-functional

- Keep p99 reads in the low milliseconds within a region.
- Bound inconsistency and make it explainable to downstream teams.
- Rebalance and repair without taking the whole service offline.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak read QPS | 800K req/s | drives in-memory indexes, replica fan-out, and hot-key behavior |
| Peak write QPS | 120K req/s | shapes quorum cost and log durability design |
| Value size | 1 to 8 KB typical, 16 KB max | decides storage engine and compaction pressure |
| Data set | 40 TB logical, 3x replicated | forces attention to repair cost and rebalance time |
| Peak factor | 4x on hot tenants or config keys | hot keys can dominate one shard even when averages look fine |

## Architecture

```text
client
  -> routing layer / partition map
  -> replica set owner for key
     -> write-ahead log
     -> memtable / SSTables
     -> quorum replication
  -> background repair + rebalance + anti-entropy
```

Core design choices:

1. Partition by key hash to spread load.
2. Replicate across racks or AZs, not just processes.
3. Use quorum reads and writes only when the workload truly needs them.
4. Run anti-entropy and hinted handoff so temporary failures do not become permanent divergence.

## Data Model & APIs

Example record:

```text
key
value
version
checksum
write_timestamp
ttl_seconds
```

Useful APIs:

- `Put(key, value, consistency, ttl)`
- `Get(key, consistency)`
- `Delete(key, consistency)`
- `BatchGet(keys, consistency)`
- `ExplainPlacement(key)`

If version conflicts matter, expose conditional writes such as compare-and-set instead of pretending last-write-wins solves business correctness.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one replica set member fails during writes | replica health and quorum failure metrics | hinted handoff, bounded degraded write mode, and fast replacement |
| repair backlog grows after a network partition | anti-entropy lag and replica divergence counters | throttle foreground traffic less than repair, prioritize hot partitions |
| one key becomes globally hot | per-key skew metrics and shard hotspot alarms | replication-aware read fan-out, request coalescing, or key-level caching |
| quorum settings acknowledge unsafe writes | post-failure data loss during incident review | tie durability claims to replica count and failure-domain placement |

## Observability

- metric: read and write latency by consistency mode
- metric: quorum success rate and coordinator timeouts
- metric: replica divergence age and repair backlog bytes
- metric: hottest keys and hottest partitions by read and write share
- log: conditional write failures and read-repair events
- trace: coordinator fan-out and slow replica contribution to p99
- SLO: region-local gets succeed within low-single-digit milliseconds while acknowledged writes meet documented durability guarantees

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| leaderless quorum replication | good regional availability and simpler write ownership | harder reasoning about staleness and repair | single leader for every partition when write availability matters more |
| LSM-style storage | strong write throughput and cheap sequential flush | compaction and read amplification | B-tree storage for append-heavy workloads |
| explicit consistency modes | lets clients pay only for what they need | exposes product teams to nuance | one fake "strong" mode that is expensive and still leaky |

## Interview It

**Google framing:** "Design a metadata KV store used by many internal services." Expect questions on consistency tiers, repair cost, and what product teams are allowed to assume.

**Cloudflare framing:** "Design a globally distributed control-plane KV path with region-local reads." Expect pressure on propagation lag, repair under partitions, and how edge consumers tolerate staleness.

**Follow-ups:**
1. What changes if some keys are configuration and must have monotonic reads?
2. How would you support multi-key conditional updates without promising full relational transactions?
3. What happens when one rack fails mid-rebalance?
4. How do you keep anti-entropy from stealing too much foreground IO?
5. What changes at 10x data size with the same staffing level?

## Ship It

- `outputs/interview-card-distributed-kv-store.md`
- `outputs/tradeoff-matrix-distributed-kv-store.md`
- `outputs/failure-checklist-distributed-kv-store.md`

## Exercises

1. **Easy** — Pick a read and write consistency policy for feature flags and justify it.
2. **Medium** — Redesign the topology for a tenant that is 25% of all traffic.
3. **Hard** — Extend the store to multi-region active-active reads without overselling freshness.

## Further Reading

- [Amazon Dynamo: Amazon's Highly Available Key-value Store](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) — canonical trade-offs for partitioning, replication, and repair  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — practical framing for quorum, storage engines, and replication behavior  
