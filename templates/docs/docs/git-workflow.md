# Git workflow

## Branches — conventional branch names

https://conventionalbranch.org/ — `<type>/<short-description>`, e.g.
`feat/user-login`, `fix/crash-on-startup`, `chore/upgrade-deps`.

- Types come from conventional commits: `feat`, `fix`, `docs`, `style`,
  `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert` (+ `hotfix`).
- `main`/`master` are exempt.
- Enforced by `githooks/pre-push` (install with `make hooks`, i.e.
  `git config core.hooksPath githooks`) and by the `pre-checks` CI gate on PRs
  (`scripts/check-branch-name.sh`).
- One PR per branch; branch off `main`; rebase before merging.

## Commits — conventional commits v1.0.0

https://www.conventionalcommits.org/en/v1.0.0/ — `<type>[scope][!]: <description>`

- Examples: `feat(api): add login endpoint`, `fix: crash on startup`,
  `feat!: drop the v1 API` (`!` marks a breaking change).
- Imperative mood, lowercase, no trailing period; one logical change per commit.
- Enforced by the `conventional-pre-commit` hook (`.pre-commit-config.yaml`).
- Release notes and changelogs (goreleaser) are generated from these types.
