---
lesson: 06-geo-consistency
---

| Workflow | Strongest useful semantic | Main cost | Common fallback |
|---------|---------------------------|-----------|-----------------|
| Profile read | bounded staleness | possible stale view | route to home region |
| User write | home-region write | write latency for distant users | queue and retry or fail closed |
| Control-plane change | strong or quorum replication | slower propagation | staged rollout with gates |
| Analytics read | eventual consistency | delayed freshness | tolerate lag |
