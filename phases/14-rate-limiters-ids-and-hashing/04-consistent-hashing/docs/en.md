# Consistent Hashing and Ring Rebalancing

> The value of consistent hashing is not the ring diagram. It is the ability to grow or fail without remapping the whole world.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Explain why consistent hashing reduces remap churn, how virtual nodes smooth skew, and what operational checks matter during node adds, removes, and failures.  
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/03-placement`, `09-partitioning-sharding-and-rebalancing/04-rebalancing`, `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`  
**Estimated time:** ~75 min  
**Primary artifact:** failure checklist + placement validator  

## The Problem

Design a placement scheme for keys across a changing fleet of cache nodes, limiter owners, or storage shards. The placement mechanism should minimize reshuffling when capacity changes and avoid making one node absorb disproportionate load.

This is a classic systems interview topic because it connects partitioning theory to operational reality: node churn, hotspot movement, ring skew, and repair cost.

## Clarify

- Are nodes mostly stable, or does the fleet scale frequently?
- Is the keyspace naturally skewed or close to uniform?
- Is minimizing remapped keys more important than perfect balancing?
- Does placement only affect cache ownership, or durable storage too?

If the interviewer is vague, assume a medium-churn fleet with read-heavy traffic, occasional node failures, and enough skew that virtual nodes matter.

## Requirements

### Functional

- Map keys to owners with bounded remap during topology change.
- Support node add, remove, and temporary failure.
- Make placement explainable for debugging and capacity work.

### Non-functional

- Keep ownership imbalance within an acceptable range.
- Limit data movement during rebalancing.
- Avoid making recovery operations unpredictable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Total keys | 500M logical keys | rebalancing can move a lot of state |
| Fleet size | 100 nodes | ring balance and virtual-node count matter |
| Growth step | add 5 to 10 nodes at once | remap percentage should stay bounded |
| Hot-key skew | top 1% keys carry 30% of traffic | placement uniformity alone is not enough |
| Rough cost | metadata churn + state migration + warmup misses | remap efficiency affects both latency and spend |

## Architecture

```text
key
  -> hash(key)
  -> walk clockwise ring
  -> choose owner vnode
  -> resolve physical node
```

Important components:

1. **Ring builder** assigns tokens or virtual nodes.
2. **Placement resolver** maps keys to owners.
3. **Rebalance planner** predicts movement before rollout.
4. **Health-aware control plane** avoids routing to dead nodes.

## Data Model & APIs

Useful state:

```text
ring -> {
  version,
  virtual_nodes,
  node_weights,
  health
}
```

Useful APIs:

- `BuildRing(nodes, vnode_count)`
- `OwnerForKey(key)`
- `PlanRebalance(old_ring, new_ring)`
- `ExplainPlacement(key, ring_version)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| too few virtual nodes create imbalance | ownership histogram drift | increase vnode count or weight nodes |
| node removal causes cold-cache storm | miss spike after rebalance | staged rollout and prewarming |
| health signal lags real failure | requests still routed to dead node | faster failure detection and fallback owner policy |
| consistent hashing is used as a magic answer to hot keys | one node still overloaded by a few keys | combine with replication or hot-key isolation |

## Observability

- metric: ownership percentage by physical node
- metric: remap percentage for planned topology changes
- metric: cache miss or migration volume after rollout
- log: sampled placement explanations for hot keys
- alert: ring version skew across the fleet

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| consistent hashing | bounded remap during change | more metadata than modulo hashing | modulo that remaps almost everything |
| many virtual nodes | smoother distribution | more control-plane state | one token per node |
| weighted placement | better heterogenous capacity usage | harder to reason about imbalance | assume all nodes are identical |

## Interview It

**Google framing:** "How do you place data or cache ownership onto a changing fleet?" Expect questions about rebalance cost and how you validate a topology change safely.

**Cloudflare framing:** "How do you assign keys or traffic slices across an edge-adjacent fleet?" Expect pressure on node churn, health signaling, and hot-key exceptions.

**Follow-ups:**
1. What changes if some nodes have 2x the capacity?
2. How do you keep node addition from creating a cold-start storm?
3. What if one key is so hot that ring balance no longer helps?
4. What metrics prove the rebalance was safe?
5. When is rendezvous hashing a cleaner answer?

## Ship It

- `outputs/failure-checklist-consistent-hashing.md`
- `outputs/interview-card-consistent-hashing.md`

## Exercises

1. **Easy** — Explain why modulo hashing is painful during fleet growth.
2. **Medium** — Design a rebalance rollout for a cache cluster adding 10 nodes.
3. **Hard** — Redesign for heterogenous node sizes and one pathological hot key.

## Further Reading

- [Consistent hashing and random trees](https://www.cs.princeton.edu/courses/archive/fall09/cos518/papers/chash.pdf) — classic paper behind the idea  
- [AWS Builders Library - Shuffle sharding](https://aws.amazon.com/builders-library/workload-isolation-using-shuffle-sharding/) — related placement thinking for blast-radius control  
