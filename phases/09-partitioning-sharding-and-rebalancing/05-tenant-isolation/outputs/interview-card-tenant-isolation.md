---
lesson: 05-tenant-isolation
focus: balanced
---

## Clarify first

- fairness target: latency, throughput, or compliance
- largest-tenant skew
- shared bottleneck resource

## Must-size numbers

- largest vs median tenant
- shard spare capacity
- burst size for premium tenants

## Core design

- per-tenant accounting
- quotas and admission control
- tiered placement and dedicated pools when justified

## Failure probes

- background work, not request QPS, causes the incident
- premium tenants need stronger SLOs
- manual isolation actions are too slow

## Trade-off summary

- stronger isolation vs higher cost
- shared efficiency vs fairness guarantees
- manual control vs automation complexity
