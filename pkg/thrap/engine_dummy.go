package thrap

import (
	"context"

	"github.com/euforia/thrap/pkg/provider/orchestrator"
)

type dummyPreparedDeploy struct {
	bytes []byte
}

func (d *dummyPreparedDeploy) Artifacts() []string {
	return []string{}
}

func (d *dummyPreparedDeploy) Bytes() []byte {
	return d.bytes
}

func (d *dummyPreparedDeploy) Spec() interface{} {
	return struct{}{}
}

// dummyEngine for testing
type dummyEngine struct {
	*engine
}

func newDummyEngine(eng *engine) *dummyEngine {
	return &dummyEngine{eng}
}

func (eng *dummyEngine) Deploy(context.Context, orchestrator.PreparedDeployment, orchestrator.RequestOptions) error {
	return nil
}
