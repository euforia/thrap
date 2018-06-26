job "one-ingest" {
  region = "us-west-2"
  datacenters = ["us-west-2"]
  type = "service"

  update {
   stagger      = "20s"
   max_parallel = 1
  }

  group "stack" {
    count = 5

    vault {
      change_mode = "noop"
      env = false
      policies = ["read-secrets"]
    }

    task "api" {

      driver = "docker"

      config {
        image = "myapp:mytag"
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

      service {
        name = "${JOB}-${TASK}"
        port = "default"
        tags = ["alert=${NOMAD_JOB_NAME}", "urlprefix-${NOMAD_JOB_NAME}-${NOMAD_TASK_NAME}.foo.com/ auth=true"]

        check {
          type     = "http"
          path     = "/healthz"
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
