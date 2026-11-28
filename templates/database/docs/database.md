# Database (Atlas)

Schema-as-code with [Atlas](https://atlasgo.io/): the database structure
history lives in `migrations/` (`atlas.sum` = integrity hash, always commit it).

## Engines

One `atlas.<db>.hcl` per selected engine (`--db postgres,mysql` for several at
scaffold time; add more later with `tools apply database --db <db>`). Make
targets take the engine via `DB` (default: first selected):

```sh
make db-new name=add_users_table DB=postgres   # atlas migrate new
make db-lint  DB=postgres                       # atlas migrate lint --env dev
make db-apply DB=postgres                       # atlas migrate apply --env dev
make db-status DB=postgres                      # atlas migrate status --env dev
```

## Credentials and .env (security)

Never commit credentials: `.env` is gitignored. Credential templates live in
`.env.example` (app) and `.env.<db>.example` (per engine, rendered with
`{{DBENV}}_DATABASE_URL` etc.). Build your local `.env` once:

```sh
make env    # creates .env from .env.example + .env.*.example if missing
```

`atlas.<db>.hcl` reads URLs from the environment (`env("{{DBENV}}_DATABASE_URL")`),
and `docker-compose.yml` reads passwords from `.env` (`${POSTGRES_PASSWORD:?...}`).
Real values go only into the local `.env`.
