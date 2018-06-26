package core

import (
	"testing"

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
	assert.Equal(t, "docker", c.orch.ID())
}
