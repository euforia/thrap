package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/utils"
	"github.com/stretchr/testify/assert"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ConfigureGlobal(t *testing.T) {
	opt := DefaultConfigureOptions()
	opt.NoPrompt = true
	opt.DataDir, _ = ioutil.TempDir("/tmp", "cg-")
	err := ConfigureGlobal(opt)
	fatal(t, err)

	cf := filepath.Join(opt.DataDir, consts.WorkDir, consts.CredsFile)
	assert.True(t, utils.FileExists(cf))
	cf = filepath.Join(opt.DataDir, consts.WorkDir, consts.ConfigFile)
	assert.True(t, utils.FileExists(cf))
}

func Test_core(t *testing.T) {
	opt := DefaultConfigureOptions()
	opt.NoPrompt = true
	err := ConfigureGlobal(opt)
	fatal(t, err)

	conf := &Config{PacksDir: "../etc/packs"}
	c, err := NewCore(conf)
	fatal(t, err)

	assert.NotNil(t, c.regs)
	assert.NotNil(t, c.sec)
	assert.NotNil(t, c.vcs)
	assert.NotNil(t, c.orch)
	assert.NotNil(t, c.packs)
	assert.Equal(t, "nomad", c.orch.ID())

	stack, err := manifest.LoadManifest("../test-fixtures/thrap.hcl")
	fatal(t, err)
	errs := stack.Validate()
	if len(errs) > 0 {
		fatal(t, utils.FlattenErrors(errs))
	}

	st := c.Stack()
	st.populatePorts(stack)
	assert.Equal(t, 1, len(stack.Components["vault"].Ports))
	assert.Equal(t, 5, len(stack.Components["consul"].Ports))

	err = st.Build(context.Background(), stack)
	if err != nil {
		t.Fatal(err)
	}

}

func Test_core_stack(t *testing.T) {

	lconf, _ := config.ReadProjectConfig("../")
	conf := &Config{
		PacksDir:    "~/.thrap/packs",
		ThrapConfig: lconf,
	}

	c, err := NewCore(conf)
	fatal(t, err)

	stack, err := manifest.LoadManifest("../thrap.hcl")
	fatal(t, err)
	errs := stack.Validate()
	if len(errs) > 0 {
		fatal(t, utils.FlattenErrors(errs))
	}

	st := c.Stack()
	sasm, err := st.Assembler(stack)
	fatal(t, err)

	err = sasm.Assemble()
	if err != nil {
		t.Fatal(err)
	}

	casm := sasm.ComponentAsm("registry")
	fmt.Println(casm.Dockerfile())

	fmt.Printf("%+v\n", stack)
}
