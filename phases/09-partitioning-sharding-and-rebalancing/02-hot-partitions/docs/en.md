# Hot Partitions and Skew

> Average load is comforting right up until one partition melts.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Diagnose hotspotting early, explain why averages hide skew, and choose mitigations that match whether the heat comes from keys, tenants, writes, or time windows.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/01-shard-key`, `02-estimation-and-cost/05-burstiness`
**Estimated time:** ~75 min
**Primary artifact:** failure checklist + interview card

## The Problem

Distributed systems fail unevenly. One celebrity post, one runaway tenant, one minute-aligned write pattern, or one overused key can make a single partition the real system bottleneck while fleet-wide dashboards still look normal.

This lesson is about recognizing skew for what it is:

- a routing problem
- a workload-shape problem
- often a product behavior problem

Hot partitions matter because they create tail latency, queueing, noisy neighbors, and cascading retries long before average utilization looks scary.

## Clarify

- Is the heat coming from reads, writes, metadata updates, or coordination locks?
- Is the skew tied to a hot key, a hot tenant, a time bucket, or one access pattern?
- Can requests be spread safely, or does correctness require serialized ownership?
- Is the system already using caches, batching, or queues that could hide or amplify the hotspot?

If the interviewer gives no extra detail, assume a multi-tenant service with a few outsized customers and bursty read traffic around a small set of hot objects.

## Requirements

### Functional

- Detect skew quickly and attribute it to the right key class or tenant.
- Keep critical traffic available even while one partition is overloaded.
- Support mitigations such as key splitting, request coalescing, batching, or dedicated placement.

### Non-functional

- Avoid retry storms and cross-tenant blast radius.
- Preserve predictable latency for the median and tail user.
- Keep mitigation operationally simple enough to deploy under incident pressure.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Fleet read QPS | 400K req/s | average load looks safe at fleet level |
| Hottest partition share | 22% of all reads | enough to melt one shard despite healthy global averages |
| Burst factor | 8x in 60 seconds | exposes queue depth and retry amplification |
| Largest tenant skew | 50x median tenant | noisy-neighbor control becomes mandatory |
| Rough cost | overprovisioning + mitigation machinery | hotspot controls can cost more than steady-state balance plans |

## Architecture

Respond in layers:

1. **Detect** whether heat is concentrated by tenant, key, region, or time bucket.
2. **Contain** the blast radius with per-partition backpressure and per-tenant isolation.
3. **Mitigate** using the least invasive change that fits the cause.

Common mitigations:

- hot-read keys: cache, coalescing, replica fanout, stale serving
- hot-write keys: key salting, finer buckets, batching, async materialization
- hot tenants: dedicated shards or tenant splits
- time-skewed writes: spread timestamps, pre-create buckets, queue smoothing

## Data Model & APIs

Useful metadata to preserve:

- tenant or customer ID on every serving request
- partition identifier in logs and traces
- key prefix or bucket classification on hot paths

Useful control-plane interfaces:

- `MoveTenant(tenant_id, target_pool)`
- `SplitKeyRange(range_id, split_point)`
- `ThrottleTenant(tenant_id, limit)`

The senior move is to connect mitigation directly to the cause. Do not answer "add more shards" unless you can explain why the heat will actually spread.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one hot key causes read pileup | top-key hit rate and per-partition queue depth spike | caching, coalescing, replication fanout |
| write hotspot on monotonic key | single partition commit latency rises while others idle | bucket or salt writes, then reassemble on read |
| one tenant starves neighbors | per-tenant share and error budget burn diverge sharply | dedicated placement, quotas, tenant throttles |
| retries amplify hotspot | request rate grows faster than user traffic | backoff, retry budgets, admission control |

## Observability

- metric: top-N partitions by QPS, CPU, queue depth, and p99 latency
- metric: tenant and key concentration ratio, such as top 1% share of load
- metric: retry rate and shed rate by partition
- log: tenant ID, partition ID, and hotspot classification for throttled or failed requests
- trace: request path tagged with partition ownership and wait time
- SLO: no single partition incident should silently consume the whole service error budget before alerting

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| overprovision hottest shards | fastest short-term relief | poor cost efficiency | pretending averages are enough |
| key splitting or salting | spreads concentrated writes | harder reads and operational complexity | scaling only vertically on one owner |
| dedicated tenant isolation | protects neighbors | raises placement and migration cost | one shared pool for every customer |

## Interview It

**Google framing:** "Your datastore is healthy overall, but p99 latency spikes during product launches." The signal is whether you investigate skew instead of asking only for more hardware.

**Cloudflare framing:** "One customer or object starts dominating a global edge-backed control plane." The signal is whether you reason about noisy neighbors and protected shared infrastructure.

**Follow-ups:**
1. How do you distinguish a hot key from a hot tenant quickly?
2. What if the hotspot is caused by retries rather than organic traffic?
3. What if read caching is not allowed because freshness is strict?
4. How do you protect small tenants during a large customer incident?
5. When is temporary overprovisioning the right answer?

## Ship It

- `outputs/failure-checklist-hot-partitions.md`
- `outputs/interview-card-hot-partitions.md`

## Exercises

1. **Easy** — Name three different causes of skew and one mitigation for each.
2. **Medium** — Redesign a write-heavy counter service that hotspots on one key per minute.
3. **Hard** — Explain a noisy-neighbor strategy for a multi-tenant analytics platform where the largest tenant is 80x the median.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — useful for understanding why uneven slowdowns dominate perceived performance
- [Google SRE books](https://sre.google/books/) — practical thinking on overload, isolation, and incident response
