# Failure Checklist — Consistent Hashing

- Did a topology change unexpectedly remap too many keys?
- Is ring version skew causing nodes to disagree on ownership?
- Are hot keys concentrated on one owner even though average balance looks fine?
- Did a node add or removal trigger a cold-cache wave or migration spike?
- Are health signals fast enough to stop routing to dead owners?

Healthy answer pattern:

- quantify remap
- explain warmup plan
- call out hot-key exceptions
- mention ring version observability
