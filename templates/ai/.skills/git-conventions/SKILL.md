---
name: git-conventions
description: Git conventions for this project — conventional commits (conventionalcommits.org v1.0.0) and conventional branch names (conventionalbranch.org). Enforced by pre-commit hook, pre-push hook and the pre-checks CI gate.
---

# Git conventions

## Commits — conventional commits v1.0.0

https://www.conventionalcommits.org/en/v1.0.0/

- Format: `<type>[optional scope][!]: <description>`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`,
  `ci`, `chore`, `revert`. `!` before the colon = breaking change.
- Examples: `feat(api): add login endpoint`, `fix: crash on startup`.
- Imperative mood, lowercase, no trailing period, one logical change per commit.
- Enforced by the `conventional-pre-commit` pre-commit hook
  (`.pre-commit-config.yaml`); changelogs/releases are generated from them.

## Branches — conventional branch names

https://conventionalbranch.org/

- Format: `<type>/<short-description>` with the same types,
  e.g. `feat/user-login`, `fix/crash-on-startup`, `chore/upgrade-deps`.
- Enforced by `githooks/pre-push` (install with `make hooks`) and by the
  `pre-checks` CI gate (`scripts/check-branch-name.sh`); `main`/`master` exempt.
- One PR per branch, branched off `main`, rebased before merge.
