package manifest

import (
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
	defaultPortLabel  = "default"
)

// MakeNomadJob returns a nomad job from the stack
func MakeNomadJob(stack *thrapb.Stack) (*api.Job, error) {
	id := stack.ID
	job := api.NewServiceJob(id, stack.Name, defaultRegion, defaultPriority)
	for _, dc := range []string{defaultRegion} {
		job = job.AddDatacenter(dc)
	}

	// By default everything goes in 1 group
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

	for _, comp := range stack.Components {
		task := makeNomadTaskDocker(id, gid, comp)

		grp = grp.AddTask(task)

	}

	job = job.AddTaskGroup(grp)
	return job, nil
}

func makeNomadTaskDocker(sid, gid string, comp *thrapb.Component) *api.Task {
	cid := sid + "." + gid + "." + comp.ID
	task := api.NewTask(cid, "docker")

	task.SetConfig("image", comp.Name+":"+comp.Version)
	task.SetConfig("labels", []map[string]interface{}{
		map[string]interface{}{
			"nomad.taskgroup": gid,
			"stack":           sid,
			"component":       comp.ID,
		},
	})

	// task.SetConfig("logging", []map[string]interface{}{
	// 	map[string]interface{}{
	// 		"type": "syslog",
	// 		"config": []map[string]interface{}{
	// 			map[string]interface{}{
	// 				"tag": cid,
	// 			},
	// 		},
	// 	},
	// })

	task.Services = make([]*api.Service, 0, len(comp.Ports))
	portmap := make(map[string]interface{}, len(comp.Ports))
	netPorts := make([]api.Port, 0, len(comp.Ports))
	for k, v := range comp.Ports {
		portmap[k] = v

		task.Services = append(task.Services, &api.Service{
			Id:        sid + "." + k,
			Name:      sid,
			Tags:      []string{k + "." + comp.ID},
			PortLabel: k,
			// Checks: []api.ServiceCheck{
			// 	api.ServiceCheck{
			// 		Path:     "/",
			// 		Method:   "GET",
			// 		Interval: 25e9,
			// 		Timeout:  3e9,
			// 	},
			// },
		})

		netPorts = append(netPorts, api.Port{Label: k})
	}

	task.SetConfig("port_map", []map[string]interface{}{portmap})

	resources := makeResources(defaultCPUMHz, defaultMemMB, defaultNetMbits)
	resources.Networks[0].DynamicPorts = netPorts
	task.Require(resources)

	return task
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

func makeResources(icpu, imem, imbits int) *api.Resources {
	cpu := icpu
	mem := imem
	mbits := imbits

	return &api.Resources{
		CPU:      &cpu,
		MemoryMB: &mem,
		Networks: []*api.NetworkResource{
			&api.NetworkResource{
				// DynamicPorts: []api.Port{api.Port{Label: portLabel}},
				MBits: &mbits,
			},
		},
	}

}

// WriteNomadJob write writes the nomad job in hcl format to the writer w
func WriteNomadJob(job *api.Job, w io.Writer) error {
	wrappedJob := hclWrapNomadJob(job)
	b, err := hclencoder.Encode(wrappedJob)
	if err == nil {
		_, err = w.Write(append(b, []byte("\n")...))
	}
	return err
}

func hclWrapNomadJob(job *api.Job) map[string]interface{} {
	key := `job "` + *job.ID + `"`
	return map[string]interface{}{
		key: job,
	}
}
