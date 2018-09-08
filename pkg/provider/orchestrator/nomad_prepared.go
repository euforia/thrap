package orchestrator

import (
	"encoding/json"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/client/driver/env"
	hargs "github.com/hashicorp/nomad/helper/args"
)

const (
	thrapServiceTagVar  = "${TASK}.${NOMAD_META_INSTANCE}"
	thrapServiceNameVar = "${NOMAD_META_PROJECT}"
)

type nomadPreparedDeployment struct {
	job *nomad.Job

	// Serialized job.  This is the final job and cannot be used
	// for a resubmit
	serialized []byte
}

func newNomadPreparedDeployment(job *nomad.Job) (*nomadPreparedDeployment, error) {
	var (
		n   = &nomadPreparedDeployment{job: job}
		err error
	)

	n.addServiceDefs()

	n.serialized, err = json.MarshalIndent(job, "", "  ")
	return n, err
}

// Bytes returns the bytes for a given deployment after all normalization
// has occured.  This is the spec that was deployed.
func (d *nomadPreparedDeployment) Bytes() []byte {
	return d.serialized
}

func (d *nomadPreparedDeployment) Spec() interface{} {
	return d.job
}

func (d *nomadPreparedDeployment) Artifacts() []string {
	// Get image names from each task
	out := make([]string, 0)
	for _, g := range d.job.TaskGroups {
		for _, t := range g.Tasks {
			if i, ok := t.Config["image"]; ok {
				image := i.(string)
				out = append(out, image)
			}
		}
	}

	// Eval sub'ing variables
	return d.evalImageNames(out)
}

func (d *nomadPreparedDeployment) evalImageNames(imgs []string) []string {
	meta := map[string]string{}
	for k, v := range d.job.Meta {
		meta[env.MetaPrefix+k] = v
	}

	replaced := make([]string, len(imgs))
	for i, img := range imgs {
		replaced[i] = hargs.ReplaceEnv(img, meta)
	}

	return replaced
}

// Adds the thrap service variables.
func (d *nomadPreparedDeployment) addServiceDefs() {
	for _, g := range d.job.TaskGroups {
		for _, t := range g.Tasks {
			addServiceDefs(t)
		}
	}
}

// Adds the thrap service variables
func addServiceDefs(task *nomad.Task) {
	services := task.Services
	add := make([]*nomad.Service, 0)

	for _, service := range services {
		if service.Name == thrapServiceNameVar {
			if !hasServiceTag(service.Tags) {
				service.Tags = append(service.Tags, thrapServiceTagVar)
			}
		} else {
			// Add an additional one
			add = append(add, &nomad.Service{
				Name:      thrapServiceNameVar,
				Tags:      []string{thrapServiceTagVar},
				PortLabel: service.PortLabel,
			})
		}
	}

	task.Services = append(task.Services, add...)
}

func hasServiceTag(tags []string) bool {
	for _, t := range tags {
		if t == thrapServiceTagVar {
			return true
		}
	}
	return false
}