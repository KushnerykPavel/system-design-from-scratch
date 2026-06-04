# Failure Checklist — Object Storage

- Multipart sessions expire and are sweepable.
- Metadata cannot transition to `complete` before integrity checks pass.
- Retention and legal hold are enforced before physical deletion.
- Listing and policy APIs do not depend on the blob retrieval path.
- Orphan blob and orphan metadata counts are continuously tracked.
- Repair backlog is prioritized by durability risk, not only by age.
- Storage-class transitions are idempotent and auditable.
