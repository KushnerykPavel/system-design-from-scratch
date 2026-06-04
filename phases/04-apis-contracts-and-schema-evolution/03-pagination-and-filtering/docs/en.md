# Pagination, Filtering, and Query Shape

> Query shape is part of the contract; bad defaults leak directly into latency, cost, and abuse risk.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Design list APIs that keep result navigation stable, filtering expressive, and backend cost predictable under growth.  
**Prerequisites:** `02-estimation-and-cost/07-bottleneck-math`, `04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events`  
**Estimated time:** ~75 min  
**Primary artifact:** query review checklist + interview card  

## The Problem

List endpoints look simple until scale arrives. Offset pagination becomes slow at deep pages, filtering turns into accidental full-table scans, and a flexible query API becomes an abuse surface.

Senior answers should treat pagination and filtering as a performance contract, not just frontend convenience.

## Clarify

- Is the user browsing recent results, jumping to arbitrary pages, or exporting large datasets?
- Which sort order is stable and meaningful to the product?
- Do filters map cleanly to indexed fields, or are they ad hoc?
- Must results be snapshot-consistent while data is mutating underneath?

## Requirements

### Functional

- Return large collections in chunks.
- Support the main product filters and sort order.
- Preserve enough state for the client to continue where it left off.

### Non-functional

- Avoid unbounded scans and expensive deep offsets.
- Keep pagination stable when concurrent writes happen.
- Prevent abusive or accidental expensive queries.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| List QPS | 35K req/s | shapes index pressure |
| Page size | 20-100 items | bounds response size and DB work |
| Export size | up to millions of rows | may require async workflow instead of normal pagination |
| Filter selectivity | 1% to 80% | determines whether indexes remain useful |
| Rough cost | DB reads + cache + search backend | query flexibility often becomes the hidden cost center |

## Architecture

Default pattern:

- use cursor pagination for user-facing feeds or time-ordered data
- keep offset pagination only for small or low-scale admin views
- restrict filters to indexed, product-approved dimensions
- move large exports into async jobs

```text
client -> list API -> query planner / validation -> indexed store or search tier
```

The most important design move is aligning the contract with the access path the storage layer can actually serve efficiently.

## Data Model & APIs

Preferred API shape:

```text
GET /items?cursor=...&limit=50&status=active&sort=created_at_desc
```

Cursor tokens often encode:

- last seen sort key
- tie-breaker ID
- optional filter signature

This prevents clients from mixing a cursor from one filter context into another.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| deep offset scans | tail latency climbs with page number | prefer cursor pagination |
| unstable ordering under writes | duplicates or skipped results reported | use stable sort key plus tie-breaker |
| flexible filters bypass indexes | high scan/read amplification | whitelist indexed filters only |
| export-style requests hit online path | long-running requests and timeouts | offload to async export job |

## Observability

- metric: latency by endpoint, page size, and query pattern
- metric: scan-to-result ratio
- metric: duplicate or skipped cursor complaints
- metric: rejected expensive query count
- log: normalized query shape and chosen access path
- SLO: common list queries stay within target latency without unbounded scans

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| cursor pagination | stable and efficient at scale | harder random page jumps | deep offset pagination everywhere |
| restricted filters | predictable performance | less query freedom | user-defined arbitrary predicates |
| async export path | protects online latency | more workflow complexity | serving huge exports through live list APIs |

## Interview It

**Google framing:** "Design a list API for an issue tracker." The signal is whether you map the query contract to indexable access patterns.

**Cloudflare framing:** "Design a query interface for logs or analytics metadata." The signal is whether you separate interactive queries from bulk exports and think about abuse limits.

**Follow-ups:**
1. What if product insists on jumping to page 10,000?
2. What if the sort order changes after launch?
3. What if clients combine filters that are individually safe but jointly expensive?
4. What if results must be snapshot-consistent for audit use cases?
5. When should search infrastructure replace direct database listing?

## Ship It

- `outputs/query-review-pagination-and-filtering.md`
- `outputs/interview-card-pagination-and-filtering.md`

## Exercises

1. **Easy** — Design the cursor fields for a reverse-chronological notifications list.  
2. **Medium** — Explain when offset pagination is still acceptable.  
3. **Hard** — Redesign an admin reporting API that currently times out on deep filtered pages.  

## Further Reading

- [Google API Design Guide](https://cloud.google.com/apis/design) — useful for list method conventions and filtering discipline  
- [GitHub REST pagination](https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api) — a practical comparison point for list APIs  
