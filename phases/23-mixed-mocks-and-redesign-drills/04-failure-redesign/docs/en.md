# Redesign After a Major Failure Mode

> A mature redesign starts from the incident you actually had, not the architecture you wish you had built.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Practice redesigning a system after a major failure mode reveals that the original architecture violated an operational assumption.  
**Prerequisites:** `10-reliability-retries-and-backpressure/08-reliability-drill`, `11-observability-slos-and-debugging/08-observability-drill`, `22-cloudflare-edge-platform-design/04-traffic-steering`  
**Estimated time:** ~75 min  
**Primary artifact:** failure redesign playbook + observability prompts  

## The Problem

The interviewer gives you a healthy-looking baseline design, then introduces a painful incident:

- retries turned a partial slowdown into a global outage
- stale config propagated to every region
- cache invalidation lag caused correctness failures
- one hot shard took down an entire tenant tier
- failover promoted bad state and exposed stale reads

Your job is not to retell the incident. Your job is to show how the system should change after learning from it.

## Clarify

- What exactly failed first: dependency latency, stale state, overload, bad rollout, or operator action?
- Was the failure localized or did the design amplify it?
- Which signal should have detected the problem earlier?
- Which guarantee must improve after the redesign: containment, correctness, recovery speed, or operator safety?

## Requirements

### Functional

- Identify the hidden assumption the incident broke.
- Explain the failure chain, not only the final symptom.
- Propose a redesign that reduces recurrence or blast radius.
- Add better observability, guardrails, and rollout discipline.

### Non-functional

- Keep the redesign proportional to the real failure.
- Avoid magical "just add circuit breakers everywhere" answers.
- Show containment and recovery trade-offs clearly.
- Make the operator experience part of the design.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Failure amplification factor | 2x to 100x over normal load or error rate | reveals why the incident became severe |
| Detection delay | seconds to tens of minutes | changes how much user damage occurs |
| Recovery objective | fast fail-closed, fail-open, or partial degrade | shapes mitigation strategy |
| Blast radius | shard, tenant, region, or global | determines the redesign target |
| Guardrail cost | additional latency, lower efficiency, or rollout friction | makes the redesign realistic |

## Architecture

Use an incident-first loop:

```text
state the failure chain
  -> identify the broken assumption
  -> redesign the amplification point
  -> add containment and rollback
  -> improve detection and operator visibility
  -> explain what trade-off got worse
```

Common redesign moves:

1. bounded retries and admission control
2. staged rollout with validation gates
3. smaller failure domains and stronger isolation
4. stale-read or stale-config safety guards
5. safer failover promotion rules

## Data Model & APIs

Useful redesign template:

```text
failure_redesign(
  incident,
  failed_assumption,
  amplification_path,
  containment_change,
  detection_change,
  rollout_change
)
```

Key questions:

- what state or policy needs versioning or validation?
- what action should become impossible without a guardrail?
- what should operators be able to explain after the redesign?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| redesign fixes only the symptom | root assumption stays unstated | force one broken assumption and one amplification path |
| incident is treated as random bad luck | no mechanism explains spread | map the propagation path explicitly |
| observability stays generic | no earlier detection signal is named | add one leading indicator and one containment metric |
| redesign worsens steady-state too much | new latency or cost is hidden | name the operational tax honestly |

## Observability

- metric: the leading indicator that would have caught the failure earlier
- metric: blast-radius or isolation effectiveness after the redesign
- metric: recovery time and rollback success
- log: operator actions, rollout gates, and reason codes
- trace: failure propagation path through dependencies or regions
- SLO: improve containment or recovery for the specific incident type

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| stronger containment guardrails | smaller outages | more friction or latency | maximum throughput with weak brakes |
| staged rollout or promotion gates | safer state changes | slower convergence | fast unsafe propagation |
| better operator explainability | faster incident response | more instrumentation and metadata | opaque automated behavior |

## Interview It

**Google framing:** Expect the interviewer to ask what assumption failed and why your redesign is the smallest fix that changes the outcome.

**Cloudflare framing:** Expect more pressure on blast radius, regional propagation, and what the edge or traffic layer should do differently during the incident.

**Follow-ups:**
1. What signal should have paged first?
2. What part of the system amplified the failure?
3. What new guardrail would have blocked the incident?
4. What gets slower or more expensive after the redesign?
5. How would you roll the redesign out safely?

## Ship It

- `outputs/failure-redesign-playbook.md`
- `outputs/observability-prompts-failure-redesign.md`

## Exercises

1. **Easy** — Pick one incident type and write the failure chain.
2. **Medium** — Add one containment change and one detection metric.
3. **Hard** — Redesign the same incident for a stricter availability requirement where fail-closed is not acceptable.

## Further Reading

- [Google SRE](https://sre.google/books/) — strong foundation for incident thinking and postmortem-driven redesign  
- [Release It!](https://pragprog.com/titles/mnee2/release-it-second-edition/) — practical source for failure amplification and stability patterns  
