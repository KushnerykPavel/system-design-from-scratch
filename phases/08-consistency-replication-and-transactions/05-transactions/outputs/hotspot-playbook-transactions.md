# Hotspot Playbook - Transactions

## Typical hotspots

- one account or wallet row
- one inventory SKU during a launch
- one tenant quota bucket
- one global counter or ordering table

## Response options

| Pattern | Helps when | Cost |
|---------|------------|------|
| per-entity queue or serialization | one key dominates | added latency and queue ops |
| reservation or escrow counters | bounded decrement-style invariants | more modeling complexity |
| partition by owner key | hotspots are separable by tenant or account | migration and skew planning |
| move noncritical work out of transaction | extra side effects dominate lock time | async coordination complexity |

## Watch closely

- lock wait time
- abort and retry rate
- hottest keys by conflict count
- p99 transaction duration
