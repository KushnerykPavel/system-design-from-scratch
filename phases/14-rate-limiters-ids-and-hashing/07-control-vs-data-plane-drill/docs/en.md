# Control-Plane vs Data-Plane Trade-offs Drill

> Mature system-design answers separate what must be fast from what must be correct eventually, then make the seam between them explicit.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Practice a full interview loop where policies, topology, and rollout intent come from a control plane, while request decisions stay in a resilient low-latency data plane.  
**Prerequisites:** `14-rate-limiters-ids-and-hashing/02-distributed-rate-limiter`, `14-rate-limiters-ids-and-hashing/05-service-discovery-placement`, `13-multi-region-cdn-and-edge-traffic/07-global-traffic-drill`  
**Estimated time:** ~60 min  
**Primary artifact:** drill worksheet + interview card  

## The Problem

Design a platform component such as a rate-limiting service, traffic router, or edge policy engine where rules are authored centrally but enforced on a high-QPS request path. The data plane must stay fast and available even when the control plane is stale, slow, or temporarily unavailable.

This drill ties the whole phase together:

- limiter primitives
- distributed ownership
- placement
- hotspot handling
- control-plane versus data-plane contracts

## Clarify

- What information must be live on the request path, and what can be cached?
- How stale can policies become before correctness breaks?
- Should the data plane fail open, fail closed, or degrade by class?
- What actions require central coordination versus local autonomy?

If the interviewer stays vague, assume edge enforcement at very high QPS with policy pushes every few minutes and stricter correctness for abuse controls than for optimization hints.

## Requirements

### Functional

- Accept centralized policy definitions and distribute them to the fleet.
- Make request decisions locally with bounded dependence on the control plane.
- Support rollout, rollback, and explainability of policy behavior.

### Non-functional

- Keep data-plane latency low and predictable.
- Bound control-plane blast radius during incidents.
- Surface policy drift quickly enough for safe operations.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Data-plane traffic | 1.5M req/s | request path cannot wait on central decisions |
| Policy objects | 200K active rules | distribution and cache footprint matter |
| Policy updates | 10K/min during incidents or launches | rollout safety becomes a real system concern |
| Fleet size | 1K edge or service nodes | version skew and staged rollout matter |
| Rough cost | replicated control plane + local caches | correctness and speed live on different budgets |

## Architecture

```text
operator
  -> control plane
     -> validation + policy compiler
     -> distribution channel
        -> local node cache
           -> data-plane evaluator
              -> allow / route / reject
```

Answer structure to practice:

1. State what belongs in control plane versus data plane.
2. Size rule count, update frequency, and fleet scale.
3. Choose caching and versioning behavior for local evaluators.
4. Define degraded mode when policy distribution is unhealthy.
5. Explain observability and rollback.

## Data Model & APIs

Useful entities:

- `PolicyRule`
- `CompiledPolicyBundle`
- `NodeVersion`
- `DriftAlert`
- `DegradedMode`

Useful APIs:

- `PublishPolicy(bundle)`
- `ValidateBundle(bundle)`
- `FetchOrWatch(bundle_version)`
- `Evaluate(request, bundle_version)`
- `ExplainDecision(request_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| control plane is down but data plane still has valid cache | publish failures but low request impact | keep local bundles cached and time-bounded |
| bad policy bundle propagates globally | reject spike after rollout | staged rollout, canary POPs, and fast rollback |
| version skew creates inconsistent customer experience | node-version divergence metrics | enforce expiry windows and bundle health checks |
| data plane reaches back to control plane on the hot path | latency and dependency blow up | compile and cache policy locally |

## Observability

- metric: bundle version skew across nodes
- metric: policy rollout latency and rollback latency
- metric: decision latency and reject rate by bundle version
- metric: degraded-mode activations and duration
- log: sampled policy decision explanation including rule ID and bundle version
- alert: control-plane publish success with no corresponding fleet adoption

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| local compiled bundle in data plane | low latency and resilient decisions | bounded staleness | live control-plane RPC per request |
| staged rollout of policy | reduces global blast radius | slower full adoption | instant worldwide rollout |
| different degraded modes by policy class | better balance of safety and availability | more policy semantics to define | one blanket fail-open or fail-closed mode |

## Interview It

**Google framing:** "Design a policy platform used by many services." Expect questions about rollout safety, ownership, and how to keep central config from becoming a hidden dependency.

**Cloudflare framing:** "Design a global edge policy system." Expect pressure on POP-local decisions, rule propagation, and stale bundle handling.

**Follow-ups:**
1. Which policies may safely fail open?
2. What if bundle propagation must finish globally within 20 seconds?
3. How do you explain inconsistent behavior across POPs?
4. What changes when policies depend on live usage counters?
5. Where do you put validation to keep bad rules out of the fleet?

## Ship It

- `outputs/design-review-control-vs-data-plane-drill.md`
- `outputs/interview-card-control-vs-data-plane-drill.md`

## Exercises

1. **Easy** — Classify a simple rate-limit rule into control-plane and data-plane pieces.
2. **Medium** — Design rollout and rollback for a global policy bundle.
3. **Hard** — Support live counters, stale bundle fallback, and customer-visible explainability together.

## Further Reading

- [Envoy architecture overview](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/intro/arch_overview) — good reference for control-plane/data-plane separation  
- [Google SRE workbook](https://sre.google/workbook/) — practical operational framing for safe rollout and blast-radius control  
