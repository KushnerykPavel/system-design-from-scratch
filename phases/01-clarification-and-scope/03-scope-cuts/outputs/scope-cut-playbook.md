---
lesson: 03-scope-cuts
focus: balanced
---

## Safe scope cuts

- one region before multi-region
- one core workflow before the full product surface
- one tenant class before all customer tiers
- one consistency mode before every possible edge case

## Unsafe scope cuts

- removing storage from a storage system
- removing synchronization from a collaboration system
- removing abuse handling from a public internet system
- removing the main write path from a write-heavy workflow

## Phrase it like this

- "I’ll keep `X` and defer `Y` so we can fully reason about the primary path and then discuss how `Y` changes the design."
