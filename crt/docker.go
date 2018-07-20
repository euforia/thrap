package crt

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/thrapb"
	"golang.org/x/net/context"
)

// RequestOptions are the common request options for crt's
// type RequestOptions struct {
// 	Output io.Writer
// }

type BuildRequest struct {
	ContextDir string
	TarOpts    *archive.TarOptions
	BuildOpts  *types.ImageBuildOptions
	Output     io.Writer
}

// PushRequest is a container image push request
type PushRequest struct {
	Image   string
	Tag     string
	Output  io.Writer
	Options types.ImagePushOptions
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

// Stop stops a running container
func (orch *Docker) Stop(ctx context.Context, containerID string) error {
	dur := 3 * time.Second
	return orch.cli.ContainerStop(ctx, containerID, &dur)
}

// HaveImage returns true if we locally have the image
func (orch *Docker) HaveImage(ctx context.Context, imageID string) bool {
	_, _, err := orch.cli.ImageInspectWithRaw(ctx, imageID)
	return err == nil
}

// Run creates and runs a container with the given config
func (orch *Docker) Run(ctx context.Context, cfg *thrapb.Container) ([]string, error) {

	resp, err := orch.cli.ContainerCreate(ctx, cfg.Container, cfg.Host, cfg.Network, cfg.Name)
	if err != nil {
		return nil, err
	}

	err = orch.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	return resp.Warnings, err
}

// Inspect returns information about the container by id
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

// ListImagesWithLabel returns a list of images that match the given label
func (orch *Docker) ListImagesWithLabel(ctx context.Context, label string) ([]types.ImageSummary, error) {
	args := filters.NewArgs(filters.Arg("label", label))

	// kvp := make([]filters.KeyValuePair, 0, len(labels))
	// for k, v := range labels {
	// 	kvp = append(kvp, filters.KeyValuePair{Key: k, Value: v})
	// }
	// args := filters.NewArgs(kvp...)

	opts := types.ImageListOptions{Filters: args}
	// fmt.Println(opts)
	return orch.cli.ImageList(ctx, opts)
}

// Remove forcibly stops and removes a container
func (orch *Docker) Remove(ctx context.Context, cid string) error {
	opts := types.ContainerRemoveOptions{Force: true}
	return orch.cli.ContainerRemove(ctx, cid, opts)
}

// Logs returns logs for a single container
func (orch *Docker) Logs(ctx context.Context, containerID string, stdout, stderr io.Writer) error {
	opts := types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	}

	clogs, err := orch.cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return err
	}
	defer clogs.Close()

	_, err = stdcopy.StdCopy(stdout, stderr, clogs)
	return err
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

// ImagePush pushes an image using the local docker engine to the remote registry
func (orch *Docker) ImagePush(ctx context.Context, req *PushRequest) error {
	rc, err := orch.cli.ImagePush(ctx, req.Image+":"+req.Tag, req.Options)
	if err != nil {
		return err
	}

	defer rc.Close()

	return jsonmessage.DisplayJSONMessagesStream(rc, req.Output, 101, true, nil)
}
