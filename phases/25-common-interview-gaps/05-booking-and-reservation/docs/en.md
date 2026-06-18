# Hotel/Flight Booking & Seat Reservation

> The seat exists. The payment succeeded. But two customers have the same seat — the reservation race is the hardest correctness problem in this design.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Design a seat reservation system that handles concurrent booking attempts atomically, separates seat locking from payment, tolerates overbooking within configurable policy, manages cancellation and waitlists, and recovers from partial failures in payment-reservation workflows.  
**Prerequisites:** `08-consistency-replication-and-transactions/05-transactions`, `08-consistency-replication-and-transactions/06-sagas`, `19-payments-wallets-and-ordering/03-inventory-reservation`  
**Estimated time:** ~75 min  
**Primary artifact:** capacity sheet

## The Problem

Design a booking system for hotels or flights. Users browse availability, select a seat or room, and complete payment. The system must guarantee no two customers receive the same seat, handle payment failures gracefully, manage cancellations and waitlists, and tolerate the reality that airlines intentionally oversell flights by 5–15%.

The senior-level insight is that the naively "simple" booking flow — check availability, charge card, write reservation — is broken under concurrency. Two users can both see the last available seat, both pass the availability check, and both be charged before the system detects the conflict. The seat assignment must be an atomic operation, decoupled from payment timing, with a time-limited hold that expires if payment does not complete.

The overbooking question separates junior from senior answers. Junior candidates try to prevent all overbooking with strict locks. Senior candidates explain that airlines intentionally oversell because cancellation rates make overselling profitable, and that the system must support a configurable overbooking tolerance per flight.

## Clarify

- Is overbooking intentionally permitted, and if so, at what tolerance level? This changes the seat locking model entirely.
- What is the seat hold timeout — how long does the system reserve a seat for a user who has started checkout but not completed payment?
- Is the waitlist first-come-first-served, or are there priority rules (frequent flyer status, original booking order)?

If the interviewer does not specify, assume overbooking allowed at up to 5% per flight, seat hold timeout of 10 minutes, and FIFO waitlist with manual upgrade to priority queue.

## Requirements

### Functional

- Display seat map and real-time availability for a flight or hotel room type.
- Reserve a seat with a time-limited hold while the user completes payment.
- Charge payment and confirm reservation atomically, or release the hold on payment failure.
- Support cancellation with configurable refund policy (full, partial, no refund by cancellation window).
- Manage a waitlist: notify waitlisted customers when a cancellation releases a seat.
- Support overbooking up to a configurable percentage per flight.

### Non-functional

- Seat hold creation latency: under 500 ms.
- No two confirmed reservations for the same physical seat (unless intentional overbooking is activated by the airline).
- Hold expiry must release the seat within 30 seconds of timeout.
- System availability: 99.99% for booking path; payment processing failures must not corrupt reservation state.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Concurrent bookings | 100K concurrent booking sessions at peak (holiday release) | determines seat lock contention; hot flights need per-seat locking granularity |
| Flights in system | 100K flights in the 180-day booking window | total inventory size; fits in a distributed KV store |
| Seats per flight | 200 seats average | each seat is a separate lockable unit; hot flight = 200 independent lock targets |
| Hold creation rate | 5K holds/s peak | seats are KV rows; 5K writes/s is trivially handled by a sharded DB |
| Cancellation rate | 15% of bookings cancel; 5% of seats go through waitlist | sizes the waitlist queue and notification fan-out |

## Architecture

```text
user -> browse (read replica or cache)
  -> flight search API: reads availability summary from cache
  -> seat map API: reads individual seat status from inventory DB

user -> select seat
  -> hold service
     -> BEGIN TRANSACTION
     -> SELECT seat WHERE status = 'available' FOR UPDATE
     -> UPDATE seat SET status = 'held', held_for = user_id, hold_expires = now + 10min
     -> INSERT hold_record (hold_id, seat_id, user_id, expires_at)
     -> COMMIT
     -> return hold_id to client

user -> complete payment
  -> payment service: charge card (async, 1-5s)
  -> on payment success:
     -> reservation service: confirm hold
        -> UPDATE seat SET status = 'confirmed', reservation_id = ...
        -> UPDATE hold SET status = 'confirmed'
        -> DELETE hold_expires job
  -> on payment failure:
     -> reservation service: release hold
        -> UPDATE seat SET status = 'available'
        -> notify waitlist service

hold expiry (background job, runs every 30s)
  -> find hold_records WHERE expires_at < now AND status = 'held'
  -> release each expired hold
  -> notify waitlist service for released seats

waitlist service
  -> ordered queue per (flight_id, seat_class)
  -> on seat release: notify head of queue via email/push notification
  -> 15-minute window for waitlisted customer to claim seat
```

## Data Model & APIs

Core entities:

```text
Flight  { flight_id, origin, destination, departure_time, aircraft_type,
          total_seats, available_count, held_count, overbooking_limit }
Seat    { seat_id, flight_id, row, column, class, status, held_for, hold_expires_at, reservation_id }
Hold    { hold_id, seat_id, user_id, created_at, expires_at, status }
Reservation { reservation_id, hold_id, user_id, seat_id, flight_id, payment_id,
              confirmed_at, status, refund_policy }
Waitlist { waitlist_id, flight_id, seat_class, user_id, joined_at, status }
```

