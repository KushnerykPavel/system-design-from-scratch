---
lesson: 02-distributed-rate-limiter
---

| Option | Strength | Weakness | Best when |
|--------|----------|----------|-----------|
| Local-only token bucket | fastest latency | weak global correctness | traffic is partitioned naturally |
| Shared central counter | simple mental model | hotspot and latency pressure | low scale or strong strictness |
| Local cache + shared store | balanced realism | bounded inconsistency | high-scale gateways |
