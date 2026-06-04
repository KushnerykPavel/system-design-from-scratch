# Secrets, Key Management, and Rotation

> Secrets are not secure because they exist. They are secure only if they can be issued, stored, used, rotated, and revoked safely.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design secret distribution and key rotation so production systems can survive credential leaks, certificate expiry, and routine rotation without outages.
**Prerequisites:** `04-apis-contracts-and-schema-evolution/06-contract-testing`, `10-reliability-retries-and-backpressure/03-circuit-breakers`, `12-security-abuse-and-multitenancy/01-auth-and-trust`
**Estimated time:** ~75 min
**Primary artifact:** rotation checklist + interview card

## The Problem

Design secret and key management for a platform with APIs, background workers, databases, and inter-service mTLS.

Senior answers should cover more than "use a secret manager." They should explain:

- which data is a secret versus a key versus configuration
- how workloads fetch or derive credentials
- how rotation happens without fleet-wide restart or outage
- how revocation, expiry, and auditability work during incidents

## Clarify

- Are we rotating API secrets, database credentials, TLS certificates, envelope-encryption keys, or all of them?
- Can workloads call a central secret store at request time, or only at startup or refresh intervals?
- What outage is more dangerous: secret-store unavailability or stale credentials?
- How quickly must the system react if a credential leaks?

Assume short-lived workload identity where possible, dynamic secret delivery, and dual-read or overlap windows during rotation.

## Requirements

### Functional

- Issue and deliver secrets or short-lived credentials to workloads.
- Rotate credentials without manual host-by-host changes.
- Revoke compromised material and audit its use.

### Non-functional

- Avoid turning the secret store into a request-path dependency.
- Support routine rotation with low operational drama.
- Keep blast radius small for leaked or stale credentials.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Workloads | 20K service instances | drives refresh fanout and cache strategy |
| Rotation cadence | daily certs, weekly service creds, monthly key policy review | determines automation pressure |
| Secret fetch latency target | startup or background refresh, not hot path | avoids central dependency on each request |
| Emergency revocation | under 15 minutes | constrains TTL and distribution design |
| Rough cost | secret manager, audit storage, and rollout tooling | rotation quality is an operational investment |

## Architecture

```text
workload identity
  -> secret / key manager
     -> short-lived credential or wrapped key
     -> local in-memory cache
     -> dependency connection
     -> audit stream
```

Recommended shape:

1. Authenticate the workload with machine identity.
2. Issue scoped, time-bounded credentials or wrapped data keys.
3. Refresh in the background before expiry.
4. Use overlap windows or trust bundles so rotation is not a cliff-edge event.

## Data Model & APIs

Core entities:

- `SecretLease`
- `CertificateBundle`
- `KeyVersion`
- `RotationPolicy`
- `RevocationEvent`

Useful APIs:

- `IssueCredential(workload, scope, ttl)`
- `RotateSecret(name, nextVersion)`
- `RevokeVersion(id)`
- `ListActiveVersions(resource)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| all clients reload the same rotated secret at once | thundering herd or secret-store saturation | jittered refresh and cached overlap windows |
| secret-store outage blocks healthy services | startup failures or connection refresh failures spike | prefetch, local cache, and bounded grace period |
| leaked credential stays valid too long | suspicious use after compromise report | short TTL plus targeted revocation and re-issue |
| certificate rotation breaks trust chain | handshake failures after bundle rollout | dual trust bundles and staged validation |

## Observability

- metric: credential issuance rate, refresh latency, and refresh failures
- metric: active versions, near-expiry workloads, and revocation propagation lag
- metric: handshake failures after cert or key rotation
- log: secret access by workload identity and scope
- trace: dependency connection setup with credential-refresh annotations
- SLO: routine rotation should complete without broad workload restart or sustained auth failures

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| short-lived credentials | smaller compromise window | more refresh traffic and automation | long-lived static secrets on disk |
| background refresh with overlap | safer rotation | more state management | rotate only during deploys |
| central manager plus local cache | strong control and auditability | cache invalidation and grace-period design | secrets baked into images |

## Interview It

**Google framing:** "How would you manage service credentials and key rotation across a large production fleet?" Expect questions about workload identity, rollout safety, and revocation speed.

**Cloudflare framing:** "How would you rotate secrets and certificates across globally distributed edge systems?" Expect focus on partial rollout, stale nodes, and trust-bundle overlap.

**Follow-ups:**
1. What secrets should never be fetched on the hot path?
2. How do you rotate certificates without breaking old peers immediately?
3. What do you do if the secret manager is degraded during an incident?
4. How fast can you contain leaked credentials?
5. What changes at 10x fleet size?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/rotation-checklist.md`
- `outputs/interview-card-secrets-and-keys.md`

## Exercises

1. **Easy** — List three credentials in a typical backend and explain which should be short-lived first.
2. **Medium** — Design zero-downtime certificate rotation for service-to-service mTLS.
3. **Hard** — Redesign the system for edge nodes that may be partially disconnected from the control plane for hours.

## Further Reading

- [Secret Manager patterns](https://cloud.google.com/architecture/security-foundations/using-secrets) — practical secret-handling guidance
- [The twelve-factor app config](https://12factor.net/config) — useful baseline, but incomplete without rotation and revocation
