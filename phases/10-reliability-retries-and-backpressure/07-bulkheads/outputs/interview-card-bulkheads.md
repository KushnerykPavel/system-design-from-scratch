---
lesson: 07-bulkheads
focus: cloudflare
---

## Clarify first

- what must fail independently
- which pools, queues, tenants, or cells map to that boundary
- whether unused capacity can be borrowed

## Core design

- dedicate capacity where blast radius matters
- isolate optional features from core traffic
- keep local failures local when possible

## Failure probes

- one cell fails and healthy cells overload
- one tenant dominates a shared pool
- fallback still shares the same bottleneck

## Trade-off summary

- resilience vs utilization
- strong partitioning vs operational complexity
- local autonomy vs central coordination
