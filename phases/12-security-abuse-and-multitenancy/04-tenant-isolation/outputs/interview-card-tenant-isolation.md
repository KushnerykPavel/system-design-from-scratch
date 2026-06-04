# Interview Card — Tenant Isolation and Blast Radius

## Strong answer shape

- Talk about data, compute, cache, queue, and control-plane isolation.
- Add per-tenant budgets and fairness, not just authz.
- Plan for heavyweight tenant carve-outs.
- Name how support and admin flows stay scoped.

## Common misses

- Treating multitenancy as only a database-schema question.
- No noisy-neighbor mitigation.
- Forgetting async jobs and caches.
- No path for regulated or massive tenants.
