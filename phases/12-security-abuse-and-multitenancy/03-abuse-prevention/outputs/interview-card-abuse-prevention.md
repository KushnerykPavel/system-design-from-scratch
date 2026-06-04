# Interview Card — Abuse Prevention and Rate Limiting Layers

## Strong answer shape

- Name the threat first.
- Separate edge, account, and tenant layers.
- Use cheap filters before expensive risk scoring.
- Design for false positives and enterprise support flow.
- Explain degraded mode for challenge or scoring outages.

## Common misses

- IP-only controls.
- No distinction between abuse and accidental burstiness.
- Treating login and read APIs the same.
- No review path for blocked healthy traffic.
