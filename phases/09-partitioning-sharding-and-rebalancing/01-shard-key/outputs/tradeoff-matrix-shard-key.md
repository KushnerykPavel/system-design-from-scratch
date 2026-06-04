---
name: shard-key-tradeoff-matrix
phase: 09
lesson: 01
---

| Candidate key | Strongest benefit | Main risk | Best fit |
|---------------|-------------------|-----------|----------|
| `tenant_id` | query locality and isolation | large-tenant skew | multi-tenant SaaS with tenant-scoped reads |
| `hashed_tenant_bucket` | smoother distribution | more indirection | uneven tenant sizes with move flexibility |
| `user_id` | balances user-owned writes | weak tenant-locality | consumer products with user-centric queries |
| `region + tenant_bucket` | legal placement and latency control | more shards and more ops | geo-sensitive enterprise systems |
| `random_object_id` | even write spread | poor workload alignment | only when queries are mostly point lookups |
