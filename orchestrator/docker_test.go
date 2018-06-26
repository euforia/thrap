package orchestrator

import (
	"os"
	"testing"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func Test_DockerOrchestrator_Build(t *testing.T) {
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

	err = bldr.Build("default", bcomp, RequestOptions{Output: os.Stdout})
	assert.Nil(t, err)

}

func Test_DockerOrchestrator_Deploy(t *testing.T) {
	st := &thrapb.Stack{ID: "test",
		Components: map[string]*thrapb.Component{
			"vault": &thrapb.Component{
				ID:   "vault",
				Name: "vault", Version: "0.10.3",
			},
			"consul": &thrapb.Component{
				ID:   "consul",
				Name: "consul", Version: "1.1.0",
			},
			"api": &thrapb.Component{
				ID:      "api",
				Name:    "api",
				Version: "latest",
				Head:    true,
				Build: &thrapb.Build{
					Dockerfile: "test.dockerfile",
					Context:    "../test-fixtures",
				},
			},
		},
	}

	bldr := &DockerOrchestrator{}
	bldr.Init(nil)

	defer bldr.tearDown(context.Background(), st)

	_, _, err := bldr.Deploy(st, RequestOptions{Output: os.Stdout})
	if err != nil {
		t.Fatal(err)
	}

	// <-time.After(2 * time.Second)

}
