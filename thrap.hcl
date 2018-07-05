manifest "thrap" {
  name = "thrap"

  components {
    nomad {
      // Image name
      name = "nomad"

      // Image version
      version = "0.8.4"
      type    = "api"

      ports {
        http = 4646
      }

      build {
        dockerfile = "nomad.dockerfile"
      }
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
      // Image name that will be used. 
      // The final full image name will be <stack.id>/<name>:<stack.version>
      name = "registry"

      type     = "api"
      language = "go"

      build {
        dockerfile = "api.dockerfile"
      }

      ports {
        http = 10000
      }

      secrets {
        destination = ".thrap/creds.hcl"
        format      = "hcl"
      }

      head = true

      env {
        file = ".env"

        vars {
          # Should be available by default  
          STACK_VERSION = "${stack.version}"
          VAULT_ADDR    = "http://${comp.vault.container.addr.default}"
          NOMAD_ADDR    = "http://${comp.nomad.container.addr.http}"
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
