# Design Review — Distributed Message Queue

Use this when reviewing a queue design answer or implementation sketch.

## Scope

- What delivery guarantee is actually promised?
- What ordering boundary is promised?
- What replay window is required?

## Architecture Checks

- Is partitioning tied to a concrete access or ordering pattern?
- Is producer durability separate from consumer acknowledgement?
- Is consumer progress explicitly tracked?
- Is poison-message handling described?

## Failure Pressure

- What happens when a consumer crashes after fetching a batch?
- What happens when one partition goes hot?
- What happens during replay after a downstream outage?
- What happens when retention expires before a slow consumer recovers?

## Observability

- backlog age by topic and consumer group
- redelivery and expired-lease counts
- partition skew and storage pressure
- replay activity with actor attribution
