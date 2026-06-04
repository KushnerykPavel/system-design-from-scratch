# Runbooks and First-Response Workflow

> The first ten minutes of an incident should not depend on memory or heroics.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Build runbooks that turn an alert into a reliable first-response workflow with diagnosis, mitigation, escalation, and communication steps.
**Prerequisites:** `03-design-framework-and-timing/07-interviewer-moves`, `10-reliability-retries-and-backpressure/07-bulkheads`, `11-observability-slos-and-debugging/05-alert-design`
**Estimated time:** ~60 min
**Primary artifact:** first-response runbook template + interview card

## The Problem

Teams often treat runbooks as documentation afterthoughts. During incidents that causes:

- responders hunting for tribal knowledge
- mitigation steps applied inconsistently
- escalations happening late or chaotically
- repeated incidents with the same confusion each time

Senior answers include the operator workflow, not just the system architecture.

## Clarify

- Is the runbook for one recurring symptom or a broad service family?
- What actions are safe for an on-call generalist versus domain expert only?
- Which mitigations are reversible?
- What communication obligations exist during customer-facing incidents?

If the interviewer is vague, assume a page on a user-facing service where the first responder may not be the service’s top expert.

## Requirements

### Functional

- Guide the responder from alert receipt to triage, mitigation, and escalation.
- Link telemetry, dashboards, and commands to concrete decision points.
- Capture communication and handoff expectations.

### Non-functional

- Keep the runbook short enough to use under stress.
- Make unsafe or irreversible actions explicit.
- Keep the workflow current as architecture changes.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| On-call expertise variance | new responder to staff engineer | runbook must help mixed experience levels |
| First-response target | first mitigation path within 10-15 minutes | encourages concise, decision-driven content |
| Major incident frequency | low, but high cost when it happens | maintenance effort is justified |
| Safe mitigation count | 3-5 standard actions | too many options slows responders |
| Rough cost | stale docs increase MTTR and risk | runbook quality has real operational payoff |

## Architecture

A practical runbook has these sections:

1. trigger and severity
2. immediate user impact check
3. first dashboards and queries
4. safe mitigations
5. escalation criteria
6. communication steps
7. recovery validation

```text
alert
  -> assess impact
  -> classify incident
  -> inspect known signals
  -> apply safe mitigation
  -> escalate if criteria met
  -> validate recovery and communicate
```

## Data Model & APIs

Useful runbook fields:

- `alert_name`
- `severity`
- `owner_team`
- `impact_check`
- `safe_actions`
- `dangerous_actions`
- `escalation_conditions`
- `status_template`

Useful interfaces:

- `LookupRunbook(alertName)`
- `StartIncident(severity, owner)`
- `RecordMitigation(action, actor, timestamp)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| runbook too long to use under stress | responders skip it and jump to ad hoc debugging | keep first-response path concise and front-loaded |
| mitigation is undocumented or ambiguous | different responders take conflicting actions | explicitly list reversible first steps and prerequisites |
| stale dashboard links or commands | runbook fails during incident | review runbooks after service changes and after incidents |
| no escalation trigger | responder spends too long alone in a severe outage | define clear time and severity thresholds for help |

## Observability

- metric: incidents with runbook usage recorded versus not
- metric: median time to first mitigation and to escalation
- metric: stale runbook age and broken-link checks
- log: incident timeline with mitigation steps taken
- trace: optional for validating recovery of the critical path after action
- SLO: strong runbooks support lower MTTR and safer response, even if not a direct user SLI

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| concise first-response runbook | usable under pressure | less encyclopedic detail | giant operations manual for live use |
| explicit safe vs dangerous actions | reduces accidental blast radius | extra maintenance discipline | undocumented judgment calls |
| embedded communication templates | faster stakeholder updates | some duplication with incident tooling | inventing updates ad hoc during outage |

## Interview It

**Google framing:** "What operational workflow would support this service when it breaks?" The signal is whether you think beyond dashboards into mitigation and escalation.

**Cloudflare framing:** "How would you support first response for a global edge incident affecting many POPs?" The signal is whether your runbook handles scope assessment, rollback, and cross-team coordination quickly.

**Follow-ups:**
1. What belongs in the runbook versus the dashboard?
2. Which actions should be reserved for experts only?
3. How do you keep runbooks current?
4. When do you escalate immediately?
5. What if the documented mitigation makes things worse?

## Ship It

- `outputs/first-response-runbook-template.md`
- `outputs/interview-card-runbooks.md`

## Exercises

1. **Easy** — Write the first five steps for a runbook responding to elevated API error rate.
2. **Medium** — Add escalation conditions for a regional dependency outage.
3. **Hard** — Redesign the runbook when the service has multiple degraded modes and only some are safe to mitigate automatically.

## Further Reading

- [Emergency Response](https://sre.google/sre-book/emergency-response/) — strong practical incident-response framing
- [Effective Troubleshooting](https://sre.google/sre-book/effective-troubleshooting/) — helpful for converting diagnosis into runbook steps
