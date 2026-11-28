---
name: sql-migrations
description: Database schema management with Atlas — migration history in migrations/ + atlas.sum, lint and apply. Works with postgres, mysql and sqlite (atlas.hcl).
---

# SQL migrations with Atlas

The project's database structure history lives in `migrations/`, managed by
[Atlas](https://atlasgo.io/):

- `NNNN_name.sql` files, created in order
- `atlas.sum` — integrity hash of the migration history, ALWAYS commit it
- `atlas.<db>.hcl` — per-environment config, one per engine (chosen with `--db`,
  comma-separated for several); URLs come from env vars, credentials live in
  `.env.<db>.example` merged into the gitignored `.env` (`make env`) — never commit them

## Workflow

- New migration: `make db-new name=add_users_table` (`atlas migrate new`)
- Lint: `make db-lint` (`atlas migrate lint --env dev`, replays on the dev db)
- Apply locally: `make db-apply` (`atlas migrate apply --env dev`)
- Status: `make db-status`; inspect schema: `atlas schema inspect --env dev`

## Rules

- Never edit or delete an already-applied/merged migration — add a new one.
- One logical change per migration, descriptive name.
- Migrations must replay cleanly on an empty dev database (`dev_url`).
- Run `make db-lint` before merging (CI gate).
