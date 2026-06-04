# Consistency Spectrum in Practice

> "Eventually consistent" is not a design answer until you say what can be stale, for whom, and for how long.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Translate abstract consistency terms into user-visible guarantees, operational costs, and interview-safe trade-offs.
**Prerequisites:** `01-clarification-and-scope/06-workload-shape`, `03-design-framework-and-timing/03-diagram-then-dive`, `06-caching-and-invalidation/07-cache-consistency`
**Estimated time:** ~75 min
**Primary artifact:** consistency guarantee worksheet

## The Problem

Candidates often jump from requirements to storage choices and say "strong consistency" or "eventual consistency" as if that closes the discussion. It does not. Senior answers map consistency to specific observations:

- who must see their own write immediately
- which users can tolerate bounded staleness
- what happens during failover or partition
- where the system intentionally allows divergence

This lesson gives you language for making those guarantees explicit without overpromising.

## Clarify

- Which action must reflect immediately in a follow-up read?
- Is stale data merely confusing, or is it financially or security harmful?
- Do users stay in one session and region, or move across replicas and regions?
- Is the main risk stale reads, lost writes, conflicting updates, or ordering confusion?

## Requirements

### Functional

- Define the consistency contract for critical user journeys.
- Separate correctness-critical entities from convenience reads.
- State how failover changes read and write behavior.

### Non-functional

- Keep guarantees understandable to product and operations teams.
- Avoid claiming a stronger model than the architecture can truly enforce.
- Bound cost and latency for paths that do not need strict coordination.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Read QPS | 300K req/s | enough to tempt broad replica or cache use |
| Write QPS | 15K req/s | enough that coordination overhead is visible |
| Cross-region sessions | 10% | exposes monotonic-read and failover questions |
| Critical write-after-read paths | 3 product flows | forces selective strong guarantees |
| Rough cost | extra coordination, replica traffic, higher tail latency | makes stricter guarantees visibly expensive |

## Architecture

Think of consistency as a menu of promises, not a binary toggle:

1. **Read-after-write** for the actor who just mutated state.
2. **Monotonic reads** so one user does not see time go backward.
3. **Prefix or causal ordering** when events must appear in sequence.
4. **Linearizable reads/writes** only where stale state is truly unacceptable.
5. **Bounded stale reads** where latency and scale matter more than immediacy.

A strong interview answer names different guarantees for different data:

- profile edits: read-after-write for the editor, bounded stale reads for others
- payments or quota: tighter coordination on the critical path
- analytics dashboards: lag is acceptable if it is explicit and observable
- abuse policy: stale reads may require fail-safe behavior

## Data Model & APIs

Useful metadata:

```text
record -> {
  value,
  version,
  committed_at,
  source_region,
  required_min_version
}
```

Helpful interface patterns:

- `Get(id, min_version)`
- `Put(id, expected_version, value)`
- `GetStrong(id)` for rare correctness-critical reads
- `GetBoundedStale(id, max_age_ms)` for cheaper paths

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| interviewer hears only vague consistency labels | design cannot explain user-visible behavior | restate guarantees as observable user outcomes |
| user sees stale state after own write | read-after-write mismatch metric rises | sticky session, leader read, or version-aware bypass |
| different replicas return conflicting snapshots | replica version skew widens | bounded stale contract, repair, or stricter routing |
| strict coordination is applied everywhere | latency and cost balloon | reserve strong guarantees for narrow critical paths |

## Observability

- metric: read-after-write mismatch rate for critical flows
- metric: replica version skew and freshness age
- metric: strong-read fraction versus total read volume
- log: stale-read incidents with entity, version, and read path
- trace: write request connected to later reads across regions
- SLO: freshness objective paired with latency target for each critical entity class

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| selective strong consistency | correctness where it matters | more complexity in a few flows | one global strict mode |
| bounded stale reads for noncritical data | lower latency and better scale | brief divergence is visible | always reading from the coordination path |
| explicit user-visible guarantees | clearer design communication | forces precise thinking | hand-waving with "eventual consistency" |

## Interview It

**Google framing:** "Design account settings and quota visibility for a global product." The signal is whether you separate critical versus convenience data and explain what each user can observe.

**Cloudflare framing:** "Design globally replicated policy reads with low-latency enforcement." The signal is whether you reason about stale policy risk, routing, and degraded safety.

**Follow-ups:**
1. Which flows deserve leader reads after mutation?
2. What if the same user switches regions within seconds?
3. What if strong reads double tail latency?
4. What if the interviewer asks whether the system is "CP or AP"?
5. How would you summarize the consistency contract in one sentence?

## Ship It

- `outputs/consistency-guarantee-worksheet.md`

## Exercises

1. **Easy** - Define a read-after-write promise for profile edits.
2. **Medium** - Split a social app into strong, monotonic, and bounded-stale paths.
3. **Hard** - Explain the consistency contract for quota enforcement plus billing visibility across regions.

## Further Reading

- [Designing Data-Intensive Applications](https://dataintensive.net/) - practical language for consistency and replication guarantees
- [System design notes](https://github.com/liquidslr/system-design-notes) - baseline interview structure before going deeper on guarantees
