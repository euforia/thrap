package deployment

import (
	"path/filepath"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/thrapb"
)

// Deployment is used to manage a single deployment instance of a project
type Deployment struct {
	// User supplied definition for this deploy
	thrapb.Deployment

	// Deploy template/configuration
	desc *thrapb.DeploymentDescriptor

	// Table to store deploy state
	state kvdb.Table
}

// NewDeployment returns a Deployment instance
func NewDeployment(desc *thrapb.DeploymentDescriptor,
	deploy *thrapb.Deployment, state kvdb.Table) *Deployment {

	return &Deployment{
		state:      state,
		desc:       desc,
		Deployment: *deploy,
	}
}

// Deploy performs a deployment and updates the internal state
func (d *Deployment) Deploy() error {
	if d.desc == nil || len(d.desc.Spec) == 0 {
		return errDeployDescNotSet
	}

	//
	// TODO: mormalize apply
	//
	// d.Spec = make([]byte, len(d.desc.Spec))
	// copy(d.Spec, d.desc.Spec)

	// Version up
	d.Version++

	//
	// TODO: do deployment
	//

	key := filepath.Join(d.Profile.ID, d.Name)
	return d.state.Update([]byte(key), &d.Deployment)
}
