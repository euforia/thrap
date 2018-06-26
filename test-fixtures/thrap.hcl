manifest "nomad" {
  name = "nomad"

  components {
    consul {
      name    = "consul"
      version = "1.1.0"
    }

    vault {
      name    = "vault"
      version = "0.10.3"
      head    = true

      env {
        vars {
          CONSUL_ADDR = "http://${comps.consul.addr}"
        }
      } // end env
    } // end vault

    app {
      name = "app"

      build {
        dockerfile = "api.dockerfile"
      }

      secrets {
        destination = "secrets.hcl"
        format      = "hcl"
      }

      head = true

      env {
        vars {
          STACK_VERSION = "${stack.version}"
          VAULT_ADDR    = "http://${comps.vault.addr}"
          CONSUL_ADDR   = "http://${comps.consul.addr}"

          # AWS_ECR_REGION = "${deps.ecr.region}"
        }
      } // end env
    } // end app
  } // end components

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

    docker {
      name    = "docker"
      version = "1.37"
    }
  }
}
