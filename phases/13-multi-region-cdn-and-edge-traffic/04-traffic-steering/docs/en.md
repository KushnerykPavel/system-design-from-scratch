# Traffic Steering, Anycast, and Regional Routing

> Getting users to the right place is a control-plane problem first and a packet-forwarding problem second.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Choose between DNS, anycast, and application-level steering, then explain how health, policy, and overload signals should change routing safely.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/03-placement`, `10-reliability-retries-and-backpressure/04-load-shedding`, `13-multi-region-cdn-and-edge-traffic/02-failover-and-rto`
**Estimated time:** ~75 min
**Primary artifact:** steering policy checklist + interview card

## The Problem

Design how global requests reach regions or POPs. The routing layer must balance latency, health, capacity, and business policy while avoiding dangerous oscillation during incidents.

Strong candidates separate three concerns:

- how traffic first arrives globally
- how unhealthy locations are avoided
- how overload, maintenance, or geography change routing decisions

## Clarify

- Is the service steered primarily by DNS, anycast, client affinity, or application redirects?
- Are requests stateless, session-affine, or write-bound to home regions?
- Is lowest latency always the goal, or do compliance and cost sometimes override it?
- How quickly must traffic move away from a bad region?

If unstated, assume internet-facing traffic, many global POPs, and a need for latency-biased routing with health and overload overrides.

## Requirements

### Functional

- Route users to a healthy location with acceptable latency.
- Support explicit policy overrides for incidents, maintenance, or compliance.
- Preserve session or data affinity where required.

### Non-functional

- Steering changes should not flap under noisy health signals.
- Control-plane propagation should be fast enough for incidents.
- Data-plane routing should keep per-request overhead minimal.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| POP count | 250+ | highlights control-plane fanout |
| Regions | 8 to 20 | changes regional override complexity |
| Control updates | 10K/day normal, bursts during incidents | routing policy must propagate safely |
| Session-affine share | 20% of traffic | shapes reroute limitations |
| Rough cost | control-plane fanout + spare regional capacity | lowest-latency routing is not the only objective |

## Architecture

```text
clients
  -> anycast / DNS ingress
     -> nearest healthy POP
        -> local routing policy
           -> preferred region
           -> alternate region
           -> degraded-mode target
```

Typical steering stack:

1. **Global ingress**
   - DNS or anycast gets traffic into the platform
2. **POP-level health and policy**
   - unhealthy or overloaded paths are suppressed
3. **Regional routing**
   - application or edge logic picks the best serving region
4. **Affinity exceptions**
   - sticky sessions, data residency, or home-region ownership constrain routing

## Data Model & APIs

Core entities:

- `SteeringPolicy`
- `RegionHealth`
- `CapacitySignal`
- `AffinityRule`
- `OverrideWindow`

Useful APIs:

- `ChooseRoute(client, policy, signals)`
- `ApplyOverride(scope, target, ttl)`
- `ExplainRoute(request_id)`
- `ValidatePolicy(policy)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| health signal flapping causes route churn | rapid route changes and latency oscillation | add hysteresis and confidence windows |
| nearest POP sends traffic to overloaded region | rising regional saturation despite normal POP health | couple routing to downstream capacity signals |
| control-plane update lags during incident | policy version skew across POPs | use versioned policy rollout and emergency override path |
| affinity-bound traffic cannot reroute cleanly | error spikes only for sticky sessions | define degraded behavior for stateful or home-region sessions |

## Observability

- metric: route choice distribution by POP, region, and policy reason
- metric: policy propagation lag and version skew
- metric: post-route regional saturation and reject rate
- log: sampled route explanations with health, latency, and override inputs
- trace: ingress POP, selected region, and fallback path
- SLO: emergency routing overrides should propagate within an explicit target window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| anycast ingress | fast and resilient global entry | harder path observability and policy intuition | pure DNS-only global routing |
| latency-biased with overload override | better user experience without melting hot regions | more dynamic control logic | static geography maps |
| affinity exceptions | protects session and data locality | limits recovery freedom | always reroute everything anywhere |

## Interview It

**Google framing:** "How would you route global user traffic to the best region?" Expect questions on health, capacity signals, and stickiness.

**Cloudflare framing:** "How do anycast and edge routing work together during an incident?" Expect pressure on control-plane safety and route explanation.

**Follow-ups:**
1. What changes if the service uses anycast ingress but region-bound writes?
2. How do you avoid route flapping?
3. What if a region is healthy but nearly full?
4. How do compliance boundaries override latency?
5. What telemetry proves your routing policy is doing the right thing?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/interview-card-traffic-steering.md`
- `outputs/failure-checklist-traffic-steering.md`

## Exercises

1. **Easy** — Compare DNS steering and anycast ingress for a stateless API.
2. **Medium** — Add home-region affinity for user writes.
3. **Hard** — Redesign for a flash crowd where the lowest-latency region is already saturated.

## Further Reading

- [Cloudflare Learning Center - Anycast network](https://www.cloudflare.com/learning/cdn/glossary/anycast-network/) — good grounding for ingress behavior
- [Google Cloud - Global load balancing overview](https://cloud.google.com/load-balancing/docs/load-balancing-overview) — useful contrast between traffic-steering layers
