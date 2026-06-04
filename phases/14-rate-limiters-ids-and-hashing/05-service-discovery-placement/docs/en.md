# Service Discovery and Placement Decisions

> Discovery answers the question "where can I send this request?" Placement answers the harder question "where should I send it right now?"

**Type:** Learn  
**Company focus:** Cloudflare  
**Learning goal:** Distinguish service discovery from placement policy, then choose between client-side, proxy-side, and control-plane-assisted routing based on freshness, latency, and blast-radius constraints.  
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/03-placement`, `10-reliability-retries-and-backpressure/07-bulkheads`, `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`  
**Estimated time:** ~60 min  
**Primary artifact:** discovery design review + observability checklist  

## The Problem

Design how services find healthy endpoints and how they choose among them. The challenge is not merely storing a service registry. The challenge is making routing decisions that reflect health, locality, capacity, and rollout intent without turning the control plane into the bottleneck.

## Clarify

- Are callers inside one region, many regions, or edge POPs worldwide?
- Does the caller need only a healthy endpoint list, or a richer placement policy?
- How fresh do health signals need to be before stale routing becomes dangerous?
- Should decisions happen in clients, sidecars, or a centralized proxy tier?

If no extra detail is given, assume a multi-region service fleet with low-latency callers and frequent deploy or failover events.

## Requirements

### Functional

- Discover healthy instances for a named service.
- Route requests using locality, health, and capacity signals.
- Support controlled rollout and fast withdrawal of bad endpoints.

### Non-functional

- Keep data-plane lookup latency low.
- Bound control-plane dependency during incidents.
- Make routing decisions observable and explainable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Services | 500 internal services | registry scale and change rate matter |
| Instances | 20K active endpoints | watch propagation and cache size matter |
| Placement updates | thousands per minute during deploys | stale config risk becomes real |
| Cross-region callers | 40% of traffic | locality and fallback logic matter |
| Rough cost | registry infra + watchers + proxy/client cache | richer placement means more control-plane complexity |

## Architecture

```text
service registry
  -> health and capacity feeds
  -> policy compiler
  -> client or proxy cache
  -> placement decision on request path
```

Common patterns:

1. **Client-side discovery** lowers central proxy cost but pushes policy into many callers.
2. **Proxy or sidecar discovery** centralizes behavior but adds extra hops or resource cost.
3. **Control-plane hints with local cached decisions** often give the best balance at scale.

## Data Model & APIs

Useful entities:

- `ServiceEndpoint`
- `PlacementPolicy`
- `HealthSignal`
- `LocalityClass`
- `RolloutVersion`

Useful APIs:

- `Resolve(service, caller_context)`
- `Watch(service)`
- `UpdateEndpointHealth(endpoint, signal)`
- `ExplainPlacement(request_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| registry is healthy but cached endpoint list is stale | version skew and routing mismatches | TTLs, push invalidation, and fast withdraw path |
| all callers stampede one recovered instance | sudden endpoint saturation | jitter, weighted ramp-up, and concurrency guards |
| placement ignores capacity and only chases locality | low latency but high error rate | couple routing to load and error signals |
| control plane outage breaks data plane | lookup failures during incident | cache placement locally and degrade gracefully |

## Observability

- metric: placement latency and cache hit ratio
- metric: endpoint version skew and stale-cache age
- metric: request distribution by locality and endpoint
- metric: health withdrawal propagation delay
- log: sampled routing explanations with locality, health, and policy version

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| local cached endpoint sets | resilient low-latency lookups | bounded staleness risk | live registry lookup on every request |
| proxy-side policy enforcement | consistent routing behavior | extra hop or sidecar overhead | every client implements its own routing |
| capacity-aware placement | better overload behavior | more signals and tuning complexity | locality-only routing |

## Interview It

**Google framing:** "How do microservices discover and choose healthy backends?" Expect follow-ups on stale config and rollout safety.

**Cloudflare framing:** "How should globally distributed traffic choose service endpoints?" Expect pressure on locality, fast failure withdrawal, and control-plane versus data-plane separation.

**Follow-ups:**
1. What if endpoint health changes faster than config propagation?
2. Should edge services call the registry directly?
3. What changes during a canary rollout?
4. How do you prevent a recovered instance from taking too much traffic too fast?
5. What makes a placement decision explainable?

## Ship It

- `outputs/design-review-service-discovery-placement.md`
- `outputs/observability-checklist-service-discovery-placement.md`

## Exercises

1. **Easy** — Compare client-side discovery with a centralized proxy.
2. **Medium** — Redesign for multiregion callers with strict locality preference.
3. **Hard** — Support fast global endpoint withdrawal without making the registry a hot dependency.

## Further Reading

- [Envoy xDS protocol docs](https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol) — good reference point for control-plane to data-plane config distribution  
- [Google SRE book](https://sre.google/sre-book/) — useful background for load balancing and resilient service communication  
