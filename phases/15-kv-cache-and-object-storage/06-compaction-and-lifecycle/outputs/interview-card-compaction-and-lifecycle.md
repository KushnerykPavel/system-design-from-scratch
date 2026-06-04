---
lesson: 06-compaction-and-lifecycle
focus: balanced
---

## Clarify first
- Main maintenance risk: compaction, tombstones, lifecycle backlog, or all three
- Acceptable lag before reclaimed data or transitioned storage must catch up
- Compliance and legal-hold constraints

## Must-size numbers
- Daily overwrite or delete volume
- Background IO budget
- Lifecycle action rate
- Maximum tolerated maintenance debt

## Core design
- Explicit maintenance schedulers
- Safety windows before destructive GC
- Dry-run and staged rollout for lifecycle policies

## Failure probes
- Early tombstone reclamation
- Compaction debt during ingest spike
- Lifecycle delete colliding with retention

## Trade-off summary
- Lower debt vs foreground IO cost
- Faster deletion vs correctness safety
- Faster rollout vs operator confidence
