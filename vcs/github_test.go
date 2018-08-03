package vcs

import (
	"os"
	"testing"

	"github.com/euforia/thrap/config"
	"github.com/stretchr/testify/assert"
)

func Test_Github(t *testing.T) {

	g := newGithubVCS(nil)
	err := g.Init(map[string]interface{}{"token": "foo"})
	assert.Nil(t, err)

	assert.Equal(t, g.IgnoresFile(), gitIgnoresFile)
}

func Test_Github_calls(t *testing.T) {
	conf, err := config.ReadCredsConfig("../.thrap/creds.hcl")
	if err != nil {
		t.Skipf("Skipping github calls: %v", err)
	}

	tdir := "/tmp/github-test"
	defer os.RemoveAll(tdir)

	c := conf.GetVCSCreds("github")
	if c == nil {
		t.Skip("Skipping github calls: no config")
	}
	if _, ok := c["token"]; !ok {
		t.Skip("Skipping github calls: no token")
	}

	g := newGithubVCS(nil)
	g.Init(map[string]interface{}{"token": c["token"]})

	var opt Option
	repo := &Repository{Name: "api-test"}
	_, _, err = g.Create(repo, opt)
	assert.Nil(t, err)

	repo.Owner = "euforia"
	_, err = g.Get(repo, opt)
	assert.Nil(t, err)

	repo.Owner = "euforia"
	err = g.Delete(repo, opt)
	assert.Nil(t, err)
}
