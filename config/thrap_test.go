package config

import (
	"fmt"
	"testing"

	"github.com/euforia/hclencoder"
	"github.com/hashicorp/hcl"
	"github.com/stretchr/testify/assert"
)

var testCreds1 = `
registry ecr {
    key = "key"
    secret = "secret"
    token = "token"
}

orchestrator nomad {
    token = ""
}

secrets vault {
    token = ""
}

vcs github {
    token = "token"
}
`

var testCreds2 = `
orchestrator nomad {
    token = "foo"
}

secrets vault {
    token = ""
}

vcs github {
    foo = "bar"
}
`

func Test_CredsConfig_Merge(t *testing.T) {
	var c1, c2 CredsConfig
	err := hcl.Decode(&c1, testCreds1)
	assert.Nil(t, err)
	err = hcl.Decode(&c2, testCreds2)
	assert.Nil(t, err)

	c1.Merge(&c2)
	fmt.Printf("%+v\n", c1)
	assert.Equal(t, "bar", c1.VCS["github"]["foo"])
	assert.Equal(t, "token", c1.VCS["github"]["token"])
}

func Test_CredsConfig(t *testing.T) {

	var cconf CredsConfig
	err := hcl.Decode(&cconf, testCreds1)
	assert.Nil(t, err)
	assert.Equal(t, "token", cconf.Registry["ecr"]["token"])
	assert.Equal(t, "key", cconf.Registry["ecr"]["key"])
	assert.Equal(t, "secret", cconf.Registry["ecr"]["secret"])
	assert.Equal(t, "", cconf.Secrets["vault"]["token"])
	assert.Equal(t, "token", cconf.VCS["github"]["token"])
	t.Logf("%+v", cconf)
}

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

func Test_ThrapConfig(t *testing.T) {

	var conf ThrapConfig
	err := hcl.Decode(&conf, testConfig1)
	assert.Nil(t, err)
	assert.Equal(t, "thrap", conf.VCS["github"].Repo.Name)
	assert.Equal(t, "http://vault:4646", conf.Secrets["vault"].Addr)
}

func Test_Config_Merge(t *testing.T) {
	var c1, c2 ThrapConfig

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

func Test_ThrapConfig_Encode(t *testing.T) {
	conf := DefaultThrapConfig()
	b, _ := hclencoder.Encode(conf)
	t.Logf("%s", b)
}
