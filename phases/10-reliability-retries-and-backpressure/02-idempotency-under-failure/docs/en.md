# Idempotency Under Partial Failure

> "Probably succeeded" is a dangerous response unless the system can prove what happened.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design idempotent write paths that stay safe under retries, ambiguous timeouts, and partial commits instead of assuming the transport layer will tell the full truth.
**Prerequisites:** `04-apis-contracts-and-schema-evolution/02-idempotency-keys`, `08-consistency-replication-and-transactions/05-transactions`, `10-reliability-retries-and-backpressure/01-timeouts-and-retries`
**Estimated time:** ~75 min
**Primary artifact:** idempotency policy checker + failure checklist

## The Problem

Partial failure is the reason idempotency matters. The caller may see:

- timeout after the server committed
- connection reset after a side effect started
- duplicate delivery from an at-least-once queue
- repeated button taps from a user who never saw the first success

If the system cannot answer "did we already apply this intent?" then retries become correctness bugs, not availability improvements.

## Clarify

- Is the operation creating, mutating, or triggering an external side effect?
- Does the client provide an idempotency key, or must the server derive one?
- How long must duplicate suppression last?
- Does the system need to return the original response body on duplicate replay?

If the interviewer is vague, assume a public write API where clients may retry for minutes, a timeout can occur after commit, and users care more about correctness than raw throughput.

## Requirements

### Functional

- Accept repeated submissions of the same logical intent without double-applying it.
- Return a stable answer for duplicate requests inside the dedupe window.
- Preserve correctness across server restarts and ambiguous network failures.

### Non-functional

- Keep dedupe storage bounded and explainable.
- Avoid blocking unrelated writes due to coarse keys.
- Make duplicate suppression debuggable during incidents.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Write QPS | 35K req/s | determines dedupe-store write volume |
| Duplicate retry window | 24 hours | drives TTL, storage size, and cleanup |
| Duplicate ratio during incidents | 8-12% | partial failure creates bursts of repeated intents |
| Response payload retained | 1-4 KB per idempotency key | response replay can dominate storage cost |
| Rough cost | durable key store + response cache + cleanup | correctness has real storage and operational cost |

## Architecture

Good idempotency design ties the request intent to a durable decision record:

```text
client request + idempotency key
  -> request normalizer
  -> durable idempotency store
  -> business write / side effect
  -> stored outcome returned on duplicates
```

Common pattern:

1. Reserve the idempotency key before the side effect.
2. Store request fingerprint and status.
3. Apply the business operation once.
4. Persist final result or durable reference.
5. On retry, return the prior result if the intent matches.

## Data Model & APIs

Useful record shape:

- `IdempotencyRecord(key, scope, request_hash, status, response_ref, created_at, expires_at)`

Useful APIs:

- `BeginIntent(key, request_hash, scope)`
- `CompleteIntent(key, response_ref)`
- `LookupIntent(key)`

Senior-level detail:

- scope the key carefully, usually `tenant + client_key + operation`
- reject key reuse with a different request fingerprint
- make the dedupe store durable enough to survive restarts and leader changes

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| key reused for a different payload | mismatch between stored fingerprint and retry payload | store request hash and return conflict |
| key TTL expires before realistic retry window ends | late duplicate creates a second side effect | size TTL from client behavior and incident recovery windows |
| server writes side effect before reserving key | duplicate side effects after timeout or crash | reserve key first or couple key and write transactionally |
| duplicate suppression depends only on in-memory cache | restart causes replay of prior intent | durable idempotency store |

## Observability

- metric: duplicate hit ratio by endpoint and tenant
- metric: idempotency-key conflicts from mismatched payloads
- metric: stuck `in_progress` intent records and their age
- metric: dedupe-store latency and storage growth
- log: idempotency key, scope, request hash, and replay reason on duplicates
- trace: original attempt and duplicate attempts correlated to the same key
- SLO: duplicate retries for the same intent should never create duplicate external effects inside the supported window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| durable idempotency store | survives crashes and ambiguous timeouts | write amplification and cleanup work | in-memory dedupe only |
| store prior response reference | duplicate callers get a consistent answer | extra storage and serialization logic | only return "already processed" |
| long TTL on keys | safer during delayed retries | higher storage footprint | short TTL that misses real incident retries |

## Interview It

**Google framing:** "Design a payment-like write API that clients can safely retry." The signal is whether you can reason about commit-before-ack ambiguity and key scoping.

**Cloudflare framing:** "Design a control-plane write API for customer configuration changes." The signal is whether retries stay safe across globally distributed clients and occasional control-plane failures.

**Follow-ups:**
1. What changes if the side effect is external email or webhook delivery?
2. What if the same key is reused with a different payload?
3. How long should the idempotency record live?
4. Can the idempotency store become a bottleneck under incident retry bursts?
5. What is your rollback story for stuck `in_progress` records?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/design-review-idempotency-under-failure.md`
- `outputs/failure-checklist-idempotency-under-failure.md`

## Exercises

1. **Easy** — Define the key scope for a single-tenant order creation API.
2. **Medium** — Extend the design to return the original response body on duplicate replay.
3. **Hard** — Explain how you keep webhook delivery idempotent when the receiver also retries.

## Further Reading

- [Handling Duplicate Requests in Distributed Systems](https://stripe.com/blog/idempotency) — practical framing for idempotent APIs
- [Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — useful context for why retries and correctness must be designed together
