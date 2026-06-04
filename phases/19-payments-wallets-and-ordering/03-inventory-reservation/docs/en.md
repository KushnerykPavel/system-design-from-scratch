# Inventory Reservation System

> Reservation systems are payments-adjacent because overselling inventory and overspending money are the same shape of bug.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design an inventory reservation service that holds limited stock safely across cart, checkout, payment, and cancellation flows.
**Prerequisites:** `09-partitioning-sharding-and-rebalancing/02-hot-partitions`, `19-payments-wallets-and-ordering/02-digital-wallet`, `10-reliability-retries-and-backpressure/05-async-backpressure`
**Estimated time:** ~75 min
**Primary artifact:** reservation-plan validator + oversell review card

## The Problem

Design an inventory reservation system for limited stock such as event tickets, flash-sale items, or regional fulfillment units. Orders may place temporary reservations, then either confirm them after payment or release them on failure or timeout.

This lesson matters because senior interviewers care about the coordination boundary: do you reserve inventory before payment, after payment, or with compensating workflows? A credible answer names the oversell risk, hot-key pressure, and cleanup path.

## Clarify

- Are items fungible within a SKU, or do they have unique identities like seats?
- Do we need temporary cart reservations, or only reservation at checkout?
- How long can reservations live before expiry?
- Is oversell tolerated briefly with later compensation, or must it be prevented strictly?

If the interviewer stays broad, assume per-SKU inventory pools, short-lived checkout holds, strict no-oversell on the authoritative path, and asynchronous release after failed payment or abandoned checkout.

## Requirements

### Functional

- Reserve inventory temporarily for a checkout attempt.
- Confirm a reservation into a committed allocation on successful payment.
- Release reservations on timeout, cancellation, or payment failure.
- Support idempotent retries from order and payment services.
- Expose inventory availability for product and operator systems.

### Non-functional

- Prevent oversell on the source-of-truth path.
- Survive hot SKUs during flash traffic.
- Make reservation leaks visible quickly.
- Keep release and confirm workflows safe under retries.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Reservation attempts | 180K/s peak during launches | drives per-SKU contention and hot-partition handling |
| Unique hot SKUs | as low as 1 for major drops | worst-case design cannot assume smooth spread |
| Reservation TTL | 2 to 15 minutes | shapes cleanup and customer fairness |
| Availability reads | 5M/s peak | separates serving path from authoritative counters |
| Peak factor | 20x over baseline | flash-sale behavior dominates design more than steady state |

## Architecture

```text
cart / checkout
  -> reservation API
  -> authoritative stock service
  -> reservation state store
  -> confirm / release workflow
  -> availability projection cache
```

Design notes:

1. Use a clear authoritative inventory counter or per-item allocator; do not let cached reads authorize reservations.
2. Model reservations as leases with TTL and downstream order reference.
3. Keep confirmation idempotent because payment callbacks and order retries will repeat.
4. Expect hot-key mitigation such as per-SKU queues, shard-local sequencing, or virtual buckets.

## Data Model & APIs

Core entities:

```text
sku_inventory(sku_id, region, total_units, reserved_units, committed_units, version)
reservation(reservation_id, sku_id, quantity, order_id, expires_at, state)
reservation_event(event_id, reservation_id, type, actor, created_at)
availability_projection(sku_id, sellable_units, stale_at)
```

Useful interfaces:

- `CreateReservation(sku_id, quantity, order_id, ttl, idempotency_key)`
- `ConfirmReservation(reservation_id, idempotency_key)`
- `ReleaseReservation(reservation_id, reason_code, idempotency_key)`
- `GetAvailability(sku_id)`
- `ListLeakedReservations(age_bucket)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| oversell under concurrent reserve requests | negative sellable invariant and version-conflict spikes | authoritative compare-and-swap or serialized per-SKU writes |
| leaked reservations reduce stock forever | aged reservation metrics and sweeper backlog | TTL-based expiry plus operator release tooling |
| hot SKU melts one partition | per-SKU latency and conflict skew | dedicated lanes, virtual shards, or admission control |
| confirm arrives after expiry and release | late-confirm count and duplicate workflow events | explicit state machine with idempotent final transitions |

## Observability

- metric: reservation success rate, conflict rate, and confirm latency
- metric: active reservation age buckets and expired-but-not-released count
- metric: availability skew between source of truth and serving cache
- log: forced inventory adjustments, manual overrides, and leak cleanups
- trace: checkout -> reservation -> payment -> confirm or release
- SLO: prevent oversell on authoritative stock while keeping checkout reservation latency within target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| strict reservation before payment | strongest oversell protection | more abandoned reservation churn | charge first and compensate later |
| TTL-based reservations | bounded leak window | possible customer frustration on expiry | indefinite cart reservations |
| per-SKU serialization for hot items | simplest correctness model | hot-spot throughput ceiling | fully optimistic globally shared writes |

## Interview It

**Google framing:** "Design inventory reservation for a high-scale commerce checkout." Expect pushback on hot SKUs, leaked holds, and how payment completion interacts with reservation expiry.

**Cloudflare framing:** "Design limited-capacity allocation for global resources such as reserved compute or bandwidth slots." Expect pressure on regional state ownership and avoiding global coordination for every read.

**Follow-ups:**
1. What changes if seats are unique rather than fungible?
2. How do you handle split inventory across regions?
3. What if payment takes longer than the reservation TTL?
4. What if product wants cart holds for anonymous users?
5. How do you migrate from eventual oversell compensation to strict reservation?

## Ship It

- `outputs/interview-card-inventory-reservation.md`

## Exercises

1. **Easy** — Pick the source-of-truth write path for reserve, confirm, and release.
2. **Medium** — Explain how you would handle one SKU generating half the platform traffic.
3. **Hard** — Redesign for reserved event seats where each seat has identity and adjacency constraints.

## Further Reading

- [Ticketmaster queueing and flash-sale style constraints](https://queue-it.com/blog/ticketing-system-design/) — useful context for hot inventory behavior
- [Google SRE overload chapters](https://sre.google/workbook/index/) — good framing for admission control under bursts
