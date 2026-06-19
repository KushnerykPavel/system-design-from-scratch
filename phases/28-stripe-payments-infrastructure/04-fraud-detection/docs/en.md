# Fraud Detection Pipeline

> Fraud scoring must be faster than the authorization timeout — you have 100ms before the card network moves on.

**Type:** Build  
**Company focus:** Stripe  
**Learning goal:** Design Stripe's Radar fraud detection pipeline that scores transactions in <100ms without blocking the payment critical path.  
**Prerequisites:** `12-security-abuse-and-multitenancy/03-abuse-prevention`, `07-queues-streams-and-workflows/01-queues-vs-streams`  
**Estimated time:** ~75 min  
**Primary artifact:** fraud signal pipeline diagram + rule engine design  

## The Problem

Stripe Radar processes every transaction flowing through the Stripe network — approximately 1 million per minute at peak. For each transaction, Radar must produce a risk score and a recommended action (allow, 3D Secure, block) before the card network authorization timeout expires. That window is approximately 100ms. Miss it and the charge fails regardless of the fraud decision.

The challenge: fraud signals are distributed across many stores (Redis for velocity counters, ML model servers, device fingerprint databases, global network data), yet all must be assembled and scored within a strict real-time latency budget.

## Clarify

- Is the question focused on the synchronous scoring path, the async enrichment pipeline, or the feedback loop that retrains models?
- What latency budget does the card network impose? (~100ms p99 for synchronous scoring)
- What is the false positive tolerance? (Too many false positives destroy legitimate revenue; too few expose the merchant to fraud liability.)
- Are we designing Radar for Stripe's own merchant base or as an API for external integrators?
- Does the merchant have custom rule overrides on top of the ML baseline?

## Requirements

### Fraud Signal Taxonomy

Radar assembles signals across several categories before scoring:

**Velocity signals** (fetched from Redis, rolling windows):
- Transactions per card in last 1 hour / 24 hours / 7 days
- Transactions per IP address in last 1 hour
- Transactions per device fingerprint in last 1 hour
- Decline rate for this card in last 24 hours (high decline rate → high risk)

**Card BIN data** (fetched from local cache):
- Issuing bank and country
- Card type: prepaid vs credit vs debit (prepaid cards have higher fraud rates)
- Commercial vs consumer card

**Device fingerprint signals**:
- Browser/device fingerprint collected via Stripe.js
- Is this device new to Stripe's global network?
- Has this device been seen with fraudulent activity at other merchants?

**Geo mismatch signals**:
- Billing address country vs IP geolocation country
- Shipping address country vs billing address country (for physical goods)
- Card issuing country vs transaction country

**Merchant history**:
- Merchant's historical chargeback rate
- Merchant's industry risk category

