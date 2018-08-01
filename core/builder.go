package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/docker/docker/api/types"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/crt"
	"github.com/euforia/thrap/metrics"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vars"
	"github.com/pkg/errors"
)

// CompBuildResult is the result of a component build
type CompBuildResult struct {
	// Labels applied to the build
	Labels map[string]string
	// Total build runtime
	Runtime *metrics.Runtime
	// Error signifying an unsuccessful build
	Error error
	// Log of the complete build used to verify the build.
	// This should eventually be an interface
	Log *crt.DockerBuildLog
	// Whether the image was published or not
	Published bool
}

// HasError returns true if the build result contains an error
func (result *CompBuildResult) HasError() bool {
	return result.Error != nil
}

type stackBuilder struct {
	conf *config.RegistryConfig
	// builder
	crt *crt.Docker

	// build and deploy common functions
	run *bdCommon

	// Total run time
	totalTime *metrics.Runtime
	// Total build time for all components
	buildTime *metrics.Runtime
	// Time to spin up dependent services
	svcTime *metrics.Runtime

	// Build result per component
	results map[string]*CompBuildResult

	// stack being built
	stack *thrapb.Stack
}

func newStackBuilder(c *crt.Docker, conf *config.RegistryConfig, stack *thrapb.Stack) *stackBuilder {
	return &stackBuilder{
		conf:      conf,
		crt:       c,
		run:       &bdCommon{c},
		totalTime: &metrics.Runtime{},
		buildTime: &metrics.Runtime{},
		results:   make(map[string]*CompBuildResult, len(stack.Components)),
		stack:     stack,
	}
}

func (bldr *stackBuilder) printResults(w io.Writer) {
	w.Write([]byte("\n"))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, "Component\tArtifact\tStatus\tDetails\n")
	fmt.Fprintf(tw, "---------\t--------\t------\t-------\n")
	for k, r := range bldr.results {
		var (
			status string
			msg    string
			art    string
		)
		if r.Error == nil {
			comp := bldr.stack.Components[k]
			status = "success"
			art = comp.Name + ":" + comp.Version
		} else {
			status = "fail"
			msg = r.Error.Error()
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", k, art, status, msg)
	}
	tw.Flush()
	w.Write([]byte("\n"))
}

func (bldr *stackBuilder) build(ctx context.Context) error {
	bldr.totalTime.Start()
	defer bldr.totalTime.End()

	err := bldr.crt.CreateNetwork(ctx, bldr.stack.ID)
	if err != nil {
		return err
	}

	defer bldr.run.destroy(ctx, bldr.stack)

	// Start containers needed for build
	bldr.svcTime, err = bldr.run.startServices(ctx, bldr.stack)
	if err != nil {
		return err
	}

	// Start non-head builds
	bldr.buildTime.Start()
	defer bldr.buildTime.End()

	err = bldr.buildByHead(ctx, false)
	if err != nil {
		return err
	}

	// Start head builds
	err = bldr.buildByHead(ctx, true)
	return err
}

// buildByHead builds based on the head flag supplied
func (bldr *stackBuilder) buildByHead(ctx context.Context, head bool) error {
	var (
		stack = bldr.stack
		err   error
	)

	// Start build containers after
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		// Build based on whether head was requested
		if comp.Head != head {
			continue
		}

		bldr.doBuild(ctx, comp)

		// Start container from image that was just built, if this component
		// is not the head
		if !comp.Head {
			err = bldr.run.startContainer(ctx, stack.ID, comp)
			if err != nil {
				break
			}
		}
	}

	return err
}

func (bldr *stackBuilder) doBuild(ctx context.Context, comp *thrapb.Component) {
	result := &CompBuildResult{
		Runtime: (&metrics.Runtime{}).Start(),
		Log:     crt.NewDockerBuildLog(os.Stdout),
	}

	fmt.Printf("\nBuilding %s:\n\n", comp.ID)

	req := bldr.makeBuildRequest(comp, result.Log)

	// Blocking
	result.Error = bldr.crt.Build(ctx, req)
	result.Runtime.End()
	result.Labels = req.BuildOpts.Labels

	// Add result
	bldr.results[comp.ID] = result
}

