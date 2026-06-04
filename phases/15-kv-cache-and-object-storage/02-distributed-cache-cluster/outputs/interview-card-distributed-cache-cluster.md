---
lesson: 02-distributed-cache-cluster
focus: balanced
---

## Clarify first
- Reconstructable versus durable data
- Working set shape and object size distribution
- Whether misses are cheap, expensive, or dangerous to the origin

## Must-size numbers
- Peak reads per second
- Hit-rate target and origin offload goal
- Hot-set size and tenant skew
- Expected node churn during scale and failure

## Core design
- Consistent hashing with virtual nodes
- Request coalescing on misses
- Explicit invalidation and TTL discipline
- Optional pool segmentation for large or noisy tenants

## Failure probes
- What happens after a node loss?
- How do you avoid mass cold-start misses during rebalance?
- Can negative caching create correctness issues?

## Trade-off summary
- Shared efficiency vs isolation
- Warmup effort vs memory utilization
- Simplicity vs strict freshness
