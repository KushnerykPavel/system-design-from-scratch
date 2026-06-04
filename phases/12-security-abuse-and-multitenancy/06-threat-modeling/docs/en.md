# Threat Modeling for Interview Design

> Threat modeling in an interview is less about naming every attack and more about proving you can prioritize realistic risks and change the design because of them.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Turn threat modeling into a fast, high-signal interview habit by mapping assets, actors, trust boundaries, abuse paths, and concrete design changes.
**Prerequisites:** `12-security-abuse-and-multitenancy/01-auth-and-trust`, `12-security-abuse-and-multitenancy/03-abuse-prevention`, `12-security-abuse-and-multitenancy/04-tenant-isolation`
**Estimated time:** ~60 min
**Primary artifact:** threat-model worksheet + interview card

## The Problem

Take any common system design prompt, such as "design a file sharing service" or "design a billing API," and produce a threat model that improves the answer instead of derailing it.

Strong interview threat modeling is short and focused:

- identify the most valuable assets
- identify realistic attacker or abuse paths
- connect those paths to trust boundaries
- name 2 to 4 design changes that materially reduce risk

## Clarify

- What assets matter most: user data, money, compute capacity, secrets, or admin actions?
- Which adversaries are realistic: external attackers, abusive customers, insiders, or buggy internal services?
- Is the system internet-facing, internally exposed, or both?
- Which failure hurts most: data leak, account takeover, cost explosion, or control-plane abuse?

If the interviewer is vague, choose the highest-value assets and the most likely threat actors instead of trying to model everything.

## Requirements

### Functional

- Enumerate assets, actors, trust boundaries, and top threats.
- Map threats to mitigations that change the architecture meaningfully.
- Keep the model concise enough for a live interview.

### Non-functional

- Prioritize realistic risk over exhaustive taxonomy.
- Preserve answer momentum instead of turning the interview into a security lecture.
- Make mitigations observable and testable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Top assets | 2 to 4 per prompt | keeps the model focused |
| Threat actors | 2 to 3 primary classes | avoids exhaustive but low-signal lists |
| Time budget | 3 to 5 interview minutes | the method must be compact |
| Mitigations | 2 to 4 architecture changes | security must change the design materially |
| Rough cost | controls, audit, and operational review | not every mitigation is worth it |

## Architecture

Fast interview flow:

1. Name the key assets.
2. Mark trust boundaries and data-entry points.
3. Pick the most likely and most damaging attacker paths.
4. Propose concrete mitigations tied to those paths.

```text
assets
  -> entry points
  -> trust boundaries
  -> attacker paths
  -> mitigations
  -> observability and response
```

## Data Model & APIs

Core entities:

- `Asset`
- `ThreatActor`
- `Boundary`
- `Threat`
- `Mitigation`

Useful APIs:

- `ScoreThreat(likelihood, impact)`
- `PrioritizeThreats(threats)`
- `MapMitigation(threat, control)`
- `ReviewCoverage(model)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| model becomes laundry list with no priorities | many threats named but no design change follows | limit to top assets and top attacker paths |
| only external attackers are considered | insider or internal-service risks absent | include at least one internal misuse path where relevant |
| mitigations are generic buzzwords | architecture stays unchanged | require each mitigation to alter boundary, identity, or policy design |
| no observability for security controls | abuse or auth drift goes unnoticed | attach metrics, logs, and audit trails to controls |

## Observability

- metric: auth failures, privilege-deny rates, and abuse-decision rates for key boundaries
- metric: control-specific signals such as challenge success, secret rotation lag, or tenant budget burn
- log: privileged actions, delete requests, and policy changes
- trace: key decision points where identity or policy alters request flow
- SLO: critical security controls should be auditable and not create silent user harm during normal operation

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| prioritize top threats only | keeps interview answer sharp | may leave edge cases unstated | exhaustive STRIDE walkthrough live |
| tie threats to architecture changes | high signal and actionable | forces trade-off choices | naming threats with no system impact |
| include operator and insider risk | more realistic | adds some control-plane complexity | only modeling anonymous internet attackers |

## Interview It

**Google framing:** "What are the main risks in this design and how would you mitigate them?" The signal is prioritization and design impact, not taxonomy recitation.

**Cloudflare framing:** "How would you think about abuse and trust boundaries for this edge service?" The signal is realistic attacker paths and edge-aware mitigations.

**Follow-ups:**
1. Which two threats are worth mentioning if time is short?
2. How did the threat model change your architecture?
3. Which risks do you accept rather than solve immediately?
4. What telemetry tells you the controls are working?
5. What changes at 10x traffic or tenant count?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/threat-model-worksheet.md`
- `outputs/interview-card-threat-modeling.md`

## Exercises

1. **Easy** — Threat-model a URL shortener in three minutes.
2. **Medium** — Add internal-service misuse and support-tool abuse to a multitenant API design.
3. **Hard** — Threat-model a billing platform where money movement and support actions both matter.

## Further Reading

- [OWASP Threat Modeling](https://owasp.org/www-community/Threat_Modeling) — compact baseline terminology
- [Google SRE - addressing cascading failures](https://sre.google/workbook/addressing-cascading-failures/) — useful for connecting security and resilience thinking
