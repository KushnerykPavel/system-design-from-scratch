---
lesson: 05-event-schema-evolution
focus: balanced
---

## Schema checks

- Is the change additive, breaking, or semantic?
- Can old consumers ignore the new field safely?
- Are replay and backfill consumers covered?
- Is a new version or event type cleaner than in-place mutation?
- Do deploy checks enforce compatibility automatically?
