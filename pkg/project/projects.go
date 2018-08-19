package project

import (
	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/storage"
	"github.com/euforia/thrap/thrapb"
)

const (
	descTableKey = "descriptor"
	projDBKey    = "project"
)

// Projects is used to interact with a set of projects
type Projects struct {
	ds kvdb.Datastore // global datastore
	db kvdb.DB        // projects scoped db
}

// NewProjects returns a Projects instance backed by the given
// datastore
func NewProjects(ds kvdb.Datastore) *Projects {
	return &Projects{
		ds: ds,
		db: ds.GetDB(projDBKey),
	}
}

// Iter iterates over all projects from the start point
func (p *Projects) Iter(start string, cb func(*thrapb.Project) error) error {
	table, _ := p.db.GetTable(descTableKey, &thrapb.Project{})

	var prefix []byte
	if start != "" {
		prefix = []byte(start)
	}

	return table.Iter(prefix, func(obj kvdb.Object) error {
		proj, ok := obj.(*thrapb.Project)
		if ok {
			return cb(proj)
		}
		return storage.ErrInvalidType
	})
}

// Create creates a new project and returns a object to manage the project
func (p *Projects) Create(proj *thrapb.Project) (*Project, error) {
	table, _ := p.db.GetTable(descTableKey, proj)

	err := table.Create([]byte(proj.ID), proj)
	if err == nil {
		return newProject(proj, p.ds, p.db), nil
	}

	return nil, err
}

// Get returns a object to manage the project
func (p *Projects) Get(id string) (*Project, error) {
	table, _ := p.db.GetTable(descTableKey, &thrapb.Project{})

	obj, err := table.Get([]byte(id))
	if err == nil {
		proj := obj.(*thrapb.Project)
		return newProject(proj, p.ds, p.db), nil
	}

	return nil, err
}
