---
lesson: 03-pagination-and-filtering
focus: balanced
---

## Query contract review

- What is the stable sort order?
- Which filters are indexed and officially supported?
- Is page size capped?
- Does the cursor encode filter context?
- Which requests should become async exports instead?

## Abuse checks

- reject unbounded page sizes
- reject unsupported filter combinations
- log query shape, not just endpoint name
