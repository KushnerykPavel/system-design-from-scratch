# Fraud and Risk Hooks Without Blocking the Core Path

> Fraud systems should influence money movement without becoming a hidden single point of checkout failure.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design the integration boundary between core payment or order flows and fraud/risk systems, including synchronous gates, asynchronous scoring, review queues, and safe degradation.
**Prerequisites:** `12-security-abuse-and-multitenancy/03-abuse-prevention`, `10-reliability-retries-and-backpressure/04-load-shedding`, `19-payments-wallets-and-ordering/04-order-state-machine`
**Estimated time:** ~60 min
**Primary artifact:** fraud-hook validator + risk integration checklist

## The Problem

Design how checkout, wallet, or order systems call into fraud and risk services. Some decisions must happen inline before accepting money movement, while richer models, graph analysis, or manual review can run later.

This lesson matters because weak answers either block every request on a giant risk service or ignore fraud altogether. Senior answers choose which checks are synchronous, what the failure policy is, and how to keep the core path safe when the risk platform degrades.

## Clarify

- Which risks matter most: stolen payment instrument, account takeover, promo abuse, or merchant fraud?
- What must block synchronously, and what can happen asynchronously after initial acceptance?
- Is false positive cost higher or lower than false negative cost for this product?
- Can the system place funds or orders into review states instead of hard declining?

If the interviewer stays broad, assume a small synchronous policy layer for high-confidence allow or deny decisions, asynchronous richer scoring, optional manual review, and clear degradation policy when fraud services time out.

## Requirements

### Functional

- Evaluate high-confidence risk checks inline on payment or order create.
- Support asynchronous risk enrichment and later review or reversal workflows.
- Allow manual-review queues for uncertain cases.
- Record model or rules version used for each decision.
- Keep risk outcomes linked to orders, payments, or wallet actions for investigation.

### Non-functional

- Keep checkout latency bounded even when risk infrastructure is slow.
- Avoid turning the fraud platform into a global availability dependency.
- Make degradation policy explicit and auditable.
- Support fast rule rollout and rollback.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Inline risk checks | 140K/s peak | shapes latency budget and cache strategy |
| Async scoring events | 500K/s with enrichment fanout | drives queueing and backfill capacity |
| Manual review rate | 0.1% to 1% of volume | affects tooling and staffing assumptions |
| Feature fetch latency budget | <20 ms inline budget | determines what can stay on the critical path |
| Peak factor | 8x during attacks or promo abuse | risk system must degrade predictably under hostile load |

## Architecture

```text
checkout / payment API
  -> lightweight synchronous policy engine
  -> allow / deny / review decision
  -> async risk enrichment pipeline
  -> case management / review tooling
  -> feedback loop to rules and models
```

Design notes:

1. Keep the inline path narrow: high-signal rules, feature cache, and bounded-time calls only.
2. Separate "block now" from "investigate later" decisions.
3. Log rule version or model version so post-incident analysis is possible.
4. Define timeout behavior deliberately; defaulting silently to "allow everything" or "deny everything" is rarely acceptable without product context.

## Data Model & APIs

Core entities:

```text
risk_decision(decision_id, subject_ref, stage, outcome, model_version, expires_at)
risk_signal(signal_id, subject_ref, signal_type, score, created_at)
review_case(case_id, subject_ref, reason_code, priority, state)
policy_version(version_id, source, activated_at, actor_id)
```

Useful interfaces:

- `EvaluateInlineRisk(subject_ref, context, idempotency_key)`
- `PublishRiskEvent(subject_ref, event_type, payload_ref)`
- `OpenReviewCase(subject_ref, reason_code, priority)`
- `ResolveReviewCase(case_id, action)`
- `ListRiskDecisionHistory(subject_ref)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| risk service times out on critical path | inline timeout rate and checkout latency regression | bounded fallback policy plus local feature cache |
| attack traffic overwhelms enrichment pipeline | queue backlog and sampling drops | isolate async path, load shed enrichment, preserve core path |
| model version causes false-positive spike | approval-rate change by version and review backlog surge | fast rollback and shadow-evaluation process |
| manual review becomes sinkhole | case age and unresolved-priority backlog | SLA tiers, auto-expiry, and better triage policies |

## Observability

- metric: inline risk latency, timeout rate, and allow/deny/review distribution
- metric: false-positive proxy signals such as manual-overturn rate by rule or model version
- metric: async enrichment lag and review-case age
- log: every policy activation, fallback-mode switch, and manual decision
- trace: checkout -> inline risk -> async enrichment -> review or reversal
- SLO: preserve checkout latency target while keeping critical risk controls active under expected attack load

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| narrow synchronous risk gate | bounded latency and availability impact | less rich context inline | blocking on full risk pipeline |
| async enrichment and review | deeper fraud coverage | delayed action on some abuse | forcing all decisions inline |
| explicit fallback policy | predictable degradation | uncomfortable product trade-offs | undefined timeout behavior |

## Interview It

**Google framing:** "Design fraud controls for checkout without destroying conversion." Expect pushback on inline versus async checks, false positives, and incident-time degradation.

**Cloudflare framing:** "Design abuse and payment-risk hooks for a global self-serve billing flow." Expect pressure on attack spikes, shared-platform isolation, and how policy rollout is audited.

**Follow-ups:**
1. What changes if fraud attacks increase 20x during a promotion?
2. How do you shadow-test a new model safely?
3. What if some enterprise tenants need stricter synchronous checks?
4. How do you reverse later if async scoring finds a problem after initial acceptance?
5. What if the risk vendor is unavailable in one region?

## Ship It

- `outputs/observability-checklist-fraud-hooks.md`

## Exercises

1. **Easy** — Choose three signals that belong on the synchronous path and justify them.
2. **Medium** — Define timeout behavior for the inline risk service and explain the trade-off.
3. **Hard** — Redesign for a marketplace where both buyer fraud and seller fraud matter differently.

## Further Reading

- [Google SRE workbook on overload](https://sre.google/workbook/handling-overload/) — strong guidance for degradation thinking
- [Stripe Radar concepts](https://stripe.com/radar/guide) — useful framing for layered fraud controls
