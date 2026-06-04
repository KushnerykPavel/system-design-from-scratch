# Abuse Layer Worksheet

## Edge

- What can be dropped cheaply?
- Which signals are safe to use at line rate?
- What should never hit origin before a coarse decision?

## Identity-aware controls

- Which paths need per-account protection?
- Which paths need per-API-key or per-tenant budgets?
- Where is challenge better than hard block?

## Cost-aware controls

- Which routes are expensive enough to deserve special limits?
- What is the fallback when challenge infrastructure is degraded?
- How will you explain a false positive?
