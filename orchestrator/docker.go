package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/euforia/thrap/crt"
	"github.com/euforia/thrap/thrapb"
)

type DockerOrchestrator struct {
	crt *crt.Docker
}

// ID returns the orchestrator id
func (orch *DockerOrchestrator) ID() string {
	return "docker"
}

// Init initializes the docker orchestrator
func (orch *DockerOrchestrator) Init(conf map[string]interface{}) (err error) {
	orch.crt, err = crt.NewDocker()
	return
}

// DeployComponent deploys a single component
func (orch *DockerOrchestrator) DeployComponent(ctx context.Context, stackID string, comp *thrapb.Component, opts RequestOptions) error {
	return orch.startContainer(ctx, stackID, comp)
}

// Deploy deploys the whole stack in the appropriate order
func (orch *DockerOrchestrator) Deploy(ctx context.Context, stack *thrapb.Stack, opts RequestOptions) (resp interface{}, job interface{}, err error) {
	// Create an isolated network for all running containers
	err = orch.crt.CreateNetwork(ctx, stack.ID)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			orch.Destroy(ctx, stack)
		}
	}()

	// Deploy services like db's etc
	err = orch.startServices(ctx, stack)
	if err != nil {
		return
	}

	// Deploy non-head containers
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		if comp.Head {
			continue
		}

		err = orch.startContainer(ctx, stack.ID, comp)
		if err != nil {
			return
		}

		fmt.Println(comp.ID)
	}

	// Start head containers
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}
		if !comp.Head {
			continue
		}

		err = orch.startContainer(ctx, stack.ID, comp)
		if err != nil {
			break
		}

		fmt.Println(comp.ID)
	}

	return
}

// Destroy removes call components of the stack from the container runtime
func (orch *DockerOrchestrator) Destroy(ctx context.Context, stack *thrapb.Stack) []*thrapb.ActionReport {
	ar := make([]*thrapb.ActionReport, 0, len(stack.Components))

	for _, c := range stack.Components {
		r := &thrapb.ActionReport{Action: thrapb.NewAction("destroy", "comp", c.ID)}
		r.Error = orch.crt.Remove(ctx, c.ID+"."+stack.ID)
		ar = append(ar, r)
	}
	return ar
}

func (orch *DockerOrchestrator) startContainer(ctx context.Context, sid string, comp *thrapb.Component) error {
	cfg := thrapb.NewContainer(sid, comp.ID)

	if comp.IsBuildable() {
		cfg.Container.Image = filepath.Join(sid, comp.Name)
	} else {
		cfg.Container.Image = comp.Name
	}

	// Add image version if present
	if len(comp.Version) > 0 {
		cfg.Container.Image += ":" + comp.Version
	}

	if comp.HasEnvVars() {
		cfg.Container.Env = make([]string, 0, len(comp.Env.Vars))
		for k, v := range comp.Env.Vars {
			cfg.Container.Env = append(cfg.Container.Env, k+"="+v)
		}
	}

	// Publish all ports for a head component.
	// TODO: May need to map this to user defined host ports
	if comp.Head {
		cfg.Host.PublishAllPorts = true
	}

	// Non-blocking
	warnings, err := orch.crt.Run(ctx, cfg)
	if err != nil {
		return err
	}

	if len(warnings) > 0 {
		for _, w := range warnings {
			fmt.Printf("%s: %s\n", cfg.Name, w)
		}
	}

	// May need this to get proper state
	//<-time.After(200*time.Millisecond)
	var cstate types.ContainerJSON
	cstate, err = orch.crt.Inspect(ctx, cfg.Name)
	if err == nil {
		if cstate.State.Dead {
			if cstate.State.ExitCode != 0 {
				err = errors.New(cstate.State.Error)
			}
		}
	}

	return err
}

// startServices starts services starts all non-build components
func (orch *DockerOrchestrator) startServices(ctx context.Context, stack *thrapb.Stack) error {
	var err error

	fmt.Printf("\nServices:\n\n")

	for _, comp := range stack.Components {
		if comp.IsBuildable() {
			continue
		}

		// Pull image if we do not locally have it
		imageID := comp.Name + ":" + comp.Version
		if !orch.crt.HaveImage(ctx, imageID) {
			err = orch.crt.ImagePull(ctx, imageID)
			if err != nil {
				break
			}
		}

		if err = orch.startContainer(ctx, stack.ID, comp); err != nil {
			break
		}

		fmt.Println(comp.ID)

	}

	return err
}
