# Retention, Deletion, and Compliance

> "Delete the data" is not a button. It is a policy, workflow, and proof problem.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Explain how retention classes, deletion workflows, legal holds, and auditability shape storage architecture beyond simple CRUD semantics.  
**Prerequisites:** `04-hot-and-cold-data`, `05-blob-metadata-separation`  
**Estimated time:** ~60 min  
**Primary artifact:** failure checklist + observability checklist  

## The Problem

Many interview answers treat data retention and deletion as a background cron job. Senior answers recognize that deletion semantics affect product behavior, storage layout, indexing, archive design, and compliance posture.

This lesson helps you reason about:

- retention classes
- user-initiated deletion
- legal hold and compliance exceptions
- auditability of delete state

## Clarify

- Is the requirement about user-visible delete, hard delete, or both?
- Are there different retention classes by tenant, region, or data type?
- Can some data be tombstoned first and physically deleted later?
- What proof does the business or regulator need after deletion completes?

## Requirements

### Functional

- Support configurable retention and expiration.
- Handle user deletion requests safely across primary and derived systems.
- Respect legal holds, audit rules, and downstream propagation needs.

### Non-functional

- Keep deletion workflows correct under retries and partial failure.
- Avoid serving deleted data from stale caches or indexes.
- Make compliance posture observable and auditable.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Deletion requests | 500K/day | enough to require durable workflow state |
| Derived systems | 4-8 downstream indexes, caches, or archives | deletion must propagate, not stop at primary store |
| Retained data | 300 TB across classes | lifecycle policy materially affects cost |
| Peak factor | 10x after policy changes or customer migrations | tests workflow durability and backlog handling |
| Rough cost | lifecycle jobs + tombstones + audit logs | deletion is an operational product, not free cleanup |

## Architecture

Typical delete flow:

```text
delete request
  -> policy check / legal hold check
  -> tombstone or delete marker
  -> async propagation to indexes, caches, archives
  -> verification and audit record
```

The strong answer distinguishes:

- user-visible disappearance
- physical byte removal
- proof of completion

## Data Model & APIs

Useful fields:

- `retention_class`
- `deleted_at`
- `delete_state`
- `legal_hold`
- `purge_after`

Useful APIs:

- `RequestDeletion(object_id, actor)`
- `ApplyRetentionPolicy(class, scope)`
- `GetDeletionStatus(object_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| primary data deleted but search index still serves it | stale-read verification detects leakage | use tombstones plus propagation acknowledgements |
| legal hold bypassed by normal delete flow | audit alerts and policy mismatches | central policy gate before destructive actions |
| physical purge runs before audit export completes | compliance gap discovered later | stage purge after downstream confirmations |
| delete backlog grows silently | queue age and completion latency rise | track per-system propagation lag and shed non-critical background work |

## Observability

- metric: deletion completion latency by retention class
- metric: downstream propagation lag for indexes, caches, and archives
- metric: stale deleted object detection count
- metric: legal hold override attempts or denials
- log: deletion state transitions with actor and policy reason
- SLO: deleted data disappears from serving paths within target time and purge completes within policy window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| tombstone before purge | safer propagation and audit trail | more workflow state and storage overhead | immediate hard delete everywhere |
| retention classes | better cost and policy alignment | more policy complexity | one global retention rule for all data |
| verified async deletion | operationally credible compliance story | slower full completion | assuming downstream systems eventually catch up |

## Interview It

**Google framing:** "Design deletion and retention for user-generated content." The signal is whether you distinguish product delete semantics from physical purge across derived systems.

**Cloudflare framing:** "Design retention policy enforcement for customer logs and configuration history." The signal is whether you include lifecycle policy, hold exceptions, and audit evidence.

**Follow-ups:**
1. What if caches still serve deleted content for a few minutes?
2. What if the customer asks for proof of deletion?
3. How do legal holds override expiry or user requests?
4. What if deletion volume spikes 20x after a migration?
5. Which parts should be synchronous versus asynchronous?

## Ship It

- `outputs/failure-checklist-retention-and-deletion.md`
- `outputs/observability-checklist-retention-and-deletion.md`

## Exercises

1. **Easy** — Define deletion states for a profile photo service.  
2. **Medium** — Design retention classes for metrics, audit logs, and uploaded documents.  
3. **Hard** — Explain how you would prove deletion completed across search, cache, archive, and backup surfaces.  

## Further Reading

- [Google Cloud data deletion concepts](https://cloud.google.com/architecture/security/deletion) — useful reference for deletion workflow thinking  
- [GDPR storage limitation principle](https://gdpr-info.eu/art-5-gdpr/) — grounding for why retention and deletion must be explicit  
