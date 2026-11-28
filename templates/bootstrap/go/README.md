# {{PROJECT}}

TODO: one-line description of {{PROJECT}}.

## Project Structure

Follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

- `cmd/{{PROJECT}}/` — application entry point (`main.go`)
- `internal/` — private application code (e.g. `internal/app/`)
- `pkg/` — public libraries safe to import from other projects
- `api/` — OpenAPI/Swagger specs, JSON schemas, protocol definitions
- `web/` — web assets, server-side templates, SPAs
- `configs/` — configuration file templates and defaults
- `build/` — packaging and CI files (`build/package`, `build/ci`)
- `deployments/` — deployment configs (docker-compose, kubernetes/helm, terraform)
- `scripts/` — build/install/analysis scripts (keeps the root Makefile small)
- `test/` — external test apps and test data
- `docs/` — design and user documentation
- root: `go.mod`, `Makefile`, `README.md`, `LICENSE`, `.gitignore`, `.editorconfig`

Do not create a `src/` directory.

## Development

```sh
make setup
make run
make test
```
