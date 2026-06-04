# Design Review Checklist — Metadata Index Service

- Name the canonical owner of metadata state.
- List the supported query shapes before naming databases.
- Separate hot listing indexes from long-retention audit history.
- Require versioned or conditional metadata updates from async workers.
- Make replay and index rebuild part of the day-one design.
- Track index lag and stale-read impact explicitly.
- Reject unsupported filter combinations instead of silently scanning.
