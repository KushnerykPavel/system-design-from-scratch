---
lesson: 05-edge-compute
---

## Put it at the edge when

- logic is stateless or cache-backed
- the request can usually finish without a remote fetch
- policy bundles are small and versionable
- rollback must be fast but operationally tractable

## Keep it central when

- writes need strong coordination
- state is large and highly mutable
- the edge would just proxy to home-region data on most requests
- compliance or audit boundaries require tighter control
