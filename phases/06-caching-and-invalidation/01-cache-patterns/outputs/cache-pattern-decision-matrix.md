# Cache Pattern Decision Matrix

| Pattern | Best fit | Main benefit | Main risk | Interview cue |
|---------|----------|--------------|-----------|---------------|
| Read-through | read-heavy serving path | simple miss amortization | hot misses can stampede origin | "cache on read, DB stays authoritative" |
| Write-through | read-after-write sensitive path | fresher post-write reads | higher write latency and partial-failure logic | "pay on writes to tighten freshness" |
| Write-behind | bursty async-friendly writes | batching and smoothing | delayed durability or lost async work | "only when delayed persistence is acceptable" |

## Quick questions

- What is the source of truth?
- Which path is hotter: reads or writes?
- What is the user-visible cost of stale reads?
- Can delayed persistence be tolerated at all?
- What failure becomes most dangerous if this pattern is chosen?
