# Cost, Latency, and Edge Cache Trade-offs

> Mature edge designs are not only fast. They are explicit about which milliseconds are worth buying, which misses are worth paying for, and which cache layers only look cheap until invalidation and egress show up.

**Type:** Learn
**Company focus:** Cloudflare
**Learning goal:** Reason explicitly about edge cost, latency, cache hit rate, and origin protection trade-offs instead of treating edge expansion as free performance.
**Prerequisites:** `02-estimation-and-cost/03-bandwidth-and-egress`, `06-caching-and-invalidation/06-cache-layers`, `22-cloudflare-edge-platform-design/02-cdn-invalidation`
**Estimated time:** ~60 min
**Primary artifact:** trade-off worksheet + cost review prompts

## The Problem

You need to explain whether a proposed edge design is actually worth its cost. The system can add more POP caching, more shielding, more routing flexibility, or stronger purge behavior, but each choice changes egress, origin load, stale risk, and operational complexity.

This lesson is about making cost and latency first-class parts of a design answer, especially for edge systems where a small miss-rate change can move a lot of traffic and money.

## Clarify

- What is the most valuable latency improvement: median, tail, or cross-region fallback latency?
- Which traffic classes are cacheable enough to matter financially?
- Is the bigger business pain origin overload, user latency, or egress cost?
- How much control-plane complexity is acceptable to gain more cache efficiency?

## Requirements

### Functional

- Compare candidate cache and routing strategies.
- Estimate how hit-rate changes affect origin load and egress.
- Explain where shielding helps beyond direct caching.
- Tie cost decisions to product-level value, not only infrastructure totals.

### Non-functional

- Avoid presenting edge footprint growth as free.
- Make invalidation and control-plane cost visible.
- Identify which latency gains matter to users versus operators.
- Keep trade-offs concrete enough for interview discussion.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Global request rate | millions QPS | small percentage shifts become huge absolute cost |
| Cache hit rate delta | 1% to 5% | can materially change egress and origin load |
| Median object size | KB to MB depending on product | changes savings per hit |
| Miss penalty | tens to hundreds of ms | affects user experience and origin saturation |
| Regional fallback factor | 2x traffic during incidents | cost and latency move together in failure mode |

## Architecture

Think in three layers:

```text
request path cost =
  edge compute
  + cache metadata/control-plane overhead
  + miss egress and origin cost
  + incident-time failover multiplier
```

Useful comparisons:

1. more POP caching versus stronger regional shield
2. lower TTL versus heavier purge control plane
3. premium low-latency routing versus wider shared efficiency
4. broader cacheability versus correctness and freshness risk

## Data Model & APIs

Useful worksheet fields:

```text
edge_option(
  option_name,
  expected_hit_rate,
  origin_qps_delta,
  egress_delta,
  latency_delta,
  control_plane_complexity
)
```

Helpful APIs:

- `EstimateHitRateBenefit(option)`
- `EstimateOriginSavings(option)`
- `CompareEdgeOption(a, b)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| design optimizes latency but ignores purge cost | control-plane spend and lag rise | include invalidation cost in option review |
| hit-rate gain is overestimated | real miss rate stays high | measure by route class, not global average only |
| cost savings harm freshness too much | stale-hit complaints rise | define freshness budget before expanding caching |
| failover path cost is ignored | incident egress spikes unexpectedly | model normal and incident cost separately |

## Observability

- metric: hit rate by route class and POP tier
- metric: origin QPS avoided by cache layer
- metric: purge volume and propagation cost
- metric: egress cost proxies by path type
- log: option assumptions and realized deltas
- trace: cache hit versus shield miss versus origin fetch path

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| broader edge caching | lower latency and less origin load | more purge complexity and freshness risk | minimal caching everywhere |
| regional shield emphasis | better origin protection | extra hop and shield cost | only POP-level cache strategy |
| premium routing class | better latency for key traffic | higher network cost and policy complexity | one size fits all pathing |

## Interview It

**Google framing:** "How would you reason about latency versus infrastructure cost in a global serving system?" Strong answers quantify miss penalties and failure-mode multipliers.

**Cloudflare framing:** "When is more edge actually worth it?" Strong answers discuss hit rate, purge cost, shield value, and incident-time economics instead of hand-waving about performance.

**Follow-ups:**
1. What if hit rate improves only 1% but object sizes are large?
2. What if purge frequency doubles after a product launch?
3. What if origin cost is cheap but user latency is the top KPI?
4. What if failover traffic dominates monthly egress spikes?

## Ship It

- `outputs/tradeoff-worksheet-cost-latency.md`
- `outputs/cost-review-prompts-edge.md`

## Exercises

1. **Easy** — Explain why a 2% hit-rate gain can still matter materially.
2. **Medium** — Compare lower TTLs with heavier purge use for a fast-changing content product.
3. **Hard** — Redesign for a case where latency goals stay fixed but network spend must drop sharply.

## Further Reading

- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — useful for understanding why tail improvements can be worth disproportionate effort
- [Content delivery network](https://en.wikipedia.org/wiki/Content_delivery_network) — background on CDN behavior and cost drivers
