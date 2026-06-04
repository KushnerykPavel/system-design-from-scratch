# Playbook — Failure Redesign

## 1. State The Incident

- what failed first
- what amplified it
- what users saw

## 2. Name The Broken Assumption

- too much retry freedom
- stale state was treated as safe
- rollout lacked validation
- one shard or tenant could consume the system

## 3. Redesign

- containment change
- detection change
- rollback or recovery change
- one honest new trade-off

## 4. Validate

- metric that would have caught it sooner
- metric that proves smaller blast radius
- drill or game day to rehearse the fix
