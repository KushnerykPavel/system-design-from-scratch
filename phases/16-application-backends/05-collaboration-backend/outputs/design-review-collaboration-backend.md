# Design Review Prompt — Collaboration Backend

## Ask first
- What document model are we editing: plain text, rich text, or structured blocks?
- Are offline edits in scope, or only live sessions with short reconnect gaps?
- What convergence model is being claimed, and why does it fit the product?

## Review for risk
- Is the operation contract explicit about versioning and retries?
- How are snapshots used to cap replay cost?
- Can session-owner failover recover from the durable log without document corruption?
- Is presence clearly separated from authoritative document state?
- How would the system degrade if real-time fanout is impaired but storage remains healthy?
