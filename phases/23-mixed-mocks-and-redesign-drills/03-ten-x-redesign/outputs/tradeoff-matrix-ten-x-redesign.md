# Trade-off Matrix — 10x Redesign

| Redesign move | Helps most when | Main benefit | Main cost |
|---------------|-----------------|--------------|-----------|
| repartitioning | hotspots dominate | spreads load better | migration complexity |
| more async work | sync path is overloaded | protects tail latency | freshness and backlog risk |
| hot/cold path split | workload classes diverge | better efficiency | more code paths |
| regionalization | one region is overloaded or distant | better locality and containment | consistency and rollout complexity |
| stricter admission control | peak bursts dominate | protects core SLOs | more rejected or delayed work |
