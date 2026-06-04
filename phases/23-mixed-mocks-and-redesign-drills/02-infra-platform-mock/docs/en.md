# Mixed Mock: Infra Platform

> Infra platform interviews reward candidates who can separate control plane from data plane without making either one magical.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Practice a mixed mock for infrastructure and platform prompts where safety, rollout, tenancy, and failure-domain reasoning matter as much as raw throughput.  
**Prerequisites:** `14-rate-limiters-ids-and-hashing/07-control-vs-data-plane-drill`, `18-messaging-and-job-platforms/07-messaging-platform-drill`, `21-google-senior-staff-system-design/06-staff-deep-dive`  
**Estimated time:** ~90 min  
**Primary artifact:** infra platform scorecard validator + design review prompts  

## The Problem

Run a full mock on an infrastructure or platform prompt such as:

- design a global configuration rollout platform
- design a multi-tenant job execution platform
- design a metrics ingestion and query service
- design a service discovery and traffic policy platform

These prompts are dangerous because many candidates speak in abstractions. Strong answers show concrete safety properties, explain blast radius, and distinguish high-QPS data paths from slower but riskier control paths.

## Clarify

- Who are the users: application teams, end users indirectly, operators, or customers?
- Which path is the high-frequency data plane and which path is the correctness-sensitive control plane?
- What failure hurts more: stale config, rejected traffic, lost telemetry, or unsafe rollout?
- How much tenant isolation and explainability is required?

## Requirements

### Functional

- Define the platform contract and primary users.
- Separate control-plane and data-plane responsibilities clearly.
- Produce rough sizing for both hot and slow paths.
- Explain safety guardrails, observability, and redesign.

### Non-functional

- Keep failure domains explicit.
- Avoid pretending rollouts or migrations are free.
- Show tenant isolation and blast-radius thinking.
- Tie abstractions to measurable operator outcomes.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Data-plane QPS | 100K to multi-million depending on prompt | determines placement, caching, and latency budget |
| Control-plane mutation rate | tens to thousands per second | shapes replication, rollout, and audit design |
| Tenant count | hundreds to tens of thousands | changes isolation, quotas, and supportability |
| Rollout latency target | seconds to minutes | determines how aggressively state must propagate |
| Failure-domain budget | one tenant, one shard, one region, or global | exposes where the real design risk lives |

## Architecture

Recommended structure:

```text
clarify platform contract
  -> size data plane and control plane separately
  -> define failure domains and tenancy boundaries
  -> propose architecture
  -> deep dive on one risky area:
       rollout safety
       metadata consistency
       multi-tenant isolation
       query versus ingest balance
  -> explain observability and incident handling
  -> redesign under a changed operational constraint
```

A strong answer usually names:

1. source of truth for policy or metadata
2. propagation path
3. local serving behavior when propagation lags
4. operator guardrails and audit trail

## Data Model & APIs

Useful answer skeleton:

```text
infra_mock_answer(
  platform_contract,
  tenants,
  control_plane,
  data_plane,
  rollout_model,
  failure_domains,
  observability,
  redesign
)
```

Scorecard fields:

- contract clarity
- control/data plane separation
- sizing quality
- safety and guardrail quality
- migration or rollout quality
- observability quality
- redesign quality

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| answer describes only data plane | rollout, metadata, or audit safety is missing | force an explicit control-plane section |
| failure domains stay vague | no blast-radius boundary is named | require one concrete isolation boundary |
| migration story is magical | state changes appear instantly everywhere | define rollout waves, versioning, and rollback |
| tenant risk is underexplored | noisy-neighbor or config mistakes are invisible | add quotas, policy validation, and scoped failure handling |

## Observability

- metric: propagation lag or rollout convergence by region and tenant
- metric: data-plane latency and error rate segmented by policy version or shard
- metric: quota, saturation, or hotspot indicators per tenant
- log: policy changes, reason codes, and rollback triggers
- trace: control-plane mutation through propagation into serving behavior
- SLO: define one operator-visible and one user-visible success signal

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit control/data plane split | clearer safety reasoning | more answer complexity | generic all-in-one service story |
| versioned rollout with rollback | safer propagation | slower convergence and more bookkeeping | instant global push |
| scoped tenant isolation | smaller blast radius | lower packing efficiency | shared everything by default |

## Interview It

**Google framing:** Expect pressure on scope, safety, and whether the abstractions map cleanly to concrete system behavior.

**Cloudflare framing:** Expect more focus on propagation lag, regional behavior, and how the platform behaves under partial control-plane failure.

**Suggested prompts:**
1. Design a global configuration rollout platform.
2. Design a service discovery and policy platform.
3. Design a metrics ingestion and query service.
4. Design a multi-tenant job execution platform.

## Ship It

- `outputs/design-review-prompts-infra-platform-mock.md`
- `outputs/skill-infra-platform-mock.md`

## Exercises

1. **Easy** — For one prompt, draw only the control plane and data plane.
2. **Medium** — Add rollout guardrails and one rollback trigger.
3. **Hard** — Redesign the system after the interviewer says one tenant must now receive policy updates in under five seconds.

## Further Reading

- [Google SRE](https://sre.google/books/) — strong reference for operational ownership and safe rollouts  
- [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md) — useful for thinking about desired state, versioning, and reconciliation  
