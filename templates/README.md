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
- `cicd/` — pipeline templates per provider and language (P5)
- `kubernetes/` — manifests / Helm chart, secrets, Skaffold (P6)
- `ai/` — agent instructions and per-language skills copied into generated projects (P7)
- `docs/` — README / CONTRIBUTING / ADR templates for generated projects (P8)

Scaffold a project with the CLI:

```sh
tools new myapp --lang go --stages dev,quality,build
tools apply cicd --lang go --dir ./myapp
```
