package manifest

import (
	"fmt"

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

func makeNomadTaskDocker(sid, gid string, comp *thrapb.Component) *api.Task {
	cid := sid + "." + gid + "." + comp.ID
	portLabel := "default"
	task := api.NewTask(cid, "docker")

	task.SetConfig("image", fmt.Sprintf("%s:%s", comp.Name, comp.Version))
	task.SetConfig("port_map", map[string]int{
		portLabel: -1,
	})
	task.SetConfig("labels", map[string]string{
		"stack":     sid,
		"group":     gid,
		"component": comp.ID,
	})
	task.SetConfig("logging", map[string]interface{}{
		"type": "syslog",
		"config": map[string]interface{}{
			"tag": cid,
		},
	})

	task.Services = []*api.Service{
		&api.Service{
			Name:      comp.ID + "." + sid,
			PortLabel: portLabel,
			// Checks: []api.ServiceCheck{
			// 	api.ServiceCheck{
			// 		Path:     "/v1/status",
			// 		Method:   "GET",
			// 		Interval: 25e9,
			// 		Timeout:  3e9,
			// 	},
			// },
		},
	}

	return task
}

func makeNomadBatchJob(id, name string) *api.Job {
	job := api.NewBatchJob(id, name, defaultRegion, defaultPriority)

	return job
}

func makeNomadJob(stack *thrapb.Stack) (*api.Job, error) {
	id := stack.ID
	job := api.NewServiceJob(id, stack.Name, defaultRegion, defaultPriority)
	for _, dc := range []string{"test-dc0", "test-dc1"} {
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

		if comp.External {
			//api.NewConstraint("${meta.hood}", operand, right)
		}

		if comp.Head {
			service := task.Services[0]
			service.Name = comp.ID + "." + id
			service.Tags = []string{}
			service.Checks = []api.ServiceCheck{
				defaultServiceCheck(),
			}
			// service.CheckRestart = &api.CheckRestart{Limit: 15}
		}
		// task.Constrain(c)

		resources := defaultResources()
		//resources.Merge(other)
		task.Require(resources)

		grp = grp.AddTask(task)

	}

	job = job.AddTaskGroup(grp)
	return job, nil
}

func defaultServiceCheck() api.ServiceCheck {
	return api.ServiceCheck{
		Type:     "http",
		Timeout:  3e9,  // 3s
		Interval: 20e9, // 20s
	}
}

func defaultResources() *api.Resources {
	cpu := defaultCPUMHz
	mem := defaultMemMB
	mbits := defaultNetMbits

	return &api.Resources{
		CPU:      &cpu,
		MemoryMB: &mem,
		Networks: []*api.NetworkResource{
			&api.NetworkResource{MBits: &mbits},
		},
	}
}

func hclWrapNomadJob(job *api.Job) map[string]interface{} {
	key := `job "` + *job.ID + `"`
	return map[string]interface{}{
		key: job,
	}
}
