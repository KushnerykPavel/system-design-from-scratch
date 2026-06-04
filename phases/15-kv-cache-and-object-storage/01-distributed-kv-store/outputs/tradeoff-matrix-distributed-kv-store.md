# Trade-off Matrix — Distributed KV Store

| Decision | Prefer When | Main Cost | Watch Metric |
|----------|-------------|-----------|--------------|
| quorum reads and writes | correctness matters more than cheapest latency | extra fan-out and tail latency | quorum timeout rate |
| local read from one replica | stale reads are acceptable | possible freshness drift | replica lag age |
| 3 replicas | standard regional durability | 3x storage and repair traffic | repair backlog bytes |
| 5 replicas | very high availability and read locality | more write fan-out and cost | write p99 and cross-AZ traffic |
| LSM engine | write-heavy or append-heavy patterns | compaction amplification | compaction debt |
| leadered partition | strict ordering per key is valuable | leader bottleneck and failover complexity | leader queue depth |
