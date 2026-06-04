# Storage Deep-Dive Drill

> A strong storage answer sounds like prioritized access paths plus explicit lifecycle choices.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Practice turning an ambiguous product prompt into a storage plan that covers model choice, keys, indexes, hot/cold behavior, and deletion semantics under interview time pressure.  
**Prerequisites:** `01-storage-models`, `03-indexes`, `07-retention-and-deletion`  
**Estimated time:** ~60 min  
**Primary artifact:** drill worksheet + scoring rubric  

## The Problem

This drill compresses the whole phase into one timed exercise. The learner must show they can:

- identify primary access patterns
- choose a system of record
- justify indexes and derived views
- handle hot data, history, and deletion

The goal is not a perfect storage platform. The goal is a credible senior-level storage deep dive.

## Clarify

- What are the top two reads and top two writes?
- Which data needs strong correctness versus eventual derived views?
- What historical access, retention, or deletion constraints matter?
- Which query shapes are explicitly out of scope for the primary path?

## Requirements

### Functional

- Choose the primary storage model and key structure.
- Support at least one list or secondary query path.
- Explain lifecycle behavior for hot data, cold history, and deletion.

### Non-functional

- Keep the answer time-boxed and prioritized.
- Make cost and write amplification visible.
- Show operational maturity around observability and failure handling.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak serving QPS | 70K req/s | enough to force clear hot-path choices |
| Writes | 18K req/s | enough to make index cost and lifecycle workflows matter |
| Historical retention | 18 months | pushes cold storage and deletion discussion |
| Peak factor | 5x on a few tenants | exposes hotspot reasoning |
| Rough cost | primary store + derived indexes + lifecycle jobs | keeps the answer grounded in operational reality |

## Architecture

Recommended drill sequence:

1. Clarify workload and invariants.
2. Rank access patterns.
3. Choose the primary store and key layout.
4. Add only the indexes or derived systems the workload truly needs.
5. Close with retention, deletion, and observability.

## Data Model & APIs

A strong drill answer usually includes:

- one primary entity and key shape
- one list or lookup index
- one explanation of what is deliberately not first-class
- one delete or retention workflow

The answer should sound like a series of constraints, not a pile of storage products.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| answer chooses a store before naming workload | architecture sounds generic and tool-driven | force pattern ranking first |
| derived systems are treated as fully consistent | stale or divergent reads are ignored | name source of truth and lag expectations |
| hot data and cold history are not separated | cost and latency trade-offs stay fuzzy | add lifecycle tiers or recent-window strategy |
| deletion semantics are skipped | compliance and stale-data risks remain hidden | include tombstones, propagation, and proof of completion |

## Observability

- metric: request latency and volume by access pattern
- metric: write amplification or index maintenance cost
- metric: stale-delete or stale-index detection count
- log: rejected query shapes and lifecycle state transitions
- trace: primary store plus derived path participation
- SLO: top serving path remains fast while lifecycle workflows complete within policy targets

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one focused scenario | easier to compare attempts over time | less breadth per session | many disconnected mini-prompts |
| explicit lifecycle discussion | produces more realistic designs | adds pressure to the answer | ignoring retention and deletion until later |
| scoring rubric | creates repeatable feedback | feels stricter than brainstorming | vague "seems fine" review style |

## Interview It

**Google framing:** "Design storage for a shared notes product with search, history, and deletion requirements." The signal is whether you keep the primary model simple while acknowledging derived needs.

**Cloudflare framing:** "Design storage for customer configuration plus audit history and export requirements." The signal is whether you separate serving, history, and policy workflows cleanly.

**Follow-ups:**
1. What if one read path is now 90% of total traffic?
2. What if customers demand bulk export by date range?
3. What if deletion must complete within 15 minutes across all derived systems?
4. What if index write cost doubles after a new feature launch?
5. Which part of the storage answer would you deep dive if time allowed?

## Ship It

- `outputs/drill-worksheet-storage-deep-dive.md`
- `outputs/scoring-rubric-storage-deep-dive.md`

## Exercises

1. **Easy** — Run the drill for a bookmark manager with tags and deletion.  
2. **Medium** — Run the drill for a media library with blob metadata and archive retention.  
3. **Hard** — Run the drill for a multitenant metrics product with hot recent reads and long-term retention.  

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview pacing to pair with this storage-specific drill  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — broad reference for tying together storage, indexes, and lifecycle trade-offs  