func (bldr *stackBuilder) makeBuildRequest(comp *thrapb.Component, output io.Writer) *crt.BuildRequest {
	req := &crt.BuildRequest{
		// Output:     crt.NewDockerBuildLog(os.Stdout),
		Output:     output,
		ContextDir: comp.Build.Context,
		BuildOpts: &types.ImageBuildOptions{
			Tags: getBuildImageTags(bldr.stack.ID, comp, bldr.conf),
			// ID to use in order to cancel the build
			// BuildID:     comp.ID,
			Dockerfile:  comp.Build.Dockerfile,
			NetworkMode: bldr.stack.ID,
			// Add labels to query later
			Labels: map[string]string{
				"stack":               bldr.stack.ID,
				"component":           comp.ID,
				vars.ComponentVersion: comp.Version,
				vars.StackVersion:     bldr.stack.Version,
			},
		},
	}

	if comp.HasEnvVars() {
		args := make(map[string]*string, len(comp.Env.Vars))

		fmt.Printf("  Arguments:\n\n")
		for k := range comp.Env.Vars {
			fmt.Println("   -", k)

			v := comp.Env.Vars[k]
			args[k] = &v
		}
		fmt.Println()

		req.BuildOpts.BuildArgs = args
	}

	return req
}

func printScopeVars(scopeVars scope.Variables) {
	fmt.Printf("\nScope:\n\n")
	for _, name := range scopeVars.Names() {
		fmt.Println(" ", name)
	}
	fmt.Println()
}

// build and deploy common functions
type bdCommon struct {
	crt *crt.Docker
}

// getBuildImageTags returns tags that should be applied to a given image build
func getBuildImageTags(sid string, comp *thrapb.Component, rconf *config.RegistryConfig) []string {
	base := filepath.Join(sid, comp.ID)
	out := []string{base}
	if len(comp.Version) > 0 {
		out = append(out, base+":"+comp.Version)
	}

	// rconf := bldr.conf //.GetDefaultRegistry()
	if rconf != nil && len(rconf.Addr) > 0 {
		rbase := filepath.Join(rconf.Addr, base)
		out = append(out, rbase)
		if len(comp.Version) > 0 {
			out = append(out, rbase+":"+comp.Version)
		}
	}
	return out
}

// startServices starts services needed to perform the build that themselves do not need
// to be built
func (c *bdCommon) startServices(ctx context.Context, stack *thrapb.Stack) (*metrics.Runtime, error) {
	var (
		runtime = (&metrics.Runtime{}).Start()
		err     error
	)

	fmt.Printf("Services:\n\n")

	for _, comp := range stack.Components {
		if comp.IsBuildable() {
			continue
		}

		// Pull image if we do not locally have it
		imageID := comp.Name + ":" + comp.Version
		if !c.crt.HaveImage(ctx, imageID) {
			err = c.crt.ImagePull(ctx, imageID)
			if err != nil {
				break
			}
		}

		if err = c.startContainer(ctx, stack.ID, comp); err != nil {
			break
		}

		fmt.Println(" -", comp.ID)

	}

	runtime.End()

	return runtime, err
}

func (c *bdCommon) startContainer(ctx context.Context, sid string, comp *thrapb.Component) error {
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
	warnings, err := c.crt.Run(ctx, cfg)
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
	cstate, err = c.crt.Inspect(ctx, cfg.Name)
	if err == nil {
		if cstate.State.Dead {
			if cstate.State.ExitCode != 0 {
				err = errors.New(cstate.State.Error)
			}
		}
	}

	return err
}

// Destroy removes call components of the stack from the container runtime
func (c *bdCommon) destroy(ctx context.Context, stack *thrapb.Stack) []*thrapb.ActionReport {
	ar := make([]*thrapb.ActionReport, 0, len(stack.Components))

	for _, comp := range stack.Components {
		r := &thrapb.ActionReport{Action: thrapb.NewAction("destroy", "comp", comp.ID)}
		r.Error = c.crt.Remove(ctx, comp.ID+"."+stack.ID)
		ar = append(ar, r)
	}
	return ar
}
