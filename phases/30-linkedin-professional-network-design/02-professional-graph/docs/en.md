# Professional Graph — Connections & Degree Traversal

> 950 million nodes, 500+ connections each, and you need 2nd-degree suggestions in under 100ms.

**Type:** Build
**Company focus:** LinkedIn
**Learning goal:** Design LinkedIn's professional graph storage and traversal system that enables People You May Know, 2nd-degree connections, and "X works at Y" social proof queries.
**Prerequisites:** `05-storage-indexing-and-access-patterns/01-storage-models`, `09-partitioning-sharding-and-rebalancing/07-cross-shard-queries`
**Estimated time:** ~90 min
**Primary artifact:** graph schema + BFS traversal design

## The Problem

Design the data storage and query system that powers LinkedIn's professional graph. The graph must support:

1. **People You May Know (PYMK)** — suggest 2nd-degree connections (friends-of-friends not yet connected) ranked by shared connection count, company, and school overlap.
2. **Social proof** — "Alice and 12 of your connections work at Acme Corp" — requires fast aggregation over the viewer's 1st-degree connections.
3. **Degree display** — show whether a member is a 1st-degree, 2nd-degree, or 3rd-degree connection on every profile view.
4. **Privacy enforcement** — blocked members must never appear in graph traversal results; connection visibility settings must be respected.

## Clarify

- Is this focused on PYMK or the full professional graph?
- What is the acceptable latency for PYMK suggestions? (Target: <100ms p99)
- How stale can PYMK suggestions be? (Batch precomputation introduces up to 24h lag)
- Does the graph include follows (asymmetric) in addition to connections (symmetric)?
- How should privacy constraints (blocked members, private profiles) interact with traversal?

## Requirements

### Functional Requirements

- Store directed and undirected professional relationships: CONNECTED (symmetric), FOLLOWS (asymmetric), WORKED_AT (member → company), EDUCATED_AT (member → school), HAS_SKILL (member → skill).
- Return 1st-degree connections for a member in <10ms.
- Return 2nd-degree PYMK candidates for a member in <100ms.
- Enforce privacy: blocked members excluded from traversal results, connection visibility respected.

### Non-functional Requirements

- 950M member nodes, ~475B edges (950M × 500 avg connections).
- ~50B new edges per year (connections + follows).
- Read QPS: ~10M/s for 1st-degree lookups (profile views, degree badges).
- Write QPS: ~500/s for new connections (mutual acceptance).

## Capacity Model

| Dimension | Estimate | Detail |
|-----------|----------|--------|
| Member nodes | 950M | one document per member in Espresso |
| Company nodes | 67M | separate keyspace |
| Avg connections per member | ~500 | symmetric; stored as adjacency list |
| Total edges (connections only) | ~475B | 950M × 500 / 2 × 2 (bidirectional) |
| Adjacency list size per member | ~5KB | 500 connections × 10 bytes (ID + timestamp + type) |
| Total adjacency list storage | ~4.75TB | 950M × 5KB |
| New connection writes/day | ~50M | both directions written atomically |

## Architecture

### Graph Model

Nodes:
- **Member** — member_id (uint64), name, headline, visibility_setting
- **Company** — company_id (uint64), name, industry
- **School** — school_id (uint64), name
- **Skill** — skill_id (uint32), name, normalized_name

Edges:
- **CONNECTED** (member ↔ member) — connection_time (unix), connection_type (direct | imported), block_flag
- **FOLLOWS** (member → member/company) — follow_time (unix)
- **WORKED_AT** (member → company) — start_year, end_year (nullable)
- **EDUCATED_AT** (member → school) — start_year, end_year
- **HAS_SKILL** (member → skill) — endorsement_count

### Storage Choice: Espresso Adjacency List

LinkedIn chose Espresso (their horizontally scalable document store built on MySQL) rather than a native graph database.

**Why not a graph database (Neo4j, JanusGraph)?**
- Graph databases store edges in a way that optimizes multi-hop traversal on a single node.
- At 950M nodes with 475B+ edges, no commercial graph database provides the throughput LinkedIn needs.
- The graph is not deeply traversed: PYMK only needs depth-2. An adjacency list with application-level BFS is sufficient.

