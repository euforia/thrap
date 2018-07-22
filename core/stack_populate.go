package core

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/euforia/thrap/thrapb"
)

func (st *Stack) populateFromImageConf(stack *thrapb.Stack) {
	confs := st.getContainerConfigs(stack)
	st.populatePorts(stack, confs)
	st.populateVolumes(stack, confs)
}

// populatePorts populates ports into the stack from the container images for ports
// that have not been defined in the stack but are in the image config
func (st *Stack) populatePorts(stack *thrapb.Stack, contConfs map[string]*container.Config) {
	for id, cfg := range contConfs {
		comp := stack.Components[id]
		if comp.Ports == nil {
			comp.Ports = make(map[string]int32, len(cfg.ExposedPorts))
		}

		if len(cfg.ExposedPorts) == 1 {
			for k := range cfg.ExposedPorts {
				if !comp.HasPort(int32(k.Int())) {
					comp.Ports["default"] = int32(k.Int())
				}
				break
			}
		} else {
			for k := range cfg.ExposedPorts {
				if !comp.HasPort(int32(k.Int())) {
					// HCL does not allow numbers as keys
					comp.Ports["port"+k.Port()] = int32(k.Int())
				}
			}
		}
	}
}

func (st *Stack) populateVolumes(stack *thrapb.Stack, contConfs map[string]*container.Config) {
	for id, cfg := range contConfs {
		comp := stack.Components[id]
		vols := make([]*thrapb.Volume, 0, len(cfg.Volumes))

		for k := range cfg.Volumes {
			if !comp.HasVolumeTarget(k) {
				vols = append(vols, &thrapb.Volume{Target: k})
			}
		}
		comp.Volumes = append(comp.Volumes, vols...)
	}
}

func (st *Stack) getContainerConfigs(stack *thrapb.Stack) map[string]*container.Config {
	out := make(map[string]*container.Config, len(stack.Components))
	for _, comp := range stack.Components {
		// Ensure we have the image locally
		err := st.crt.ImagePull(context.Background(), comp.Name+":"+comp.Version)
		if err != nil {
			continue
		}
		// Get image config
		ic, err := st.crt.ImageConfig(comp.Name + ":" + comp.Version)
		if err != nil {
			continue
		}

		out[comp.ID] = ic
	}
	return out
}
