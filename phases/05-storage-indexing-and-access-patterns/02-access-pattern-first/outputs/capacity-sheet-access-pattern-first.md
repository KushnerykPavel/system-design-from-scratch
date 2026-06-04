---
lesson: 02-access-pattern-first
---

| Metric | Value | Notes |
|--------|-------|-------|
| Primary read QPS | 90K req/s | top path should shape keys and indexes |
| Secondary read QPS | 8K req/s | important but should not dominate the model |
| Writes | 15K req/s | enough to make every extra index count |
| Peak factor | 6x | stress the hottest tenant and list patterns |
| Export volume | 2 TB/day | move to async path |
