| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Time-bucket partitions | cheap retention and recent reads | cross-bucket complexity | one giant append table |
| Raw plus rollups | fast long-range queries | compaction and lag | querying raw points forever |
| Query window limits | protects serving path | less analyst flexibility | unlimited interactive scans |
