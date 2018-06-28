package core

import (
	"context"
	"testing"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/utils"
	"github.com/stretchr/testify/assert"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_core(t *testing.T) {
	opt := DefaultConfigureOptions()
	opt.NoPrompt = true
	err := ConfigureGlobal(opt)
	fatal(t, err)

	conf := &CoreConfig{PacksDir: "../etc/packs"}
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
