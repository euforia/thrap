package thrap

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/pkg/provider/orchestrator"
	"github.com/euforia/thrap/pkg/storage"
)

var (
	errDeployAlreadyTriggered = errors.New("deployment already triggered")
	errReqVarNotSet           = errors.New("required variable not set")
)

// Vault token login
// recursive template copy
// secondary service block

// DeployRequest is the user facing request
type DeployRequest struct {
	Descriptor string
	Vars       map[string]string
	Dryrun     bool
}

// Deployment is used to manage a single deployment instance of a project
type Deployment struct {
	// Project being deployed
	proj pb.Project

	// Template/configuration tied to the deployment
	desc *pb.DeploymentDescriptor

	// User supplied definition for this deploy
	depl pb.Deployment

	// Engine loaded with the deployment profile
	eng Engine

	deploys storage.DeploymentStorage
	descs   storage.DeployDescStorage
}

func newDeployment(proj pb.Project, deploy *pb.Deployment,
	eng Engine,
	descs storage.DeployDescStorage,
	deploys storage.DeploymentStorage) *Deployment {

	return &Deployment{
		proj:    proj,
		depl:    *deploy,
		eng:     eng,
		deploys: deploys,
		descs:   descs,
	}
}

// Deployable returns the deploy object tied to this deployment.
func (d *Deployment) Deployable() pb.Deployment {
	depl := d.depl.Clone()
	return *depl
}

// PrepareDeploy prepares a deployment.  This must be called before a
// call to Deploy can be made
func (d *Deployment) PrepareDeploy(ctx context.Context, req *DeployRequest) (orchestrator.PreparedDeployment, error) {
	var err error
	if len(req.Descriptor) == 0 {
		d.desc, err = d.descs.Get(d.proj.ID)
	} else {
		d.desc, err = d.descs.GetVersion(d.proj.ID, req.Descriptor)
	}
	if err != nil {
		return nil, err
	}

	if d.isVoid() {
		return nil, fmt.Errorf("deployment cannot be prepared in state=%s status=%s",
			d.depl.State.String(), d.depl.Status.String())
	}

	if err := d.setupPrepare(req.Vars); err != nil {
		return nil, err
	}

	orchReq := &provider.Request{
		Project:    d.proj,
		Deployment: d.depl,
	}

	prepared, err := d.eng.PrepareDeploy(ctx, orchReq)
	if err != nil {
		d.depl.Status = pb.DeployStateStatus_FAILED
		d.depl.StateMessage = err.Error()
		d.Sync()

		return nil, err
	}

	// Write the final spec that is deployed back to the deploy object
	// in Prepare state
	// d.depl.Desc.Spec, err = prepared.Bytes(d.depl.Desc.Mime)
	d.depl.Desc.Spec, err = json.Marshal(prepared.Spec())
	if err != nil {
		d.depl.Status = pb.DeployStateStatus_FAILED
		d.depl.StateMessage = err.Error()
		d.Sync()
		return nil, err
	}

	d.depl.Status = pb.DeployStateStatus_OK
	return prepared, d.Sync()
}

// Deploy performs a deployment and updates the internal state
func (d *Deployment) Deploy(ctx context.Context, req *DeployRequest) (*pb.Deployment, error) {
	// TODO: lock()
	// TODO: defer unlock()
	prepared, err := d.PrepareDeploy(ctx, req)
	if err != nil {
		return nil, err
	}

	opt := orchestrator.RequestOptions{Dryrun: req.Dryrun}

	d.depl.State = pb.DeploymentState_DEPLOY
	err = d.eng.Deploy(ctx, prepared, opt)
	if err != nil {
		d.depl.Status = pb.DeployStateStatus_FAILED
		d.depl.StateMessage = err.Error()
		d.Sync()
		return nil, err
	}

	d.depl.Status = pb.DeployStateStatus_OK

	err = d.Sync()
	return &d.depl, err
}

