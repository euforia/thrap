package orchestrator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/provider"
	nomad "github.com/hashicorp/nomad/api"
)

type nomadOrchestrator struct {
	client *nomad.Client
}

func (orch *nomadOrchestrator) ID() string {
	return "nomad"
}

// Environment Variables:
// NOMAD_ADDR
func (orch *nomadOrchestrator) Init(c *provider.Config) error {

	var (
		nomadClientConf = nomad.DefaultConfig()
		err             error
	)

	if c.Addr != "" {
		nomadClientConf.Address = c.Addr
	}

	nomadClientConf.SetTimeout(3 * time.Second)

	orch.client, err = nomad.NewClient(nomadClientConf)

	return err
}

func (orch *nomadOrchestrator) PrepareDeploy(ctx context.Context, req *provider.Request) (p PreparedDeployment, err error) {
	// job, err := jobspec.Parse(bytes.NewBuffer(req.Deployment.Spec))

	var apiJob nomad.Job
	err = json.Unmarshal(req.Deployment.Spec, &apiJob)
	if err != nil {
		return
	}

	var (
		jobID = req.Project.ID + "-" + req.Deployment.Name
		job   = &apiJob
	)

	job.ID = &jobID
	job.Name = &jobID

	// Inject platform variables
	profile := req.Deployment.Profile
	if profile.Meta != nil {
		if job.Meta == nil {
			job.Meta = make(map[string]string)
		}
		for k, v := range profile.Meta {
			job.Meta[k] = v
		}
	}

	// Inject user defined variables
	if profile.Variables != nil {
		if job.Meta == nil {
			job.Meta = make(map[string]string)
		}
		for k, v := range profile.Variables {
			job.Meta[k] = v
		}
	}

	job.Canonicalize()

	p, err = newNomadPreparedDeployment(job)

	return
}

func (orch *nomadOrchestrator) Deploy(ctx context.Context, d PreparedDeployment, opt RequestOptions) (resp interface{}, err error) {
	iface := d.Spec()
	job, ok := iface.(*nomad.Job)
	if !ok {
		err = errors.New("not a nomad job")
		return
	}
	// b, _ := json.MarshalIndent(job, " ", "  ")
	// fmt.Printf("%s\n", b)

	jobs := orch.client.Jobs()
	q := &nomad.WriteOptions{
		Region: opt.Region,
	}

	if opt.Dryrun {
		planOpts := &nomad.PlanOptions{Diff: true}
		resp, _, err = jobs.PlanOpts(job, planOpts, q)
	} else {
		regOpts := &nomad.RegisterOptions{}
		resp, _, err = jobs.RegisterOpts(job, regOpts, q)
	}
	return
}

// func (orch *nomadOrchestrator) Status(ctx context.Context, stack *thrapb.Stack) []*thrapb.CompStatus {
// 	nomadJobs := orch.client.Jobs()
// 	qopts := &nomad.QueryOptions{AllowStale: false}

// 	// summary, _, err := nomadJobs.Summary(stack.ID, qopts)
// 	// summary.Summary
// 	taskGroups, _, err := nomadJobs.Allocations(stack.ID, true, qopts)
// 	if err != nil {

// 	}

// 	groups := make(map[string][]*thrapb.CompStatus)

// 	for _, group := range taskGroups {
// 		// thrapb.CompStatus{
// 		// 	ID:
// 		// }
// 		i := strings.LastIndex(group.Name, "[")
// 		key := group.Name[i+1 : len(group.Name)-1]

// 		if _, ok := groups[key]; !ok {
// 			groups[key] = make([]*thrapb.CompStatus, 0)
// 		}

// 		fmt.Println(group.TaskGroup, group.Name)
// 		for k, v := range group.TaskStates {
// 			// fmt.Println(" ", k, v.State)
// 			status := &thrapb.CompStatus{
// 				ID:     k,
// 				Status: v.State,
// 			}
// 			groups[key] = append(groups[key], status)
// 		}
// 	}

// 	for k, stati := range groups {
// 		fmt.Println(k)
// 		for _, s := range stati {
// 			fmt.Println(" ", s)
// 		}
// 		fmt.Println()
// 	}

// 	return nil
// }

// func (orch *nomadOrchestrator) Destroy(ctx context.Context, stack *thrapb.Stack, opts RequestOptions) error {
// 	jobs := orch.client.Jobs()

// 	writeOpt := &nomad.WriteOptions{
// 		Region: opts.Region,
// 	}

// 	_, _, err := jobs.Deregister(stack.ID, opts.Purge, writeOpt)

// 	return err
// }

// default service block added to each task
// func serviceBlock(instance, name, comp, portLabel string) *nomad.Service {
// 	prefix := comp + "." + instance
// 	return &nomad.Service{
// 		Id:        prefix + "." + name,
// 		Name:      name,
// 		PortLabel: portLabel,
// 		Tags:      []string{prefix},
// 	}
// }
