manifest "test" {
  name = "test"

  components {
    consul {
      name    = "consul"
      version = "1.1.0"
      type    = "api"

      ports {
        http = 8500
      }
    }

    vault {
      name    = "vault"
      version = "0.10.3"
      type    = "api"

      env {
        vars {
          CONSUL_ADDR = "http://${comp.consul.container.addr.http}"
        }
      } // end env
    } // end vault

    app {
      name     = "app"
      type     = "api"
      language = "go"

      build {
        dockerfile = "test.dockerfile"
        context    = "../test-fixtures"
      }

      secrets {
        destination = "secrets.hcl"
        format      = "hcl"
      }

      head = true

      env {
        vars {
          STACK_VERSION = "${stack.version}"
          VAULT_ADDR    = "http://${comp.vault.container.addr.default}"
          CONSUL_ADDR   = "http://${comp.consul.container.addr.http}"

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
