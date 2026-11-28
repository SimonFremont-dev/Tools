# Conventions

## Layout

- `templates/<stage>/` — files copied into a generated project for that stage.
  Stages: `bootstrap`, `dev`, `quality`, `build`, `cicd`, `kubernetes`, `ai`, `docs`.
- `templates/<stage>/<lang>/` — optional language overlay, merged over `templates/<stage>/`.
- `templates/database/<db>/` — database overlay (`postgres`, `mysql`, `sqlite`),
  merged like a language overlay when `--db` is set. In general, top-level
  subdirectories of a stage named after a language or a database are overlays,
  never content. In the `database` stage each engine overlay ships
  `atlas.{{DB}}.hcl` and `.env.{{DB}}.example` (credential templates; `.env` itself stays gitignored).
- `templates/bootstrap/common/` — language-independent bootstrap files (LICENSE,
  .gitattributes, .editorconfig), merged before `bootstrap/<lang>/`; not listed as a language.
- `prompts/` — prompt templates given to an AI agent; they are not copied into projects.
- `docs/` — documentation for this repository itself.

## Placeholders

Rendered by the scaffold engine (`scaffold.go`) before files are written:

- `{{PROJECT}}` — project name
- `{{MODULE}}` — module / import path (defaults to the project name)
- `{{LANG}}` — language key (`go`, `node`, `python`)
- `{{YEAR}}` — current year
- `{{IMAGE}}` — container image reference: `ghcr.io/<owner>/<repo>` when
  `{{MODULE}}` is a `github.com/...` path, else the project name
- `{{DB}}` — database engine (`postgres` | `mysql` | `sqlite`, default `postgres`);
  `--db` accepts a comma-separated list for several engines, each rendered with its own `{{DB}}`
- `{{DBENV}}` — uppercased `{{DB}}`, used in env var names (`{{DBENV}}_DATABASE_URL`)

Placeholders also work in file and directory names, e.g. `cmd/{{PROJECT}}/main.go`.

## Rules

- Every file under `templates/` must render with no leftover `{{...}}`;
  the engine fails the whole copy when a placeholder is unresolved.
- Keep templates plain text; dotfiles are embedded via `all:templates`.
- A template file ending in `.tmpl` is written without the suffix. Use it
  for files this module's own Go tooling must not pick up (a template
  `go.mod`, template `*.go` sources).
- Template `*.sh` files and anything under `githooks/` are written executable (0755).
- Generated projects must be standalone: do not make them import this module.
- Go templates follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):
  `cmd/`, `internal/`, `pkg/`, `api/`, `configs/`, `build/`, `deployments/`, `scripts/`,
  `test/`, `docs/` — never a `src/` directory. Later stages map onto it: the `build`
  stage puts `Dockerfile`/`.dockerignore`/release config at the project root (the
  docker build context), `cicd/` puts pipeline files under `.github/` or `build/ci`,
  `kubernetes/` puts manifests under `deployments/`, `docs/` fills `docs/`.
