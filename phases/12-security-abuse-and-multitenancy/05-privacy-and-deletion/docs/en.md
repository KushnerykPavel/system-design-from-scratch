# Privacy, Retention, and Deletion Semantics

> Deletion is easy to promise and hard to make true across replicas, caches, backups, and derived systems.

**Type:** Build
**Company focus:** Google
**Learning goal:** Design retention and deletion flows that respect privacy requirements without pretending distributed systems can delete every copy instantaneously.
**Prerequisites:** `05-storage-indexing-and-access-patterns/07-retention-and-deletion`, `07-queues-streams-and-workflows/06-outbox-and-cdc`, `08-consistency-replication-and-transactions/06-sagas`
**Estimated time:** ~60 min
**Primary artifact:** deletion semantics checklist + interview card

## The Problem

Design privacy-aware retention and deletion for a user data platform with primary databases, search indexes, caches, analytics pipelines, backups, and blob storage.

Strong answers do not say "just delete the row." They explain:

- which copies are user-visible versus internal versus backup-only
- what is deleted synchronously, asynchronously, or after retention expiry
- how tombstones, legal holds, and derived data are handled
- what guarantees can honestly be promised to users and regulators

## Clarify

- Is the requirement user-initiated deletion, retention expiry, legal hold, or all three?
- Which stores serve live product reads versus batch analytics or backups?
- What deletion latency matters: immediate product invisibility or full backend purge?
- Are we allowed to soft delete first and hard delete later?

Assume user-facing deletion should hide data quickly from product surfaces, with asynchronous fanout to downstream indexes and a separate backup-retention policy.

## Requirements

### Functional

- Support user-initiated deletion and retention-based expiry.
- Remove data from product reads promptly and fan out deletion to downstream systems.
- Track deletion state, exceptions, and audit evidence.

### Non-functional

- Avoid inconsistent partial deletion across derived systems.
- Preserve compliance evidence without retaining deleted user content unnecessarily.
- Make guarantees honest about backups and asynchronous propagation.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| User records | 5B primary records | drives async deletion pipeline scale |
| Deletion requests | 200K/day | fanout and backlog handling matter |
| Derived systems | search, cache, analytics, blob metadata, backups | deletion is a graph, not one table |
| Backup retention | 30 to 90 days | changes what "fully deleted" can mean operationally |
| Rough cost | tombstones, fanout jobs, audit history | privacy work adds storage and operational overhead |

## Architecture

```text
deletion request
  -> authoritative delete state
  -> product read suppression
  -> async fanout to indexes / caches / blobs / analytics
  -> completion tracking
  -> retention and backup lifecycle
```

Recommended pattern:

1. Mark an authoritative deletion state.
2. Hide content from user-facing reads immediately.
3. Fan out purge events to derived systems.
4. Track completion and exceptions per downstream target.

## Data Model & APIs

Core entities:

- `DeletionRequest`
- `RetentionPolicy`
- `DeletionTarget`
- `Tombstone`
- `LegalHold`

Useful APIs:

- `RequestDeletion(subjectID)`
- `ResolveRead(subjectID)`
- `ApplyDeletionTarget(target, subjectID)`
- `GetDeletionStatus(requestID)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| primary row deleted but search index still serves content | delete-status incomplete while search hits remain | authoritative tombstone plus fanout tracking |
| cache serves stale deleted data | cache hit after delete request | explicit invalidation and tombstone-aware read path |
| analytics pipeline retains raw identifiers indefinitely | downstream retention report mismatch | minimize raw PII and apply retention by pipeline stage |
| backup restore reintroduces deleted data | restored snapshot conflicts with tombstones | replay deletion ledger after recovery |

## Observability

- metric: deletion request volume, backlog age, and per-target completion lag
- metric: stale-read violations after deletion
- metric: legal hold count and retention exception count
- log: deletion state transitions and downstream failures
- trace: end-to-end delete fanout across product and derived systems
- SLO: user-facing deletion should hide content quickly while downstream completion remains measurable and bounded

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| immediate logical hide plus async purge | fast user-visible effect | more state tracking | wait for every store before hiding |
| deletion ledger / tombstones | prevents resurrection and stale reads | extra metadata and read-path logic | hard delete everywhere first |
| explicit backup policy language | honest compliance posture | less magical marketing | claiming instant deletion from immutable backups |

## Interview It

**Google framing:** "Design account deletion for a large data platform." Expect follow-ups on derived systems, backups, and what guarantee language is honest.

**Cloudflare framing:** "Design deletion and retention for logs or customer configuration at global scale." Expect focus on replicated stores, retention windows, and operational realism.

**Follow-ups:**
1. What becomes invisible immediately versus later?
2. How do backups affect the deletion guarantee?
3. How do you prevent deleted data from resurfacing after restore?
4. What data should never enter analytics in raw form?
5. What changes at 10x deletion volume?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/deletion-semantics-checklist.md`
- `outputs/interview-card-privacy-and-deletion.md`

## Exercises

1. **Easy** — List every place a user profile may exist outside the primary database.
2. **Medium** — Add a deletion ledger that survives backup restore.
3. **Hard** — Redesign for a product with strict legal holds and region-specific retention policies.

## Further Reading

- [Google Cloud data deletion guidance](https://cloud.google.com/architecture/security-foundations/data-protection-and-compliance) — good framing for lifecycle controls
- [GDPR Article 17](https://gdpr-info.eu/art-17-gdpr/) — the legal right to erasure that shapes many design discussions
