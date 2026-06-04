# Real-Time Leaderboard Trade-off Matrix

| Decision | Faster | More Exact | Lower Cost | More Complex |
|----------|--------|------------|------------|--------------|
| Precompute top-N slices | Yes | Usually | Mixed | Yes |
| Exact around-me rank on demand | No | Yes | No | Mixed |
| Approximate around-me rank | Yes | No | Yes | Mixed |
| Delay suspicious scores before publication | Mixed | Yes | Mixed | Yes |
| Publish all scores immediately | Yes | No | Yes | No |

## Review Prompts

- Which readers need exactness, and which only need freshness?
- Can corrections happen after publication?
- Is segmentation a core product feature or a nice-to-have?
- Where does anti-cheat sit relative to the user-visible rank?
