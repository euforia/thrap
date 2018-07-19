package vcs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
)

func Test_GetRepoVersion(t *testing.T) {
	t.Log(GetRepoVersion("../"))
}

func Test_VCS(t *testing.T) {
	conf := &Config{Provider: "xxx"}
	_, err := New(conf)
	assert.NotNil(t, err)

	var v VCS

	p := []string{"github", "git"}
	for i := range p {
		conf.Provider = p[i]
		v, err = New(conf)
		assert.Nil(t, err)
	}

	_, err = v.Create(&Repository{Name: "fail"}, Option{})
	assert.NotNil(t, err)

	defer os.RemoveAll("/tmp/remote-test")
	_, err = v.Create(&Repository{Name: "fail"}, Option{Path: "/tmp/remote-test", Remote: "git://foo"})
	assert.Contains(t, err.Error(), "unsupported URL scheme")

	_, ok := v.(*GitVCS)
	assert.True(t, ok)
}

func Test_git(t *testing.T) {
	testdir := "/tmp/test-git"
	defer os.RemoveAll(testdir)

	vcs := NewGitVCS()

	// Create
	rmt := "ssh://git@github.com" + testdir
	opt := Option{Path: testdir, Remote: rmt}
	cr := &Repository{Name: testdir[1:]}
	r, err := vcs.Create(cr, opt)
	assert.Nil(t, err)

	repo := r.(*git.Repository)

	// Check remotes
	c, _ := repo.Config()
	origin, _ := c.Remotes["origin"]
	assert.Equal(t, rmt, origin.URLs[0])

	// Check directory
	stat, err := os.Stat(testdir + "/.git")
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())

	// Check exists error
	_, err = vcs.Create(&Repository{Name: testdir[1:]}, opt)
	// assert.Equal(t, git.ErrRemoteExists, err)
	assert.Nil(t, err)
}

func Test_git_init(t *testing.T) {
	vcs := NewGitVCS()
	err := vcs.Init(nil)
	assert.Nil(t, err)
	// assert.NotEmpty(t, vcs.globalUser)
	// t.Log(vcs.globalUser)
}

func Test_parseURL(t *testing.T) {
	_, err := parseURL("git@github.com:foo/bar")
	assert.NotNil(t, err)
}

func Test_git_remote(t *testing.T) {
	testdir := "/tmp/test-git0"
	defer os.RemoveAll(testdir)

	conf := &Config{Provider: "git"}
	vcs, _ := New(conf)

	// No remote
	opt := Option{Path: testdir}
	_, err := vcs.Create(&Repository{Name: testdir[1:]}, opt)
	assert.Nil(t, err)

	// Open
	opt.Remote = "ssh://git@github.com" + testdir
	_, err = vcs.Create(&Repository{Name: testdir[1:]}, opt)
	assert.Nil(t, err)

	// Parse error
	opt.Remote = "git@github.com:" + testdir[1:]
	_, err = vcs.Create(&Repository{Name: testdir[1:]}, opt)
	assert.NotNil(t, err)

}
