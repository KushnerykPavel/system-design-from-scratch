# Interview Card — Logs, Traces, and Correlation IDs

## Minimum strong answer

- ingress request ID
- propagated trace context
- structured logs
- async causation metadata
- sampling policy
- redaction policy

## Good phrasing

- "Metrics tell me that something is wrong; traces and logs tell me which path failed."
- "Queues are not observability dead ends, so I would carry origin identifiers into job envelopes."
- "I would bias tracing toward failures and tail latency to control cost."
