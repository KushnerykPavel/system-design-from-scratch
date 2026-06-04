# Map Tiles and Read-Mostly Geo Data

> Map systems are a reminder that "low latency" sometimes means winning on immutability, cacheability, and rollout discipline rather than on a fancy write path.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Design a map-tile system that serves read-mostly geo data efficiently with immutable versions, CDN layering, and safe large-scale tile publishing.
**Prerequisites:** `06-caching-and-invalidation/06-cache-layers`, `13-multi-region-cdn-and-edge-traffic/03-cdn-layering`, `15-kv-cache-and-object-storage/03-object-storage`
**Estimated time:** ~60 min
**Primary artifact:** tile-policy validator + cache rollout checklist

## The Problem

Design a system that serves raster or vector map tiles worldwide. Reads dominate heavily, data updates are batchy or region-scoped, and users expect fast pan-and-zoom behavior. The system must keep tiles cacheable while still shipping data updates safely.

This lesson matters because many answers overcomplicate read-mostly geo systems. Senior answers emphasize immutable versions, zoom-level storage behavior, hot-region prewarm, and how to roll out new tile sets without cache confusion.

## Clarify

- Are we serving raster tiles, vector tiles, or both?
- How frequently does the underlying map data change?
- Do clients tolerate small regional rollout differences during a publish?
- Are tiles public cacheable assets or personalized responses?

If left open, assume globally cached vector tiles, mostly read-only access, region-based updates a few times per day, and immutable versioned tile assets published through a CDN.

## Requirements

### Functional

- Serve map tiles by version, zoom, and coordinate.
- Publish new tile versions safely.
- Support CDN and browser caching for immutable assets.
- Allow selective region or layer updates without rebuilding everything unnecessarily.
- Provide metadata for tile-set version selection.

### Non-functional

- Keep p99 tile latency low from global edges.
- Prevent cache invalidation storms during tile publishes.
- Bound origin load during hot zoom-level demand.
- Make rollback easy if a bad tile set ships.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Tile requests | 2M req/s peak global | CDN and edge hit rate dominate cost |
| Hot zooms | 70% of reads at zoom 12 to 16 | prewarm and cache layering should focus there |
| Tile size | 20 to 80 KB vector average | bandwidth and browser cache behavior matter |
| Publish cadence | 6 partial region publishes/day | rollout mechanics matter more than raw write throughput |
| Cold-region share | 90% of tiles rarely requested | object storage and long-tail caching policy matter |

## Architecture

```text
tile build pipeline
  -> versioned tile manifests
  -> object storage
  -> CDN / edge cache
  -> browser cache
  -> tile metadata API for version selection
```

Design notes:

1. Prefer immutable versioned tile paths so serving can use long TTLs.
2. Separate tile payload delivery from the lighter metadata that tells clients which version to request.
3. Prewarm only the hot zooms and hot regions where traffic concentration justifies it.
4. Roll forward and roll back by version reference, not by deleting mutable cached paths everywhere.

## Data Model & APIs

Core records:

```text
tile_set(version, region, layer, created_at, state)
tile_manifest(version, zoom_range, shards, checksum)
tile_request(z, x, y, version)
```

Useful interfaces:

- `GET /tiles/{version}/{z}/{x}/{y}.mvt`
- `GET /v1/tile-manifest?region=...`
- `PublishTileSet(version, region, layers[])`
- `RollbackTileSet(version)`

Strong answers separate immutable tile blobs from the lighter control-plane decision about which version is active.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| bad publish corrupts one region's tiles | checksum mismatch and client error spikes | canary publish, manifest validation, and fast version rollback |
| cache invalidation stampede on publish | origin surge and edge miss-rate spike | immutable versioning and manifest flips instead of broad purge |
| one zoom level becomes overwhelmingly hot | per-zoom cache hit and origin offload metrics | targeted prewarm and multi-layer caching |
| tile build lag leaves regions inconsistent too long | manifest age and publish duration metrics | staged rollout and clear active-version tracking |

## Observability

- metric: edge hit rate by zoom, region, and version
- metric: origin requests, bytes served, and publish duration
- metric: tile error rate and checksum validation failures
- metric: active-version skew across regions
- log: publish approvals, manifest flips, and rollback actions
- trace: tile request through edge, origin, and object fetch on miss
- SLO: 99.9% of tile requests are served within target latency from an approved active version

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| immutable versioned tiles | strong cacheability and easy rollback | more storage for overlapping versions | mutable tile URLs with global purge |
| separate manifest/control plane | small fast version switch | another piece of state to operate | embedding active state inside every tile path rewrite |
| targeted prewarm | protects hot paths cost-effectively | misses can still happen in cold regions | prewarming the full world on every publish |

## Interview It

**Google framing:** "Design a global map tile service." Expect follow-ups on cacheability, version rollout, and hot zoom-level behavior.

**Cloudflare framing:** "Design a read-mostly geodata delivery system optimized for edge caching." Expect pressure on purge avoidance, regional rollout, and origin offload.

**Follow-ups:**
1. What changes if tiles become personalized by user layer settings?
2. How do you publish road-closure updates faster than the normal daily cadence?
3. What if one metro region is 100x hotter than rural regions?
4. How do you roll back a bad tile set without purging the whole CDN?
5. What changes if vector tiles are too large for mobile bandwidth goals?

## Ship It

- `outputs/cache-rollout-checklist-map-tiles.md`

## Exercises

1. **Easy** — Explain why immutable tile versions are usually better than mutable tile URLs.
2. **Medium** — Compare broad purge against manifest-based rollout.
3. **Hard** — Redesign the system when a subset of layers must update every few minutes but the rest remain daily.

## Further Reading

- [Cloudflare cache concepts](https://developers.cloudflare.com/cache/) — useful grounding for layered cache behavior
- [OpenStreetMap](https://www.openstreetmap.org/) — useful reference point for map-data update cadence and tile ecosystems
