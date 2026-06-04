# Trade-off Matrix — News Feed / Timeline

| Decision | Prefer When | Main Cost | Watch Metric |
|----------|-------------|-----------|--------------|
| push fanout | follower sets are modest and read latency dominates | write amplification | fanout queue depth |
| pull fanout | authors have huge follower counts or ranking is highly dynamic | expensive read-time merge | timeline p99 |
| mixed fanout | graph skew is real and workloads vary widely | control-plane complexity | celebrity read amplification |
| materialized timelines | product needs predictable low-latency reads | delete propagation and storage cost | time-to-hide deleted content |
| ranking fallback | uptime matters more than perfect relevance | lower engagement quality | stale-feature fallback rate |
