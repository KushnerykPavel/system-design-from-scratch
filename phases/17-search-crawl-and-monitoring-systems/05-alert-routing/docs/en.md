# Alert Routing and On-Call Signal Quality

> The goal of alerting is not to emit alerts; it is to wake the right human rarely, quickly, and with enough context to act.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design an alert routing system that emphasizes ownership, deduplication, escalation, suppression, and operator trust instead of treating paging as a simple notification fanout.
**Prerequisites:** `10-reliability-retries-and-backpressure/07-bulkheads`, `11-observability-slos-and-debugging/05-alert-design`, `11-observability-slos-and-debugging/06-runbooks`
**Estimated time:** ~60 min
**Primary artifact:** policy validator + interview card

## The Problem

Design an alert routing platform that receives alerts from monitoring systems, deduplicates and groups them, routes them to the correct on-call rotations, escalates when there is no acknowledgement, and reduces noisy or low-value pages.

This lesson matters because alerting design is often underrated in interviews. Senior answers talk about ownership data, severity policies, maintenance windows, escalation contracts, and how to measure whether the system is helping responders or burning them out.

## Clarify

- Are we routing only pages, or also tickets, chat notifications, and webhooks?
- Do teams own their own routing policies, or is there a central platform standard?
- What acknowledgement and escalation expectations exist by severity?
- Do we require maintenance windows, quiet hours, or incident-level suppression?

If the interviewer stays broad, assume a platform-owned core with team-managed routes, paging plus chat/ticket targets, and strong emphasis on reducing duplicate human interrupts.

## Requirements

### Functional

- Receive alerts from multiple monitoring backends.
- Deduplicate, group, and suppress related signals.
- Route alerts to owning teams and escalation targets.
- Track acknowledgements, escalations, and incident linkage.
- Enforce policy quality such as runbook presence and ownership metadata.

### Non-functional

- Keep routing latency low for high-severity pages.
- Preserve delivery and acknowledgement audit trails.
- Reduce false-positive human wakeups.
- Avoid one notification provider becoming a single point of paging failure.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Incoming alerts | 250K alerts/min peak | bursts happen during wide incidents |
| Active rotations | 8K | ownership lookup and config distribution must scale |
| Notification fanout | 4 targets average | retries and provider diversity matter |
| Ack target | under 5 minutes for sev-1 | routing and escalation paths must stay fast |
| Peak factor | 20x during correlated incidents | grouping quality is essential to human survival |

## Architecture

```text
monitoring systems
  -> alert intake
  -> dedup / grouping
  -> routing policy engine
  -> notifier adapters
  -> ack / escalate tracker
  -> audit store + analytics
```

Design notes:

1. Separate intake normalization from routing policy so multiple alert sources can share one platform.
2. Group aggressively for correlated failures, but keep the grouping rule explainable to responders.
3. Make ownership data explicit and validated before a sev-1 discovers an unowned service.
4. Track delivery outcomes per provider so fallback channels are evidence-based.

## Data Model & APIs

Core records:

```text
alert_fingerprint
severity
service
owner_team
group_key
policy_id
delivery_targets
ack_state
incident_id
```

Useful interfaces:

- `POST /v1/alerts`
- `POST /v1/alerts/{id}/ack`
- `POST /v1/policies/validate`
- `POST /v1/incidents/{id}/suppress`

The design is stronger when it clearly explains who owns routing metadata and how stale configs are prevented.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| unowned service emits sev-1 | policy-validation failures and unroutable alert count | block bad configs, default fallback team, ownership audits |
| one provider fails to deliver pages | provider-specific delivery errors and ack gap | multi-provider fallback and delivery health routing |
| incident storm pages every host separately | page fanout spike and duplicate fingerprint ratio | grouping, suppression, and incident-level correlation |
| stale on-call schedule routes to the wrong person | schedule sync lag and bounce metrics | cached last-good schedules and sync freshness alerts |

## Observability

- metric: routing latency by severity and target type
- metric: duplicate alerts collapsed per group key
- metric: acknowledgement time and escalation rate by team and severity
- metric: provider delivery failures and fallback usage
- log: routing decision with policy ID, owner, group key, and chosen targets
- trace: alert intake through policy engine to provider delivery attempt
- SLO: sev-1 pages reach at least one valid target within the contractual routing window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| aggressive grouping | fewer human interrupts | risk of over-grouping unrelated issues | one page per raw alert |
| mandatory ownership metadata | safer routing during incidents | more platform governance work | best-effort team lookup from naming conventions |
| multi-provider notification | better resilience | more operational integration cost | single paging vendor dependency |

## Interview It

**Google framing:** "Design internal alert routing for thousands of services and teams." Expect questions on ownership, escalation correctness, and signal quality metrics.

**Cloudflare framing:** "Design incident paging for a global edge platform." Expect questions on correlated failures, provider fallback, and keeping responders from being spammed during regional incidents.

**Follow-ups:**
1. How do you measure alert quality, not just alert volume?
2. What changes if teams demand custom grouping logic?
3. How do you keep schedule data fresh without making routing depend on a slow source system?
4. How do you prevent maintenance windows from hiding real incidents too broadly?
5. What is your fallback when routing metadata is missing at the worst moment?

## Ship It

- `outputs/interview-card-alert-routing.md`

## Exercises

1. **Easy** — Define the minimum fields required for a sev-1 routable alert.
2. **Medium** — Design a grouping rule for host-level alerts during a regional network failure.
3. **Hard** — Redesign the platform so customers can route their own alerts safely in a multi-tenant SaaS.

## Further Reading

- [Google SRE Workbook: Alerting on SLOs](https://sre.google/workbook/alerting-on-slos/) — great framing for signal quality
- [PagerDuty Incident Response Guide](https://response.pagerduty.com/) — practical grounding for escalation and acknowledgement design
