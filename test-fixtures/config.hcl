
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
    addr = "937954377342.dkr.ecr.us-west-2.amazonaws.com"
    provider = "ecr"
    config {
      region = "us-west-2"
    }
  }
  shared {
    addr = "583623634314.dkr.ecr.us-west-2.amazonaws.com"
    provider = "ecr"
    config {
      region = "us-west-2"
    }
  }
}

secrets {
  file {}
}