---
lesson: 03-pagination-and-filtering
focus: balanced
---

## Clarify first

- browse vs random access vs export
- stable sort key
- filter cardinality
- consistency expectations during mutation

## Must-size numbers

- list QPS
- page size
- deep-page frequency
- export volume

## Core design

- cursor pagination for large mutable lists
- whitelisted indexed filters
- async path for bulk export

## Failure probes

- duplicates or skipped items
- deep offset scans
- unbounded filters
