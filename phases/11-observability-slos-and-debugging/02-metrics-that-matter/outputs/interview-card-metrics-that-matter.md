# Interview Card — Metrics That Actually Explain the System

## Core categories

- user outcome: success and latency
- workload: QPS, concurrency, queue ingress
- saturation: pool exhaustion, queue age, utilization
- dependencies: downstream latency and errors

## Strong habits

- prefer metrics that narrow causes, not just describe activity
- use bounded labels by default
- separate sync and async path metrics
- mention tails, not only averages

## Typical mistake

"We’ll monitor CPU, memory, and QPS" is too generic unless you connect those signals to user impact and bottleneck hypotheses.
