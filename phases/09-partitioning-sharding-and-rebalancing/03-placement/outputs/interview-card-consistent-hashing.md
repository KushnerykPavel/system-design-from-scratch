---
lesson: 03-placement
focus: balanced
---

## Clarify first

- what is being placed: requests, keys, ranges, or replicas
- how often topology changes
- whether capacity is heterogeneous

## Must-size numbers

- acceptable remap on node add
- ownership imbalance tolerance
- warmup or migration cost per moved key

## Core design

- deterministic placement
- bounded churn on topology change
- fault-domain-aware replication if durability matters

## Failure probes

- stale ring view
- few virtual nodes cause imbalance
- replicas land together

## Trade-off summary

- simpler modulo hashing vs safer remap behavior
- smoother balance vs larger control-plane state
- weighted flexibility vs rollout complexity
