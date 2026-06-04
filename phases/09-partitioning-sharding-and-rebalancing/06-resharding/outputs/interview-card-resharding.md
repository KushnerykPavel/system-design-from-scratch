---
lesson: 06-resharding
focus: balanced
---

## Clarify first

- range split vs key-shape change
- client compatibility constraints
- active writes during migration

## Must-size numbers

- duplicate storage needed
- write rate during dual phase
- cohort size and blast radius

## Core design

- old and new layouts coexist
- backfill first
- validate parity
- cut over by cohort

## Failure probes

- clients still use old shard map
- indexes are incomplete
- parity mismatch persists

## Trade-off summary

- compatibility safety vs temporary complexity
- smaller cohorts vs slower finish
- stronger validation vs extra read cost
