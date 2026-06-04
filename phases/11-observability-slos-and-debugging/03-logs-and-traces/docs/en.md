# Logs, Traces, and Correlation IDs

> Metrics tell you that something is wrong. Logs and traces explain which path was wrong.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design log and trace instrumentation that preserves request narrative across services without exploding cost or leaking sensitive data.
**Prerequisites:** `04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events`, `07-queues-streams-and-workflows/06-outbox-and-cdc`, `11-observability-slos-and-debugging/02-metrics-that-matter`
**Estimated time:** ~75 min
**Primary artifact:** correlation checker + tracing checklist

## The Problem

Interview answers often say "we'll add logs and traces" as if that is enough. In practice:

- logs without request identity are hard to join
- traces without sampling strategy are expensive or incomplete
- mixed sync/async systems break narrative continuity
- careless logging can leak secrets or PII

This lesson focuses on preserving enough narrative to debug cross-service failures quickly.

## Clarify

- Is the main path synchronous RPC, event-driven workflow, or both?
- Which hops are controlled by your team versus third-party dependencies?
- Do you need per-request debugging or fleet-level statistical tracing?
- Are there privacy constraints on payload logging?

If the interviewer is vague, assume a front-door API that triggers multiple RPCs and an asynchronous queue-backed side effect.

## Requirements

### Functional

- Propagate correlation context across request hops.
- Emit structured logs with enough fields to reconstruct failures.
- Sample traces in a way that keeps rare failures explainable.

### Non-functional

- Bound telemetry cost and cardinality.
- Prevent secrets and sensitive payloads from entering logs.
- Preserve context across asynchronous boundaries and retries.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Peak request rate | 100K req/s | forces selective logging and trace sampling |
| Average service fanout | 5 downstream hops | increases value of shared trace context |
| Async work share | 30% of requests enqueue follow-up jobs | request identity must survive into workers |
| Debug sample rate | 0.1%-1% normal, higher on error | cost and diagnostic depth trade off directly |
| Rough cost | log storage + indexing + trace retention | over-instrumentation can become an infrastructure incident |

## Architecture

```text
client request
  -> gateway adds request_id / trace_id
  -> service A logs structured fields
  -> service B / service C receive propagated context
  -> async job carries correlation metadata
  -> traces and logs join under shared identifiers
```

Design principles:

1. one stable request identifier at the entry point
2. child spans for important downstream hops
3. structured logs, not free-form narrative only
4. payload redaction rules by default
5. error-triggered trace retention or higher sample policy

## Data Model & APIs

Useful fields:

- `trace_id`
- `span_id`
- `request_id`
- `tenant_tier`
- `route_class`
- `dependency`
- `attempt`
- `error_class`

Useful interfaces:

- `InjectContext(headers, traceContext)`
- `ExtractContext(headers)`
- `LogEvent(fields)`
- `ShouldSample(routeClass, statusClass, latencyMs)`

For async systems, carry:

- `causation_id`
- `job_id`
- `origin_request_id`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| correlation ID not propagated across one hop | logs stop joining after a dependency boundary | enforce middleware for propagation and test it |
| full payload logging leaks secrets | security review or incident reveals sensitive fields in logs | structured logging with explicit allowlists and redaction |
| sampling drops the interesting failures | low trace availability for slow/error requests | tail-based or error-biased sampling |
| async worker logs cannot be tied back to request | queue event has no origin metadata | carry origin request and causation IDs in job envelope |

## Observability

- metric: sample rate, dropped log count, and trace ingestion volume
- metric: propagation failure count where downstream span lacks parent context
- log: structured event fields with stable request and dependency identifiers
- trace: end-to-end critical path including queue enqueue/dequeue spans
- SLO: observability completeness can have its own internal objective for critical paths

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| structured logs | queryable and joinable events | schema discipline required | human-only free-text logs |
| error-biased tracing | better incident usefulness per dollar | incomplete coverage of healthy traffic | trace everything forever |
| correlation through async envelopes | preserves narrative across workflows | slightly larger messages and more code paths | treating queues as observability dead ends |

## Interview It

**Google framing:** "How would you instrument a multi-service request path so an on-call engineer can debug user failures?" The signal is whether you discuss propagation, schema, and sampling together.

**Cloudflare framing:** "How would you trace and log a global edge request through shared services and asynchronous systems?" The signal is whether you handle partial control, privacy, and cost realistically.

**Follow-ups:**
1. How would you trace a request that fans out to five services and one queue?
2. When is tail-based sampling worth the extra complexity?
3. What do you absolutely avoid logging?
4. How do retries affect trace readability?
5. What changes when a third-party dependency will not propagate your context?

## Ship It

- `code/main.go`
- `code/main_test.go`
- `outputs/correlation-id-checklist.md`
- `outputs/interview-card-logs-and-traces.md`

## Exercises

1. **Easy** — Define the minimum structured log fields you want at an API gateway and explain why.
2. **Medium** — Show how a request ID should survive from synchronous API call into an async worker.
3. **Hard** — Redesign the approach when full traces are too expensive but incidents still require cross-service debugging.

## Further Reading

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/) — practical baseline for logs, metrics, and tracing conventions
- [Distributed Tracing in Practice](https://queue.acm.org/detail.cfm?id=3526967) — useful trade-off framing for real-world tracing
