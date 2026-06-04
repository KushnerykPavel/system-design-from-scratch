# Pub/Sub Fanout and Subscription Isolation

> Fanout systems look simple until one slow subscriber quietly turns shared infrastructure into a retention and latency incident.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Design a pub/sub platform where topic fanout, per-subscription progress, filtering, replay, and subscriber isolation are explicit design choices.
**Prerequisites:** `07-queues-streams-and-workflows/03-consumer-groups`, `11-observability-slos-and-debugging/05-alert-design`, `16-application-backends/06-fanout-patterns`
**Estimated time:** ~75 min
**Primary artifact:** pubsub-plan validator + trade-off matrix

## The Problem

Design a publish/subscribe platform that distributes event streams to many independent subscribers. Some subscribers need all events immediately, some apply filters, and some may be offline for hours. The system must provide topic-scale fanout without allowing one lagging or misbehaving subscription to degrade everyone else.

This lesson matters because pub/sub answers often confuse topic durability with subscription delivery. Senior candidates talk about independent subscriber state, backlog budgets, filtered delivery cost, and how replay interacts with storage tiers.

## Clarify

- Are subscribers internal services, external tenants, or both?
- Is the system push-based, pull-based, or hybrid?
- Do subscribers need independent replay windows?
- How much filtering happens at publish time versus subscribe time?

If left broad, assume durable topics, independent subscriber offsets, optional server-side filtering, and mixed-latency subscribers with different backlog tolerances.

## Requirements

### Functional

- Accept event publishes durably by topic.
- Let many subscriptions consume at independent speeds.
- Support subscription filters and replay within retention windows.
- Track per-subscription lag and delivery state.
- Isolate repeatedly failing or extremely slow subscribers.

### Non-functional

- Scale fanout without copying data naively for each subscription on the write path.
- Prevent a cold or stuck subscriber from forcing unlimited retention growth.
- Make subscription-specific backlog visible to operators.
- Preserve predictable delivery for high-priority subscribers.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Publish rate | 6M events/s peak | write path must stay efficient before fanout |
| Subscriptions | 120K active | per-subscription state becomes a major control-plane concern |
| Average subscribers per topic | 40 | shapes fanout, filtering, and backlog accounting |
| Replay window | 24 hours standard, 7 days premium | storage tiering and product policy are linked |
| Peak factor | 3x during incident broadcasts | surge fanout often hits when subscribers are already stressed |

## Architecture

```text
publishers
  -> topic ingress
  -> partitioned durable log
  -> subscription registry + filter metadata
  -> delivery workers / pull readers
  -> per-subscription cursors and backlog controls
  -> replay / dead-letter controls
```

Design notes:

1. Keep topics durable and subscriber progress separate so subscribers can move independently.
2. Avoid materializing a full copy per subscriber unless the product truly requires it.
3. Give each subscription explicit backlog budgets and pause behavior.
4. Treat filtering as both a compute-cost decision and an isolation decision.

## Data Model & APIs

Core records:

```text
topic
partition_id
subscription_id
filter_expression
committed_cursor
backlog_age
retention_tier
delivery_attempt
```

Useful interfaces:

- `Publish(topic, payload, attributes)`
- `CreateSubscription(topic, filter, delivery_mode)`
- `Pull(subscription_id, max_batch)`
- `Ack(subscription_id, partition, cursor)`
- `SeekSubscription(subscription_id, cursor)`
- `PauseSubscription(subscription_id, reason)`

The system is stronger when topic durability, subscription state, and replay operations are clearly separated.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one subscription falls far behind | backlog age and retained-by-subscriber bytes | pause, quota, or force retention downgrade |
| server-side filters become expensive | delivery CPU by filter and per-topic fanout lag | precomputed indexes for common filters or filter limits |
| replay by one subscriber starves live traffic | replay read bandwidth and live-delivery lag | isolated replay lanes and replay budgets |
| failing subscriber keeps retrying poison events | repeated delivery attempts per cursor | dead-letter path and bounded retries |

## Observability

- metric: publish latency versus delivery latency by subscription tier
- metric: backlog age and retained bytes per subscription
- metric: filter evaluation cost and drop rate
- metric: replay traffic versus live traffic
- log: subscription create, pause, seek, and retention-tier changes
- trace: publish through subscription delivery for sampled events
- SLO: premium subscriptions receive 99% of matching events within the target delivery delay under normal load

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| independent subscription cursors | flexible consumer recovery | more control-plane state | global topic progress |
| server-side filtering | lower subscriber work and egress | platform CPU and complexity | all subscribers filter locally |
| backlog budgets per subscription | bounded shared-platform cost | subscriber experience varies by tier | infinite retention for everyone |

## Interview It

**Google framing:** "Design a pub/sub system for internal service events." Expect follow-ups on subscription state, replay, and why lagging consumers should not hurt healthy ones.

**Cloudflare framing:** "Design multi-tenant fanout for platform notifications or streaming events." Expect pressure on noisy neighbors, filtered delivery cost, and storage blowups from slow tenants.

**Follow-ups:**
1. When would you precompute filtered streams instead of filtering on read?
2. How do you stop one premium tenant from replaying huge windows repeatedly?
3. What changes for push delivery to public webhooks?
4. How would you tier retention by subscription plan?
5. How do you migrate a subscription to more partitions safely?

## Ship It

- `outputs/tradeoff-matrix-pubsub-fanout.md`

## Exercises

1. **Easy** — Explain why per-subscription cursors are different from topic offsets.
2. **Medium** — Add a budget policy for subscribers that stay offline too long.
3. **Hard** — Redesign the fanout path when 5% of subscriptions need strict low latency and the rest can lag.

## Further Reading

- [Google Cloud Pub/Sub subscriber overview](https://cloud.google.com/pubsub/docs/subscriber) — useful for thinking about independent subscription state
- [NATS JetStream concepts](https://docs.nats.io/nats-concepts/jetstream) — helpful background on streams, consumers, and replay
