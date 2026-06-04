# Partitioning and Consumer Groups

> Consumer groups scale only when partitioning, ownership, and skew are treated as first-class design inputs.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design partitioned event consumption with explicit reasoning about ordering, rebalance cost, skew, and per-key parallelism limits.
**Prerequisites:** `07-queues-streams-and-workflows/01-queues-vs-streams`, `09-partitioning-sharding-and-rebalancing/01-shard-key`
**Estimated time:** ~75 min
**Primary artifact:** partition planning sheet + assignment simulator

## The Problem

It is easy to say "we will use consumer groups to scale horizontally." The harder questions are:

- what defines the partition key?
- how much ordering do we need?
- what happens when one tenant or entity becomes hot?
- how disruptive are rebalances?

This lesson focuses on the real scaling boundary in stream systems: parallelism is limited by partitioning decisions, not by the number of pods you can launch.

## Clarify

- What ordering is required: global, per tenant, per user, or per aggregate root?
- What is the expected skew between cold keys and hot keys?
- How often do consumers join, leave, or roll?
- Is replay common, and if so, is it selective or whole-topic replay?

## Requirements

### Functional

- Preserve the required ordering scope.
- Scale reads through partitioned parallelism and consumer groups.
- Allow safe recovery when consumers crash or rebalance.

### Non-functional

- Avoid hot partitions becoming the true throughput ceiling.
- Bound rebalance disruption and catch-up time.
- Keep partition ownership and lag debuggable during incidents.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Ingest rate | 500K events/s | enough to require many partitions |
| Partitions | 128 initial, grow to 512 | sets parallelism ceiling and operational overhead |
| Hot-key skew | top 0.5% keys generate 35% of traffic | exposes shard-key quality and hotspot risk |
| Consumer fleet | 16 to 80 instances | makes rebalance behavior meaningful |
| Rough cost | broker partitions + consumer memory + rebalance churn | pushes trade-offs beyond raw throughput |

## Architecture

```text
producers -> topic partitions
          -> consumer group members
          -> offset commits
          -> lag monitoring
```

Design rules:

1. Choose a key that matches the required ordering boundary.
2. Create enough partitions for near-term scale, but not so many that coordination dominates.
3. Measure skew and rebalance time, not just average throughput.
4. Treat partition count changes and consumer rollouts as operational events.

## Data Model & APIs

Key concepts:

- topic
- partition
- consumer group
- member assignment
- committed offset

Useful APIs:

- `Append(topic, key, event)`
- `Poll(group, member_id)`
- `Commit(group, partition, offset)`
- `DescribeLag(group)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hot key saturates one partition | per-partition lag and CPU diverge sharply | change key design, split heavy tenants, or isolate hot keys |
| frequent rebalances stall progress | rebalance duration and consumer idle time spike | cooperative rebalancing and slower rollout cadence |
| too few partitions cap throughput | consumers stay idle while lag grows | increase partitions with migration plan |
| wrong ordering scope chosen | downstream state conflicts appear | restate required ordering and repartition appropriately |

## Observability

- metric: lag by partition, group, and tenant class
- metric: rebalance count and duration
- metric: events per partition and processing skew
- log: assignment changes and offset commit failures
- trace: one message through publish, consume, commit, and side effect
- SLO: bounded lag on hot paths plus bounded rebalance recovery time

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| per-entity partition key | preserves local ordering | hotspot risk for celebrity entities | random key that breaks correctness |
| more partitions up front | more headroom for parallelism | more coordination and metadata overhead | too few partitions with early ceiling |
| cooperative rebalancing | less stop-the-world disruption | more complex coordination | eager full reassignment on every change |

## Interview It

**Google framing:** "Design event consumption for account updates where ordering matters per account." The signal is whether you translate the ordering scope into a key and notice skew risks.

**Cloudflare framing:** "Design log or policy stream consumers running across many POP-adjacent workers." The signal is whether you reason about rebalances, regional skew, and operational visibility.

**Follow-ups:**
1. What if one tenant suddenly becomes 20x hotter than the rest?
2. What if the topic needs to grow from 128 to 512 partitions?
3. What if some consumers are much slower because of downstream rate limits?
4. What if product changes the ordering requirement from per-user to per-team?
5. How do you roll a new consumer version without causing a full stall?

## Ship It

- `outputs/partition-planning-sheet.md`
- `outputs/interview-card-consumer-groups.md`

## Exercises

1. **Easy** — Pick a partition key for a notification stream with per-user ordering.
2. **Medium** — Explain why adding more consumers does not help when partition count is already the bottleneck.
3. **Hard** — Redesign the consumer group after a hot tenant starts dominating one partition.

## Further Reading

- [Kafka consumer design](https://docs.confluent.io/platform/current/clients/consumer.html) — useful details on groups, offsets, and rebalancing
- [Sticky assignor KIP-54](https://cwiki.apache.org/confluence/display/KAFKA/KIP-54+-+Sticky+Partition+Assignment+Strategy) — good context for rebalance trade-offs
