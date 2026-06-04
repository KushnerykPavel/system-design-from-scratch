# Trade-off Matrix: Eviction Policies

| Policy | Strongest on | Weakest on | Operational note |
|--------|--------------|------------|------------------|
| FIFO | simplicity | burst recency and popularity shifts | easiest to explain, often weakest hit rate |
| LRU | short-term recency | large scans and pollution | great default when recent usually means valuable |
| LFU | stable long-term popularity | trend shifts and stale frequency | usually needs decay or aging in real systems |

## Review prompts

- What is the estimated working set size?
- Are objects similarly sized?
- Does the workload contain one-hit scans?
- Is origin offload more important than raw hit ratio?
- Would admission control help more than changing eviction?
