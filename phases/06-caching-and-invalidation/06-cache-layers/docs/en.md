# CDN, Browser, and Edge Cache Layers

> A cache hierarchy only helps when each layer has a clear job.

**Type:** Learn
**Company focus:** Cloudflare
**Learning goal:** Design layered caching across browser, CDN, edge, and origin-adjacent caches without losing control of freshness or debugging.
**Prerequisites:** `02-freshness-models`, `05-negative-caching`, `13-multi-region-cdn-and-edge-traffic/03-cdn-layering`
**Estimated time:** ~75 min
**Primary artifact:** layered cache review card

## The Problem

Many systems use more than one cache:

- browser cache
- CDN or POP cache
- application or edge-worker cache
- origin-adjacent shared cache

The hard part is not adding more layers. It is deciding what each layer caches, how it revalidates, and where to look when users report stale or inconsistent results.

This lesson teaches you to treat layered caching as a hierarchy of responsibilities rather than a stack of magic speedups.

## Clarify

- Which responses are public, private, or tenant-specific?
- Which layer is closest to the user, and which is easiest to invalidate globally?
- Are responses immutable objects, mutable metadata, or personalized views?
- Do you need global consistency, regional freshness, or simply lower origin cost?

## Requirements

### Functional

- Cache safe public objects as close to the user as possible.
- Prevent private or user-specific data from leaking across tenants or browsers.
- Support revalidation or purge behavior for mutable objects.

### Non-functional

- Lower origin latency and egress cost.
- Keep debugging paths understandable across layers.
- Bound stale windows during content or config changes.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Global read QPS | 800K req/s | justifies edge and browser offload |
| Cacheable public share | 70% | determines how much the hierarchy can save |
| Mutable content share | 15% | drives purge and revalidation complexity |
| Regions / POPs | 200+ edge sites | makes propagation and hit ratio nontrivial |
| Rough cost | CDN bandwidth, purge traffic, origin misses | grounds the multi-layer choice |

## Architecture

A healthy layered design often looks like this:

1. **Browser cache** for immutable assets and conditional revalidation.
2. **CDN/edge cache** for public responses with cache-control and purge support.
3. **Edge compute or service cache** for small derived data or policy lookups.
4. **Origin-adjacent cache** for protecting the system of record from repeated misses.

Do not push the same object through every layer with identical policy unless there is a good reason. The strongest answers deliberately split:

- static assets
- dynamic public content
- personalized responses
- internal metadata

## Data Model & APIs

Useful control signals:

- `Cache-Control`
- `ETag`
- `Last-Modified`
- surrogate keys or purge tags
- tenant and authorization context in the cache key

The most dangerous bug class is key ambiguity. If the layer cannot distinguish public from personalized data, the cache hierarchy becomes a data-leak risk.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| private response cached publicly | cache-key audit and security incidents | vary by auth context or disable shared caching |
| global purge lags across POPs | purge propagation latency metric | versioned assets, purge confirmation, bounded TTL |
| browser caches stale object after CDN is fresh | user reports with old ETag or Age headers | clear browser directives, conditional revalidation |
| one layer masks another during debugging | inconsistent headers and cache status visibility | standard cache-debug headers per layer |

## Observability

- metric: hit ratio by layer and content class
- metric: purge propagation latency across regions
- metric: origin offload percentage and miss amplification
- log: purge requests, cache-key decisions, and revalidation outcomes
- trace: request annotated with which layers were hit, missed, or bypassed
- SLO: public-content latency plus purge/freshness objectives for mutable classes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| browser + CDN for immutable assets | cheapest and fastest path | harder client-side debugging | serving all assets from origin |
| edge cache for public dynamic content | large origin offload and lower tail latency | purge and revalidation complexity | no edge caching for mutable responses |
| strict no-store for personalized responses unless proven safe | lowers privacy risk | less cache efficiency | optimistic shared caching with weak keys |

## Interview It

**Google framing:** "Design caching for global document thumbnails and metadata." The signal is whether you separate immutable assets from mutable metadata and personalized permissions.

**Cloudflare framing:** "Design a global edge caching strategy for API and static traffic." The signal is whether you reason cleanly about cache keys, purge propagation, and safe layering.

**Follow-ups:**
1. Which content classes should never enter a shared cache?
2. What if a global purge must take effect within 30 seconds?
3. What if browsers keep stale responses after edge caches are updated?
4. What if origin egress becomes the dominant cost?
5. What headers would you expose to debug cache behavior safely?

## Ship It

- `outputs/layered-cache-review-card.md`

## Exercises

1. **Easy** — Classify static assets, product pages, and personalized dashboards across cache layers.
2. **Medium** — Design purge strategy for dynamic public content with regional traffic skew.
3. **Hard** — Explain how you would prevent data leakage in a multitenant API behind a CDN.

## Further Reading

- [HTTP Caching](https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching) — essential cache-control and revalidation semantics
- [Caching best practices](https://web.dev/http-cache/) — useful browser and CDN behavior guide
