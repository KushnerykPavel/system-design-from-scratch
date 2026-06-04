# Layered Cache Review Card

## Content classes

- immutable public assets
- mutable public content
- personalized responses
- internal or tenant-specific metadata

## Review prompts

- Which layers may cache each class safely?
- What is the cache key and auth context at each layer?
- How is mutable content revalidated or purged?
- Which headers reveal where the response was served from?
- What layer should an on-call engineer inspect first during stale-data reports?
