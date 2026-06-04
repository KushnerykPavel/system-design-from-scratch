# Idempotency Keys and Safe Retries

> A retry without idempotency is often just a duplicate side effect with better branding.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design write APIs that can survive client timeouts, gateway retries, and partial failures without creating duplicate work or conflicting state.  
**Prerequisites:** `04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events`, `03-design-framework-and-timing/07-interviewer-moves`  
**Estimated time:** ~75 min  
**Primary artifact:** failure checklist + design review prompt  

## The Problem

Many strong-looking designs still fail a basic production test: what happens if the client times out after the server already committed the write and retries the request?

Payment creation, order submission, quota reservation, and provisioning APIs all need a duplicate-suppression plan. This lesson teaches the design pattern and, equally important, its limits.

## Clarify

- What operation must be safe to repeat: create, charge, reserve, or mutate?
- Who generates the idempotency key: client, gateway, or server?
- How long must the system remember keys?
- What counts as "the same request" if headers or metadata differ?

## Requirements

### Functional

- Accept retried write requests without duplicating the side effect.
- Return the original result when a repeated request is recognized.
- Reject key reuse when the payload no longer matches.

### Non-functional

- Keep duplicate detection fast on the write path.
- Bound storage growth for remembered keys.
- Make it debuggable when retries, conflicts, or stale keys appear.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Write QPS | 20K create requests/s | drives lookup and persistence volume |
| Retry rate | 0.5% normal, 8% during incidents | duplicate suppression load spikes during failures |
| Key retention | 24 hours | defines storage footprint and cleanup policy |
| Payload hash size | 32 bytes | lets you detect conflicting key reuse cheaply |
| Rough cost | hot store plus TTL cleanup | key retention is not free at large scale |

## Architecture

Minimal safe flow:

1. Client sends `Idempotency-Key`.
2. Service checks key store.
3. If key exists with same payload hash, return stored result.
4. If key exists with different payload hash, reject.
5. If key is new, execute the side effect and atomically store outcome.

```text
client -> API -> idempotency lookup -> business write -> stored response
```

The subtle point: the idempotency record and the side effect must align. A duplicate-protection record written too early or too late creates new failure modes.

## Data Model & APIs

Example record:

```text
idempotency_key -> {
  request_hash,
  status,
  response_ref,
  created_at,
  expires_at
}
```

API shape:

- `POST /payments` with `Idempotency-Key` header
- conflict response if the key is reused with a different semantic request

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| write committed but result not returned | repeated key with existing success record | return original stored outcome |
| key stored before side effect succeeds | high rate of stuck pending records | use transactional write or reconciler |
| same key reused for different payload | request hash mismatch metrics | reject with explicit conflict |
| retention too short | duplicates reappear after TTL expiry | align TTL with realistic retry windows and user behavior |

## Observability

- metric: idempotency hit rate
- metric: payload-hash conflict count
- metric: pending record age
- metric: duplicate write attempts during incident windows
- log: key, request hash, and outcome class for sampled failures
- SLO: duplicate side effects remain below an explicit threshold during retries

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| client-generated keys | stable across retries through proxies | shifts responsibility to clients | server guesses duplicates heuristically |
| storing prior response | simple repeat behavior | storage overhead and privacy review | re-running the write path |
| bounded TTL retention | manageable storage growth | very late retries may miss protection | infinite retention |

## Interview It

**Google framing:** "Design a create-order API that stays safe under retries." The signal is whether you notice partial failure windows and not just happy-path deduping.

**Cloudflare framing:** "Protect an API from duplicate writes when gateways retry aggressively under origin timeouts." The signal is whether you think through retry amplification and storage pressure.

**Follow-ups:**
1. What changes if the client cannot generate keys reliably?
2. What if the first attempt is still in flight when the retry arrives?
3. What if duplicate protection must span regions?
4. What if the response body is too large to store with every key?
5. How do you expire old keys without creating a thundering cleanup job?

## Ship It

- `outputs/failure-checklist-idempotency-keys.md`
- `outputs/design-review-idempotency-keys.md`

## Exercises

1. **Easy** — Design the minimum idempotency record for a `POST /orders` API.  
2. **Medium** — Explain why a database unique constraint alone is not always enough.  
3. **Hard** — Redesign the pattern for an event-driven write that may be processed more than once downstream.  

## Further Reading

- [Stripe idempotent requests](https://stripe.com/docs/idempotency) — a practical public API example  
- [RFC 7231](https://www.rfc-editor.org/rfc/rfc7231) — useful background on HTTP method semantics  
