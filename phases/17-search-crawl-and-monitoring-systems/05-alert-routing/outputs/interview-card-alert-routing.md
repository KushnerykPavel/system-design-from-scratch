# Interview Card — Alert Routing and On-Call Signal Quality

---
lesson: 05-alert-routing
focus: balanced
---

## Clarify first
- Which severities page humans versus create tickets or chat notifications?
- Who owns routing metadata and schedule freshness?
- What grouping, suppression, and maintenance-window behavior is required?

## Must-size numbers
- Peak raw alert rate during correlated incidents
- Ack target and escalation timing by severity
- Number of active teams, rotations, and delivery targets
- Duplicate collapse ratio needed to protect responders

## Core design
- Normalize alerts from many sources into one intake model
- Deduplicate and group before human fanout
- Route through validated ownership and schedule data
- Track delivery, ack, escalation, and provider fallback outcomes

## Failure probes
- What happens if the paging provider fails?
- How do you route sev-1 alerts from an unowned service?
- How do you keep one regional outage from paging thousands of hosts separately?

## Trade-off summary
- Aggressive grouping vs over-merging unrelated issues
- Team flexibility vs central policy quality
- Fast paging vs rich context assembly
