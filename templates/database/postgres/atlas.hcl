# Atlas project configuration - https://atlasgo.io/
# Database: PostgreSQL. Schema history: migrations/ (commit atlas.sum).

env "dev" {
  url      = "postgres://dev:dev@localhost:5432/{{PROJECT}}?sslmode=disable"
  dev_url  = "docker://postgres:16"
  migration {
    dir = "file://migrations"
  }
}

env "staging" {
  url      = "env('STAGING_DATABASE_URL')"
  dev_url  = "docker://postgres:16"
  migration {
    dir = "file://migrations"
  }
}

lint {
  latest  = 1
  dev_url = "docker://postgres:16"
}
