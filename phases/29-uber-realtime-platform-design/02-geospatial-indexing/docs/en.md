# Real-Time Geospatial — H3 Indexing & Proximity Search
> Finding the nearest driver is a proximity problem — H3 turns it into a set membership problem.

**Type:** Build
**Company focus:** Uber
**Learning goal:** Design a geospatial index that finds all available drivers within 2km of a rider request in <50ms at 5M concurrent driver locations.
**Prerequisites:** `20-low-latency-location-and-market-systems/01-proximity-service`, `14-rate-limiters-ids-and-hashing/04-consistent-hashing`
**Estimated time:** ~90 min
**Primary artifact:** H3 cell hierarchy diagram + proximity query design

---

## The Problem

Given 5 million active drivers each emitting GPS coordinates at 4 Hz:
1. How do you index their locations so that "find all drivers within 2km of point P" answers in <50ms?
2. How do you update the index at ~20M events/sec without write bottlenecks?
3. How do you shard the index across machines?

A naive `SELECT * WHERE distance(lat, lon, ?, ?) < 2000` on 5M rows with a B-tree index is O(N) for large radii and requires expensive trigonometry. The answer is **spatial partitioning**.

---

## H3: Hierarchical Hexagonal Grid

H3 was open-sourced by Uber in 2018. It divides the Earth's surface into hexagonal cells arranged in a hierarchy. Each cell has a fixed-length 64-bit ID encoding its resolution and position.

### Resolution Table

| Resolution | Edge length | Area | Use case |
|---|---|---|---|
| Res 15 | ~0.5m | 0.9 m² | Individual parking spaces |
| Res 12 | ~9m | 300 m² | Building footprint |
| Res 9 | ~174m | 0.1 km² | City block (driver lookup) |
| Res 7 | ~1.2km | 5.2 km² | Neighborhood (surge zones) |
| Res 5 | ~8.5km | 253 km² | District (supply forecasting) |
| Res 3 | ~59km | 12,393 km² | Metropolitan area |

### Why Hexagons?

In a square grid, the 4 edge-adjacent neighbors are at distance `d`, but the 4 diagonal neighbors are at distance `d√2 ≈ 1.41d`. This means a "radius-1 ring" search has non-uniform coverage.

In a hexagonal grid, all 6 neighbors are exactly `d` away from the center cell. A k-ring query returns cells whose centers are within k hops — uniform distance coverage.

```
Square grid k=1:          Hexagonal grid k=1:
  D  E  D                     E  E
  E  C  E              E  E  C  E  E
  D  E  D                  E  E
  (D = diagonal, 41% farther)   (all 6 E cells equidistant)
```

---

## H3 Proximity Query Design

### Data Structure

For each active H3 cell at res 9, maintain a Redis set of driver IDs:

```
Key:   h3cell:{cell_id}
Value: SET { "driver_101", "driver_204", ... }
TTL:   30s (auto-expire stale drivers)
```

### Proximity Query

To find all drivers within ~2km of a rider:

1. Compute the rider's H3 cell at res 9: `rider_cell = h3_cell(lat, lon, res=9)`
2. Compute k-ring(1): returns the center cell + its 6 neighbors = 7 cells total
3. For each of the 7 cells, execute `SMEMBERS h3cell:{cell_id}`
4. Union the results, filter for AVAILABLE status, compute ETAs

```
k-ring(1) = 7 cells × avg 10 drivers/cell = ~70 driver candidates
SMEMBERS is O(N) per cell — fast when cells are small
```

For a 4km search radius, use k-ring(2) = 19 cells.

### Why Not Redis GEORADIUS?

`GEORADIUS` is O(N + log M) where N is the number of results and M is the total number of elements in the sorted set. If you put all 5M drivers in one sorted set, a 2km radius query over a dense city returns thousands of results and scans a large portion of the set.

H3 cells bound the scan: each res-9 cell covers ~0.1 km². A 2km radius touches at most ~125 cells, but k-ring(1) gives 7 well-bounded cells with O(1) key lookup per cell. The tradeoff: H3 cells are hexagonal, not circular — but the approximation is excellent for dispatch.

---

## Driver Location Update Protocol

When a driver emits a GPS update at 4 Hz:

