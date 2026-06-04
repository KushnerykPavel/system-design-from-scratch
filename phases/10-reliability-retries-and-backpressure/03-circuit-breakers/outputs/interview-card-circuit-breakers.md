---
lesson: 03-circuit-breakers
focus: balanced
---

## Clarify first

- which dependencies are optional versus critical
- what degraded behavior is acceptable
- how recovery will be probed safely

## Core design

- trip policy from error, timeout, or saturation signals
- explicit fallback mode per request class
- limited half-open probes before reopening fully

## Failure probes

- fallback uses the same bottleneck
- breaker trips too late
- half-open floods a recovering dependency

## Trade-off summary

- availability vs freshness
- protection vs false trips
- recovery speed vs safety
