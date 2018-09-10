package thrap

import (
	"bytes"
	"crypto/sha256"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/thrapb"
)

// Project is used to manage a project
type Project struct {
	t *Thrap

	*thrapb.Project
	// project data hash
	hash []byte
	// whole datastore passed to deployments
	ds kvdb.Datastore
	// underlying db project/{id}
	db kvdb.DB
}

func newProject(t *Thrap, proj *thrapb.Project) *Project {
	p := &Project{
		t:       t,
		ds:      t.ds,
		db:      t.ds.GetDB(projDBKey),
		Project: proj,
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

	// Update state
	table, _ := p.db.GetTable("descriptor", &thrapb.Project{})
	return table.Update([]byte(p.ID), p.Project)
}

// Deployments returns a Deployments instance to manage the project's
// deployments
func (p *Project) Deployments() *Deployments {
	return NewDeployments(p.t, *p.Clone())
}
