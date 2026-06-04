# Consistency Checklist - Google Serving and Storage

## Define First

- what is the source of truth?
- which writes must be immediately correct?
- which reads can be stale?
- which anomaly is acceptable?
- which anomaly is forbidden?

## Verify Next

- replica lag is measurable
- failover cannot silently promote bad state
- strong paths are narrow and justified
- caches have freshness or invalidation rules

## Say Out Loud

- what gets slower
- what gets less available
- what gets more operationally complex
