package thrap

import (
	"context"
	"errors"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/thrapb"
)

// const (
// 	deplDBKey    = "deployment"
// 	instTableKey = "instance"
// )

var (
	errDeployDescNotSet = errors.New("deployment descriptor not set")
)

// Deployments is used to manage a single projects deployments
type Deployments struct {
	// Project associated to the deployments are being managed
	proj thrapb.Project

	// Deployment descriptor i.e config/template used for all deploys across
	// profiles and instances.  in-mem cache
	desc *thrapb.DeploymentDescriptor

	// Thrap core instance
	t *Thrap
}

// NewDeployments returns a new Deployments instance for a given project
func NewDeployments(t *Thrap, proj thrapb.Project) *Deployments {
	depl := &Deployments{
		t:    t,
		proj: proj,
	}

	depl.loadDescriptor()

	return depl
}

// List returns a list of all deployments for the project
func (d *Deployments) List() ([]*thrapb.Deployment, error) {
	return d.t.deploys.List(d.proj.ID, "")
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

	err = d.t.deploys.Create(d.proj.ID, profileID, nd)
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

	return newDeployment(d.proj, d.desc, nd, eng), nil
}

// Get returns an existing deployment given the profile and instance name,
// which can then be used to perform deployments.
func (d *Deployments) Get(ctx context.Context, profID, instance string) (*Deployment, error) {
	eng, err := d.t.Engine(ctx, profID)
	if err != nil {
		return nil, err
	}

	dp, err := d.t.deploys.Get(d.proj.ID, profID, instance)
	if err != nil {
		return nil, err
	}

	// Table with proj id. diff from above
	return newDeployment(d.proj, d.desc, dp, eng), nil
}

// Descriptor returns the current loaded deployment descriptor
func (d *Deployments) Descriptor() *thrapb.DeploymentDescriptor {
	return d.desc
}

// SetDescriptor sets the deployment descriptor in the store.
func (d *Deployments) SetDescriptor(desc *thrapb.DeploymentDescriptor) error {
	return d.t.descs.Set(d.proj.ID, desc)
}

// loadDescriptor loads the deployment descriptor from hard state
func (d *Deployments) loadDescriptor() error {
	desc, err := d.t.descs.Get(d.proj.ID)
	if err == nil {
		d.desc = desc
	}
	return err
}
