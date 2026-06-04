# Freshness Budget Worksheet

## Tier the reads

| Entity / flow | Maximum stale age | Default read path | Fresh fallback | Why |
|---------------|-------------------|-------------------|----------------|-----|
|               |                   |                   |                |     |
|               |                   |                   |                |     |
|               |                   |                   |                |     |

## Ask before committing

- Which user action requires read-after-write?
- Can this user roam across regions or replicas?
- If freshness is violated, should we reroute, block, or serve stale with a banner?
- Which metric pages us before customers complain?

## Routing rules

- `stale-ok`: follower or cache
- `fresh-soon`: follower if within budget, else leader
- `fresh-now`: leader or version-aware bypass

## Incident notes

| Scenario | Temporary policy |
|----------|------------------|
| one follower lags | |
| whole region lags | |
| failover in progress | |
