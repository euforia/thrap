package deployment

import (
	"path/filepath"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/provider/orchestrator"
	"github.com/euforia/thrap/thrapb"
)

// Deployment is used to manage a single deployment instance of a project
type Deployment struct {
	// User supplied definition for this deploy
	thrapb.Deployment

	// Project being deployed
	proj thrapb.Project

	// Deploy template/configuration
	desc *thrapb.DeploymentDescriptor

	// Table to store deploy state
	state kvdb.Table

	orch orchestrator.Orchestrator
}

// NewDeployment returns a Deployment instance
func NewDeployment(proj thrapb.Project, desc *thrapb.DeploymentDescriptor,
	deploy *thrapb.Deployment, state kvdb.Table) *Deployment {

	return &Deployment{
		Deployment: *deploy,
		proj:       proj,
		desc:       desc,
		state:      state,
	}
}

// Deploy performs a deployment and updates the internal state
func (d *Deployment) Deploy(opt orchestrator.RequestOptions) error {
	if d.desc == nil || len(d.desc.Spec) == 0 {
		return errDeployDescNotSet
	}

	// We copy the spec to the deployment as input.  The orchestrator consumes and
	// returns the deployed/edited spec.
	d.Spec = make([]byte, len(d.desc.Spec))
	copy(d.Spec, d.desc.Spec)

	//
	// TODO: mormalize apply
	//

	// Version up
	d.Version++

	req := &orchestrator.DeploymentRequest{
		Project:    d.proj,
		Deployment: d.Deployment,
	}

	prepared, err := d.orch.Prepare(req)
	if err != nil {
		return err
	}

	err = d.checkArtifacts(prepared)
	if err != nil {
		return err
	}

	//
	// TODO: do deployment
	//

	key := filepath.Join(d.Profile.ID, d.Name)
	return d.state.Update([]byte(key), &d.Deployment)
}

func (d *Deployment) checkArtifacts(p orchestrator.PreparedDeployment) error {
	arts := p.Artifacts()

}
