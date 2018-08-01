package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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
	assert.NotNil(t, c.orchs)
	assert.NotNil(t, c.packs)
	assert.NotNil(t, c.orchs["docker"])

	_, err = c.Stack(&Profile{Orchestrator: "foo"})
	assert.Contains(t, err.Error(), errOrchNotLoaded.Error())
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
		Orchestrator: map[string]*config.OrchestratorConfig{
			"docker": &config.OrchestratorConfig{},
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

	st, err := c.Stack(DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}
	err = st.Build(context.Background(), stack, BuildOptions{Workdir: "../"})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Core_populateFromImageConf(t *testing.T) {

	if !utils.FileExists("/var/run/docker.sock") {
		t.Skip("Skipping: docker file descriptor not found")
	}

	tmpdir, _ := ioutil.TempDir("/tmp", "core-")
	defer os.RemoveAll(tmpdir)

	opt := DefaultConfigureOptions()
	opt.NoPrompt = true
	opt.DataDir = tmpdir
	err := ConfigureGlobal(opt)
	if err != nil {
		t.Fatal(err)
	}

	conf := &Config{DataDir: tmpdir, ThrapConfig: &config.ThrapConfig{
		Registry: map[string]*config.RegistryConfig{
			"ecr": &config.RegistryConfig{
				ID:   "ecr",
				Addr: "foobar.com",
			},
		},
		Orchestrator: map[string]*config.OrchestratorConfig{
			"docker": &config.OrchestratorConfig{},
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

	stack.Validate()

	st, err := c.Stack(DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}

	st.populateFromImageConf(stack)
	assert.Equal(t, 1, len(stack.Components["vault"].Ports))
	assert.Equal(t, 5, len(stack.Components["consul"].Ports))
	assert.True(t, stack.Components["consul"].HasVolumeTarget("/consul/data"))
}

func Test_Core_Assembler(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("/tmp", "core.stack-")
	defer os.RemoveAll(tmpdir)

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

	st, err := c.Stack(DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}
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
