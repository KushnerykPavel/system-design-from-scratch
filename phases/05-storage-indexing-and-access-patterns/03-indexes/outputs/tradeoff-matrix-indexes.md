| Choice | Benefit | Cost | Use when |
|--------|---------|------|----------|
| Primary index only | cheapest writes | weak secondary query support | one lookup dominates |
| Targeted secondary indexes | fast important reads | write amplification | a few bounded list paths matter |
| Covering index | avoids base-row reads | more storage and maintenance | one hot query needs minimal hops |
| External search/index | flexible query surface | more systems and lag | many filters or text search matter |
