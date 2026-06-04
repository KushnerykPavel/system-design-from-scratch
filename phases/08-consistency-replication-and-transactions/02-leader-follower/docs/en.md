# Leader-Follower Replication

> Leader-follower is not "one writer, many readers" and done. The real answer is who reads where under lag, failover, and partial partitions.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Evaluate leader-follower topologies by durability, failover behavior, read routing, and operational risk.
**Prerequisites:** `01-consistency-spectrum`, `05-storage-indexing-and-access-patterns/01-storage-models`, `07-queues-streams-and-workflows/06-outbox-and-cdc`
**Estimated time:** ~75 min
**Primary artifact:** topology review card + failover checklist

## The Problem

Leader-follower replication is often the first answer candidates reach for when they want scalable reads and simpler writes. That is fine, but the interviewer is listening for the hidden questions:

- how many acknowledgements are needed before commit
- whether follower reads are allowed on critical paths
- what freshness is promised during lag
- how leader failover works without split brain or silent data loss

This lesson gives you a practical way to discuss those trade-offs.

## Clarify

- Is the leader single-region, zonal, or globally replicated?
- Are follower reads allowed for all entities or only convenience data?
- What write loss is acceptable during leader failure?
- Must failover be automatic, or can operators gate it?

## Requirements

### Functional

- Support durable writes through a leader path.
- Scale read traffic with followers where freshness allows.
- Fail over to a healthy leader without ambiguous ownership.

### Non-functional

- Keep failover understandable and auditable.
- Bound replica lag before follower reads become product-visible bugs.
- Avoid split brain and unclear commit semantics.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Write QPS | 20K req/s | enough to make ack policy visible |
| Read QPS | 250K req/s | follower routing becomes attractive |
| Followers per leader | 3-5 | affects lag, fan-out, and failure options |
| Failover target | <60 seconds | constrains election and readiness checks |
| Rough cost | extra replicas, replication bandwidth, failover tooling | shows why stronger durability is not free |

## Architecture

Typical flow:

```text
client
  -> leader
     -> local durable write
     -> replication to followers
     -> ack based on commit policy
```

Strong answer sequence:

1. Name the leader as the write serialization point.
2. Explain the commit rule, such as local durable write only versus leader plus one follower.
3. Separate follower-read entities from leader-read entities.
4. Describe failover gating, including replication position and fencing.

## Data Model & APIs

Helpful replication metadata:

```text
replica -> {
  role,
  commit_index,
  applied_index,
  region,
  healthy
}
```

Useful surfaces:

- `Write(record, required_acks)`
- `Read(record_id, consistency=leader|follower|min_version)`
- `Promote(replica_id, fence_epoch)`
- `ReplicationStatus()`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| follower serves stale critical read | freshness age or version lag alert | route critical reads to leader or require min version |
| leader fails before followers persist a write | commit gap after failover | stronger ack policy for critical paths |
| automatic failover promotes stale node | commit index mismatch | promotion fencing and lag thresholds |
| split brain during partition | dual-leader heartbeat or fencing violation | single writer lease, quorum-based promotion, operator gate |

## Observability

- metric: replication lag by follower in bytes, time, or log index
- metric: leader commit latency by ack policy
- metric: follower-read freshness age by entity class
- log: failover decisions with prior leader epoch and promoted replica index
- trace: write path with local commit and replication acknowledgment milestones
- SLO: critical read freshness target plus write durability target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| follower reads for noncritical data | cheaper read scale | stale windows become real | all reads through leader |
| stronger write ack policy | better durability on failover | slower writes and more coupling to replica health | leader-local acks for everything |
| gated failover | lower split-brain risk | slower recovery in some incidents | blind auto-promotion |

## Interview It

**Google framing:** "Design a strongly owned metadata service with heavy reads and moderate writes." The signal is whether you explain commit semantics and what reads can safely hit followers.

**Cloudflare framing:** "Design globally distributed policy storage with a single write leader per shard." The signal is whether you reason about read freshness, promotion safety, and fail-safe routing.

**Follow-ups:**
1. Which user paths can read from followers?
2. What changes if one region becomes isolated but still serves reads?
3. What if write durability matters more than p99 latency?
4. How do you explain promotion safety without saying "consensus" and stopping there?
5. What if lag rises only for one follower in a remote region?

## Ship It

- `outputs/interview-card-leader-follower.md`
- `outputs/failover-checklist-leader-follower.md`

## Exercises

1. **Easy** - Mark which reads in a user profile system can safely use followers.
2. **Medium** - Compare local-durable ack versus leader-plus-one-follower ack for order metadata.
3. **Hard** - Redesign for a control-plane policy store where stale follower reads can be dangerous.

## Further Reading

- [Designing Data-Intensive Applications](https://dataintensive.net/) - strong replication and failover background
- [System design notes](https://github.com/liquidslr/system-design-notes) - useful interview framing before replication nuance
