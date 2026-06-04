---
name: reliability-drill-sheet
phase: 10
lesson: 08
---

## Prompt

Design a global webhook delivery platform with durable ingest, safe retries, bounded backlog, and strong tenant isolation.

## Clarify

- delivery guarantee and freshness target
- tenant skew and abusive endpoint assumptions
- replay and operator-control expectations

## Must-size numbers

- ingest QPS
- delivery fanout
- duplicate submit rate during incidents
- maximum useful backlog age

## Score yourself

- Did you separate sync ingest from async delivery?
- Did you make retries safe and bounded?
- Did you explain overload behavior explicitly?
- Did you name the blast-radius boundary?
- Did you include the metrics that would prove the design works?
