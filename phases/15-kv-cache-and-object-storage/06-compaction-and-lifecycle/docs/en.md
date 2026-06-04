# Compaction, GC, and Lifecycle Policies

> Storage engines and storage products both accumulate debt. Compaction and lifecycle are how you keep that debt from becoming an outage or runaway bill.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Reason about background maintenance in storage systems, including compaction debt, tombstone cleanup, retention rules, and lifecycle transitions.  
**Prerequisites:** `05-storage-indexing-and-access-patterns/07-retention-and-deletion`, `10-reliability-retries-and-backpressure/05-async-backpressure`, `15-kv-cache-and-object-storage/03-object-storage`  
**Estimated time:** ~60 min  
**Primary artifact:** lifecycle-policy linter + rollout checklist  

## The Problem

Storage systems rarely fail only because of foreground traffic. They also fail when background maintenance falls behind. In KV stores that means compaction debt and tombstone buildup. In object stores that means transition, expiration, and delete workflows that quietly stop working.

A senior system-design answer should show that background maintenance is part of the architecture, not a janitorial afterthought.

## Clarify

- Is the main pressure on the system write amplification, stale data retention, or both?
- Which deletes are user-visible immediately versus only eventually reclaimed?
- Are there legal hold or compliance constraints on lifecycle deletion?
- How much background work can the platform do before serving latency suffers?

If nothing else is specified, assume a write-heavy LSM-backed store paired with object lifecycle transitions and retention-aware deletion.

## Requirements

### Functional

- Reclaim obsolete versions and tombstones safely.
- Support retention, transition, and deletion lifecycle rules.
- Allow operators to pause or throttle maintenance safely during incidents.

### Non-functional

- Prevent compaction debt from silently degrading latency.
- Keep lifecycle policy execution auditable.
- Bound maintenance impact on foreground reads and writes.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Daily overwrite volume | 15 TB/day | drives tombstones and compaction pressure |
| SSTable growth | 8x per day before merge | compaction backlog can explode under sustained writes |
| Lifecycle actions | 100M transitions or expirations/day | requires scalable asynchronous policy engine |
| Delete retention | 7 to 90 days by class | legal and operational semantics change GC timing |
| Background IO budget | 20 to 30% of disk/network | maintenance competes directly with serving work |

## Architecture

```text
foreground writes
  -> log / memtable / immutable segments
  -> compaction scheduler
  -> tombstone GC after safety window

metadata + policy engine
  -> lifecycle evaluator
  -> transition / expire / delete workers
  -> audit log
```

Core rule:

- reclamation must respect correctness windows first, then chase efficiency

## Data Model & APIs

Lifecycle rule:

```text
scope
transition_after_days
expire_after_days
delete_grace_days
legal_hold_required
```

Useful APIs:

- `EstimateCompactionDebt()`
- `PauseCompaction(reason)`
- `ApplyLifecycleRule(rule)`
- `DryRunLifecycle(rule)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| compaction falls behind after ingest spike | compaction debt and read amplification rise | priority scheduling and temporary ingest shaping |
| tombstones are GC'd too early | resurrected or inconsistent reads after replica lag | enforce GC grace periods tied to repair assumptions |
| lifecycle rule deletes retained data | audit alerts and policy violation checks | dry-run, approval gates, and retention-aware engine |
| background jobs starve serving path | disk queue and p99 latency spikes | bandwidth caps and adaptive throttling |

## Observability

- metric: compaction debt bytes and read amplification
- metric: tombstone age distribution and GC backlog
- metric: lifecycle transition backlog by rule and action type
- metric: policy dry-run diff versus actual actions
- log: destructive lifecycle actions with actor and rule version
- SLO: maintenance stays within a bounded resource budget while keeping debt and backlog under target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| aggressive compaction | lower read amplification and smaller footprint | higher immediate IO cost | letting debt accumulate until latency breaks |
| long GC grace period | safer under lag and repair delay | more storage used by tombstones | early reclamation that risks correctness |
| dry-run before lifecycle rollout | safer destructive-policy changes | slower rollout and more operator ceremony | direct delete policy push in production |

## Interview It

**Google framing:** "Explain how a storage platform avoids accumulating maintenance debt that hurts serving performance." Expect questions on compaction, tombstones, and safe deletion semantics.

**Cloudflare framing:** "Explain lifecycle and retention policies for a high-scale storage product." Expect questions on async backlogs, policy safety, and observability during large migrations.

**Follow-ups:**
1. What if compaction and repair both need the same disks during an incident?
2. How do you set GC grace periods when replica lag is highly variable?
3. What changes when lifecycle transitions involve moving PBs between classes?
4. How do you test destructive retention rules safely?
5. What happens if customers suddenly shorten retention by 90%?

## Ship It

- `outputs/lifecycle-rollout-checklist.md`
- `outputs/failure-checklist-compaction-and-lifecycle.md`
- `outputs/interview-card-compaction-and-lifecycle.md`

## Exercises

1. **Easy** — Pick a safe GC grace period assumption and justify it.
2. **Medium** — Design a dry-run workflow for bucket-wide lifecycle policy changes.
3. **Hard** — Explain how you would prioritize compaction, repair, and lifecycle work during a sustained ingest spike plus node loss.

## Further Reading

- [Bigtable: A Distributed Storage System for Structured Data](https://research.google/pubs/pub27898/) — classic SSTable and compaction-oriented storage framing  
- [Apache Cassandra tombstones](https://cassandra.apache.org/doc/stable/cassandra/managing/operating/compaction/tombstones.html) — concrete example of deletion semantics and GC grace trade-offs  
