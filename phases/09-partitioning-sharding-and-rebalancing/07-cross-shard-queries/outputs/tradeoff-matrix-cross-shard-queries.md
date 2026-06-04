---
name: cross-shard-query-tradeoff-matrix
phase: 09
lesson: 07
---

| Query style | Best benefit | Main cost | Best use |
|-------------|--------------|-----------|----------|
| shard-local query | cheapest latency path | product scope must stay local | core user-serving reads |
| bounded scatter-gather | simple for low-rate global reads | merge latency and fanout cost | admin or moderate-rate control queries |
| materialized view | fast common global reads | lag and pipeline complexity | dashboards and rollups |
| search / analytics backend | flexible filtering and exploration | more systems and sync work | exploratory or broad global queries |
