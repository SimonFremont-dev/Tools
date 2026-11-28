# Atlas project configuration - https://atlasgo.io/
# Database: MySQL. Schema history: migrations/ (commit atlas.sum).
# Credentials come from the environment - see .env.{{DB}}.example; never commit them.

env "dev" {
  url      = env("{{DBENV}}_DATABASE_URL")
  dev_url  = "docker://mysql:8"
  migration {
    dir = "file://migrations"
  }
}

env "staging" {
  url      = env("{{DBENV}}_STAGING_DATABASE_URL")
  dev_url  = "docker://mysql:8"
  migration {
    dir = "file://migrations"
  }
}

lint {
  latest  = 1
  dev_url = "docker://mysql:8"
}
