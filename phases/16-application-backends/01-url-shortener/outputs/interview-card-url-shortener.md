# Interview Card — URL Shortener

---
lesson: 01-url-shortener
focus: balanced
---

## Clarify first
- Redirect volume versus create volume
- Need for custom aliases, edits, and expiration
- Whether analytics or billing sits on the synchronous path

## Must-size numbers
- Peak redirect QPS and viral-link multiplier
- Link-create QPS and retry rate
- Metadata cardinality and retention window
- Abuse-screening throughput and revocation propagation target

## Core design
- Thin redirect path backed by cache and durable metadata
- Idempotent create path with alias reservation and policy screening
- Async click-event pipeline for analytics and abuse feedback

## Failure probes
- What happens when one link becomes 30% of global traffic?
- How are malicious links revoked quickly through caches?
- What breaks if the analytics pipeline lags for an hour?

## Trade-off summary
- Fast redirect path vs fresh revocation visibility
- Human-friendly aliases vs namespace collisions
- Strong synchronous accounting vs user-facing latency
