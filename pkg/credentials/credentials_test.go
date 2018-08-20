package credentials

import (
	"fmt"
	"testing"

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
	var c1, c2 Credentials
	err := hcl.Decode(&c1, testCreds1)
	assert.Nil(t, err)
	err = hcl.Decode(&c2, testCreds2)
	assert.Nil(t, err)

	c1.Merge(&c2)
	fmt.Printf("%+v\n", c1)
	assert.Equal(t, "bar", c1.VCS["github"]["foo"])
	assert.Equal(t, "token", c1.VCS["github"]["token"])
}

func Test_Credentials(t *testing.T) {

	var cconf Credentials
	err := hcl.Decode(&cconf, testCreds1)
	assert.Nil(t, err)
	assert.Equal(t, "token", cconf.Registry["ecr"]["token"])
	assert.Equal(t, "key", cconf.RegistryCreds("ecr")["key"])
	assert.Equal(t, "secret", cconf.RegistryCreds("ecr")["secret"])
	assert.Equal(t, "", cconf.Secrets["vault"]["token"])
	assert.Equal(t, "token", cconf.VCS["github"]["token"])
	t.Logf("%+v", cconf)
}
