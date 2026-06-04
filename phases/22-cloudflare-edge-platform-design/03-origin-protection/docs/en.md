# Origin Protection and Health-Based Failover

> The origin is where elegant edge stories go to die if failover, retries, and shielding are not designed like first-class reliability features.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Design origin shielding and health-based failover policies that protect origins during partial failures without causing retry storms or cross-region overload.
**Prerequisites:** `10-reliability-retries-and-backpressure/01-timeouts-and-retries`, `13-multi-region-cdn-and-edge-traffic/02-failover-and-rto`, `22-cloudflare-edge-platform-design/01-global-api-edge-gateway`
**Estimated time:** ~75 min
**Primary artifact:** failover policy validator + failure checklist

## The Problem

Design origin protection for a global edge platform serving traffic to customer origins in multiple regions. POPs should route to healthy pools, shield fragile backends, and fail over only when the mitigation path is safer than staying local.

The real challenge is not "add health checks." It is avoiding false failovers, retry amplification, cost explosions, and invisible customer pain when one origin becomes slow but not fully down.

## Clarify

- Is the system protecting one origin cluster, multiple origin pools, or customer-defined origins?
- Is the higher priority lowest latency, strongest availability, or strongest origin safety?
- Are health decisions local to each POP, centralized, or hybrid?
- What is the tolerated failover cost in latency and egress?

## Requirements

### Functional

- Detect origin health degradation and full failure.
- Route requests through a shield or directly to origins based on policy.
- Bound retries, hedging, and failover decisions.
- Surface explainable reasons for origin selection and rejection.

### Non-functional

- Prevent retry storms and cascading overload.
- Avoid flapping between origin pools.
- Preserve per-tenant isolation when customer origins behave badly.
- Keep operators able to explain why traffic moved.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Global QPS | 1M to 3M | shapes how expensive wrong failover is |
| Healthy-to-failover multiplier | up to 2x on surviving region | drives reserve capacity |
| Health signal lag | seconds | determines false-positive risk |
| Retry budget | 5% to 15% extra traffic | caps incident amplification |
| Shield hit ratio | highly workload dependent | affects origin pressure and egress |

## Architecture

```text
POP request
  -> route policy
  -> local health snapshot
  -> origin shield or direct pool
  -> bounded retry / failover policy
  -> regional or cross-region origin pool
```

Key ideas:

1. Health signals should influence routing gradually, not as a binary toggle only.
2. Shielding is useful when it smooths traffic and localizes origin pressure.
3. Retry and failover policies need budgets, cooldowns, and jitter.
4. Explainability is part of the product because customers will ask why traffic shifted.

## Data Model & APIs

Useful records:

```text
origin_pool(
  name,
  region,
  healthy_threshold,
  failover_threshold,
  cooldown_seconds,
  shield_enabled
)
```

Helpful APIs:

- `EvaluateOriginHealth(pool)`
- `SelectUpstream(route, health_snapshot)`
- `ExplainFailover(request_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| origin is slow but not dead | tail latency rises before error rate does | outlier detection and gradual traffic shift |
| POPs disagree on health too aggressively | divergent failover rates by POP | hybrid health model with shared signals |
| retries overload weak origin further | retry budget burn spikes | strict retry budgets and circuit breaking |
| surviving region saturates after failover | queueing and saturation rise there | headroom policy and admission control |
| customer misconfigures origin pool | unusual failover and 5xx spikes after config change | staged rollout and config validation |

## Observability

- metric: origin selection rate by pool and POP
- metric: failover count and failback count
- metric: retry budget burn and hedge rate
- metric: shielded versus direct origin traffic
- log: origin rejection reasons and health threshold transitions
- trace: POP decision -> shield -> origin response -> retry or failover outcome

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| gradual health scoring | fewer false failovers | more policy complexity | binary up/down routing only |
| shield layer | smoother origin load and locality | extra hop and more capacity to manage | every POP directly to every origin |
| strict failover cooldown | prevents flapping | slower recovery to preferred path | instant failback on every healthy signal |

## Interview It

**Google framing:** "Design origin failover for a globally distributed serving layer." Good answers discuss health quality, blast radius, and overload propagation.

**Cloudflare framing:** "Protect fragile customer origins while keeping edge traffic available." Strong answers must cover shielding, explainable failover, and why retries are treated as dangerous.

**Follow-ups:**
1. What changes if origins are customer-operated and often misconfigured?
2. What if health checks say healthy but request latency is still bad?
3. What if cross-region failover is very expensive?
4. What if one tenant's origin pool is much weaker than the rest?

## Ship It

- `outputs/failure-checklist-origin-protection.md`
- `outputs/interview-card-origin-protection.md`

## Exercises

1. **Easy** — Compare binary health checks with score-based failover.
2. **Medium** — Add per-tenant failover guardrails for noisy or fragile origins.
3. **Hard** — Redesign for long-lived connections where failover semantics are much harsher.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — useful framing for tail latency and retry amplification
- [Cloudflare engineering blog](https://blog.cloudflare.com/) — practical context on edge-to-origin behavior
