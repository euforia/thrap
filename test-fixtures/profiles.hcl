default = "local"

profiles {
    local {
        orchestrator = "docker"
        secrets      = "file"
        registry     = "docker"
        vcs          = "git"
    }
    sandbox {
        orchestrator = "docker"
        secrets      = "file"
        registry     = "sandbox"
    }
    // Remote pull (predefined)
    dev {
        orchestrator = "nomad"
        registry     = "shared"
    }
    // Remote pull (predefined)
    live {
        orchestrator = "nomad"
        registry     = "shared"
    }
}