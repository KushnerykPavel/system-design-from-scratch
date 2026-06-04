# Error Budget Checklist

- What exact user journey does this SLI measure?
- Does the SLI stay meaningful if one dependency degrades partially?
- Is latency part of the objective, not just binary success?
- Is the denominator defined clearly?
- Are exclusions explicit and reviewable?
- Is the target strict enough to matter but loose enough to be achievable?
- What short-window and long-window burn views exist?
- What changes when the budget burns: rollout pace, paging, or reliability work?
