# Netflix Interview Rubric & Strong Signals

> Generic system design answers get you to the next round. Netflix-specific signals get you an offer.

**Type:** Concept  
**Company focus:** Netflix  
**Learning goal:** Understand what Netflix actually evaluates in a system design loop — culture of chaos, extreme microservices scale, streaming data pipelines, and an experimentation mindset. Recognize the difference between a strong-hire answer and a competent generic answer.  
**Prerequisites:** `10-reliability-retries-and-backpressure/03-circuit-breakers`, `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`  
**Estimated time:** ~75 min  
**Primary artifact:** rubric card + strong-hire signal checklist  

## The Problem

You are preparing for a Netflix senior/staff system design interview. Netflix's bar is high and highly specific. Interviewers are engineers who built the actual systems. Generic distributed-systems answers will not impress them — they need to see that you understand Netflix's operational reality: massive scale, intentional failure injection, polyglot microservices, and every decision backed by data from experiments.

## Clarify

- Which domain is the question in? Streaming pipeline, recommendation, experimentation, data infrastructure, or platform reliability?
- Is the focus on new-system design or on improving a known bottleneck?
- Does the interviewer want depth-first (one subsystem deeply) or breadth-first (end-to-end sketch)?

## Requirements

### What Netflix Tests

Netflix interviewers are assessing several dimensions simultaneously:

1. **Scale intuition** — Can you size the problem correctly? Netflix streams to 270M+ subscribers across 190 countries. At peak, hundreds of terabits per second leave their CDN. Numbers matter.
2. **Microservices decomposition** — Can you decompose a system into independently deployable services with clear ownership? Netflix runs 1000+ microservices. Monolith answers are a yellow flag.
3. **Failure-first thinking** — Do you design for failure before you design for the happy path? Netflix invented Chaos Monkey precisely because engineers assumed too much reliability.
4. **Data-driven decisions** — Can you instrument your design for A/B testing? Every major Netflix decision is tested. Answers that skip experimentation show a gap.
5. **Streaming pipeline familiarity** — Do you understand Kafka, Flink, or event-driven architectures? Netflix's data plane is event-driven end to end.

### Non-functional Signals Netflix Cares About

- Availability over consistency for most user-facing paths (a stale recommendation list is better than no homepage).
- Latency percentiles, not just averages. Netflix targets p99 for playback start.
- Regional isolation — one region's failure should not cascade globally.
- Graceful degradation with explicit fallback chains.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Subscribers | 270M+ | drives recommendation and account system scale |
| Concurrent streams | ~15M peak | drives CDN, encoding, and bitrate-switching scale |
| Catalog size | ~36,000 titles, millions of encoded variants | drives storage and transcoding pipeline sizing |
| Events per second | billions per day (playback, clicks, impressions) | drives streaming data pipeline sizing |
| A/B experiments running simultaneously | 1,000+ | drives assignment service and isolation requirements |

## Architecture: What Strong Looks Like

### Weak-Hire Answer Pattern

- Draws a CDN in front of an origin server.
- Mentions S3 for video storage.
- Says "use load balancers for scaling."
- No failure modes, no fallback chains, no capacity numbers.
- "We can shard the database" as the answer to every scaling question.

### Strong-Hire Answer Pattern

- Starts with clarifying questions and a capacity estimate before touching architecture.
- Decomposes the system into services with clear ownership boundaries.
- Identifies the top 2-3 failure modes and names explicit fallbacks for each.
- Mentions at least one measurement or experimentation hook.
- Discusses trade-offs explicitly: "I chose X over Y because Z, at the cost of W."
- Uses Netflix-specific vocabulary appropriately: Open Connect, Chaos Monkey, EVCache, Zuul, Hystrix, Mantis.

### Strong-Hire Milestone Map

