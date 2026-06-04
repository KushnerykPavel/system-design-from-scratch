---
lesson: 01-http-vs-grpc-vs-events
focus: balanced
---

| Interface | Best when | Main strength | Main risk | Watch closely |
|-----------|-----------|---------------|-----------|---------------|
| HTTP | public clients, broad compatibility, easy debugging | ecosystem reach | weaker typing, chatty payloads | p95 latency, client error mix |
| gRPC | internal low-latency service calls | typed contracts and efficiency | tighter release coupling | p99 latency, version skew |
| Events | fanout, async side effects, replay | decoupling and durability | eventual consistency, harder debugging | lag, redelivery, consumer failures |

## Decision prompts

- Does the caller need an immediate answer?
- How many consumers need the result?
- Who owns the contract and rollout?
- What failure mode is acceptable: timeout now or lag later?
