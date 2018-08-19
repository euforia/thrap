package registry

import (
	"testing"

	"github.com/euforia/thrap/config"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	conf := &Config{
		ID:       "ecr",
		Provider: "ecr",
	}
	reg, err := New(conf)
	assert.Nil(t, err)

	conf.Config = map[string]interface{}{"region": "us-west-2"}
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

func Test_Registry_ECR(t *testing.T) {
	credsfile := "../../../.thrap/creds.hcl"
	cc, err := config.ReadCredentials(credsfile)
	fatal(t, err)

	// Add creds
	ecrCreds := cc.RegistryCreds("sandbox")
	conf := &Config{
		Provider: "ecr",
		Config: map[string]interface{}{
			"region": "us-west-2",
		},
	}
	for k, v := range ecrCreds {
		conf.Config[k] = v
	}

	treg, _ := New(conf)
	reg := treg.(*awsContainerRegistry)
	err = reg.Init(conf)

	assert.Nil(t, err)
	// assert.NotEmpty(t, os.Getenv("AWS_ACCESS_KEY_ID"))
	// assert.NotEmpty(t, os.Getenv("AWS_SECRET_ACCESS_KEY"))

	_, err = reg.GetRepo("keycloak")
	assert.Nil(t, err)

	repoName := "test-thrap/test-comp"
	_, err = reg.CreateRepo(repoName)
	assert.Nil(t, err)

	_, err = reg.GetImageManifest("test-comp", "notfound")
	assert.NotNil(t, err)

	_, err = reg.GetRepo(repoName)
	assert.Nil(t, err)
	_, err = reg.DeleteRepo(repoName)
	assert.Nil(t, err)

	_, err = reg.GetRepo(repoName)
	assert.NotNil(t, err)
	// nr := nrepo.(*ecr.Repository)
}
