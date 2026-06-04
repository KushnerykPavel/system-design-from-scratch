# Global API Edge Gateway

> Edge systems are judged twice: once by latency, and again by how safely they fail when origins get weird.

**Type:** Build  
**Company focus:** Cloudflare  
**Learning goal:** Design an edge gateway that terminates TLS, enforces auth and rate limits, and protects regional origins while preserving low tail latency.  
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`, `10-reliability-retries-and-backpressure/03-circuit-breakers`, `12-security-abuse-and-multitenancy/01-auth-and-trust`  
**Estimated time:** ~90 min  
**Primary artifact:** design review prompt + observability checklist  

## The Problem

You need a global API edge gateway. Requests arrive at many POPs. The gateway must terminate TLS, validate auth, apply policy, and forward to regional origins. It should keep origin latency and blast radius low, degrade safely during failures, and expose enough signals for operators to understand regional incidents.

This is a Cloudflare-style prompt because the interesting parts are not only routing and auth. They are also origin shielding, cacheability, retries, POP/regional blast radius, and abuse resistance.

## Clarify

- Are responses cacheable, or is every request a dynamic origin call?
- Is the goal lowest latency, strongest consistency, strongest origin protection, or lowest cost?
- Can POPs forward to multiple regions, or is there a preferred home region?
- Should failures fail open for auth and rate limiting, or fail closed?

## Requirements

### Functional

- TLS termination at the edge.
- Authentication and authorization enforcement.
- Per-tenant or per-route policy checks.
- Regional origin forwarding with retries and failover controls.

### Non-functional

- Low p95/p99 latency.
- Limited blast radius during partial origin or regional failure.
- Strong observability by POP, region, tenant, and origin pool.
- Abuse and misconfiguration resistance.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak QPS | 2M global | drives POP sizing, auth hot path, and origin shielding |
| Request size | small JSON APIs, 2-20 KB typical | shapes egress and TLS overhead |
| Cache hit rate | low to moderate depending on endpoint | affects origin load and shielding value |
| Regional failover factor | 2x when one region is unhealthy | drives headroom requirements |
| Rough cost | edge compute + egress + origin capacity | keeps retries and forwarding strategy realistic |

## Architecture

```text
client
  -> nearest POP
     -> TLS termination
     -> auth / policy / rate limit
     -> optional cache
     -> origin shield / regional router
     -> regional origin pool
```

Key ideas:

1. Keep policy decisions close to the client.
2. Avoid retry amplification into unhealthy origins.
3. Separate **edge request handling** from **control-plane policy rollout**.
4. Preserve a clean signal path for POP-level and region-level debugging.

## Data Model & APIs

Gateway policy model:

```text
listener -> route -> auth policy -> rate policy -> origin pool -> retry budget
```

Helpful APIs:
- `Evaluate(request)`
- `SelectOrigin(route, region, health_state)`
- `Explain(request_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one regional origin gets slow, not fully down | latency and retry counters rise | outlier detection, retry budget, failover threshold |
| POP policy versions drift | policy version metrics differ by POP | versioned rollout and convergence alarms |
| aggressive retries amplify origin failure | retry counts increase while success drops | bounded retry budgets and circuit breaking |
| auth service dependency degrades | auth decision latency spikes | local token validation, cached metadata, fail-closed only for high-risk routes |
| abuse floods one tenant or route | per-tenant reject/allow anomalies | tiered limits, challenge paths, and targeted mitigation |

## Observability

- metric: request latency by POP, route, tenant, and origin pool
- metric: retry budget consumption
- metric: auth decision latency and error rate
- metric: origin shield hit ratio or forwarded request ratio
- metric: policy version skew across POPs
- log: sampled reject and failover decisions with reason codes
- trace: request correlation from edge decision to origin response

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| regional origin shielding | protects fragile origins and improves cache locality | adds one more forwarding layer | direct POP-to-origin traffic for every request |
| bounded retries with circuit break | reduces retry storms | can surface more errors briefly | unbounded retries that worsen incidents |
| local token validation when possible | lower latency and fewer dependency hops | policy freshness is harder | central auth check on every request |

## Interview It

**Google framing:** "Design a globally distributed API serving layer." Strong answers still cover retries, auth, and observability, but may focus less on POP semantics and more on serving architecture.

**Cloudflare framing:** "Design a global API gateway at the edge." Strong answers must discuss origin protection, POP/regional routing, layered policy enforcement, and how to stay incident-safe during partial origin slowness.

**Follow-ups:**
1. What changes if some endpoints become cacheable?
2. What if one region is healthy but expensive?
3. What if auth policy updates must propagate in under 5 seconds?
4. What if a misconfigured customer route creates retry storms?
5. What changes if requests are long-lived streaming connections?

## Ship It

- `outputs/design-review-global-api-edge-gateway.md`
- `outputs/observability-checklist-global-api-edge-gateway.md`
- `outputs/interview-card-global-api-edge-gateway.md`

## Exercises

1. **Easy** — Add one route type that must fail closed on auth errors.  
2. **Medium** — Add an origin shielding layer with a bounded cache for semi-cacheable endpoints.  
3. **Hard** — Redesign for WebSocket or SSE traffic where retries are much more dangerous.  

## Further Reading

- [Cloudflare engineering blog](https://blog.cloudflare.com/) — useful edge and routing design context  
- [Gateway API](https://gateway-api.sigs.k8s.io/) — helpful separation of listeners, routes, and policies  
