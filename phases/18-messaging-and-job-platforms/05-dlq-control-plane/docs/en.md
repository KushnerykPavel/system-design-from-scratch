# Dead-Letter and Replay Control Plane

> Dead letters are not a graveyard; they are an operations interface for deciding what to trust, retry, repair, or discard.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design the control plane around dead-letter queues and replay so operators can recover safely without turning replay into an outage multiplier.
**Prerequisites:** `07-queues-streams-and-workflows/04-dlq-and-replay`, `10-reliability-retries-and-backpressure/07-bulkheads`, `18-messaging-and-job-platforms/01-distributed-message-queue`
**Estimated time:** ~60 min
**Primary artifact:** replay-request validator + replay runbook

## The Problem

Design the dead-letter and replay controls for a messaging platform. Messages land in a dead-letter queue after bounded retry failure, schema incompatibility, or unsafe downstream conditions. Operators need tools to inspect failures, replay only the right subset, and protect live traffic while doing recovery.

This lesson matters because many designs add a DLQ and stop there. Strong answers explain how messages reach the DLQ, how failures are classified, how replay is scoped, and why replay needs actor attribution, dry-run tools, and rate limiting.

## Clarify

- Why do messages enter the DLQ: poison payloads, temporary dependency failure, schema drift, or policy rejection?
- Do replays trigger side effects that could duplicate business actions?
- Is replay self-serve for application teams or tightly controlled by operators?
- Do we need partial replay by tenant, partition, key, or time window?

If the prompt stays open, assume multiple failure classes, side-effect risk on replay, scoped operator-driven replay, and strict isolation from live traffic.

## Requirements

### Functional

- Route repeatedly failing or invalid messages into DLQ storage with cause metadata.
- Let operators search and sample failed messages before replay.
- Support scoped replay by topic, tenant, key, partition, or time window.
- Allow dry-run estimation and replay throttling.
- Preserve audit history of who replayed what and why.

### Non-functional

- Prevent replay from overwhelming live consumers or brokers.
- Keep DLQ retention long enough for incident recovery without becoming indefinite junk storage.
- Preserve enough metadata for debugging root cause.
- Make irrecoverable failures distinguishable from transient ones.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| DLQ ingress | 0.05% steady-state, 5% during incidents | incident-mode scale can be orders of magnitude larger |
| Replay batch size | up to 100M messages | replay needs quotas and staging |
| DLQ retention | 14 days standard, 30 days audited | storage and search need policy tiers |
| Operator actions | hundreds/day normal, spikes during incidents | UI and auditability matter operationally |
| Peak factor | 20x on bad deploy rollback | replay tooling often activates during the worst moments |

## Architecture

```text
main delivery path
  -> retry policy
  -> failure classifier
  -> dlq storage + index
  -> inspection UI / APIs
  -> dry-run estimator
  -> throttled replay workers
```

Design notes:

1. Record failure reason and last processing context with each dead-lettered message.
2. Separate DLQ storage from replay execution so operators can inspect before acting.
3. Require replay scope, rate limit, and actor identity up front.
4. Support "repair then replay" for schema or policy failures instead of binary replay-or-drop thinking.

## Data Model & APIs

Core records:

```text
dlq_message_id
origin_topic
origin_partition
original_cursor
tenant_id
failure_reason
last_error_class
first_failed_at
last_failed_at
replay_status
```

Useful interfaces:

- `ListDeadLetters(filter, page_token)`
- `SampleDeadLetters(filter, limit)`
- `CreateReplay(scope, rate_limit, reason, actor)`
- `EstimateReplay(scope)`
- `CancelReplay(replay_id)`
- `MarkDiscarded(dlq_message_id, reason)`

The design gets stronger when replay is treated as a controlled job with validation, not as an unrestricted button.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| replay reintroduces poison messages endlessly | replay failure ratio and repeated DLQ re-entry | replay caps, dedupe, and "repair required" classification |
| broad replay saturates live consumers | live lag increase during replay | isolated replay workers and rate limits |
| operators replay the wrong scope | dry-run estimate mismatch and audit review | preview samples, explicit filters, confirmation guardrails |
| DLQ becomes junk drawer with no triage | old-age DLQ inventory and unknown-reason counts | retention tiers, required reason codes, and cleanup workflow |

## Observability

- metric: DLQ ingress by failure class and topic
- metric: replay throughput, replay failure ratio, and live-lag impact
- metric: DLQ age distribution and untriaged message count
- metric: discard rate versus successful recovery rate
- log: every replay request with actor, scope, rate limit, and reason
- trace: sampled message lifecycle from primary failure to replay outcome
- SLO: replay operations must not push live subscription lag beyond the protected platform threshold

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| scoped replay with dry run | lower blast radius | slower operator workflow | global replay button |
| long DLQ retention | better investigation and recovery | higher storage and policy burden | immediate discard after retry exhaustion |
| separate replay workers | protects live traffic | extra infrastructure | replay on same hot path as normal delivery |

## Interview It

**Google framing:** "Design the replay and dead-letter controls for an internal event platform." Expect follow-ups on replay scope, auditability, and how operators avoid duplicating side effects.

**Cloudflare framing:** "Design recovery controls for a shared messaging platform under tenant incidents." Expect pressure on blast radius, noisy-neighbor protection, and safe multi-tenant tooling.

**Follow-ups:**
1. How do you distinguish transient downstream failure from truly bad payloads?
2. What if one team's replay keeps harming a shared dependency?
3. How do you let tenants self-serve without giving them dangerous broad replay power?
4. When should messages be discarded instead of replayed?
5. How would you support message repair before replay?

## Ship It

- `outputs/replay-runbook-dlq-control-plane.md`

## Exercises

1. **Easy** — Define the minimum metadata you need on a dead-lettered message.
2. **Medium** — Design a dry-run preview for replaying only one tenant's failures from the last hour.
3. **Hard** — Redesign the control plane when replayed messages can trigger irreversible billing side effects.

## Further Reading

- [Enterprise Integration Patterns: Dead Letter Channel](https://www.enterpriseintegrationpatterns.com/patterns/messaging/DeadLetterChannel.html) — useful baseline concept
- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — good operational framing for safe recovery tooling
