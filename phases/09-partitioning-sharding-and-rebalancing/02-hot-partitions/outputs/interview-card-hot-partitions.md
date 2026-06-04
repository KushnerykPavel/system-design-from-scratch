---
lesson: 02-hot-partitions
focus: balanced
---

## Clarify first

- hotspot source: key, tenant, or time pattern
- read vs write vs retry amplification
- freshness and correctness constraints

## Must-size numbers

- top partition share of traffic
- burst factor
- top tenant vs median tenant volume

## Core design

- diagnose concentration before scaling
- isolate neighbors
- choose targeted mitigation

## Failure probes

- retries worsen the incident
- the hot tenant keeps growing
- cache is unavailable or stale reads are disallowed

## Trade-off summary

- fast relief vs durable fix
- balance vs locality
- isolation strength vs cost