**Why Espresso?**
- Horizontal scaling: each member's adjacency list is a self-contained document, partitioned by member_id.
- Sub-10ms read latency for single-key lookups (fetch one member's connections).
- Operationally proven at LinkedIn scale.

**Schema:**

```
Keyspace: member-connections
Key: member_id (uint64)
Value: {
  connections: [
    { peer_id: uint64, connection_time: int64, type: string, block_flag: bool }
  ],
  updated_at: int64
}
```

### BFS Traversal Design

**Level 1 (1st-degree) — real-time:**

```
fetch document for member_id → return connections list
```

Latency: ~5ms (Espresso single-key read).

**Level 2 (2nd-degree) — real-time with cardinality cap:**

```
1. fetch connections for member_id → set L1 (up to ~500 peers)
2. for each peer in L1:
     fetch connections for peer → set L2_candidates (up to 500 per peer)
3. union L2_candidates, deduplicate
4. exclude members already in L1 (already connected)
5. exclude member_id itself
6. exclude blocked members (check block_flag at read time)
7. return deduplicated candidates (capped at 10K for hot members)
```

Latency: ~50–100ms for average members (500 Espresso reads in parallel batches).

**PYMK Pipeline: Batch Precomputation**

Real-time BFS is only feasible for members with ~500 connections. For the system as a whole (30M DAU × 500 connections × 500 = 7.5T operations/day), batch precomputation is essential.

```
Offline (Spark job, runs daily via Azkaban):
  1. Load full adjacency list from Espresso snapshot to HDFS
  2. For each member: compute 2nd-degree candidates
  3. Score candidates: shared_connections_count × 3 + same_company × 5 + same_school × 4
  4. Keep top-100 candidates per member
  5. Write results to Espresso PYMK keyspace: member_id → [candidate_id, score]

Online (PYMK API, runs per request):
  1. Fetch precomputed candidates from Espresso PYMK keyspace (~5ms)
  2. Fetch member's current 1st-degree connections (~5ms)
  3. Exclude already-connected candidates (real-time filter)
  4. Fetch Venice features: shared_connections, same_company, same_school
  5. Re-score and return top-10
```

### Graph Partitioning

Espresso partitions member documents by member_id via consistent hashing. This means:

- A member's full adjacency list is on one shard — fast for single-member reads.
- BFS level 2 requires reads across many shards (each peer on potentially different shard) — this is parallelized.
- **Hot shard problem**: a member with 30K+ connections has a large document. Mitigation: read-only replicas for high-degree members.

### Privacy Enforcement

- **Blocked members**: block_flag stored on the edge; checked at traversal time before returning candidates.
- **Private profiles**: members with private visibility are excluded from 2nd-degree results.
- **GDPR deletion**: deleting a member requires removing their document and updating all peer adjacency lists that reference them (tombstone write, propagated via Kafka).

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| BFS timeout on highly-connected member | PYMK API latency p99 spikes; timeout errors | Cap BFS at 10K 2nd-degree candidates; precompute PYMK for top-1% connected members |
| Hot Espresso shard for celebrity member | QPS spike on single shard; increased latency | Read-only replica of large adjacency list documents |
| Stale PYMK shows already-connected member | Member reports seeing existing connection in PYMK | Real-time 1st-degree exclusion filter applied after precomputed candidate fetch |
| Batch Spark job fails | PYMK suggestions stale beyond 24h | Alerting on Azkaban job failure; fallback to real-time BFS for active members |
| Block not honored in traversal | Blocked member appears in PYMK or degree badge | block_flag check is applied at Espresso read time, not at batch time; tested in integration suite |

## Observability

- metric: PYMK API latency p50/p95/p99 broken down by precomputed vs real-time path
- metric: BFS traversal depth distribution (% of requests hitting depth-2 cap)
- metric: Espresso read latency per keyspace (connections vs PYMK candidates)
- metric: Azkaban Spark job duration and success rate
- metric: PYMK cache hit rate (precomputed candidate freshness)
- log: BFS traversal with member_id, depth reached, candidate count before/after exclusion
- alert: PYMK p99 latency > 150ms for 5 consecutive minutes
- alert: Azkaban PYMK Spark job not completed within 26 hours

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Espresso adjacency list over graph DB | Horizontal scaling; proven at LinkedIn scale; <10ms single-key reads | BFS is application-level; no native graph query language | Graph DB (Neo4j) cannot handle 950M nodes + 475B edges at required QPS |
| Batch PYMK precomputation over real-time BFS | Avoids 7.5T operations/day; <10ms candidate fetch | Up to 24h staleness; requires Spark + Azkaban infra | Real-time BFS is infeasible at 30M DAU with 500+ connections per member |
| Bidirectional edge storage (both directions per Espresso doc) | O(1) connection lookup from either side; no join required | 2× storage for symmetric CONNECTED edges | Single-direction + join is slower and adds cross-shard reads |
| Cap 2nd-degree BFS at 10K candidates | Bounded latency even for super-connectors | Some 2nd-degree candidates never seen for high-degree members | Unbounded BFS has unpredictable latency; worst case hours for 30K-connection members |

## Interview It

**LinkedIn framing:** "Design People You May Know." Strong answers decompose into: graph storage (Espresso adjacency list), BFS traversal (real-time for 1st-degree, batch precomputed for 2nd-degree), PYMK scoring pipeline (Venice features + ranking), and privacy enforcement (block_flag at traversal time).

**Follow-ups:**

1. A new member connects with 50 people on their first day. How quickly will they appear as a PYMK suggestion for 2nd-degree connections?
2. How do you handle a member who blocks another: which queries must immediately reflect the block?
3. Your Azkaban PYMK batch job fails. What is your fallback strategy?
4. How do you design the "X and 5 of your connections work at Acme" social proof query?
5. A super-connector with 30K connections posts a job. How does their PYMK list differ from an average member's?

## Ship It

- `outputs/graph-schema-linkedin.md`
- `outputs/bfs-traversal-design-linkedin.md`
- `outputs/pymk-pipeline-diagram.md`

## Exercises

1. **Easy** — Write the Espresso document schema for a member's adjacency list including all edge types.
2. **Medium** — Design the real-time BFS traversal for 2nd-degree connections with privacy enforcement and cardinality capping.
3. **Hard** — Design the full PYMK batch precomputation pipeline from Espresso snapshot ingestion to scored candidates in Espresso, including failure handling.

## Further Reading

- [LinkedIn's People You May Know pipeline](https://engineering.linkedin.com/blog/2021/optimizing-people-you-may-know) — PYMK optimization deep dive
- [Espresso: LinkedIn's Distributed Document Store](https://engineering.linkedin.com/espresso/espresso-linkedins-distributed-document-store) — storage architecture
- [Graph-Based People Recommendations at LinkedIn](https://engineering.linkedin.com/recommendations/graph-based-people-recommendations) — recommendation algorithms
- [Azkaban Workflow Scheduler](https://github.com/azkaban/azkaban) — batch job orchestration
