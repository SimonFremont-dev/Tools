# AGENTS.md — working on the Tools repo itself

Tools is a Go library plus a template collection used to scaffold and run
projects quickly.

## Rules

- Templates live under `templates/<stage>/[<lang>/]`; always follow
  `docs/conventions.md`: placeholders `{{PROJECT}}`, `{{MODULE}}`, `{{LANG}}`,
  `{{YEAR}}`, `{{IMAGE}}`; `.tmpl` suffix is stripped on render; `*.sh` templates
  are written executable; `bootstrap/common/` is merged before `bootstrap/<lang>/`;
  language overlays are top-level `<lang>/` dirs of a stage and are merged over it.
- No template may keep an unresolved `{{PLACEHOLDER}}` — the engine fails on it.
- Prompts for AI agents live in `prompts/`; AI assets shipped into generated
  projects live in `templates/ai/` (skills under `.skills/`).
- Go templates for generated projects follow the Standard Go Project Layout
  (`cmd/`, `internal/`, ... — never `src/`).

## Checks before committing

```sh
gofmt -l .        # must be empty
go vet ./...
go test ./...
go run ./cmd/tools list
go run ./cmd/tools new demo --lang go --db postgres --module github.com/me/demo --stages dev,quality,build,cicd,kubernetes,ai,docs,database --dir /tmp/tools-smoke
```
