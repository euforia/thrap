package storage

import (
	"errors"

	"github.com/euforia/thrap/thrapb"
)

var (
	// ErrInvalidType is used when there is a type mismatch
	ErrInvalidType = errors.New("invalid type")
)

type IdentityStorage interface {
	Get(id string) (*thrapb.Identity, error)
	Update(proj *thrapb.Identity) error
	Create(proj *thrapb.Identity) error
	Delete(id string) error
	Iter(start string, cb func(*thrapb.Identity) error) error
}

// type DeploymentStorage interface {
// 	Get(id uint64) (*thrapb.Deployment, error)
// 	Update(proj *thrapb.Deployment) error
// 	Create(proj *thrapb.Deployment) error
// 	Delete(id uint64) error
// 	Iter(id uint64, cb func(*thrapb.Deployment) error) error
// }

// type ProjectStorage interface {
// 	Get(id string) (*thrapb.Project, error)
// 	Update(proj *thrapb.Project) error
// 	Create(proj *thrapb.Project) error
// 	Delete(id string) error
// 	Iter(start string, cb func(*thrapb.Project) error) error
// }

// Datastore is used to managed all datastores. This is the top-level
// object
// type Datastore struct {
// 	ds kvdb.Datastore
// }

// func NewDatastore(dir string) (*Datastore, error) {
// 	// kvdb.NewBadgerDatastoreFromDB()
// 	ds, err := kvdb.NewBadgerDatastore(dir)
// 	if err == nil {
// 		return &Datastore{ds: ds}, nil
// 	}
// 	return nil, err
// }

// func (s *Datastore) Identity() (IdentityStorage, error) {
// 	db := s.ds.GetDB("")
// 	table, err := db.GetTable("identity", &thrapb.Identity{})
// 	if err == nil {
// 		return &KvdbIdentityStorage{table: table}, nil
// 	}
// 	return nil, err
// }

// func (s *Datastore) Project() (ProjectStorage, error) {
// 	db := s.ds.GetDB("")
// 	table, err := db.GetTable("project", &thrapb.Project{})
// 	if err == nil {
// 		return &KvdbProjectStorage{table: table}, nil
// 	}

// 	return nil, err
// }
