package orchestrator

import (
	"errors"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/thrapb"
	"golang.org/x/net/context"
)

type DockerOrchestrator struct {
	cli *client.Client
}

func (orch *DockerOrchestrator) Init(config map[string]interface{}) error {
	os.Setenv("DOCKER_API_VERSION", "1.37")
	cli, err := client.NewEnvClient()
	if err == nil {
		orch.cli = cli
	}
	return err
}

func (orch *DockerOrchestrator) Deploy(stack *thrapb.Stack, opts DeployOptions) (interface{}, interface{}, error) {
	return nil, nil, errors.New("tbi")
}

func (orch *DockerOrchestrator) Build(comp *thrapb.Component, output io.Writer) error {
	bc := comp.Build

	tarOpt := &archive.TarOptions{}

	ign, err := dockerfile.ReadIgnoresFile(bc.Context)
	if err == nil {
		tarOpt.ExcludePatterns = ign
	}

	rdc, err := archive.TarWithOptions(bc.Context, tarOpt)
	if err != nil {
		return err
	}
	defer rdc.Close()

	opt := types.ImageBuildOptions{
		Tags:       []string{},
		BuildID:    comp.ID,
		Dockerfile: bc.Dockerfile,
		// BuildArgs:  make(map[string]*string),
	}

	ctx := context.Background()
	resp, err := orch.cli.ImageBuild(ctx, rdc, opt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, output, 100, true, nil)
	return err
}
