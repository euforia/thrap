package config

import (
	"testing"

	"github.com/euforia/hclencoder"
	"github.com/hashicorp/hcl"
	"github.com/stretchr/testify/assert"
)

var testConfig1 = `
registry {
    ecr {}
}

orchestrator {
    nomad {}
}

secrets {
    vault {
        addr = "http://vault:4646"
    }
}

vcs {
    github {
        addr = "github.com"
        username = "euforia"
        repo {
            name = "thrap"
            owner = "euforia"
        }
    }
}
`

var testConfig2 = `
registry {
    ecr {
        config {
            region = "us-west-2"
        }
    }
}

orchestrator {
    nomad {
        addr = "http://nomad:4646"
    }
    docker {}
}

secrets {
    vault {}
}

vcs {
    github {
        addr = "github.com"
        username = "euforia"
        repo {
            name = "thrap"
            owner = "new"
        }
    }
}
`

func Test_Config(t *testing.T) {

	var conf Config
	err := hcl.Decode(&conf, testConfig1)
	assert.Nil(t, err)
	assert.Equal(t, "thrap", conf.VCS["github"].Repo.Name)
	assert.Equal(t, "http://vault:4646", conf.Secrets["vault"].Addr)
}

func Test_Config_Merge(t *testing.T) {
	var c1, c2 Config

	err := hcl.Decode(&c1, testConfig1)
	assert.Nil(t, err)

	err = hcl.Decode(&c2, testConfig2)
	assert.Nil(t, err)

	c1.Merge(&c2)

	vcsp := c1.DefaultVCS()
	secp := c1.DefaultSecrets()
	regp := c1.DefaultRegistry()

	assert.Equal(t, "thrap", vcsp.Repo.Name)
	assert.Equal(t, "http://vault:4646", secp.Addr)
	assert.Equal(t, "new", vcsp.Repo.Owner)
	assert.Equal(t, "us-west-2", regp.Config["region"])
	assert.NotNil(t, c1.Orchestrator["docker"])
}

func Test_Config_Encode(t *testing.T) {
	conf := DefaultConfig()
	b, _ := hclencoder.Encode(conf)
	t.Logf("%s", b)
}
