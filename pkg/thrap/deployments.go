package thrap

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/thrapb"
)

const (
	deplDBKey    = "deployment"
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

	// Project deployment db
	db kvdb.DB

	// Thrap core instance
	t *Thrap
}

// NewDeployments returns a new Deployments instance for a given project
func NewDeployments(t *Thrap, proj thrapb.Project) *Deployments {
	// Deployment db
	pdb := t.ds.GetDB(deplDBKey)

	depl := &Deployments{
		t:    t,
		proj: proj,
		db:   pdb,
	}

	depl.loadDescriptor()

	return depl
}

// List returns a list of all deployments for the project
func (d *Deployments) List() ([]*thrapb.Deployment, error) {
	tableKey := filepath.Join(instTableKey, d.proj.ID)
	table, _ := d.db.GetTableVersion(tableKey, &thrapb.Deployment{}, d.t.hashFunc)

	out := make([]*thrapb.Deployment, 0)
	err := table.Iter(nil, func(obj kvdb.ObjectVersion) error {
		depl := obj.(*thrapb.Deployment)
		out = append(out, depl)
		return nil
	})

	return out, err
}

// Create creates a new deploy for the project with the given profile and instance
// name. This is simply an initialization record. This must called before an
// actual deploy can be performed.
func (d *Deployments) Create(ctx context.Context, profileID, instanceName string) (*Deployment, error) {

	nd := thrapb.NewDeployment(profileID, instanceName)
	nd.State = thrapb.DeploymentState_CREATE
	nd.Status = thrapb.DeployStateStatus_OK

	err := nd.Validate()
	if err != nil {
		return nil, err
	}

	// Ensure valid profile and load associated engine
	eng, err := d.t.Engine(ctx, profileID)
	if err != nil {
		return nil, err
	}

	tableKey := filepath.Join(instTableKey, d.proj.ID)
	table, _ := d.db.GetTableVersion(tableKey, nd, d.t.hashFunc)

	// key := d.getkey(dpl.Profile.ID, dpl.Name)
	key := filepath.Join(profileID, instanceName)

	_, err = table.Create([]byte(key), nd)
	if err != nil {
		return nil, err
	}

	seedReq := &provider.Request{
		Project: d.proj,
		Deployment: thrapb.Deployment{
			Name:    instanceName,
			Profile: &thrapb.Profile{ID: profileID},
		},
	}

	if err := eng.SeedSecrets(seedReq); err != nil {
		d.t.log.Println("Not seeding secrets:", err)
	}

	// dtable, _ := d.db.GetTableVersion(filepath.Join(instTableKey, d.proj.ID), dpl, d.t.hashFunc)

	return newDeployment(d.proj, d.desc, nd, eng, table), nil
}

// Get returns an existing deployment given the profile and instance name,
// which can then be used to perform deployments.
func (d *Deployments) Get(ctx context.Context, profID, instance string) (*Deployment, error) {
	eng, err := d.t.Engine(ctx, profID)
	if err != nil {
		return nil, err
	}

	tableKey := filepath.Join(instTableKey, d.proj.ID)

	td := &thrapb.Deployment{}
	table, _ := d.db.GetTableVersion(tableKey, td, d.t.hashFunc)

	key := filepath.Join(profID, instance)
	obj, _, err := table.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	dp := obj.(*thrapb.Deployment)

	// Table with proj id. diff from above
	return newDeployment(d.proj, d.desc, dp, eng, table), nil
}

// Descriptor returns the current loaded deployment descriptor
func (d *Deployments) Descriptor() *thrapb.DeploymentDescriptor {
	return d.desc
}

// SetDescriptor sets the deployment descriptor in the store.
func (d *Deployments) SetDescriptor(desc *thrapb.DeploymentDescriptor) error {
	table, _ := d.db.GetTable(descTableKey, desc)

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

// loadDescriptor loads the deployment descriptor from hard state
func (d *Deployments) loadDescriptor() error {
	// db := d.ds.GetDB(deplDBKey)
	table, _ := d.db.GetTable(descTableKey, &thrapb.DeploymentDescriptor{})
	obj, err := table.Get([]byte(d.proj.ID))
	if err == nil {
		d.desc = obj.(*thrapb.DeploymentDescriptor)
	}
	return err
}
