# Bot Mitigation and Abuse Control Planes

> Abuse systems are judged less by their cleverest detector than by how safely they roll out, how clearly they explain decisions, and how rarely they hurt legitimate traffic at scale.

**Type:** Learn
**Company focus:** Cloudflare
**Learning goal:** Design bot mitigation and abuse-control architectures that separate fast edge enforcement from slower model and rule management while keeping false positives, customer trust, and rollout safety visible.
**Prerequisites:** `12-security-abuse-and-multitenancy/03-abuse-prevention`, `14-rate-limiters-ids-and-hashing/02-distributed-rate-limiter`, `22-cloudflare-edge-platform-design/01-global-api-edge-gateway`
**Estimated time:** ~75 min
**Primary artifact:** enforcement matrix + rollout checklist

## The Problem

Design a bot mitigation platform for traffic passing through a global edge network. The platform should classify suspicious behavior, enforce challenges or blocks, and let customers choose policy levels without turning model rollouts or rule changes into widespread false-positive incidents.

The hard part is balancing detection quality with latency, transparency, and safe rollout. A beautiful detector that burns good traffic or cannot be debugged is not a good platform.

## Clarify

- Is the product defending login, content scraping, API abuse, or all of them?
- Are mitigations challenge, rate limit, score, log-only, or hard block?
- How much per-request latency budget exists for detection?
- Who owns the final policy: platform defaults, customer rules, or hybrid?

## Requirements

### Functional

- Score or classify requests using edge-visible features and optional shared intelligence.
- Apply mitigations such as allow, challenge, throttle, or block.
- Support customer-specific policy tuning with safe defaults.
- Expose enough reason codes for investigation and appeals.

### Non-functional

- Keep false positives and rollout blast radius bounded.
- Separate fast data-plane decisions from slower rule and model management.
- Resist adversarial adaptation and rule abuse.
- Make enforcement explainable and auditable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Request volume | 10M+ requests/sec globally | forces cheap hot-path decisions |
| Feature lookup latency budget | single-digit milliseconds | determines what can run synchronously |
| Rule or model rollout frequency | many times per day | makes control-plane safety critical |
| False-positive tolerance | very low on login and API flows | drives staged enforcement |
| Tenant policy count | 100K+ customers | requires scalable policy composition |

## Architecture

```text
request at POP
  -> feature extraction
  -> local heuristics / cached model inputs
  -> score or class decision
  -> enforcement action
  -> event stream to abuse intelligence control plane
  -> model/rule updates back to POPs
```

Key ideas:

1. Keep the hot path cheap and deterministic enough to debug.
2. Treat model and rule rollout like a safety-sensitive control plane.
3. Separate detection confidence from enforcement severity.
4. Preserve per-request reason codes even when using learned systems.

## Data Model & APIs

Useful records:

```text
decision(
  request_id,
  tenant_id,
  score,
  action,
  reason_codes,
  policy_version
)
```

Helpful APIs:

- `ScoreRequest(features)`
- `SelectMitigation(score, tenant_policy)`
- `ExplainDecision(request_id)`
- `PublishPolicyVersion(version)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| new rule causes false-positive spike | allow-to-block mix and support tickets jump | staged rollout and shadow mode |
| model is accurate but too slow | latency budget burn at POPs | cache features, simplify hot path, async enrichments |
| customers create conflicting policies | policy compile or decision anomalies | normalized precedence rules and validation |
| adversary adapts to one visible challenge | attack success rebounds after mitigation | layered signals and rotating mitigations |
| reason codes become opaque | operators cannot explain blocks | mandatory explainability schema |

## Observability

- metric: decision mix by action and tenant
- metric: challenge solve rate and false-positive proxies
- metric: decision latency by detector stage
- metric: policy version adoption across POPs
- log: sampled decisions with feature and reason summaries
- trace: request -> score -> action -> post-action outcome

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| score separate from action | safer policy tuning | more product complexity | one hard-coded block threshold |
| shadow rollout for rules | lower false-positive risk | slower full deployment | immediate global enforcement |
| cached local signals | low latency | freshness complexity | remote intelligence lookup on every request |

## Interview It

**Google framing:** "Design abuse detection for a large-scale API or login system." Strong answers still discuss feedback loops, false positives, and rollout safety.

**Cloudflare framing:** "Design bot mitigation at the edge." Strong answers must cover fast local decisions, policy/version rollout, customer tuning, and explainability under adversarial pressure.

**Follow-ups:**
1. What if customers want different challenge tolerance by path?
2. What if model quality improves but latency doubles?
3. What if the safest action is log-only during rollout?
4. What if an attacker learns the challenge policy and shifts tactics?

## Ship It

- `outputs/enforcement-matrix-bot-mitigation.md`
- `outputs/rollout-checklist-bot-mitigation.md`

## Exercises

1. **Easy** — Compare rate limiting with challenge-based mitigation for login protection.
2. **Medium** — Design a shadow rollout for a new bot score model.
3. **Hard** — Redesign for an API product where false positives are more damaging than moderate bot leakage.

## Further Reading

- [CAPTCHA](https://en.wikipedia.org/wiki/CAPTCHA) — useful baseline on challenge-based mitigation
- [Cloudflare engineering blog](https://blog.cloudflare.com/) — practical security and abuse mitigation context
