# Interview Card — Cache Hit Rate

## Say this early

- "I want to translate hit rate into origin QPS, because that is the real capacity number."
- "A 15-point hit-rate drop can more than triple origin load depending on the starting point."

## What to compute

1. read QPS
2. hit ratio
3. origin miss QPS
4. miss latency or miss cost

## Common misses

- using only global hit rate
- ignoring hot-object skew
- ignoring invalidation events
- assuming misses are cheap
