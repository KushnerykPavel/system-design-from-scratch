# Workflow Engine for Long-Running Jobs

> Long-running workflows are mostly a state-management problem that happens to call other systems.

**Type:** Build
**Company focus:** Google
**Learning goal:** Design a workflow engine that handles durable state, timers, retries, compensation, and code evolution for long-running business or infrastructure processes.
**Prerequisites:** `07-queues-streams-and-workflows/05-workflow-engines`, `08-consistency-replication-and-transactions/06-sagas`, `10-reliability-retries-and-backpressure/02-idempotency-under-failure`
**Estimated time:** ~75 min
**Primary artifact:** workflow-plan validator + failure checklist

## The Problem

Design a workflow engine for processes that can run for minutes to days: onboarding, infrastructure provisioning, payment review, or recovery playbooks. Steps may wait on timers, human approval, or external systems, and the platform must recover from crashes without losing workflow state or replaying dangerous side effects blindly.

This lesson matters because many candidates reduce workflow engines to "queue + workers." Senior answers explain why durable state machines, timer handling, activity idempotency, and workflow versioning matter once work outlives a single process.

## Clarify

- Are workflows mostly machine-driven, human-in-the-loop, or mixed?
- Do steps call external side-effecting systems that need compensation?
- What is the longest execution time we must support?
- Do workflow definitions change while instances are still running?

If the prompt stays broad, assume machine-driven workflows with occasional waits, at-least-once activity execution, hours-to-days runtime, and rolling workflow-definition upgrades.

## Requirements

### Functional

- Start, persist, resume, and inspect long-running workflows.
- Schedule timers, retries, and wakeups without keeping workers pinned.
- Execute activities with retry policy and idempotency guidance.
- Support compensation or rollback-aware orchestration for partial failure.
- Allow operators to pause, resume, cancel, or replay workflows safely.

### Non-functional

- Survive worker and coordinator crashes without losing progress.
- Prevent one workflow family from starving others.
- Keep operator visibility high for stuck and flapping workflows.
- Evolve workflow code without corrupting in-flight executions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| New workflow starts | 300K/min peak | drives coordinator throughput and state writes |
| Concurrent open workflows | 150M | storage, indexing, and timer scale dominate |
| Timer events | 20M/min | timer-wheel or scheduled-queue design matters |
| Average activities per workflow | 12 | retry fanout can multiply write load |
| Peak factor | 4x during incident automation | workflows often surge exactly when dependencies are unhealthy |

## Architecture

```text
workflow API
  -> workflow coordinator
  -> durable workflow state store
  -> timer queue
  -> task queues for activities
  -> worker pools
  -> event history / operator console
```

Design notes:

1. Separate workflow state transitions from activity execution so workers stay stateless and replaceable.
2. Model timers as persisted events, not sleeping threads.
3. Require idempotent or fenced activities because retries are routine, not exceptional.
4. Treat workflow versioning as core platform functionality if runs can outlive a deploy.

## Data Model & APIs

Core records:

```text
workflow_id
workflow_type
workflow_version
state
current_step
next_wakeup_time
attempt
activity_token
compensation_state
```

Useful interfaces:

- `StartWorkflow(type, input, idempotency_key)`
- `PollActivityTask(queue, worker_identity)`
- `CompleteActivity(task_token, result)`
- `HeartbeatActivity(task_token, progress)`
- `SignalWorkflow(workflow_id, signal_name, payload)`
- `CancelWorkflow(workflow_id, reason)`

A strong answer distinguishes workflow history, current materialized state, and operator-facing controls.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| worker dies during long activity | heartbeat timeout and stuck-attempt age | task timeout plus safe retry or compensation |
| timer service falls behind | wakeup lag and overdue-timer count | shard timer queues and prioritize overdue expirations |
| workflow code changes break old executions | version mismatch and replay divergence | workflow version pinning and compatibility gates |
| one workflow family floods the engine | queue saturation by workflow type | per-type quotas, reserved pools, and admission control |

## Observability

- metric: workflow-start latency and state-transition throughput
- metric: open workflows, overdue timers, and stuck-activity count
- metric: retry attempts and compensation execution rate
- metric: workflow-type saturation and worker-pool backlog
- log: operator actions such as pause, cancel, force retry, and replay
- trace: workflow start through activity chain with state-transition annotations
- SLO: 99% of timer wakeups fire within the workflow platform target for standard-priority workflows

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| durable event history | strong recovery and audit | storage and replay overhead | ephemeral in-memory orchestration |
| explicit compensation hooks | safer side-effect recovery | workflow author complexity | pretending rollback is automatic |
| versioned workflow definitions | safe long-lived execution | more platform machinery | overwriting definitions in place |

## Interview It

**Google framing:** "Design a workflow engine for internal business processes and infrastructure operations." Expect focus on retries, timers, determinism, and safe evolution of running workflows.

**Cloudflare framing:** "Design an orchestration engine for long-running platform operations." Expect follow-ups on blast radius, operator intervention, and how unhealthy dependencies affect retry storms.

**Follow-ups:**
1. How do you prevent duplicate side effects during activity retries?
2. What changes if workflows need manual approval steps?
3. How do you migrate millions of in-flight workflows to a new definition?
4. When would you use a queue alone instead of a workflow engine?
5. How do you handle workflows that must span regions?

## Ship It

- `outputs/failure-checklist-workflow-engine.md`

## Exercises

1. **Easy** — Explain why timers belong in durable state, not sleeping workers.
2. **Medium** — Add a compensation plan for a provisioning workflow that fails halfway through.
3. **Hard** — Redesign the engine so workflow definitions can evolve safely while old instances keep running for weeks.

## Further Reading

- [Temporal durable execution](https://docs.temporal.io/workflows) — useful mental model for durable workflow orchestration
- [Sagas](https://microservices.io/patterns/data/saga.html) — background on compensation and multi-step failure recovery
