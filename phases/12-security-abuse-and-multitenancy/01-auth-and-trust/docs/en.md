# Authentication, Authorization, and Trust Boundaries

> A secure design starts by saying who is trusted, where that trust ends, and what must be re-checked on every hop.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Distinguish authentication from authorization, draw trust boundaries explicitly, and choose where identities are verified, propagated, and constrained in a distributed system.
**Prerequisites:** `01-clarification-and-scope/04-assumption-logging`, `04-apis-contracts-and-schema-evolution/07-api-safety-defaults`, `10-reliability-retries-and-backpressure/07-bulkheads`
**Estimated time:** ~75 min
**Primary artifact:** trust-boundary checklist + interview card

## The Problem

Design the auth model for a multi-tenant API platform with edge gateways, internal services, background workers, and an admin control plane.

Weak interview answers say "use OAuth" and move on. Strong answers clarify:

- who authenticates users, services, and operators
- where tokens are validated versus merely forwarded
- what each service is allowed to trust from upstream
- how tenant identity and role information cross service boundaries

## Clarify

- Are users calling directly with browser sessions, API keys, OAuth tokens, or service credentials?
- Which actions are tenant-scoped, which are global admin actions, and which require break-glass access?
- Do internal services trust gateway-added identity headers, or do they re-verify signed tokens?
- Is the bigger risk impersonation, confused deputy behavior, lateral movement, or over-broad admin access?

If the interviewer is vague, assume an internet-facing API gateway, short-lived user or service tokens, tenant-scoped RBAC, and explicit re-validation at critical boundaries.

## Requirements

### Functional

- Authenticate end users, service-to-service traffic, and operators separately.
- Authorize requests by tenant, role, action, and resource scope.
- Propagate identity context safely across synchronous and asynchronous boundaries.

### Non-functional

- Keep request-path auth checks low-latency.
- Limit blast radius if one service is compromised.
- Make privilege decisions explainable during audits and incidents.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak API traffic | 250K req/s | drives token validation placement and cache strategy |
| Internal fanout | 3 to 5 service hops | each hop is a trust decision, not just a network call |
| Tenants | 200K active | affects policy-cardinality and role modeling |
| Admin actions | low volume, high impact | deserves stronger checks and audit logging |
| Rough cost | auth cache, policy store, and audit pipeline | security choices change latency and operating cost |

## Architecture

```text
user or service
  -> edge gateway
     -> identity validation
     -> policy lookup / cached authz
     -> internal services
        -> scoped identity propagation
        -> resource check
     -> audit log
```

Recommended shape:

1. Validate external identity at the first internet-facing boundary.
2. Convert identity into a small signed internal context or short-lived token.
3. Re-check authorization at the service that owns the resource.
4. Treat async jobs as a fresh boundary and carry actor plus tenant context explicitly.

## Data Model & APIs

Core entities:

- `Principal`
- `Credential`
- `RoleBinding`
- `Permission`
- `TenantScope`
- `AuditEvent`

Useful APIs:

- `Authenticate(credential)`
- `Authorize(principal, action, resource, tenant)`
- `MintInternalIdentity(principal, scopes, ttl)`
- `ExplainDecision(requestID)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| gateway injects identity that downstream trusts blindly | privileged actions succeed with malformed or forged context | sign internal identity and re-verify on critical hops |
| service performs authn but not resource-level authz | cross-tenant reads appear in audit logs | enforce ownership checks in resource-owning service |
| admin roles are too broad | audit trail shows wide read or write access | split break-glass, support, and routine admin roles |
| async worker loses actor context | deletions or writes become unattributable | propagate actor, tenant, and reason fields in job payloads |

## Observability

- metric: auth success and failure by credential type and boundary
- metric: authorization deny rate by action and tenant tier
- metric: token verification cache hit ratio and latency
- log: structured auth decision with principal type, tenant, action, and policy version
- trace: identity propagation and authz check spans on the request path
- SLO: privileged actions should have verifiable audit coverage and bounded decision latency

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| re-authorize at resource owner | stronger tenant isolation | more repeated checks and policy distribution | trust gateway only |
| short-lived internal identity | limits replay and lateral movement | token minting and clock-skew complexity | long-lived shared service credentials |
| explicit support/admin split | smaller blast radius | more operational friction | one broad operator role |

## Interview It

**Google framing:** "Design auth for an internal platform used by many services." Expect pushback on identity propagation, service credentials, and confused deputy risks.

**Cloudflare framing:** "Design auth and authorization for a global API edge product." Expect emphasis on low-latency validation, tenant safety, and how edge and origin trust differ.

**Follow-ups:**
1. When is gateway-only authorization acceptable, and when is it dangerous?
2. What changes for service-to-service calls versus end-user calls?
3. How do you model break-glass access safely?
4. What identity context crosses an async job boundary?
5. What changes at 10x tenant count?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/trust-boundary-checklist.md`
- `outputs/interview-card-auth-and-trust.md`

## Exercises

1. **Easy** — Draw the trust boundaries for an API gateway, worker, and admin console flow.
2. **Medium** — Redesign the system so downstream services can authorize without calling a central policy service on every request.
3. **Hard** — Add cross-region failover while preserving tenant-scoped authz and auditable break-glass actions.

## Further Reading

- [Google Zanzibar paper](https://research.google/pubs/zanzibar-googles-consistent-global-authorization-system/) — strong mental model for authorization at scale
- [Cloudflare Access](https://www.cloudflare.com/zero-trust/products/access/) — practical example of identity at the edge
