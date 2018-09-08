job "thrap" {
    datacenters = ["us-west-2"]
    type = "service"

    group "head" {

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
                {{ with printf "${NOMAD_META_SECRETS_PATH}" | secret }}
                orchestrator {
                    nomad {
                        provider = "nomad"
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
                        addr = "{{.Data.SandboxECRRegistry}}"
                        provider = "ecr"
                        config {
                            region = "us-west-2"
                        }
                    }
                    shared {
                        addr = "{{.Data.SharedECRRegistry}}"
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
                {{ end }}
                EOF
                destination   = "local/conf/config.hcl"
                change_mode   = "signal"
                change_signal = "SIGINT"
            }
            // Creds for provider
            template {
                data = <<EOF
                {{ with printf "${NOMAD_META_SECRETS_PATH}" | secret }}
                registry {
                    sandbox {
                        key    = "{{.Data.SandboxRegKey}}"
                        secret = "{{.Data.SandboxRegSecret}}"
                    }
                    shared {
                        key    = "{{.Data.SharedRegKey}}"
                        secret = "{{.Data.SharedRegSecret}}"
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
                default = "dev"
                profiles {
                    dev {
                        name         = "Development"
                        orchestrator = "nomad"
                        secrets      = "dev"
                        registry     = "shared"
                        meta {
                            PUBLIC_TLD   = ""
                            TLD          = ""
                            SECRETS_PATH = ""
                            INSTANCE     = ""
                            PROJECT      = ""
                        }
                        variables {
                            APP_VERSION = ""
                        }
                    }

                    int {
                        name         = "Integration"
                        orchestrator = "nomad"
                        secrets      = "int"
                        registry     = "shared"
                        meta {
                            PUBLIC_TLD   = ""
                            TLD          = ""
                            SECRETS_PATH = ""
                            INSTANCE     = ""
                            PROJECT      = ""
                        }
                        variables {
                            APP_VERSION = ""
                        }
                    }

                    prod {
                        name         = "Production"
                        orchestrator = "nomad"
                        secrets      = "prod"
                        registry     = "shared"
                        meta {
                            PUBLIC_TLD   = ""
                            TLD          = ""
                            SECRETS_PATH = ""
                            INSTANCE     = ""
                            PROJECT      = ""
                        }
                        variables {
                            APP_VERSION = ""
                        }
                    }
                }
                EOF
                destination   = "local/conf/profiles.hcl"
                change_mode   = "signal"
                change_signal = "SIGINT"
            }

            service {
                name = "${TASK}"
                port = "default"
            }
        }
    }

}