# Stripe Interview Rubric & Strong Signals

> Payments infrastructure interviews reward precision, not enthusiasm. Know the exact failure mode before you name the fix.

**Type:** Concept  
**Company focus:** Stripe  
**Learning goal:** Understand what Stripe evaluates in a system design loop — idempotency correctness, ledger consistency, developer API quality, fraud hooks, and global money movement. Recognize the difference between a strong-hire answer and a competent generic answer.  
**Prerequisites:** `19-payments-wallets-and-ordering/01-payment-fundamentals`, `08-consistency-replication-and-transactions/02-distributed-transactions`  
**Estimated time:** ~75 min  
**Primary artifact:** rubric card + strong-hire signal checklist  

## The Problem

You are preparing for a Stripe senior/staff system design interview. Stripe's bar is high and domain-specific. Interviewers are engineers who built Payment Intents, Radar, Connect, or Treasury. Generic distributed-systems answers will not land — they want to see that you understand Stripe's operational reality: every mutating API call must be idempotent, every money movement must balance a double-entry ledger, and every failure must leave the system in a recoverable state, never a half-charged one.

Stripe processes payments in 135+ currencies across 46+ countries, handles ~1M API requests per minute at peak, stores ~250M cards on file, and operates at 99.99%+ uptime. A candidate who cannot size the problem or reason about its failure modes precisely will not pass the bar.

## Clarify

- Which product surface is the question targeting? Payment Intents, Connect multi-party flows, Radar fraud, Treasury (banking-as-a-service), Issuing, or Sigma?
- Is the focus on the charge path (online payments), the payout path (payouts to connected accounts), or the reconciliation path (ledger correctness)?
- Is the interviewer probing breadth (end-to-end system sketch) or depth (idempotency implementation, exactly-once myths)?
- What is the scale point of interest — startup volume or Stripe's current peak?

## Requirements

### What Stripe Tests

Stripe interviewers assess several dimensions simultaneously:

