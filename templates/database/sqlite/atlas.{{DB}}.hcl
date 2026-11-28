# Atlas project configuration - https://atlasgo.io/
# Database: SQLite (embedded, file-based). Schema history: migrations/ (commit atlas.sum).
# No server credentials needed; the database file is gitignored via *.db if added to .gitignore.

env "dev" {
  url      = "sqlite://file:{{PROJECT}}.db?_fk=1"
  dev_url  = "sqlite://file:dev.db?_fk=1"
  migration {
    dir = "file://migrations"
  }
}

env "staging" {
  url      = env("{{DBENV}}_STAGING_DATABASE_URL")
  dev_url  = "sqlite://file:dev.db?_fk=1"
  migration {
    dir = "file://migrations"
  }
}

lint {
  latest  = 1
  dev_url = "sqlite://file:dev.db?_fk=1"
}
