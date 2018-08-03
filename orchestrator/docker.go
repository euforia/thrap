package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	fmt.Printf("\nApplication:\n\n")

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

		fmt.Printf(" - %s:%s\n", comp.ID, comp.Version)
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

		fmt.Printf(" - %s:%s\n", comp.ID, comp.Version)
	}
	fmt.Println()

	return
}

// Status returns a CompStatus slice containing the status of each component
// in the stack
func (orch *DockerOrchestrator) Status(ctx context.Context, stack *thrapb.Stack) []*thrapb.CompStatus {
	out := make([]*thrapb.CompStatus, 0, len(stack.Components))
	for _, comp := range stack.Components {
		id := comp.ID + "." + stack.ID
		ss := orch.getCompStatus(ctx, id)
		ss.ID = comp.ID

		out = append(out, ss)
	}

	return out
}

func (orch *DockerOrchestrator) getCompStatus(ctx context.Context, id string) *thrapb.CompStatus {

	ss := &thrapb.CompStatus{}
	ss.Details, ss.Error = orch.crt.Inspect(ctx, id)

	if ss.Error == nil {
		if ss.Details.State.Status == "exited" {
			s := ss.Details.State
			ss.Error = fmt.Errorf("code=%d", s.ExitCode)
		}

	} else {
		ss.Details = types.ContainerJSON{
			ContainerJSONBase: &types.ContainerJSONBase{
				State: &types.ContainerState{Status: "failed"},
			},
			Config: &container.Config{},
		}
	}

	return ss
}

// Destroy removes call components of the stack from the container runtime
func (orch *DockerOrchestrator) Destroy(ctx context.Context, stack *thrapb.Stack) []*thrapb.ActionResult {
	ar := make([]*thrapb.ActionResult, 0, len(stack.Components))

	for _, c := range stack.Components {
		r := &thrapb.ActionResult{
			Action:   "destroy",
			Resource: c.ID,
			Error:    orch.crt.Remove(ctx, c.ID+"."+stack.ID),
		}
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

		fmt.Printf(" - %s:%s\n", comp.ID, comp.Version)
	}

	return err
}
