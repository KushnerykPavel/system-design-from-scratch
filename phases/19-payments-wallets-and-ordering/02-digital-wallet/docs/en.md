# Digital Wallet with Holds and Settlements

> Wallets fail when "available" and "committed" money are treated as the same thing.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a wallet system with holds, expirations, settlements, and release flows that prevents double-spend under retries and delayed downstream confirmation.
**Prerequisites:** `19-payments-wallets-and-ordering/01-payment-ledger`, `08-consistency-replication-and-transactions/06-sagas`, `10-reliability-retries-and-backpressure/01-timeouts-and-retries`
**Estimated time:** ~75 min
**Primary artifact:** wallet-plan validator + hold-lifecycle checklist

## The Problem

Design a digital wallet used for stored value, promo credits, or internal prepaid balances. Users can add funds, place holds before an order is finalized, settle the hold when the order succeeds, and release it on cancellation or timeout.

This lesson matters because many interview answers collapse wallet design into a single balance number. Senior answers distinguish available funds from reserved funds, make hold expiry explicit, and show how wallet state stays safe when the order service or payment service is slow.

## Clarify

- Is the wallet system the monetary source of truth or a derived balance layer on top of a ledger?
- Can users spend across currencies, or is each wallet single-currency?
- Are holds short-lived authorization-style reserves, or long-running marketplace escrow holds?
- What happens if the downstream order service never confirms completion?

If the interviewer stays broad, assume single-currency wallets backed by the phase ledger, short-lived holds, idempotent wallet operations, and asynchronous settlement or release after the order path finishes.

## Requirements

### Functional

- Support wallet credit, debit, hold creation, hold release, and hold settlement.
- Prevent spending more than the available balance.
- Expire abandoned holds safely.
- Expose wallet activity history for customer support and reconciliation.
- Allow downstream order or payment workflows to retry safely.

### Non-functional

- Keep the synchronous hold path low latency.
- Avoid orphaned holds that permanently reduce spendable balance.
- Make settlement and release idempotent.
- Preserve auditability across manual intervention.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Wallet operations | 250K/s peak | shapes partitioning and hot-account strategy |
| Balance reads | 3M/s peak | drives serving cache and snapshot design |
| Hold TTL range | 5 minutes to 24 hours | determines cleanup and stuck-fund behavior |
| Active holds | 900M globally | affects storage, indexing, and expiry scans |
| Peak factor | 6x at campaign launch | hot wallets and promotional bursts must not overspend |

## Architecture

```text
wallet API
  -> idempotency store
  -> balance/hold service
  -> ledger posting service
  -> hold-expiry sweeper
  -> events to order, risk, and support systems
```

Design notes:

1. Represent wallet state with at least `available`, `held`, and `posted` semantics.
2. Make holds first-class objects with TTL, reason code, and downstream reference.
3. Expire or release abandoned holds asynchronously, but keep those actions auditable.
4. Keep overspend protection at the authoritative write path, not only in caches.

## Data Model & APIs

Core entities:

```text
wallet(wallet_id, owner_id, currency, status)
wallet_hold(hold_id, wallet_id, reference_id, amount_minor, expires_at, state)
wallet_activity(activity_id, wallet_id, type, amount_minor, hold_id, ledger_batch_id)
wallet_snapshot(wallet_id, available_minor, held_minor, posted_minor, version)
```

Useful interfaces:

- `CreateHold(wallet_id, amount_minor, reference_id, ttl, idempotency_key)`
- `ReleaseHold(hold_id, reason_code, idempotency_key)`
- `SettleHold(hold_id, final_amount_minor, idempotency_key)`
- `CreditWallet(wallet_id, amount_minor, source_ref, idempotency_key)`
- `GetWalletBalance(wallet_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| duplicate settle or release retry | idempotency conflicts and repeated downstream reference use | operation-scoped idempotency keys and hold state machine |
| hold never resolved after order timeout | aged active holds and hold-expiry backlog | sweeper plus operator tooling for stuck holds |
| hot wallet receives many concurrent requests | optimistic-lock conflicts and latency spikes | per-wallet serialization or shard-local sequencing |
| cache says funds are available when source disagrees | read/write version mismatch metrics | authoritative write checks and projection invalidation |

## Observability

- metric: hold-create latency, settle latency, and release latency
- metric: active hold age percentiles and expired-hold backlog
- metric: insufficient-funds rejection rate and optimistic-lock conflict rate
- log: manual release, forced settlement, and hold override actions with actor identity
- trace: order request -> wallet hold -> settlement or release -> ledger posting
- SLO: 99.9% of valid hold requests complete within target latency without overspend

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit hold objects | clear reservation lifecycle | more storage and cleanup logic | only tracking one mutable balance |
| async hold expiry sweeper | practical cleanup for abandoned flows | temporary stale reservations | blocking request path on remote order state |
| authoritative per-wallet concurrency control | prevents overspend | can create hot-account bottlenecks | eventually consistent balance checks |

## Interview It

**Google framing:** "Design a prepaid wallet for a commerce platform where orders place temporary holds before final settlement." Expect pushback on overspend protection, stale holds, and what happens during retries.

**Cloudflare framing:** "Design customer credit wallets for usage billing and promotional funds." Expect pressure on global scale, tenant isolation, and how ledger truth and serving balances interact.

**Follow-ups:**
1. What changes for marketplace escrow holds that last days instead of minutes?
2. How do you support partial settlement against a larger hold?
3. What if one enterprise wallet receives extreme concurrent debits?
4. How do you rebuild wallet balances after projection corruption?
5. What if risk review must freeze settlement but not balance reads?

## Ship It

- `outputs/failure-checklist-digital-wallet.md`

## Exercises

1. **Easy** — Model the state transitions for create hold, release hold, and settle hold.
2. **Medium** — Explain how to handle a partial settle after inventory shrinks.
3. **Hard** — Redesign the wallet for cross-region writes where the same user can spend from multiple edges.

## Further Reading

- [Stripe authorization and capture](https://docs.stripe.com/payments/place-a-hold-on-a-payment-method) — useful mental model for reserve then settle flows
- [Martin Kleppmann on transactions](https://martin.kleppmann.com/2015/09/26/transactions-at-stripe-scale.html) — good framing for practical correctness