| Time mark | Expected progress |
|-----------|-------------------|
| 5 min | Clarified scope, named top 3 requirements, gave rough capacity numbers |
| 15 min | Sketched high-level architecture with named services, not just boxes |
| 25 min | Drilled into the hardest subsystem with data flow and failure modes |
| 35 min | Addressed the top interviewer follow-up (they will pivot you deliberately) |
| 45 min | Summarized trade-offs and named what you would instrument first |

## Netflix-Specific Vocabulary

| Term | What it signals |
|------|-----------------|
| Open Connect | Netflix's own CDN embedded in ISP networks — shows CDN depth |
| EVCache | Memcached-based distributed cache used pervasively — shows caching literacy |
| Zuul | API gateway layer — shows familiarity with edge routing |
| Hystrix/Resilience4j | Circuit breaker library — shows failure-first thinking |
| Chaos Monkey / Simian Army | Deliberate fault injection — shows operational maturity |
| Mantis | Real-time stream processing platform — shows data pipeline depth |
| Titanic | Netflix's A/B testing platform (also referred to as XP) — shows experimentation mindset |
| Hollow | Netflix's in-memory data propagation system — shows data distribution literacy |

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Recommendation service is slow | Latency spike on homepage load | Fallback to trending/popular list; cache last known-good result |
| CDN origin fetch fails | Cache miss errors spike | Retry with exponential backoff; serve stale if within TTL |
| A/B experiment leaks across buckets | Metric anomalies in holdout group | Assignment service audit log; deterministic bucketing verification |
| Kafka consumer lag grows | Lag metric alert | Scale consumer group; prioritize high-value event types |
| Region-level partial failure | Error rate and latency diverge by region | Route around via traffic steering; shed non-critical background work |

## Observability

- metric: playback start success rate by region and device type
- metric: recommendation API latency at p50/p95/p99
- metric: CDN cache hit ratio by content popularity tier
- metric: A/B experiment assignment coverage vs target
- log: fallback activation with reason and affected user count
- trace: user request from edge to recommendation to homepage render

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Availability over consistency for recommendations | Users always see a homepage | Stale data for a short window | Strong consistency would cause timeouts during service degradation |
| Microservices over modular monolith | Independent scaling and deployment | Operational complexity, distributed tracing overhead | Monolith easier to develop but cannot scale individual hot paths |
| Own CDN (Open Connect) over third-party | Cost control, deeper ISP integration, proactive fill | Massive capital investment, ongoing ISP relationship management | Akamai/CloudFront viable at smaller scale but cost prohibitive at Netflix volume |

## Interview It

**Netflix framing:** "Design Netflix's streaming system." Strong answers decompose into at least: ingestion/transcoding pipeline, CDN/delivery, recommendation, and experimentation. Weak answers stay at the CDN + database level.

**Follow-ups:**
1. What happens to the homepage when the recommendation service is down for 30 seconds?
2. How would you design the system to detect a bad video encode before it reaches 10M subscribers?
3. What is your fallback strategy if a CDN region goes dark mid-stream?
4. How does your design support running 500 simultaneous A/B experiments without contamination?
5. Walk me through how a playback event flows from the device into a recommendation signal.

## Ship It

- `outputs/rubric-card-netflix.md`
- `outputs/strong-hire-checklist-netflix.md`
- `outputs/interview-card-netflix-rubric.md`

## Exercises

1. **Easy** — List five Netflix-specific terms and explain the trade-off each one represents.  
2. **Medium** — Sketch a fallback chain for the Netflix homepage that degrades gracefully through three levels.  
3. **Hard** — Write a capacity model for the Netflix playback event pipeline: from device tap to Kafka topic to recommendation update.  

## Further Reading

- [Netflix Tech Blog](https://netflixtechblog.com/) — primary source for every Netflix-specific architecture decision  
- [Chaos Engineering book](https://www.oreilly.com/library/view/chaos-engineering/9781492043850/) — co-authored by Netflix engineers  
- [Netflix Open Source Software](https://netflix.github.io/) — Hystrix, Zuul, Eureka, and more  
