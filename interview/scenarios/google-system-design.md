# Google System Design Scenario Bank

Use for `/mock-interview google-system-design`. Reveal only the candidate prompt unless feedback mode is active.

## Scenario: collaborative-doc-backend

Candidate prompt:
Design the backend for collaborative documents where many users can edit the same document concurrently. Focus on storage, consistency, and conflict handling.

Hidden interviewer notes:
- Strong candidates clarify latency vs correctness expectations.
- Good deep dives: write path, conflict resolution, document model, presence vs persistent state.
- Follow-up: what changes if offline edits must merge later?

## Scenario: quota-enforcement-platform

Candidate prompt:
Design a multi-tenant quota platform used by internal services. It must enforce per-customer usage limits and expose low-latency checks on the critical path.

Hidden interviewer notes:
- Look for rate limiting vs quota distinction.
- Good deep dives: data freshness, hot-tenant isolation, abuse resistance, degraded mode.
- Follow-up: how do you support hourly, daily, and monthly limits together?

## Scenario: regional-metrics-platform

Candidate prompt:
Design a metrics platform for internal services across multiple regions. Teams need dashboards, alerts, and 30-day retention.

Hidden interviewer notes:
- Strong answers separate ingest, storage, query, and alerting.
- Good deep dives: cardinality control, rollups, alert freshness, cross-region design.
- Follow-up: what changes if cost must be cut by 40%?
