# Durability Tiers and Data Repair

> Durability is not a boolean. It is a menu of promises, failure assumptions, and repair deadlines.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design storage durability tiers and repair workflows that align cost, recovery speed, and data criticality instead of pretending every byte deserves the same policy.  
**Prerequisites:** `08-consistency-replication-and-transactions/02-leader-follower`, `13-multi-region-cdn-and-edge-traffic/02-failover-and-rto`, `15-kv-cache-and-object-storage/03-object-storage`  
**Estimated time:** ~60 min  
**Primary artifact:** tier-policy evaluator + trade-off matrix  

## The Problem

Storage platforms often serve different classes of data: ephemeral cache backups, user media, audit logs, and compliance-critical records. Treating them all the same either overspends or undershoots the real durability promise.

This lesson focuses on the system-design skill of matching data classes to replication, erasure coding, repair urgency, and geo redundancy while staying honest about failure assumptions.

## Clarify

- Which data classes are rebuildable, user-critical, or legally protected?
- What failure domains matter: disk, host, rack, AZ, or region?
- What is the maximum acceptable window of vulnerability before repair must restore redundancy?
- Is availability during repair as important as long-term durability?

If the prompt is broad, assume three tiers: hot operational data, standard user content, and compliance-critical archival data.

## Requirements

### Functional

- Offer multiple durability tiers with documented guarantees.
- Detect lost or degraded redundancy quickly.
- Repair damaged objects or shards back to policy.

### Non-functional

- Keep cost proportional to data value.
- Avoid creating long windows of under-replicated data.
- Make durability claims auditable and understandable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Hot tier volume | 500 TB | faster repair and high availability may justify full replicas |
| Standard tier volume | 20 PB | cost pushes toward mixed replication and erasure coding |
| Archive tier volume | 80 PB | geo durability and long repair windows dominate economics |
| Daily degraded events | thousands of component failures | automation, not heroics, must drive repair |
| Repair budget | restore policy within minutes to hours by tier | durability is a time-based promise |

## Architecture

```text
write path
  -> choose durability tier
  -> placement policy
  -> replication or erasure coding
  -> durability auditor
  -> repair scheduler
  -> rebuild / replicate workers
```

Key principle:

- durability policy and repair policy are inseparable

A storage platform that promises three copies but repairs slowly during common failures is weaker than a platform with fewer copies but faster and better-isolated repair.

## Data Model & APIs

Tier policy:

```text
tier_name
placement_scope
replica_count
erasure_scheme
geo_redundant
max_repair_hours
```

Useful APIs:

- `PutObject(object, durability_tier)`
- `ExplainDurability(object_id)`
- `ListDegradedReplicas(tier)`
- `RepairNow(object_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| degraded replicas linger for days | under-replicated object age metrics | tier-based repair queues and paging thresholds |
| archive tier claims geo durability but replicas share a control-plane dependency | blast-radius review and regional fault drills | independence audits and regional restore tests |
| repair traffic overwhelms hot serving paths | IO saturation and read latency spikes | background bandwidth caps and priority scheduling |
| cheap tier accidentally stores irreplaceable data | policy mismatch audit and customer incidents | explicit class selection plus safe defaults |

## Observability

- metric: objects below target durability by tier and age
- metric: repair backlog bytes and mean time to restored redundancy
- metric: restore success rate during regional or AZ failure drills
- metric: cost by tier versus stored logical bytes
- log: tier-policy overrides and repair escalations
- SLO: degraded objects return to policy within documented time for each tier

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| full replicas for hot tier | simple reads and fast repair | higher storage cost | erasure coding for latency-critical hot path |
| erasure coding for colder large data | lower storage overhead | slower rebuild and more CPU/network during repair | full replication for multi-PB cold storage |
| geo redundancy only for selected tiers | cost aligned with business need | multiple durability classes to explain | forcing every object into cross-region replication |

## Interview It

**Google framing:** "Design durability classes for a shared storage platform used by many products." Expect questions about cost justification and honest repair-time guarantees.

**Cloudflare framing:** "Design storage classes for edge-adjacent durable content and control data." Expect questions about failure domains, rebuild traffic, and regional independence.

**Follow-ups:**
1. When is erasure coding the wrong answer despite lower storage cost?
2. How do you prove a geo-redundant tier is truly independent?
3. What if repair bandwidth is constrained during an outage?
4. How do you migrate data between durability tiers safely?
5. What changes if customers can choose their own durability class?

## Ship It

- `outputs/tradeoff-matrix-durability-tiers.md`
- `outputs/failure-checklist-durability-tiers.md`
- `outputs/interview-card-durability-tiers.md`

## Exercises

1. **Easy** — Define a safe default durability tier for user uploads and justify it.
2. **Medium** — Redesign the archive tier under a hard cost cap.
3. **Hard** — Explain how repair prioritization should change during a regional incident affecting multiple tiers at once.

## Further Reading

- [Backblaze Vault: Cloud Storage Architecture](https://www.backblaze.com/blog/vault-cloud-storage-architecture/) — practical erasure coding and durability trade-offs  
- [Google File System paper](https://research.google/pubs/pub51/) — classic framing for repair, replication, and large-scale storage operations  
