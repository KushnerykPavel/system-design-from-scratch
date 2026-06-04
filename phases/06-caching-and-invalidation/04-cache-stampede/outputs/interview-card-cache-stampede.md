# Interview Card: Cache Stampede

## Strong answer shape

- separate steady-state hit ratio from miss-path behavior
- describe hot-key expiry and why it can overload the origin
- add request coalescing first
- add TTL jitter for synchronized expiries
- decide whether stale-while-revalidate is acceptable

## Senior signals

- "One miss is not the issue. Ten thousand identical misses at once are."
- "Coalescing protects the origin, but it adds waiter management and timeout choices."
- "If stale reads are not acceptable, I still need a strict miss-path control story."
