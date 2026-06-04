# Public API Safety Defaults

> Safe defaults are part of system design because misuse is a predictable workload, not a surprise edge case.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Design public and shared APIs with defensive defaults around limits, validation, auth scope, and expensive operations so misuse does not become an outage.  
**Prerequisites:** `04-apis-contracts-and-schema-evolution/02-idempotency-keys`, `04-apis-contracts-and-schema-evolution/03-pagination-and-filtering`  
**Estimated time:** ~60 min  
**Primary artifact:** safety checklist + abuse review prompt  

## The Problem

Public interfaces are not consumed only by ideal clients. They are used by buggy SDKs, aggressive retries, exploratory scripts, and sometimes malicious actors. A system that depends on perfect client behavior is not really production-ready.

This lesson focuses on contract-level safety defaults:

- bounded page sizes
- explicit auth scope
- timeouts and retry guidance
- idempotent write expectations
- rate limits and error semantics

## Clarify

- Is the API public to third parties, internal to many teams, or both?
- Which operations are most expensive or abuse-prone?
- What are the default limits on request size, page size, and retry behavior?
- What should fail closed versus degrade gracefully?

## Requirements

### Functional

- Expose a usable API without enabling runaway cost or abuse.
- Provide clear error responses and safe retry semantics.
- Bound expensive query and mutation shapes by default.

### Non-functional

- Protect backend capacity from accidental client misuse.
- Keep onboarding simple enough that good clients succeed quickly.
- Make safety limits observable and adjustable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Public QPS | 80K req/s | abuse and accidental misuse scale with adoption |
| Tenants | 50K active | fairness and per-tenant controls matter |
| Max page size | 100 items default | shapes backend fanout and payload cost |
| Retry burst | 10x during incidents | safe error codes and idempotency matter |
| Rough cost | rate limiting, auth checks, validation, observability | small per-request safety cost avoids much larger outages |

## Architecture

Safe API surface design usually includes:

- authentication and authorization before expensive work
- bounded limits on page size, body size, and time range
- explicit idempotency guidance for writes
- versioned error codes and retry semantics
- rate limits matched to tenant or token scope

The best answer here sounds operational:

"We do not let the default query shape become an unbounded scan, and we do not let the default retry behavior amplify a failure."

## Data Model & APIs

Examples of safe defaults:

- `limit` defaults to 50 and caps at 100
- write endpoints accept `Idempotency-Key`
- expensive export paths require async job creation
- auth scopes map to least-privilege operations
- error responses distinguish retryable from non-retryable failures

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| unlimited query shapes | high scan amplification and latency | cap parameters and validate ranges |
| ambiguous error semantics | clients retry unsafe operations blindly | document retryability and status code meaning |
| broad auth scopes | tenant data exposure or destructive misuse | least-privilege scopes and explicit scope checks |
| safety defaults too strict | good clients fail to onboard smoothly | measure limit rejections and tune with intent |

## Observability

- metric: rejected requests by safety policy type
- metric: per-tenant rate-limit hits
- metric: expensive query denial count
- log: auth scope, normalized query shape, and reject reason
- trace: validation and authorization path for slow or denied requests
- SLO: safety mechanisms prevent expensive misuse without degrading healthy clients beyond target levels

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| strict parameter caps | predictable backend load | some clients need extra negotiation | open-ended defaults |
| explicit retry semantics | safer client behavior | more documentation and contract discipline | vague 500-style guidance |
| least-privilege scopes | smaller blast radius | more auth design work | broad catch-all tokens |

## Interview It

**Google framing:** "Design a public API for a developer platform." The signal is whether you think about accidental misuse, not only happy-path ergonomics.

**Cloudflare framing:** "Design default safety controls for an internet-facing API product." The signal is whether you connect contract shape to abuse prevention and origin protection.

**Follow-ups:**
1. Which limits should be tenant-specific versus global?
2. How do you expose clearer retry semantics to SDKs?
3. What if large customers genuinely need bigger limits?
4. Which errors should be retryable, and how should clients know?
5. How do you keep safety defaults from becoming product friction?

## Ship It

- `outputs/api-safety-checklist.md`
- `outputs/abuse-review-api-safety-defaults.md`

## Exercises

1. **Easy** — Define safe defaults for a list and create API in a multi-tenant SaaS product.  
2. **Medium** — Redesign a query API that currently allows unbounded filters and body sizes.  
3. **Hard** — Explain the contract-level protections you would add before exposing an internal API publicly.  

## Further Reading

- [OWASP API Security Top 10](https://owasp.org/API-Security/) — useful for thinking about API misuse and safety defaults  
- [Google API Design Guide](https://cloud.google.com/apis/design) — solid conventions for safer public interfaces  
