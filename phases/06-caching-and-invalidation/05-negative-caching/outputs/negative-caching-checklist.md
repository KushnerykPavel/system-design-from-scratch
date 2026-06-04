# Negative Caching Checklist

- Separate stable absence from transient backend failure.
- Use shorter TTLs when objects may appear soon after an initial miss.
- Invalidate cached `404` on create or restore.
- Be cautious caching `403` unless policy versioning is clear.
- Avoid caching `500` by default, or keep it extremely short-lived.
- Measure stale-not-found and stale-deny incidents.
- Pair backend protection with abuse controls for random-key traffic.
