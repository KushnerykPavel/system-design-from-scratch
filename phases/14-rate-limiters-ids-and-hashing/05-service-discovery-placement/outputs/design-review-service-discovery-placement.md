# Design Review — Service Discovery and Placement

Checklist:

- Does the design separate endpoint discovery from request placement policy?
- Where does routing logic live: client, sidecar, proxy, or edge service?
- How stale can cached endpoint or policy data become before it is unsafe?
- How are health, capacity, and rollout signals combined?
- What happens if the control plane is down but the fleet still has cached state?

Interview pushback:

- "Nearest is not always best if the nearest endpoint is already saturated."
- "A registry is a control-plane component; the data plane should degrade gracefully without it."
