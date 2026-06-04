# Failure Checklist — Compaction and Lifecycle

- Compaction debt and read amplification have paging thresholds.
- Tombstone GC grace is tied to replica lag and repair assumptions.
- Lifecycle rules support dry-run and versioning.
- Destructive actions are logged with rule version and actor.
- Background bandwidth caps can be changed quickly during incidents.
- Repair, compaction, and lifecycle workers have clear priority order.
- Legal hold and retention exceptions block deletion reliably.
