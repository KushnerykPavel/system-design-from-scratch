# Failure Checklist — Workflow Engine

Review these when pressure-testing a workflow-engine design.

## State

- What persists the current step and retry state?
- Can the coordinator crash without losing progress?
- How are stuck workflows found?

## Activities

- Are activities retried?
- Are activities idempotent or fenced?
- How are long-running steps heartbeated?

## Timers

- Are timers durable?
- Can overdue timers be measured?
- What happens when timer shards fall behind?

## Change Management

- How do old workflows survive new code?
- Can operators pause, resume, and cancel safely?
- Is replay or reset scoped and auditable?
