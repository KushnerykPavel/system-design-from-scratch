# Order Management & Fulfillment Pipeline

> An order must complete fully or reverse fully — partial fulfillment is a customer support nightmare.

**Type:** Build
**Company focus:** Amazon
**Learning goal:** Design the order management system that handles 1.6M orders/day with inventory reservation, payment, and fulfillment coordination.
**Prerequisites:** `19-payments-wallets-and-ordering/04-order-state-machine`, `08-consistency-replication-and-transactions/06-sagas`
**Estimated time:** ~90 min
**Primary artifact:** order state machine + saga choreography diagram

---

## The Problem

Amazon processes 1.6M+ orders per day (Prime Day: ~8M orders/day = ~93 orders/sec peak). Each order touches at least four independent services — catalog, inventory, payment, and fulfillment — each of which can fail independently. The order must still complete fully or reverse fully; a customer paying but never receiving goods is a business and legal problem.

The core challenge is **distributed coordination without a global transaction**. You cannot use 2PC across microservices at this scale. The solution is a saga with compensating transactions.

---

## Scale Envelope

| Metric | Normal | Prime Day Peak |
|---|---|---|
| Orders/day | 1.6M | ~8M |
| Orders/sec | ~19 | ~93 |
| Fulfillment centers | 175+ globally | — |
| SKUs per order (avg) | ~2.5 | ~2.5 |
| Payment gateway calls/sec | ~19 | ~93 |
| Inventory reservation ops/sec | ~48 | ~233 |

Storage: order records ~2KB each; 1.6M/day × 365 = ~584M orders/year × 2KB = ~1.1TB/year. Retain 7 years for legal/tax = ~8TB. Use DynamoDB with TTL for active orders, archive to S3 Glacier for old orders.

---

## Order State Machine

An order moves through a well-defined set of states. Illegal transitions are rejected.

```
PLACED
  └─► PAYMENT_PENDING
        └─► PAYMENT_CONFIRMED
              └─► INVENTORY_RESERVED
                    └─► PICKING
                          └─► PACKED
                                └─► SHIPPED
                                      └─► OUT_FOR_DELIVERY
                                            └─► DELIVERED
```

**Cancellation paths** (customer or system-initiated):
- PLACED → CANCELLED
- PAYMENT_PENDING → CANCELLED
- PAYMENT_CONFIRMED → CANCELLED (triggers payment refund saga)
- INVENTORY_RESERVED → CANCELLED (triggers inventory release + refund)

**Returns** (post-delivery):
- DELIVERED → RETURN_REQUESTED → RETURNED

**Key invariants:**
- Once SHIPPED, order cannot be CANCELLED (handed to carrier)
- Transitions are idempotent: applying the same transition twice is a no-op, not an error
- Every state change is persisted atomically with a timestamp and actor (system/user)

---

## Saga Pattern for Distributed Coordination

A saga breaks a distributed transaction into a sequence of local transactions, each with a compensating transaction for rollback.

### Happy Path (Choreography-based)

```
Order Service          Payment Service        Inventory Service      Fulfillment Service
     │                       │                       │                       │
     │── OrderPlaced ────────►│                       │                       │
     │                       │── PaymentCaptured ────►│                       │
     │                       │                       │── InventoryReserved ──►│
     │                       │                       │                       │── FulfillmentAssigned
     │◄──────────────────────────────────────────────────────────── FulfillmentAssigned
```

Each service subscribes to events via SQS/SNS, performs its local transaction, and emits the next event.

### Failure & Compensation

| Failure point | Compensating action |
|---|---|
| Payment capture fails | Emit `OrderPaymentFailed` → Order moves to CANCELLED |
| Inventory unavailable after payment | Emit `InventoryReservationFailed` → Payment Service issues refund → Order CANCELLED |
| Fulfillment center offline | Retry with next nearest center; after 3 attempts escalate → manual queue |
| Carrier API timeout during shipping | Retry with exponential backoff; idempotency key prevents duplicate shipment labels |

### Idempotency

Every operation carries an **idempotency key** (composed of `order_id + step_name`). Payment service stores processed keys in Redis with TTL 24h. Re-delivery of the same SQS message is a no-op.

---

## Inventory Reservation

Inventory reservation is the most race-condition-prone step: two customers can simultaneously try to purchase the last unit.

### Optimistic Locking in DynamoDB

```
UpdateItem
  Key: { sku_id: "WIDGET-123", warehouse_id: "SEA-1" }
  UpdateExpression: SET reserved = reserved + :qty, available = available - :qty
  ConditionExpression: available >= :qty
  ExpressionAttributeValues: { :qty: 1 }
```

If `available < qty`, DynamoDB throws `ConditionalCheckFailedException`. The order service catches this, emits `InventoryReservationFailed`, and triggers the refund saga.

**Double-reservation prevention:** The idempotency key `{order_id}-inventory-reserve` is checked before writing. If already present, return success (idempotent re-try).

### Oversell Window

Between the customer seeing "In Stock" on the product page (from the 5s Redis cache) and the reservation, the item may sell out. This ~5s window is acceptable for most SKUs. For high-demand launches (new console, concert tickets), pre-reservation queues limit oversell.

---

## Payment Flow

1. Customer submits order → Order Service creates order in PLACED state
2. Order Service publishes `OrderPlaced` event to SQS
3. Payment Service subscribes, charges the payment method via Stripe/internal gateway
4. On success: Payment Service publishes `PaymentConfirmed` → Order moves to PAYMENT_CONFIRMED
5. On failure: Payment Service publishes `PaymentFailed` → Order moves to CANCELLED

