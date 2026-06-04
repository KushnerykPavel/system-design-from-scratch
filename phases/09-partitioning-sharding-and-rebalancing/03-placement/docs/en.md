# Consistent Hashing and Placement

> Good placement spreads ownership without remapping the world.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Explain when consistent hashing is the right placement primitive, what virtual nodes buy you, and how to reason about remap cost when nodes are added or removed.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/01-shard-key`, `14-rate-limiters-ids-and-hashing/04-consistent-hashing`
**Estimated time:** ~75 min
**Primary artifact:** placement simulator + interview card

## The Problem

Once data or request ownership spans many machines, you need a placement function. The naive choice, `hash(key) % N`, looks simple until `N` changes and a large fraction of keys move at once.

Consistent hashing exists to make ownership changes cheaper:

- adding one node should move only a slice of keys
- removing one node should redistribute only that node's slice
- virtual nodes should smooth imbalance

The lesson is not "always use a ring." It is "use a placement strategy whose remap behavior matches your operational needs."

## Clarify

- Are we placing requests, cache keys, partitions, or durable storage ranges?
- How often do nodes join or leave, and how much data movement is acceptable each time?
- Do we need weighted placement for heterogeneous capacity?
- Is the main concern balance, remap cost, fault domains, or replica placement?

If unclear, assume a cache or stateless ownership layer where nodes scale regularly and minimizing remapped keys matters.

## Requirements

### Functional

- Route each key to an owner deterministically.
- Add or remove nodes without remapping most keys.
- Support weighting or virtual nodes when capacity differs.

### Non-functional

- Keep load spread reasonably even.
- Bound remap churn during scaling events.
- Make ownership observable enough to debug hotspots and imbalance.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active keys | 500M keys | enough ownership movement to make remap cost expensive |
| Node count | 200 nodes | balancing quality matters at this fleet size |
| Weekly scaling events | 5-20 adds or removes | placement instability becomes operationally visible |
| Acceptable remap on single add | ideally ~1/N of keys | frames why modulus hashing is too disruptive |
| Rough cost | background data copy + cache warmup + control-plane churn | remap cost is part of scaling cost |

## Architecture

Common pattern:

```text
key
  -> hash
  -> ring lookup
  -> owner node / replica set
```

Use cases:

- cache clusters
- request affinity
- partition owner selection
- replica placement seeds

Important senior-level notes:

- virtual nodes reduce imbalance from unlucky hash spacing
- weighting helps mixed-capacity fleets
- consistent hashing helps placement churn, not every query pattern

## Data Model & APIs

Useful control-plane state:

- `Node(id, capacity_weight, fault_domain, state)`
- `PlacementEpoch(version, ring_checksum)`

Useful interfaces:

- `ResolveOwner(key, epoch) -> node`
- `AddNode(node_spec)`
- `DrainNode(node_id)`

If durability matters, pair placement with replication rules. A ring alone does not guarantee safe replica diversity.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| modulus hashing remaps too many keys on scale-up | cache miss storm or mass copy volume after small capacity change | use consistent hashing or stable range assignment |
| too few virtual nodes creates imbalance | per-node ownership variance remains high | increase vnode count or use weighted placement |
| replicas land in same failure domain | outage takes out primary and backup together | enforce rack, AZ, or region diversity in placement |
| inconsistent ring view across clients | conflicting owners and split traffic | versioned ring distribution and epoch checks |

## Observability

- metric: key ownership share per node and standard deviation across nodes
- metric: percent of keys remapped after a topology change
- metric: cache miss surge or migration bytes after ring updates
- log: ring epoch, node state change, and ownership checksum
- trace: owner resolution and redirect or stale-epoch retries
- SLO: node adds and drains should not cause disproportionate latency or miss spikes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| consistent hashing ring | low remap churn on node changes | more complex routing than modulo hashing | `hash(key) % N` |
| many virtual nodes | smoother balance | larger ring state and update overhead | one token per node |
| weighted placement | better capacity utilization | more control-plane complexity | assuming identical hardware forever |

## Interview It

**Google framing:** "Design a distributed cache cluster that scales horizontally." The signal is whether you talk about remap churn and balancing, not just hashing.

**Cloudflare framing:** "Design request placement across edge caches or ownership nodes." The signal is whether you understand fault domains, warmup cost, and topology churn.

**Follow-ups:**
1. What if one node has twice the memory of the others?
2. What if a single ring view mismatch causes conflicting owners?
3. How do you keep replicas out of the same availability zone?
4. What if adding one node still causes visible imbalance?
5. When would range-based placement beat a ring?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-consistent-hashing.md`
- `outputs/observability-checklist-placement.md`

## Exercises

1. **Easy** — Compare expected remap behavior of modulo hashing and consistent hashing when adding one node.
2. **Medium** — Explain how virtual nodes improve balance on a 20-node cluster.
3. **Hard** — Design weighted placement with fault-domain-aware replicas across three availability zones.

## Further Reading

- [Dynamo paper](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) — classic reference for ring-based placement
- [Memcached at Facebook](https://engineering.fb.com/2013/04/15/core-infra/scaling-memcache-at-facebook/) — practical lessons about placement, churn, and hot keys
