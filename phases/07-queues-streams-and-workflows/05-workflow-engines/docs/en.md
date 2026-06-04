# Workflow Engines and Long-Running State

> Retries are not a workflow model. They are only one behavior inside a workflow model.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Recognize when durable orchestration is required and design explicit workflow state, timers, compensation, and operator visibility.
**Prerequisites:** `07-queues-streams-and-workflows/01-queues-vs-streams`, `08-consistency-replication-and-transactions/06-sagas`
**Estimated time:** ~75 min
**Primary artifact:** workflow design review prompt

## The Problem

Some business processes last minutes, hours, or days. They may include:

- external callbacks
- human approval
- retries with backoff
- waiting on inventory, fraud review, or regional failover
- compensating actions when later steps fail

Trying to encode that only with queue retries usually creates hidden state spread across databases, cron jobs, and manual runbooks. Workflow engines make that state explicit.

## Clarify

- How long can the process remain active: seconds, minutes, days, or weeks?
- Does the product need a user-visible status page or only eventual completion?
- Which steps are reversible, and what compensation is possible?
- Are waits event-driven, timer-driven, human-driven, or all three?

## Requirements

### Functional

- Represent durable step-by-step progress for long-running operations.
- Support retries, timers, external signals, and compensation.
- Expose workflow status for users and operators.

### Non-functional

- Survive worker crashes and restarts without losing progress.
- Bound duplicate step execution and external side effects.
- Keep debugging simple enough that operators can explain stuck workflows.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Workflow starts | 20K/min | drives orchestration throughput |
| Active workflows | 5M | long-running state dominates storage and indexing |
| Average duration | 2 hours, tail to 14 days | makes timers and persistence central |
| External callbacks | 8 per workflow on average | raises ordering and idempotency needs |
| Rough cost | workflow state store + timer service + worker fleet | orchestration is a platform choice, not a free abstraction |

## Architecture

```text
API
  -> workflow engine
  -> durable workflow state
  -> task queues / activity workers
  -> timers + signals
  -> status API
```

Design goals:

1. Separate deterministic workflow state from side-effecting activities.
2. Make step transitions durable and inspectable.
3. Treat compensation as an explicit design path, not an afterthought.
4. Prefer resumable workflows over ad hoc cron repair.

## Data Model & APIs

Core entities:

- workflow instance
- step state
- timer
- signal or callback
- compensation action

Useful APIs:

- `StartWorkflow(input)`
- `Signal(workflow_id, event)`
- `QueryStatus(workflow_id)`
- `Cancel(workflow_id)`
- `RetryStep(workflow_id, step_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| worker crashes during external activity | activity timeout and retry count rise | durable state plus idempotent activities |
| workflow waits forever for missing callback | workflow age or stuck-step alerts fire | deadlines, timers, and escalation path |
| compensation logic is underspecified | inconsistent cross-system state after failure | define reversible boundaries and manual fallback |
| status API reads stale or partial progress | user support tickets and workflow-state mismatch | query durable workflow state directly or use versioned projections |

## Observability

- metric: workflow starts, completions, cancellations, and failures by type
- metric: active workflow count and oldest step age
- metric: timer backlog and callback timeout rate
- log: workflow ID, step transitions, and compensation decisions
- trace: one workflow across activities and waits
- SLO: critical workflow completion within target time percentile

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit workflow state | debuggable long-running progress | more platform and state complexity | hidden state across jobs and tables |
| durable timers and signals | resilient waits and callbacks | timer-service overhead | sleep loops and cron polling everywhere |
| compensation-aware design | safer partial-failure handling | more product and engineering effort | pretending all steps are atomic |

## Interview It

**Google framing:** "Design document export or provisioning workflow with user-visible progress." The signal is whether you make state and timers explicit rather than hand-waving retries.

**Cloudflare framing:** "Design certificate issuance or edge policy rollout workflow." The signal is whether you reason about long-running orchestration, approvals, and rollback paths.

**Follow-ups:**
1. What if a human approval step can take days?
2. What if one activity is irreversible, like sending an email or contacting an external bank?
3. What if product needs to cancel in-flight workflows cleanly?
4. What if status reads must be globally visible with low latency?
5. When is a queue plus table enough, and when is a workflow engine justified?

## Ship It

- `outputs/workflow-design-review.md`

## Exercises

1. **Easy** — List the minimum state you need to expose workflow status safely.
2. **Medium** — Explain how workflow timers differ from consumer retry backoff.
3. **Hard** — Redesign a fragile cron-driven provisioning flow into explicit durable orchestration.

## Further Reading

- [Temporal concepts](https://docs.temporal.io/workflows) — practical language for workflow state, activities, and signals
- [Sagas](https://microservices.io/patterns/data/saga.html) — useful background on compensation-driven process design
