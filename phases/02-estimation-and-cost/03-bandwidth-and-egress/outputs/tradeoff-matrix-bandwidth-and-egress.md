# Trade-off Matrix — Bandwidth and Egress

| Lever | Benefit | Cost | Best when |
|------|---------|------|-----------|
| higher cache hit rate | lowers origin bytes and egress bills | harder invalidation and cache warming | content is reusable |
| compression | reduces transfer size | adds CPU and some latency | payload is text-heavy |
| edge compute or serving | protects origin and reduces central links | more distributed operational surface | latency and origin cost both matter |
| smaller payloads | cuts bandwidth at the source | may require product or schema changes | bloated responses dominate bytes |
