# Cloudflare Full Mock Loop

> Cloudflare-style interviews reward the same structure as strong system design interviews everywhere, but they are quicker to punish hand-wavy edge claims, unbounded failover, and magical control planes.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Run a full Cloudflare-style mock that forces explicit edge trade-offs around POP behavior, origin protection, routing, abuse controls, and operational safety.
**Prerequisites:** `03-design-framework-and-timing/08-full-loop-drill`, `22-cloudflare-edge-platform-design/01-global-api-edge-gateway`, `22-cloudflare-edge-platform-design/05-bot-mitigation`
**Estimated time:** ~90 min
**Primary artifact:** Cloudflare mock scorecard validator + practice card

## The Problem

Run a timed full mock on prompts that feel plausible for an edge, platform, or traffic-heavy infrastructure interview.

Example prompts:

- design global API gateway at the edge
- design CDN invalidation and purge propagation
- design origin protection for customer backends
- design traffic steering for a global edge platform
- design bot mitigation for API and login traffic

This lesson is not about memorizing one Cloudflare answer. It is about showing that you can stay structured while discussing POPs, shields, purge fanout, policy rollout, and incident safety.

## Clarify

- Which user-visible property matters most: latency, origin safety, freshness, or abuse control quality?
- Which path is data plane and which path is control plane?
- What degraded mode is acceptable if a subsystem is slow but not fully down?
- Which deep dive will most improve confidence in the design?

If the prompt stays broad, choose one hot path, one control-plane risk, and one failure amplification risk to organize the answer.

## Requirements

### Functional

- Clarify the workload and edge-specific priorities.
- Build a small but useful capacity model.
- Present a high-level design before diving deep.
- Choose one meaningful Cloudflare-style deep dive.
- Close with failure modes, observability, and redesign.

### Non-functional

- Keep edge claims concrete and operationally realistic.
- Avoid generic "just route to nearest region" answers.
- Show control-plane safety, not only data-plane speed.
- Preserve explainability for operator and customer-facing behavior.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification window | 3 to 5 min | enough to define the hot path and risk |
| Sizing window | 5 to 7 min | purge fanout, failover multipliers, and QPS should be visible |
| High-level design window | 10 to 12 min | needed to separate control and data planes |
| Deep-dive window | 8 to 12 min | where edge-specific judgment becomes visible |
| Close and redesign window | 5 to 8 min | tests operational ownership |

## Architecture

Use this loop:

```text
clarify edge goal
  -> size traffic and amplification risk
  -> define data plane and control plane
  -> propose architecture
  -> choose one deep dive
  -> explain failure, observability, and rollout
  -> redesign after a changed edge constraint
```

Strong deep-dive choices for this phase:

1. purge propagation safety
2. origin failover guardrails
3. policy rollout visibility across POPs
4. abuse enforcement false-positive control
5. cost-aware traffic steering

## Data Model & APIs

Main answer artifact:

```text
cloudflare_mock_answer(
  scope,
  priorities,
  sizing,
  control_plane,
  data_plane,
  deep_dive,
  failure_plan,
  observability,
  redesign
)
```

Scorecard fields:

- scope clarity
- edge-specific prioritization
- sizing quality
- control-plane and data-plane separation
- deep-dive quality
- failure and observability coverage
- redesign quality

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| answer treats edge as generic stateless proxying | no POP, shield, or propagation detail appears | force one edge-specific deep dive |
| control plane is ignored | policy rollout and invalidation safety are missing | require explicit control-plane section |
| failover story is unsafe | retries and region shifts are unbounded | ask for budgets, cooldowns, and headroom |
| explanation is impossible to debug | no reason codes or path explainability | require operator-visible decisions |

## Observability

- metric: whether the answer named POP, region, and origin-level signals when relevant
- metric: whether the design included one amplification-risk estimate
- metric: whether control-plane rollout visibility was addressed
- metric: whether each major failure named a detection signal
- log: assumptions, degraded modes, and redesign triggers
- trace: clarify -> size -> design -> deep dive -> close -> redesign
- SLO: produce a coherent Cloudflare-style design answer inside mock time without hiding edge trade-offs

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit control-plane section | shows platform maturity | uses answer time | only describing request forwarding |
| one strong edge deep dive | reveals real judgment | less breadth elsewhere | generic component tour |
| clear degraded mode | makes incident handling credible | weaker headline guarantees | magical always-fast and always-fresh claims |

## Interview It

**Cloudflare framing:** Expect follow-ups that test whether your system behaves safely when the network is healthy but origins are weird, the control plane lags, or one tenant creates outsized risk. Strong answers sound operational, not ornamental.

**Suggested prompts:**
1. Design a global edge API gateway.
2. Design a purge system for CDN content.
3. Design origin protection for fragile backends.
4. Design a bot mitigation platform at the edge.

## Ship It

- `outputs/skill-cloudflare-full-mock.md`
- `outputs/interview-card-cloudflare-full-mock.md`

## Exercises

1. **Easy** — Time-box a 45-minute Cloudflare-style answer into six stages.
2. **Medium** — Pick the best deep dive for "design global CDN purge" and justify it.
3. **Hard** — Re-answer "design global edge API gateway" after the interviewer changes the top priority from latency to origin protection.

## Further Reading

- [Cloudflare engineering blog](https://blog.cloudflare.com/) — practical edge and operations context
- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) — useful reminder of the core interview loop this mock adapts
