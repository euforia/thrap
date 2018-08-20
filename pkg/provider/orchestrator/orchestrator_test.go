package orchestrator

import (
	"testing"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func Test_Orchestrator(t *testing.T) {
	conf := &provider.Config{Provider: "nomad"}
	orch, err := New(conf)
	assert.Nil(t, err)
	o := orch.(*nomadOrchestrator)
	assert.NotNil(t, o.client)

	conf.Provider = "foo"
	_, err = New(conf)
	assert.Contains(t, err.Error(), "unsupported")
}
