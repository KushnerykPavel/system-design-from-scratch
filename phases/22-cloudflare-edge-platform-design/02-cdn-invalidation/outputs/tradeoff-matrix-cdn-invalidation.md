---
lesson: 02-cdn-invalidation
---

## Purge Scope Trade-offs

| Scope | Benefit | Risk | Best use |
|-------|---------|------|----------|
| Object | precise and cheap | many calls for bulk updates | individual asset correction |
| Prefix | easier bulk targeting | broader blast radius | deployment path or content subtree |
| Tag | flexible grouping | metadata and cardinality cost | product catalogs or campaign assets |
| Full zone | simple for customer | very expensive and disruptive | emergency rollback only |
