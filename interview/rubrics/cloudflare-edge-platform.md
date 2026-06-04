# Cloudflare Edge / Platform System Design Rubric

Score each dimension from 1-4. Total: 36.

| Dimension | 1 (Poor) | 2 (Weak) | 3 (Strong) | 4 (Exceptional) |
|-----------|----------|----------|------------|-----------------|
| **Traffic realism** | Ignores traffic shape | Mentions scale vaguely | Reasons about hot paths, bursts, and edge traffic | Predicts line-rate or tail-latency breakpoints |
| **Cache and edge semantics** | Misunderstands caching | Basic cache discussion | Correct freshness and invalidation reasoning | Handles layered caches, regional drift, and abuse cases confidently |
| **Origin protection** | No protection model | Generic retries and LB | Good failover, shielding, and retry-budget thinking | Strong degradation story under slow-origin or partial-failure conditions |
| **Routing / regional reasoning** | No routing model | High level only | Correct POP/region traffic narrative | Connects routing decisions to cost, latency, and failure containment |
| **Observability** | No signal plan | Mentions logs or metrics | Actionable metrics, logs, traces, and alerts | Strong cardinality discipline and dashboard usefulness |
| **Abuse / security awareness** | Ignores attackers | Mentions WAF vaguely | Good abuse controls and trust boundaries | Anticipates adversarial behavior without overbuilding |
| **Performance / cost trade-offs** | No cost awareness | Generic "optimize later" | Ties design to latency, egress, and compute cost | Uses trade-offs to justify topology and cache strategy |
| **Communication** | Rambling | Some structure | Clear, paced, and practical | Feels like a senior engineer walking an incident-safe design review |

## Strong-positive signals

- Talks about cache hit rate and origin load together.
- Uses retry budgets instead of blind retries.
- Mentions blast radius by POP, region, or origin pool.
- Handles abuse and cost as first-class constraints.
- Explains what fails open vs fails closed.

## Strong-negative signals

- Assumes infinite origin capacity.
- Ignores cache invalidation semantics.
- Treats edge routing as just another load balancer.
- No story for partial regional failure.
- Mentions security only at the end as a bolt-on.
