# Full-Loop Drill: Design a URL Shortener

> Frameworks matter only if you can use them under pressure.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Integrate the entire phase by running a full interview loop on a familiar product prompt and checking whether the answer includes clarification, sizing, architecture, depth, risks, and redesign.  
**Prerequisites:** `03-design-framework-and-timing/07-interviewer-moves`  
**Estimated time:** ~60 min  
**Primary artifact:** full-loop drill sheet + interview card  

## The Problem

A URL shortener is simple enough to fit in one session, but rich enough to test most of the interview habits in this phase: scope control, sizing, pacing, high-level architecture, targeted deep dive, wrap-up quality, and redesign agility.

The point is not to invent the fanciest shortener. The point is to run the process cleanly.

## Clarify

- Do we support only create-and-redirect, or also analytics, custom aliases, and expiration?
- What scale are we targeting for redirects and new link creation?
- Is read latency or write correctness the dominant requirement?
- What abuse, spam, or safety constraints matter?

If the interviewer stays vague, assume a core redirect service with heavy reads, moderate writes, and basic abuse protection.

## Requirements

### Functional

- Create a short link for a long URL.
- Redirect users quickly and safely.
- Support basic lifecycle controls such as expiration or deletion.

### Non-functional

- Low redirect latency.
- High availability for reads.
- Safe handling of malicious or abusive links.
- Reasonable cost under read-heavy load.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Redirect QPS | 100K peak | drives cache and stateless tier size |
| Create QPS | 2K peak | affects ID generation and write path |
| Average URL metadata | ~1 KB | shapes storage growth |
| Peak factor | 5x normal | affects cache and overload plans |
| Rough cost | dominated by redirects and egress | keeps the design realistic |

## Architecture

A strong loop for this prompt:

1. Clarify product scope and abuse expectations.
2. Size read-heavy traffic and redirect latency targets.
3. Present high-level design:
   - API tier
   - short-code generation
   - metadata store
   - cache
   - abuse or validation service
4. Deep dive on one area:
   - ID generation
   - cache hit path and fallback
   - abuse detection and quarantine
5. Wrap up with failure modes, metrics, and redesign options.

## Data Model & APIs

Core entities:

- `ShortLink`
- `RedirectEvent`
- `AbuseStatus`

Core interfaces:

- `POST /links`
- `GET /{shortCode}`
- `DELETE /links/{id}`

Deep-dive notes:

- idempotency for create retries
- cache TTL for redirect path
- safe handling of deleted or quarantined URLs

The code artifact checks whether a practice answer covered the full interview loop.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| cache miss storm on a popular link | redirect latency spike and origin read surge | hot-key protection and request coalescing |
| ID collisions or bad allocator behavior | create failures or duplicate code alarms | strong uniqueness checks and retry path |
| malicious redirect target abuse | abuse alerts or user reports | validation, quarantine, and blocklists |
| datastore outage | elevated redirect errors | cache fallback and graceful degradation for known entries |

## Observability

- metric: redirect latency p50/p95/p99
- metric: cache hit rate and origin fallback rate
- metric: create success rate and collision retries
- metric: abuse quarantine counts and false-positive investigations
- SLO: successful redirect rate under peak read traffic

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| cache-heavy redirect path | low read latency | invalidation and hot-key complexity | direct datastore reads on every redirect |
| random or snowflake-like ID generation | simple write path | collision or predictability concerns depending on scheme | sequential DB IDs exposed publicly |
| basic abuse checks in v1 | faster launch | weaker protection depth | fully synchronous heavy scanning on create path |

## Interview It

**Google framing:** "Design a tiny URL service." Expect follow-ups on scale, data model, cache behavior, and trade-offs between availability and write-path rigor.

**Cloudflare framing:** "Design a globally distributed redirect service." Expect edge-cache questions, abuse handling, and origin protection under very high redirect volume.

**Follow-ups:**
1. How would the design change at 10M redirect QPS?
2. What if custom aliases become the dominant create path?
3. How would you support regional deletion compliance?
4. What if malicious links become a major abuse vector?
5. How would you redesign for multi-region writes?

## Ship It

- `outputs/full-loop-drill-sheet-url-shortener.md`
- `outputs/interview-card-url-shortener-loop.md`

## Exercises

1. **Easy** — Deliver a 15-minute version of the URL shortener answer using the same phase structure.  
2. **Medium** — Re-run the drill with analytics and custom aliases added to scope.  
3. **Hard** — Redesign the shortener for global edge redirects with multi-region writes and compliance-driven deletions.  

## Further Reading

- [System design notes: URL shortener](https://github.com/liquidslr/system-design-notes) — useful baseline prompt for repeated drills  
- [Google SRE books](https://sre.google/books/) — helpful perspective for reliability, observability, and rollout thinking  
