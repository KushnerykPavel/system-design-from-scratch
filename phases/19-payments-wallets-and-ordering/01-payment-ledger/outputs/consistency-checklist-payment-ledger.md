---
lesson: 01-payment-ledger
focus: balanced
---

## Ledger invariants

- Every posting batch is balanced before commit.
- Money movement is recorded append-only.
- Corrections are modeled as compensating entries.
- Idempotency keys are scoped and durable.

## Read-path checks

- State which balance view is authoritative.
- Explain how snapshots are rebuilt from ledger entries.
- Say how stale serving balances are detected.

## Failure probes

- What happens after client retry during partial timeout?
- How is settlement drift detected and triaged?
- How do operators inspect manual adjustments?

## Redesign prompts

- Add multi-currency settlement.
- Add cross-region writes.
- Add processor outage with 4x retries.
