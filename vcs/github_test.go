package vcs

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Github(t *testing.T) {
	tdir := "/tmp/github-test"
	defer os.RemoveAll(tdir)

	b, err := ioutil.ReadFile("../secrets.txt")
	if err != nil {
		t.Fatal(err)
	}
	token := strings.TrimSuffix(string(b), "\n")

	g := newGithubVCS(nil)
	g.Init(map[string]interface{}{"token": token})

	assert.Equal(t, g.IgnoresFile(), gitIgnoresFile)

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
