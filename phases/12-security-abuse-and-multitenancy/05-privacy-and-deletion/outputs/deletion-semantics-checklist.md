# Deletion Semantics Checklist

## Scope

- Which stores hold live user-visible data?
- Which stores are derived, cached, or analytical?
- Which copies live only in backups?

## Guarantees

- What becomes invisible immediately?
- What is deleted asynchronously?
- What exceptions exist for legal hold or backup retention?

## Safety

- Is there a tombstone or deletion ledger?
- Can restore replay deleted state?
- How will you detect stale deleted reads?
