# Consistency Guarantee Worksheet

Use this when a design prompt reaches storage, replication, or cache freshness decisions.

## Entity map

| Entity / flow | User-visible action | Required guarantee | Allowed stale window | Why |
|---------------|---------------------|--------------------|----------------------|-----|
|               |                     |                    |                      |     |
|               |                     |                    |                      |     |
|               |                     |                    |                      |     |

## Ask explicitly

- Who must see their own write immediately?
- Can another user see older state for a few seconds?
- Can one user move across regions or replicas within a session?
- What failure becomes dangerous if stale data is served?

## Narrow guarantee menu

- `read-after-write`
- `monotonic reads`
- `causal or ordered visibility`
- `linearizable read or write`
- `bounded stale read`

## Failure review

| Failure | User sees | Detection | Mitigation |
|---------|-----------|-----------|------------|
| replica lag | | | |
| regional failover | | | |
| stale cache after write | | | |
| conflicting concurrent updates | | | |

## Close with one sentence

`For <entity/flow>, we promise <guarantee>, tolerate <stale window>, and detect violations with <metric>.`
