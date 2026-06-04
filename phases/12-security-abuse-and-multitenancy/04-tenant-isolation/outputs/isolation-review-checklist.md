# Isolation Review Checklist

## Data boundaries

- Are storage and cache keys tenant-scoped?
- Do async jobs preserve tenant context?
- Can support tooling bypass resource ownership checks?

## Performance boundaries

- Are there per-tenant quotas or budgets?
- Does worker scheduling prevent starvation?
- Can one tenant dominate one shard or queue?

## Escalation path

- Which tenants may need dedicated capacity?
- How will you detect that need early?
- How will you migrate without outage?
