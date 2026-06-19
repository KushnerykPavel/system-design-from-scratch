# Global Money Movement & Multi-Currency

> Money moves through a chain of correspondent banks — design for failures at every link.

**Type:** Build  
**Company focus:** Stripe  
**Learning goal:** Design Stripe's global money movement infrastructure that handles payouts in 135+ currencies across 46+ countries without currency conversion errors.  
**Prerequisites:** `19-payments-wallets-and-ordering/01-payment-ledger`, `08-consistency-replication-and-transactions/06-sagas`  
**Estimated time:** ~75 min  
**Primary artifact:** multi-currency ledger design + payout routing diagram  

## The Problem

Stripe processes $640B+ in total payment volume per year across 135+ currencies in 46+ countries. Every transaction involves currency conversion decisions, banking partner routing, and regulatory compliance checks. A single wrong floating-point operation can cause the ledger to be off by one cent — at Stripe's scale, that drifts into millions of dollars of reconciliation failures per year.

The payout side is where money physically leaves Stripe and reaches merchants' bank accounts. The path goes through correspondent banking networks (SWIFT, ACH, SEPA, Faster Payments) that each have their own failure modes, settlement windows, and regulatory requirements.

## Clarify

- Is the question focused on charge currency conversion, payout routing, or the multi-currency balance ledger?
- What payout speed tiers does the merchant expect? (Instant minutes, Standard 1-3 days, Wire same/next day)
- Which banking networks are in scope? (ACH/US, SEPA/EU, Faster Payments/UK, SWIFT/international)
- Does the merchant hold multi-currency balances, or does Stripe convert everything to a single settlement currency?
- What is the AML/KYC requirement level for this payout type?

## Requirements

### Money Representation: Never Use Float

**The cardinal rule of financial systems: store money as integers in the smallest currency unit.**

| Currency | Smallest unit | ISO 4217 | Example: $10.00 |
|----------|--------------|----------|-----------------|
| USD | cent (1/100) | USD | 1000 |
| EUR | euro cent | EUR | 1000 |
| GBP | penny | GBP | 1000 |
| JPY | yen (no subunit) | JPY | 1000 (= ¥1000) |
| BHD | fils (1/1000) | BHD | 10000 (= 10.000 BHD) |
| KWD | fils (1/1000) | KWD | 10000 |

Zero-decimal currencies (JPY, KRW, VND, etc.): `1` in the API = one full unit (one yen).  
Three-decimal currencies (BHD, KWD, OMR): `1` in the API = 0.001 of the unit.  
Two-decimal currencies: most common (USD, EUR, GBP).

`float64` cannot represent most decimal fractions exactly. `0.1 + 0.2 = 0.30000000000000004` in IEEE 754. At 1M transactions/day, this causes ledger drift. Use `int64` always.

### Payout Types

| Type | Network | Speed | Cost | Use case |
|------|---------|-------|------|---------|
| Instant Payout | Visa Direct / MC Send (card push) | Minutes | ~1.5% + $0.25 | Urgency payouts, gig economy |
| Standard Payout | ACH (US), SEPA (EU), Faster Payments (UK) | 1-3 business days | Low flat fee | Default merchant payout |
| Wire | SWIFT | Same/next day | $10-25 fee | Large amounts, international |

### Local vs Cross-Border Acquiring

Stripe prefers local acquiring (Stripe entity in the same country as the cardholder) because:
- Lower interchange fees: ~0.1% for local vs ~1-3% for cross-border
- Better approval rates (issuers trust local acquirers more)
- Avoids FX conversion on the acquiring side

When Stripe has no local entity in a country, it routes through a cross-border acquiring partner and charges the merchant the FX fee.

### Currency Accounts and FX Conversion

**Single settlement currency (default):** Stripe converts all charges to merchant's settlement currency (usually USD) at the point of charge. The FX rate applied is: mid-market rate + Stripe's spread (typically 2%).

**Multi-currency balances (Stripe Balance product):** Merchant holds separate balances per currency. Example: a merchant might hold USD, EUR, GBP separately and only convert when they choose. Avoids conversion fees for merchants with natural multi-currency revenue and expenses.

**FX rate stale risk:** Rates are updated from market data providers and cached with a maximum TTL of 30 seconds. A stale rate means the merchant pays a different spread than expected — the drift must be reconciled in the FX settlement file.

### Correspondent Banking Chain

```
Merchant's customer pays (card charge)
           │
           ▼
Stripe acquires the transaction
           │
           ▼
Card network settles funds to Stripe (T+2 for most cards)
           │
           ▼
Stripe credits merchant's balance (internal ledger entry)
           │
           ▼
Merchant requests payout
           │
           ▼
Stripe's bank submits payout instruction:
  US domestic → ACH via Stripe's bank (Wells Fargo / JPMorgan)
  EU → SEPA Credit Transfer via Stripe's EU banking license
  UK → Faster Payments via Stripe's UK FCA authorization
  International → SWIFT wire via correspondent banking network
           │
           ▼
Correspondent bank(s) relay funds
           │
           ▼
Recipient bank credits merchant account (T+1 to T+3)
```

### AML / KYC Checks Before Every Payout

Every payout must pass:
1. **OFAC SDN list check** — US Treasury sanctions list; cannot pay sanctioned entities
2. **Bank account verification** — Stripe verifies bank accounts via:
   - Micro-deposits (send 2 small amounts, merchant confirms values)
   - Plaid integration (instant OAuth-based account verification)
3. **KYC verification** — Identity verification for merchants above certain volume thresholds (required by banking partners and regulations)

### Idempotent Payout Design

