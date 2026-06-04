# Interview Card — Rate Limiter Primitives

## Fast answer frame

1. Clarify whether the limiter is protecting a backend, shaping traffic, or enforcing fairness.
2. Size burstiness and hot-key skew before naming the primitive.
3. Pick token bucket, sliding window, or leaky bucket based on desired behavior.
4. Name the hidden trade-off: fairness, memory, latency, or queue buildup.
5. Close with failure mode and observability.

## High-value phrases

- "Token bucket is a strong default when controlled burst is acceptable."
- "Sliding window buys closer rolling fairness at higher state cost."
- "Leaky bucket is useful when the downstream path needs smoothing, not just counting."
- "I would validate the choice with reject rate, backend saturation, and hot-key metrics."
