package thrap

import (
	"context"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/pkg/storage"
	"github.com/pkg/errors"
)

var (
	errDeployDescNotSet = errors.New("deployment descriptor not set")
)

// Deployments is used to manage a single projects deployments
type Deployments struct {
	// Project associated to the deployments are being managed
	proj pb.Project

	// Store with current deployables
	deploys storage.DeploymentStorage

	// Store holding deploy descriptors
	descs storage.DeployDescStorage

	// Thrap core instance
	t *Thrap
}

// NewDeployments returns a new Deployments instance for a given project
func NewDeployments(t *Thrap, proj pb.Project) *Deployments {
	return &Deployments{
		t:       t,
		proj:    proj,
		deploys: t.store.Deployment(),
		descs:   t.store.DeployDesc(),
	}
}

// List returns a list of all deployments for the project
func (d *Deployments) List() ([]*pb.Deployment, error) {
	return d.deploys.List(d.proj.ID, "")
}

// Create creates a new deploy for the project with the given profile and instance
// name. This is simply an initialization record. This must called before an
// actual deploy can be performed.
func (d *Deployments) Create(ctx context.Context, profileID, instanceName string) (*Deployment, error) {
	nd := pb.NewDeployment(profileID, instanceName)
	err := nd.Validate()
	if err != nil {
		return nil, err
	}

	// Ensure valid profile and load associated engine
	eng, err := d.t.Engine(ctx, profileID)
	if err != nil {
		return nil, err
	}

	err = d.deploys.Create(d.proj.ID, profileID, nd)
	if err != nil {
		return nil, err
	}

	seedReq := &provider.Request{
		Project: d.proj,
		Deployment: pb.Deployment{
			Name:    instanceName,
			Profile: &pb.Profile{ID: profileID},
		},
	}

	if err := eng.SeedSecrets(seedReq); err != nil {
		d.t.log.Println("Not seeding secrets:", err)
	}

	return newDeployment(d.proj, nd, eng, d.descs, d.deploys), nil
}

// Get returns an existing deployment given the profile and instance name,
// which can then be used to perform deployments.
func (d *Deployments) Get(ctx context.Context, profID, instance string) (*Deployment, error) {
	eng, err := d.t.Engine(ctx, profID)
	if err != nil {
		return nil, err
	}

	dp, err := d.deploys.Get(d.proj.ID, profID, instance)
	if err == nil {
		return newDeployment(d.proj, dp, eng, d.descs, d.deploys), nil
	}

	return nil, err
}

// Descriptor returns the current loaded deployment descriptor
func (d *Deployments) Descriptor(version string) (*pb.DeploymentDescriptor, error) {
	return d.descs.GetVersion(d.proj.ID, version)
}

// SetDescriptor sets the deployment descriptor in the store.
func (d *Deployments) SetDescriptor(desc *pb.DeploymentDescriptor) error {
	return d.descs.Set(d.proj.ID, desc)
}

// Descriptors lists all descriptors available in project
func (d *Deployments) Descriptors() ([]string, error) {
	return d.descs.ListVersions(d.proj.ID)
}
