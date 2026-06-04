# Interview Card - Leader-Follower Replication

## Say first

`The leader is the write serialization point. I want to define the commit rule, which reads may use followers, and how promotion avoids stale leadership.`

## Clarify

- Which reads are allowed to be stale?
- What write loss is acceptable during failover?
- Is promotion automatic or operator-gated?
- Does the same user bounce across replicas or regions?

## Commit choices

| Commit policy | Benefit | Cost | Best fit |
|---------------|---------|------|----------|
| leader local durable write | low latency | more failover loss risk | low-value metadata |
| leader plus 1 follower | better durability | slower writes | core application metadata |
| quorum-like durable replication | stronger safety | higher coordination cost | correctness-critical control plane |

## Close strong

`For critical reads I stay on the leader or require a minimum version. For cheaper reads I use followers, but I watch freshness age and gate promotion on replica position plus fencing.`
