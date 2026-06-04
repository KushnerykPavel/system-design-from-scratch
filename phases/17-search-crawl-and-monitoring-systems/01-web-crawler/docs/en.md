# Web Crawler

> Crawl speed is easy to ask for and easy to get wrong; the hard part is freshness without becoming a bad internet citizen.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a large-scale crawler that balances frontier throughput, politeness, deduplication, and recrawl freshness instead of treating fetchers as a single giant queue.
**Prerequisites:** `07-queues-streams-and-workflows/03-consumer-groups`, `09-partitioning-sharding-and-rebalancing/02-hot-partitions`, `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`
**Estimated time:** ~75 min
**Primary artifact:** crawl-plan validator + design review prompt

## The Problem

Design a web crawler that continuously discovers pages, fetches them, respects robots and host politeness, and feeds an indexing pipeline. The system should handle huge URL volume, avoid duplicate fetches, and recover cleanly from fetcher failures.

This lesson matters because senior interviewers quickly move past "have workers fetch URLs." They want to hear how you shape per-host fairness, dedup state, crawl freshness, content-change detection, and how indexers are insulated from the messiness of the public web.

## Clarify

- Are we crawling the public web, a controlled partner set, or internal documents?
- Is freshness more important than total coverage?
- Do we need full-page render support for JavaScript-heavy pages, or only HTML fetch?
- Is the primary product web search, domain monitoring, or compliance archiving?

If the interviewer stays broad, assume HTML-first crawling of the public web, eventual coverage, stronger freshness for high-value domains, and asynchronous rendering for a limited subset.

## Requirements

### Functional

- Accept seed URLs and continuously expand discovered links.
- Enforce robots.txt and per-host politeness budgets.
- Deduplicate URLs and near-duplicate content.
- Schedule recrawls based on change rate and page importance.
- Emit parsed page content and metadata into indexing pipelines.

### Non-functional

- Sustain large crawl throughput without overloading individual hosts.
- Recover from fetcher crashes without losing frontier ownership.
- Keep high-priority domains fresh within a tighter recrawl window.
- Preserve auditability for blocked, skipped, or delayed fetch decisions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active URLs in frontier | 40B | requires partitioned frontier and compact dedup state |
| Fetch throughput | 8M fetches/min peak | shapes queueing, bandwidth, and politeness budgeting |
| Average page size | 600 KB | drives bandwidth, parsing cost, and storage pressure |
| High-priority recrawl set | 200M URLs | freshness policy cannot be uniform |
| Peak factor | 4x during scheduled partner crawls | bursts stress host-level fairness and worker leases |

## Architecture

```text
seed loader
  -> URL normalizer + dedup
  -> partitioned crawl frontier
  -> host politeness scheduler
  -> fetch workers
  -> parser + link extractor
  -> content dedup + change detector
  -> indexing and recrawl planner
```

Design notes:

1. Partition the frontier by host or host-hash so politeness is enforceable without global coordination.
2. Keep URL dedup separate from content dedup because redirects and canonicalization make them different problems.
3. Lease frontier work to fetchers so crashed workers do not permanently strand URLs.
4. Drive recrawl cadence from observed change rate, page importance, and crawl budget.

## Data Model & APIs

Core records:

```text
normalized_url
host_key
discovery_time
last_fetch_time
next_fetch_time
fetch_status
content_fingerprint
robots_policy_version
priority_tier
```

Useful interfaces:

- `POST /v1/seeds {url, priority}`
- `ClaimFrontierBatch(host_key, limit, lease_seconds)`
- `ReportFetchResult(url, status, discovered_links, content_fingerprint)`
- `UpdateRecrawlPolicy(url, next_fetch_time, reason)`

The crawler should make crawl policy explicit. "Why was this URL not fetched?" is often an operational question, not a debugging afterthought.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one domain dominates frontier partitions | host skew metrics and queue-age spread | per-host token buckets and priority caps |
| crashed fetchers strand leased URLs | lease expiry backlog and stale claim age | expiring leases plus reclaim workers |
| robots policy changes but old workers keep crawling | robots version drift and policy violation audit | short-lived policy cache and policy-version checks |
| parser floods frontier with duplicate URLs | dedup hit-rate drop and frontier growth anomaly | stronger normalization, bloom filters, and domain caps |

## Observability

- metric: fetch success rate by host tier and response class
- metric: frontier age percentiles by priority tier
- metric: dedup hit rate for discovered URLs and content fingerprints
- metric: robots disallow count and politeness delay utilization
- log: fetch decision with URL, host, lease ID, robots verdict, and retry reason
- trace: claim to fetch to parse to enqueue pipeline for sampled URLs
- SLO: 99% of high-priority URLs are recrawled within their freshness target over a rolling day

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| host-partitioned frontier | easier politeness and fairness | less even worker distribution under skew | pure global FIFO queue |
| lease-based frontier claims | clean recovery from worker loss | duplicate work is still possible after timeouts | sticky ownership without expiry |
| freshness tiers by value | better product utility per fetch dollar | extra scheduling complexity | uniform recrawl interval |

## Interview It

**Google framing:** "Design the crawl system behind a search engine." Expect pressure on crawl budget, duplicate suppression, and which pages deserve frequent recrawls.

**Cloudflare framing:** "Design an internet-scale fetch platform for monitoring site changes safely." Expect questions on distributed fairness, noisy domains, and how you keep bad hosts from dominating the system.

**Follow-ups:**
1. What changes if some pages require headless rendering?
2. How would you crawl fast-moving news sites more often without starving the long tail?
3. How do you handle canonical tags and redirect chains in dedup?
4. What if legal policy requires honoring takedown requests within minutes?
5. How do you isolate one abusive or broken domain from your fleet?

## Ship It

- `outputs/design-review-web-crawler.md`

## Exercises

1. **Easy** — Sketch the normalization rules you would apply before URL dedup.
2. **Medium** — Redesign the frontier when one host suddenly produces 20% of all discovered URLs.
3. **Hard** — Add a dual-pipeline design where HTML fetch and JS rendering have different budgets and freshness policies.

## Further Reading

- [Google Search Central: How Search works - crawling](https://developers.google.com/search/docs/fundamentals/how-search-works) — useful grounding for crawl-discovery concepts
- [The Architecture of Open Source Search Engines](https://nlp.stanford.edu/IR-book/html/htmledition/web-crawling-and-indexes-1.html) — solid background on frontier and duplicate challenges
