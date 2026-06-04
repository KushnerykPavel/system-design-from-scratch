# Global Traffic Drill

> Senior answers do not just name global components; they rebalance latency, failover, cache behavior, and consistency when the constraints move.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Practice running a full multi-region and edge design loop with explicit clarification, sizing, architecture choice, failure handling, and redesign under changing traffic constraints.
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive`, `13-multi-region-cdn-and-edge-traffic/03-cdn-layering`, `13-multi-region-cdn-and-edge-traffic/06-geo-consistency`
**Estimated time:** ~60 min
**Primary artifact:** drill worksheet + interview card

## The Problem

Design a global API and content platform for a product with users in North America, Europe, and Asia. Static assets should be fast everywhere, read-heavy API traffic should stay responsive, and the system must survive a regional outage without turning into an undefined recovery scramble.

This drill forces you to combine the whole phase:

- topology choice
- failover objectives
- cache hierarchy
- traffic steering
- edge placement
- geo-consistency

## Clarify

- What portion of traffic is static asset delivery versus API reads and writes?
- Which API actions require fresh reads after writes?
- What regional outage target is the business actually asking for?
- Are there compliance or residency rules for some users?

If you get no extra detail, assume 80% static or cacheable traffic, globally distributed reads, and writes that can be home-region scoped for most entities.

## Requirements

### Functional

- Deliver static assets with low global latency.
- Serve read-heavy API traffic close to users when possible.
- Survive a regional outage with explicit degraded-mode behavior.

### Non-functional

- Keep user-visible latency low on the common path.
- Preserve correctness for workflows that require fresher data.
- Make routing, failover, and cache decisions explainable during incidents.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Global traffic | 400K req/s peak | large enough that POP and cache layers matter |
| Static/cacheable share | 80% | justifies CDN hierarchy and shield layers |
| Write share | 5% of requests | allows selective stronger consistency |
| Peak regional loss | lose one 35% traffic region | forces spare capacity and failover planning |
| Rough cost | edge POP cache + regional app capacity + standby headroom | global resilience requires explicit budget |

## Architecture

```text
client
  -> anycast or DNS ingress
     -> nearest edge POP
        -> edge cache / worker
        -> shield or tiered cache
        -> preferred serving region
           -> stateless API tier
           -> regional caches
           -> home-region write path for owned entities
           -> replicated data and async repair
```

Suggested answer arc:

1. Clarify traffic mix and fresh-read requirements.
2. Size the cacheable share and regional headroom.
3. Choose active-active or active-passive for the right components, not blindly for everything.
4. Name which logic runs at the edge and which data stays regional.
5. Define outage behavior, failover gates, and freshness trade-offs.

## Data Model & APIs

Useful entities:

- `TrafficClass`
- `RegionOwnership`
- `CachePolicy`
- `FailoverObjective`
- `ConsistencyClass`

Useful APIs:

- `RouteRequest(request)`
- `ReadContent(id, freshness_bound)`
- `WriteOwnedEntity(region, payload)`
- `PurgeContent(tag)`
- `ExecuteFailover(plan)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| regional outage overwhelms surviving capacity | post-cutover saturation spike | reserve headroom and predefine degraded mode |
| cache hierarchy hides stale or incorrect data | purge delay and stale-read complaints | explicit freshness classes and purge observability |
| route policy sends traffic to healthy-but-full region | regional saturation despite low ingress latency | couple steering to capacity signals |
| interview answer becomes a bag of global buzzwords | no ownership or failure behavior stated | force explicit write ownership, RTO/RPO, and freshness choices |

## Observability

- metric: latency and success rate by POP, region, and traffic class
- metric: hit ratio by cache layer and origin offload rate
- metric: route-policy version skew and failover timing
- metric: replication lag and freshness-bound breach rate
- log: sampled route, cache, and failover explanations on one request ID
- SLO: define latency, availability, and freshness targets per major traffic class

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| edge + shield for cacheable traffic | low latency and strong origin protection | more invalidation and observability complexity | all traffic goes direct to origin region |
| home-region writes with local reads | simpler correctness and good read latency | some distant writes are slower | fully multi-writer global writes |
| degraded read-only mode during failover | better outage experience | temporary feature reduction | undefined behavior or full outage |

## Interview It

**Google framing:** "Design a global service that serves both content and APIs." Expect follow-ups on consistency classes and recovery realism.

**Cloudflare framing:** "Design the edge and regional traffic path for a global platform." Expect pressure on POP behavior, routing policy, and purge correctness.

**Follow-ups:**
1. What changes if writes rise from 5% to 25% of traffic?
2. What if some countries require data to remain in-region?
3. What if purge must propagate globally within 15 seconds?
4. What if the lowest-latency region is already saturated?
5. How would you explain the design differently for a control plane versus a media path?

## Ship It

- `outputs/design-review-global-traffic-drill.md`
- `outputs/interview-card-global-traffic-drill.md`

## Exercises

1. **Easy** — Run the drill for a static-content-heavy website.
2. **Medium** — Repeat the drill for a read-heavy API with home-region writes.
3. **Hard** — Redesign for a multiregion product with compliance-bound users and aggressive failover objectives.

## Further Reading

- [Google SRE - Handling Overload](https://sre.google/sre-book/handling-overload/) — useful when a surviving region must absorb more than planned
- [Cloudflare - Building resilient systems at global scale](https://blog.cloudflare.com/) — broad source of edge-scale case studies and trade-off intuition
