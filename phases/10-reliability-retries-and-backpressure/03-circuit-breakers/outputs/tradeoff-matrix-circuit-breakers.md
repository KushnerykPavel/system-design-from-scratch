---
name: circuit-breaker-tradeoff-matrix
phase: 10
lesson: 03
---

| Choice | Benefit | Cost | Best fit |
|--------|---------|------|----------|
| stale cache fallback | preserves availability | weaker freshness | optional read features |
| hard fail when open | clear correctness boundary | more visible user errors | auth, billing, security-critical paths |
| small half-open probe set | safe recovery | slower confirmation | fragile or costly dependencies |
| per-dependency breaker | smaller blast radius | more policy tuning | mixed dependency criticality |
