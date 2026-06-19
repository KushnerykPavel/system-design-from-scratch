# Stripe Full Mock Loop

> Strong-hire answers show you understand payments as a regulated financial system, not just a software architecture problem.

**Type:** Mock Interview  
**Company focus:** Stripe  
**Learning goal:** Demonstrate end-to-end Stripe system design proficiency in a 45-minute mock interview covering marketplace payments, Connect architecture, multi-party ledger, and fraud/compliance hooks.  
**Prerequisites:** All prior lessons in Phase 28 (01–06), `19-payments-wallets-and-ordering/01-payment-ledger`  
**Estimated time:** 45 min (mock) + 30 min (debrief)  
**Primary artifact:** Scored mock interview with annotated strong/weak signals  

## The Prompt

> "Design Stripe's payment processing API for a marketplace like Airbnb or DoorDash, where funds must be split between the platform and sellers. The marketplace processes 1 million transactions per day with an average order of $50. Payouts to sellers happen on a T+2 schedule."

This is the **Stripe Connect** use case. The platform (Airbnb) collects from the buyer, Stripe holds the funds, and routes the split between the platform's account and the seller's (connected) account.

Do not look at the design below until you have attempted the mock. Set a 45-minute timer.

---

## Timeline Milestones

| Time | Expected progress |
|------|-------------------|
| 5 min | Clarified scope — Connect vs basic charges, split mechanics, KYC requirements for sellers, payout schedule |
| 15 min | High-level design — PaymentIntent creation, Stripe Connect topology (platform + connected accounts), split routing |
| 25 min | Deep dive — idempotency on every operation, double-entry ledger entries for split charge, async payout design |
| 35 min | Failure recovery + compliance — fraud hooks, OFAC on payout, chargeback propagation to split accounts |
| 45 min | Observability, trade-offs summary, API versioning mention |

---

## The Problem

### Scale

- 1M transactions/day = ~11.6 transactions/second average
- Average order: $50 → $50M/day in Gross Merchandise Volume (GMV)
- T+2 payouts: funds collected on Monday pay out to sellers on Wednesday
- Connected accounts (sellers): could be thousands to millions of individual sellers

### Stripe Connect Architecture

Stripe Connect has three account topologies. This marketplace use case uses **Destination Charges**:

```
Buyer's card charge → Platform Stripe account
                           │
                           ├── Platform keeps: application_fee_amount
                           │
                           └── destination: seller connected account receives the rest
```

**Destination charge flow:**
1. Platform creates a `PaymentIntent` with `transfer_data.destination = connected_account_id`
2. Stripe charges the buyer's card for the full amount
3. The full amount lands in the platform's balance first
4. Stripe immediately creates a `Transfer` to the connected account for `amount - application_fee`
5. The connected account's balance grows; they receive a payout on their own schedule

**Alternative: Separate Charges + Transfers** (more flexible but more complex):
- Platform charges the buyer's card separately
- Platform manually creates a Transfer to the connected account later
- Useful when the split amount is not known at charge time

### Double-Entry Ledger for Split Charge

For a $50 charge where the platform takes $5 (application fee) and the seller gets $45:

```
Ledger entries (all in USD cents):

CREDIT  buyer_card_funding_account          5000   (funds from card network)
DEBIT   platform_receivable                 5000

DEBIT   platform_receivable                 5000
CREDIT  platform_balance                     500   (application fee: $5.00)
CREDIT  seller_connected_balance            4500   (destination: $45.00)
```

Every credit has a matching debit. The ledger never goes out of balance. This is a hard invariant.

### Connected Account Onboarding (KYC)

Before a seller can receive payouts, Stripe requires KYC (Know Your Customer) verification:
- Name, date of birth, address, SSN last 4 (US) or national ID
- Business verification for business accounts
- Bank account for payout
- Stripe performs identity verification and OFAC SDN list check

KYC requirements vary by country. Stripe handles the complexity via Connect Onboarding (hosted UI) or embedded components (iframe). Platforms should not build their own KYC flow.

### Payout Scheduling (T+2)

Stripe holds funds in the platform or connected account's balance until the payout date:
- T = date of charge settlement (card network settles T+1 or T+2 typically)
- T+2 = platform's configured payout schedule (can be T+2, weekly, monthly, manual)
- Connected accounts can have their own payout schedule, different from the platform

Payout design:
- Idempotent payout with stable payout ID
- OFAC check on bank account before each payout
- ACH / SEPA / Faster Payments based on seller's country
- Payout state machine: `pending → in_transit → paid | failed`

### Chargeback Propagation

When a buyer disputes a charge:
1. Card network notifies Stripe with chargeback
2. Stripe debits the platform's balance for the full charge amount + dispute fee (~$15)
3. Platform is responsible for responding with evidence
4. Platform can optionally "claw back" the seller's portion from the connected account balance
5. If the dispute is won, the funds are restored

The platform must decide: absorb chargebacks as platform risk, or push them to sellers.

---

## Strong-Hire Indicators

A strong-hire candidate demonstrates all of the following:

