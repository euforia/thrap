package thrap

import (
	"bytes"
	"crypto/sha256"

	"github.com/euforia/thrap/pkg/storage"
	"github.com/euforia/thrap/thrapb"
)

// Project is used to manage a project
type Project struct {
	// Project data
	*thrapb.Project

	// Project data hash used to know if a write should be
	// performed
	hash []byte

	t *Thrap

	store storage.ProjectStorage
}

func newProject(t *Thrap, proj *thrapb.Project) *Project {
	p := &Project{
		t:       t,
		Project: proj,
		store:   t.store.Project(),
	}

	p.init()

	return p
}

func (p *Project) init() {
	p.hash = p.Project.Hash(sha256.New())
}

// Sync persists all project changes to the db
func (p *Project) Sync() error {
	hash := p.Project.Hash(sha256.New())
	// Nothings changed
	if bytes.Compare(p.hash, hash) == 0 {
		return nil
	}

	return p.store.Update(p.Project)
}

// Deployments returns a Deployments instance to manage the project's
// deployments
func (p *Project) Deployments() *Deployments {
	return NewDeployments(p.t, *p.Clone())
}
