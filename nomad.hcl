job "thrap" {
    datacenters = ["us-west-2"]
    type = "service"

    group "head" {

        vault {
            change_mode = "noop"
            env = false
            policies = ["read-secrets"]
        }

        task "api" {
            driver = "docker"
            config {
                image = "${NOMAD_META_REGISTRY}/${NOMAD_META_PROJECT}/api:${NOMAD_META_APP_VERSION}"
                port_map {
                    default = 10000
                }
                volumes = [
                    "local/conf:/thrap/conf",
                    "thrap:/thrap/data"
                ]
            }

            // Configured providers
            template {
                data = <<EOF
                {{ with printf "%s" (env "NOMAD_META_SECRETS_PATH") | secret }}
                orchestrator {
                    dev {
                        provider = "nomad"
                        addr = "{{.Data.DevNomadAddr}}"
                    }
                    int {
                        provider = "nomad"
                        addr = "{{.Data.IntNomadAddr}}"
                    }
                    prod {
                        provider = "nomad"
                        addr = "{{.Data.ProdNomadAddr}}"
                    }
                }

                registry {
                    dockerhub {
                        provider = "dockerhub"
                    }
                    docker {
                        provider = "docker"
                    }
                    sandbox {
                        addr = "{{.Data.SandboxRegistryAddr}}"
                        provider = "ecr"
                        config {
                            region = "us-west-2"
                        }
                    }
                    shared {
                        addr = "{{.Data.SharedRegistryAddr}}"
                        provider = "ecr"
                        config {
                            region = "us-west-2"
                        }
                    }
                }

                secrets {
                    dev {
                        addr = "{{.Data.DevVaultAddr}}"
                        provider = "vault"
                    }
                    int {
                        addr = "{{.Data.IntVaultAddr}}"
                        provider = "vault"
                    }
                    prod {
                        addr = "{{.Data.ProdVaultAddr}}"
                        provider = "vault"
                    }
                }

                iam {
                    dev {
                        addr = "{{.Data.DevVaultAddr}}"
                        provider = "vault"
                    }
                    int {
                        addr = "{{.Data.IntVaultAddr}}"
                        provider = "vault"
                    }
                    prod {
                        addr = "{{.Data.ProdVaultAddr}}"
                        provider = "vault"
                    }
                }

                storage {
                    provider = "consul"
                    addr = "http://consul.service.{{ env "NOMAD_META_TLD" }}:8500"
                }
                {{ end }}
                EOF
                destination   = "local/conf/config.hcl"
                change_mode   = "signal"
                change_signal = "SIGINT"
            }
            // Creds for provider
            template {
                data = <<EOF
                {{ with printf "%s" (env "NOMAD_META_SECRETS_PATH") | secret }}
                registry {
                    sandbox {
                        key    = "{{.Data.SandboxRegistryKey}}"
                        secret = "{{.Data.SandboxRegistrySecret}}"
                    }
                    shared {
                        key    = "{{.Data.SharedRegistryKey}}"
                        secret = "{{.Data.SharedRegistrySecret}}"
                    }
                }
                vcs {
                    github {
                        token = ""
                    }
                }
                {{ end }}
                EOF
                destination   = "local/conf/creds.hcl"
                change_mode   = "signal"
                change_signal = "SIGINT"
            }
            // Profiles based on configured providers
            template {
                data = <<EOF
                {{ with printf "%s" (env "NOMAD_META_SECRETS_PATH") | secret }}
                default = "dev"
                profiles {
                    dev {
                        name         = "Development"
                        orchestrator = "dev"
                        secrets      = "dev"
                        registry     = "shared"
                        iam          = "dev"
                        meta {
                            PUBLIC_TLD    = "{{.Data.DevPublicTLD}}"
                            PUBLIC_DOMAIN = "{{.Data.DevPublicDomain}}"
                            TLD           = "{{.Data.DevPrivateTLD}}"
                            SECRETS_PATH  = ""
                            INSTANCE      = ""
                            PROJECT       = ""
                        }
                        variables {
                            APP_VERSION = ""
                        }
                    }

                    int {
                        name         = "Integration"
                        orchestrator = "int"
                        secrets      = "int"
                        registry     = "shared"
                        iam          = "int"
                        meta {
                            PUBLIC_TLD    = "{{.Data.IntPublicTLD}}"
                            PUBLIC_DOMAIN = "{{.Data.IntPublicDomain}}"
                            TLD           = "{{.Data.IntPrivateTLD}}"
                            SECRETS_PATH  = ""
                            INSTANCE      = ""
                            PROJECT       = ""
                        }
                        variables {
                            APP_VERSION = ""
                        }
                    }

                    prod {
                        name         = "Production"
                        orchestrator = "prod"
                        secrets      = "prod"
                        registry     = "shared"
                        iam          = "prod"
                        meta {
                            PUBLIC_TLD    = "{{.Data.ProdPublicTLD}}"
                            PUBLIC_DOMAIN = "{{.Data.ProdPublicDomain}}"
                            TLD           = "{{.Data.ProdPrivateTLD}}"
                            SECRETS_PATH  = ""
                            INSTANCE      = ""
                            PROJECT       = ""
                        }
                        variables {
                            APP_VERSION = ""
                        }
                    }
                }
                {{end}}
                EOF
                destination   = "local/conf/profiles.hcl"
                change_mode   = "signal"
                change_signal = "SIGINT"
            }

            service {
                name = "${TASK}"
                port = "default"
                tags = ["urlprefix-${NOMAD_META_PROJECT}.${NOMAD_META_DOMAIN}"]
                check {
                    type     = "http"
                    path     = "/v1/status"
                    interval = "20s"
                    timeout  = "3s"
                }
            }

            resources {
                cpu     = 200
                memory  = 128

                network {
                    mbits = 1
                    port "default" {
                        static = 10000
                    }
                }
            }
        }
    }

}