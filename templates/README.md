# Templates

Project templates organized by stage. See [docs/conventions.md](../docs/conventions.md)
for the layout rules and placeholders.

- `bootstrap/common/` — language-independent files (LICENSE, .gitattributes, .editorconfig)
- `bootstrap/<lang>/` — project skeleton per language (`go`, `node`, `python`);
  `go` follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- `dev/` — local development environment: `Makefile` + per-language `Makefile.lang`
  overlay, `.tool-versions`, `.env.example`, `docker-compose.yml`, devcontainer
- `quality/` — pre-commit, markdownlint; per language: `.golangci.yml` (all
  linters enabled, info/warning/error severity map + `scripts/lint-report.sh`
  health summary via `make lint-report`), eslint/prettier/commitlint, ruff/mypy
- `build/` — per-language multi-stage `Dockerfile` + `.dockerignore` at the project
  root, Go release config (`.goreleaser.yaml`)
- `cicd/` — standardized pipeline: `ci.yml` with 6 merge gates (pre-checks,
  format, lint, security, test, build — all via make targets), `docker.yml`
  (publish to ghcr), `cd.yml` (validate + deploy staging), go-only `release.yml`
  (goreleaser), `PIPELINE.md` documenting the gates, Dependabot per language,
  PR/issue templates
- `kubernetes/` — kustomize base + dev/staging/prod overlays under `deployments/`,
  root `skaffold.yaml`, secrets guidance in `deployments/README.md`
- `ai/` — `AGENTS.md` + `.skills/` skill set: `project-spec` (interactive spec
  Q&A, scaffolds via the CLI, alerts on tools without conventions in Tools),
  practice skills `tdd` / `clean-code` / `clean-architecture` (language- and
  tool-agnostic, linkable to any project's skill set) + per-language skills
  (e.g. `go`)
- `docs/` — `docs/` index + `spec.md` (project spec), ADR template, runbook,
  `git-workflow.md` (conventional commits + branches), and a generated-project `CONTRIBUTING.md`
- `database/` — Atlas schema management: `migrations/` + README, per-engine
  overlay (`postgres`, `mysql`, `sqlite`; `--db a,b` for several) with
  `atlas.{{DB}}.hcl` + `.env.{{DB}}.example` credential templates, `make db-*` / `make env`

Scaffold a project with the CLI:

```sh
tools new myapp --lang go --db postgres --stages dev,quality,build,cicd,kubernetes,ai,docs,database
```
