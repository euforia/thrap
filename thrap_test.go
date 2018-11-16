package thrap

import (
	"fmt"
	"testing"

	"github.com/euforia/thrap/pkg/provider/vcs"
	"github.com/euforia/thrap/utils"
	"gopkg.in/src-d/go-git.v4"
)

func Test_foo(t *testing.T) {
	vcsp := vcs.NewGitVCS()
	vcsp.Init(nil)
	cwd, _ := utils.GetLocalPath("")
	r, _ := vcsp.Open(&vcs.Repository{Name: "thrap"}, vcs.Option{Path: cwd})
	repo := r.(*git.Repository)
	wt, _ := repo.Worktree()
	fmt.Println(wt.Status())
}
