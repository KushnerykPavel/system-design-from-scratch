# HTTP vs gRPC vs Event Interfaces

> Choose the interface that matches the coordination pattern, not the one that sounds modern.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Decide when request-response APIs, strongly typed RPC, or asynchronous events are the right boundary for a system and explain the operational consequences.  
**Prerequisites:** `03-design-framework-and-timing/03-diagram-then-dive`, `03-design-framework-and-timing/04-deep-dive-selection`  
**Estimated time:** ~75 min  
**Primary artifact:** trade-off matrix + interview card  

## The Problem

Senior interview answers get weaker when every boundary becomes "an API" without naming whether the caller needs immediate acknowledgement, strong contracts, fanout, replay, or loose coupling.

This lesson gives you a decision frame:

- `HTTP` when humans, browsers, or broad ecosystem compatibility dominate
- `gRPC` when low-latency service-to-service contracts and typed clients matter
- `events` when decoupling, fanout, and asynchronous workflows dominate

## Clarify

- Does the caller need a synchronous answer before it can continue?
- Is the main risk latency, team coordination, or replayable state change?
- Will many consumers need the same signal, or is this mostly one caller to one callee?
- Are external clients involved, or is this an internal boundary you control tightly?

## Requirements

### Functional

- Select an interface style for a given interaction.
- Explain what happens to retries, error handling, and schema management.
- Distinguish command, query, and event boundaries.

### Non-functional

- Keep latency and operational complexity visible.
- Avoid accidental coupling across teams.
- Make failure behavior clear under overload and partial delivery.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Front-door QPS | 150K HTTP req/s | browsers and external SDKs bias toward HTTP semantics |
| Internal RPC fanout | 4-10 downstream calls per request | typed service contracts and latency overhead matter |
| Event fanout | 5-30 consumers per domain event | favors asynchronous delivery and replay |
| Peak factor | 4x during launches | tests backpressure and queue depth choices |
| Rough cost | gateway + service mesh + event bus | multiple interface types are justified only if they reduce bigger risks |

## Architecture

Think in interaction shapes:

1. **Query or command with immediate user impact**  
   Prefer HTTP or gRPC.
2. **Low-latency internal service hop**  
   Prefer gRPC if you control both sides.
3. **State change that many systems react to later**  
   Prefer events.

Typical mix:

```text
client -> HTTP API gateway
gateway -> gRPC user/profile/order services
services -> event bus for "order_created", "user_updated", "quota_exhausted"
```

The strong answer is usually not "pick one everywhere." It is "use synchronous boundaries on the critical path and events for downstream side effects."

## Data Model & APIs

Example boundaries:

- `POST /orders` over HTTP for public client compatibility
- `InventoryService/ReserveStock` over gRPC for internal low-latency reservation
- `order.created.v2` event for email, analytics, and fulfillment consumers

Questions to answer explicitly:

- who owns the contract
- which side can retry
- how errors are surfaced
- whether consumers can replay or reprocess

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| event chosen for a synchronous user path | user sees delayed or ambiguous confirmation | keep user-facing commit path synchronous |
| gRPC used for many external clients | SDK churn and proxy/tooling friction | keep public surface HTTP unless strong reason otherwise |
| HTTP used for high-fanout domain propagation | webhook storms or N-way coupling | emit a durable event instead |
| interface choice hides ownership | many teams block on one contract change | separate public, internal, and event contracts explicitly |

## Observability

- metric: per-interface latency and error rate
- metric: downstream fanout count for synchronous requests
- metric: event lag and redelivery count
- log: contract version and caller identity on failures
- trace: request path plus async continuation IDs
- SLO: synchronous critical path meets target latency without offloading correctness to eventual consumers

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| HTTP for public APIs | wide compatibility and debuggability | more verbose contracts and weaker typing | exposing raw gRPC to third parties |
| gRPC for internal hot paths | efficient, typed, low-latency calls | stronger deployment coordination | JSON/HTTP between tightly coupled services |
| events for side effects | decoupling and replay | eventual consistency and harder debugging | chaining synchronous calls for every downstream reaction |

## Interview It

**Google framing:** "Design the interfaces for a checkout system." The signal is whether you separate the synchronous payment decision from asynchronous downstream work.

**Cloudflare framing:** "Design the control and data plane interfaces for an edge product." The signal is whether you distinguish fast request-path RPC from eventual configuration propagation.

**Follow-ups:**
1. What changes if mobile clients and third-party integrators use the same API?
2. What if one downstream consumer needs strict ordering?
3. What if the synchronous path is missing an acknowledgement budget?
4. What if you need replay for compliance or backfill?
5. When is it worth supporting more than one interface style in the same design?

## Ship It

- `outputs/tradeoff-matrix-http-vs-grpc-vs-events.md`
- `outputs/interview-card-http-vs-grpc-vs-events.md`

## Exercises

1. **Easy** — Classify five interactions in a file-upload product as HTTP, gRPC, or event driven.  
2. **Medium** — Redesign a chat system that currently uses synchronous HTTP for every notification side effect.  
3. **Hard** — Explain the control-plane vs data-plane interface mix for a global API gateway.  

## Further Reading

- [gRPC documentation](https://grpc.io/docs/) — useful for understanding typed RPC trade-offs  
- [CloudEvents](https://cloudevents.io/) — helpful framing for event contracts  
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview structure for selecting interfaces  
