# Observability Checklist — Distributed Cache Cluster

- Track hit rate and miss rate by caller, tenant, and object class.
- Break misses down into expired, evicted, cold-node, and invalidated causes.
- Monitor origin offload as a first-class success metric, not just cache CPU.
- Page on miss storms only when paired with dangerous origin pressure.
- Watch eviction churn and average object age for signs of poor admission policy.
- Measure key movement and hit-rate regression after each topology change.
- Sample miss explanations so teams can tune keys and TTLs without guessing.
