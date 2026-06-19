# Trip State Machine & Recovery
> A trip must always be in exactly one state — ambiguity means incorrect billing and angry customers.

**Type:** Build
**Company focus:** Uber
**Learning goal:** Design the trip state machine that tracks a ride from request to completion, handles driver and rider app crashes, and ensures correct billing even under partial failures.
**Prerequisites:** `08-consistency-replication-and-transactions/06-sagas`, `10-reliability-retries-and-backpressure/02-idempotency-under-failure`
**Estimated time:** ~75 min
**Primary artifact:** trip state machine diagram + recovery runbook

---

## Why State Machines for Trips?

A trip passes through many phases: request, matching, pickup, ride, completion. Without an explicit state machine, you get:
- Double billing (charge fired twice because two workers both think the trip just completed)
- Ghost trips (driver thinks COMPLETED, server thinks IN_PROGRESS)
- Incorrect cancellation fees (cancelled after the free window but the system didn't record the window start time)

The discipline of "one authoritative state, explicit transitions only, idempotent events" eliminates these bugs by design.

---

## Trip States

```
                              ┌────────────────────────────┐
                              │         REQUESTED          │
                              │  rider submits trip request │
                              └────────────┬───────────────┘
                                           │ dispatch finds candidates
                                           ▼
                              ┌────────────────────────────┐
                              │          MATCHING          │
                              │  DISCO evaluating matches  │
                              └────────────┬───────────────┘
                                           │ driver accepts offer
                                           ▼
                              ┌────────────────────────────┐
                              │      DRIVER_ASSIGNED       │
                              │  driver committed to trip  │◄──── free cancellation window starts
                              └────────────┬───────────────┘
                                           │ driver starts navigating
                                           ▼
                              ┌────────────────────────────┐
                              │      DRIVER_EN_ROUTE       │
                              │  driver heading to pickup  │
                              └────────────┬───────────────┘
                                           │ driver arrives
                                           ▼
                              ┌────────────────────────────┐
                              │       DRIVER_ARRIVED       │
                              │  driver at pickup location │
                              └────────────┬───────────────┘
                                           │ rider enters vehicle
                                           ▼
                              ┌────────────────────────────┐
                              │        IN_PROGRESS         │◄──── heartbeat required (GPS + app)
                              │    trip is underway        │
                              └────────────┬───────────────┘
                                           │ driver ends trip
                                           ▼
                              ┌────────────────────────────┐
                              │         COMPLETED          │◄──── billing triggered
                              └────────────────────────────┘
```

### Terminal States

| State | Triggered By |
|---|---|
| `COMPLETED` | Driver taps "End Trip" |
| `CANCELLED_BY_RIDER` | Rider cancels |
| `CANCELLED_BY_DRIVER` | Driver cancels |
| `CANCELLED_BY_SYSTEM` | Dead trip detection, no available drivers |

Once a trip reaches a terminal state, no further transitions are permitted.

---

## Valid Transitions

| From | To | Trigger |
|---|---|---|
| REQUESTED | MATCHING | Dispatch begins search |
| MATCHING | DRIVER_ASSIGNED | Driver accepts offer |
| MATCHING | CANCELLED_BY_SYSTEM | No driver found in window |
| DRIVER_ASSIGNED | DRIVER_EN_ROUTE | Driver taps "Navigate" |
| DRIVER_ASSIGNED | CANCELLED_BY_RIDER | Rider cancels (within free window) |
| DRIVER_ASSIGNED | CANCELLED_BY_DRIVER | Driver cancels |
| DRIVER_EN_ROUTE | DRIVER_ARRIVED | Driver arrives at pickup |
| DRIVER_EN_ROUTE | CANCELLED_BY_RIDER | Rider cancels (fee applies) |
| DRIVER_EN_ROUTE | CANCELLED_BY_DRIVER | Driver cancels |
| DRIVER_ARRIVED | IN_PROGRESS | Rider enters vehicle |
| DRIVER_ARRIVED | CANCELLED_BY_RIDER | Rider no-show |
| IN_PROGRESS | COMPLETED | Driver ends trip |
| IN_PROGRESS | CANCELLED_BY_SYSTEM | Dead trip (no heartbeat 10 min) |

Any transition not in this table is **rejected with an error**.

---

## Idempotency Rule

State transitions must be idempotent:

```
Transition(trip, COMPLETED) → success
Transition(trip, COMPLETED) → success (no-op, not an error)
Transition(trip, DRIVER_ASSIGNED) → error: invalid transition COMPLETED → DRIVER_ASSIGNED
```

Why? The driver app may retry a "trip ended" event due to network failure. The billing worker must only fire once. The second `Transition(trip, COMPLETED)` call returns success without triggering billing again.

**Implementation**: Store the trip document in Docstore (Uber's document store) with optimistic locking:
1. Read trip document + version number
2. Validate transition
3. Write new state only if version number unchanged
4. If version conflict → re-read and retry (someone else already transitioned, check if result is idempotent)

---

## Cancellation Windows & Fees

| Time Since DRIVER_ASSIGNED | Policy |
|---|---|
| 0–2 minutes | Free cancellation |
| > 2 minutes | Cancellation fee applied |

The 2-minute window start is stored on the trip document at assignment time. The cancellation handler computes `now - assignedAt` to determine the fee.

**Edge case**: clock skew between device and server. Always use the server-side timestamp recorded at `DRIVER_ASSIGNED` transition, not the device timestamp.

---

## Billing Trigger

Billing is triggered by the `COMPLETED` → billing worker event flow:

```
Trip transitions to COMPLETED
        │
        ▼
Kafka: trip-completed event (tripID, driverID, riderID, route_trace)
        │
        ▼
Billing worker:
  1. Compute final fare (base fare + distance + time + surge multiplier)
  2. Charge pre-authorized payment method (idempotency key = tripID)
  3. Record earnings to driver payout ledger
  4. Emit receipt to rider
```

The idempotency key on the charge is the `tripID`. If the billing worker crashes after charging but before recording, a retry re-sends the same key to the payment processor, which returns the already-processed result rather than double-charging.

---

## Recovery Scenarios

### Scenario 1: Driver App Crash During Trip

```
Trip state: IN_PROGRESS (on server)
Driver GPS heartbeat stops at T=0.

T+30s: Redis driver location key expires → driver marked UNREACHABLE
T+10m: Dead trip detection fires:
  - Last known GPS: stationary for >10 min
  - Trip in IN_PROGRESS with no heartbeat
  - System sends driver notification: "Are you still on this trip?"
  - If no response within 2 min: transition to CANCELLED_BY_SYSTEM

Recovery if driver reconnects before 10 min:
  - Driver app re-fetches trip state on startup → gets IN_PROGRESS
  - App resumes GPS streaming
  - No state change needed (idempotent reconnect)
```

### Scenario 2: Rider App Crash During Trip

```
Trip state: IN_PROGRESS (on server)
Rider app crashes at T=0.

Trip continues normally — driver is unaffected.
Rider reconnects:
  - App fetches /trips/current → receives IN_PROGRESS trip with driver ETA
  - App resumes real-time tracking via WebSocket
```

### Scenario 3: Server Crash Mid-Trip

```
Trip state: IN_PROGRESS stored in Docstore (persistent, not in-memory)
Server process crashes.

Cadence/Temporal workflow is the authoritative state keeper:
  - Each transition is recorded as a workflow activity
  - On crash, Temporal re-runs from the last completed activity checkpoint
  - Trip resumes from exact last state without re-running billing or duplicate state transitions
```

### Scenario 4: State Split (Driver ≠ Server)

```
Driver app shows: COMPLETED
Server shows: IN_PROGRESS

Cause: driver sent COMPLETED event, server acknowledged, but
       server's state write was lost before persistence (edge case).

Resolution:
  1. Post-trip reconciliation job runs every 5 minutes
  2. Detects trips in IN_PROGRESS with no heartbeat for >15 min
  3. Fetches GPS trace from Cassandra for the trip window
  4. If GPS trace shows trip ended (driver location returned to base),
     transition to COMPLETED using the GPS end time
  5. Triggers billing with reconciled end time
```

---

## Dead Trip Detection

A trip is considered **dead** when both conditions hold:
1. Trip is in `IN_PROGRESS` state
2. No driver GPS heartbeat received for > 10 minutes AND last GPS location is stationary

Dead trip action:
1. Notify driver via push notification
2. Wait 2 minutes for driver response
3. If no response: transition to `CANCELLED_BY_SYSTEM`, trigger partial billing (time-based up to last known GPS point), trigger rider refund for uncompleted portion

---

## Observability

| Metric | SLO | Alert |
|---|---|---|
| Trip state transition latency P99 | <200ms | >1s |
| Billing trigger latency (COMPLETED → charged) | <30s | >120s |
| Dead trip rate | <0.01% | >0.1% |
| State split detection rate | <0.001% | >0.01% |
| Idempotent retry rate | <1% of transitions | >5% (retry storm indicator) |

---

## Failure Modes Summary

| Failure | Impact | Mitigation |
|---|---|---|
| State machine splits | Ghost trip in one system | Post-trip reconciliation + GPS trace |
| Billing double-charge | Rider overcharged | Idempotency key on payment processor call |
| Cancellation fee on cancelled trip | Rider disputes | Server-side cancellation window clock |
| Dead trip no-show billing | Rider charged for incomplete trip | Partial billing capped at last GPS point |
| Docstore write conflict | Transition lost | Optimistic locking retry with idempotency check |
