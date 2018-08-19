package thrapb

import (
	"errors"

	"github.com/euforia/kvdb"
)

func (d *Deployment) New() kvdb.Object {
	return &Deployment{}
}

func (d *Deployment) Validate() error {
	if d.Name == "" {
		return errors.New("name required")
	}
	if d.Profile == nil || len(d.Profile.ID) == 0 {
		return errors.New("profile id required")
	}
	return nil
}

func (desc *DeploymentDescriptor) New() kvdb.Object {
	return &DeploymentDescriptor{}
}
