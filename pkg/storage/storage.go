package storage

import (
	"errors"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/thrapb"
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
	Get(id string) (*pb.Profile, error)
	List() []*pb.Profile
	GetDefault() *pb.Profile
	SetDefault(id string) error
}

// ProjectStorage implements storage to persist projects
type ProjectStorage interface {
	Get(id string) (*pb.Project, error)
	Update(proj *pb.Project) error
	Create(proj *pb.Project) error
	Delete(id string) error
	Iter(start string, cb func(*pb.Project) error) error
}

// DeploymentStorage implements storage to persist deployment instances
type DeploymentStorage interface {
	Get(project, profile, id string) (*pb.Deployment, error)
	Update(string, string, *pb.Deployment) error
	Create(string, string, *pb.Deployment) error
	Delete(project, profile, id string) error
	List(project string, prefix string) ([]*pb.Deployment, error)
}

// DeployDescStorage implements storage to persist deployment descriptors
type DeployDescStorage interface {
	Get(string) (*pb.DeploymentDescriptor, error)
	Set(string, *pb.DeploymentDescriptor) error
	Delete(string) error
}

// Storage implements an all encompasing storage
type Storage interface {
	Deployment() DeploymentStorage
	Project() ProjectStorage
	DeployDesc() DeployDescStorage
}
