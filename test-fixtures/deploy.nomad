job "one-ingest" {
  region = "us-west-2"
  datacenters = ["us-west-2"]
  type = "service"

  constraint {
      attribute = "${meta.hood}"
      value = "io"
  }

  constraint {
      attribute = "${meta.env_type}"
      value = "<ENV_TYPE>"
  }

  update {
   stagger      = "20s"
   max_parallel = 1
  }

  group "stack" {
    count = 2

    constraint {
      operator = "distinct_hosts"
    }

    constraint {
      operator  = "distinct_property"
      attribute = "${attr.platform.aws.placement.availability-zone}"
    }

    vault {
      change_mode = "noop"
      env = false
      policies = ["read-secrets"]
    }

    task "api" {

      driver = "docker"

      config {
        image = "foo/bar:<APP_VERSION>"
        port_map {
            default = 9000
        }
        labels {
            service = "${NOMAD_JOB_NAME}"
        }
        logging {
          type = "syslog"
          config {
            tag = "${NOMAD_JOB_NAME}-${NOMAD_TASK_NAME}"
          }
        }
      }

      env {
        APP_VERSION =  "<APP_VERSION>"
      }

      service {
        name = "${JOB}-${TASK}"
        port = "default"
        tags = ["alert=${NOMAD_JOB_NAME}", "urlprefix-${NOMAD_JOB_NAME}-${NOMAD_TASK_NAME}.domain.<PUBLIC_TLD>/ auth=true"]

        check {
          type     = "http"
          path     = "/v1/status"
          interval = "10s"
          timeout  = "3s"
        }
      }

      template {
        data = <<EOF
{{ with printf "secret/%s" (env "NOMAD_JOB_NAME") | secret }}{{ range $k, $v := .Data }}{{ $k }}={{ $v }}
{{ end }}
        DB_NAME = "{{.Data.DB_NAME}}"
{{ end }}
EOF
        destination = ".env"
        env = true
      }

      resources {
        cpu    = 100 # MHz
        memory = 512 # MB

        network {
          mbits = 2
          port "default" {}
        }
      }
    }
  }
}
