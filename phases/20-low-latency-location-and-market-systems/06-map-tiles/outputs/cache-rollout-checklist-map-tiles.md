# Map Tiles Cache Rollout Checklist

## Versioning

- Are tile asset URLs immutable?
- What manifest or metadata record points clients to the active version?
- How many active versions may coexist during rollout?

## Publish Safety

- Was a regional canary validated first?
- Are tile manifests checksummed?
- What is the rollback step if one region serves bad tiles?

## Cache Health

- Which zooms and regions are hot enough to prewarm?
- What is the edge hit rate by version after rollout?
- Did origin load spike because a cache purge was too broad?
