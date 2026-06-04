---
lesson: 06-retry-budgets
focus: google
---

## Clarify first

- whether the goal is success rate, p99, or both
- which requests are safe to hedge
- how much spare capacity exists for speculative work

## Must-size numbers

- baseline QPS
- extra-attempt budget ratio
- hedge trigger percentile and win rate

## Core design

- separate retry and hedge budgets
- hedge only late safe requests
- cancel losers aggressively

## Failure probes

- correlated replicas
- no real cancellation
- budget burn stays high even after p99 improves
