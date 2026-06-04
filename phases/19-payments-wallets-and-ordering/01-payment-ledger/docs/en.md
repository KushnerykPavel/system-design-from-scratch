# Payment Ledger

> In financial systems, the fastest way to lose trust is to make balance history impossible to reconstruct.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design an append-only payment ledger with idempotent writes, immutable auditability, and operationally credible reconciliation boundaries.
**Prerequisites:** `08-consistency-replication-and-transactions/05-transactions`, `10-reliability-retries-and-backpressure/02-idempotency-under-failure`, `18-messaging-and-job-platforms/06-exactly-once-myths`
**Estimated time:** ~90 min
**Primary artifact:** ledger-plan validator + consistency checklist

## The Problem

Design the ledger behind card charges, bank payouts, refunds, and internal adjustments. Product teams want fast payment state changes, finance wants exact books, and operators need to prove what happened after duplicates, retries, or downstream settlement drift.

This lesson matters because "store the current balance" is not enough in interviews. Strong answers explain the source of truth, posting rules, reconciliation boundaries, and how to survive ambiguous payment retries without inventing money.

## Clarify

- Is this a system of record for money movement, or a serving layer above an external processor?
- Do we need double-entry accounting, or only an auditable event stream with derived balances?
- What currencies, settlement windows, and correction workflows matter?
- Which writes must be synchronous before returning success?

If the interviewer stays broad, assume an internal ledger of record for platform balances, multi-currency metadata but same-currency postings per entry set, idempotent client retries, and asynchronous reconciliation with external processors.

## Requirements

### Functional

- Record every financial movement as immutable ledger entries.
- Guarantee that each posting batch is balanced before commit.
- Support idempotent retries for payment create, capture, refund, and reversal APIs.
- Derive account balances from postings or trusted snapshots plus postings.
- Support reconciliation jobs against processor and bank statements.

### Non-functional

- Never acknowledge a successful posting that cannot be recovered after crash.
- Preserve strict auditability for corrections and reversals.
- Keep balance lookups fast without making the cache the source of truth.
- Separate customer-visible payment status from ledger finality where necessary.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Posting requests | 120K/s peak during flash sales | drives partitioning, commit path, and idempotency storage |
| Entry fanout | 2 to 8 ledger lines per business action | determines actual write amplification |
| Retention | 7 years hot-queryable metadata, cold archive beyond | compliance and audit posture dominate storage planning |
| Reconciliation volume | 500M external settlement rows/day | makes batch validation and drift detection a first-class workflow |
| Peak factor | 4x retry amplification during processor incidents | retries must not create duplicate money movement |

## Architecture

```text
clients
  -> payment API with idempotency store
  -> posting service
  -> balanced-transaction validator
  -> append-only ledger log / ledger DB
  -> balance materializer and snapshots

external processor + bank reports
  -> reconciliation pipeline
  -> discrepancy queue
  -> finance operations tooling
```

Design notes:

1. Keep business workflows and ledger posting separate: an order may be pending while the ledger records an authorization hold or a reserve.
2. Model corrections as compensating entries, never in-place mutation.
3. Make idempotency keys scoped by operation and actor so retries are safe but accidental collisions stay visible.
4. Treat reconciliation as part of the product, not a cleanup job.

## Data Model & APIs

Core entities:

```text
account(account_id, owner_type, currency, status)
posting_batch(batch_id, idempotency_key, reference_id, state, created_at)
ledger_entry(entry_id, batch_id, account_id, direction, amount_minor, currency, code)
balance_snapshot(account_id, as_of_entry_id, available_minor, pending_minor)
reconciliation_item(source, external_ref, expected_amount, observed_amount, status)
```

Useful interfaces:

- `CreatePostingBatch(reference_id, idempotency_key, entries[])`
- `GetAccountBalance(account_id)`
- `ReversePostingBatch(batch_id, reason_code)`
- `ImportSettlementFile(source, file_id)`
- `ListDiscrepancies(status, age_bucket)`

Senior answers explicitly state that `sum(debits) == sum(credits)` within a committed batch, and that payment status APIs may read from serving projections while finance trusts the ledger.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| duplicate client retry posts the same charge twice | idempotency hit rate and duplicate reference alarms | durable idempotency store keyed before posting |
| batch commits with unbalanced lines | invariant-check failures and reconciliation drift | validate batch totals atomically before append |
| cached balance diverges from ledger | snapshot drift metrics and reconciliation jobs | rebuild projections from append-only source |
| settlement file disagrees with internal ledger | discrepancy counts by processor and age | hold funds, create investigation queue, compensate with explicit entries |

## Observability

- metric: posting commit latency, idempotency-hit rate, and failed-balance-invariant count
- metric: unreconciled dollars and discrepancy aging by processor
- metric: projection lag between latest ledger entry and serving balance snapshot
- log: every reversal, manual adjustment, and reconciliation override with operator identity
- trace: payment request -> posting batch -> ledger commit -> balance materialization
- SLO: 99.99% of accepted posting batches are durably committed and queryable within the ledger target window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| append-only double-entry ledger | strong auditability and correction model | more modeling overhead for product teams | mutable balance row as source of truth |
| async reconciliation with explicit discrepancy workflow | practical integration with external processors | temporary status mismatch risk | blocking every payment on processor finality |
| snapshots plus replay for balances | fast reads with recoverability | projection maintenance complexity | recomputing every balance from genesis on each read |

## Interview It

**Google framing:** "Design the ledger behind payments and refunds for a large commerce platform." Expect pushback on idempotency scope, correction semantics, and how finance proves exactness after outages.

**Cloudflare framing:** "Design a global billing ledger for usage charges, credits, and payouts." Expect pressure on multi-tenant isolation, audit trails, and operating safely when downstream settlement feeds are delayed.

**Follow-ups:**
1. What changes if the product expands to multi-currency settlements?
2. How do you handle partial captures and partial refunds?
3. What if a region requires customer data deletion but financial records retention?
4. What happens if projection rebuild falls hours behind?
5. How would you migrate from mutable balances to an append-only ledger?

## Ship It

- `outputs/consistency-checklist-payment-ledger.md`

## Exercises

1. **Easy** — Sketch the ledger entries for authorize, capture, refund, and chargeback on one order.
2. **Medium** — Explain how idempotency keys should differ for charge creation versus refund creation.
3. **Hard** — Redesign the ledger for cross-region active-active writes where some balance views can be stale but books cannot be corrupted.

## Further Reading

- [Modern Treasury ledger design overview](https://www.moderntreasury.com/journal/how-to-scale-a-ledger-part-i) — strong grounding on ledger invariants
- [Stripe idempotent requests](https://docs.stripe.com/api/idempotent_requests) — practical framing for retry-safe payment APIs
