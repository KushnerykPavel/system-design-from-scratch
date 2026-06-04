---
lesson: 01-http-vs-grpc-vs-events
focus: balanced
---

## Clarify first

- Is the interaction synchronous or asynchronous?
- Are external clients involved?
- Is there fanout or replay?

## Must-size numbers

- front-door QPS
- downstream fanout
- lag tolerance
- target latency budget

## Core design

- HTTP for public compatibility
- gRPC for controlled internal hot paths
- events for downstream side effects and fanout

## Failure probes

- what happens on retry
- what happens on duplicate delivery
- how consumers recover from lag

## Trade-off summary

- request-response is simpler to reason about
- events are better for decoupling but harder to debug
