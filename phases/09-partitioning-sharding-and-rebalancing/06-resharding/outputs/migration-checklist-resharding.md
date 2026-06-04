---
name: resharding-migration-checklist
phase: 09
lesson: 06
---

1. New shard map exists and is versioned.
2. Old and new layouts can both route requests safely.
3. Backfill progress is measurable by cohort.
4. Dual-write or change replay path is idempotent.
5. Parity checks cover primary data and derived indexes.
6. Cutover is cohort-based with pause and rollback.
