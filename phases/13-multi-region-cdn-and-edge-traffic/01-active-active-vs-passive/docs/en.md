# Active-Active vs Active-Passive

> The topology choice is really a failure, consistency, and operations choice wearing an availability label.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Choose between active-active and active-passive regional topologies by tying failover goals, write semantics, and operator burden to a concrete design.
**Prerequisites:** `08-consistency-replication-and-transactions/04-replication-lag`, `09-partitioning-sharding-and-rebalancing/04-rebalancing`, `10-reliability-retries-and-backpressure/07-bulkheads`
**Estimated time:** ~75 min
**Primary artifact:** topology review checklist + interview card

## The Problem

Design a globally available service that must survive a full-region outage. You need to decide whether both regions serve traffic at the same time or whether one region stays warm until failover.

This comes up constantly in interviews because candidates often say "multi-region" without explaining the real cost: duplicate capacity, conflict handling, replication lag, operator burden, and risk during failover.

## Clarify

- Are writes allowed in every region, or only reads?
- What are the RTO and RPO targets for a regional failure?
- Is the system user-facing latency-sensitive, money-moving, or batch-heavy?
- Can the business tolerate degraded functionality during failover?

If the interviewer is vague, assume read traffic is global, writes are latency-sensitive, and the business wants low RTO but not unlimited complexity.

## Requirements

### Functional

- Serve traffic during loss of a full region.
- Support predictable traffic shifts during failover or maintenance.
- Keep data access semantics understandable for users and operators.

### Non-functional

- Normal-path latency should stay low for the largest user populations.
- Recovery should meet explicit RTO and RPO goals.
- Topology should avoid hidden operator traps such as split-brain or stale failback.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak read QPS | 250K global | determines steady-state regional capacity |
| Peak write QPS | 25K global | drives replication and conflict cost |
| Regional split | 55% / 35% / 10% | shapes how much spare capacity each region needs |
| Failure reserve | 1 extra region worth of headroom | active-passive and active-active size this differently |
| Rough cost | 1.4x to 2x single-region spend | topology choice is fundamentally a cost trade-off |

## Architecture

```text
clients
  -> global DNS / traffic manager
     -> region A
        -> stateless app tier
        -> regional cache
        -> primary data services
     -> region B
        -> stateless app tier
        -> regional cache
        -> replica or peer data services
```

Two common shapes:

1. **Active-passive**
   - one region serves writes
   - passive region is warm, replicated, and ready for promotion
   - simpler write semantics, slower recovery if automation is weak
2. **Active-active**
   - both regions serve at least some live traffic
   - faster traffic absorption and better latency locality
   - requires stronger thinking about data ownership, conflicts, and failback

## Data Model & APIs

Useful control-plane entities:

- `RegionProfile`
- `TrafficMode`
- `FailoverPolicy`
- `ReplicationMode`
- `RecoveryPlan`

Useful APIs:

- `ValidateTopology(profile)`
- `PromoteRegion(region)`
- `DrainRegion(region)`
- `GetRecoveryReadiness()`

For stateful systems, the decisive question is often: which writes are single-writer, multi-writer, or replay-safe?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| passive region is "warm" in name only | failover drills show cold caches or lagging replicas | automate readiness checks and regular game days |
| active-active causes conflicting writes | divergence or reconciliation queue growth | constrain ownership, use per-entity home region, or accept eventual consistency explicitly |
| failback reintroduces stale state | replication lag or version mismatch after recovery | require controlled re-sync and staged failback |
| spare capacity is insufficient during outage | regional saturation and queue growth after cutover | reserve headroom and load-test degraded mode |

## Observability

- metric: per-region success rate, latency, saturation, and queue depth
- metric: replication lag, conflict rate, and replication backlog
- metric: failover readiness score covering health, version, and data freshness
- log: control-plane actions such as drain, promote, and policy override
- trace: request region, chosen dependency region, and fallback path
- SLO: regional failover should meet declared RTO and bounded data-loss targets

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| active-passive for write-heavy systems | simpler correctness model | slower cutover and idle capacity | fully active-active multi-writer everywhere |
| active-active for read-heavy global traffic | better user latency and faster absorption | higher reconciliation and ops complexity | one global primary region |
| warm passive instead of cold standby | realistic recovery | higher steady-state spend | cheaper but operationally risky cold region |

## Interview It

**Google framing:** "Design a service that stays available during a regional outage." Expect pressure on RTO/RPO, write ownership, and failback.

**Cloudflare framing:** "How would you run global traffic across regions without making recovery fragile?" Expect questions on traffic absorption, control-plane safety, and degraded operation.

**Follow-ups:**
1. What changes if the system is read-heavy and writes are rare?
2. What changes if user writes must remain strongly ordered?
3. What if the passive region must become active in under two minutes?
4. How do you prove the passive region is truly ready?
5. What is your failback plan after the failed region returns?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/failure-checklist-active-active-vs-passive.md`
- `outputs/interview-card-active-active-vs-passive.md`

## Exercises

1. **Easy** — Pick active-active or active-passive for a media CDN control plane and defend the choice.
2. **Medium** — Redesign for a payments ledger where only one region may accept writes.
3. **Hard** — Support regional isolation for compliance while still surviving a full-region loss.

## Further Reading

- [Google SRE - Disaster Recovery Testing](https://sre.google/sre-book/testing-reliability/) — useful framing for proving failover claims
- [Cloudflare - How we built Rate Limiting capable of scaling to millions of domains](https://blog.cloudflare.com/how-we-built-rate-limiting-capable-of-scaling-to-millions-of-domains/) — good intuition for distributed control and global traffic realities
