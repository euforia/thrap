default = "local"

profiles {
    local {
        orchestrator = "docker"
        secrets = "file"
        registry = "docker"
    }
    remote-registry {
        orchestrator = "docker"
        secrets = "file"
        registry = "ecr"
    }
    // Remote pull (predefined)
    dev {
        orchestrator = "nomad"
        registry = "ecr"
    }
    // Remote pull (predefined)
    live {
        orchestrator = "nomad"
        registry = "ecr"
    }
    // Example 
    custom {
        orchestrator = "nomad"
        secrets  = "vault"
        registry = "ecr"
    }
}