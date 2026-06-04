# Trade-off Matrix — Distributed Cache Cluster

| Decision | Benefit | Cost | Best Watch Metric |
|----------|---------|------|-------------------|
| shared cluster | high memory efficiency | noisy-neighbor risk | tenant hit rate |
| segmented pools | better isolation | lower overall utilization | eviction churn per pool |
| replication in cache | better availability on node loss | more memory cost and write work | warm failover hit rate |
| request coalescing | protects origin during stampedes | coordination on miss path | duplicate miss suppression count |
| short TTL | stronger freshness | more origin pressure | origin QPS after expiry waves |
