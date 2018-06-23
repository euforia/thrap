package orchestrator

import (
	"os"
	"testing"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	"github.com/stretchr/testify/assert"
)

func Test_dockerImageBuilder(t *testing.T) {
	bldr := &DockerOrchestrator{}
	err := bldr.Init(nil)
	if err != nil {
		t.Fatal(err)
	}

	st, err := manifest.LoadManifest("../test-fixtures/builder.yml")
	if err != nil {
		t.Fatal(err)
	}
	st.Validate()

	var bcomp *thrapb.Component
	for _, comp := range st.Components {
		if comp.IsBuildable() {
			bcomp = comp
			bcomp.Build.Context = "../test-fixtures"
			break
		}
	}

	err = bldr.Build(bcomp, os.Stdout)
	assert.Nil(t, err)

}
