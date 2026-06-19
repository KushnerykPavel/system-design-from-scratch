# Developer API Design — Webhooks & Versioning

> Breaking changes break trust — version everything, deprecate slowly, and deliver events reliably.

**Type:** Concept  
**Company focus:** Stripe  
**Learning goal:** Understand how Stripe designs its developer-facing API to be stable, versioned, and self-documenting, with reliable webhook delivery as a first-class primitive.  
**Prerequisites:** `04-apis-contracts-and-schema-evolution/04-api-versioning`, `04-apis-contracts-and-schema-evolution/05-event-schema-evolution`  
**Estimated time:** ~75 min  
**Primary artifact:** API versioning strategy + webhook delivery design  

## The Problem

Stripe's API is its product. Millions of developers have integrated Stripe code into their applications. A breaking API change silently corrupts those integrations — payment flows stop working, charges fail, webhooks deliver mismatched schemas. Unlike typical software where a major version is a migration, a payment API breaking change can mean real money lost in production.

Stripe's API must be:
- **Stable**: a version pinned in 2018 must still work in 2026.
- **Evolvable**: Stripe must be able to add new payment methods, fraud signals, and features.
- **Reliable**: webhook events must be delivered even when merchant endpoints are temporarily down.
- **Secure**: webhook payloads must be authenticated so merchants can trust their source.

## Clarify

- Is the question about the external REST API design or the internal service API design?
- Which surface: Payment Intents, Connect, Webhooks, or the overall versioning strategy?
- What does "breaking change" mean precisely? (Removing a field? Changing a field type? Changing behavior?)
- How do we handle merchants who never upgrade their pinned version?

## Requirements

### API Versioning Strategy

**Date-based versions:** Stripe uses calendar date versions (e.g., `2023-10-16`). Each release date is a snapshot of the API contract at that point.

**Version pinning per API key:** When a merchant creates an API key or first authenticates, they are pinned to the current version. Their integration sees exactly that API version forever — field names, types, behavior — until they explicitly opt in to a newer version.

**Rules for compatibility:**
- Adding new fields to a response: always allowed, not a breaking change.
- Adding new event types: allowed; merchants receive events they didn't register for only if they subscribed to `*`.
- Removing a field: breaking change — requires new version, old version keeps the field.
- Changing a field type (e.g., string → object): breaking change — new version only.
- Changing behavior of an existing endpoint: breaking change — new version only.
- Stripe NEVER removes deprecated fields within a live version.

**Version header:** Merchants can override their pinned version per-request using the `Stripe-Version` header. This allows gradual migration without changing the key's pinned version.

**Deprecation lifecycle:**
1. Feature added in version `2023-10-16`.
2. Deprecated in version `2024-06-01` with migration guide published.
3. Removal announced 12+ months before it happens, with email notifications to affected API key holders.
4. Removal only happens in a new version — old versions keep the deprecated field.

### API Design Principles

**Noun-based resources:** Stripe models everything as resources: `PaymentIntent`, `Customer`, `Subscription`, `Invoice`, `Refund`. Actions are expressed as state transitions on these resources, not as verbs.

| Bad (RPC style) | Good (REST resource style) |
|-----------------|---------------------------|
| `POST /charge` | `POST /v1/payment_intents` |
| `POST /cancelCharge` | `POST /v1/payment_intents/:id/cancel` |
| `POST /refundCharge` | `POST /v1/refunds` |

**Consistent error format:**
```json
{
  "error": {
    "type": "card_error",
    "code": "card_declined",
    "message": "Your card was declined.",
    "param": "card",
    "decline_code": "insufficient_funds",
    "doc_url": "https://stripe.com/docs/error-codes/card-declined"
  }
}
```

Error `type` values: `api_error`, `card_error`, `idempotency_error`, `invalid_request_error`, `authentication_error`, `rate_limit_error`.

**Expandable objects:** By default, Stripe returns object IDs for related resources. Merchants can request inline expansion using `expand[]=`.

```
GET /v1/payment_intents/pi_xxx?expand[]=customer&expand[]=payment_method
```

This avoids N+1 API calls without requiring merchants to always receive large payloads.

