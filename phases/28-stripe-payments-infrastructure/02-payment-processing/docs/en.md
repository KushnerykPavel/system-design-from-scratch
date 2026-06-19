# Payment Processing & Idempotency
> A charge must succeed exactly once — network retries are a fact of life, double charges are not.

**Type:** Build
**Company focus:** Stripe
**Learning goal:** Design a payment processing API where every mutating operation is idempotent and partial failures result in a clean retry path.
**Prerequisites:** `04-apis-contracts-and-schema-evolution/02-idempotency-keys`, `08-consistency-replication-and-transactions/05-transactions`
**Estimated time:** ~90 min
**Primary artifact:** idempotency key design + payment flow diagram

---

## Payment Intents API

Stripe models a charge as a stateful object called a **Payment Intent**. The lifecycle is:

1. **Create intent** — client calls `POST /v1/payment_intents` with amount, currency, and payment method. Returns a `payment_intent_id` and `client_secret`.
2. **Confirm** — client confirms the intent (attaches payment method, triggers authorization). This is when the card network is contacted.
3. **Capture or decline** — for manual capture flows, a separate `POST /v1/payment_intents/:id/capture` moves funds. For automatic capture, confirmation and capture happen atomically.

This two-step model separates authorization (card network check, hold funds) from capture (actual charge), which is required for many merchant workflows (e.g., hotel pre-auth, marketplace payouts).

---

## Idempotency Keys

Every mutating Stripe API call accepts an `Idempotency-Key` header. The design:

- **Client-generated UUID** — the caller picks a key (usually a UUID v4) and stores it locally before sending the request.
- **Scoped to (Stripe-Account, endpoint)** — the key is unique within one account's operations on one endpoint. The same key on a different endpoint is a different operation.
- **Stored with request hash + response** — the server hashes the request body. On retry, if the hash matches, the cached response is returned. If the hash differs, a 422 is returned (the caller is reusing a key for a different operation).
- **24-hour dedup window** — after 24 hours the key expires and a new request with the same key is treated as a fresh operation.

### Why client-generated keys?

The client must generate the key *before* sending the request. If the server generated the key, a timeout on the first call would give the client no key to retry with. Client generation means the key survives network failures.

---

## Charge Flow

```
Client
  │
  ▼
Stripe API Gateway
  │  (validates auth, rate limits, deserializes)
  ▼
Payment Router
  │  (selects card network, applies Radar fraud score)
  ▼
Card Network (Visa / Mastercard)
  │  (routes to issuing bank)
  ▼
Issuing Bank
  │  (authorization decision: approve / decline / 3DS)
  ▼
Authorization Response → back to Stripe → back to Client
```

**Database:** PostgreSQL for the payment ledger. Strong consistency is required — authorization decisions read current balance, and a stale read could permit an overdraft.

**Optimistic locking on Payment Intent status transitions** — the `payment_intents` table has a `status` column and a `version` integer. Updates use `WHERE id = $1 AND version = $2` and fail if another process already advanced the state.

---

## Webhook Delivery

After a payment completes (or fails), Stripe emits events to merchant-configured webhook endpoints:

- Stripe sends `POST` with a signed JSON payload.
- Merchant must respond with HTTP 200 within 30 seconds.
- On failure (non-200, timeout), Stripe retries with **exponential backoff** up to **72 hours**.
- After 72 hours without ack, the event moves to a **dead-letter queue** for manual review.

Merchants should make their webhook handlers idempotent — Stripe may deliver the same event more than once.

---

## PCI Scope Isolation

Raw card numbers (PANs) never touch Stripe's main application servers. The flow:

1. Browser renders a Stripe-hosted iframe (`Stripe.js`).
2. Card number is typed into the iframe — it never enters the merchant's DOM.
3. Stripe.js tokenizes the card directly against Stripe's **Cardholder Data Environment (CDE)**, a network-isolated system with PCI DSS Level 1 certification.
4. Only the resulting token (`pm_xxx`) is returned to the merchant and passed to Stripe's API.

This design limits PCI scope: only the CDE handles raw PANs. The API tier, application databases, and merchant systems only ever see tokens.

---

## Failure Scenarios

| Scenario | Behavior |
|---|---|
| Gateway timeout on charge | Client retries with same idempotency key. Server returns cached result if charge completed, or processes once if not. |
| Bank decline | Terminal failure. No retry. Return `payment_intent.payment_failed` event with decline code. |
| Network partition during capture | Distributed lock prevents double-capture. Compensation job reverses any partial state after lock timeout. |
| Idempotency key hash mismatch | Return 422 — caller is using same key for different request body. |

---

## Trade-offs

- **Strong consistency vs latency:** PostgreSQL serializable isolation prevents overdrafts but adds ~5ms lock contention under high write load. Stripe accepts this cost for correctness.
- **Synchronous fraud scoring vs latency:** Radar runs synchronously in the charge critical path. This adds ~10-20ms but enables real-time block/3DS decisions before authorization is sent to the card network.
- **Webhook at-least-once vs exactly-once:** Stripe guarantees at-least-once delivery. Exactly-once would require two-phase commit with every merchant endpoint, which is not feasible. Merchants must handle deduplication.
