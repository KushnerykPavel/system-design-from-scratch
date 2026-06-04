---
lesson: 05-audit-and-compliance
focus: balanced
---

| Decision | Benefit | Cost | Watch-out |
|----------|---------|------|-----------|
| Separate PII from finance records | smaller blast radius | more joins and policy links | token mapping becomes critical |
| Immutable audit store | trustworthy investigations | more storage and access tooling | log coverage must include reads |
| Policy-driven retention engine | consistent enforcement | higher implementation complexity | misclassification becomes dangerous |
| Archive cold records | lower hot-storage cost | slower retrieval | need searchable metadata |
