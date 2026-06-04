# Interview Card — Secrets, Key Management, and Rotation

## Strong answer shape

- Separate secrets, certificates, and encryption keys.
- Use workload identity to fetch scoped short-lived credentials.
- Keep the secret manager off the request hot path.
- Design overlap windows, revocation, and staged rollout.
- Explain how the system behaves if the secret manager is degraded.

## Common misses

- "Use Vault" as the whole answer.
- Rotating only during deploys.
- No revocation path for leaked material.
- No monitoring for near-expiry or stale-version use.
