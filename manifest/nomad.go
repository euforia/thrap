package manifest

import (
	"fmt"
	"io"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/nomad/api"
)

const (
	defaultPriority   = 50
	defaultRegion     = "us-west-2"
	defaultGroupCount = 3
	defaultNetMbits   = 1
	defaultCPUMHz     = 200 //mhz
	defaultMemMB      = 256 //mb
)

func MakeNomadJob(stack *thrapb.Stack) (*api.Job, error) {
	id := stack.ID
	job := api.NewServiceJob(id, stack.Name, defaultRegion, defaultPriority)
	for _, dc := range []string{defaultRegion} {
		job = job.AddDatacenter(dc)
	}

	gid := "0"
	grp := api.NewTaskGroup(id+"."+gid, defaultGroupCount)
	// grp.SetMeta(key, val)
	// grp = grp.Constrain(c)

	grp.RestartPolicy = &api.RestartPolicy{}
	// grp.RestartPolicy.Merge(restartPolicy)

	grp.Update = api.DefaultUpdateStrategy()
	//grp.Update.Merge(other)

	grp.ReschedulePolicy = api.NewDefaultReschedulePolicy(api.JobTypeService)
	// grp.ReschedulePolicy.Merge(reschedPolicy)

	// var (
	//     enclave string
	// )

	for _, comp := range stack.Components {
		task := makeNomadTaskDocker(id, gid, comp)

		resources := makeResources(defaultCPUMHz, defaultMemMB, defaultNetMbits)

		if comp.External {
			//api.NewConstraint("${meta.hood}", operand, right)
		}

		if comp.Head {
			service := task.Services[0]
			// service.Name = comp.ID + "." + id
			service.Name = stack.ID
			service.Tags = []string{comp.ID}
			service.Checks = []api.ServiceCheck{
				defaultServiceCheck(),
			}
			resources.Networks[0].DynamicPorts = []api.Port{
				api.Port{Label: "default"},
			}
			// service.CheckRestart = &api.CheckRestart{Limit: 15}
		}
		// task.Constrain(c)

		//resources.Merge(other)
		task.Require(resources)

		grp = grp.AddTask(task)

	}

	job = job.AddTaskGroup(grp)
	return job, nil
}

func WriteNomadJob(job *api.Job, w io.Writer) error {
	wrappedJob := hclWrapNomadJob(job)
	b, err := hclencoder.Encode(wrappedJob)
	if err == nil {
		_, err = w.Write(append(b, []byte("\n")...))
	}
	return err
}

func makeNomadTaskDocker(sid, gid string, comp *thrapb.Component) *api.Task {
	cid := sid + "." + gid + "." + comp.ID
	portLabel := "default"
	task := api.NewTask(cid, "docker")

	task.SetConfig("image", fmt.Sprintf("%s:%s", comp.Name, comp.Version))

	task.SetConfig("labels", []map[string]interface{}{
		map[string]interface{}{
			"stack":     sid,
			"group":     gid,
			"component": comp.ID,
		},
	})

	task.SetConfig("logging", []map[string]interface{}{
		map[string]interface{}{
			"type": "syslog",
			"config": []map[string]interface{}{
				map[string]interface{}{
					"tag": cid,
				},
			},
		},
	})

	task.SetConfig("port_map", []map[string]interface{}{
		map[string]interface{}{
			portLabel: 10000,
		},
	})

	task.Services = []*api.Service{
		&api.Service{
			// Name:      sid,
			// Tags:      []string{comp.ID},
			PortLabel: portLabel,
			// 		Checks: []api.ServiceCheck{
			// 			api.ServiceCheck{
			// 				Path:     "/",
			// 				Method:   "GET",
			// 				Interval: 25e9,
			// 				Timeout:  3e9,
			// 			},
			// 		},
		},
	}

	return task
}

func makeNomadBatchJob(id, name string) *api.Job {
	job := api.NewBatchJob(id, name, defaultRegion, defaultPriority)

	return job
}

func defaultServiceCheck() api.ServiceCheck {
	return api.ServiceCheck{
		Type:     "http",
		Path:     "/",
		Method:   "GET",
		Timeout:  3e9,  // 3s
		Interval: 20e9, // 20s
	}
}

func makeResources(cpu, mem, mbits int) *api.Resources {
	// cpu := defaultCPUMHz
	// mem := defaultMemMB
	// mbits := defaultNetMbits

	return &api.Resources{
		CPU:      &cpu,
		MemoryMB: &mem,
		Networks: []*api.NetworkResource{
			&api.NetworkResource{
				// DynamicPorts: []api.Port{api.Port{Label: portLabel}},
				MBits: &mbits},
		},
	}

	// if portLabel != "" {
	// 	rsrc.Networks[0].DynamicPorts = []api.Port{api.Port{Label: portLabel}}
	// }

	// return rsrc
}

func hclWrapNomadJob(job *api.Job) map[string]interface{} {
	key := `job "` + *job.ID + `"`
	return map[string]interface{}{
		key: job,
	}
}
