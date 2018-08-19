package registry

import (
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/crt"
)

var (
	// ErrAuthNotRequired is used for registries that do not require
	// auth to signal the caller.
	ErrAuthNotRequired = errors.New("registry auth not required")
)

// localDocker implements the registry interface for a local docker runtime.
// This is mostly no-ops to satisfy the Registry interface for consistency
type localDocker struct {
	id  string
	crt *crt.Docker
}

func (reg *localDocker) ID() string {
	return reg.id
}

// Initialize the registry provider
func (reg *localDocker) Init(conf *config.RegistryConfig) (err error) {
	reg.id = conf.ID
	reg.crt, err = crt.NewDocker()
	return
}

func (reg *localDocker) GetAuthConfig() (types.AuthConfig, error) {
	var auth types.AuthConfig
	return auth, ErrAuthNotRequired
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
	return reg.crt.ImageConfig(name + ":" + tag)
}

// Name of the image with the registry. This is needed for deployments. In the
// case of local docker it is the image name itself.
func (reg *localDocker) ImageName(name string) string {
	return name
}
