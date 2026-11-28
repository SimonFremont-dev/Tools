# Migrations (Atlas)

- One logical change per migration, named descriptively: `NNNN_add_users_table.sql`
- `atlas.sum` is the integrity hash of the whole history — always commit it.
- Never edit or delete a migration that was already applied or merged; add a new one.
- Migrations are replayed from scratch against the dev database (`dev_url` in
  atlas.hcl), so they must be safe to re-run on an empty schema.

Commands (see `make db-*`):

```sh
atlas migrate new <name>     # create
atlas migrate lint --env dev # lint (CI gate)
atlas migrate apply --env dev
atlas migrate status --env dev
atlas schema inspect --env dev
```
