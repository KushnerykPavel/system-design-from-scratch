# Rotation Checklist

## Before rollout

- Is the credential short-lived where feasible?
- Can workloads refresh without restart?
- Is there a staged overlap window?
- Is revocation tested, not just documented?

## During rollout

- Watch refresh latency and failure rate.
- Track old-version versus new-version usage.
- Look for handshake failures and auth spikes.
- Keep rollout scoped by service or region.

## After rollout

- Confirm old versions are no longer accepted.
- Review stragglers still using stale material.
- Record time-to-rotate and time-to-revoke.
- Capture what would break during an emergency rotation.
