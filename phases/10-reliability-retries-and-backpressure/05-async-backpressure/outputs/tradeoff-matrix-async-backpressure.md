---
name: async-backpressure-tradeoff-matrix
phase: 10
lesson: 05
---

| Choice | Benefit | Cost | Best fit |
|--------|---------|------|----------|
| bounded backlog | broker stays healthy | producers may be throttled | high-volume pipelines |
| queue-age dropping | preserves fresh value | some old work is lost | freshness-sensitive workloads |
| separate topics by class | stronger isolation | more operational overhead | mixed-value traffic |
| publish credits | clear producer feedback | more coupling and control logic | internal producer fleets |
