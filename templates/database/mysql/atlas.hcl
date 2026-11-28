# Atlas project configuration - https://atlasgo.io/
# Database: MySQL. Schema history: migrations/ (commit atlas.sum).

env "dev" {
  url      = "mysql://dev:dev@localhost:3306/{{PROJECT}}"
  dev_url  = "docker://mysql:8"
  migration {
    dir = "file://migrations"
  }
}

env "staging" {
  url      = "env('STAGING_DATABASE_URL')"
  dev_url  = "docker://mysql:8"
  migration {
    dir = "file://migrations"
  }
}

lint {
  latest  = 1
  dev_url = "docker://mysql:8"
}
