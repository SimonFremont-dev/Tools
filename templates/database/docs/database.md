# Database (Atlas)

This file is scaffolded into `docs/database.md` of generated projects.

Schema-as-code with [Atlas](https://atlasgo.io/): the database structure
history lives in `migrations/` (`atlas.sum` = integrity hash, always commit it),
configured per environment in `atlas.hcl`. Database engine selected with
`--db postgres|mysql|sqlite` (see the per-engine `atlas.hcl` overlays).

```sh
make db-new name=add_users_table   # atlas migrate new
make db-lint                       # atlas migrate lint --env dev (needs dev db)
make db-apply                      # atlas migrate apply --env dev
make db-status                     # atlas migrate status --env dev
```

Install atlas via mise/asdf (`mise use atlas`) or see .tool-versions.
