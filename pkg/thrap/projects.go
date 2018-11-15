package thrap

import (
	"context"
	"log"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/pkg/storage"
	"github.com/pkg/errors"
)

const (
	projDBKey = "project"
)

var (
	// ErrProjectSourceRequired ...
	ErrProjectSourceRequired = errors.New("project source required")
)

// ProjectCreateRequest is a request to create a new project
type ProjectCreateRequest struct {
	Project *pb.Project
	// Secrets to add to the project during creation
	Secrets map[string]interface{}
}

// Projects is used to interact with a set of projects
type Projects struct {
	t     *Thrap
	store storage.ProjectStorage
}

// NewProjects returns a Projects instance backed by the given
// datastore
func NewProjects(t *Thrap) *Projects {
	return &Projects{
		t:     t,
		store: t.store.Project(),
	}
}

// Iter iterates over all projects from the start point
func (p *Projects) Iter(start string, cb func(*pb.Project) error) error {
	return p.store.Iter(start, cb)
}

// Create creates a new project and returns a object to manage the project
func (p *Projects) Create(ctx context.Context, req *ProjectCreateRequest) (*Project, error) {
	proj := req.Project
	if proj.Name == "" {
		proj.Name = proj.ID
	}
	if proj.Source == "" {
		return nil, ErrProjectSourceRequired
	}

	err := p.store.Create(proj)
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

		log.Printf("Secrets setup project=%s profile=%s", projID, profile.ID)

		er = eng.SetupSecrets(projID, secrets)
		if er != nil {
			err = errors.Wrapf(err, "failed to setup secrets for project %s", projID)
		}
	}

	return
}

// Get returns a object to manage the project
func (p *Projects) Get(id string) (*Project, error) {
	proj, err := p.store.Get(id)
	if err != nil {
		return nil, err
	}

	return newProject(p.t, proj), nil
}

// Delete deletes the project from the db. Currently it does not purge the
// underlying resources that were created in the Create call
func (p *Projects) Delete(id string) error {
	return p.store.Delete(id)
}
