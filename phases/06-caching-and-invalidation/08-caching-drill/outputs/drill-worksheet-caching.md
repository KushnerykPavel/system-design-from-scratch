# Drill Worksheet: Caching

## Prompt

Design caching for:

____________________________

## Clarify

- hottest read path:
- most mutable object:
- strict read-after-write path:
- acceptable stale window:
- allowed cache layers:

## Capacity

| Dimension | Estimate |
|-----------|----------|
| peak reads | |
| writes | |
| hot-key skew | |
| object size | |
| rough origin cost of a miss | |

## Design

- cache pattern:
- freshness model:
- invalidation trigger:
- hot-key protection:
- user-visible consistency promise:

## Failure review

- what happens on hot expiry?
- what happens on missed invalidation?
- what happens after user writes?
- what metrics prove the cache is healthy?
