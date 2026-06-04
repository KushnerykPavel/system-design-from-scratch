# Cache Stampedes and Request Coalescing

> The miss path matters most when many callers miss together.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Prevent origin overload when hot keys expire or miss under bursty traffic by combining request coalescing, jitter, and bounded stale behavior.
**Prerequisites:** `01-cache-patterns`, `02-estimation-and-cost/05-burstiness`, `10-reliability-retries-and-backpressure/04-load-shedding`
**Estimated time:** ~75 min
**Primary artifact:** coalescing simulator + failure checklist

## The Problem

A cache can look healthy until a hot key expires. Then thousands of callers may stampede the origin at once, turning one miss into a localized outage.

This lesson focuses on the miss path under burst:

- request coalescing
- TTL jitter
- stale-while-revalidate behavior
- degraded-mode protection

The included Go helper simulates how many origin fetches occur with and without coalescing for the same burst.

## Clarify

- How many concurrent requests can hit the same key after expiry?
- Can the system serve slightly stale data while one refresh is in flight?
- Is one slow origin fetch enough to queue or timeout many callers?
- Which path fails first: origin CPU, DB, thread pool, or network slots?

## Requirements

### Functional

- Allow one refill to satisfy many waiting callers when possible.
- Bound origin fanout for the hottest keys.
- Explain degraded behavior when refreshes are slow or failing.

### Non-functional

- Preserve latency for most callers even during refill windows.
- Avoid self-inflicted origin overload after synchronized expiry.
- Make hot-key incidents visible to operators.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Total read QPS | 250K req/s | enough to make cache central to availability |
| Hot-key burst | 12K requests in 2 seconds | drives stampede risk |
| Origin fetch latency | 80-300 ms | determines concurrent in-flight pressure |
| TTL alignment risk | many keys expiring on exact minute boundaries | reveals synchronized refill danger |
| Rough cost | extra origin fanout vs stale-serving controls | frames mitigation choices |

## Architecture

A robust miss path often uses four layers of defense:

1. **Request coalescing** so one refresh fills many waiters.
2. **TTL jitter** so hot keys do not all expire together.
3. **Stale-while-revalidate** for paths where slight staleness is safer than overload.
4. **Admission control** when origin or refresh workers are already saturated.

The lesson artifact simulates coalescing to show the difference between:

- every concurrent request causing its own origin fetch
- one in-flight refresh serving the rest

## Data Model & APIs

Useful behavior:

- `Get(key)` returns cached value if fresh
- `GetOrRefresh(key)` joins an in-flight refresh when present
- `Refresh(key)` publishes one shared result to waiters

Useful metadata:

```text
key -> {
  value,
  expires_at,
  refreshing
}
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| hot key expires and all callers refill independently | origin QPS spikes sharply after expiry | request coalescing and stale serving |
| many TTLs align | periodic origin thundering herd | jitter TTLs and stagger warming jobs |
| origin fetch is slow so waiters pile up | refresh queue depth and waiter counts grow | timeout, fallback, stale-if-error, or shed refresh load |
| coalescing lock becomes a bottleneck | lock contention or long queue on one hot key | sharded coordination and bounded waiter caps |

## Observability

- metric: origin fetches per cache miss or per hot key
- metric: number of coalesced waiters by key class
- metric: stale serves during refresh or error
- log: refresh start, timeout, failure, and bypass reasons
- trace: one expired request linked to the shared refresh and waiting callers
- SLO: protect origin availability while keeping the hot read path within latency bounds

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| request coalescing | dramatically lowers origin fanout on hot misses | coordination overhead on miss path | every caller refills independently |
| stale-while-revalidate | preserves latency during refresh | bounded stale reads | blocking all readers until refresh completes |
| TTL jitter | prevents synchronized expiry | less deterministic expiry timing | fixed TTL boundaries for every key |

## Interview It

**Google framing:** "Design a cache in front of a product service where one hot item can get thousands of reads per second." The signal is whether you talk about miss-path collapse, not just steady-state hit ratio.

**Cloudflare framing:** "Design edge caching for a hot object under global traffic spikes." The signal is whether you cover request collapsing, refill safety, and stale-serving under origin distress.

**Follow-ups:**
1. What if the origin fetch itself takes 2 seconds?
2. What if stale reads are unacceptable for one class of data?
3. What if many different hot keys expire together?
4. What if one POP can coalesce locally but not across POPs?
5. How would you explain whether the cache incident was caused by expiry, origin slowness, or key skew?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/failure-checklist-cache-stampede.md`
- `outputs/interview-card-cache-stampede.md`

## Exercises

1. **Easy** — Compare origin fetch counts with and without coalescing for 100 simultaneous requests.
2. **Medium** — Add stale-while-revalidate and explain which endpoints should use it.
3. **Hard** — Redesign the strategy for multiple regions where each POP can only coalesce local traffic.

## Further Reading

- [Caching at Scale at Facebook](https://engineering.fb.com/2013/04/15/core-infra/scaling-memcache-at-facebook/) — practical examples of hot-key and refill pressure
- [System design notes](https://github.com/liquidslr/system-design-notes) — useful interview framing for explaining cache failure modes
