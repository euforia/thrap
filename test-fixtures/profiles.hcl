default = "local"

profiles {
    local {
        orchestrator = "docker"
        secrets      = "local"
        registry     = "docker"
        vcs          = "git"
    }
    sandbox {
        orchestrator = "docker"
        secrets      = "local"
        registry     = "sandbox"
    }
    // Remote pull (predefined)
    dev {
        orchestrator = "nomad"
        secrets      = "local"
        registry     = "shared"
        meta {
            PUBLIC_TLD = "com"
            TLD        = "local"
            ENV_TYPE   = "test" 
        }
        variables {
            APP_VERSION = ""
        }
    }
    // Remote pull (predefined)
    live {
        orchestrator = "nomad"
        registry     = "shared"
        secrets      = "local"
    }
}