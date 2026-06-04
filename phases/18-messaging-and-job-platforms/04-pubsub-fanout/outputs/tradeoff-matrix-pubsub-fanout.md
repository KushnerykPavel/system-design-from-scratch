# Trade-off Matrix — Pub/Sub Fanout

| Decision | Benefit | Cost | When to prefer |
|----------|---------|------|----------------|
| Pull delivery | simple subscriber pacing | more polling overhead | internal service subscribers |
| Push delivery | lower subscriber complexity | retry, webhook, and isolation complexity | external integrations |
| Server-side filtering | lower subscriber CPU and egress | platform compute cost | common reusable filters |
| Local subscriber filtering | simpler platform | more duplicate egress and client work | few subscribers or simple internal use |
| Long replay windows | stronger recovery and audit | higher storage cost | premium or compliance-sensitive subscriptions |
| Tight backlog budgets | bounded shared cost | more subscriber drop/pause events | noisy multi-tenant environments |
