# Contributing

Thanks for contributing. The course gets better when readers send fixes,
port lessons to new languages, add test vectors, or contribute new lessons.

## What's welcome

- **Bug fixes** — broken code, wrong test vectors, factual errors.
- **New language ports** — add `code/main.<ext>` next to existing `code/main.<primary>`.
- **Test vectors** — RFC / NIST / academic sources only. Cite the source.
- **Attack demos** — adding an "Attack It" section to a primitive lesson (if course style includes it).
- **Glossary entries** — clarify a confused term in `glossary/terms.md`.
- **Quiz questions** — `quiz.json` additions per lesson.

## What's not welcome

- Production-grade rewrites of educational primitives (course is for learning, not deployment).
- Lesson reorderings without roadmap discussion first.
- Adding dependencies on niche or unmaintained libraries.

## Lesson rules

- Code must run. CI runs every `code/main.*` per lesson.
- Code must pass test vectors in `tests/vectors.json` if the lesson implements a primitive.
- No comments unless the *why* is non-obvious.
- Cite sources.

## Workflow

1. Open an issue first for non-trivial changes.
2. Branch off `main`.
3. Use `scripts/scaffold-lesson.sh` for new lessons.
4. One lesson per PR.

## Code review

PRs need:
- Passing CI
- One maintainer approval
- Test vectors for any primitive (if course requires)
- Updated `ROADMAP.md` row if adding a lesson
