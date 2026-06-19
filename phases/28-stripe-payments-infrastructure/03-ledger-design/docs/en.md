# Ledger Design & Reconciliation
> Every debit has a matching credit — if they don't match, something is wrong.

**Type:** Build
**Company focus:** Stripe
**Learning goal:** Design an immutable double-entry ledger that supports real-time balance queries, multi-currency accounts, and nightly reconciliation against card network settlement files.
**Prerequisites:** `19-payments-wallets-and-ordering/02-digital-wallet`, `05-storage-indexing-and-access-patterns/06-time-series`
**Estimated time:** ~90 min
**Primary artifact:** ledger schema + reconciliation job design

---

## Double-Entry Bookkeeping

Every financial transaction is recorded as **two or more entries**: a debit on one account and a credit on another. The invariant is:

```
sum(debits) == sum(credits)  for every transaction
```

This means:
- Entries are **immutable** — they are never updated or deleted.
- A mistake is corrected by a new reversing entry, not by modifying the original.
- The ledger is an **append-only audit log** of every balance change.

---

## Ledger Entry Schema

```sql
CREATE TABLE ledger_entries (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id     UUID        NOT NULL,
    currency       CHAR(3)     NOT NULL,           -- ISO 4217 (USD, EUR, GBP)
    amount_cents   BIGINT      NOT NULL CHECK (amount_cents > 0),
    direction      TEXT        NOT NULL CHECK (direction IN ('DEBIT', 'CREDIT')),
    transaction_id UUID        NOT NULL,           -- groups the debit+credit pair
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    metadata       JSONB
);

CREATE INDEX idx_ledger_account_currency ON ledger_entries (account_id, currency, created_at DESC);
CREATE INDEX idx_ledger_transaction ON ledger_entries (transaction_id);
```

Key design decisions:
- `amount_cents` is always positive — direction is encoded in the `direction` column, not by sign.
- `transaction_id` groups all entries that belong to the same double-entry transaction.
- `metadata` stores context (payment_intent_id, fee breakdown, etc.) without schema changes.

---

## Balance Computation

Balance is not stored as a field. It is computed on read:

```sql
SELECT
    SUM(CASE WHEN direction = 'CREDIT' THEN amount_cents ELSE 0 END) -
    SUM(CASE WHEN direction = 'DEBIT'  THEN amount_cents ELSE 0 END) AS balance_cents
FROM ledger_entries
WHERE account_id = $1
  AND currency   = $2;
```

For high-read workloads, a **materialized balance** can be maintained:
- A `account_balances` table stores `(account_id, currency, balance_cents, version)`.
- On each ledger insert, the balance row is updated with optimistic locking (`WHERE version = $old_version`).
- Reads hit the balance table; the ledger is the source of truth for reconciliation.

---

## Multi-Currency

Each currency has a **separate logical ledger**. There is no implicit conversion:
- `(account_id=A, currency=USD)` and `(account_id=A, currency=EUR)` are independent balances.
- Currency conversion is itself a transaction: debit USD account, credit EUR account, debit/credit a FX-spread account.
- This prevents silent conversion errors and makes FX revenue explicit in the ledger.

---

## Reconciliation Job

Nightly (or near-real-time), a batch job downloads settlement files from card networks and banks, then matches them against the internal ledger:

```
1. Download settlement file (CSV, SWIFT MT940, ISO 20022 camt.053)
2. Parse each line into: (external_ref, amount_cents, currency, settled_at)
3. For each external line:
   a. Look up matching ledger entry by external_ref / transaction_id
   b. Compare amount and currency
   c. Mark as MATCHED or flag as UNMATCHED
4. Produce reconciliation report:
   - MATCHED: count, total amount
   - MISSING_INTERNAL: charge exists in bank file but not in ledger (revenue leak)
   - MISSING_EXTERNAL: charge in ledger but not in bank file (processing error)
   - AMOUNT_MISMATCH: found in both but amounts differ (fee error, FX drift)
   - CURRENCY_MISMATCH: found in both but currencies differ (conversion error)
```

**Stripe Sigma** exposes the ledger as SQL-queryable data, allowing merchants to run their own reconciliation queries against Stripe's data warehouse.

---

## Failure Modes

| Failure | Cause | Fix |
|---|---|---|
| Balance drift under concurrent writes | Two processes read balance, both apply delta, one overwrites the other | Serializable isolation or optimistic lock on `account_balances` |
| Reconciliation false positive | Timezone mismatch between Stripe timestamps (UTC) and bank settlement timestamps (local) | Normalize all timestamps to UTC before comparison; use date windows not exact timestamps |
| Bulk file parse error blocks reconciliation | Malformed line in settlement CSV aborts the job | Parse line-by-line with per-line error handling; flag bad lines, continue processing rest |
| Reversals not matched | Refund appears as a credit in bank file but ledger has debit-then-credit pair | Reconciliation logic must handle net settlement: match on (transaction_id, refund_id) |

---

## Trade-offs

- **Computed balance vs cached balance:** Computing from entries is always correct but O(n) per query. A cached balance is O(1) but requires careful consistency management to avoid drift. Stripe uses the cached approach with the immutable ledger as the reconciliation source of truth.
- **Append-only vs updateable entries:** Append-only entries enable audit trails and simplify replication (no UPDATE conflicts), but corrections require compensating entries which increase row count over time. Stripe partitions old entries to cold storage.
- **Nightly vs real-time reconciliation:** Nightly is simpler but surfaces discrepancies with up to 24h lag. Real-time reconciliation requires streaming settlement data (not all networks support this) but catches errors faster.
