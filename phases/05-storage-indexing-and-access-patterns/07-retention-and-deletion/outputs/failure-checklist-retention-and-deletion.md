## Policy checks
- Are retention classes explicit by data type or tenant?
- Can legal hold block normal purge?

## Workflow checks
- Is delete tracked with tombstones or state transitions?
- Do caches, indexes, and archives receive delete propagation?

## Proof checks
- Can operators show deletion status to users or auditors?
- Is there a metric for stale deleted data still being served?
