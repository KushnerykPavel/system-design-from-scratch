# Interview Card — Authentication, Authorization, and Trust Boundaries

## Strong answer shape

- Define caller types first: user, service, operator, worker.
- Separate authentication from authorization explicitly.
- Name each trust boundary: edge, internal service, async worker, admin plane.
- Re-authorize at the resource-owning service for tenant-sensitive operations.
- Explain identity propagation, auditing, and break-glass handling.

## High-signal phrases

- "The gateway authenticates, but the resource owner authorizes."
- "Async boundaries need actor and tenant context, not just object IDs."
- "Trust should shrink as requests move inward, not expand."

## Common misses

- Assuming internal traffic is automatically trusted.
- Treating RBAC as enough without resource ownership checks.
- Forgetting admin-plane blast radius and auditability.
- Ignoring how background jobs preserve identity context.
