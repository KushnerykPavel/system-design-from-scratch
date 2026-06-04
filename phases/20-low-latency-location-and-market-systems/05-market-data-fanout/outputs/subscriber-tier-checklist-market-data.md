# Market Data Subscriber Tier Checklist

## Tier Definitions

- What latency target belongs to each tier?
- What replay window belongs to each tier?
- Which quotas are enforced on symbols, sessions, and retained bytes?

## Isolation

- Can a slow subscriber retain unlimited history?
- Is replay isolated from live delivery?
- Can one reconnect storm hurt unrelated regions?

## Policy

- Which feeds are allowed to be delayed or sampled?
- Which subscribers may receive raw depth versus aggregated products?
- What operator action pauses or downgrades abusive subscribers?
