package registry

// Local docker engine as the registry

import (
	"context"
	"errors"
	"os"

	"github.com/docker/docker/client"
)

const defaultDockerRegAddr = "https://registry.hub.docker.com"

// localDocker implements a container registry using the local docker
// registry
type localDocker struct {
	cli *client.Client
}

func (local *localDocker) Init(config Config) error {
	os.Setenv("DOCKER_API_VERSION", "1.37")

	cli, err := client.NewEnvClient()
	if err == nil {
		local.cli = cli
	}

	return err
}

func (local *localDocker) Type() Type {
	return TypeContainer
}

func (local *localDocker) ID() string {
	return "local.docker"
}

// Create a new repository
func (local *localDocker) Create(string) (interface{}, error) {
	return nil, errors.New("tbi")
}

// Get a repository manifest
func (local *localDocker) GetManifest(name string, tag string) (interface{}, error) {
	inf, _, err := local.cli.ImageInspectWithRaw(context.Background(), name+":"+tag)
	return inf, err
}
