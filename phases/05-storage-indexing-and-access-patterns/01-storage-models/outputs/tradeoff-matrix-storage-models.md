| Store type | Best fit | Main strength | Main risk | Good fallback |
|------------|----------|---------------|-----------|---------------|
| Relational | invariant-heavy workflows | transactions and constraints | scaling complex joins and hot rows | add caches or derived indexes |
| KV | narrow lookup paths | predictable low-latency reads | poor support for rich queries | pair with relational source of truth |
| Document | evolving nested objects | natural object retrieval | index sprawl and weak query discipline | derive search or analytics views |