// Stop the deployment
func (d *Deployment) Stop(ctx context.Context) error {
	d.depl.State = pb.DeploymentState_STOP
	err := d.eng.Stop(ctx, d.proj.ID, d.depl.Name)
	if err != nil {
		d.depl.StateMessage = err.Error()
		d.depl.Status = pb.DeployStateStatus_FAILED
	} else {
		d.depl.Status = pb.DeployStateStatus_OK
	}
	d.Sync()
	return err
}

// Destroy stops and destroys the deployment
func (d *Deployment) Destroy(ctx context.Context) error {
	d.depl.State = pb.DeploymentState_STOP
	err := d.eng.Destroy(ctx, d.proj.ID, d.depl.Name)
	if err != nil {
		d.depl.StateMessage = err.Error()
		d.depl.Status = pb.DeployStateStatus_FAILED
	} else {
		d.depl.Status = pb.DeployStateStatus_OK
	}
	d.Sync()
	return err
}

// Sync persists the current deployment to the store
func (d *Deployment) Sync() error {
	d.depl.ModifiedAt = time.Now().UnixNano()
	return d.deploys.Update(d.proj.ID, d.depl.Profile.ID, &d.depl)
}

func (d *Deployment) isVoid() bool {
	state := d.depl.State
	status := d.depl.Status

	switch state {
	case pb.DeploymentState_CREATE, pb.DeploymentState_DEPLOY, pb.DeploymentState_STOP:
		return false
	case pb.DeploymentState_PREPARE:
		if status == pb.DeployStateStatus_FAILED {
			return false
		}
	}

	return true
}

func (d *Deployment) setupPrepare(vars map[string]string) error {
	prof := d.eng.Profile()
	d.depl.Profile = &prof
	d.depl.Nonce = rand.Uint64()
	d.depl.StateMessage = ""

	// Version up
	d.depl.Version++

	// We copy the descriptor to the deployment as input.  The orchestrator consumes and
	// returns the deployed/edited spec.
	d.depl.Desc = d.desc.Clone()
	// d.depl.Desc.Spec = make([]byte, len(d.desc.Spec))
	// copy(d.depl.Spec, d.desc.Spec)

	// Set supplied and internal vars
	d.setProfileVars(vars)

	// Check all vars
	err := d.validateProfileMetaAndVars()
	if err != nil {
		return err
	}

	d.depl.State = pb.DeploymentState_PREPARE

	// Sync here so another concurrent call would bail
	return d.Sync()
}

func (d *Deployment) setProfileVars(vars map[string]string) {
	profile := d.depl.Profile
	if profile.Meta == nil {
		profile.Meta = make(map[string]string)
	}

	// Add our internally injected project and deployment variables.
	profile.Meta[InstanceVarName] = d.depl.Name
	profile.Meta[ProjectVarName] = d.proj.ID
	profile.Meta[SecretsPathVarName] = d.eng.SecretsPath(d.proj.ID, d.depl.Name)
	profile.Meta[DeployVerVarName] = fmt.Sprintf("%d", d.depl.Version)

	// Add user supplied variables to profile
	if len(vars) > 0 {
		if profile.Variables == nil {
			profile.Variables = make(map[string]string)
		}
		for k, v := range vars {
			profile.Variables[k] = v
		}
	}
}

// Make sure all required variables have been supplied i.e no key
// can have an empty value
func (d *Deployment) validateProfileMetaAndVars() error {
	profile := d.depl.Profile
	// Check user supplied/required vars
	if len(profile.Variables) > 0 {
		for k, v := range profile.Variables {
			if v == "" {
				return errors.Wrap(errReqVarNotSet, k)
			}
		}
	}

	// Check internal metadata vars
	if len(profile.Meta) > 0 {
		for k, v := range profile.Meta {
			if v == "" {
				return errors.Wrap(errReqVarNotSet, k)
			}
		}
	}

	return nil
}
