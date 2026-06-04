# Traffic Steering Across POPs and Regions

> Traffic steering only looks elegant from a distance. Up close it is a constant negotiation between latency, health quality, cost, and how much local autonomy you trust.

**Type:** Learn
**Company focus:** Cloudflare
**Learning goal:** Design steering policies across POPs and regions that balance latency, origin health, isolation, and operational safety instead of treating routing as a static shortest-path problem.
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`, `22-cloudflare-edge-platform-design/01-global-api-edge-gateway`, `22-cloudflare-edge-platform-design/03-origin-protection`
**Estimated time:** ~75 min
**Primary artifact:** steering decision card + failure playbook

## The Problem

Design traffic steering for a global edge platform. Requests enter anycast POPs, then may be steered to local caches, regional shields, or remote origins depending on policy and health. The platform must stay fast, stable, and economically reasonable even during regional incidents or skewed traffic.

This prompt is about routing judgment, not just network trivia. Strong answers explain what signals drive steering, how much local POPs can override global policy, and how to avoid turning routing changes into cascading incidents.

## Clarify

- Is the main goal lowest user latency, strongest origin protection, or lowest transit cost?
- Which routing choices are made at BGP/anycast level versus application-level forwarding?
- How quickly can steering policy change during incidents?
- Do all tenants share the same steering policy, or are some premium routes treated differently?

## Requirements

### Functional

- Route users into suitable POPs.
- Steer requests from POPs to shields or regional origins based on policy and health.
- Support incident-time overrides and controlled rollback.
- Expose steering decisions for debugging and customer explanation.

### Non-functional

- Minimize routing oscillation.
- Avoid hidden dependence on stale control-plane data.
- Respect cost and capacity headroom during failover.
- Limit tenant blast radius from custom routing policies.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| POP count | 100s | steering state and rollout fanout matter |
| Routing update frequency | seconds to minutes | determines control-plane design |
| Regional traffic skew | 2x to 5x burst during incidents | drives reserve capacity |
| Premium route classes | small but sensitive | may need differentiated policy |
| Propagation lag tolerance | seconds | affects safety of fast steering changes |

## Architecture

```text
client
  -> anycast POP selection
  -> local steering policy
  -> shield / regional forwarding choice
  -> origin pool selected from health, latency, cost, and capacity signals
```

Design guidance:

1. Keep a clear separation between control-plane policy publication and data-plane local decisions.
2. Use local health and latency observations, but constrain them with bounded policy.
3. Prefer safe dampening over hyper-reactive routing.
4. Decide ahead of time which traffic classes may trade latency for isolation or cost.

## Data Model & APIs

Useful records:

```text
steering_policy(
  route_class,
  preferred_regions,
  max_failover_cost,
  health_weight,
  latency_weight,
  capacity_weight,
  dampening_seconds
)
```

Helpful APIs:

- `GetSteeringPolicy(route_class)`
- `SelectRegion(pop, route_class, signal_set)`
- `ExplainPath(request_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| route oscillation between regions | failover and failback churn spikes | dampening and cooldown windows |
| stale control-plane policy at some POPs | policy version skew across POPs | versioned rollout and lag alarms |
| lowest-latency path overloads | latency and saturation both rise | incorporate capacity into steering |
| premium tenant policy harms shared fleet | per-tenant routing anomalies | hard isolation and route-class budgets |
| local POP overreacts to noisy measurements | local divergence from regional pattern | smoothing and shared signal blending |

## Observability

- metric: selected region distribution by POP and route class
- metric: routing override frequency
- metric: steering policy version skew
- metric: path latency versus chosen path cost
- log: steering decisions with ranked candidate reasons
- trace: entry POP -> steering choice -> upstream region -> response

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| local plus global signals | faster incident response with guardrails | more tuning complexity | control-plane-only routing |
| dampened steering | more stability | slower reaction to sudden changes | instant reroute on every metric twitch |
| route classes with differentiated policy | better product fit | more operator complexity | one global routing policy for all traffic |

## Interview It

**Google framing:** "Design global traffic steering for a serving system." Strong answers discuss policy/data-plane separation, regional overload, and rollback safety.

**Cloudflare framing:** "How would you steer traffic across POPs and regions at the edge?" Strong answers cover anycast entry, application-level forwarding, explainability, and cost-aware failover.

**Follow-ups:**
1. What if the cheapest path becomes the noisiest path?
2. What if a POP loses fresh control-plane policy but still serves traffic?
3. What if a premium customer buys lower latency at higher cost?
4. What if a regional outage causes 5x traffic to a neighboring region?

## Ship It

- `outputs/interview-card-traffic-steering.md`
- `outputs/failure-playbook-traffic-steering.md`

## Exercises

1. **Easy** — Explain when latency should lose to capacity in a routing decision.
2. **Medium** — Add a dampening policy for a region that flaps every few minutes.
3. **Hard** — Redesign for a case where data sovereignty rules forbid some cross-region steering choices.

## Further Reading

- [Anycast](https://en.wikipedia.org/wiki/Anycast) — baseline routing concept
- [Cloudflare engineering blog](https://blog.cloudflare.com/) — practical routing and edge design context
