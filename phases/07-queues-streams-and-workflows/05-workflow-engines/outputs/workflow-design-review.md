---
lesson: 07-workflow-engines
focus: balanced
---

# Workflow Design Review

## Clarify First

- Maximum workflow lifetime
- Status visibility requirements
- External callbacks and human approvals
- Reversible versus irreversible steps

## Must Model Explicitly

- Workflow instance state
- Step transitions
- Timer or deadline behavior
- Signal and callback handling
- Cancellation and compensation

## Failure Review

- Worker crash during activity
- Missing callback
- Duplicate signal
- Partial rollback after irreversible step

## Observability

- Oldest active workflow age
- Stuck step counts by type
- Timer backlog
- Compensation rate