**Idempotency-Key header:** Every POST endpoint accepts an `Idempotency-Key` header. The server stores the response for 24 hours and replays it on retry. This is mandatory, not optional — Stripe's SDKs automatically generate idempotency keys.

### Webhook Architecture

Webhooks are Stripe's mechanism for push-based event notification. When a payment succeeds, a subscription renews, or a dispute is filed, Stripe sends an HTTP POST to the merchant's registered endpoint.

**Event schema:**
```json
{
  "id": "evt_xxx",
  "type": "payment_intent.succeeded",
  "created": 1698412800,
  "data": {
    "object": {
      "id": "pi_xxx",
      "status": "succeeded",
      "amount": 2000,
      "currency": "usd"
    }
  },
  "api_version": "2023-10-16"
}
```

**Delivery guarantees:** At-least-once delivery. Merchants must handle duplicate events idempotently (using `evt_xxx` as the idempotency key).

**Signature verification (HMAC-SHA256):**
The `Stripe-Signature` header contains:
```
t=1698412800,v1=<hmac_hex>
```
Where `<hmac_hex>` = HMAC-SHA256(secret, `t=1698412800` + `.` + request_body).

Merchant must:
1. Extract timestamp `t` from the header.
2. Check that `|current_time - t| < 300s` (tolerance window) to prevent replay attacks.
3. Compute expected HMAC using the webhook signing secret.
4. Compare to received HMAC using constant-time comparison.

**Retry schedule:** Stripe retries webhook delivery with exponential backoff for up to 72 hours (3 days). After 3 days of failures, the endpoint is automatically disabled. Merchant is notified by email.

**Webhook registration:**
- Per-endpoint: merchant registers an HTTPS URL and selects event types.
- Per-account or Connect app: platform-level webhook listeners for Connect events.
- Test vs. live mode: separate webhook endpoints for each mode.

**Event versioning:** Webhook events are versioned using the merchant's webhook endpoint version (configurable separately from the API key version). This allows merchants to migrate webhook consumers independently from their API key.

**Stripe CLI for local testing:**
```bash
stripe listen --forward-to localhost:3000/webhook
stripe trigger payment_intent.succeeded
```
The CLI creates a temporary HTTPS tunnel, forwards events, and captures the signing secret for local development.

### SDK Responsibilities

Official Stripe SDKs (10+ languages) provide:
- Automatic idempotency key generation on POST requests
- Automatic retry with exponential backoff on network errors or 429/500
- Webhook signature verification helper
- Strong types for all API resources
- API version pinning via SDK configuration

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Webhook events/day | ~100M | Drives webhook delivery system throughput |
| Webhook retry attempts | Up to 8 per event | Drives queue depth and retry storage |
| Active API versions in production | ~6-10 simultaneously | Drives version compatibility layer complexity |
| Webhook delivery latency SLO | <30s for first attempt | Drives webhook worker fleet sizing |
| HMAC verification overhead | ~0.1ms per payload | Negligible; do not skip it |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- "Just use API keys and version the URL as `/v1/`, `/v2/`." — URL versioning is incompatible with Stripe's per-key pinning model.
- Says webhooks are "fire and forget" — misses at-least-once delivery guarantee.
- Does not mention HMAC signature verification — the merchant has no way to verify event authenticity.
- Proposes semantic versioning (v1, v2, v3) — creates cliff-edge migrations instead of gradual ones.
- Does not mention the replay attack protection (timestamp tolerance window).

### Strong-Hire Answer Pattern

