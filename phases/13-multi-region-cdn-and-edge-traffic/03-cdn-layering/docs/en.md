# CDN Layering and Cache Hierarchies

> One cache layer is a performance feature. A cache hierarchy is a control system for latency, origin protection, and correctness drift.

**Type:** Learn
**Company focus:** Cloudflare
**Learning goal:** Explain how browser, edge, shield, and origin-adjacent caches interact, and choose where freshness and invalidation complexity should live.
**Prerequisites:** `06-caching-and-invalidation/06-cache-layers`, `06-caching-and-invalidation/07-cache-consistency`, `13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive`
**Estimated time:** ~75 min
**Primary artifact:** cache-layer trade-off matrix

## The Problem

Design a global content delivery path for static assets, APIs, or dynamic fragments. The system needs low latency worldwide, strong origin protection, and predictable invalidation behavior.

Interview answers often stop at "use a CDN." Stronger answers explain the cache hierarchy: browser cache, edge POP, regional shield, and origin, plus who owns invalidation and stale behavior.

## Clarify

- Is the workload mostly immutable static assets, semi-static HTML, or dynamic API responses?
- What freshness guarantees do users actually notice?
- Are purge events rare and explicit, or frequent and automated?
- Is origin protection more important than the last few milliseconds of latency?

If the interviewer is vague, assume mixed content with a heavy static majority and a smaller set of purge-sensitive dynamic responses.

## Requirements

### Functional

- Serve cacheable traffic from the nearest practical layer.
- Protect origins from global request bursts and cache misses.
- Support explicit purge or version-based invalidation.

### Non-functional

- Keep edge hit ratio high enough to reduce origin load materially.
- Bound stale-content exposure after purge or content update.
- Avoid making invalidation control more expensive than the traffic it saves.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Global request volume | 2B/day | shows why origin protection matters |
| Cacheable share | 85% static or semi-static | determines whether hierarchical caching pays off |
| Purge events | 50K/day | invalidation frequency changes control-plane needs |
| Origin fan-in | 300+ POPs to a few origins | shield layers can collapse miss traffic |
| Rough cost | POP cache memory + shield egress + purge plane | layering trades latency for origin savings |

## Architecture

```text
browser cache
  -> edge POP cache
     -> regional shield / tiered cache
        -> origin-adjacent cache
           -> origin services + object storage
```

Why multiple layers:

1. Browser cache eliminates repeated user fetches.
2. Edge POP cache keeps latency low near the user.
3. Shield or tiered cache reduces duplicate misses toward origin.
4. Origin-adjacent cache smooths expensive backend fetches.

## Data Model & APIs

Core entities:

- `CacheKey`
- `FreshnessPolicy`
- `PurgeRequest`
- `VariantRule`
- `ShieldPolicy`

Useful APIs:

- `ResolveCacheKey(request)`
- `PurgeByTag(tag)`
- `SetCachePolicy(route, ttl, stale_policy)`
- `ExplainCacheDecision(request_id)`

The key API design concern is whether invalidation is explicit by object, by tag, by prefix, or mostly versioned in URLs.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| low edge hit ratio despite high cacheable share | miss spikes and high origin QPS | fix cache keys, TTLs, and Vary behavior |
| purge storms overload the control plane | purge backlog and delayed invalidation | batch purges, tag content, and prefer versioned assets where possible |
| shield becomes a bottleneck | shield saturation or elevated miss latency | shard shields and allow controlled bypass |
| stale or personalized content leaks | cache-key anomalies or user reports | separate personalized routes and tighten cache-key rules |

## Observability

- metric: hit ratio per layer and per content class
- metric: origin request collapse rate from tiered cache
- metric: purge backlog, propagation delay, and stale serve duration
- log: sampled cache-decision explanations including layer, key, and freshness reason
- trace: request path through browser, edge, shield, and origin on misses
- SLO: purge-sensitive objects should clear globally within an explicit propagation target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| tiered cache or shield | dramatically lowers origin fan-in | adds another layer to debug | every POP misses straight to origin |
| versioned static assets | simple correctness and cheap invalidation | changes deployment workflow | frequent hard purges for everything |
| caching dynamic fragments selectively | lowers latency for semi-static pages | higher cache-key complexity | no caching for anything personalized |

## Interview It

**Google framing:** "How would you reduce origin load for a globally used service?" Expect questions about invalidation and what really belongs in cache.

**Cloudflare framing:** "Design the cache hierarchy for a global CDN-backed application." Expect pressure on cache keys, shield layers, and stale behavior.

**Follow-ups:**
1. What changes if most traffic is API reads instead of static assets?
2. How do you keep personalized responses out of shared caches?
3. What if purge must propagate globally within 30 seconds?
4. When is tiered cache worth the extra hop?
5. What changes if origin egress cost dominates?

## Ship It

- `outputs/tradeoff-matrix-cdn-layering.md`
- `outputs/observability-checklist-cdn-layering.md`

## Exercises

1. **Easy** — Sketch a cache hierarchy for a static asset-heavy web app.
2. **Medium** — Redesign for semi-dynamic HTML with per-country variations.
3. **Hard** — Support frequent purges from a CMS without overloading the purge plane.

## Further Reading

- [Cloudflare Learning Center - What is tiered caching?](https://www.cloudflare.com/learning/cdn/tiered-caching/) — practical vocabulary for shield layers
- [MDN - HTTP caching](https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching) — concise grounding for cache-control semantics
