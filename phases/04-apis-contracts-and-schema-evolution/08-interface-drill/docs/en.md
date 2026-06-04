# Interface Design Drill

> A mature interface answer sounds like a series of deliberate constraints, not a random set of endpoint names.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Practice turning an ambiguous product prompt into an interface plan that covers transport choice, API shape, compatibility, retries, and safe defaults under interview time pressure.  
**Prerequisites:** `04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events`, `04-apis-contracts-and-schema-evolution/07-api-safety-defaults`  
**Estimated time:** ~60 min  
**Primary artifact:** drill worksheet + scoring rubric  

## The Problem

By the end of this phase, the learner should be able to answer questions like:

- what interface type fits this boundary?
- how are writes retried safely?
- how does listing stay efficient?
- how will the contract evolve without breaking clients?

This drill compresses those moves into one timed exercise.

## Clarify

- Which interaction is on the user-critical path?
- Which clients are public versus internal?
- What query patterns or writes are likely to hurt the system?
- What change pressure will the contract face over the next year?

## Requirements

### Functional

- Pick interface styles for key boundaries.
- Design one write contract with retry safety.
- Design one list contract with bounded query shape.

### Non-functional

- Keep the answer time-boxed and prioritized.
- Make operational risks visible, not just payload shapes.
- Leave room for versioning and rollout discussion.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Core API QPS | 50K req/s | enough scale to make contract shape matter |
| Internal fanout | 3-8 downstream actions | interface choice changes latency and coupling |
| Retry burst | 5x during incidents | forces idempotency discussion |
| Query page size | 20-100 items | shapes list safety defaults |
| Rough cost | gateway, validation, versioning, and observability overhead | interfaces are part of the cost model too |

## Architecture

Recommended drill sequence:

1. Clarify the user path and team boundaries.
2. Choose HTTP, gRPC, or event interfaces for the main interactions.
3. Specify one write API with idempotency.
4. Specify one list API with pagination and filter constraints.
5. Close with compatibility, failure modes, and observability.

## Data Model & APIs

A strong drill answer usually includes:

- one public request/response example
- one internal or async contract example
- idempotency key handling
- page size cap and cursor shape
- a brief versioning or schema-evolution policy

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| interface choice justified weakly | answer sounds tool-driven, not workload-driven | tie each boundary to timing and ownership |
| write contract ignores retries | duplicate side effects appear in review | add idempotency key and replay semantics |
| list contract allows unbounded scans | interviewer pushes on scale and abuse | cap parameters and move exports async |
| compatibility plan is missing | future changes sound hand-wavy | name additive strategy and sunset policy |

## Observability

- metric: retry rate and duplicate suppression hits
- metric: query rejection count by safety policy
- metric: request volume by version or schema
- log: normalized request shape and reject reason
- trace: sync request plus async continuation identifiers
- SLO: interface design preserves latency and correctness under expected retries and growth

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one focused drill scenario | easier to compare answers over time | less breadth per session | many tiny disconnected prompts |
| operational rubric | trains senior-level answers | feels stricter than endpoint-only reviews | API-shape-only feedback |
| deliberate time-boxing | improves interview pacing | some depth must be deferred | open-ended design with no wrap-up |

## Interview It

**Google framing:** "Design the interface layer for a task management platform." The signal is whether you connect contract choices to scale, client diversity, and future evolution.

**Cloudflare framing:** "Design the public and internal interfaces for an edge configuration product." The signal is whether you separate public API concerns from internal propagation and safety controls.

**Follow-ups:**
1. What if the public write path now needs aggressive client retries?
2. What if the list endpoint becomes the cost hotspot?
3. What if a downstream team needs replay and audit history?
4. What if the API must support old mobile clients for a year?
5. Which part of the answer would you deep dive if time allowed?

## Ship It

- `outputs/drill-worksheet-interface-design.md`
- `outputs/scoring-rubric-interface-design.md`

## Exercises

1. **Easy** — Do the drill on a comments API with one write and one list endpoint.  
2. **Medium** — Redo the drill for a partner-facing analytics export product.  
3. **Hard** — Run the drill on an edge-control product with public APIs and internal propagation events.  

## Further Reading

- [Google API Design Guide](https://cloud.google.com/apis/design) — a useful reference when reviewing the drill outcome  
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview pacing and deep-dive discipline  
