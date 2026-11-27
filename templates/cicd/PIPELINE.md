# Pipeline

Standardized CI/CD for `{{PROJECT}}`. See `.github/workflows/` for the
implementation. Language- and tool-specific behavior plugs in through the
`*_CMD` variables in `Makefile.lang`, so the workflows stay generic.

## Merge gates (run on every pull request)

Protect `main` and require these six status checks:

| # | Gate | CI job | Default implementation |
|---|------|--------|------------------------|
| 1 | pre-checks (hygiene) | `pre-checks` | `pre-commit run --all-files` (whitespace, yaml/json, merge conflicts, conventional commits) |
| 2 | format | `format` | `make fmt-check` (`FMT_CHECK_CMD` per language) |
| 3 | lint / static analysis | `lint` | `make lint` (`LINT_CMD` per language; go also has `make lint-report` for the info/warning/error health summary) |
| 4 | security | `security` | gitleaks (secrets, all languages) + `make security` (dependency vulnerabilities, `SECURITY_CMD` per language) |
| 5 | tests | `test` | `make test` (`TEST_CMD` per language) |
| 6 | build | `build` | `make build` + `make docker-build` (also validates the Dockerfile) |

## After merge / on tags

| Stage | Workflow | Trigger | What it does |
|-------|----------|---------|--------------|
| publish image | `docker.yml` | push to `main`, tags `v*` | buildx build + push to `ghcr.io/<owner>/<repo>` (semver/sha/latest tags, GHA cache) |
| validate + deploy staging | `cd.yml` | push to `main` | `kustomize build` + `kubectl --dry-run=client` on the staging overlay, then the deploy step (fill in: kubectl/ArgoCD/Flux) |
| release | `release.yml` (go overlay) | tags `v*` | `goreleaser release --clean` |

## Customizing for a project

- Per language: set the `*_CMD` variables in `Makefile.lang` (dev stage of the
  Tools repo) and the linter configs (quality stage).
- Per tool, when a make target is not enough: add a job or a whole workflow
  under `templates/cicd/<lang>/.github/workflows/` in the Tools repo, then
  re-run `tools apply cicd --lang <lang> --dir . --force`.
- Required checks: in GitHub, Settings > Branches > Branch protection for
  `main` > "Require status checks to pass": `pre-checks`, `format`, `lint`,
  `security`, `test`, `build`.
