---
lesson: 04-rebalancing
focus: balanced
---

## Clarify first

- durable data vs cache ownership
- active writes during move
- spare bandwidth and concurrency budget

## Must-size numbers

- bytes to move
- writes per second on moved ranges
- max concurrent moves

## Core design

- copy in bounded batches
- track epochs
- verify before cutover
- keep rollback alive

## Failure probes

- stale routing cache
- target lags ongoing writes
- move traffic hurts p99 latency

## Trade-off summary

- faster completion vs smaller blast radius
- dual-write safety vs implementation complexity
- retained source vs temporary extra capacity
