# Contributing to {{PROJECT}}

## Setup

```sh
make setup        # install toolchains/deps (see .tool-versions)
make docker-up    # local dependencies (postgres, redis, ...)
```

All commands go through the Makefile: `make run`, `make test`, `make lint`,
`make fmt-check`, `make security`, `make build`, `make docker-build`.

## Pull requests

1. Branch from `main`; use conventional commits (`feat:`, `fix:`, `docs:`, ...).
2. The merge gates in `PIPELINE.md` must be green: pre-checks, format, lint,
   security, test, build (`.github/workflows/ci.yml`).
3. Fill in `.github/pull_request_template.md`.
4. Follow the practice skills in `.skills/` (tdd, clean-code, clean-architecture).

## Conventions

See `AGENTS.md`, `docs/spec.md` and `docs/adr/`. This project is scaffolded from
the Tools library (https://github.com/SimonFremont-dev/Tools) — prefer changing
a template there over local one-offs.
