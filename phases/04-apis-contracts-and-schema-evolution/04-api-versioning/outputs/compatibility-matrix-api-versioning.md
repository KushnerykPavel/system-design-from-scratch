---
lesson: 04-api-versioning
focus: balanced
---

| Change | Usually compatible? | Notes |
|--------|----------------------|-------|
| add optional field | yes | clients ignore unknown fields if contract allows |
| remove field | no | breaks readers still depending on it |
| change field meaning | no | semantic break even if schema compiles |
| change sort semantics | usually no | affects behavior, paging, and client assumptions |

## Review prompts

- What versions are officially supported?
- What telemetry exists by version?
- How is deprecation announced and enforced?
