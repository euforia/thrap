package vcs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/euforia/thrap/utils"
	"github.com/hashicorp/hcl"
	"github.com/stretchr/testify/assert"
)

func Test_Github(t *testing.T) {

	g := newGithubVCS(nil)
	err := g.Init(map[string]interface{}{"token": "foo"})
	assert.Nil(t, err)

	assert.Equal(t, g.IgnoresFile(), gitIgnoresFile)
}

func Test_Github_calls(t *testing.T) {
	if !utils.FileExists("../secrets.hcl") {
		t.Skip("Skipping github calls: no creds")
	}

	tdir := "/tmp/github-test"
	defer os.RemoveAll(tdir)

	b, err := ioutil.ReadFile("../secrets.hcl")
	if err != nil {
		t.Fatal(err)
	}

	var c map[string]interface{}
	err = hcl.Unmarshal(b, &c)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := c["token"]; !ok {
		t.Skip("Skipping github calls: no token")
	}
	// token := strings.TrimSuffix(string(b), "\n")

	g := newGithubVCS(nil)
	g.Init(c)

	var opt Option
	repo := &Repository{Name: "api-test"}
	_, err = g.Create(repo, opt)
	assert.Nil(t, err)

	repo.Owner = "euforia"
	_, err = g.Get(repo, opt)
	assert.Nil(t, err)

	repo.Owner = "euforia"
	err = g.Delete(repo, opt)
	assert.Nil(t, err)
}
