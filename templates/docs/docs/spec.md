# {{PROJECT}} — project spec

Filled in by the `project-spec` skill (`.skills/project-spec`). To change the
setup, re-run the skill — it asks the questions and re-applies the stages.

- Name: {{PROJECT}}
- Language: {{LANG}}
- Module / import path: {{MODULE}}
- Container image: {{IMAGE}}
- Framework / main libraries: TODO
- Database(s): {{DB}} — one `atlas.<db>.hcl` per engine, history in `migrations/`;
  credentials in `.env.<db>.example` merged into gitignored `.env` via `make env`
- Other data stores & infrastructure: TODO (local services in `docker-compose.yml`)
- CI/CD: github-actions — see `PIPELINE.md` and `.github/workflows/`
- Deployment target: TODO (`none` | `docker-compose` | `kubernetes` — see `deployments/`)
- Stages applied from the Tools repo: bootstrap, dev, quality, build, cicd, kubernetes, ai, docs, database
- Skill set (`.skills/`): project-spec, tdd, clean-code, clean-architecture, sql-migrations, git-conventions, {{LANG}}
- Coverage alerts — tools used here with NO conventions in the Tools repo
  (researched good practices should be added as templates): TODO
