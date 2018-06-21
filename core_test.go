package thrap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_core(t *testing.T) {
	c, err := NewCore(nil)
	fatal(t, err)

	assert.NotNil(t, c.vcs)
	assert.NotNil(t, c.reg)
	assert.NotNil(t, c.sec)
}

func Test_foo(t *testing.T) {

	repo, _ := git.PlainOpen(".")

	tags, _ := repo.Tags()
	// last,tags.Next()
	tags.ForEach(func(arg1 *plumbing.Reference) error {
		fmt.Println(arg1)
		return nil
	})

	cmtIter, _ := repo.CommitObjects()
	cmtIter.ForEach(func(arg1 *object.Commit) error {
		fmt.Println(arg1)
		return nil
	})
}
