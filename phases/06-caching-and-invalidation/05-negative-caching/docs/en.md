# Negative Caching and Error Caching

> Not found and temporarily broken are not the same thing.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Use negative caching and short-lived error caching to protect origins without turning transient failures into persistent user-visible bugs.
**Prerequisites:** `02-freshness-models`, `10-reliability-retries-and-backpressure/01-timeouts-and-retries`
**Estimated time:** ~60 min
**Primary artifact:** negative caching checklist

## The Problem

When a system receives many repeated misses for the same absent key, the cache can protect the backend by remembering the absence. That is negative caching. But engineers often misuse the same mechanism for transient errors, turning a 2-second outage into a 2-minute stale failure.

This lesson separates:

- caching `not found`
- caching authorization or policy denials
- caching temporary upstream errors very carefully

The senior-level move is to cache absence and failure differently because they represent different truths.

## Clarify

- Is `not found` a stable answer or can the object appear shortly after?
- Are misses user-generated typos, abuse traffic, or expected lookups for eventually created objects?
- Which error classes are safe to cache briefly, if any?
- Is a wrong deny or stale error more harmful than extra origin traffic?

## Requirements

### Functional

- Reduce repeated backend lookups for stable absent objects.
- Preserve correct behavior when objects are created shortly after an initial miss.
- Avoid amplifying transient backend errors.

### Non-functional

- Protect backend capacity under abusive or accidental repeated misses.
- Keep negative and error caching policies explainable by status class.
- Limit stale-denial blast radius when policies or objects change.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Total read QPS | 90K req/s | shows the cacheable serving base |
| Repeated miss share | 18% of requests | enough to justify negative caching |
| Attack spike | 15x on random or guessed IDs | exposes backend protection value |
| Object creation delay | seconds to minutes | determines safe TTL for `not found` entries |
| Rough cost | backend miss traffic vs short-lived cache entries | frames the memory trade-off |

## Architecture

A practical policy usually looks like this:

- cache stable `404 not found` for a short or medium TTL
- cache `403 denied` only when policy changes are infrequent and the denial is safe to repeat
- avoid caching `500` by default
- if caching transient failures at all, use very short TTLs with explicit status tagging

Key idea:

```text
cache key + response class + policy version
```

That prevents a negative answer from silently outliving a major policy or data change.

## Data Model & APIs

Negative cache entry:

```text
key -> {
  result_class,   // not_found, denied, upstream_error
  cached_at,
  expires_at,
  version
}
```

Useful behavior:

- shorter TTL for `not found` when creation races are possible
- shortest TTL, or no caching, for upstream `5xx`
- explicit invalidation when object creation or permission changes occur

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| newly created object still returns cached `404` | create-after-miss mismatch counters | invalidate on create and shorten not-found TTL |
| temporary upstream outage becomes sticky | sustained cached-error serves after recovery | do not cache 5xx by default, or use sub-second TTL |
| denial cached across policy change | policy version mismatch on denied responses | include policy version in cache key or invalidate on policy rollout |
| random-ID abuse floods backend anyway | miss QPS remains high with low negative-cache hit ratio | rate limit abusive patterns and widen negative cache coverage where safe |

## Observability

- metric: negative-cache hit ratio by result class
- metric: cached error serve count and duration
- metric: create-after-negative-cache conflict rate
- log: policy-version mismatch and negative-entry invalidations
- trace: miss path annotated with result class and whether the cache served it
- SLO: backend protection without prolonged stale failure responses

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| cache `404` for stable misses | reduces repeated backend lookups | risk of stale absence after object creation | never caching misses |
| avoid caching `500` by default | limits sticky outage behavior | less backend shielding during incidents | caching all failures uniformly |
| version-aware deny caching | protects policy backends | more invalidation complexity | status-only caching with no policy context |

## Interview It

**Google framing:** "Design user lookup by username where many guessed usernames do not exist." The signal is whether you reduce miss load without making newly created accounts invisible.

**Cloudflare framing:** "Design config or route lookups at the edge where some paths are invalid and some failures are temporary." The signal is whether you separate absence, deny, and backend outage semantics.

**Follow-ups:**
1. What if usernames can be created immediately after a miss?
2. When is caching `403` more dangerous than caching `404`?
3. What if abusive clients intentionally spray random keys?
4. What if backend recovery is fast but cached errors linger?
5. How would you explain this policy to an on-call engineer at 3 a.m.?

## Ship It

- `outputs/negative-caching-checklist.md`

## Exercises

1. **Easy** — Choose a TTL for `404` on a mostly static catalog.
2. **Medium** — Design negative caching for user handles that can be created at any time.
3. **Hard** — Explain whether you would cache `403` and `500` in an edge policy system and why.

## Further Reading

- [HTTP response status codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status) — good grounding for separating response classes
- [Caching best practices](https://web.dev/http-cache/) — useful when reasoning about whether failures should be cached at all