**Idempotency throughout:**
- PaymentIntent creation is idempotent (client sends `Idempotency-Key` on POST)
- Transfer to connected account is idempotent
- Payout creation is idempotent
- Explicit statement: "I would add idempotency keys to every mutating operation"

**Correct ledger entries for split charge:**
- Can draw the debit/credit pair for a split charge
- States that `application_fee_amount + destination_amount = total_charge_amount` as a hard invariant
- Mentions append-only ledger with reconciliation

**Async payout design:**
- Payouts are not synchronous with the charge — they are scheduled
- Payout state machine explicitly named
- OFAC check before payout dispatch
- Retry on ACH/SEPA failure with appropriate return code handling

**Regulatory considerations:**
- KYC for connected account onboarding (do not build your own — use Stripe's Connect Onboarding)
- Money transmission licenses (Stripe holds them; platform does not need their own)
- OFAC on every payout

**Webhook design for split outcome notification:**
- `payment_intent.succeeded` → platform webhook
- `transfer.created` → connected account webhook
- `payout.paid` or `payout.failed` → connected account webhook
- `charge.dispute.created` → platform webhook for chargeback handling

---

## Weak-Hire Indicators

- Treating the split as a simple percentage split in application code (no mention of Stripe Connect)
- No mention of ledger entries — "we just call the Stripe API and record it in our DB"
- No idempotency on the Transfer or Payout creation
- Synchronous payout (waiting for the bank transfer to complete in-request)
- No KYC mention
- No OFAC / AML consideration
- Treating chargebacks as an edge case to handle later

---

## Scoring Rubric

| Dimension | Weight | What earns full marks |
|-----------|--------|----------------------|
| Idempotency design | 20 pts | Idempotency keys on PaymentIntent, Transfer, and Payout; replay semantics described |
| Ledger correctness | 20 pts | Double-entry for split charge; invariant: fee + destination = total; append-only postings |
| API design quality | 15 pts | Resource-based design; error format; expandable objects; versioning mention |
| Fraud/compliance hooks | 15 pts | Radar for charge scoring; OFAC on payout; KYC for connected account onboarding |
| Failure recovery | 15 pts | Transfer idempotency; payout retry on ACH failure; chargeback propagation design |
| Observability | 15 pts | Named metrics: charge success rate, payout failure rate, ledger reconciliation drift |

**Hire thresholds:**
- 85–100: strong-hire
- 65–84: hire
- 45–64: mixed
- < 45: no-hire

---

## Sample Strong Answer (Abbreviated)

"I'll start by clarifying scope. This is a Connect marketplace — platform collects from buyers, splits funds to sellers. I'll use destination charges: the platform creates a PaymentIntent with `transfer_data.destination` and `application_fee_amount`. Every operation gets an idempotency key: the PaymentIntent creation, the Transfer, and each Payout.

For the ledger, when a $50 charge splits as $5 platform fee + $45 seller: I post a debit to the card funding account and a credit to the platform receivable for the full $50. Then I post credits to platform balance ($5) and seller connected balance ($45) — the debits and credits always balance.

Payouts are async and scheduled. Sellers receive T+2 payouts via ACH or SEPA depending on their country. Before dispatching each payout, I run an OFAC SDN check. If ACH returns R02 (account closed), I fail the payout and notify the seller to update their banking details.

For connected account onboarding, I use Stripe Connect Onboarding — I do not build my own KYC flow. Stripe handles identity verification, bank account verification, and the compliance requirements by country.

For chargebacks, I design the platform to be notified via `charge.dispute.created` webhook. The platform responds with evidence and decides whether to claw back from the seller's balance.

Observability: I'd instrument `charge.success_rate`, `transfer.failure_rate`, `payout.failure_rate` by ACH return code, and a reconciliation job that checks `platform_balance + all_connected_balances = total_settled_from_card_networks`."

---

## Follow-Up Questions

1. Your seller's ACH payout fails with return code R16 (account frozen). What does your system do, and who is responsible for resolving it?
2. A buyer disputes a $200 charge that was split $20 platform / $180 seller. How does the chargeback propagate, and how does the platform decide whether to claw back the seller?
3. Walk me through the idempotency semantics if the Transfer creation call to the connected account times out before you receive a response. What happens on retry?
4. How would you design observability to detect when the total of all connected account balances diverges from what the card networks have settled to Stripe?
5. The platform wants to delay payouts for new sellers for 7 days to reduce fraud risk. How would you implement this without breaking Stripe's standard payout flow?

---

## Ship It

After this mock, you should be able to:
- Draw the destination charge flow with correct ledger entries for a split charge.
- Explain why idempotency is required on Transfer and Payout creation, not just PaymentIntent.
- Describe the connected account onboarding requirements and why you delegate KYC to Stripe.
- Name the four webhook events relevant to a split charge lifecycle.
- Score your own answer using the rubric before reviewing the debrief.

## Further Reading

- Stripe Connect documentation: stripe.com/docs/connect
- Destination charges: stripe.com/docs/connect/destination-charges
- Connect account types: stripe.com/docs/connect/accounts
- Stripe Connect onboarding: stripe.com/docs/connect/connect-embedded-components
- Chargeback handling: stripe.com/docs/disputes
