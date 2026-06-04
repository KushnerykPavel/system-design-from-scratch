# Market Data Fanout and Subscriber Tiers

> The exchange core can be elegant and still fail the interview if your market-data layer lets one premium firehose user destabilize everyone else.

**Type:** Build
**Company focus:** Cloudflare
**Learning goal:** Design a market-data distribution system with tiered subscribers, bounded replay, and low-latency fanout that does not leak downstream pain back into the trading core.
**Prerequisites:** `18-messaging-and-job-platforms/04-pubsub-fanout`, `20-low-latency-location-and-market-systems/04-stock-exchange`, `11-observability-slos-and-debugging/05-alert-design`
**Estimated time:** ~60 min
**Primary artifact:** market-data policy validator + subscriber-tier checklist

## The Problem

Design the market-data layer that publishes quotes, trades, depth changes, and derived feed products to many subscribers. Some subscribers want low-latency raw feeds. Others accept delayed or sampled data. The platform must protect itself from slow consumers and replay abuse.

This lesson matters because senior answers separate matching from fanout, but they also separate subscriber tiers from one another. The system must define what is real-time, what is delayed, and what happens when subscribers lag.

## Clarify

- Are subscribers internal services, brokers, retail clients, or external data vendors?
- Which products are premium low-latency versus delayed public feeds?
- Do subscribers need replay windows or only live delivery?
- How much server-side filtering or aggregation is expected?

If broad, assume multiple feed tiers: premium near-real-time subscribers, standard delayed subscribers, and downstream analytics consumers with replay capability.

## Requirements

### Functional

- Fan out market-data events to many subscribers.
- Support tiered latency and retention policies.
- Allow replay within bounded windows for eligible tiers.
- Track per-subscriber lag and enforce controls on slow consumers.
- Support filtering by symbol set or feed product.

### Non-functional

- Protect the trading core from downstream fanout failures.
- Bound shared infrastructure cost from replay and retention.
- Keep premium latency predictable under bursty symbol traffic.
- Avoid giving one subscriber unlimited backlog power.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Live feed events | 20M updates/s peak | raw fanout volume dominates network design |
| Subscribers | 60K active | subscriber state and lag accounting matter |
| Premium share | 5% of subscribers, 40% of traffic value | tiering is a product and SLO decision |
| Replay window | 15 minutes premium, 1 minute standard | retention cost and storage shape differ by tier |
| Symbol filters | average 200 symbols per subscriber | filtering cost can dominate the delivery plane |

## Architecture

```text
matching outputs
  -> feed normalizer
  -> partitioned market-data log
  -> tier-aware delivery layer
  -> edge fanout / regional relays
  -> subscriber sessions and replay controls
```

Design notes:

1. Keep market-data fanout asynchronous from the matching engine.
2. Give each tier explicit backlog, delay, and replay limits.
3. Decide whether filtering happens at publish time, relay time, or subscriber pull time based on cost and latency.
4. Treat relays and edge nodes as distribution helpers, not new sources of truth.

## Data Model & APIs

Core records:

```text
feed_event(seq, symbol, product, payload, created_at)
subscriber(subscriber_id, tier, symbols, mode, replay_window)
delivery_state(subscriber_id, partition, cursor, lag_ms, retained_bytes)
```

Useful interfaces:

- `Subscribe(feed, tier, symbols)`
- `AckCursor(subscriber_id, partition, seq)`
- `Replay(subscriber_id, from_seq)`
- `PauseSubscriber(subscriber_id, reason)`

Strong answers say exactly how premium and standard subscribers differ operationally, not only commercially.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one subscriber lags and retains too much history | lag and retained-by-subscriber bytes | tier quotas, pause, or replay downgrade |
| symbol-filter computation becomes the bottleneck | filter CPU and per-product delivery lag | precomputed bundles for common symbol sets |
| replay storms collide with live feed delivery | replay bandwidth versus live-latency metrics | separate replay lanes and replay budgets |
| regional relay outage creates concentrated reconnect bursts | reconnect rate and relay session churn | staged reconnect, jitter, and multi-relay failover |

## Observability

- metric: premium versus standard delivery latency
- metric: subscriber lag, retained bytes, and replay traffic
- metric: filter CPU and per-symbol fanout skew
- metric: relay health, reconnect storms, and session churn
- log: subscriber tier changes, pause actions, and replay-limit violations
- trace: normalized feed event through relay and subscriber delivery
- SLO: premium subscribers receive matching feed products within the target delay without trading-core coupling

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| tier-specific backlog budgets | protects shared platform | uneven subscriber experience | one unlimited policy for everyone |
| separate replay lanes | protects live feed latency | more infrastructure and policy logic | replay on the same hot path as live traffic |
| regional relays | lowers client latency and origin load | extra failover paths | every subscriber connects to one central source |

## Interview It

**Google framing:** "Design market-data fanout or a similar high-rate subscription system." Expect follow-ups on replay, subscriber isolation, and product tiering.

**Cloudflare framing:** "Design a global distribution layer for latency-sensitive event feeds." Expect pressure on relay topology, reconnect storms, and multi-tenant fairness.

**Follow-ups:**
1. What changes if a premium subscriber wants thousands of symbols instead of hundreds?
2. How do you enforce replay limits without surprising paying users?
3. What if retail delayed feeds can be sampled but institutional feeds cannot?
4. How would you migrate subscribers between relays safely?
5. What if a region disconnects and then reconnects tens of thousands of sessions at once?

## Ship It

- `outputs/subscriber-tier-checklist-market-data.md`

## Exercises

1. **Easy** — Explain why replay should usually be isolated from the live fanout hot path.
2. **Medium** — Compare filter-on-read versus precomputed product bundles.
3. **Hard** — Redesign the fanout layer when one exchange event causes a 50x reconnect storm and simultaneous replay demand.

## Further Reading

- [NATS JetStream concepts](https://docs.nats.io/nats-concepts/jetstream) — helpful background on streams, consumers, and replay
- [Cloudflare blog](https://blog.cloudflare.com/) — useful operational framing for global distribution and edge relays
