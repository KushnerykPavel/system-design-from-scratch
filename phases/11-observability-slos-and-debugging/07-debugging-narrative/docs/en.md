# Incident Debugging Narrative

> A senior engineer debugs by shrinking the hypothesis space in public.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Explain incident debugging as a disciplined narrative of hypotheses, evidence, narrowing moves, and mitigation choices instead of random dashboard hopping.
**Prerequisites:** `03-design-framework-and-timing/04-deep-dive-selection`, `10-reliability-retries-and-backpressure/08-reliability-drill`, `11-observability-slos-and-debugging/06-runbooks`
**Estimated time:** ~60 min
**Primary artifact:** debugging narrative worksheet + interview card

## The Problem

In interviews, candidates often say:

- "I’d check the logs."
- "I’d look at the dashboard."
- "I’d scale the service."

Those are actions, not a debugging story. A strong debugging narrative says:

1. what symptom was observed
2. what classes of causes are plausible
3. what evidence separates them
4. what mitigation is safe before certainty
5. how the fix will be validated

## Clarify

- What is the user-visible symptom: errors, slowness, backlog growth, or data inconsistency?
- Did the incident start suddenly or degrade over time?
- Was there a recent deploy, config push, traffic shift, or dependency event?
- Which mitigations are safe before root cause is proven?

If the interviewer is vague, assume a latency or error spike on a service with at least one dependency and one recent change candidate.

## Requirements

### Functional

- Build a hypothesis-driven incident narrative from symptom to mitigation.
- Use telemetry to narrow root-cause classes quickly.
- Explain how to validate recovery and capture follow-up work.

### Non-functional

- Keep the narrative understandable to interviewers and teammates.
- Avoid overconfident root-cause claims too early.
- Prefer reversible mitigations when certainty is low.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Possible root-cause classes | deploy, dependency, capacity, data skew, abuse | debugging quality depends on narrowing quickly |
| Triage window | 5-15 minutes for first directional answer | forces deliberate prioritization |
| Change frequency | multiple deploys or config pushes per day | recent-change correlation matters |
| Data sources | metrics, traces, logs, runbooks, deploy timeline | narrative should connect them coherently |
| Rough cost | bad debugging increases MTTR and blast radius | reasoning quality is part of reliability |

## Architecture

A useful debugging flow:

```text
symptom
  -> classify user harm
  -> list plausible cause buckets
  -> choose highest-signal evidence
  -> narrow to one or two likely causes
  -> apply safest mitigation
  -> validate recovery
  -> record follow-up and prevention work
```

Good evidence order often looks like:

- user-impact and SLO panels
- scope by region / route / tenant class
- recent deploy or config timeline
- dependency comparisons
- traces or logs on the abnormal path

## Data Model & APIs

Useful incident entities:

- `SymptomSnapshot`
- `Hypothesis`
- `Evidence`
- `Mitigation`
- `RecoveryCheck`

Useful fields:

- `scope`
- `suspected_component`
- `confidence`
- `reversible`
- `validation_metric`

Useful interfaces:

- `ListRecentChanges(window)`
- `CompareRegions(metric)`
- `CaptureHypothesis(note, confidence)`
- `ValidateMitigation(metric, expectedDirection)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| random walk debugging | responders jump tools without narrowing causes | require explicit hypothesis and next discriminating check |
| early fix bias | first plausible cause is treated as proven | compare evidence against at least one competing hypothesis |
| unsafe mitigation under uncertainty | outage worsens after manual action | prefer reversible actions and clear rollback plan |
| incident closes without learning | same failure repeats later | record narrative, evidence, and prevention items in postmortem |

## Observability

- metric: time to first plausible hypothesis and to first mitigation
- metric: percentage of incidents with deploy or dependency correlation identified
- log: incident timeline entries with hypothesis and evidence notes
- trace: sampled slow/error paths to confirm bottleneck location
- SLO: debugging quality is indirectly reflected in MTTR and repeat-incident rate

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| hypothesis-first workflow | faster narrowing and better communication | requires discipline under stress | tool-hopping without structure |
| reversible mitigation preference | safer under uncertainty | may recover more slowly than drastic action | risky irreversible fix attempt |
| public narrative during response | better team alignment | slight overhead to document while debugging | silent solo debugging |

## Interview It

**Google framing:** "A service latency SLO is burning quickly. Walk me through how you would debug it." The signal is whether your reasoning is structured and evidence-driven.

**Cloudflare framing:** "Users report intermittent failures from some regions. How would you investigate?" The signal is whether you narrow scope and separate edge, control-plane, and origin hypotheses cleanly.

**Follow-ups:**
1. What do you check first after confirming user impact?
2. When do you roll back before proving the root cause?
3. How do you tell whether the issue is one region or global?
4. Which evidence can mislead you during overload?
5. What should the postmortem preserve from the debugging story?

## Ship It

- `outputs/debugging-narrative-worksheet.md`
- `outputs/interview-card-debugging-narrative.md`

## Exercises

1. **Easy** — Write a three-step debugging story for a sudden API error spike.
2. **Medium** — Compare how your narrative changes for a slow-burn backlog incident instead of a sudden outage.
3. **Hard** — Explain how you would debug a partial regional failure when logs are delayed and traces are heavily sampled.

## Further Reading

- [Effective Troubleshooting](https://sre.google/sre-book/effective-troubleshooting/) — excellent structure for debugging narratives
- [Postmortem Culture](https://sre.google/sre-book/postmortem-culture/) — useful for connecting incident reasoning to learning
