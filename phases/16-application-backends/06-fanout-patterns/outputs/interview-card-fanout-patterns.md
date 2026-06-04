# Interview Card — Application Fanout Patterns

---
lesson: 06-fanout-patterns
focus: balanced
---

## Clarify first
- Audience size and skew per event
- Read-to-write ratio
- Freshness target and tolerance for partial results

## Must-size numbers
- Recipients per logical event
- Read amplification on a cache miss
- Queue depth or merge depth during peaks
- Time from source event to user-visible result

## Core design
- Push when read speed matters and audience sizes stay bounded
- Pull when write amplification dominates
- Mixed when skew or product classes differ enough that one rule fails

## Failure probes
- What breaks first under celebrity-scale skew?
- Can the product show stale or partial results safely?
- Which metrics tell us the current strategy is wrong?
