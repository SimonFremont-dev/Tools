# Tools

A library of tools to setup and run projects quickly.

Two parts:

- a Go library: `env` (load `.env` files), `command` (run/capture external commands), and the root `tools` package (the scaffold engine)
- a `templates/` collection organized by project stage, rendered by the CLI into new projects

## Usage

```sh
go build -o bin/tools ./cmd/tools

# see what is available
./bin/tools list

# scaffold a new Go project with the local-dev stage
./bin/tools new myapp --lang go --module github.com/me/myapp --stages dev --git

# add a stage to an existing project
./bin/tools apply dev --lang go --dir ./myapp
```

Placeholders (`{{PROJECT}}`, `{{MODULE}}`, ...) are documented in [docs/conventions.md](docs/conventions.md).
The master prompt template for having an AI set up a project from this repo is [prompts/setup-project.md](prompts/setup-project.md).

## Project Structure

tools/
- cmd/tools/		- CLI entry point (list, new, apply)
- command/		- run and capture external commands
- env/			- load KEY=VALUE pairs from env files
- prompts/		- prompt templates to give an AI (setup todos, reviews)
- templates/		- project templates by stage (bootstrap, dev, quality, build, cicd, kubernetes, ai, docs)
- docs/			- conventions for templates and placeholders
- scaffold.go		- embed + render engine for templates/
- go.mod			- Go module file
- LICENSE			- Apache 2.0
- README.md		- This file

## Development

```sh
go test ./...
```
