# Templates

Project templates organized by stage. See [docs/conventions.md](../docs/conventions.md)
for the layout rules and placeholders.

- `bootstrap/<lang>/` — project skeleton per language (`go`, `node`, `python`)
- `dev/` — local development environment (Makefile, `.env.example`, `docker-compose.yml`, devcontainer)
- `quality/` — linters, formatters, pre-commit (P3)
- `build/` — build targets, release config, Dockerfiles, `.dockerignore` (P4)
- `cicd/` — pipeline templates per provider and language (P5)
- `kubernetes/` — manifests / Helm chart, secrets, Skaffold (P6)
- `ai/` — agent instructions and per-language skills copied into generated projects (P7)
- `docs/` — README / CONTRIBUTING / ADR templates for generated projects (P8)

Scaffold a project with the CLI:

```sh
tools new myapp --lang go --stages dev
tools apply quality --lang go --dir ./myapp
```
