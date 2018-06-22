package builder

import (
	"fmt"
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

// Builder implements a code builder ie. test and build/compile/assemble etc
type Builder interface {
	Build(*thrapb.Component, io.Writer) error
}

// New returns a new builder given the config
func New(conf interface{}) (Builder, error) {
	return newDockerImageBuilder()
}

type dockerImageBuilder struct {
	cli *client.Client
}

func newDockerImageBuilder() (*dockerImageBuilder, error) {
	os.Setenv("DOCKER_API_VERSION", "1.37")
	cli, err := client.NewEnvClient()
	if err == nil {
		return &dockerImageBuilder{cli}, nil
	}
	return nil, err
}

func (builder *dockerImageBuilder) Build(comp *thrapb.Component, output io.Writer) error {
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

	fmt.Println(comp.ID)
	opt := types.ImageBuildOptions{
		Tags:       []string{},
		BuildID:    comp.ID,
		Dockerfile: bc.Dockerfile,
	}

	ctx := context.Background()
	resp, err := builder.cli.ImageBuild(ctx, rdc, opt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, output, 100, true, nil)
	return err
}