- Explains date-based versioning pinned to API key creation; distinguishes from URL versioning.
- Names exactly what a breaking change is (field removal, type change, behavior change) vs non-breaking (field addition, new event type).
- Designs webhook delivery as a queue with exponential backoff, dead-letter after 72h.
- Explains HMAC-SHA256 signature with timestamp tolerance window to prevent replay attacks.
- Mentions at-least-once delivery and that merchants must handle idempotency on their side.
- Notes that webhook event version is configurable separately from API key version.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Merchant endpoint returns 500 | Non-2xx response code | Retry with exponential backoff; 72h window; dead-letter after |
| Replay attack on webhook | Timestamp `t` > 300s old | Reject payload; log and alert |
| HMAC mismatch (tampered payload) | Computed HMAC ≠ received HMAC | Reject payload with 400; do not process event |
| Endpoint disabled after 3 days | Consecutive failures exceed window | Email merchant; auto-disable endpoint; event log retained for 30 days |
| SDK auto-retry causes duplicate processing | Same event delivered twice | Merchant deduplicates on `evt.id`; at-least-once is documented guarantee |
| Deprecated field removed too early | Merchants still using field | Monitor API field usage analytics; never remove until usage = 0 |

## Observability

- `stripe.webhook.delivery_latency_p99` — SLO: <30s for first attempt
- `stripe.webhook.retry_queue_depth` — rising depth indicates broad endpoint outage or attack
- `stripe.webhook.failure_rate` by endpoint — merchant-specific failure spikes
- `stripe.api.version_distribution` — track what fraction of requests use each pinned version
- `stripe.api.deprecated_field_usage` — detect merchants still using deprecated fields before removal
- `stripe.idempotency.replay_rate` — rising replay rate indicates client retry storm
- Alert on HMAC mismatch rate > 0.01% — may indicate signing key compromise

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Date-based versioning per API key | Merchant integrations never break; API evolves freely | Version proliferation over years; compatibility layer complexity grows | URL versioning (/v1, /v2) — cliff-edge migrations; merchants must change all URLs |
| At-least-once webhook delivery | Simple delivery guarantee; recovers from transient failures | Merchants must handle duplicate events idempotently | Exactly-once — impossible in distributed systems without 2PC overhead |
| HMAC-SHA256 with timestamp tolerance | Cryptographic authenticity + replay prevention in one header | Merchant must implement signature verification (SDK handles this) | No signature — any caller can forge a payment.succeeded event |
| Exponential backoff for 72h | Handles extended merchant maintenance windows | Events can be 3 days old by final delivery | Fixed interval retry — thundering herd on endpoint recovery |
| Separate webhook version from API version | Independent migration paths for server and event consumers | Two version settings to manage | Single version for both — forces synchronized migration of API + webhook consumers |

## Interview It

**Stripe framing:** "Design Stripe's webhook delivery system. Stripe sends 100M webhook events per day. How do you ensure reliable, secure, at-least-once delivery?"

**Follow-ups:**
1. A merchant's endpoint goes down for 6 hours during maintenance. What happens to the webhooks Stripe tried to deliver during that window?
2. How would you design the HMAC signature verification to be resistant to replay attacks?
3. A merchant reports they are receiving duplicate `payment_intent.succeeded` events. Is this a bug or expected behavior?
4. How do you handle a merchant who is still using a deprecated field that Stripe wants to remove in the next version?
5. How would you design the Stripe CLI's local webhook forwarding without requiring the merchant to set up a public HTTPS endpoint?

## Ship It

After this lesson, you should be able to:
- Explain date-based API versioning and why it is superior to URL-based major versions for SDK-heavy APIs.
- List the five fields in a Stripe error response and their purpose.
- Describe the HMAC-SHA256 webhook signature with timestamp tolerance and why both are needed.
- Explain at-least-once webhook delivery and what merchants must do on their side.
- Describe the exponential backoff retry schedule and what happens after 72 hours.

## Exercises

1. Write the HMAC-SHA256 signature verification in Go, including timestamp tolerance check.
2. Design the webhook delivery queue schema: what fields does each queue entry need?
3. Sketch the webhook endpoint state machine: `active` → `disabled` and the recovery path.
4. Design the API version compatibility layer: how would you serve a 2019-pinned merchant the correct response schema for a resource that has changed 4 times since then?

## Further Reading

- Stripe API versioning documentation: stripe.com/docs/api/versioning
- Stripe webhook documentation: stripe.com/docs/webhooks
- "Designing robust and predictable APIs with idempotency" (stripe.com/blog/idempotency)
- HMAC-SHA256: RFC 2104 (tools.ietf.org/html/rfc2104)
- Stripe API changelog: stripe.com/docs/upgrades
