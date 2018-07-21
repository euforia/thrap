package crt

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/utils"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func Test_Docker_Build(t *testing.T) {
	if !utils.FileExists("/var/run/docker.sock") {
		t.Skip("Skipping: docker file descriptor not found")
	}

	dkr, _ := NewDocker()

	lr := NewDockerBuildLog(os.Stdout)
	req := &BuildRequest{
		Output:     lr,
		ContextDir: "../test-fixtures",
		BuildOpts: &types.ImageBuildOptions{
			Dockerfile: "multi-stage.Dockerfile",
		},
	}

	err := dkr.Build(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	logSteps, err := lr.Steps()
	if err != nil {
		t.Fatal(err)
	}

	raw, _ := dockerfile.ParseFile("../test-fixtures/multi-stage.Dockerfile")
	df := dockerfile.ParseRaw(raw)
	assert.Equal(t, df.StepCount(), len(logSteps))

	for i, step := range logSteps {
		_, err := dkr.ImageConfig(step.ID())
		assert.Nil(t, err)
		if i < len(logSteps)-1 {
			assert.Equal(t, "", step.Log())
		}

		fmt.Println(step.ID(), step.Cmd())

	}
}
