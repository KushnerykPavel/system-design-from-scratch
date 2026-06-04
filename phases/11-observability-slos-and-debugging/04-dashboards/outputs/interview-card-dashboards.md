# Interview Card — Dashboards and Cardinality Discipline

## Triage panel order

1. user impact / SLO
2. volume and latency
3. errors by scope
4. dependencies
5. saturation or queueing
6. recent changes

## Safe default labels

- region
- route_class
- dependency_name
- status_class
- tenant_tier

## Avoid by default

- raw path
- tenant_id
- user_id
- session_id
- unbounded error strings
