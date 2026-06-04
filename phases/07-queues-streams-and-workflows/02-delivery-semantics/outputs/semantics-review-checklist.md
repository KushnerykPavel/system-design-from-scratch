---
lesson: 07-delivery-semantics
focus: balanced
---

# Delivery Semantics Review Checklist

## Boundary

- What exact side effect is protected?
- Which part is broker-guaranteed versus application-guaranteed?
- Can the producer distinguish success, failure, and uncertainty?

## Producer

- Stable message ID on retry
- Retry budget or backoff policy
- Publish confirmation behavior documented

## Consumer

- Idempotent handler for repeated delivery
- Durable processed record or effect reference
- Ack only after durable effect boundary

## Failure Review

- Crash after side effect but before ack
- Producer timeout after successful publish
- Replay after dedupe TTL expiry
- External provider effect that is not transactional

## Observability

- Duplicate delivery rate
- Duplicate suppression hit rate
- Ack latency and oldest unacked message
- Reconciliation queue depth or manual repair count
