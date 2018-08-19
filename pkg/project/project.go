package project

import (
	"bytes"
	"crypto/sha256"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/deployment"
	"github.com/euforia/thrap/thrapb"
)

// Project is used to manage a project
type Project struct {
	*thrapb.Project
	// project data hash
	hash []byte
	// whole datastore
	ds kvdb.Datastore
	// underlying db
	db kvdb.DB
}

func newProject(proj *thrapb.Project, ds kvdb.Datastore, db kvdb.DB) *Project {
	p := &Project{
		ds:      ds,
		db:      db,
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

	// Update
	table, _ := p.db.GetTable("descriptor", &thrapb.Project{})
	return table.Update([]byte(p.ID), p.Project)
}

// Deployments returns a Deployments instance to manage the project's
// deployments
func (p *Project) Deployments() *deployment.Deployments {
	return deployment.New(*p.Clone(), p.ds)
}
