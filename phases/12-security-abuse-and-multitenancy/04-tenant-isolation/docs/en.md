# Tenant Isolation and Blast Radius

> Multitenancy is successful only when one customer's growth, bugs, or compromise do not quietly become everyone else's problem.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Reason about multitenant isolation across compute, storage, cache, queues, and control plane decisions, not just authorization checks.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/05-tenant-isolation`, `10-reliability-retries-and-backpressure/07-bulkheads`, `12-security-abuse-and-multitenancy/01-auth-and-trust`
**Estimated time:** ~75 min
**Primary artifact:** isolation review checklist + interview card

## The Problem

Design a multitenant platform for APIs and background processing where free-tier, self-serve, and enterprise customers share core infrastructure.

Strong answers explain:

- where tenants share resources and where they must be separated
- how noisy neighbors are detected and contained
- which state is logically isolated versus physically isolated
- when premium tenants justify dedicated capacity or partitions

## Clarify

- Are tenants sharing compute only, or also databases, caches, and queues?
- What is the strongest isolation requirement: performance, data security, compliance, or all three?
- Do enterprise tenants need dedicated shards or only stronger quotas and observability?
- Is the dominant risk a noisy neighbor, a software bug, or a compromised tenant account?

Assume shared infrastructure by default with explicit isolation at authz, quota, storage keying, and workload scheduling layers.

## Requirements

### Functional

- Keep tenant data scoped correctly in reads, writes, caches, and async jobs.
- Prevent one tenant from exhausting shared compute or storage.
- Support higher-isolation tiers when justified.

### Non-functional

- Preserve cost efficiency for the long tail of small tenants.
- Bound blast radius for bugs, abuse, or operational mistakes.
- Make isolation guarantees observable and testable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active tenants | 100K | drives policy scale and noisy-neighbor monitoring |
| Top tenant share | 20% of total traffic | hotspot and carve-out planning matter |
| Shared worker pool | 5K concurrent jobs | queue fairness and budget enforcement matter |
| Storage skew | top 1% hold 60% of bytes | tiering and dedicated partitions may be needed |
| Rough cost | shared pool plus carve-outs for premium tiers | isolation is a business and architecture decision |

## Architecture

```text
tenant request
  -> authz and tenant context
  -> per-tenant quotas / budgets
  -> shared or dedicated compute lane
  -> tenant-keyed storage and cache
  -> audit and isolation telemetry
```

Isolation layers:

1. Logical isolation via tenant-scoped authz and storage keys.
2. Performance isolation via quotas, fair schedulers, and bulkheads.
3. Blast-radius isolation via partitions, dedicated lanes, or regional segmentation for large tenants.

## Data Model & APIs

Core entities:

- `TenantClass`
- `IsolationPolicy`
- `QuotaBudget`
- `PlacementDecision`
- `NoisyNeighborEvent`

Useful APIs:

- `CheckTenantBudget(tenant, resource)`
- `PlaceTenantWork(tenantClass, workload)`
- `ResolveStorageNamespace(tenant, object)`
- `EscalateIsolationTier(tenant)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| cache key misses tenant dimension | cross-tenant reads or support incidents | tenant-prefix all shared cache keys and test for leakage |
| one tenant monopolizes workers | queue latency for everyone else rises | fair scheduling, quotas, and dedicated overflow lanes |
| shared database shard becomes dominated by one tenant | shard CPU and storage skew spike | carve out heavyweight tenants or reshard |
| support tooling bypasses tenant scope | audit logs show broad access paths | admin-plane guardrails and scoped support workflows |

## Observability

- metric: per-tenant budget burn, queue wait, and cache share for top tenants
- metric: cross-tenant access-deny anomalies and storage skew
- metric: worker fairness and dedicated-lane spillover
- log: tenant isolation decisions, carve-outs, and support actions
- trace: tenant context propagation through async and sync paths
- SLO: one tenant should not materially degrade unrelated tenants beyond defined fairness thresholds

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| shared by default, isolate heavy tenants selectively | cost-efficient for long tail | requires good telemetry and migration tools | dedicated stack per tenant |
| fair schedulers and quotas | strong performance isolation | scheduling complexity | first-come first-served shared pools |
| tenant-keyed everything | safer data boundaries | more cardinality and developer discipline | implicit tenant scoping |

## Interview It

**Google framing:** "Design a multitenant internal platform or SaaS backend." Expect follow-ups on noisy neighbors, premium tiers, and operational migration of heavy tenants.

**Cloudflare framing:** "Design tenant isolation for a shared edge or platform product." Expect focus on abuse, fairness, and how large customers avoid hurting everyone else.

**Follow-ups:**
1. When do you keep tenants shared versus move them to dedicated capacity?
2. What is your strongest protection against cross-tenant data leaks?
3. How do you keep one tenant's background jobs from starving others?
4. Which telemetry tells you isolation is failing before customers complain?
5. What changes at 10x tenant count?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/isolation-review-checklist.md`
- `outputs/interview-card-tenant-isolation.md`

## Exercises

1. **Easy** — List three places tenant scope is easy to forget in a design.
2. **Medium** — Add a plan for moving one fast-growing tenant to dedicated capacity.
3. **Hard** — Redesign the architecture for regulated enterprise customers that require stricter physical isolation.

## Further Reading

- [Google SRE workbook - handling overload](https://sre.google/workbook/handling-overload/) — useful fairness and isolation thinking
- [AWS SaaS tenant isolation strategies](https://docs.aws.amazon.com/whitepapers/latest/saas-architecture-fundamentals/tenant-isolation.html) — practical multitenancy patterns
