package config

import (
	"testing"

	"github.com/euforia/hclencoder"
	"github.com/hashicorp/hcl"
	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	c, err := ParseFile("../etc/config.hcl")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", c)
}

func Test_CredsConfig(t *testing.T) {
	testIn := `
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
    token = ""
}
`
	var cconf CredsConfig
	err := hcl.Decode(&cconf, testIn)
	assert.Nil(t, err)
	assert.Equal(t, "token", cconf.Registry["ecr"]["token"])
	assert.Equal(t, "key", cconf.Registry["ecr"]["key"])
	assert.Equal(t, "secret", cconf.Registry["ecr"]["secret"])
	assert.Equal(t, "", cconf.Secrets["vault"]["token"])
	assert.Equal(t, "", cconf.VCS["github"]["token"])
	t.Logf("%+v", cconf)
}

func Test_ThrapConfig(t *testing.T) {
	testIn := `

registry {
    id = "ecr"
}

orchestrator {
    id = "docker"
    addr = ""
}

secrets {
    id = "vault"
    addr = "http://vault.service:4646"
}

vcs {
    id = "github"
    addr = "github.com"
    username = "euforia"
    repo {
        name = "thrap"
        owner = "euforia"
    }
}
`

	var conf ThrapConfig
	err := hcl.Decode(&conf, testIn)
	assert.Nil(t, err)
	assert.Equal(t, "thrap", conf.VCS.Repo.Name)
	assert.Equal(t, "vault", conf.Secrets.ID)
}

func Test_ThrapConfig_Encode(t *testing.T) {
	conf := DefaultThrapConfig()
	b, _ := hclencoder.Encode(conf)
	t.Logf("%s", b)
}