Payouts must be idempotent:
- Each payout has a stable payout ID
- The banking instruction includes a unique reference that the bank can deduplicate
- If the banking partner ACKs the same reference twice, Stripe ignores the second ACK
- SWIFT payments include a UETR (Unique End-to-end Transaction Reference) for tracking

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Total payment volume | $640B/year | Drives FX exposure and banking partner relationships |
| Currencies | 135+ | Drives multi-currency ledger complexity and FX rate management |
| Payout countries | 46+ | Drives banking partner count and regulatory compliance |
| FX rate cache TTL | 30s max | Stale rates cause reconciliation drift |
| ACH batch windows | 3 per day | US ACH batches submit at fixed times; missed window = +1 day delay |
| SWIFT settlement | 1-3 business days | International wires slower than domestic |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- Stores money amounts as `float64` — fundamental correctness failure
- Treats all currencies the same (ignores zero-decimal and three-decimal currencies)
- Says "just use a currency conversion API" without addressing rate staleness or reconciliation
- Treats payout as a simple API call without mentioning correspondent banking, idempotency, or failure modes
- No mention of OFAC / AML requirements

### Strong-Hire Answer Pattern

- Immediately states int64 minor units, explains zero-decimal currencies (JPY), three-decimal (BHD)
- Describes FX conversion at charge time vs payout time trade-off
- Names ACH/SEPA/Faster Payments/SWIFT by name with their settlement windows
- Explains correspondent banking chain and where failures can occur
- Designs idempotent payout with UETR reference for SWIFT
- Mentions OFAC check as a hard requirement on every payout

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| FX rate staleness | Monitor cache age; alert at 25s | Refresh from market data provider; fail-safe to last known rate with wider spread |
| ACH file rejected by bank | ACH return code received (R01-R99) | Parse return code; map to user-facing error; retry with verified account |
| SWIFT wire delayed by correspondent | SWIFT MT199 status message not received | Idempotent retry via alternative correspondent; track with UETR |
| Payout routing failure (no banking partner) | No route found for currency+country | Fallback to USD wire if merchant accepts; otherwise fail with clear error |
| OFAC false positive | Legitimate merchant flagged | Manual review queue; compliance team clears within 24h |
| Duplicate payout instruction | Same payout ID submitted twice to bank | Bank deduplicates on unique reference; Stripe idempotency layer prevents double-submission |

## Observability

- `stripe.payout.latency_p99` by network (ACH, SEPA, SWIFT) — SLO: ACH <2 days, SEPA <1 day, Instant <5 min
- `stripe.fx.rate_age_seconds` — alert if > 25s; page if > 30s
- `stripe.payout.failure_rate` by return code — R02 (account closed) vs R09 (insufficient funds)
- `stripe.ledger.balance_drift_cents` — reconciliation job populates; any nonzero value pages on-call
- `stripe.payout.ofac_hold_count` — rising count may indicate data quality issue with name matching
- Distributed traces on every payout instruction submission, tagged with `payout_id` and `network`

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| int64 minor units | Zero floating-point rounding errors | Requires explicit decimal display logic per currency | float64 — 0.1+0.2 drift; catastrophic at scale |
| Convert at charge time (default) | Merchant knows exact settlement amount at sale time | Merchant absorbs FX risk; cannot benefit from rate moves | Convert at payout — merchant gets rate uncertainty for days |
| Multi-currency balances | Merchants avoid conversion fees for natural hedges | More complex ledger (one balance per currency per merchant) | Single currency — cheaper to operate; merchant loses FX flexibility |
| Local acquiring over cross-border | 0.1% vs 1-3% savings on interchange | Requires Stripe to incorporate legal entities in each country | Single global acquiring — simpler but 10-30x higher FX costs |
| UETR for SWIFT | ISO 20022 standard; enables end-to-end tracking | Only useful for SWIFT; each network has its own reference scheme | Proprietary Stripe reference — not usable for external tracing |

## Interview It

**Stripe framing:** "Design Stripe's payout infrastructure. Stripe pays out to merchants in 135 currencies across 46 countries. Walk me through the design — specifically how money moves from a charge settlement to a merchant's bank account in Germany."

**Follow-ups:**
1. A merchant in Germany wants to receive EUR payouts. Walk me through how Stripe routes the SEPA Credit Transfer.
2. How do you handle a Japanese yen payout where the currency has no subunit? What does `1` mean in your API?
3. An ACH return code R02 (account closed) comes back 3 business days after payout. What does your system do?
4. How do you ensure a payout is not sent twice if your payout service crashes after submitting to the bank but before recording the confirmation?
5. What does FX rate staleness risk look like at Stripe's scale, and how do you mitigate it?

## Ship It

After this lesson, you should be able to:
- Explain why float64 is forbidden for money and how int64 minor units solve it.
- Name the four major payout networks (ACH, SEPA, Faster Payments, SWIFT) with their settlement windows.
- Explain the difference between zero-decimal and two-decimal currencies with examples.
- Describe the OFAC check requirement and where it sits in the payout flow.
- Design an idempotent payout with a bank-level unique reference.

## Exercises

1. Write the SQL schema for a multi-currency balance ledger. Each merchant can hold balances in multiple currencies.
2. Implement a safe integer rounding function for FX conversion that rounds to the nearest minor unit (no float64 in the rounding logic).
3. Design the payout state machine from `created` to `paid` or `failed`. What states exist for an ACH payout?
4. Sketch the reconciliation job that compares Stripe's internal ledger against ACH settlement files. What mismatches can occur?

## Further Reading

- Stripe Payouts documentation: stripe.com/docs/payouts
- ISO 4217 currency codes: iso.org/iso-4217-currency-codes.html
- ACH return codes: nacha.org/ach-network/return-codes
- SWIFT UETR: swift.com/swift-resource/248576/download
- "Falsehoods programmers believe about money" (github.com/kdeldycke/awesome-falsehood)
