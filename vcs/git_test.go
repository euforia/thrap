package vcs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func committer() *git.CommitOptions {
	return &git.CommitOptions{
		Author: &object.Signature{
			Name:  "thrap",
			Email: "thrap",
			When:  time.Now(),
		},
	}
}

func Test_GetRepoVersion(t *testing.T) {

	tmpdir, _ := ioutil.TempDir("/tmp", "grv-")
	defer os.RemoveAll(tmpdir)

	ver := GetRepoVersion(tmpdir)
	assert.Equal(t, "v0.0.0", ver.String())

	_, repo, _ := SetupLocalGitRepo("test", "me", tmpdir, "foo.com")
	ver, _ = getRepoVersion(repo)
	assert.Equal(t, "v0.0.0", ver.String())

	wt, _ := repo.Worktree()
	wt.AddGlob("dkfjdkfj")

	wt.Commit("Initial commit", committer())

	ver, _ = getRepoVersion(repo)
	assert.Equal(t, 1, ver.Count)

	for i, s := range []string{"foo", "bar", "baz", "bees", "wx"} {
		wt.AddGlob(s)
		ch, _ := wt.Commit("message", committer())

		if (i%3) == 0 && i != 0 {
			// tag
			refName := plumbing.ReferenceName(fmt.Sprintf("refs/tags/tag%d", i))
			ref := plumbing.NewHashReference(refName, ch)
			repo.Storer.SetReference(ref)

			ver, _ = getRepoVersion(repo)
			assert.Zero(t, ver.Count)
		}

	}

	ver, _ = getRepoVersion(repo)
	assert.Equal(t, "tag3", ver.Tag)

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

	_, _, err = v.Create(&Repository{Name: "fail"}, Option{})
	assert.NotNil(t, err)

	defer os.RemoveAll("/tmp/remote-test")
	_, _, err = v.Create(&Repository{Name: "fail"}, Option{Path: "/tmp/remote-test", Remote: "git://foo"})
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
	r, _, err := vcs.Create(cr, opt)
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
	_, _, err = vcs.Create(&Repository{Name: testdir[1:]}, opt)
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
	_, _, err := vcs.Create(&Repository{Name: testdir[1:]}, opt)
	assert.Nil(t, err)

	// Open
	opt.Remote = "ssh://git@github.com" + testdir
	_, _, err = vcs.Create(&Repository{Name: testdir[1:]}, opt)
	assert.Nil(t, err)

	// Parse error
	opt.Remote = "git@github.com:" + testdir[1:]
	_, _, err = vcs.Create(&Repository{Name: testdir[1:]}, opt)
	assert.NotNil(t, err)

}
