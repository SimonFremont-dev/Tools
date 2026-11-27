# AGENTS.md — agent instructions for {{PROJECT}}

This project was scaffolded from the Tools library
(https://github.com/SimonFremont-dev/Tools). Keep it that way.

## Conventions

- Language: {{LANG}}. Layout, lint config, Makefile targets and pipeline come
  from the Tools templates — do not invent parallel conventions.
- All commands go through the Makefile: `make setup run test lint fmt-check
  security build docker-build k8s-deploy lint-report` (see `make help`-style
  targets in the Makefile; per-language commands in `Makefile.lang`).
- CI/CD: see `PIPELINE.md`. The merge gates are the jobs of
  `.github/workflows/ci.yml`; never weaken them.
- Templates and placeholder rules: `docs/conventions.md` in the Tools repo.

## Changing the project setup with an agent

- Skills live in `.skills/` — start with `.skills/project-spec` to redefine
  the project spec (it asks questions, re-applies stages with the `tools` CLI,
  and alerts on tools that have no conventions in the Tools repo).
- Re-apply a stage: `tools apply <stage> --lang {{LANG}} --dir . [--force]`
  (build the CLI from the Tools repo: `go build -o bin/tools ./cmd/tools`).
