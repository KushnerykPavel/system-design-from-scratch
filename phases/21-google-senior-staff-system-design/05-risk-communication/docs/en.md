# Communicating Risks, Rollouts, and Fallbacks

> Senior design answers improve when they stop pretending deployment is instant and failure is clean.

**Type:** Learn
**Company focus:** Google
**Learning goal:** Learn how to close a design with credible risk communication, rollout sequencing, and fallback behavior instead of ending at the diagram.
**Prerequisites:** `10-reliability-retries-and-backpressure/03-circuit-breakers`, `11-observability-slos-and-debugging/06-runbooks`, `13-multi-region-cdn-and-edge-traffic/02-failover-and-rto`
**Estimated time:** ~60 min
**Primary artifact:** rollout and fallback checklist

## The Problem

Many otherwise good answers weaken near the end. The candidate describes the architecture well but has little to say about how it ships, how it is rolled out safely, what happens during partial failure, or how operators retreat if the new design misbehaves.

Google senior/staff interviewers often see this as a maturity gap. Owning a system includes owning its introduction and its failure modes.

## Clarify

- Is this a net-new system, a migration, or a replacement of an existing path?
- Which failure is most dangerous: data loss, downtime, bad latency, or policy mistakes?
- Can we fail open, fail closed, or degrade to a simpler mode?
- Is rollout global at once, regional, tenant-based, or traffic-sliced?

If rollout context is missing, assume a staged migration with rollback checkpoints and one reduced-functionality fallback mode.

## Requirements

### Functional

- Identify the top operational risks in the proposed design.
- Define a staged rollout strategy and blast-radius controls.
- Explain the primary fallback or degradation path.
- State what metrics or alerts would halt the rollout.
- Show how the system returns to a safe state if the new path fails.

### Non-functional

- Avoid all-or-nothing launch plans.
- Avoid vague "we can monitor it" language without specific signals.
- Keep fallback behavior consistent with product and security requirements.
- Make rollback cost and data migration risk explicit.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Rollout stages | 3 to 5 | enough to reduce blast radius meaningfully |
| Health gates per stage | 2 to 4 | rollout should pause on concrete signals |
| Fallback latency hit | 1.2x to 3x typical | degraded modes often protect correctness at a speed cost |
| Migration overlap window | hours to weeks | dual-run and rollback burden can dominate operations |

## Architecture

Close with an operator view:

```text
new path behind flag
  -> canary scope
  -> health checks and SLO gates
  -> staged traffic expansion
  -> fallback or rollback path
```

Good closing questions to answer yourself:

1. What goes wrong first if the new design is flawed?
2. What signal detects that failure earliest?
3. How do we reduce blast radius while learning?
4. What simpler fallback preserves the most important guarantee?

Strong candidates distinguish:

- rollout safety from steady-state design
- fallback behavior from full rollback
- degraded correctness from degraded performance

## Data Model & APIs

Useful rollout model:

```text
rollout_stage(percent, scope, health_gate, rollback_action)
fallback_mode(mode, preserved_guarantee, degraded_property)
```

Helpful verbal APIs:

- `Canary(scope)`
- `Promote(stage)`
- `Pause(reason)`
- `Fallback(mode)`
- `Rollback(version)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| rollout is described as one big launch | no stage or blast-radius control exists | define tenant, region, or traffic canaries |
| fallback is vague | no preserved guarantee is named | state exactly what remains correct in degraded mode |
| migration risk is ignored | dual-write or backfill behavior is missing | define overlap window and reconciliation checks |
| stopping criteria are unclear | rollout has no health gates | tie promotion to SLO, error rate, or correctness metrics |

## Observability

- metric: rollout health by stage, region, and tenant slice
- metric: new-path versus old-path latency and error deltas
- metric: fallback activation count and duration
- metric: migration lag, dual-write mismatch, or reconciliation backlog
- log: pause, rollback, and fallback decisions with reasons
- trace: request behavior through old path, new path, and degraded mode
- SLO: rollout only expands when the preserved guarantee remains within target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| staged rollout | smaller blast radius | slower launch and more control logic | global cutover |
| explicit fallback mode | safer incidents | lower functionality or higher latency | undefined failure behavior |
| dual-run or overlap | better migration confidence | extra cost and complexity | immediate replacement with no verification |

## Interview It

**Google framing:** "How would you deploy this safely?" A strong answer should sound like an engineer who expects partial failure, bad assumptions, and the need to retreat gracefully.

**Follow-ups:**
1. What if the fallback path is more expensive?
2. What if rollback is hard because writes are not reversible?
3. What if only one region is unhealthy?
4. What if the new path is correct but slower than expected?

## Ship It

- `outputs/failure-checklist-risk-rollout-fallback.md`

## Exercises

1. **Easy** - Define a canary plan for a new read path.
2. **Medium** - Explain the difference between fallback and rollback for a migrated write path.
3. **Hard** - Design rollout gates for a global API where one region can fail independently.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) - useful for change management and progressive delivery thinking
- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) - important for reasoning about why partial regressions matter during rollout