**Async vs synchronous:** Payment is async. The customer sees "Processing payment…" and receives an email confirmation when PAYMENT_CONFIRMED. This decouples order acceptance latency (~50ms to PLACED) from payment gateway latency (~500ms–2s).

**Payment timeout handling:** If the payment gateway does not respond within 30s, the Payment Service emits `PaymentTimedOut`. Order moves to PAYMENT_PENDING with a retry scheduled. After 3 retries over 10 minutes, moves to CANCELLED.

---

## Fulfillment Assignment

Once inventory is reserved, the order must be assigned to a fulfillment center.

**Assignment algorithm:**
1. Find all fulfillment centers that hold the reserved SKUs
2. Score each by: shipping distance to customer address + current workload
3. Assign to lowest-score center
4. Publish `FulfillmentAssigned` event to the assigned center's SQS queue

**Multi-warehouse split shipment:** If no single center has all SKUs, split the order into sub-orders, each fulfilled independently. Each sub-order runs its own state machine. Customer sees a single order with multiple tracking numbers.

**SQS + Workers:** Each fulfillment center has a dedicated SQS queue. Workers (EC2 Auto Scaling group) poll the queue, pick and pack items, and publish `ItemShipped` when the carrier label is generated.

---

## Shipping & Tracking

- Carrier API integration (UPS, FedEx, USPS) called by fulfillment center worker
- Tracking token stored in Order DynamoDB record
- **Status updates:** Carrier uses webhooks to push status events to Amazon's Carrier Integration Service → updates Order state to OUT_FOR_DELIVERY / DELIVERED
- Fallback: if webhook is unreliable, a polling job checks carrier API every 4h for orders older than expected delivery

---

## Observability

**Key metrics (CloudWatch + Prometheus):**
- `order.state_transitions` (counter by from_state, to_state) — detect stuck orders
- `order.saga_compensation_rate` — % of orders requiring compensating transaction; spike = upstream issue
- `payment.capture_latency_p99` — SLA for payment gateway
- `inventory.reservation_conflict_rate` — high rate = demand spike or bug
- `fulfillment.assignment_queue_depth` — per-center backlog
- `order.e2e_latency` (PLACED → SHIPPED p50/p99)

**Alarms:**
- `saga_compensation_rate > 2%` → PagerDuty
- `inventory_conflict_rate > 5%` → investigate inventory data freshness
- `fulfillment_queue_depth > 1000` → scale worker fleet

---

## Failure Modes

### 1. Inventory Reservation Race Condition
- **Risk:** Two customers buy the last unit simultaneously; both succeed due to race
- **Mitigation:** DynamoDB conditional update (optimistic lock); only one thread wins; loser triggers refund saga

### 2. Payment Gateway Timeout
- **Risk:** Gateway does not respond; order stuck in PAYMENT_PENDING; customer charged but order not confirmed
- **Mitigation:** Timeout after 30s → retry 3× with idempotency key → if all fail, CANCELLED + refund; idempotency prevents double charge

### 3. Fulfillment Center Offline
- **Risk:** Assigned fulfillment center goes dark (power outage, network partition)
- **Mitigation:** After 15 min without acknowledgment, re-assign to next nearest center; original reservation released at first center

---

## Trade-offs

| Decision | Chosen | Rejected | Reason |
|---|---|---|---|
| Coordination pattern | Saga (choreography) | 2PC | 2PC requires a coordinator; single point of failure; does not work across microservices at scale |
| Order DB | DynamoDB | PostgreSQL | DynamoDB handles 93 writes/sec trivially; append-only event log maps to document model |
| Inventory locking | Optimistic (conditional update) | Pessimistic (SELECT FOR UPDATE) | Pessimistic locking serializes writes; at 233 ops/sec this creates bottleneck; optimistic fails fast instead |
| Payment | Async via event | Synchronous inline | Async decouples order acceptance latency from gateway latency; improves p99 for customers |
| Fulfillment routing | Score-based nearest-center | Global round-robin | Round-robin ignores shipping cost and center backlog; score-based minimizes delivery time and cost |

---

## Follow-up Questions

1. **Split shipment from multiple warehouses:** The order has 3 items, each in a different fulfillment center. How does the state machine handle this? (Answer: create child order records; parent order reaches DELIVERED only when all children are DELIVERED)

2. **Prime Now 2-hour delivery:** Requires inventory at a local urban delivery station (not a warehouse). How does the assignment algorithm change for 2-hour SLA? (Answer: geo-radius filter first, strict time-window constraint, dedicated last-mile carrier integration)

3. **Returns pipeline:** Customer requests a return. How does the system generate a return label, track the inbound package, and issue a refund? (Answer: reverse saga — issue return label → track via carrier webhook → on scan at fulfillment center, trigger refund payment saga)

4. **Fraud detection integration:** Where in the order flow do you add a fraud check, and what happens if the fraud model flags the order? (Answer: after PLACED, before payment capture; if flagged → manual review queue → order stays in PAYMENT_PENDING; timeout triggers CANCELLED)

5. **Exactly-once delivery guarantees:** SQS delivers at-least-once. How do you ensure the fulfillment worker does not pick and pack the same order twice? (Answer: idempotency key per step stored in DynamoDB; worker checks before acting; SQS message deduplication ID as secondary guard)
