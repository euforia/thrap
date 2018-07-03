manifest "thrap" {
  name = "thrap"

  components {
    nomad {
      name    = "thrap/nomad"
      version = "0.8.4"
      type    = "api"

      ports {
        http = 4646
      }

      # build {
      #   dockerfile = "nomad.dockerfile"
      # }
    }

    vault {
      name    = "vault"
      version = "0.10.3"
      type    = "api"

      ports {
        default = 8200
      }

      env {
        vars {
          VAULT_DEV_ROOT_TOKEN_ID = "myroot"
        }
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
          # Should be available by default  
          STACK_VERSION = "${stack.version}"
          VAULT_ADDR    = "http://${comps.vault.container.ip}:${comps.vault.container.default.port}"
          NOMAD_ADDR    = "http://${comps.nomad.container.ip}:${comps.nomad.container.http.port}"
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

    nomad {
      name    = "vault"
      version = "0.10.3"
    }

    docker {
      name    = "docker"
      version = "1.37"
    }
  }
}
