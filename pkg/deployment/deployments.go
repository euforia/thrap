package deployment

import (
	"errors"
	"path/filepath"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/storage"
	"github.com/euforia/thrap/thrapb"
)

const (
	deplDBKey    = "deployment"
	descTableKey = "descriptor"
	instTableKey = "instance"
)

var (
	errDeployDescNotSet = errors.New("deployment descriptor not set")
)

// Deployments is used to manage a single projects deployments
type Deployments struct {
	// Deployment descriptor i.e config/template used for all deploys across
	// profiles and instances
	desc *thrapb.DeploymentDescriptor

	// Project who's deployments are being managed
	proj thrapb.Project

	// Profiles
	profiles storage.ProfileStorage

	// Project deployment db
	db kvdb.DB

	// Global datastore
	ds kvdb.Datastore
}

// New returns a new Deployments instance for a given project
func New(proj thrapb.Project, ds kvdb.Datastore) *Deployments {
	// Deployment db for project
	pdb := ds.GetDB(deplDBKey)

	depl := &Deployments{
		proj: proj,
		db:   pdb,
		ds:   ds,
	}

	depl.loadDescriptor()

	return depl
}

// List returns a list of all deployments for the project
func (d *Deployments) List() ([]*thrapb.Deployment, error) {
	table, _ := d.db.GetTable(instTableKey, &thrapb.Deployment{})

	out := make([]*thrapb.Deployment, 0)
	err := table.Iter(nil, func(obj kvdb.Object) error {
		depl := obj.(*thrapb.Deployment)
		out = append(out, depl)
		return nil
	})

	return out, err
}

// Create creates a new deploy for the project with the given profile and instance
// name. This is simply an initialization record. This must called before an
// actual deploy can be performed.
func (d *Deployments) Create(dpl *thrapb.Deployment) (*Deployment, error) {
	err := dpl.Validate()
	if err != nil {
		return nil, err
	}

	dpl.Version = 0
	table, _ := d.db.GetTable(instTableKey, dpl)
	key := d.getkey(dpl.Profile.ID, dpl.Name)

	err = table.Create([]byte(key), dpl)
	if err == nil {
		dtable, _ := d.db.GetTable(filepath.Join(instTableKey, d.proj.ID), dpl)
		return NewDeployment(d.proj, d.desc, dpl, dtable), nil
	}

	return nil, err
}

// Get returns an existing deployment given the profile and instance name,
// which can then be used to perform deployments.
func (d *Deployments) Get(profID, instance string) (*Deployment, error) {
	td := &thrapb.Deployment{}
	table, _ := d.db.GetTable(instTableKey, td)

	key := d.getkey(profID, instance)
	obj, err := table.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	dp := obj.(*thrapb.Deployment)

	// Table with proj id. diff from above
	dtable, _ := d.db.GetTable(filepath.Join(instTableKey, d.proj.ID), td)
	return NewDeployment(d.proj, d.desc, dp, dtable), nil
}

// Descriptor returns the current loaded deployment descriptor
func (d *Deployments) Descriptor() *thrapb.DeploymentDescriptor {
	return d.desc
}

// SetDescriptor sets the deployment descriptor in the store.
func (d *Deployments) SetDescriptor(desc *thrapb.DeploymentDescriptor) error {
	db := d.ds.GetDB(deplDBKey)
	table, _ := db.GetTable(descTableKey, desc)

	var err error
	if d.desc == nil {
		err = table.Create([]byte(d.proj.ID), desc)
	} else {
		err = table.Update([]byte(d.proj.ID), desc)
	}

	if err == nil {
		d.desc = desc
	}

	return err
}

func (d *Deployments) getkey(prof, instance string) string {
	return filepath.Join(d.proj.ID, prof, instance)
}

// loadDescriptor loads the deployment descriptor from hard state
func (d *Deployments) loadDescriptor() error {
	db := d.ds.GetDB(deplDBKey)
	table, _ := db.GetTable(descTableKey, &thrapb.DeploymentDescriptor{})
	obj, err := table.Get([]byte(d.proj.ID))
	if err == nil {
		d.desc = obj.(*thrapb.DeploymentDescriptor)
	}
	return err
}