**Global network signals** (Stripe's unique advantage):
- Has this card been marked fraudulent at any other Stripe merchant?
- Has this card been used successfully at many other Stripe merchants? (positive signal)
- Global card velocity across the entire Stripe network, not just this merchant

### Synchronous Scoring Path (must complete in <100ms)

```
PaymentIntent.confirm()
       │
       ▼
  Idempotency check (Redis, ~1ms)
       │
       ▼
  Radar signal assembly (parallel fetches, ~10-20ms)
  ┌────────────────────────────────────┐
  │  Velocity counters (Redis)         │
  │  Card BIN lookup (local cache)     │
  │  Device fingerprint (Redis/DB)     │
  │  Global network signals (DB)       │
  └────────────────────────────────────┘
       │
       ▼
  ML model inference (~20-40ms)
  - Gradient boosting or neural network
  - Input: assembled feature vector
  - Output: risk score 0.0–1.0
       │
       ▼
  Rule engine (merchant custom rules, ~1-2ms)
  - Merchant can add rules that override ML score
  - Examples: "block all prepaid cards", "always allow my own email domain"
       │
       ▼
  Risk action decision
  - score < 0.3 → allow
  - 0.3 ≤ score < 0.7 → 3D Secure redirect
  - score ≥ 0.7 → block
       │
       ▼
  Continue to card network authorization
```

### 3D Secure (3DS) Redirect

When score falls in the medium-risk band, Radar triggers a 3DS redirect:
- Cardholder authenticates directly with their issuing bank
- Shifts liability from merchant to card network on successful authentication
- Adds ~5-10 seconds to checkout UX
- Reduces false positive cost by allowing risky-but-legitimate transactions to proceed with extra auth

### Post-Decision Async Enrichment

After the synchronous decision, a background pipeline enriches transaction data for future model training:
- OFAC list check (US Treasury sanctions list) — legally required, async is acceptable
- Extended BIN database lookup
- Shipping address risk scoring
- IP reputation database cross-reference

This enrichment does NOT affect the current transaction decision but feeds the retraining pipeline.

### Feedback Loop: Disputed Charges → Fraud Labels

Chargebacks and manual dispute outcomes are the primary labels for model retraining:
- A chargeback marked "fraudulent" → negative training label for that transaction's feature vector
- A successfully resolved dispute → positive label
- Batch retraining on a regular cadence (daily or weekly)
- Online learning for fast-adapting attack patterns

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Transactions/min at peak | ~1M | Drives Radar throughput requirement |
| Scoring latency SLO | <100ms p99 | Hard constraint from card network authorization timeout |
| Velocity counter reads/tx | ~6 (card, IP, device, 3 time windows each) | Drives Redis QPS: ~6M/min = ~100K/s |
| ML model inference | ~20-40ms | Must fit in 100ms budget with signal assembly |
| False positive rate target | <1% | Every false positive is a legitimate revenue loss |
| Chargeback rate target | <0.1% | Visa/MC chargeback thresholds: 1% triggers fines |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- Scores transactions synchronously using only rule-based checks (no ML mention)
- Puts OFAC check on the synchronous path (OFAC lookups are slow; they block authorization)
- Does not separate the scoring latency budget from the authorization timeout
- Ignores false positive cost — treats any fraud block as a success
- No feedback loop from chargebacks to model retraining

### Strong-Hire Answer Pattern

- Explicitly budgets 100ms: signal assembly ~20ms + ML inference ~30ms + rules ~2ms = ~52ms, leaving headroom
- Uses Redis for velocity counters with TTL-based rolling windows (`INCR` + `EXPIRE`)
- Describes ML inference as a separate gRPC call to a model server, with fallback to rules-only if the model times out
- Explains 3DS as a liability shift mechanism, not just an auth step
- Knows that OFAC must be async (post-decision) because it's too slow for the synchronous path
- Describes the feedback loop: chargebacks → fraud labels → model retraining pipeline

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Redis velocity counter stale on failover | Redis replica lag metric > threshold | Temporarily disable velocity features; fall back to rules-only scoring |
| ML model latency spike | p99 inference latency exceeds budget | Circuit breaker: fall back to rules-only mode automatically |
| False positive rate spike from model drift | Monitor chargeback-to-block ratio | Alert on-call; freeze model rollout; roll back to previous model version |
| Global network signal DB overloaded | DB CPU/latency metric | Serve stale cache (TTL: 30s) instead of blocking on DB read |
| OFAC match async enrichment fails | Dead-letter queue depth rising | Page compliance team; re-queue for processing |

## Observability

- `radar.score.distribution` histogram — watch for bimodal shift (attack in progress)
- `radar.action.block_rate` gauge — sudden spike = new attack; sudden drop = rule misconfiguration
- `radar.scoring.latency_p99` — must stay under 80ms (leaves 20ms buffer)
- `radar.model.fallback_rate` — rising fallback = ML model unhealthy
- `radar.3ds.trigger_rate` — watch for merchant-specific spikes (may indicate card testing attack)
- `radar.false_positive_rate` — computed weekly from dispute outcomes
- Distributed traces on every scoring pipeline execution, tagged with `payment_intent_id`

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| ML + rule engine hybrid | Rules allow merchant override; ML captures unknown patterns | Two systems to maintain and tune | Rules-only — misses unknown patterns; ML-only — no merchant control |
| Async OFAC check | Stays off synchronous path; meets 100ms budget | Rare window where OFAC-listed card completes before check | Synchronous OFAC — too slow; fails latency budget |
| 3DS for medium-risk | Reduces false positives while shifting liability | Adds checkout friction (~10s); some cardholders abandon | Block all medium-risk — too many false positives |
| Redis for velocity counters | Sub-millisecond reads; atomic INCR for accuracy | Counter lost on Redis restart (short-term accuracy gap) | DB-backed counters — too slow for synchronous path |
| Global network signals | Unique signal strength (250M cards) | Requires cross-merchant signal sharing (privacy/regulatory constraints) | Per-merchant signals only — weaker model |

## Interview It

**Stripe framing:** "Design Stripe Radar. Stripe processes 1M transactions/minute. Every transaction must be fraud-scored before the card network authorization times out. Walk me through the design."

**Follow-ups:**
1. Your velocity counter Redis cluster fails over. What happens to fraud scoring during the 30-second failover window?
2. A merchant reports that Radar is blocking 15% of their transactions — much higher than their industry average. Walk me through your investigation process.
3. How would you design the 3DS trigger logic to minimize false positives without increasing chargebacks?
4. The ML model shows a sudden spike in false negatives (fraud getting through). What is your rollback strategy?
5. How do you handle the OFAC list check without violating your 100ms latency budget?

## Ship It

After this lesson, you should be able to:
- Name the six categories of fraud signals Radar uses and where each is stored.
- Explain the 100ms latency budget and how to allocate it across signal assembly + ML inference + rule engine.
- Describe why 3DS is a liability shift mechanism, not just an authentication step.
- Explain why OFAC checks are async and what the risk window is.
- Design the feedback loop from chargebacks to model retraining.

## Exercises

1. Sketch the Redis data model for velocity counters using rolling windows. What key schema and TTL strategy would you use?
2. Design the feature vector for the ML model. List 10 features and their data types.
3. Write the rule engine schema: how would a merchant express "block all prepaid cards from non-US IPs"?
4. Design the async enrichment pipeline: what queue technology, what retry policy, what dead-letter queue schema?

## Further Reading

- Stripe Radar documentation: stripe.com/docs/radar
- "How Stripe builds fraud prevention with machine learning" (stripe.com/blog/radar)
- "3D Secure 2 — what it means for merchants" (stripe.com/guides/3d-secure-2)
- Redis INCR and expiry patterns for rate limiting (redis.io/commands/incr)
