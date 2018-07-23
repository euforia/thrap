package orchestrator

import (
	"context"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	nomad "github.com/hashicorp/nomad/api"
	//"github.com/hashicorp/nomad/nomad/structs"
)

type nomadOrchestrator struct {
	client *nomad.Client
}

func (orch *nomadOrchestrator) ID() string {
	return "nomad"
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

func (orch *nomadOrchestrator) Deploy(ctx context.Context, st *thrapb.Stack, opts RequestOptions) (resp interface{}, job interface{}, err error) {
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

	regOpts := &nomad.RegisterOptions{}
	resp, _, err = jobs.RegisterOpts(njob, regOpts, q)
	return
}

func (orch *nomadOrchestrator) Destroy(ctx context.Context, stack *thrapb.Stack) []*thrapb.ActionReport {
	jobs := orch.client.Jobs()
	q := &nomad.WriteOptions{}
	_, _, err := jobs.Deregister(stack.ID, true, q)

	ar := make([]*thrapb.ActionReport, 0, len(stack.Components))
	for _, c := range stack.Components {
		r := &thrapb.ActionReport{Action: thrapb.NewAction("destroy", "comp", c.ID)}
		r.Error = err
		ar = append(ar, r)
	}

	return ar
}