Key APIs:

- `GET /v1/flights/{flight_id}/seats` — returns seat map with status (available, held, confirmed)
- `POST /v1/holds` — body: `{seat_id, user_id}`; returns `{hold_id, expires_at}`; fails if seat unavailable
- `POST /v1/reservations` — body: `{hold_id, payment_method_id}`; triggers payment and confirms reservation
- `DELETE /v1/reservations/{reservation_id}` — cancel with refund calculation per policy
- `POST /v1/waitlist` — body: `{flight_id, seat_class, user_id}`; joins waitlist

The `POST /v1/holds` endpoint is the critical serialization point: it must use a database-level lock (`SELECT FOR UPDATE`) or optimistic locking with a CAS to prevent double-holds.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Payment service times out after seat hold is created | payment status unknown; hold remains active | implement saga: if payment status unknown after 60s, query payment service; confirm or release based on result |
| Hold expiry job crashes before releasing expired holds | expired holds not released; seats appear unavailable to new bookers | hold expiry job is idempotent; restart releases same holds again; seat remains lockable via DB constraint |
| Two users select same seat simultaneously | one transaction wins lock; other gets lock wait timeout | serialized by DB row lock; loser gets HTTP 409 and must select a different seat |
| Overbooking threshold exceeded due to concurrent confirmations | confirmed_count > total_seats × (1 + overbooking_limit) alert | count confirmed reservations atomically at confirmation time; block new confirmations above threshold |

## Observability

- metric: seat hold creation latency p99 — tracks database lock contention on hot flights
- metric: hold-to-confirmation conversion rate and average time — measures checkout funnel health
- metric: hold expiry release rate — confirms background job is running and releasing seats correctly
- metric: overbooking ratio per flight (confirmed / capacity) — alerts when approaching legal or policy limits
- log: every seat status transition with seat_id, previous status, new status, hold_id/reservation_id, and user_id
- SLO: 99.99% of completed payments result in a confirmed reservation within 5 seconds; no physical seat double-booking

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Pessimistic row-level locking (SELECT FOR UPDATE) for seat hold | simple, correct serialization with no application-level retry | lock contention on high-demand flights (e.g. 1000 users competing for 1 seat causes queue) | optimistic concurrency (CAS on seat version) — reduces contention but causes many application-level retry loops |
| Time-limited hold before payment | prevents payment latency from blocking other buyers; abandoned carts release seats automatically | adds hold management complexity; payment must complete within hold window | reserve seat only after payment — payment takes 3-10s, seat unavailable to other buyers during this time |
| Configurable overbooking percentage per flight | matches airline business model; oversell reduces empty seats | requires voluntary bump management (compensation system); not usable for events where physical capacity is absolute | zero overbooking — simpler but leaves revenue on the table for airlines |

## Interview It

**Google framing:** "Design the seat reservation backend for Google Flights booking integration. How do you handle 100K users attempting to book the same flight at holiday schedule release?" Expect pushback on lock contention, payment-reservation atomicity, and overbooking policy.

**Cloudflare framing:** "How would you cache flight availability data at the edge while keeping the booking path strongly consistent?" Expect questions on stale cache serving incorrect availability and the cache invalidation strategy when seats are held or released.

**Follow-ups:**
1. A payment processor takes 30 seconds to respond due to an outage. How does your hold expiry system interact with this, and can a customer lose their seat during payment?
2. How would you handle a flight cancellation — all 200 reservations must be refunded and all waitlisted users notified?
3. How would you implement seat upgrades where a confirmed customer can upgrade to a higher class if one becomes available?
4. What changes if you need to support group bookings where 10 seats must be reserved atomically?
5. How would you detect and prevent seat scalping bots that hold seats without purchasing?

## Ship It

- `outputs/capacity-sheet-booking-and-reservation.md`

## Exercises

1. **Easy** — Model the hold expiry background job. What query identifies expired holds? How do you ensure it is idempotent so restarting the job does not cause double-releases?
2. **Medium** — Design the saga for payment-reservation atomicity. Draw the state machine for a hold that moves through: created → payment initiated → payment succeeded → confirmed. Include the compensating action for each failure point.
3. **Hard** — A major airline sells 50K seats simultaneously at 09:00 on New Year's Day (schedule release). Design the database sharding and locking strategy to handle 100K concurrent hold attempts without a global queue or single-point bottleneck.

## Further Reading

- https://martinfowler.com/eaaDev/timeoutcommand.html — Martin Fowler's analysis of timeout patterns in reservation systems
- https://shopify.engineering/building-resilient-payment-systems — Shopify's approach to payment-reservation atomicity with saga patterns, directly applicable to booking flows
- https://www.postgresql.org/docs/current/explicit-locking.html — PostgreSQL locking documentation covering SELECT FOR UPDATE and advisory locks, the foundation of database-level seat locking
