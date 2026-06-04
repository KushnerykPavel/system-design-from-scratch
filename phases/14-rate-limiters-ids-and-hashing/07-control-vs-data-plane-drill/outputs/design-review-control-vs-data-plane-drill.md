# Design Review — Control Plane vs Data Plane Drill

Use this worksheet during practice:

- What state must be live in the data plane, and what can arrive asynchronously?
- What is the acceptable policy staleness window?
- Which policy classes fail open, fail closed, or degrade differently?
- How is bundle validation, rollout, rollback, and versioning handled?
- How would you explain one inconsistent customer experience across two POPs?

Pushback to use in self-review:

- "If the control plane times out, does the request path still work?"
- "What proves a new bundle actually reached the fleet safely?"
