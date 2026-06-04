# Skill Sheet — Messaging Platform Drill

## Opening Template

1. Clarify whether this is queue, workflow, scheduler, or pub/sub.
2. State the guarantee boundary.
3. Estimate throughput, retention, and retry or replay amplification.
4. Propose the high-level design.
5. Choose one deep dive.

## Good Deep-Dive Targets

- partitioning and lag recovery
- timer durability and workflow state
- retry amplification and missed-run policy
- subscriber isolation and replay control

## Close Strong

- name two key trade-offs
- name the most likely failure mode
- state the main observability signals
- adapt the design when a constraint changes
