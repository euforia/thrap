package storage

import (
	"errors"
	"os"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/consul/api"
)

var (
	// ErrInvalidType is used when there is a type mismatch
	ErrInvalidType = errors.New("invalid type")
)

// IdentityStorage implements a storage interface for Identity
type IdentityStorage interface {
	Get(id string) (*thrapb.Identity, error)
	Update(proj *thrapb.Identity) error
	Create(proj *thrapb.Identity) error
	Delete(id string) error
	Iter(start string, cb func(*thrapb.Identity) error) error
}

// ProfileStorage implements a storage interface for profiles
type ProfileStorage interface {
	Get(id string) (*thrapb.Profile, error)
	List() []*thrapb.Profile
	GetDefault() *thrapb.Profile
	SetDefault(id string) error
}

// ProjectStorage implements storage to persist projects
type ProjectStorage interface {
	Get(id string) (*thrapb.Project, error)
	Update(proj *thrapb.Project) error
	Create(proj *thrapb.Project) error
	Delete(id string) error
	Iter(start string, cb func(*thrapb.Project) error) error
}

// DeploymentStorage implements storage to persist deployment instances
type DeploymentStorage interface {
	Get(project, profile, id string) (*thrapb.Deployment, error)
	Update(string, string, *thrapb.Deployment) error
	Create(string, string, *thrapb.Deployment) error
	Delete(project, profile, id string) error
	List(project string, prefix string) ([]*thrapb.Deployment, error)
}

// DeployDescStorage implements storage to persist deployment descriptors
type DeployDescStorage interface {
	Get(string) (*thrapb.DeploymentDescriptor, error)
	Set(string, *thrapb.DeploymentDescriptor) error
	Delete(string) error
}

type Storage interface {
	Deployment() DeploymentStorage
	Project() ProjectStorage
	DeployDesc() DeployDescStorage
}

type ConsulStorage struct {
	client *api.Client
}

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
