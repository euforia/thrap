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
	if err != nil {
		t.Fatal(err)
	}

	cf := filepath.Join(opt.DataDir, consts.CredsFile)
	assert.True(t, utils.FileExists(cf))
	cf = filepath.Join(opt.DataDir, consts.ConfigFile)
	assert.True(t, utils.FileExists(cf))
}

func Test_NewCore(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("/tmp", "core-")

	opt := DefaultConfigureOptions()
	opt.NoPrompt = true
	opt.DataDir = tmpdir
	err := ConfigureGlobal(opt)
	if err != nil {
		t.Fatal(err)
	}

	conf := &Config{DataDir: tmpdir}
	c, err := NewCore(conf)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, c.regs)
	assert.NotNil(t, c.sec)
	assert.NotNil(t, c.vcs)
	assert.NotNil(t, c.orch)
	assert.NotNil(t, c.packs)
	assert.Equal(t, "nomad", c.orch.ID())
}

func Test_Core_Build(t *testing.T) {
	if !utils.FileExists("/var/run/docker.sock") {
		t.Skip("Skipping: docker file descriptor not found")
	}

	tmpdir, _ := ioutil.TempDir("/tmp", "core-")

	opt := DefaultConfigureOptions()
	opt.NoPrompt = true
	opt.DataDir = tmpdir
	err := ConfigureGlobal(opt)
	if err != nil {
		t.Fatal(err)
	}

	// tmpdir = filepath.Join(tmpdir, "packs")
	conf := &Config{DataDir: tmpdir, ThrapConfig: &config.ThrapConfig{
		Registry: map[string]*config.RegistryConfig{
			"ecr": &config.RegistryConfig{
				ID:   "ecr",
				Addr: "foobar.com",
			},
		},
	}}
	c, err := NewCore(conf)
	if err != nil {
		t.Fatal(err)
	}

	stack, err := manifest.LoadManifest("../test-fixtures/thrap.hcl")
	if err != nil {
		t.Fatal(err)
	}

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

func Test_Core_Assembler(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("/tmp", "core.stack-")
	opt := DefaultConfigureOptions()
	opt.NoPrompt = true
	opt.DataDir = tmpdir
	err := ConfigureGlobal(opt)
	fatal(t, err)

	lconf, err := config.ReadProjectConfig("../")
	fatal(t, err)
	conf := &Config{
		DataDir:     tmpdir,
		ThrapConfig: lconf,
	}

	c, err := NewCore(conf)
	if err != nil {
		t.Fatal(err)
	}

	stack, _ := manifest.LoadManifest("../thrap.hcl")
	stack.Validate()

	st := c.Stack()
	sasm, err := st.Assembler("../", stack)
	if err != nil {
		t.Fatal(err)
	}

	err = sasm.Assemble()
	if err != nil {
		t.Fatal(err)
	}

	casm := sasm.ComponentAsm("registry")
	assert.NotNil(t, casm)
	fmt.Println(sasm.ComponentAsm("nomad").Dockerfile())
}
