# Abuse Prevention and Rate Limiting Layers

> Abuse control is not one limiter. It is a stack of fast cheap filters and slower smarter decisions.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Design layered abuse defenses that separate human traffic from automation, protect expensive paths early, and avoid punishing healthy customers under attack conditions.
**Prerequisites:** `02-estimation-and-cost/05-burstiness`, `10-reliability-retries-and-backpressure/04-load-shedding`, `14-rate-limiters-ids-and-hashing/02-distributed-rate-limiter`
**Estimated time:** ~75 min
**Primary artifact:** abuse-layer worksheet + interview card

## The Problem

Design abuse prevention for an API and web platform that faces credential stuffing, scraping, brute force login attempts, and tenant-specific bursts.

This lesson matters because "add rate limiting" is rarely enough. Strong answers separate:

- cheap edge filters from expensive behavioral checks
- account protection from network protection
- accidental client spikes from adversarial traffic
- enforcement for anonymous traffic, authenticated users, and tenants

## Clarify

- Is the main threat scraping, login abuse, fraudulent account creation, API key sharing, or DDoS-style overload?
- Which requests are expensive enough to deserve stricter protection?
- Do we need to protect user accounts, backend cost, origin capacity, or all three?
- How much false-positive pain is acceptable for free-tier versus enterprise traffic?

Assume layered protections: edge reputation and rate limits first, account-aware checks second, and tenant-aware safeguards for high-cost APIs.

## Requirements

### Functional

- Throttle abusive traffic across anonymous, account, and tenant dimensions.
- Support escalation from soft friction to hard blocking.
- Preserve an explainable path for false-positive investigation.

### Non-functional

- Keep cheap filters at the earliest boundary.
- Avoid large collateral damage during attacks.
- Allow policies to adapt quickly as adversaries shift behavior.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak traffic | 500K req/s | determines how much enforcement must happen at the edge |
| Login QPS | 20K req/s with 50x attack bursts | login paths need tighter and separate controls |
| API tenants | 50K active | motivates tenant-aware quotas and fairness |
| Bot variability | IP churn plus account churn | one-dimensional rate limits are insufficient |
| Rough cost | edge counters, reputation signals, and review tooling | layered defense trades cost for lower false positives |

## Architecture

```text
client
  -> edge reputation and coarse rate limits
  -> challenge / friction layer
  -> account or API-key aware controls
  -> tenant and route protection
  -> origin service
```

Good layering usually looks like:

1. Drop obviously bad traffic cheaply.
2. Slow suspicious traffic with challenges or proof-of-work style friction.
3. Enforce per-account or per-API-key limits on sensitive flows.
4. Protect expensive backend actions with tenant-aware budgets and anomaly rules.

## Data Model & APIs

Core entities:

- `SignalKey`
- `AbusePolicy`
- `Decision`
- `ChallengeState`
- `EscalationRule`

Useful APIs:

- `EvaluateEdge(request)`
- `EvaluateAccount(account, action)`
- `EvaluateTenantBudget(tenant, route)`
- `ExplainBlock(requestID)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| all traffic from one NAT is blocked | high enterprise deny rate behind shared IPs | combine identity, session, and behavioral keys |
| attack bypasses one weak dimension | origin cost rises while IP limits look normal | add multi-dimensional and route-specific controls |
| challenge service fails and blocks healthy users | challenge error rate spikes | fail open for low-risk paths, fail closed for credential attacks |
| policies overfit one incident | false positives rise after emergency rule | time-box emergency rules and review outcomes |

## Observability

- metric: allows, challenges, throttles, and hard blocks by policy layer
- metric: login failure ratios, tenant budget burn, and backend cost by route
- metric: false-positive review rate and challenge solve success
- log: sampled block explanations with key dimension and rule version
- trace: expensive protected route with abuse-decision annotations
- SLO: abuse controls should protect origin capacity without broad enterprise collateral damage

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| layered enforcement | catches more abuse with cheaper early drops | more control-plane complexity | one global limiter |
| friction before hard block | lowers false positives | more product complexity and UX cost | immediate hard deny for all suspicious traffic |
| tenant-aware limits | protects shared infrastructure fairly | more cardinality and policy tuning | IP-only enforcement |

## Interview It

**Google framing:** "Design abuse prevention for a login or API platform." Expect questions about false positives, shared IPs, and adaptive defenses.

**Cloudflare framing:** "Protect a global edge service from scraping and credential stuffing." Expect emphasis on edge-first cheap decisions and how to preserve good traffic under attack.

**Follow-ups:**
1. How do you protect login without locking out users behind one office NAT?
2. When should a policy challenge instead of block?
3. Which expensive routes deserve stricter budgets?
4. How do you investigate a false positive for an enterprise customer?
5. What changes at 10x attack volume?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/abuse-layer-worksheet.md`
- `outputs/interview-card-abuse-prevention.md`

## Exercises

1. **Easy** — Sketch the cheapest three filters you would apply before origin.
2. **Medium** — Add a design for credential stuffing protection without making login unusable from shared networks.
3. **Hard** — Redesign the system for a large scraper that rotates IPs, user agents, and accounts aggressively.

## Further Reading

- [Cloudflare bot management](https://www.cloudflare.com/application-services/products/bot-management/) — practical layered abuse-control framing
- [Google reCAPTCHA Enterprise](https://cloud.google.com/recaptcha-enterprise/docs/best-practices) — useful examples of friction and risk scoring
