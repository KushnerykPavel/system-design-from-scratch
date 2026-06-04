# Cached Consistency Checklist

- State the source of truth.
- Name the user journey that needs the strongest freshness guarantee.
- Distinguish bounded stale reads from read-after-write guarantees.
- Decide whether monotonic reads matter across sessions or regions.
- Explain how invalidation or version metadata enforces the promise.
- Measure version skew and read-after-write mismatch.
- Say what degraded mode looks like when invalidation lags or fails.
