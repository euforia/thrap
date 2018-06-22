package orchestrator

import (
	"errors"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	nomad "github.com/hashicorp/nomad/api"
	//"github.com/hashicorp/nomad/nomad/structs"
)

// DeployOptions holds available deployment options
type DeployOptions struct {
	// If true only a report of actions to be taken is generated. An actual
	// deploy is not performed
	Dryrun bool
}

type nomadOrchestrator struct {
	client *nomad.Client
}

// Environment Variables:
// NOMAD_ADDR
func (orch *nomadOrchestrator) Init(conf map[string]interface{}) error {
	var (
		config = nomad.DefaultConfig()
		err    error
	)

	if iaddr, ok := conf["addr"]; ok {
		if addr, ok := iaddr.(string); ok {
			config.Address = addr
		}
	}

	orch.client, err = nomad.NewClient(config)

	return err
}

func (orch *nomadOrchestrator) Deploy(st *thrapb.Stack, opts DeployOptions) (resp interface{}, job interface{}, err error) {
	var njob *nomad.Job
	njob, err = manifest.MakeNomadJob(st)
	if err != nil {
		return
	}
	njob.Canonicalize()

	job = njob

	jobs := orch.client.Jobs()
	q := &nomad.WriteOptions{
		// Region:"",
	}

	// _, _, err = jobs.Validate(njob, q)
	// if err != nil {
	// 	return nil, njob, err
	// }
	// else {
	// r.Warnings
	// r.Error
	// for i:=range r.ValidationErrors{}
	// }

	if opts.Dryrun {
		planOpts := &nomad.PlanOptions{Diff: true}
		resp, _, err = jobs.PlanOpts(njob, planOpts, q)
		return
	}

	// TODO:
	// regOpts := &nomad.RegisterOptions{}
	// resp, _, err = jobs.RegisterOpts(njob, regOpts, q)
	// return
	return nil, nil, errors.New("to be implemented")
}
