package registry

import (
	"os"
	"testing"

	"github.com/euforia/thrap/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func Test_NewContainerRegistry(t *testing.T) {
	conf := &Config{TypeContainer, "ecr", nil}
	reg, err := New(conf)
	assert.Nil(t, err)
	assert.Equal(t, TypeContainer, reg.Type())

	conf.Conf = map[string]interface{}{"region": "us-west-2"}
	reg, err = New(conf)
	assert.Nil(t, err)
	r := reg.(*awsContainerRegistry)
	assert.Equal(t, "us-west-2", *r.sess.Config.Region)

	conf.Provider = "unsupported"
	_, err = New(conf)
	assert.Contains(t, err.Error(), "unsupported")
}

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ECR(t *testing.T) {
	credsfile, _ := homedir.Expand("~/.thrap/creds.hcl")
	cc, err := config.ReadCredsConfig(credsfile)
	fatal(t, err)

	ecrCreds := cc.GetRegistryCreds("ecr")

	conf := &Config{
		Type:     TypeContainer,
		Provider: "ecr",
		Conf: map[string]interface{}{
			"region": "us-west-2",
		},
	}

	for k, v := range ecrCreds {
		conf.Conf[k] = v
	}

	treg, _ := New(conf)
	reg := treg.(*awsContainerRegistry)
	err = reg.Init(*conf)
	assert.Nil(t, err)
	assert.NotEmpty(t, os.Getenv("AWS_ACCESS_KEY_ID"))
	assert.NotEmpty(t, os.Getenv("AWS_SECRET_ACCESS_KEY"))

	_, err = reg.Get("keycloak")
	assert.Nil(t, err)

	repoName := "test-thrap/test-comp"
	_, err = reg.Create(repoName)
	assert.Nil(t, err)
	_, err = reg.Get(repoName)
	assert.Nil(t, err)
	_, err = reg.Delete(repoName)
	assert.Nil(t, err)

	_, err = reg.Get(repoName)
	assert.NotNil(t, err)
	// nr := nrepo.(*ecr.Repository)
}
