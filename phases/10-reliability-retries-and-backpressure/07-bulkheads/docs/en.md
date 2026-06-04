# Bulkheads and Failure Isolation

> Shared capacity is efficient right up until it becomes shared failure.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Partition capacity, dependencies, and control surfaces so one noisy feature, tenant, or dependency failure does not drown the entire system.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/05-tenant-isolation`, `10-reliability-retries-and-backpressure/04-load-shedding`, `10-reliability-retries-and-backpressure/05-async-backpressure`
**Estimated time:** ~60 min
**Primary artifact:** isolation policy evaluator + failure checklist

## The Problem

Many outages are not total failures. They are blast-radius failures:

- one premium tenant overwhelms a shared pool
- one optional feature saturates worker threads
- one cell or POP becomes unhealthy but the rest are fine
- one control-plane dependency slows down the data plane

Bulkheads are the design answer to "what should fail together?" If the answer is "everything," the system is under-partitioned.

## Clarify

- Are we isolating by tenant, feature, dependency, region, or execution pool?
- Which traffic classes deserve dedicated capacity?
- What control-plane dependencies are allowed on the data path?
- Is the system already partitioned into cells, pools, or shards we can leverage?

If the interviewer is vague, assume a multi-tenant service with uneven customer sizes, optional features, and mixed control-plane and data-plane traffic.

## Requirements

### Functional

- Prevent one traffic class or feature from consuming all shared capacity.
- Keep healthy cells or regions serving even when one segment is degraded.
- Make failover and recovery possible without global coordination bottlenecks.

### Non-functional

- Bound blast radius explicitly.
- Preserve utilization that is good enough in normal periods.
- Keep the isolation model explainable to operators and interviewers.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Total fleet capacity | 100 units | makes reserved slices and spare headroom visible |
| Largest tenant share | 18% of peak traffic | one customer can dominate a shared pool |
| Optional feature load | 30% of CPU in spikes | feature-level isolation can preserve core flows |
| Cell count | 12 cells or POP groups | blast radius depends on how failure domains are shaped |
| Rough cost | reserved idle capacity + more pools + more ops | isolation is a deliberate efficiency trade |

## Architecture

Bulkheads can live at several levels:

- separate thread or worker pools
- per-tenant or per-class quotas
- cell-based routing
- dedicated storage or dependency pools

```text
request
  -> traffic classifier
  -> tenant / feature / region bulkhead
  -> isolated worker pool or cell
  -> shared dependencies only where justified
```

## Data Model & APIs

Useful policy fields:

- `pool_name`
- `reserved_capacity`
- `max_burst_capacity`
- `tenant_quota`
- `cell_affinity`

Useful interfaces:

- `RouteToPool(request_class, tenant)`
- `CanBorrowCapacity(from, to)`
- `EvictOptionalTraffic(pool)`

Senior-level detail:

- isolation is not only threads; it also includes queues, caches, and downstream quotas
- borrowing unused capacity can be useful, but must be revocable
- control-plane reads on the serving path deserve special suspicion

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| isolation is only logical, not capacity-backed | one feature still saturates shared CPU or DB pools | dedicated worker pools, quotas, or cells |
| fallback path lands in the same shared bottleneck | degraded traffic still collapses with the primary | isolate fallback and essential dependencies too |
| one cell failover overloads all remaining cells | healthy-cell saturation spikes immediately after reroute | spare headroom and bounded failover plans |
| too much isolation strands idle capacity | utilization stays low while some pools reject | controlled borrowing and periodic policy review |

## Observability

- metric: utilization, rejects, and latency by pool, cell, and tenant class
- metric: capacity borrowing events and duration
- metric: failover traffic concentration across healthy cells
- log: isolation-policy changes, quota violations, and cell reroutes
- trace: tags for selected pool, cell, and degraded-mode path
- SLO: core traffic should remain within target even when optional traffic or one tenant experiences pathological load

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| dedicated pools for critical work | strong blast-radius control | lower average utilization | one global worker pool |
| cell-based isolation | region or POP failures stay bounded | more placement and operational complexity | monolithic global fleet |
| revocable capacity borrowing | better efficiency in calm periods | more policy complexity | hard partitions with no flexibility |

## Interview It

**Google framing:** "One dependency class and one customer tier are noisy. How do you keep the rest of the service healthy?" The signal is whether you can bound blast radius, not just add more autoscaling.

**Cloudflare framing:** "Design an edge platform where one POP, product feature, or customer class should not drown the global fleet." The signal is whether you think in cells, pools, and failure domains.

**Follow-ups:**
1. Which resources need bulkheads first: CPU, threads, queues, or storage?
2. When is dedicated capacity worth the efficiency loss?
3. How do you isolate optional features from the core data plane?
4. What if one cell fails and all traffic must shift?
5. How do you explain the blast-radius boundary clearly in an interview?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/failure-checklist-bulkheads.md`
- `outputs/interview-card-bulkheads.md`

## Exercises

1. **Easy** — Name one concrete resource that should get its own pool in a mixed-priority API service.
2. **Medium** — Explain how you would isolate one very large tenant without redesigning the whole fleet.
3. **Hard** — Redesign the system when a control-plane dependency on the serving path is the biggest shared failure risk.

## Further Reading

- [Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — useful foundation for blast-radius thinking
- [Cloudflare architecture docs](https://developers.cloudflare.com/reference-architecture/architectures/) — relevant edge and failure-domain intuition
