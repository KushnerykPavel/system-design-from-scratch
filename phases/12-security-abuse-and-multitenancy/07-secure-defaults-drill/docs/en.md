# Secure Defaults Drill

> The best security answers do not bolt controls on later. They choose defaults that make the safe path the easy path.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Synthesize Phase 12 by reviewing a system design for safe-by-default choices across auth, secrets, abuse controls, tenant isolation, and deletion semantics.
**Prerequisites:** `12-security-abuse-and-multitenancy/06-threat-modeling`
**Estimated time:** ~60 min
**Primary artifact:** secure-defaults rubric + interview card

## The Problem

Use this prompt:

"Design the control plane and request path for a multitenant developer platform that exposes public APIs, background jobs, admin tooling, customer secrets, and per-tenant usage limits."

This drill is strong because it forces you to answer:

- what is secure by default for new tenants and new engineers
- what is denied until explicitly allowed
- which controls belong on the request path versus the control plane
- how to fail safely without making the product unusable

## Clarify

- Which actions are public, authenticated, tenant-scoped, or operator-only?
- Which defaults should prioritize product adoption versus safety?
- What does a brand-new tenant get automatically?
- Which degraded modes should fail open versus fail closed?

If the interviewer is vague, assume a public API, internal admin plane, shared multitenant workers, and a need for secure defaults without hand-managed onboarding.

## Requirements

### Functional

- Define secure defaults for identity, privilege, secret handling, quotas, and deletion.
- Show how operators override defaults safely when needed.
- Explain degraded-mode behavior for critical security controls.

### Non-functional

- Keep the default path simple enough that teams do not bypass it.
- Minimize blast radius for new services or tenant misconfiguration.
- Make security posture inspectable and reviewable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| New tenants | 5K per month | defaults must scale without manual review |
| Public API QPS | 150K req/s | request-path defaults must stay cheap |
| Internal services | 200+ | secure service onboarding must be repeatable |
| Operator actions | low volume, high consequence | override paths need strong guardrails |
| Rough cost | secure defaults reduce incident cost but add setup friction | default choices are product and ops trade-offs |

## Architecture

A strong drill answer usually includes:

1. deny-by-default authz for new capabilities
2. short-lived credentials and managed secret delivery
3. default quotas, abuse guards, and audit logging
4. tenant-scoped storage and worker budgets
5. deletion and retention policies defined up front

```text
tenant or operator
  -> authenticated entry point
  -> policy-enforced request path
  -> tenant-scoped services and data
  -> audit / quotas / deletion workflows
```

## Data Model & APIs

Core entities:

- `DefaultPolicy`
- `TenantTier`
- `OverrideRequest`
- `AuditRecord`
- `DeletionPolicy`

Useful APIs:

- `ProvisionTenantDefaults(tier)`
- `ReviewOverride(change)`
- `ValidateSecureDefaults(plan)`
- `AssessDrill(answer)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| product teams bypass secure defaults because onboarding is too hard | shadow credentials or ad hoc exceptions appear | make paved path fast and audited |
| new tenants get over-broad permissions | privilege review finds excess default grants | deny-by-default and explicit role templates |
| emergency override remains permanent | override age and review backlog grow | require expiry and post-incident review |
| deletion, authz, and quotas are defined inconsistently across services | policy drift metrics or audit gaps appear | central templates plus service-level enforcement |

## Observability

- metric: default-policy adoption, override count, and override expiry age
- metric: deny rates, quota trips, and secret-refresh failures for newly onboarded services
- log: privileged overrides, tenant provisioning, and policy template changes
- trace: request-path decisions that show which default policy applied
- SLO: new tenants and services should onboard securely without requiring manual exception churn

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| deny-by-default | smaller initial blast radius | more explicit onboarding work | permissive defaults with later cleanup |
| paved secure onboarding | fewer ad hoc exceptions | platform investment | leaving every team to solve security alone |
| expiring overrides | limits permanent drift | operational review burden | indefinite emergency exceptions |

## Interview It

**Google framing:** "What defaults would you choose for a new internal platform or API product?" Expect emphasis on safe service onboarding and how defaults reduce long-term toil.

**Cloudflare framing:** "How would you design secure defaults for a multitenant edge or developer platform?" Expect focus on tenant safety, abuse controls, and operational override discipline.

**Follow-ups:**
1. Which default matters most for brand-new tenants?
2. What do you deny until explicitly allowed?
3. Which failures should fail closed and which should fail open?
4. How do you keep overrides from becoming permanent security debt?
5. What changes at 10x tenant and service count?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/secure-defaults-rubric.md`
- `outputs/interview-card-secure-defaults-drill.md`

## Exercises

1. **Easy** — Review a simple API design and name three unsafe defaults.
2. **Medium** — Add a secure onboarding path for internal services using workload identity and managed secrets.
3. **Hard** — Redesign the drill for enterprise tenants that require stricter defaults without harming self-serve adoption too much.

## Further Reading

- [Secure by design principles](https://www.cisa.gov/securebydesign) — good framing for choosing defaults that prevent entire classes of mistakes
- [Google BeyondProd](https://cloud.google.com/beyondprod) — useful mental model for identity and trust in modern production systems
