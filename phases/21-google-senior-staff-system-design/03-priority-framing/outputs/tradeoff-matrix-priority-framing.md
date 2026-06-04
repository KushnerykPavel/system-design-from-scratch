# Trade-off Matrix - Priority Framing

| Priority That Wins | Typical Design Direction | Common Cost |
|--------------------|--------------------------|-------------|
| lowest read latency | caching, precompute, regional replicas | staleness and higher invalidation complexity |
| strongest write correctness | tighter write quorum, stronger source of truth | higher latency and lower availability during faults |
| lowest cost | simpler topology, fewer replicas, less precompute | weaker latency or resiliency margins |
| fastest product iteration | modular boundaries, simpler contracts first | more cleanup and migration later |
| strongest regional isolation | per-region data and control boundaries | duplication and operational overhead |
