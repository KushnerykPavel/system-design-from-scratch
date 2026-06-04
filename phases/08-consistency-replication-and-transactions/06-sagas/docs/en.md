# Sagas and Compensating Actions

> A saga is not a free transaction replacement. It is a workflow that accepts intermediate state and pays for reversibility.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Decide when to replace a wide transaction with a saga and explain compensation, idempotency, and user-visible state transitions.
**Prerequisites:** `05-transactions`, `07-queues-streams-and-workflows/05-workflow-engines`, `07-queues-streams-and-workflows/06-outbox-and-cdc`
**Estimated time:** ~75 min
**Primary artifact:** saga design review worksheet

## The Problem

Once a design crosses service boundaries, the temptation is to keep one big atomic story. In practice that is often brittle or impossible. Sagas are the alternative, but they introduce their own demands:

- each step needs durable intent
- compensation must be safe and realistic
- users may observe in-progress or partially completed state
- retries and duplicates must not multiply side effects

This lesson helps you talk about sagas as an operational workflow, not just a buzzword.

## Clarify

- Which side effects can be undone, and which only need reconciliation?
- Can users see an in-progress state, or must completion appear atomic?
- Is orchestration centralized, choreographed by events, or mixed?
- What is the worst business impact of a stuck partially completed flow?

## Requirements

### Functional

- Coordinate multi-step workflows across service boundaries.
- Persist step state and retry safely.
- Compensate or reconcile when downstream steps fail.

### Non-functional

- Keep long-running state inspectable and restartable.
- Bound duplicate side effects through idempotency.
- Make irreversibility explicit rather than pretending compensation is universal.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Workflow starts | 5K/s peak | enough to require durable workflow state |
| Average steps | 4-7 | enough branching to need visibility |
| Failure/retry rate | 1-3% normal | compensation is not rare |
| Long-tail duration | up to 30 minutes | state recovery matters |
| Rough cost | workflow storage, retries, audits, compensation logic | shows why sagas are deliberate, not free |

## Architecture

A strong saga answer names:

1. **Workflow state machine**
2. **Per-step idempotency**
3. **Forward action and compensation or reconciliation**
4. **Timeout, retry, and stuck-workflow handling**

Typical states:

- pending
- reserved
- charged
- fulfilled
- compensating
- failed

## Data Model & APIs

Helpful interfaces:

- `StartWorkflow(workflow_id, input)`
- `AdvanceStep(workflow_id, step_id, idempotency_key)`
- `Compensate(workflow_id, failed_step)`
- `GetWorkflowStatus(workflow_id)`

Helpful metadata:

```text
workflow -> {
  workflow_id,
  status,
  step_states,
  retry_count,
  last_error
}
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one step succeeds twice | duplicate side effect metric or audit mismatch | idempotency keys and step result store |
| compensation is impossible in reality | stuck workflows and manual queue grow | redesign step order or add reconciliation path |
| users see confusing partial state | support incidents and status mismatch | explicit workflow status and progress model |
| orchestrator crash loses progress | missing heartbeat or state transition stall | durable workflow state and replayable step intents |

## Observability

- metric: workflow completion, failure, and compensation rate
- metric: retries per step and oldest in-progress workflow age
- metric: manual intervention queue size
- log: state transitions with workflow ID, step ID, and idempotency key
- trace: one workflow across all step executions and compensations
- SLO: completion success and bounded stuck-workflow age

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| saga instead of wide transaction | lower synchronous coupling | intermediate state and compensation complexity | distributed transaction across many services |
| orchestrated workflow | clearer visibility and control | central workflow state dependency | pure choreography with diffuse ownership |
| explicit user-visible progress | honest product contract | more UX and state handling | pretending the flow is atomic when it is not |

## Interview It

**Google framing:** "Design checkout or provisioning across several internal services." The signal is whether you know what can be local atomic work versus workflow coordination.

**Cloudflare framing:** "Design a long-running control-plane rollout with rollback and retries." The signal is whether you handle irreversible steps, status visibility, and operator recovery.

**Follow-ups:**
1. Which steps can truly be compensated?
2. What if one step calls an external provider that is not fully idempotent?
3. What if users need progress updates for a 20-minute flow?
4. When is choreography acceptable?
5. What metric most quickly reveals a stuck saga fleet?

## Ship It

- `outputs/saga-design-review.md`

## Exercises

1. **Easy** - Break a checkout flow into local transactions and compensating steps.
2. **Medium** - Explain how to retry an external callback safely.
3. **Hard** - Design a control-plane rollout saga where rollback is partial and some actions require manual remediation.

## Further Reading

- [Microservices patterns - sagas](https://microservices.io/patterns/data/saga.html) - useful conceptual grounding for compensation and workflow choices
- [System design notes](https://github.com/liquidslr/system-design-notes) - general interview structure before workflow-specific nuance
