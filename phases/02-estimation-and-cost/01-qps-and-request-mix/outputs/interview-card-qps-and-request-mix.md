# Interview Card — QPS and Request Mix

## Default flow

1. Pick one traffic anchor.
2. Convert to average QPS.
3. Apply a visible peak factor.
4. Split into read, write, and background load.
5. Name the first bottleneck this creates.

## Fast prompts

- "I’ll estimate global and peak QPS first so the topology is grounded."
- "I’m assuming 80/20 read/write unless you want a write-heavy workload."
- "The number that matters most here is not average traffic, but burst-adjusted origin load."

## Common misses

- using DAU without requests per user
- sizing only average traffic
- forgetting internal fanout
- ignoring regional skew
