# Edge Compute and Data Gravity

> Running code at the edge is easy to pitch and hard to justify unless you know which state follows the request and which state must stay home.

**Type:** Learn
**Company focus:** Cloudflare
**Learning goal:** Decide what logic belongs at the edge, what data must remain regional or central, and how state placement changes latency, consistency, and cost.
**Prerequisites:** `05-storage-indexing-and-access-patterns/02-access-pattern-first`, `06-caching-and-invalidation/06-cache-layers`, `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`
**Estimated time:** ~60 min
**Primary artifact:** edge placement worksheet

## The Problem

Design an edge-enhanced system for tasks like request normalization, auth checks, personalization hints, image transforms, bot filtering, or lightweight API composition.

The trap is pushing too much stateful logic outward without asking where the source of truth lives, how cold-start or deploy fanout works, and whether remote data fetches erase the latency win.

## Clarify

- Is the edge logic stateless, cache-friendly, or stateful per user?
- What data must be fetched synchronously from a home region?
- Are writes allowed at the edge or only reads and transformations?
- Is the primary goal latency, origin offload, isolation, or extensibility?

If the interviewer does not answer, assume edge code is best for stateless or cache-backed request processing with a small amount of replicated reference data.

## Requirements

### Functional

- Run selected logic near the user.
- Keep authoritative state in a clearly owned region or backend.
- Fall back safely when edge dependencies are stale or unavailable.

### Non-functional

- Edge execution should reduce end-user latency or origin load materially.
- Deployment fanout and rollback should stay operationally safe.
- The design should avoid hidden cross-region data fetches on the hot path.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| POP execution volume | 500K req/s global | determines whether per-request remote fetches are affordable |
| Edge state footprint | tens of MB to low GB per POP | small replicated data is feasible, large mutable state is not |
| Code rollout frequency | dozens per day | deployment propagation and rollback matter |
| Remote fetch share | target under 5% of edge requests | otherwise edge compute loses much of its latency value |
| Rough cost | edge CPU + replicated data + control plane | compute at every POP is expensive if the logic is weakly justified |

## Architecture

```text
client
  -> edge worker / function
     -> local config + cache
     -> optional regional data lookup
     -> origin or service backend
```

Good edge candidates:

- request authn/authz prechecks with cached policy
- bot scoring and request normalization
- static or semi-static personalization hints
- image resizing or response shaping

Poor edge candidates:

- strongly consistent multi-row writes
- chatty backend composition across many origins
- large mutable per-user state

## Data Model & APIs

Key entities:

- `EdgePolicyBundle`
- `ReferenceDataset`
- `HomeRegionBinding`
- `ExecutionGuardrail`
- `FallbackMode`

Useful APIs:

- `ExecuteAtEdge(request)`
- `FetchReferenceData(key)`
- `ResolveHomeRegion(user)`
- `ExplainFallback(request_id)`

The key modeling question is what data is replicated outward, what is cached with TTL, and what remains strictly home-region.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| edge code still calls home region for most requests | high remote-fetch ratio and weak latency gains | move only truly local logic to edge or replicate the needed reference data |
| stale policy bundle at some POPs | version skew across POPs | versioned rollouts and expiry on old bundles |
| edge write path creates cross-region inconsistency | divergence or retry storms | keep writes centralized or constrain to idempotent append-safe patterns |
| rollout bug propagates globally | simultaneous error spike across many POPs | staged deployment and fast rollback gates |

## Observability

- metric: edge execution latency and remote-fetch ratio
- metric: per-POP policy/data version skew
- metric: offload rate to origin and fallback frequency
- log: sampled fallback reasons and remote lookup causes
- trace: time spent at edge versus regional backend hops
- SLO: edge path should reduce latency enough to justify added operational surface

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| edge for stateless transforms | low latency and origin offload | control-plane fanout and distributed deploy risk | central-only request processing |
| replicate small reference data | fewer remote hops | memory and version-management cost | synchronous regional fetch on every request |
| centralize authoritative writes | simpler correctness | some latency remains centralized | fully stateful edge writes |

## Interview It

**Google framing:** "What logic would you move to the edge and what would you keep central?" Expect questions on state placement and hidden remote calls.

**Cloudflare framing:** "Design an edge-compute path for a globally used product." Expect pressure on rollout safety, config distribution, and data gravity.

**Follow-ups:**
1. What data is worth replicating to every POP?
2. When does edge compute stop being a latency win?
3. How do you roll back a bad global edge deploy quickly?
4. What if compliance requires some users to stay in-region?
5. Which write patterns, if any, are safe at the edge?

## Ship It

- `outputs/skill-edge-placement.md`
- `outputs/interview-card-edge-compute.md`

## Exercises

1. **Easy** — Place auth prechecks, image resizing, and analytics beacons on an edge or central map.
2. **Medium** — Redesign an edge personalization system when 20% of requests still need fresh user state.
3. **Hard** — Support regional privacy boundaries while keeping a low-latency edge path.

## Further Reading

- [Cloudflare Workers - How Workers works](https://developers.cloudflare.com/workers/reference/how-workers-works/) — strong intuition for edge execution constraints
- [Martin Kleppmann - Turning the database inside out](https://martin.kleppmann.com/2015/11/05/database-inside-out-at-strange-loop.html) — useful for thinking about moving computation toward data and vice versa
