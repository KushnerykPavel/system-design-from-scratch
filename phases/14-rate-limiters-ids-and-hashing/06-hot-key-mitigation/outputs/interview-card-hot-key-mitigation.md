# Interview Card — Hot-Key Mitigation

## Fast answer frame

1. Call out that average balance is not enough.
2. Separate read-heavy from write-heavy hotspots.
3. Choose among replication, coalescing, isolation, and admission control.
4. State the consistency trade-off explicitly.
5. Name the key metrics: top-key QPS, owner saturation, mitigation activation.

## Good closing line

"I would not overprovision the entire fleet for one pathological key if I can isolate or replicate it safely."
