---
lesson: 07-cross-shard-queries
focus: balanced
---

## Clarify first

- critical-path vs analytical query
- exactness and freshness target
- expected shard fanout

## Must-size numbers

- shard count
- cross-shard QPS
- merge CPU and network cost

## Core design

- keep core reads shard-local when possible
- bound fanout when acceptable
- precompute or derive common global views

## Failure probes

- query cost rises with shard count
- materialized view lags silently
- pagination duplicates or drops items

## Trade-off summary

- simple live reads vs scalable derived views
- exactness vs freshness and cost
- broad query flexibility vs primary-store protection
