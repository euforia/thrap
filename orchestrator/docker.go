package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/thrapb"
	"golang.org/x/net/context"
)

// DockerOrchestrator implements a docker backed orchestrator
type DockerOrchestrator struct {
	cli *client.Client
}

// ID returns the docker orchestrator id
func (orch *DockerOrchestrator) ID() string {
	return "docker"
}

// Init initializes the docker orchestrator.  It satisfies the orchestrator
// interface
func (orch *DockerOrchestrator) Init(config map[string]interface{}) error {
	os.Setenv("DOCKER_API_VERSION", "1.37")
	cli, err := client.NewEnvClient()
	if err == nil {
		orch.cli = cli
	}
	return err
}

func (orch *DockerOrchestrator) containerCreateStart(ctx context.Context, cfg *thrapb.Container, opts RequestOptions) error {
	resp, err := orch.cli.ContainerCreate(ctx, cfg.Container, cfg.Host, cfg.Network, cfg.Name)
	if err != nil {
		return err
	}

	logline := cfg.Name + ": created\n"
	// Append warnings
	if len(resp.Warnings) > 0 {
		for _, w := range resp.Warnings {
			logline += "  " + cfg.Name + ": " + w + "\n"
		}
	}
	// Write log line
	opts.Output.Write([]byte(logline))

	err = orch.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err == nil {
		fmt.Fprintf(opts.Output, "%s: started\n", cfg.Name)
	}

	return err
}

// setupNetwork sets up a user-defined bridge network
func (orch *DockerOrchestrator) setupNetwork(ctx context.Context, netID string) error {
	// Check if it already exists
	_, err := orch.cli.NetworkInspect(ctx, netID, types.NetworkInspectOptions{})
	if err != nil {
		// Create if doesn't exist
		_, err = orch.cli.NetworkCreate(ctx, netID, types.NetworkCreate{})
	}

	return err
}

// Deploy deploys the stack. Each non-buildable stack component is deployed first. Then
// each buildable component is built
// TODO: add dryrun support
func (orch *DockerOrchestrator) Deploy(stack *thrapb.Stack, opts RequestOptions) (interface{}, interface{}, error) {
	var (
		ctx = context.Background()
		err = orch.setupNetwork(ctx, stack.ID)
	)

	if err != nil {
		return nil, nil, err
	}

	// Deploy non-buildable components
	for _, comp := range stack.Components {

		cfg := thrapb.NewContainer(stack.ID, comp.ID)
		if comp.IsBuildable() {
			// Set fully qualified name including stack id for components
			// being built
			cfg.Container.Image = filepath.Join(stack.ID, comp.Name) + ":" + comp.Version
		} else {
			// Set user provided name otherwise
			cfg.Container.Image = comp.Name + ":" + comp.Version
		}

		err = orch.containerCreateStart(ctx, cfg, opts)
		if err != nil {
			break
		}

	}

	return nil, nil, err
}

// tearDown removes all non-build containers as cleanup.
// TODO: refactor
func (orch *DockerOrchestrator) tearDown(ctx context.Context, st *thrapb.Stack) {

	opts := types.ContainerRemoveOptions{Force: true}

	for _, c := range st.Components {
		if c.IsBuildable() {
			continue
		}

		if err := orch.cli.ContainerRemove(ctx, c.ID+"."+st.ID, opts); err != nil {
			fmt.Println(err)
		}
	}
}

// Build builds a single component of a stack using 'docker build'
func (orch *DockerOrchestrator) Build(stackID string, comp *thrapb.Component, opts RequestOptions) error {
	bc := comp.Build

	tarOpt := &archive.TarOptions{}

	ign, err := dockerfile.ReadIgnoresFile(bc.Context)
	if err == nil {
		tarOpt.ExcludePatterns = ign
	}

	rdc, err := archive.TarWithOptions(bc.Context, tarOpt)
	if err != nil {
		return err
	}
	defer rdc.Close()

	name := filepath.Join(stackID, comp.ID)
	opt := types.ImageBuildOptions{
		Tags:        []string{name},
		BuildID:     comp.ID,
		Dockerfile:  bc.Dockerfile,
		NetworkMode: stackID,
		// BuildArgs:  make(map[string]*string),
		// Remove:      true,
	}

	ctx := context.Background()
	resp, err := orch.cli.ImageBuild(ctx, rdc, opt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, opts.Output, 100, true, nil)
	if err == nil {
		fmt.Fprintf(opts.Output, "%s: built\n", name)
	}
	return err
}
