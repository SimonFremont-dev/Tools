---
name: go-conventions
description: Go conventions for projects scaffolded from the Tools library — Standard Go Project Layout, lint severity map, make targets, pipeline gates.
---

# Go conventions

- Layout follows the Standard Go Project Layout: `cmd/<app>/main.go` is a thin
  entry point calling `internal/...`; public reusable code goes in `pkg/`;
  never create a `src/` directory.
- Lint: `.golangci.yml` enables ALL linters with an error/warning/info severity
  map. `make lint` runs it; `make lint-report` prints the per-severity health
  summary (error = must fix, warning = should fix, info = opinionated).
- Format: `make fmt-check` (CI gate) / `make fmt` (golangci-lint fmt).
- Tests: `make test` (`go test ./...`); keep `internal/` packages unit-tested.
- Build/release: `make build` → `bin/<project>`; `make release` runs a
  goreleaser snapshot; real releases happen on `v*` tags via `release.yml`.
- Pipeline: see `PIPELINE.md` — the six CI gates must stay green to merge.
- Security: `make security` runs `govulncheck ./...`.
