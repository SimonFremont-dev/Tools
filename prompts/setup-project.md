# Prompt template: set up a new project with the Tools library

Fill in the placeholders, then give this todo list to your AI agent.
The agent should read `docs/conventions.md` and `templates/README.md` in
the Tools repository first.

## Project specification

- Project name: `{{PROJECT}}`
- Language: `{{LANGUAGE}}` (`go` | `node` | `python`)
- Module / import path: `{{MODULE}}`
- Framework: `{{FRAMEWORK}}`
- Database and dependencies: `{{DATABASE}}`
- CI/CD provider: `{{CICD}}` (`github-actions` | `gitlab-ci`)
- Deployment target: `{{DEPLOY_TARGET}}` (`none` | `docker-compose` | `kubernetes`)

Alternatively, hand this file to an agent that has the `project-spec` skill
(`templates/ai/.skills/project-spec`) — it asks the questions itself, applies
the stages, and alerts on tools missing from the Tools templates.

## Todo list (prioritized, execute in order)

1. Bootstrap the project: `tools new {{PROJECT}} --lang {{LANGUAGE}} --module {{MODULE}} --stages dev --git`.
   Go projects follow the Standard Go Project Layout (`cmd/`, `internal/`, `pkg/`, ...).
2. Fill in the local dev environment from `templates/dev` (Makefile targets for `{{LANGUAGE}}`, `.env.example`, `docker-compose.yml` for `{{DATABASE}}`).
3. Add quality tooling for `{{LANGUAGE}}` from `templates/quality`, including the conventional-branch pre-push hook (`make hooks`) and the conventional-commit pre-commit hook; see `templates/docs/docs/git-workflow.md`.
4. Add database management from `templates/database` with `--db {{DATABASE}}` (Atlas `atlas.hcl` + `migrations/`; postgres/mysql/sqlite).
5. Add build, release and container files from `templates/build` (multi-stage Dockerfile, `.dockerignore`, release config).
6. Add the CI/CD pipeline for `{{CICD}}` from `templates/cicd`, matching `{{LANGUAGE}}`.
7. If `{{DEPLOY_TARGET}}` is `kubernetes`, add manifests from `templates/kubernetes`; if `docker-compose`, extend `docker-compose.yml` with the app service.
8. Add AI assets from `templates/ai` (`AGENTS.md`, per-language and practice skills for `{{LANGUAGE}}`).
9. Complete the docs from `templates/docs` (README, CONTRIBUTING, git-workflow, first ADR for `{{FRAMEWORK}}`).
10. Verify: run lint, tests and build; resolve every `TODO` left by the templates.
