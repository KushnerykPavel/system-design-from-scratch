# Tenant Isolation and Noisy Neighbors

> Shared infrastructure is only multi-tenant if small customers still survive big ones.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Design shard pools that preserve fairness and predictable performance under uneven tenant behavior, and explain when to escalate from quotas to dedicated placement.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/02-hot-partitions`, `12-security-abuse-and-multitenancy/04-tenant-isolation`
**Estimated time:** ~60 min
**Primary artifact:** design-review prompt + interview card

## The Problem

Multi-tenant systems save cost by pooling infrastructure, but pooled capacity creates shared failure domains. One customer can exhaust CPU, memory, cache working set, compaction budget, or request concurrency and quietly degrade everyone nearby.

This lesson is about designing for fairness:

- place tenants with awareness of size and growth
- enforce quotas and admission control
- create escalation paths to stronger isolation

The weak answer is "we shard by tenant so we are safe." A tenant can still be too large for its shard or too spiky for its pool.

## Clarify

- Is isolation needed for performance, security, compliance, or all three?
- Are tenants similar in size, or does the largest customer differ by orders of magnitude?
- Which resource is most likely to become the shared bottleneck: CPU, storage IO, cache, background compaction, or connection count?
- Can large tenants tolerate dedicated placement or higher cost tiers?

If details are missing, assume a pooled SaaS control plane where the largest tenants are 20-100x larger than the median and latency isolation matters more than perfect cost efficiency.

## Requirements

### Functional

- Prevent one tenant from exhausting shared shard resources.
- Offer a path to move or isolate oversized tenants.
- Preserve per-tenant accounting for rate limits, quotas, and debugging.

### Non-functional

- Keep smaller tenants stable during large-customer bursts.
- Avoid over-isolating every tenant and destroying cost efficiency.
- Make fairness mechanisms transparent enough to explain in interviews and incidents.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Tenant count | 300K active tenants | large long-tail pool with few outsized customers |
| Largest-to-median traffic ratio | 60x | raw pooling becomes risky |
| Shared shard CPU headroom | 30% normal spare | not enough to absorb large bursts safely |
| Bursty tenant launch factor | 10x in 5 minutes | fairness controls must react quickly |
| Rough cost | shared pool + premium isolation tiers | isolation decisions are product and cost decisions, not only technical ones |

## Architecture

A practical isolation ladder:

1. **Fair-share accounting**
   Track per-tenant usage on critical resources.
2. **Soft controls**
   Quotas, token buckets, and priority-aware admission control.
3. **Placement controls**
   Separate heavy tenants, pin premium tenants, or create shard classes.
4. **Dedicated isolation**
   One tenant gets its own shard pool or reserved capacity.

Good answer shape:

```text
request
  -> auth identifies tenant
  -> quota / priority check
  -> routing layer chooses shard class
  -> serving tier emits per-tenant usage metrics
```

## Data Model & APIs

Useful control-plane objects:

- `TenantClass(tenant_id, service_tier, placement_pool, limits)`
- `TenantUsage(tenant_id, qps, cpu_share, storage_bytes, background_jobs)`
- `IsolationAction(tenant_id, action, reason, expires_at)`

Useful interfaces:

- `UpdateTenantLimits(tenant_id, limits)`
- `MoveTenantPool(tenant_id, pool_id)`
- `ThrottleTenant(tenant_id, policy)`

Data model note:

- usage attribution must exist on the hot path or operators will only see shard-level symptoms, not who caused them

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| big tenant thrashes shared cache or CPU | one tenant's usage spike aligns with neighbor latency | quotas, dedicated pool, or stronger cache partitioning |
| fairness only enforced on request QPS | background jobs still starve neighbors | account for async and storage-heavy work too |
| isolation actions are manual and slow | repeated noisy-neighbor incidents with long mitigation time | automate move or throttle triggers |
| too much hard isolation destroys economics | low average utilization across many dedicated pools | use tiered isolation, not one-size-fits-all dedication |

## Observability

- metric: per-tenant resource share on each shard or shard class
- metric: neighbor latency before and after large-tenant bursts
- metric: throttled requests and shed work by tenant class
- log: isolation action reason, policy, and affected pool
- trace: tenant-tagged critical path with wait time or throttling decisions
- SLO: small and medium tenants should not see large error-budget burn because of one oversized neighbor

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| shared pool with quotas | strong cost efficiency | weaker isolation under extreme skew | dedicated infrastructure for every tenant |
| tiered shard classes | better placement for heavy users | more fleet complexity | one homogeneous shard pool |
| dedicated tenant placement | strongest neighbor protection | higher cost and migration overhead | leaving heavy tenants in common pools indefinitely |

## Interview It

**Google framing:** "Design the storage tier for a multi-tenant enterprise product with a few very large customers." The signal is whether you connect fairness, quotas, and placement instead of just talking about sharding.

**Cloudflare framing:** "Protect a shared customer-control system from one noisy zone or account." The signal is whether you reason about shared control-plane safety, not just edge throughput.

**Follow-ups:**
1. What if background compaction, not request QPS, is the real noisy-neighbor source?
2. How do you justify dedicated placement economically?
3. What if premium customers demand stronger SLO isolation?
4. How would you detect slow-burn degradation instead of sharp spikes?
5. When should tenant throttling be automatic versus human-approved?

## Ship It

- `outputs/design-review-tenant-isolation.md`
- `outputs/interview-card-tenant-isolation.md`

## Exercises

1. **Easy** — Name three shared resources that can create noisy-neighbor incidents.
2. **Medium** — Design a tiered isolation policy for free, standard, and enterprise tenants.
3. **Hard** — Explain how to keep a control plane cost-efficient while offering dedicated placement to the largest customers.

## Further Reading

- [Google SRE books](https://sre.google/books/) — good operational framing for fairness and overload control
- [Kubernetes multi-tenancy concepts](https://kubernetes.io/docs/concepts/security/multi-tenancy/) — useful analogies for isolation layers and policy enforcement
