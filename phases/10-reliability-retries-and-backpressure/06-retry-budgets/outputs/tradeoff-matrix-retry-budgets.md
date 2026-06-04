---
name: retry-budget-tradeoff-matrix
phase: 10
lesson: 06
---

| Choice | Benefit | Cost | Best fit |
|--------|---------|------|----------|
| explicit retry budget | bounded overload risk | policy tuning effort | services with tight capacity headroom |
| delayed hedging | better p99 without doubling all traffic | some tails still miss rescue | read-heavy latency-sensitive paths |
| aggressive cancellation | less wasted backend work | needs end-to-end support | RPC systems with cancel semantics |
| disable hedging on saturation | protects dependencies | less tail improvement during incidents | overload-prone environments |
