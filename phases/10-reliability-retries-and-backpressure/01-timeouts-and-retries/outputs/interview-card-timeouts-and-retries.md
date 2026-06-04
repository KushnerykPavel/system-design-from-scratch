---
lesson: 01-timeouts-and-retries
focus: balanced
---

## Clarify first

- who owns the end-to-end deadline
- which failures are actually safe to retry
- how much fanout sits behind one caller request

## Must-size numbers

- baseline QPS and fanout width
- p95 and p99 dependency latency
- max attempts and worst-case retry amplification

## Core design

- propagate deadlines, not only timeouts
- derive hop-level timeouts from remaining budget
- cap retries and add jitter

## Failure probes

- every layer retries independently
- timeouts are shorter than real tails
- non-idempotent work is retried after ambiguity

## Trade-off summary

- transient recovery vs overload risk
- shorter timeout vs false timeout rate
- simplicity vs budget-aware policy
