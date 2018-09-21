
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
  local {
    # addr = "http://127.0.0.1:8200"
    provider = "vault"
  }
}

storage {
  provider = "consul"
  addr = "http://127.0.0.1:8500"
}