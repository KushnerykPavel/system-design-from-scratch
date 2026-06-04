# Replication Lag and Read Freshness

> Lag is not a storage detail. It is a user experience, enforcement, or billing detail waiting to happen.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Explain how replication lag changes read routing, product guarantees, and incident response.
**Prerequisites:** `01-consistency-spectrum`, `02-leader-follower`, `06-caching-and-invalidation/07-cache-consistency`
**Estimated time:** ~75 min
**Primary artifact:** freshness budget worksheet

## The Problem

Many designs assume replicas exist and that reads can use them, but do not define how much lag is acceptable. That is where real systems break. A replica that is only two seconds behind may be harmless for a public profile page and unacceptable for quota, inventory, or policy enforcement.

This lesson teaches you to discuss lag as a freshness budget with routing consequences.

## Clarify

- Which reads can tolerate seconds of staleness, and which cannot?
- How will the client or server know that a follow-up read needs a minimum version?
- Is lag local to one follower, one region, or a whole replica class?
- Does failover temporarily expand the stale window?

## Requirements

### Functional

- Route reads based on freshness needs, not only topology convenience.
- Detect and surface replica lag explicitly.
- Provide a path for critical reads to bypass stale replicas.

### Non-functional

- Keep freshness contracts measurable.
- Avoid silent lag becoming a correctness or trust incident.
- Balance origin or leader pressure against user-facing consistency needs.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Follower read share | 85% of reads | large enough that freshness policy matters |
| Critical fresh-read share | 5% of reads | low enough to justify selective stronger routing |
| Acceptable lag | 0 ms to 10 s by entity | defines freshness tiers |
| Regional lag spikes | rare but 30-60 s during incidents | forces degraded-mode planning |
| Rough cost | fresh reads raise leader load and latency | keeps the guarantee selective |

## Architecture

A strong freshness design includes:

1. **Lag measurement** in time, log position, or version.
2. **Entity tiers** that define maximum acceptable stale age.
3. **Routing logic** that picks leader, follower, or cached response based on freshness need.
4. **Degraded mode** when all cheap replicas exceed the freshness budget.

Example tiers:

- `fresh-now`: balance, inventory, recent security changes
- `fresh-soon`: user settings, profile metadata
- `stale-ok`: feed ranking, recommendations, analytics views

## Data Model & APIs

Helpful interfaces:

- `Get(id, max_staleness_ms)`
- `GetAfterVersion(id, min_version)`
- `ReplicaStatus(region, replica_id)`

Helpful metadata:

```text
read_result -> {
  value,
  version,
  replica_id,
  freshness_age_ms
}
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| replica used outside freshness budget | freshness age breaches for entity tier | reroute to leader or fail safely |
| lag metric exists but is not used in routing | product incidents despite healthy dashboards | make freshness an admission input for reads |
| failover temporarily increases stale reads | mismatch rate rises during promotion | stricter read routing during recovery window |
| one remote region lags badly and hides it | regional skew grows without traffic shift | isolate stale replicas and alert on skew |

## Observability

- metric: freshness age by entity tier and read path
- metric: replica lag by region and replica role
- metric: leader-read fallback rate caused by stale followers
- log: reads denied or rerouted because freshness budget was exceeded
- trace: mutation to follow-up read path with version or staleness checks
- SLO: freshness objective for critical entity classes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| freshness-tiered routing | clear match between product need and cost | more routing logic | all reads treated the same |
| selective leader bypass | protects critical correctness | higher tail latency and leader load | follower reads for every path |
| explicit stale budget | auditable contract | requires more metrics and version metadata | vague "replicas are usually close enough" |

## Interview It

**Google framing:** "Design read scaling for account and inventory metadata." The signal is whether you create freshness tiers instead of one blanket read policy.

**Cloudflare framing:** "Design globally replicated policy or configuration reads." The signal is whether stale reads are treated as an enforcement risk rather than a normal cache miss.

**Follow-ups:**
1. What if only 3% of reads truly need leader freshness?
2. How do you protect read-after-write for one user without sending all traffic to the leader?
3. What if lag measurement itself is delayed?
4. What if the cheapest remote replica is often the stalest?
5. What metric would page you before customers notice?

## Ship It

- `outputs/freshness-budget-worksheet.md`

## Exercises

1. **Easy** - Assign freshness tiers to four common product entities.
2. **Medium** - Design min-version reads after a user changes settings.
3. **Hard** - Explain how a global enforcement system should react when a whole region is 45 seconds behind.

## Further Reading

- [Designing Data-Intensive Applications](https://dataintensive.net/) - useful mental models for lag and stale-read behavior
- [System design notes](https://github.com/liquidslr/system-design-notes) - baseline interview scaffolding for replicated read paths
