# Failure Checklist — Distributed KV Store

- Are replicas placed across real failure domains rather than only across hosts?
- Is the durability claim consistent with the configured replica count and write quorum?
- Can the system explain whether a returned read was potentially stale?
- Is anti-entropy lag observable and tied to paging thresholds?
- Are hot partitions detectable before they become full outages?
- Is rebalance bandwidth capped so repair does not starve serving traffic?
- Can the control plane revoke a bad partition map safely?
