package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/euforia/thrap/manifest"
	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

func Test_nomad_env(t *testing.T) {
	conf := &Config{Provider: "nomad"}
	orch, err := New(conf)
	assert.Nil(t, err)

	norch := orch.(*nomadOrchestrator)
	assert.Equal(t, norch.client.Address(), os.Getenv("NOMAD_ADDR"))
}

func Test_nomad_dryrun(t *testing.T) {
	conf := &Config{Provider: "nomad", Conf: map[string]interface{}{
		"addr": os.Getenv("NOMAD_ADDR"),
	}}
	orch, err := New(conf)
	assert.Nil(t, err)

	st, _ := manifest.LoadManifest("../thrap.hcl")
	st.Validate()

	_, ijob, err := orch.Deploy(st, RequestOptions{Dryrun: true})
	assert.Nil(t, err)

	ijob.(*api.Job).Canonicalize()
	b, _ := json.MarshalIndent(ijob, "", "  ")
	fmt.Printf("%s\n", b)
}
