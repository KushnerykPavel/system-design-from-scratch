# Trust Boundary Checklist

Use this when reviewing an interview answer or a real design.

## Identity

- Who authenticates end users?
- Who authenticates services?
- Which identities are short-lived versus static?
- Which boundary first validates untrusted external input?

## Authorization

- Which service owns the final resource-level authorization?
- Are tenant scope and resource scope both explicit?
- Are support and break-glass permissions separate?
- Can any upstream component accidentally over-assert permissions?

## Propagation

- What identity context is forwarded synchronously?
- What security context is attached to async jobs?
- Are forwarded claims signed, scoped, and time-bounded?
- Can downstream services explain why they allowed a request?

## Failure Review

- If the gateway is compromised, what still protects tenant data?
- If policy rollout is inconsistent, where is drift visible?
- If a worker replays a job, how is actor attribution preserved?
