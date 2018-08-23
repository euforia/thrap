package orchestrator

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/thrapb"
)

var (
	errOrchProviderMissing      = errors.New("orchestration provider missing")
	errOrchProviderNotSupported = errors.New("orchestration provider not supported")
)

// RequestOptions holds available deployment options
type RequestOptions struct {
	// If true only a report of actions to be taken is generated. An actual
	// deploy is not performed
	Dryrun bool

	// Progress output
	Output io.Writer

	// This is only for destroy request
	Purge bool

	// Region to make the request
	Region string
}

// DeploymentRequest is a raw deployment request.
type DeploymentRequest struct {
	Project    thrapb.Project
	Deployment thrapb.Deployment
}

// PreparedDeployment implements a prepared deployment that the orchestrator
// then uses to perform a deploy
type PreparedDeployment interface {
	// Artifacts needed for the deployment
	Artifacts() []string
	// Orchestrator specific spec
	Spec() interface{}
	// Serialized bytes of a deploy after all normalization.  This is stored
	// along with the deployment
	Bytes() []byte
}

// Orchestrator implements an application/project deployment orchestrator
type Orchestrator interface {
	// Init is called to initialize the orchestrator with the given config
	Init(config *provider.Config) error

	// ID of the orchestrator. Each orchestrator must have a unique id
	ID() string

	// Prepares a deployment.  The output of this is used to call deploy
	PrepareDeploy(req *DeploymentRequest) (PreparedDeployment, error)

	// Deploy should deploy the stack returning the response, deploy object
	// based on the orchestrator or an error
	Deploy(context.Context, PreparedDeployment, RequestOptions) (interface{}, error)

	// // Destroy the stack returning results for each component
	// Destroy(ctx context.Context, stack *thrapb.Stack) []*thrapb.ActionResult

	// // Status of all comps
	// Status(ctx context.Context, stack *thrapb.Stack) []*thrapb.CompStatus
}

// New returns a new orchestrator based on the given config
func New(conf *provider.Config) (Orchestrator, error) {
	var (
		orch Orchestrator
		err  error
	)

	switch conf.Provider {
	case "nomad":
		orch = &nomadOrchestrator{}

	case "":
		err = errors.Wrapf(errOrchProviderMissing, conf.ID)

	default:
		err = errors.Wrap(errOrchProviderNotSupported, conf.Provider)

	}

	if err == nil {
		err = orch.Init(conf)
	}

	return orch, err
}