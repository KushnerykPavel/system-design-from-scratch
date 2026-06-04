# Clock Skew, Ordering, and Time Assumptions

> Wall-clock time is a useful signal and a dangerous source of false confidence.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Explain where time-based ordering works, where clock skew breaks assumptions, and how to design safer ordering and expiry logic.
**Prerequisites:** `03-quorums`, `04-replication-lag`, `07-queues-streams-and-workflows/02-delivery-semantics`
**Estimated time:** ~60 min
**Primary artifact:** time-assumption checklist

## The Problem

System design answers often rely on timestamps for conflict resolution, expiry, id generation, or event ordering without mentioning what happens when clocks drift. That omission matters because:

- different machines do not agree perfectly on time
- network delay can reorder observations
- expiry and lease logic can fail dangerously under skew
- timestamp order is not the same as causality

This lesson gives you a safer vocabulary for time in distributed systems.

## Clarify

- Is time being used for ordering, lease ownership, expiry, or user-visible history?
- What happens if two writers have skewed clocks?
- Is "latest write wins" acceptable, or can it lose intent?
- Are there correctness-critical decisions tied to TTL or lease expiry?

## Requirements

### Functional

- Preserve useful ordering for writes or events without overtrusting wall clocks.
- Make expiry, lease, and timeout semantics resilient to modest skew.
- Explain when logical or version-based ordering is safer than timestamps alone.

### Non-functional

- Keep time assumptions auditable and observable.
- Avoid correctness bugs caused by hidden clock drift.
- Bound operational risk when NTP or regional timing is impaired.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Ordering-sensitive writes | 8K/s | enough to expose ambiguous ordering |
| Lease holders | 50K active | makes expiry safety important |
| Tolerated skew | single-digit ms normally, much larger during incidents | defines margin in TTL and leases |
| Event fan-out | 200K events/s | ordering shortcuts become tempting |
| Rough cost | extra metadata, logical versions, safety buffers | safer ordering costs some simplicity and latency |

## Architecture

A strong answer separates time usages:

1. **Display time** for users
2. **Logical version or sequence** for correctness ordering
3. **Lease and timeout buffers** for operational safety
4. **Event-time versus processing-time** for analytics or streams

Safer patterns:

- version numbers or leader-assigned sequence for critical ordering
- hybrid logical clocks or monotonic tokens where partial causality matters
- lease buffers and fencing tokens for ownership changes
- TTL margins so small skew does not create premature expiration

## Data Model & APIs

Helpful metadata:

```text
event -> {
  event_id,
  logical_version,
  wall_clock_ts,
  source_region,
  fence_epoch
}
```

Useful interfaces:

- `Write(key, expected_version)`
- `AcquireLease(owner, ttl_ms, fence_epoch)`
- `Compare(version_a, version_b)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| last-write-wins drops the real newer intent | conflicting updates with skewed timestamps | use versions or causal metadata instead of wall time alone |
| lease expires too early on one node | dual-owner or fence violations | skew margin and fencing tokens |
| TTL-based deletion happens before consumers catch up | replay or reader misses rise | delay delete, use tombstones, or version gating |
| analytics order differs from business order | event-time versus processing-time skew | define which notion of time each pipeline uses |

## Observability

- metric: clock offset and synchronization health by node or region
- metric: fence violations, dual-ownership attempts, or stale lease use
- metric: timestamp-order conflicts versus logical-version order
- log: writes rejected due to version or fencing mismatch
- trace: ownership changes and lease renewals with epoch values
- SLO: maximum safe skew budget for timing-dependent control paths

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| logical versions for correctness | safer ordering | more metadata and coordination | pure timestamp ordering everywhere |
| TTL or lease safety buffers | fewer skew-related bugs | slightly slower failover or cleanup | razor-thin expiry thresholds |
| explicit event-time semantics | clearer analytics behavior | more pipeline complexity | pretending processing order equals event order |

## Interview It

**Google framing:** "Design ordering for collaborative edits, leases, or distributed locks." The signal is whether you know wall clocks are not enough for correctness.

**Cloudflare framing:** "Design globally propagated policy state and rollout ownership." The signal is whether you discuss fencing, expiry safety, and the risks of skewed time at the edge.

**Follow-ups:**
1. When is last-write-wins acceptable?
2. Why does lease ownership often need fencing in addition to TTL?
3. What if NTP is healthy 99.9% of the time but occasionally drifts badly?
4. How do event-time and processing-time differ for analytics?
5. What is the simplest safe ordering story for a critical write path?

## Ship It

- `outputs/time-assumption-checklist.md`

## Exercises

1. **Easy** - Identify one feature where timestamps are display-only, not correctness metadata.
2. **Medium** - Explain why lease expiry needs a skew margin.
3. **Hard** - Redesign last-write-wins policy storage for a multi-region control plane with skew and delayed delivery.

## Further Reading

- [Time, clocks, and the ordering of events in a distributed system](https://lamport.azurewebsites.net/pubs/time-clocks.pdf) - classic grounding for ordering without perfect clocks
- [Designing Data-Intensive Applications](https://dataintensive.net/) - practical guidance on clocks, ordering, and time-based pitfalls
