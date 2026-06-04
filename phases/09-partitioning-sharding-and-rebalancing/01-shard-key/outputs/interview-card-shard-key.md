---
lesson: 01-shard-key
focus: balanced
---

## Clarify first

- what must stay single-shard
- whether isolation is per tenant, user, or region
- how uneven the largest customers are

## Must-size numbers

- top tenant share of traffic
- scatter-gather rate on critical reads
- write concentration risk by key prefix

## Core design

- route from workload backward
- keep the shard key in primary and secondary data paths
- name the directory or placement layer if the key is indirect

## Failure probes

- one customer becomes 100x larger
- common reads become cross-shard
- routing metadata goes stale during moves

## Trade-off summary

- locality vs balance
- isolation vs routing simplicity
- future migration flexibility vs self-locating lookup speed
