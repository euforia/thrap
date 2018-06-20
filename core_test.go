package thrap

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
	c, err := NewCore(nil)
	fatal(t, err)

	assert.NotNil(t, c.vcs)
	assert.Nil(t, c.reg)
	assert.NotNil(t, c.sec)
}
