# Time Assumption Checklist

- Is wall-clock time being used for display, ordering, expiry, or ownership?
- What breaks if two nodes disagree on time by more than expected?
- Does this path need logical versioning or a leader-assigned sequence instead?
- Are TTL and lease thresholds padded with a safety margin?
- Do ownership changes use fencing tokens?
- Which metrics show skew before it becomes a product incident?
- What is the fallback behavior if time sync health degrades?
