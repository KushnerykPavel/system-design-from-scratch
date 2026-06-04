# Cardinality Discipline Checklist

- Is the primary dashboard optimized for first-response triage?
- Are labels bounded by taxonomy rather than raw identifiers?
- Is `route_class` used instead of raw `path` where possible?
- Are `tenant_id`, `user_id`, or arbitrary error strings kept off default high-volume metrics?
- Does the dashboard begin with user impact before implementation detail?
- Are there region and dependency drill-downs without series explosion?
- Is query latency of the dashboard itself measured?
- Is there a review step for new labels before rollout?
