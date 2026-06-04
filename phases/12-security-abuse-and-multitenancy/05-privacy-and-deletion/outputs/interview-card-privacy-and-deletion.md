# Interview Card — Privacy, Retention, and Deletion Semantics

## Strong answer shape

- Separate immediate product invisibility from full backend purge.
- Track deletion as a workflow, not a single SQL statement.
- Explain how indexes, caches, analytics, blobs, and backups behave.
- Be precise about promises and exceptions.

## Common misses

- "Delete the row" as the entire plan.
- No restore or replay story.
- No discussion of backups.
- No observability for stale deleted reads.
