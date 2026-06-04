# Interview Card — Consumer Backend Mock

## Use This For

- social feed
- chat or group messaging
- notifications
- collaboration or presence-heavy prompts

## Answer Order

1. Name the main user-visible promise.
2. Size read/write mix, skew, and retention.
3. Draw write path, read path, and async path.
4. Pick one deep dive that matches the real bottleneck.
5. Close with degraded mode, observability, and redesign.

## Good Deep Dives

- fanout strategy under skew
- source-of-truth and consistency boundary
- session or presence state placement
- ranking or enrichment pipeline

## Red Flags

- generic cache plus queue answer with no product contract
- no estimate for amplification
- no tolerated anomaly named
- no user-visible degraded behavior
