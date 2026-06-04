# Drill Worksheet - Consistency, Replication, and Transactions

## Prompt

Use one ambiguous system prompt and fill this out in under 12 minutes.

## Guarantee map

| Entity / flow | Risk if stale or duplicated | Needed guarantee | Read path | Write path |
|---------------|-----------------------------|------------------|-----------|------------|
|               |                             |                  |           |            |
|               |                             |                  |           |            |
|               |                             |                  |           |            |

## Narrow decisions

- Replication model:
- Transaction boundary:
- Saga or compensation boundary:
- Ordering or time assumptions:

## Failure review

| Failure | Detection | Immediate response |
|---------|-----------|--------------------|
| lag on critical read path | | |
| hotspot on critical write key | | |
| partial workflow failure | | |
| unsafe failover or skew event | | |

## One-sentence close

`We use <model> for <critical flow> because it protects <invariant>, while cheaper paths accept <bounded stale/async behavior> and are monitored by <metric>.`
