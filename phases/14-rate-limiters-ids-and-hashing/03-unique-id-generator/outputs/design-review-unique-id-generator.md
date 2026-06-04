# Design Review — Unique ID Generator

Use this checklist when reviewing an ID strategy:

- What exact properties are required: uniqueness, ordering, locality, or low guessability?
- Is the chosen strategy a hot dependency at projected peak write rate?
- What happens during clock regression, node restart, or worker-ID collision?
- Are public IDs and storage keys intentionally the same, or just accidentally the same?
- What metrics would detect generator regressions before customers notice?

Common pushback:

- "We do not need sortable IDs if nothing consumes order from them."
- "A public identifier should not leak internal sequence unless that trade-off is deliberate."