```
1. Compute new H3 cell: new_cell = h3_cell(new_lat, new_lon, res=9)
2. Read cached old cell from driver index: old_cell = GET driver:{driver_id}:cell
3. If new_cell != old_cell:
     SREM h3cell:{old_cell} driver_{id}
     SADD h3cell:{new_cell} driver_{id}
     SET driver:{driver_id}:cell new_cell
4. Always refresh TTL: EXPIRE h3cell:{new_cell} 30
```

Cell transitions are infrequent. At 4 Hz with ~174m cells, a driver moving at 40 km/h crosses a cell boundary every ~15 seconds — about once per 60 updates. So most updates are TTL refreshes only.

### Write Volume

| Operation | Rate |
|---|---|
| Total GPS events | 20M/sec |
| Cell boundary crossings | ~333K/sec (1/60 of events) |
| SADD + SREM per crossing | 666K ops/sec |
| TTL refresh (SET) | 20M/sec |

Redis pipelining and batching reduce round-trips. A single Redis node handles ~1M ops/sec; a cluster of 20+ shards is sufficient.

---

## Sharding the Cell Store

Shard Redis by H3 cell ID. H3 cell IDs are 64-bit integers with structure:
- Bits 63–60: resolution
- Bits 59–0: position encoding

A simple shard key: `shard = cell_id % num_shards`. This distributes cells evenly because H3 IDs are designed to be spatially diverse at the bit level.

Alternatively, shard by the first hex digit of the cell ID string (0–F = 16 shards). Geographic shards make it easy to route all queries for a city to the same shard set.

### Capacity Estimate

```
5M drivers × 1 cell membership × 8 bytes/driver_id = 40 MB raw data
Plus Redis set overhead: ~64 bytes/element → 5M × 64B = 320 MB
Plus driver→cell index: 5M × ~50 bytes = 250 MB
Total: ~650 MB — fits on a single large Redis instance; sharding is for throughput, not capacity
```

---

## Geofencing with H3

Surge zones, airport pickup zones, and restricted areas are defined as polygons. H3 provides `polyfill(polygon, resolution)` which returns all cells covering the polygon interior.

To test if a driver is in a surge zone:
1. Polyfill the surge polygon at res 9 → set of cell IDs
2. Store the set: `SADD surge_zone:{zone_id} {cell_id_1} {cell_id_2} ...`
3. Driver enters a cell: `SISMEMBER surge_zone:{zone_id} {driver_cell_id}` — O(1)

This turns polygon containment (expensive) into set membership (O(1)).

---

## Failure Modes

| Failure | Symptoms | Mitigation |
|---|---|---|
| Redis shard failure | Drivers in N cells invisible | Replica failover; dispatch gap accepted |
| Stale driver location | Driver moved but cell not updated | 30s TTL auto-expires; dispatch skip |
| H3 cell boundary artifact | Driver straddling two cells | k-ring covers both cells naturally |
| Thundering herd on cell | Popular pickup zone: one cell with thousands of drivers | Rare at res 9 (~0.1 km²); use res 10 if needed |
| Redis memory pressure | TTL expiry not fast enough | Monitor `used_memory`; evict with `allkeys-lru` as safety valve |

---

## Complete System Flow

```
Driver GPS (4 Hz)
  → Kafka topic: driver-location (partitioned by driver_id)
  → Location Consumer Service
      → compute H3 cell
      → if cell changed: SREM old cell, SADD new cell
      → SET driver:cell index
      → EXPIRE cell set TTL
  → Redis Cell Store (sharded by cell_id)

Rider Request
  → API Gateway
  → Dispatch Service
      → compute rider H3 cell
      → k-ring(1) → 7 cell IDs
      → Redis pipeline: SMEMBERS × 7
      → collect driver candidates
      → fetch driver status + ETA
      → rank and send to matching engine
```

---

## Design Checklist

- [ ] H3 resolution chosen: res 9 for dispatch, res 7 for surge zones
- [ ] Redis cell set structure: `h3cell:{cell_id}` → SET of driver IDs
- [ ] Proximity query: k-ring(1) = 7 cells, pipeline SMEMBERS
- [ ] Location update: SADD/SREM on cell change, TTL refresh always
- [ ] Sharding strategy: by cell_id prefix or modulo
- [ ] Geofencing: polyfill → SISMEMBER
- [ ] Failure handling: replica failover, 30s TTL for stale drivers
- [ ] Observability: cell membership size histograms, Redis op latency P99
