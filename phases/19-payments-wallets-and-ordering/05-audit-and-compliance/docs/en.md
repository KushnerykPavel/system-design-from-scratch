# Audit, Compliance, and Data Retention Constraints

> Compliance is not a side chapter; it changes data ownership, deletion semantics, and operational design.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design the audit and retention layer for financial systems where immutable records, privacy deletion, access review, and legal holds all pull in different directions.
**Prerequisites:** `12-security-abuse-and-multitenancy/05-privacy-and-deletion`, `12-security-abuse-and-multitenancy/06-threat-modeling`, `19-payments-wallets-and-ordering/01-payment-ledger`
**Estimated time:** ~60 min
**Primary artifact:** compliance-plan validator + retention matrix

## The Problem

Design the audit and compliance posture for payments, wallet, and order systems. Finance needs long-lived records, privacy teams need deletion workflows, security needs access traceability, and legal may require selective holds.

This lesson matters because strong system design does not stop at the happy path. Senior answers explain which records are immutable, which customer data can be minimized or detached, and how retention or deletion policies propagate across storage systems and backups.

## Clarify

- Which regimes matter most: PCI scope reduction, privacy deletion, financial retention, or regional residency?
- Are we retaining raw payment instrument data, tokens, or references only?
- What records must stay queryable online versus archived?
- How should legal hold override normal deletion behavior?

If the interviewer stays broad, assume payment instrument tokenization, long-lived financial event retention, privacy-driven deletion or pseudonymization for customer PII, and full audit logs for privileged access and manual adjustments.

## Requirements

### Functional

- Preserve immutable financial and audit records for required retention windows.
- Separate sensitive PII from long-lived financial records where possible.
- Support privacy deletion or anonymization workflows without corrupting books.
- Record access to privileged financial data and manual interventions.
- Apply legal hold and regional policy overrides to retention workflows.

### Non-functional

- Minimize PCI and sensitive-data blast radius.
- Make retention policy enforceable across hot, cold, and backup storage.
- Keep audit trails tamper-evident.
- Allow investigations without exposing more customer data than necessary.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Financial records | 20B rows/year | archive tiering and query strategy dominate storage decisions |
| Access audit events | 3M/day | operator visibility must scale too |
| Deletion requests | 500K/day globally | privacy workflows need throughput and verification |
| Retention windows | 90 days to 7+ years depending on record class | policy engine needs class-aware behavior |
| Peak factor | 10x after regulatory or incident-driven backfills | batch policy jobs must not disrupt core systems |

## Architecture

```text
product systems
  -> event / record classification
  -> immutable audit store
  -> PII vault / tokenization boundary
  -> retention and deletion policy engine
  -> archive tier + legal-hold controls
```

Design notes:

1. Store minimal customer identity in financial records; join to PII through tokens or indirection when possible.
2. Separate "delete customer-identifying data" from "destroy accounting evidence."
3. Treat retention as policy-driven by record class, region, and legal hold.
4. Log privileged reads and data exports, not just writes.

## Data Model & APIs

Core entities:

```text
record_class(class_id, retention_days, pii_level, residency_policy)
audit_event(event_id, actor_id, action, resource_ref, justification, created_at)
pii_subject(subject_id, token_ref, deletion_state, legal_hold)
retention_job(job_id, class_id, storage_tier, cutoff_time, state)
```

Useful interfaces:

- `ClassifyRecord(resource_ref, class_id)`
- `RequestSubjectDeletion(subject_id, reason_code)`
- `ApplyLegalHold(subject_id_or_case_id)`
- `ListPrivilegedAccess(resource_ref, time_range)`
- `RunRetentionJob(class_id, cutoff_time)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| PII copied into immutable ledger records | data classification scans and DLP alerts | strict schema boundaries and tokenization |
| deletion workflow removes data still under legal hold | policy conflict metrics and audit alarms | hold-aware deletion engine with deny-by-default behavior |
| privileged read is not logged | audit coverage gaps and control tests | mandatory access proxy and tamper-evident audit sink |
| archive tier makes investigations too slow | retrieval latency and unresolved-case age | searchable metadata and staged hot restores |

## Observability

- metric: privileged-access audit coverage and unlogged-access violations
- metric: deletion request backlog, legal-hold conflicts, and retention job lag
- metric: data classification violations by storage system
- log: every legal hold, deletion approval, export, and manual retention override
- trace: deletion request -> policy evaluation -> storage actions -> verification
- SLO: policy-controlled retention and deletion jobs complete within the compliance target window with full audit evidence

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| tokenize and separate PII from finance records | smaller compliance blast radius | more joins and operational boundaries | storing customer identity inline everywhere |
| immutable audit store | strong accountability | more storage and access tooling | mutable admin logs |
| policy engine across tiers | consistent enforcement | implementation complexity | ad hoc per-system retention scripts |

## Interview It

**Google framing:** "Design the compliance and audit posture for a payment and wallet platform." Expect pressure on privacy deletion versus financial retention and how operator actions are traced.

**Cloudflare framing:** "Design auditability and retention controls for a global billing system with regional constraints." Expect pressure on data residency, privileged access logging, and blast-radius reduction.

**Follow-ups:**
1. What if backups retain data longer than hot storage?
2. How do you support investigations without exposing raw PII widely?
3. What if one country requires in-region storage for selected billing data?
4. How do you verify deletion happened across derived systems?
5. What changes if PCI scope must shrink quickly?

## Ship It

- `outputs/tradeoff-matrix-audit-and-compliance.md`

## Exercises

1. **Easy** — Split one payment record into finance-safe data and PII-bound data.
2. **Medium** — Design a deletion workflow that preserves accounting evidence.
3. **Hard** — Redesign for conflicting regional residency and retention obligations.

## Further Reading

- [PCI DSS overview](https://www.pcisecuritystandards.org/standards/pci-dss/) — useful for scope-reduction thinking
- [Google Cloud data deletion guidance](https://cloud.google.com/architecture/architect-for-data-deletion) — strong framing for policy-driven deletion
