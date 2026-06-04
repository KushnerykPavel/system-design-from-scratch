---
lesson: 07-pipeline-backpressure
focus: balanced
---

# Backpressure Response Checklist

## Detect Early

- Oldest backlog age
- Ingest versus drain gap
- Critical completion delay
- Replay share of consumer capacity

## Protect First

- Critical traffic lane
- Producer throttling or quotas
- Replay pause or rate limit
- Downstream circuit protection

## Review Questions

- Which work can be delayed safely?
- Which work can be dropped safely?
- Which producers can slow down?
- What backlog age becomes a product incident?
