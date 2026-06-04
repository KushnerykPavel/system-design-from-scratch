# Trade-off Matrix — Unique ID Generator

| Strategy | Benefit | Cost | Good fit |
|----------|---------|------|----------|
| DB sequence | simple, strongly ordered, easy to reason about | central bottleneck and dependency risk | modest write rate or strong ordering requirement |
| Snowflake-style | scalable local generation with rough time order | clock and worker coordination complexity | high write rate with sortability needs |
| Random ID | no central allocator and hard to guess | poor order and weaker storage locality | public IDs where uniqueness matters more than order |

Closing interview line:

"I would choose the minimum coordination that still satisfies the downstream assumptions people are quietly making about these IDs."
