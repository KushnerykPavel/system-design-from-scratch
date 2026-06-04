# Trade-off Matrix — Durability Tiers

| Tier Decision | Benefit | Main Cost | Key Metric |
|---------------|---------|-----------|------------|
| 3 full replicas in one region | simple reads and repairs | high storage overhead | under-replicated age |
| erasure-coded cold tier | better cost efficiency at scale | slower rebuild and more coordination | repair completion time |
| geo-redundant archive | stronger disaster durability | higher egress and control complexity | regional restore success |
| fast repair SLO | shorter vulnerability window | more background capacity reserved | mean time to restored redundancy |
| customer-selectable tiers | business flexibility | support and policy complexity | invalid tier selection rate |
