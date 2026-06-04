# Checklist — 10x Redesign

## Before Redesigning

- What exact dimension increased by 10x?
- Which original estimate is now invalid?
- What is the new dominant bottleneck?
- Which guarantee still matters most?

## During Redesign

- choose the smallest meaningful architecture change
- say what gets cheaper or faster
- say what gets more complex or riskier
- define migration and rollback
- add one new metric for the new bottleneck

## After Redesign

- validate the old bottleneck moved
- compare cost before and after
- check degraded mode under the new shape
