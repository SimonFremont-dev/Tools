---
name: project-spec
description: Interactively define or update a project's spec by asking questions (language, framework, tools, CI/CD, deployment), scaffold or update the project with the Tools templates, then alert on any tool that has no conventions in the Tools repo.
---

# Project spec

## 1. Ask the questions

Ask the user, batching related questions together:

- Project name and module/import path (e.g. `github.com/me/myapp`)?
- Language: `go` | `node` | `python`?
- Framework and main libraries (e.g. gorilla/mux, express, fastapi)?
- Database: `postgres` | `mysql` | `sqlite` | `none`? (drives the `--db` flag,
  the `atlas.hcl` datasource and the `database` stage)
- Other data stores and infrastructure dependencies (redis, kafka, ...)?
- CI/CD provider: `github-actions` (default) | `gitlab-ci` | other?
- Deployment target: `none` | `docker-compose` | `kubernetes`?
- Which stages to (re)apply: `dev`, `quality`, `build`, `cicd`, `kubernetes`, `ai`, `docs`?
- Which practice skills should this project enforce (the linkable skill set)?
  Available in `.skills/`: `tdd`, `clean-code`, `clean-architecture` (plus
  language skills like `go`). Remove the SKILL.md dirs of the ones the project
  does not follow, and list the chosen ones in `AGENTS.md`.

Record the answers as the project spec (a short markdown summary the user can
edit later, e.g. `docs/spec.md`).

## 2. Scaffold or update with the Tools templates

Tools repo: https://github.com/SimonFremont-dev/Tools (or the local checkout).
Build the CLI once: `go build -o bin/tools ./cmd/tools` inside the Tools repo.

- New project:
  `bin/tools new <name> --lang <lang> --db <db> --module <module> --stages <comma-separated stages> --git`
- Existing project:
  `bin/tools apply <stage> --lang <lang> --db <db> --dir . [--force]`

Then fill in what templates cannot know: framework wiring, business code,
real `.env` values, the deploy step in `cd.yml`, the ingress host.

## 3. Audit conventions coverage — ALERT on gaps

List every tool/technology now present in the project, detected from:
`go.mod` / `package.json` / `pyproject.toml`, `Dockerfile`, `docker-compose.yml`,
`.github/workflows/`, `deployments/`, `Makefile` targets, `.skills/`, plus the
answers from step 1 (frameworks, databases).

Check each against the Tools repo `templates/` tree:
`bootstrap/<lang>/`, `<stage>/` (dev, quality, build, cicd, kubernetes, ai, docs)
and their `<lang>` overlays.

For every tool WITHOUT a matching template or config, print an explicit alert:

```
⚠ no conventions in Tools for: <tool> (<why it matters>)
  → search good practices for <tool> and add a template under
    templates/<stage>/[<lang>/] in the Tools repo so future projects get it
```

Collect all alerts in a final "Coverage alerts" section of your answer.
