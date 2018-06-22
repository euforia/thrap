
manifest "thrap" {
  name = "thrap"

  components {
    api {
      name     = "${registry.addr}/thrap/api"
      type     = "api"
      language = "go"

      build {
        dockerfile = "api.dockerfile"
      }

      secrets {
        destination = "secrets.hcl"
        format      = "hcl"
      }

      head = true

      env {
        file = ".env"

        vars {
          APP_VERSION = ""
        }
      }
    }
  }

  dependencies {
    github {
      name = "github"
      version = "v3"
      external = true
      config {}
    }

    ecr {
      external = true
    }

    vault {
        name = "vault"
        version = "0.10.2"
    }

  }

}
