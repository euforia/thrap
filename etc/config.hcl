#
# Backend config
#

stack {
    requirements {
        files {
            README.md = "README.md"
            Makefile  = "Makefile"
        }
    }
}

webservers {
    nginx {
        Versions = ["~> 1.15-alpine"]
        Image = "nginx"
        DefaultVersion = "1.15-alpine"
    }
}

datastores {
    mysql {
        Versions = []
        Image    = ""
    }
    postgresql {
        Versions = []
        Image    = ""
    }
    elasticsearch {
        Versions = []
        Image    = ""
    }
    cockroach {
        Versions = ["~>  v2.0.2"]
        DefaultVersion = "v2.0.2"
        Image = "cockroachdb/cockroach"
    }
}
