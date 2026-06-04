# Site

Optional static site that renders the course phases, lessons, and progress.

The template does not ship a site implementation — wire your own, or
copy `site/` from
[ai-engineering-from-scratch](https://github.com/rohitg00/ai-engineering-from-scratch/tree/main/site)
(MIT) and adapt.

Minimum pages used by upstream courses:

- `index.html` — landing page
- `catalog.html` — phase + lesson grid
- `lesson.html` — single-lesson reader
- `prereqs.html` — prereq map between phases
- `glossary.html` — pulls from `glossary/terms.md`
- `data.js` — generated from `phases/**/docs/en.md`
