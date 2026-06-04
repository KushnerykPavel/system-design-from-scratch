# Stock Exchange Matching Engine

> Matching engines are where "eventually consistent later" stops sounding clever and starts sounding like regulatory trouble.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a low-latency matching engine with explicit order-book correctness, deterministic sequencing, risk boundaries, and market-data side effects.
**Prerequisites:** `08-consistency-replication-and-transactions/05-transactions`, `18-messaging-and-job-platforms/01-distributed-message-queue`, `19-payments-wallets-and-ordering/01-payment-ledger`
**Estimated time:** ~90 min
**Primary artifact:** matching-engine validator + order-path checklist

## The Problem

Design the core of a stock exchange or exchange-like marketplace that accepts orders, matches buys and sells, and publishes fills and market data. Traders care about fairness and latency. Operators care about deterministic recovery. Regulators care about reconstructability.

This lesson matters because low-latency interviews tempt candidates into magical "single thread solves everything" answers. Strong answers explain symbol partitioning, pre-trade risk, commit boundaries, recovery from journal replay, and what must happen before an ack is safe.

## Clarify

- Is this equities, crypto, internal ad bidding, or a simplified exchange-like system?
- Do we require strict price-time priority?
- Is the scope one symbol, one venue, or many symbols?
- What latency target matters for order acceptance and execution acknowledgment?

If the interviewer stays broad, assume an exchange with multiple symbols, strict price-time priority per symbol, sub-millisecond matching inside a symbol partition, and durable journaling before final acknowledgment.

## Requirements

### Functional

- Accept new, cancel, and modify order requests.
- Enforce pre-trade validation and basic risk checks.
- Match orders using deterministic price-time priority.
- Persist the authoritative order/event history for replay.
- Publish executions and market-data updates after matching.

### Non-functional

- Preserve deterministic replay for post-incident reconstruction.
- Keep the per-symbol hot path extremely short.
- Avoid cross-symbol coupling on the execution path.
- Ensure acknowledgments reflect a durable correctness boundary.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Orders | 5M orders/s peak venue-wide | partitioning by symbol and gateway fan-in matter |
| Hot symbols | top 1% produce 40% of traffic | hot-symbol isolation dominates the design |
| Cancel ratio | 60% of messages | order-book mutation path is not just new orders |
| Market-data fanout | 15M subscriber updates/s derived | execution path must not block on fanout |
| Recovery target | replay one symbol partition in minutes | journal format and snapshot cadence matter |

## Architecture

```text
clients
  -> session gateway
  -> pre-trade validation / risk check
  -> symbol router
  -> single-writer matching engine per symbol partition
  -> durable journal + snapshots
  -> fill publisher + market-data publisher
```

Design notes:

1. Keep one authoritative sequencing point per symbol partition so price-time priority stays deterministic.
2. Make pre-trade validation fast and bounded, but do not let market-data fanout sit on the same critical path as matching.
3. Use append-only journaling for reconstructability and warm restart.
4. Be explicit about what is acknowledged after validation versus after durable sequencing.

## Data Model & APIs

Core records:

```text
order(order_id, account_id, symbol, side, price, quantity, tif, state)
match_event(seq, symbol, resting_order_id, taking_order_id, price, quantity)
journal_record(seq, symbol, event_type, payload)
order_book_snapshot(symbol, seq, bids, asks)
```

Useful interfaces:

- `NewOrderSingle`
- `CancelOrder`
- `ModifyOrder`
- `GetTopOfBook(symbol)`
- `ReplayFrom(seq)`

The most credible answers distinguish the deterministic matching core from surrounding risk, journaling, and market-data systems.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one hot symbol overwhelms its partition | per-symbol queue depth and match latency | symbol-level isolation, admission controls, and hot-symbol hardware tiering |
| journal append stalls the ack path | ack latency and journal fsync histograms | fast local append log, batching, and safe degrade before accepting new load |
| market-data publisher lags behind trades | seq-gap metrics between execution and feed | isolate fanout from matching and allow downstream catch-up |
| replay produces a different book than live | replay checksum mismatch and snapshot divergence | deterministic event order, pure matching logic, and verified snapshots |

## Observability

- metric: order accept latency, match latency, and cancel latency by symbol
- metric: journal append latency and snapshot age
- metric: per-symbol queue depth and dropped/admission-controlled requests
- metric: execution-to-market-data sequence lag
- log: all manual halts, symbol routing changes, and recovery events
- trace: bounded request path through gateway, risk, router, and matching partition
- SLO: accepted orders for healthy symbols are durably sequenced and acknowledged within target latency while replay remains deterministic

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| single-writer per symbol partition | deterministic sequencing and simpler fairness story | limited vertical scale per symbol | multi-writer concurrent matching on one book |
| durable journal before final ack | strong recovery boundary | added hot-path latency | acknowledge before authoritative persistence |
| separate market-data fanout path | protects execution latency | more downstream lag handling | synchronous fanout in the match path |

## Interview It

**Google framing:** "Design a matching engine or auction system with strict fairness and low latency." Expect follow-ups on deterministic ordering, partitioning, and recovery.

**Cloudflare framing:** "Design a globally distributed ingress and fanout layer around a latency-critical central core." Expect pressure on isolation, backpressure, and operational degradation.

**Follow-ups:**
1. What changes if one symbol becomes 100x hotter than the rest?
2. How do you recover one partition after a crash without reopening with a corrupted book?
3. What if pre-trade risk checks need shared account state?
4. How would you support read replicas for top-of-book without slowing the core?
5. Which guarantees would you keep if latency budgets got even tighter?

## Ship It

- `outputs/order-path-review-stock-exchange.md`

## Exercises

1. **Easy** — Explain why single-writer per symbol is attractive even if venue-wide throughput is massive.
2. **Medium** — Compare acknowledging after validation versus after durable sequencing.
3. **Hard** — Redesign for a venue where one symbol becomes permanently hotter than a single machine can handle.

## Further Reading

- [LMAX Disruptor](https://lmax-exchange.github.io/disruptor/) — useful background on latency-sensitive single-writer event processing
- [The Tail at Scale](https://research.google/pubs/the-tail-at-scale/) — good framing for why fanout cannot sit inside the critical core path
