package registry

import (
	"fmt"

	"github.com/euforia/thrap/config"
)

// // Type is the registry type
// type Type uint8

// func (t Type) String() string {
// 	var s string

// 	switch t {

// 	case TypeContainer:
// 		s = "Container"

// 	case TypeDeployment:
// 		s = "Deployment"

// 	default:
// 		s = "Unknown"

// 	}

// 	return s
// }

// const (
// 	// TypeContainer is a container registry
// 	TypeContainer Type = iota + 1
// 	// TypeDeployment is a deployment registry
// 	TypeDeployment
// )

// // Config holds the registry config
// type Config struct {
// 	Type     Type
// 	Provider string
// 	Conf     map[string]interface{}
// }

// // DefaultConfig returns a default registry config
// func DefaultConfig() *Config {
// 	return &Config{
// 		Type: TypeContainer,
// 		Conf: make(map[string]interface{}),
// 	}
// }

// Registry implements a registry interface
type Registry interface {
	ID() string
	// Initialize the registry provider
	Init(conf *config.RegistryConfig) error
	// Create a new repository
	Create(string) (interface{}, error)
	// Get repo info
	Get(string) (interface{}, error)
	// Get image manifest
	GetManifest(name, tag string) (interface{}, error)
	// Name of the image with the registry. Needed for deployments
	ImageName(string) string
}

// New returns a new registry based on the config.
// It returns an error if an unsupported provider is supplied or fails to
// initialize the underlying registry provider
func New(conf *config.RegistryConfig) (Registry, error) {
	var (
		reg Registry
		err error
	)

	// switch conf.Type {
	// case TypeContainer:
	// default:
	// 	return nil, fmt.Errorf("unsupported registry type: '%v'", conf.Type)
	// }

	switch conf.Provider {
	case "ecr":
		reg = &awsContainerRegistry{}

	case "docker":
		reg = &localDocker{}

	case "dockerhub":
		reg = &dockerHub{}

	default:
		err = fmt.Errorf("unsupported container registry: '%s'", conf.Provider)

	}

	if err != nil {
		return nil, err
	}

	err = reg.Init(conf)
	return reg, err
}
