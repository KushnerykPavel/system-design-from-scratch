# Regional Failover and Recovery Objectives

> Failover is not a button. It is a timed promise about what still works, how much data is lost, and who is allowed to pull the lever.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Translate RTO and RPO targets into concrete multi-region failover design, automation, and operational guardrails.
**Prerequisites:** `10-reliability-retries-and-backpressure/01-timeouts-and-retries`, `11-observability-slos-and-debugging/06-runbooks`, `13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive`
**Estimated time:** ~75 min
**Primary artifact:** recovery worksheet + failure checklist

## The Problem

You have already decided the service will span multiple regions. Now you must answer the harder question: how quickly do you fail over, how much data can you lose, and what evidence lets you act safely?

Candidates often say "we fail over to another region" without converting that into timers, dependencies, and operator workflows. This lesson makes those promises explicit.

## Clarify

- What are the target RTO and RPO values?
- Is failover automatic, operator-approved, or customer-triggered?
- Which dependencies are regional and which are globally shared?
- Is degraded read-only service acceptable during recovery?

If the interviewer gives no numbers, assume an interactive user-facing service with RTO in minutes and RPO measured in seconds to low minutes.

## Requirements

### Functional

- Detect regional failure and decide whether to fail over.
- Promote traffic and data roles without making recovery worse.
- Support validation, rollback, and controlled failback.

### Non-functional

- Detection must avoid flapping on partial regional incidents.
- Automation must be safer than ad hoc operator action.
- Recovery steps should be auditable and rehearsable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Detection window | 30 to 90 seconds | sets lower bound on realistic RTO |
| Target RTO | 5 minutes | drives automation and pre-provisioning |
| Target RPO | under 30 seconds | drives replication mode and promotion rules |
| Critical dependencies | 5 to 12 | failover breaks if even one is forgotten |
| Rough cost | extra standby + drills + control plane | recovery promises are not free |

## Architecture

```text
health probes + synthetic checks
  -> incident classifier
     -> failover controller
        -> traffic manager
        -> data-role promotion
        -> readiness verification
        -> operator approval or auto-cutover
```

Key steps:

1. Detect whether the issue is local, zonal, regional, or dependency-wide.
2. Verify the target region is healthy enough to absorb traffic.
3. Shift traffic, promote state where needed, and enforce degraded-mode policy.
4. Track whether recovery met the promised RTO and RPO.

## Data Model & APIs

Core entities:

- `DependencyStatus`
- `RecoveryObjective`
- `FailoverPlan`
- `PromotionGate`
- `DrillResult`

Useful APIs:

- `AssessIncident(scope)`
- `CanPromote(region)`
- `ExecuteFailover(plan)`
- `RecordDrill(result)`

The most important data model question is which signals are authoritative enough to promote a new primary safely.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| false positive regional failover | failover started while local-region metrics were mixed | combine black-box probes, dependency status, and manual override gates |
| target region is healthy but underprovisioned | surge causes saturation after cutover | reserve capacity and run absorption drills |
| data promotion violates RPO | promoted region missing recent writes | gate cutover on replication freshness and known replay procedures |
| failover controller becomes single point of failure | no safe actuation during control-plane incident | isolate control plane and provide limited manual emergency path |

## Observability

- metric: detection time, classification confidence, and time-to-promote
- metric: replication freshness and promotion eligibility by region
- metric: post-cutover saturation, error rate, and queue growth in target region
- log: every cutover decision, approval, and rollback event
- trace: synthetic request path before and after traffic shift
- SLO: measured RTO and RPO during drills should stay within target bands

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| automated failover with approval gates | faster recovery without blind action | more control-plane complexity | purely manual cutover under pressure |
| degraded read-only mode during recovery | preserves partial availability | reduced functionality | full write availability at any cost |
| aggressive detection thresholds | lower outage duration | higher risk of flapping | conservative thresholds that miss true failures |

## Interview It

**Google framing:** "How would you design disaster recovery for this service?" Expect questions on RTO realism, promotion safety, and operator workflow.

**Cloudflare framing:** "How do you shift global traffic safely when a region is unhealthy?" Expect pressure on detection quality, actuation safety, and capacity absorption.

**Follow-ups:**
1. What part of the RTO budget is usually underestimated?
2. What if the data tier has fresher replicas than the traffic manager thinks?
3. How would you keep the system available in read-only mode?
4. How do you stop failover from flapping?
5. How often do you drill, and what do you measure?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/capacity-sheet-failover-and-rto.md`
- `outputs/failure-checklist-failover-and-rto.md`

## Exercises

1. **Easy** — Budget a 5-minute RTO across detection, approval, cutover, and warm-up.
2. **Medium** — Redesign when RPO must be near-zero for a subset of writes.
3. **Hard** — Support three regions where one dependency is only dual-region and may lag promotion.

## Further Reading

- [Google SRE Workbook - Emergency Response](https://sre.google/workbook/emergency-response/) — practical operational framing for incident action
- [AWS - Disaster Recovery of Workloads on AWS](https://docs.aws.amazon.com/whitepapers/latest/disaster-recovery-workloads-on-aws/disaster-recovery-workloads-on-aws.html) — helpful comparison vocabulary for warm, pilot-light, and active setups
