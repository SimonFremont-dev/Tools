---
name: clean-code
description: Clean code practices — naming, function size, error handling, no dead code, comments that explain why not what. Language-agnostic.
---

# Clean code

## Naming

- Names reveal intent: `CalculateInvoiceTotal` over `Calc`; `isValid` over
  `check`. Spend the time here — it is the cheapest documentation.
- Booleans read as questions (`isReady`, `hasPermission`); collections as
  plurals (`cars`, `items`); avoid `data`, `info`, `manager`, `util`.
- Keep names consistent across the codebase (one word per concept).

## Functions and structure

- Small functions that do one thing and one level of abstraction.
- Early returns instead of deep nesting; guard clauses first.
- Parameters: few (0–2 ideal, 3 acceptable); no boolean flag parameters —
  split the function instead.
- No dead code, no commented-out code; version control remembers.

## Errors and comments

- Handle errors at the level that can decide what to do; never swallow them
  (`_ =`). Wrap with context when propagating.
- Comments explain **why**, not what: business rules, workarounds, non-obvious
  constraints. If the code says it, don't comment it.
- Keep the public API small and documented.

## Discipline

- `make fmt` before committing; `make fmt-check`, `make lint`,
  `make security` must stay green (CI gates).
- DRY, but not prematurely: duplicate twice, abstract on the third occurrence.
- Boy-scout rule: leave the code you touch cleaner than you found it.
