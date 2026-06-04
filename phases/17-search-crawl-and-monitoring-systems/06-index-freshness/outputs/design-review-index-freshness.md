# Design Review Prompt — Index Freshness and Ranking Updates

Use this when reviewing a search-index update design.

## Clarify
- What freshness target applies to high-priority content?
- Are ranking updates on the same cadence as document visibility?
- How fast must deletes disappear from user-facing results?

## Core checks
- Document publish, ranking refresh, and delete handling are separate versioned paths.
- Rollouts use canaries, dual reads, or blue/green validation before broad promotion.
- Backfills do not silently starve live update freshness.
- Schema evolution is checked before publish, not after bad results appear.

## Failure probes
- How do you roll back a bad index snapshot?
- What happens when ranking features and document snapshots drift?
- How do you remove unsafe documents quickly without full reindex?
- How do you compare new and old result sets safely?
