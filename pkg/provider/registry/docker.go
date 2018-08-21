package registry

import (
	"context"
	"errors"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/euforia/thrap/pkg/provider"
)

var (
	// ErrAuthNotRequired is used for registries that do not require
	// auth to signal the caller.
	ErrAuthNotRequired = errors.New("registry auth not required")
)

// localDocker implements the registry interface for a local docker runtime.
// This is mostly no-ops to satisfy the Registry interface for consistency
type localDocker struct {
	id string
	// crt *crt.Docker
	cli *client.Client
}

func (reg *localDocker) ID() string {
	return reg.id
}

// Initialize the registry provider
func (reg *localDocker) Init(conf *provider.Config) (err error) {
	reg.id = conf.ID

	// Set default
	os.Setenv("DOCKER_API_VERSION", "1.37")

	// Set from config if available
	if val, ok := conf.Config["DOCKER_API_VERSION"]; ok {
		if apiVersion, ok := val.(string); ok {
			os.Setenv("DOCKER_API_VERSION", apiVersion)
		}
	}

	reg.cli, err = client.NewEnvClient()

	return
}

func (reg *localDocker) GetAuthConfig() (types.AuthConfig, error) {
	var auth types.AuthConfig
	return auth, ErrAuthNotRequired
}

func (reg *localDocker) Name() string {
	return ""
}

// Create is a no-op that satisfies the Register interface, as there's
// nothing to do for local docker
func (reg *localDocker) CreateRepo(string) (interface{}, error) {
	return nil, nil
}

// Get is a no-op that satisfies the Register interface
func (reg *localDocker) GetRepo(string) (interface{}, error) {
	return nil, nil
}

// Get image manifest
func (reg *localDocker) GetImageManifest(name, tag string) (interface{}, error) {
	image := name + ":" + tag

	inf, _, err := reg.cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		return nil, err
	}

	return inf.Config, nil
}

// Name of the image with the registry. This is needed for deployments. In the
// case of local docker it is the image name itself.
func (reg *localDocker) ImageName(name string) string {
	return name
}
