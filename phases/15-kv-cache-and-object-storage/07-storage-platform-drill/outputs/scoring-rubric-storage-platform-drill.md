# Scoring Rubric — Storage Platform Drill

## Strong answer
- Clarifies where blob bytes, metadata, and cache each belong
- Sizes upload, metadata, and hot-read paths separately
- Names durability promises and repair timing explicitly
- Explains at least one deep dive with operational detail
- Covers lifecycle, observability, and migration risks

## Weak answer
- Treats object, metadata, and cache layers as one fuzzy store
- Skips sizing and says only "use CDN" or "use S3"
- Describes durability only as replica count
- Ignores repair, retention, or destructive policy safety
- Lacks metrics that detect real failures
