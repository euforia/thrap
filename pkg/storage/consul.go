package storage

import (
	"os"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/hashicorp/consul/api"
)

// ConsulStorage implements a consul backed Storage interface
type ConsulStorage struct {
	client *api.Client
}

// NewConsulStorage returns a new ConsulStorage instance
func NewConsulStorage(conf *provider.Config) (*ConsulStorage, error) {
	cc := api.DefaultConfig()
	cc.Address = os.Getenv("CONSUL_ADDR")
	if conf.Addr != "" {
		cc.Address = conf.Addr
	}

	client, err := newConsulClient(cc)
	if err == nil {
		return &ConsulStorage{client: client}, nil
	}
	return nil, err
}

func (s *ConsulStorage) Project() ProjectStorage {
	return NewConsulProjectStorageFromClient(s.client, "thrap/project")
}

func (s *ConsulStorage) Deployment() DeploymentStorage {
	return NewConsulDeployStorageFromClient(s.client, "thrap/deployment")
}

func (s *ConsulStorage) DeployDesc() DeployDescStorage {
	return NewConsulDeployDescStorageFromClient(s.client, "thrap/deployment")
}
