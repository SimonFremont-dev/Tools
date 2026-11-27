---
name: tdd
description: Test-driven development workflow — red/green/refactor, behavior-focused tests, when and how to run them in this project (make test).
---

# TDD — test-driven development

## The cycle (red → green → refactor)

1. **Red**: write one small failing test for the next behavior. Run
   `make test` (or the focused equivalent) and confirm it fails for the right
   reason.
2. **Green**: write the minimum code to make it pass. Run again — green.
3. **Refactor**: improve the code (naming, structure) with tests green.
   Run again — still green. Never refactor without green tests.

## Rules

- One behavior per test; test names describe the behavior
  (e.g. `TestCreateCar_RejectsDuplicateNumberplate`).
- Test behavior, not implementation: no asserting on private internals when a
  public behavior covers it.
- Arrange/Act/Assert structure; keep fixtures in `testdata/` where relevant.
- Unit tests for domain/use-case logic; integration tests against real
  dependencies only at the edges (storage, HTTP), using the
  `docker-compose.yml` services (`make docker-up`).
- Coverage is a signal, not a target: use it to find untested behavior, not to
  hit a number. `make lint-report` shows the project's health by severity.
- If a test is slow or flaky, fix it or move it — a suite people skip is
  worse than no suite.
