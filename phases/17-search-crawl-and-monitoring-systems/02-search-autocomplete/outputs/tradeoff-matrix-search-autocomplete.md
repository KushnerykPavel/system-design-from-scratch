---
lesson: 02-search-autocomplete
---

| Option | Strength | Weakness | Best when |
|--------|----------|----------|-----------|
| Snapshot-only serving | very low, stable latency | freshness bounded by publish cadence | query popularity moves in minutes, not seconds |
| Hybrid snapshot + live trend features | better trending relevance | more serving-path dependencies | trends matter and some feature lag is acceptable |
| Heavy personalization | higher user-specific relevance | memory and feature-fetch complexity | product value strongly depends on user history |
| Global shared suggestions only | simpler rollout and policy review | weaker local relevance | policy safety and operational simplicity dominate |
