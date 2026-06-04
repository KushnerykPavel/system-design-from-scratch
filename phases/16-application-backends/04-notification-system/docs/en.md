# Notification System

> Notification systems fail not when they cannot send, but when they send the wrong thing to the wrong user at the wrong time.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design a multi-channel notification backend that handles preferences, dedupe, retries, provider outages, and fatigue controls without treating all events as equal.  
**Prerequisites:** `12-security-abuse-and-multitenancy/05-privacy-and-deletion`, `07-queues-streams-and-workflows/02-delivery-semantics`, `10-reliability-retries-and-backpressure/06-retry-budgets`  
**Estimated time:** ~75 min  
**Primary artifact:** notification policy validator  

## The Problem

Design a service that turns application events into email, push, SMS, or in-app notifications. The challenge is not "put messages on a queue." It is deciding whether to send, through which channel, with what dedupe key, under what rate limits, and what to do when a provider is flaky.

Senior candidates distinguish event ingestion from policy evaluation and delivery orchestration, and they explain how user trust can be damaged even when the system is technically available.

## Clarify

- Which channels matter: push, email, SMS, in-app, or all of them?
- Do we need guaranteed delivery for critical events, or best-effort engagement notifications?
- Are user preference checks and quiet hours strict blockers before delivery?

If open-ended, assume multi-channel support, strict preference enforcement, strong guarantees for security and billing events, and best-effort delivery for low-priority engagement traffic.

## Requirements

### Functional

- Ingest product events and classify them by notification policy.
- Respect user preferences, locale, quiet hours, and compliance rules.
- Dedupe repeated events and support retries through external providers.
- Track delivery state and allow channel fallback for high-priority events.

### Non-functional

- Avoid duplicate sends under retries or provider uncertainty.
- Keep policy-evaluation latency low enough to process bursts.
- Contain blast radius when one downstream provider degrades.
- Make fatigue and saturation visible, not just transport success.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Event ingress | 4M events/min peak | drives policy engine and queue partitioning |
| Delivery attempts | 9M attempts/min across channels | reflects retries and fallback amplification |
| Preference reads | 1 per policy decision, often cached | shapes profile-store and cache design |
| Quiet-hour suppression | 20% of consumer traffic | means "not sent" is a normal outcome |
| Provider count | 2 to 3 per channel class | supports resilience and routing control |

## Architecture

```text
product events
  -> event bus
  -> policy engine
  -> dedupe / idempotency store
  -> channel orchestrator
     -> push provider
     -> email provider
     -> sms provider
     -> in-app inbox
  -> delivery state store
```

Important design split:

1. Policy decides whether a notification should exist.
2. Orchestration decides how and when to deliver it.
3. Providers are replaceable dependencies, not the source of truth.

## Data Model & APIs

Core records:

```text
notification_intent(intent_id, user_id, event_type, priority, dedupe_key, created_at)
user_preferences(user_id, channel_rules, quiet_hours, locale)
delivery_attempt(intent_id, channel, provider, state, error_code, updated_at)
```

APIs:

- `POST /v1/events`
- `GET /v1/notifications/{intent_id}`
- `POST /v1/preferences/{user_id}`
- `POST /v1/notifications/{intent_id}/cancel`

A useful answer explicitly names idempotency and dedupe boundaries so the same event is not re-sent three times during retries.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| provider starts timing out but not fully failing | timeout-rate increase and delivery state skew | circuit breaking, fallback provider, retry budget |
| dedupe store unavailable | duplicate-send alarms and policy replay backlog | fail-safe buffering or limited degraded mode for critical traffic |
| preference cache stale after user opt-out | opt-out violation reports and cache-age metrics | versioned preference cache with fast invalidation |
| low-priority traffic crowds out urgent notifications | queue lag by priority class | separate queues and budgets by notification class |

## Observability

- metric: send attempts and success rate by channel, provider, and priority
- metric: duplicate suppression count and retry-budget burn
- metric: preference opt-out violation rate
- metric: queue lag split by priority class
- log: provider response code and dedupe key for failed or ambiguous sends
- trace: event ingestion to final delivery-state transition
- SLO: critical notification intents reach an accepted provider state within the defined latency target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| explicit notification intents | clear audit trail and retry control | more storage and state transitions | fire-and-forget provider calls |
| per-priority queues | protects urgent traffic | more scheduling complexity | one global queue |
| provider fallback for critical classes | higher resilience | higher cost and possible duplicate ambiguity | single provider per channel |

## Interview It

**Google framing:** "Design a notification platform used by many product teams." Expect questions on policy ownership, priority classes, and how the platform avoids becoming a spam cannon.

**Cloudflare framing:** "Design a notification service for operational and security events." Expect pressure on urgency tiers, provider isolation, and observability for partial downstream failures.

**Follow-ups:**
1. What changes for password-reset or fraud alerts?
2. How do you stop duplicate notifications after a provider times out but maybe delivered?
3. What if one country allows only certain channels or sending windows?
4. How do you protect users from engagement-notification fatigue?
5. What changes if product teams want custom templates and experiments?

## Ship It

- `outputs/observability-checklist-notification-system.md`

## Exercises

1. **Easy** — Define three priority classes and the delivery guarantees each should get.
2. **Medium** — Design a dedupe-key strategy for comment, mention, and security notifications.
3. **Hard** — Redesign the system for strict regional residency of user preferences and delivery logs.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — practical guidance for alert quality, overload, and service degradation  
- [Building Reliable Distributed Systems in the Presence of Software Errors](https://queue.acm.org/detail.cfm?id=945134) — useful mindset for retry ambiguity and downstream dependency design  
