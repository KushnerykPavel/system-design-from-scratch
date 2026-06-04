---
lesson: 08-sharding-drill
focus: balanced
---

## Clarify first

- dominant locality boundary
- tenant skew
- freshness on reads and writes

## Must-size numbers

- read/write QPS
- largest tenant share
- cross-shard admin query rate

## Core design

- choose shard key from workload
- use placement indirection if moves matter
- isolate heavy tenants
- derive global views

## Failure probes

- giant tenant appears
- stale routing during migration
- admin dashboards hit primaries directly

## Trade-off summary

- routing flexibility vs extra indirection
- stronger isolation vs cost
- migration safety vs rollout speed
