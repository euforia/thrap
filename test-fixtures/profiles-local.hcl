default = "local"

profiles {
    local {
        orchestrator = "docker"
        secrets      = "file"
        registry     = "docker"
        vcs          = "git"
    }
    // Remote pull (predefined)
    dev {
        orchestrator = "docker"
        secrets      = "file"
        registry     = "docker"
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
        orchestrator = "docker"
        secrets      = "file"
        registry     = "docker"
    }
}