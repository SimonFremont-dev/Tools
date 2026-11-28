# Contributing

- Add or edit templates under `templates/<stage>/[<lang>/]` following [docs/conventions.md](docs/conventions.md).
- Keep every template free of leftover `{{PLACEHOLDER}}` after rendering; `go test ./...` covers this.
- Add or edit AI prompt templates under `prompts/`.
- Before committing: `gofmt -l .` (clean), `go vet ./...`, `go test ./...`.
- Open a pull request against `main`.
