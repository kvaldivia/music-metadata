data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev?search_path=public"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "dev" {
  url = "postgresql://music-metadata:v9qsJRuw6e@db:5432/music-metadata?sslmode=disable"
}

env "prod" {
  url = "postgresql://music-metadata:v9qsJRuw6e@localhost:5432/music-metadata"
}
