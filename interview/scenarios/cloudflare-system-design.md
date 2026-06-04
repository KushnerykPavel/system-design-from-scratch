# Cloudflare Edge / Platform Scenario Bank

Use for `/mock-interview cloudflare-system-design`. Reveal only the candidate prompt unless feedback mode is active.

## Scenario: global-api-gateway

Candidate prompt:
Design a global API gateway that terminates TLS at the edge, applies auth and rate limits, and forwards to regional origin clusters with low tail latency.

Hidden interviewer notes:
- Strong answers reason about POPs, cacheability, origin shielding, retries, and observability.
- Follow-up: what if one region becomes slow but not fully down?

## Scenario: cache-invalidation-control-plane

Candidate prompt:
Design a cache invalidation system that can purge customer content across a global edge within seconds while preventing abuse.

Hidden interviewer notes:
- Good deep dives: propagation, auth, batching, rate limits, regional drift, and verification.
- Follow-up: what changes if some customers need stronger freshness guarantees than others?

## Scenario: bot-mitigation-decision-engine

Candidate prompt:
Design a decision engine that evaluates incoming requests at the edge and chooses whether to allow, challenge, or block them.

Hidden interviewer notes:
- Strong answers discuss signal freshness, false positives, latency budget, and bypass safety.
- Follow-up: what changes if models update every few minutes?
