# Correlation ID Checklist

- Is a stable request identifier created at ingress?
- Is trace context propagated across all RPC hops?
- Do async messages carry origin request or causation identifiers?
- Are logs structured with stable fields instead of free-form text only?
- Are retries and attempt numbers recorded explicitly?
- Are sensitive fields redacted or allowlisted?
- Is trace sampling biased toward errors or slow requests?
- Can one on-call engineer reconstruct a failed request path quickly?
