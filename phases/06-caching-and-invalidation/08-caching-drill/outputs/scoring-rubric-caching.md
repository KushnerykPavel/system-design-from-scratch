# Scoring Rubric: Caching Drill

| Dimension | Strong | Weak |
|-----------|--------|------|
| Clarification | asks about mutability, hot paths, and freshness | asks only generic product questions |
| Freshness | states a clear stale-read contract | hides behind one arbitrary TTL |
| Miss path | explains hot-key expiry and refill control | talks only about hit ratio |
| Consistency | names read-after-write or bounded staleness explicitly | says "eventual consistency" without detail |
| Observability | picks metrics for skew, misses, and invalidation | says only "monitor the cache" |
| Trade-offs | names latency, cost, and correctness tension | lists components without trade-offs |
