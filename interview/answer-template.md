# System Design Interview Answer Template

Use this for every practice round. It blends the common 4-step system design loop with mandatory senior-level depth.

## 1. Clarify (2-4 minutes)

- Who are the users?
- What is the core user journey?
- What is explicitly in scope?
- What is explicitly out of scope for v1?
- What scale are we designing for?
- Which non-functional requirements matter most: latency, availability, consistency, durability, cost, abuse resistance?

## 2. Restate and prioritize

State the problem back in one or two sentences.

Then rank:
1. primary functional requirement
2. top two non-functional constraints
3. one assumption that could change the design later

## 3. Back-of-the-envelope sizing

Fill rough numbers before architecture:

| Metric | Rough value | Notes |
|--------|-------------|-------|
| DAU / MAU | | |
| Peak QPS | | |
| Read/write ratio | | |
| Storage growth | | |
| Bandwidth / egress | | |
| Latency target | | |

If you do not know the exact numbers, choose a plausible range and say why.

## 4. High-level design

Draw the main components first:

```text
client -> edge -> app tier -> storage/cache/queue
```

Explain:
- request path
- write path
- read path
- one major async path

## 5. Deep dive on 1-2 critical areas

Pick only the areas that matter most for this prompt, such as:
- data model and indexes
- cache strategy
- partitioning and hotspot control
- consistency model
- retries and idempotency
- abuse prevention
- observability

Say why you chose these deep dives.

## 6. Failure modes

Name at least 3:
- component outage
- overload / hotspot
- stale or inconsistent data
- replay / duplicate writes
- regional failure

For each, state:
- how you detect it
- how you mitigate it
- whether the system degrades or fails closed/open

## 7. Observability and SLOs

Cover:
- key SLIs
- one or two SLOs
- alert conditions
- dashboard signals
- tracing / request correlation

## 8. Trade-offs

Use this structure:

| Decision | Why chosen | Main downside | Alternative rejected |
|----------|------------|---------------|----------------------|
|          |            |               |                      |

Good answers explicitly mention what they are giving up.

## 9. Wrap-up

State:
- current design strengths
- biggest known limitation
- one next-step improvement if given more time

## 10. Constraint-change follow-up

Always prepare for:
- 10x traffic
- stricter latency
- stricter consistency
- lower cost target
- new compliance or abuse requirement

Answer with: "What changes first, and what stays the same?"
