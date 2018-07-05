package crt

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/thrapb"
	"golang.org/x/net/context"
)

// RequestOptions are the common request options for crt's
type RequestOptions struct {
	Output io.Writer
}

type BuildRequest struct {
	ContextDir string
	TarOpts    *archive.TarOptions
	BuildOpts  *types.ImageBuildOptions
	Output     io.Writer
}

// Docker implements a docker backed orchestrator
type Docker struct {
	cli *client.Client
}

func NewDocker() (*Docker, error) {
	os.Setenv("DOCKER_API_VERSION", "1.37")
	cli, err := client.NewEnvClient()
	if err == nil {
		return &Docker{cli}, nil
	}
	return nil, err
}

func (orch *Docker) Stop(ctx context.Context, containerID string) error {
	dur := 3 * time.Second
	return orch.cli.ContainerStop(ctx, containerID, &dur)
}

// Run creates and runs a container with the given config
func (orch *Docker) Run(ctx context.Context, cfg *thrapb.Container) ([]string, error) {

	resp, err := orch.cli.ContainerCreate(ctx, cfg.Container, cfg.Host, cfg.Network, cfg.Name)
	if err != nil {
		return nil, err
	}

	//logline := cfg.Name + ": created\n"
	// Append warnings
	// if len(resp.Warnings) > 0 {
	// 	var logline string
	// 	for _, w := range resp.Warnings {
	// 		logline += "  " + cfg.Name + ": " + w + "\n"
	// 	}
	// 	opts.Output.Write([]byte(logline))
	// }
	// Write log line

	err = orch.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	// if err == nil {
	// 	fmt.Fprintf(opts.Output, "%s: started\n", cfg.Name)
	// }

	return resp.Warnings, err
}

func (orch *Docker) Inspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return orch.cli.ContainerInspect(ctx, containerID)
}

// CreateNetwork sets up a user-defined bridge network only if one does not
// exist by the given id
func (orch *Docker) CreateNetwork(ctx context.Context, netID string) error {
	// Check if it already exists
	_, err := orch.cli.NetworkInspect(ctx, netID, types.NetworkInspectOptions{})
	if err != nil {
		// Create if doesn't exist
		_, err = orch.cli.NetworkCreate(ctx, netID, types.NetworkCreate{})
	}

	return err
}

// Deploy deploys each component in the stack. The components must exist for deploy to succeed
// func (orch *Docker) Deploy(stack *thrapb.Stack, opts RequestOptions) (interface{}, interface{}, error) {
// 	var (
// 		ctx = context.Background()
// 		err = orch.CreateNetwork(ctx, stack.ID)
// 	)

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Deploy non-buildable components
// 	for _, comp := range stack.Components {
// 		err = orch.DeployComponent(ctx, stack.ID, comp, opts)
// 		if err != nil {
// 			break
// 		}

// 	}

// 	return nil, nil, err
// }

// func (orch *Docker) DeployComponent(ctx context.Context, cfg *thrapb.Container, opts RequestOptions) error {

// 	return orch.containerCreateStart(ctx, cfg, opts)
// }

// tearDown removes all non-build containers as cleanup.
// TODO: refactor
// func (orch *Docker) tearDown(ctx context.Context, st *thrapb.Stack) {

// 	opts := types.ContainerRemoveOptions{Force: true}

// 	for _, c := range st.Components {
// 		if c.IsBuildable() {
// 			continue
// 		}

// 		if err := orch.cli.ContainerRemove(ctx, c.ID+"."+st.ID, opts); err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }

// Remove forcibly stops and removes a container
func (orch *Docker) Remove(ctx context.Context, cid string) error {
	opts := types.ContainerRemoveOptions{Force: true}
	return orch.cli.ContainerRemove(ctx, cid, opts)
}

func (orch *Docker) Build(ctx context.Context, req *BuildRequest) error {
	ign, err := dockerfile.ReadIgnoresFile(req.ContextDir)
	if err != nil {
		return err
	}

	if req.TarOpts == nil {
		req.TarOpts = &archive.TarOptions{}
	}

	if len(req.TarOpts.ExcludePatterns) == 0 {
		req.TarOpts.ExcludePatterns = ign
	} else {
		return errors.New("merging ignores currently not supported")
	}
	// TODO: dedup ignores and append

	rdc, err := archive.TarWithOptions(req.ContextDir, req.TarOpts)
	if err != nil {
		return err
	}
	defer rdc.Close()

	resp, err := orch.cli.ImageBuild(ctx, rdc, *req.BuildOpts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// params: reader, output writer, descriptor, isTerminal, auxCallback
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, req.Output, uintptr(rand.Uint32()), true, nil)
	return err
}

// Build builds a single component of a stack using 'docker build'
// func (orch *Docker) BuildComponent(ctx context.Context, stackID string, comp *thrapb.Component, opts RequestOptions) error {
// 	bc := comp.Build

// 	tarOpt := &archive.TarOptions{}

// 	ign, err := dockerfile.ReadIgnoresFile(bc.Context)
// 	if err == nil {
// 		tarOpt.ExcludePatterns = ign
// 	}

// 	rdc, err := archive.TarWithOptions(bc.Context, tarOpt)
// 	if err != nil {
// 		return err
// 	}
// 	defer rdc.Close()

// 	name := filepath.Join(stackID, comp.ID)
// 	opt := types.ImageBuildOptions{
// 		Tags:        []string{name},
// 		BuildID:     comp.ID, // todo add vcs repo version
// 		Dockerfile:  bc.Dockerfile,
// 		NetworkMode: stackID,
// 		// Remove:      true,
// 	}

// 	if comp.HasEnvVars() {
// 		opt.BuildArgs = make(map[string]*string, len(comp.Env.Vars))
// 		for k := range comp.Env.Vars {
// 			v := comp.Env.Vars[k]
// 			opt.BuildArgs[k] = &v
// 		}
// 	}

// 	//fmt.Printf("%+v\n", opt)

// 	resp, err := orch.cli.ImageBuild(ctx, rdc, opt)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	// params: reader, output writer, descriptor, isTerminal, auxCallback
// 	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, opts.Output, uintptr(rand.Uint32()), true, nil)
// 	return err
// }

// ImagePull pulls in image from the docker registry using docker. This uses
// dockers built in mechanism to communicate to the registry
func (orch *Docker) ImagePull(ctx context.Context, ref string) error {
	var options types.ImagePullOptions
	rd, err := orch.cli.ImagePull(ctx, ref, options)
	if err != nil {
		return err
	}

	defer rd.Close()

	return jsonmessage.DisplayJSONMessagesStream(rd, os.Stdout, 100, true, nil)
}

// ImageConfig returns an image config for the given name and tagged image
func (orch *Docker) ImageConfig(name, tag string) (*container.Config, error) {
	inf, _, err := orch.cli.ImageInspectWithRaw(context.Background(), name+":"+tag)
	if err != nil {
		return nil, err
	}
	return inf.Config, nil
}
