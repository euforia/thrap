package registry

import (
	"errors"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/crt"
)

var (
	errNotImplemented = errors.New("not implemented")
)

// localDocker implements the registry interface for a local docker runtime
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

// Create a new repository
func (reg *localDocker) Create(string) (interface{}, error) {
	return nil, errNotImplemented
}

// Get repo info
func (reg *localDocker) Get(string) (interface{}, error) {
	return nil, errNotImplemented
}

// Get image manifest
func (reg *localDocker) GetManifest(name, tag string) (interface{}, error) {
	return reg.crt.ImageConfig(name + ":" + tag)
}

// Name of the image with the registry. Needed for deployments
func (reg *localDocker) ImageName(name string) string {
	return name
}
