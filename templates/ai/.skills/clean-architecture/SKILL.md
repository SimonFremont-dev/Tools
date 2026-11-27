---
name: clean-architecture
description: Clean architecture — the dependency rule, layers (domain, use cases, adapters, frameworks), ports and adapters, constructor injection. Language-agnostic.
---

# Clean architecture

## The dependency rule

**Source code dependencies point inward, never outward.** Inner layers know
nothing about outer layers. In Go this is enforced structurally: domain code
must not import `api/`, storage drivers, or frameworks.

## Layers (inner → outer)

1. **Domain / entities** — business rules and types that change for business
   reasons only. No framework, no I/O. (Go: `internal/domain/...`)
2. **Use cases / application** — orchestrate entities to fulfill a business
   operation. Define the ports (interfaces) it needs. (Go: `internal/app/...`)
3. **Interface adapters** — controllers/handlers, presenters, gateways that
   convert between use-case formats and outer-world formats. (Go: `api/`,
   `internal/storage/`, `internal/...` adapters)
4. **Frameworks & drivers** — web framework, database, message queue, external
   APIs, the `main` wiring. (Go: `cmd/<app>/main.go` does the composition root)

## Ports and adapters

- Use cases depend on **interfaces they own** (ports), implemented by
  adapters in outer layers: `type CarRepository interface { ... }` defined next
  to the use case, implemented in `internal/storage/`.
- Inject dependencies via constructors — no globals, no service locators, no
  `init()` side effects. `main.go` wires the concrete implementations.

## Practical rules

- Business logic changes rarely; framework/DB changes often — keep them
  separable so swapping a driver or framework touches only the outer layer.
- I/O at the edges only: use cases return values/errors, adapters do the
  HTTP/DB/queue work.
- Tests: domain and use cases are pure unit tests (TDD skill); adapters get
  integration tests against real infrastructure (`make docker-up`).
- See `docs/adr/` for the recorded decisions; add an ADR when you deviate.
