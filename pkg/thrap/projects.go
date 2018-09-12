package thrap

import (
	"context"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/storage"
	"github.com/euforia/thrap/thrapb"
	"github.com/pkg/errors"
)

const (
	projDBKey = "project"
)

// ProjectCreateRequest is a request to create a new project
type ProjectCreateRequest struct {
	Project *thrapb.Project
	// Secrets to add to the project during creation
	Secrets map[string]interface{}
}

// Projects is used to interact with a set of projects
type Projects struct {
	t *Thrap
	// global datastore
	ds kvdb.Datastore
	// projects scoped db
	db kvdb.DB
}

// NewProjects returns a Projects instance backed by the given
// datastore
func NewProjects(t *Thrap) *Projects {
	return &Projects{
		t:  t,
		ds: t.ds,
		db: t.ds.GetDB(projDBKey),
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
func (p *Projects) Create(ctx context.Context, req *ProjectCreateRequest) (*Project, error) {
	proj := req.Project
	table, _ := p.db.GetTable(descTableKey, proj)

	if proj.Name == "" {
		proj.Name = proj.ID
	}

	err := table.Create([]byte(proj.ID), proj)
	if err != nil {
		return nil, err
	}

	err = p.setupProjectSecrets(ctx, proj.ID, req.Secrets)

	return newProject(p.t, proj), err
}

// setupProjectSecrets sets up project secrets across all profiles. It does not
// stop on errors rather completes all operations returning the last encountered
// error
func (p *Projects) setupProjectSecrets(ctx context.Context, projID string, secrets map[string]interface{}) (err error) {
	profiles := p.t.profiles.List()

	// Setup secrets for each profile
	for _, profile := range profiles {
		eng, er := p.t.Engine(ctx, profile.ID)
		if er != nil {
			err = errors.Wrapf(err, "failed to initialize engine for profile %s", profile.ID)
			continue
		}

		er = eng.SetupSecrets(projID, secrets)
		if er != nil {
			err = errors.Wrapf(err, "failed to setup secrets for project %s", projID)
		}
	}

	return
}

// Get returns a object to manage the project
func (p *Projects) Get(id string) (*Project, error) {
	table, _ := p.db.GetTable(descTableKey, &thrapb.Project{})

	obj, err := table.Get([]byte(id))
	if err == nil {
		proj := obj.(*thrapb.Project)
		return newProject(p.t, proj), nil
	}

	return nil, err
}

// Delete deletes the project from the db. Currently it does not purge the
// underlying resources that were created in the Create call
func (p *Projects) Delete(id string) error {
	table, _ := p.db.GetTable(descTableKey, &thrapb.Project{})
	return table.Delete([]byte(id))
}
