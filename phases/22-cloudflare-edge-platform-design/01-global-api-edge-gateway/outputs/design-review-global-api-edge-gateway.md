---
name: global-api-edge-gateway-review
description: Review prompt for a Cloudflare-style global API edge gateway design.
phase: 22
lesson: 01
---

Review the design against:

1. Where is TLS terminated and why?
2. Which policy checks run locally at the POP and which require remote state?
3. How is origin shielding handled?
4. How are retries bounded during partial origin slowness?
5. Can operators explain decisions by POP, region, route, and tenant?
