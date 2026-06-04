# Quorum Planning Sheet

## Start with intent

- Is the system read-heavy, write-heavy, or balanced?
- Is stale data acceptable for any user path?
- Are conflicts acceptable, rare, or forbidden?
- Does the product prefer availability during node loss or tighter freshness?

## Choose the set

| N | R | W | Overlap? | Read latency impact | Write latency impact | Best fit |
|---|---|---|----------|---------------------|----------------------|----------|
|   |   |   |          |                     |                      |          |
|   |   |   |          |                     |                      |          |

## Repair plan

| Divergence source | Detection | Repair path | Cost |
|-------------------|-----------|-------------|------|
| replica missed write | | | |
| conflicting versions | | | |
| delete tombstone lag | | | |

## Close with one sentence

`We choose N=<>, R=<>, W=<> because the system values <freshness/availability>, and we handle divergence with <versioning + repair strategy>.`
