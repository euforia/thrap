manifest "thrap" {
  name = "thrap"

  components {
    vault {
      name    = "vault"
      version = "0.10.3"
      type    = "api"

      ports {
        default = 8200
      }
    }

    registry {
      name     = "registry"
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
          APP_VERSION = "${stack.version}"
          VAULT_ADDR  = "${comps.vault.container.default.addr}"
          NOMAD_ADDR  = ""
        }
      }
    }
  }

  dependencies {
    github {
      name     = "github"
      version  = "v3"
      external = true
      config   = {}
    }

    ecr {
      external = true
    }

    vault {
      name    = "vault"
      version = "0.10.3"
    }

    docker {
      name    = "docker"
      version = "1.37"
    }
  }
}
