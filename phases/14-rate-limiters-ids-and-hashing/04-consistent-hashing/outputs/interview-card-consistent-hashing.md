# Interview Card — Consistent Hashing

## Core answer

1. Explain why modulo hashing causes excessive remap on fleet changes.
2. Introduce consistent hashing as a bounded-remap placement strategy.
3. Add virtual nodes for smoother balance.
4. Mention operational details: health, ring versioning, staged rebalance, hot-key mitigation.
5. Close with the metric: remap percentage and ownership skew.

## Good closing line

"Consistent hashing solves churn better than it solves skew, so I would pair it with hot-key handling and a safe rebalance plan."
