package thrapb

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// Container holds a container configuration.  For now everything is
// translated into and out of this structure
type Container struct {
	// Name of deployed container.  This must be unique per container
	// engine instance
	Name      string
	Container *container.Config
	Host      *container.HostConfig
	Network   *network.NetworkingConfig
}

// NewContainer returns a container configuration using the stack and
// component ids
func NewContainer(stackID, compID string) *Container {
	return &Container{
		Name:      compID + "." + stackID,
		Container: &container.Config{},
		Host: &container.HostConfig{
			// Create isolated network by stack id
			NetworkMode: container.NetworkMode(stackID),
		},
		Network: &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				// Map access to isolated network from above
				stackID: &network.EndpointSettings{},
			},
		},
	}
}
