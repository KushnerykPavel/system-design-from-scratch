# URL Shortener

> The interesting part is not generating short codes; it is surviving read skew, abuse, and redirect latency at massive scale.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design a redirect-heavy backend that separates write-path correctness from read-path speed, and explain how alias generation, hot objects, and abuse controls shape the system.  
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/01-shard-key`, `14-rate-limiters-ids-and-hashing/03-unique-id-generator`, `15-kv-cache-and-object-storage/02-distributed-cache-cluster`  
**Estimated time:** ~75 min  
**Primary artifact:** redirect topology validator  

## The Problem

Design a URL shortener that lets users create shortened links and resolve them with very low redirect latency. The workload is asymmetric: write volume is modest, but read traffic can spike dramatically when one link gets shared widely.

This lesson is useful because many candidates treat it as a toy CRUD service. Senior answers instead focus on slug-generation guarantees, cache behavior, hot-link mitigation, expiration, analytics side effects, and what the redirect path is allowed to do synchronously.

## Clarify

- Are we optimizing for custom aliases, random generated aliases, or both?
- Is click analytics required on the synchronous redirect path, or can it be buffered asynchronously?
- Do links expire, support edits, or need abuse and malware screening before they become active?

If the interviewer leaves details open, assume generated aliases are the default, custom aliases exist for paid users, redirect latency matters more than immediate analytics visibility, and abuse screening is mandatory.

## Requirements

### Functional

- Create a shortened URL from a long destination.
- Resolve a short code to the current destination.
- Support optional custom aliases and expiration time.
- Record click events for analytics and abuse investigation.

### Non-functional

- Keep redirect p99 under 50 ms in-region.
- Avoid duplicate alias assignment under retries.
- Survive hot-link traffic spikes without melting the primary store.
- Reject malicious or policy-violating destinations before activation.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Link creates | 25K writes/s peak | shapes ID generation, idempotency, and write-store sizing |
| Redirects | 2.5M reads/s peak | pushes caching, hot-key mitigation, and CDN use |
| Average URL size | 180 bytes input, 10-byte code output | affects storage and validation cost |
| Data retention | 30B links total, 90 days click analytics hot | separates metadata store from analytics pipeline |
| Peak factor | 500x for viral links | one code can dominate fleet traffic |

## Architecture

```text
create client
  -> API gateway
  -> idempotent write service
  -> alias generator / uniqueness check
  -> URL policy scanner
  -> metadata store

redirect client
  -> edge / CDN
  -> redirect service
  -> cache
  -> metadata store
  -> async click event bus
  -> analytics + abuse systems
```

Design notes:

1. Keep redirect resolution simple: code lookup, policy check, redirect response.
2. Push analytics emission off the critical path unless the prompt explicitly requires synchronous billing.
3. Treat hot codes as cache problems first, not database-scaling proofs.
4. Support idempotency keys on create requests so retries do not mint multiple aliases.

## Data Model & APIs

Core metadata:

```text
short_code
destination_url
owner_id
created_at
expires_at
status
custom_alias
policy_version
```

Suggested APIs:

- `POST /v1/links {destination_url, custom_alias?, expires_at?, idempotency_key}`
- `GET /r/{short_code}`
- `PATCH /v1/links/{short_code} {destination_url?, expires_at?, status?}`
- `GET /v1/links/{short_code}/stats`

If custom aliases are allowed, the write path must clearly separate "name reservation" from "metadata write" so collision handling is explicit.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| viral link overloads one cache shard | top-code skew metrics and cache shard saturation | replicate hottest keys, edge caching, per-code request coalescing |
| alias generator hands out duplicates under retry | uniqueness violation alarms and idempotency mismatch logs | transactional uniqueness check plus idempotency key storage |
| malware list update marks active links unsafe | policy version drift and block-rate jump | serve interstitial or block page, re-scan active set asynchronously |
| analytics pipeline lags behind redirect traffic | click ingestion backlog and consumer lag | decouple redirect success from analytics completion |

## Observability

- metric: redirect p50, p95, and p99 latency by cache hit or miss
- metric: create-path uniqueness conflicts and idempotency replay rate
- metric: hottest short codes as share of total redirect traffic
- metric: abuse scan latency and blocked-link activation attempts
- log: redirect decision with code, cache source, policy verdict, and request region
- trace: create request across generator, scanner, and metadata store
- SLO: 99.9% of redirect requests succeed under the documented latency target for active non-blocked links

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| random generated aliases by default | fast writes and predictable entropy | not human-friendly | only custom aliases, which causes hotspot and namespace issues |
| async analytics on redirect | lower redirect latency and simpler scaling | analytics can lag | synchronous logging in the redirect path |
| cache-first redirect service | absorbs hot reads cheaply | invalidation and stale entry handling | direct database lookup on every click |

## Interview It

**Google framing:** "Design a globally used short-link service for internal docs and external sharing." Expect pushback on slug uniqueness, analytics side effects, and what happens during viral traffic.

**Cloudflare framing:** "Design a redirect-heavy backend with strong edge caching and abuse enforcement." Expect questions on cache invalidation, region routing, and safe blocking behavior.

**Follow-ups:**
1. What changes if premium users demand custom aliases with instant availability?
2. How would you expire links automatically without scanning the full keyspace?
3. What if click events drive billing and cannot be dropped?
4. How would you protect the system from link-flood abuse?
5. What changes at 10x redirect traffic with the same storage budget?

## Ship It

- `outputs/interview-card-url-shortener.md`

## Exercises

1. **Easy** — Size the alias space for 8-character base62 IDs at 10 years of growth.
2. **Medium** — Redesign the service so paid aliases are globally unique and instantly reserved.
3. **Hard** — Add per-link revocation that propagates through edge caches in under one minute.

## Further Reading

- [TinyURL's Design and Scalability](https://www.usenix.org/legacy/event/nsdi12/tech/full_papers/DeCandia.pdf) — not a URL shortener paper, but useful for key-value and hot-key intuition  
- [Caching at Scale at Facebook](https://engineering.fb.com/2013/04/15/core-infra/tao-the-power-of-the-graph/) — helpful for thinking about skew and read-heavy application backends  
