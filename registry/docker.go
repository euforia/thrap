package registry

// Local docker engine as the registry

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

// DockerRuntime implements a container registry using the local docker runtime
type DockerRuntime struct {
	cli *client.Client
}

func (local *DockerRuntime) Init(config Config) error {
	os.Setenv("DOCKER_API_VERSION", "1.37")

	cli, err := client.NewEnvClient()
	if err == nil {
		local.cli = cli
	}

	return err
}

// ImagePull pulls in image from the docker registry using docker. This uses
// dockers built in mechanism to communicate to the registry
func (local *DockerRuntime) ImagePull(ctx context.Context, ref string) error {
	var options types.ImagePullOptions
	rd, err := local.cli.ImagePull(ctx, ref, options)
	if err != nil {
		return err
	}

	defer rd.Close()

	return jsonmessage.DisplayJSONMessagesStream(rd, os.Stdout, 100, true, nil)
}

// ImageConfig returns an image config for the given name and tagged image
func (local *DockerRuntime) ImageConfig(name, tag string) (*container.Config, error) {
	inf, _, err := local.cli.ImageInspectWithRaw(context.Background(), name+":"+tag)
	if err != nil {
		return nil, err
	}
	return inf.Config, nil
}

// func (local *DockerRuntime) Type() Type {
// 	return TypeContainer
// }

// func (local *DockerRuntime) ID() string {
// 	return "local.docker"
// }

// Create a new repository
// func (local *DockerRuntime) Create(string) (interface{}, error) {
// 	return nil, errors.New("tbi")
// }

// Get a repository manifest
// func (local *DockerRuntime) GetManifest(name string, tag string) (interface{}, error) {
// 	inf, _, err := local.cli.ImageInspectWithRaw(context.Background(), name+":"+tag)
// 	return inf, err
// }