1. **Idempotency correctness** — Every mutating API call must have an idempotency key. Keys are scoped to `(account_id, key_string)`. The system stores the original response and replays it verbatim on retry. Mismatched request body with the same key returns a 400. The dedup window is 24 hours.
2. **Double-entry ledger discipline** — Every credit has a matching debit. Postings are immutable and append-only. Balance never goes negative without explicit overdraft. Reconciliation runs to catch drift between ledger sums and processor statement balances.
3. **Exactly-once myth awareness** — Stripe uses at-least-once delivery with idempotency to *approximate* exactly-once semantics. There is no true exactly-once in a distributed system. Candidates who claim otherwise reveal a gap.
4. **Failure atomicity** — A charge either fully completes or fully reverses. No half-charged states. Every payment workflow has a compensating transaction path.
5. **PCI DSS compliance** — Cardholder data is isolated in a vault (Stripe's own or third-party tokenizer). PANs never touch application code or logs. Audit logging is mandatory on every cardholder data access.
6. **Developer API quality** — Stripe's API is its product. Versioned HTTP APIs, webhook event delivery (at-least-once with signed payloads), idempotency on POST, resource-based REST semantics, and rich error types.

### Non-functional Signals Stripe Cares About

- Consistency over availability for financial transactions — a stale balance is worse than a momentary 503.
- Latency at p99, not just mean — Stripe targets <200ms p99 for charge creation.
- Auditability — every state transition on a payment is logged with timestamp, actor, and old/new state.
- Regional isolation — Stripe operates in multiple regions; a region failure must not corrupt global ledger state.
- Webhook reliability — retries up to 72 hours (3 days) with exponential backoff; endpoint signing (HMAC-SHA256, `Stripe-Signature` header with timestamp tolerance).

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| API requests/min | ~1M peak | drives rate limiting, gateway, and API server sizing |
| Cards on file | ~250M | drives vault storage and tokenization throughput |
| Currencies supported | 135+ | drives FX conversion, rounding rules, and settlement partitioning |
| Countries | 46+ | drives compliance, tax, and banking partner topology |
| Uptime SLO | 99.99% | <53 min downtime/year; drives redundancy and failover design |
| Idempotency key dedup store | billions of keys/year | drives TTL-keyed storage (Redis or distributed KV) |
| Ledger entries/day | ~hundreds of millions | drives append-only storage sizing and reconciliation batch design |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- "Use Stripe's SDK and call `/v1/charges`." The candidate is describing Stripe's *external* API, not designing the internals.
- Stores money amounts as `float64` — instant red flag. Stripe uses integer minor units (cents, pence, yen).
- Says "add a queue to handle duplicates" without defining idempotency key semantics.
- No mention of ledger, reconciliation, or failure compensation.
- Treats webhooks as fire-and-forget.

### Strong-Hire Answer Pattern

- Starts with clarifying questions and a capacity estimate before touching architecture.
- Explicitly states that money amounts are `int64` minor units (never float).
- Describes idempotency key scoping: `(account_id, key)`, 24h TTL, response stored and replayed.
- Names the charge state machine: `requires_payment_method → requires_confirmation → processing → succeeded/failed`.
- Identifies double-entry ledger requirement; posts debit+credit atomically.
- Describes compensating transactions for failures; no half-charged states.
- Mentions Radar for fraud (ML + rule engine + 3DS fallback), even if not asked.
- Discusses API versioning strategy (Stripe adds fields, never removes; date-based versions).
- Names at least two Stripe-specific technologies: Sorbet, Veneur, Goose, Hydroplane, Temporal-like workflows, PostgreSQL for ledger.

### Strong-Hire Milestone Map

| Time mark | Expected progress |
|-----------|-------------------|
| 5 min | Clarified scope, named the charge path vs payout path, gave capacity numbers |
| 15 min | Sketched high-level: API gateway, charge service, idempotency store, ledger, card vault |
| 25 min | Drilled into idempotency — key scoping, dedup window, replay semantics, mismatch 400 |
| 35 min | Addressed failure atomicity, compensating transactions, and ledger reconciliation |
| 45 min | Summarized trade-offs, named observability signals, discussed API versioning |

## Stripe-Specific Vocabulary

| Term | What it means |
|------|---------------|
| Payment Intent | Stripe's object representing a payment lifecycle; replaces deprecated Charges API |
| Sources API | Deprecated; replaced by Payment Methods + Payment Intents |
| Connect | Stripe's platform product for multi-party payments (marketplace, SaaS billing) |
| Radar | Stripe's fraud detection — ML models + merchant-defined rules + 3D Secure |
| Treasury | Stripe's banking-as-a-service product (financial accounts, money movement) |
| Sigma | Stripe's SQL analytics product over merchant transaction data |
| Atlas | Stripe's business incorporation product |
| Terminal | Stripe's in-person payments (hardware + SDK) |
| Sorbet | Stripe's internal Ruby gradual type checker |
| Veneur | Stripe's open-source metrics aggregation daemon (StatsD-compatible) |
| Goose | Stripe's distributed locking library |
| Hydroplane | Stripe's database migration tool |

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Duplicate charge from client retry | Idempotency key lookup on every POST | Replay stored response; never re-execute |
| Card processor timeout | Timeout after 10s; no 2xx received | Mark charge `processing`; poll processor for terminal state |
| Ledger posting partial write | Atomic transaction on ledger DB | If commit fails, abort; retry will re-enter with idempotency key |
| Webhook delivery failure | Endpoint returns non-2xx or times out | Retry with exponential backoff for up to 72 hours |
| Region failover corrupts ledger | Dual-write or single-region writer with read replicas | Ledger writes go to single-region PostgreSQL; reads fan out to replicas |
| PAN leaks in logs | Automated log scanning for card number patterns | Tokenize at ingestion; PANs never enter application layer |

## Observability

- `stripe.charge.created` counter with currency and country dimensions
- `stripe.idempotency.replay_rate` gauge — rising replay rate indicates client retry storm
- `stripe.ledger.reconciliation_drift_cents` — any nonzero value pages on-call
- `stripe.webhook.delivery_latency_p99` — SLO: <30s for first attempt
- `stripe.radar.fraud_block_rate` — watch for sudden spikes (attack) or drops (rule misconfiguration)
- Distributed traces on every charge creation, tagged with `payment_intent_id`
- SLO: 99.99% successful charge creation; error budget tracked via error rate on `/v1/payment_intents`

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Integer minor units for amounts | No floating-point rounding errors on money | More complex display formatting | `float64` — causes $0.001 drift at scale |
| At-least-once + idempotency | Simple delivery guarantees; works with standard queues | Requires idempotency key discipline on every caller | Exactly-once — impossible in distributed systems |
| Single-region ledger writer | Avoids split-brain ledger corruption | Adds latency for cross-region reads | Multi-master ledger — two concurrent writers risk double-credit |
| Append-only ledger postings | Full audit trail; reconciliation catches any drift | Storage grows unbounded; compaction needed | Mutable balance row — no audit trail, hard to reconcile |
| Date-based API versioning | Merchants pinned to known behavior | Version proliferation over years | Semantic versioning — hard to communicate behavioral changes |

## Interview It

**Stripe framing:** "Design the payment processing API for Stripe. Start anywhere." Common pushbacks: "Your idempotency store is a single Redis — how do you handle its failure?" / "How do you guarantee the ledger is correct after a charge processor timeout?" / "What happens if the same webhook fires twice to the same endpoint?"

**Follow-ups:**
1. How does your design handle a charge that succeeds at the card network but your service crashes before writing the ledger entry?
2. Walk me through the idempotency semantics when a client sends the same key with a different request body.
3. How would you design the refund path so that a partial refund cannot exceed the original charge amount?
4. How does Connect change your design — specifically for a marketplace where the platform and the seller both need ledger entries?
5. How would you instrument the charge path so you could detect a processor degradation within 60 seconds?

## Ship It

After this lesson, you should be able to:
- Describe the charge lifecycle state machine from `requires_payment_method` to `succeeded`.
- Implement a simple idempotency key store with replay semantics in under 50 lines of Go.
- Explain double-entry ledger invariants and why append-only is required.
- Name five Stripe-specific technologies and describe their role.

## Exercises

1. Draw the Payment Intent state machine with all terminal states and transitions.
2. Write the SQL schema for a minimal double-entry ledger (accounts, transactions, postings) using integer amounts.
3. Design the idempotency key table: what columns does it need, what index, what TTL strategy?
4. Sketch the webhook delivery retry state machine from `pending` to `delivered` or `abandoned`.

## Further Reading

- Stripe API documentation: Payment Intents guide (stripe.com/docs/payments/payment-intents)
- "Designing robust and predictable APIs with idempotency" (stripe.com/blog/idempotency)
- "Advanced distributed systems design" — Kyle Kingsbury (Jepsen) on exactly-once myths
- Sorbet: gradual typing for Ruby (sorbet.org)
- Veneur: distributed metrics aggregation (github.com/stripe/veneur)
