---
lesson: 07-consumer-groups
focus: balanced
---

# Partition Planning Sheet

| Question | Answer |
|----------|--------|
| Required ordering scope | |
| Proposed partition key | |
| Initial partition count | |
| Near-term max consumer count | |
| Expected hot-key skew | |
| Rebalance strategy | |
| Lag SLO | |

## Review Prompts

- If one key becomes 20x hotter than average, what breaks first?
- If consumers double, is partition count still the ceiling?
- If product changes the ordering scope, how expensive is repartitioning?
- What metrics tell you the partition plan is wrong before users complain?
