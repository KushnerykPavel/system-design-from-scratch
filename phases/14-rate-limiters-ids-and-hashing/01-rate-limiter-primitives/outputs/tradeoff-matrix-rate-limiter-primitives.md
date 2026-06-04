# Trade-off Matrix — Rate Limiter Primitives

| Primitive | Best for | Main benefit | Main cost | Watch-out |
|-----------|----------|--------------|-----------|-----------|
| Token bucket | burst-tolerant API limits | cheap state and intuitive burst semantics | rolling fairness is approximate | can still allow short spikes that hurt a fragile backend |
| Sliding window | stricter fairness over time | closer to true rolling-window behavior | more state or more approximation work | can get expensive for very hot keys |
| Leaky bucket | smoothing toward downstream workers | steadier output rate | queue delay becomes part of the system | shaping is not the same as abuse enforcement |

Use this matrix in interviews:

- Start with workload shape: burst tolerance, fairness, or smoothing.
- State the hot-path cost of the chosen primitive.
- Mention what extra mechanism is needed if one primitive does not solve the whole problem.
