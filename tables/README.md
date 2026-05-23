# SIDC source tables

TSV tables vendored from upstream MIT-licensed projects by Måns Beckman:

- `app6b/` — from [spatialillusions/stanag-app6](https://github.com/spatialillusions/stanag-app6) (`tsv-tables/app6b/`).
- `app6d/` — from [spatialillusions/stanag-app6](https://github.com/spatialillusions/stanag-app6) (`tsv-tables/app6d/`).
- `app6e/` — from [spatialillusions/milstandard-e](https://github.com/spatialillusions/milstandard-e) (`tsv-tables/`).

These are the source of truth for the generated entity and modifier constants in the
`app6b` and `app6d` packages. Run `go generate ./...` to regenerate after updating.
