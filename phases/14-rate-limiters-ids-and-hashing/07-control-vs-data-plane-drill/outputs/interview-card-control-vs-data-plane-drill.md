# Interview Card — Control Plane vs Data Plane Drill

## Fast answer frame

1. Define the boundary: policy authoring and distribution versus request evaluation.
2. Size the request path and policy-update path separately.
3. Cache compiled policy locally in the data plane.
4. State degraded modes and version-skew observability.
5. Close with rollout safety and rollback speed.

## Good closing line

"I want the control plane to be powerful without becoming a per-request dependency."
