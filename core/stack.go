package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"github.com/pkg/errors"

	"github.com/euforia/pseudo"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/crt"
	"github.com/euforia/thrap/metrics"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/vars"

	"github.com/euforia/thrap/asm"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/registry"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
)

var (
	// ErrStackAlreadyRegistered is used when a stack is already registered
	ErrStackAlreadyRegistered = errors.New("stack already registered")
	errComponentNotBuildable  = errors.New("component not buildable")
)

// BuildOptions are options to perform a build-publish
type BuildOptions struct {
	// Workdir is the root of a git repo.  This is used to decide
	// if we can auto-publish
	Workdir string
	// If true the build is published despite the auto-publish check,
	// essentially a force publish
	Publish bool
}

// StackStore implements a thrap stack store
type StackStore interface {
	Get(id string) (*thrapb.Stack, error)
	Iter(prefix string, f func(*thrapb.Stack) error) error
	Register(stack *thrapb.Stack) (*thrapb.Stack, []*thrapb.ActionReport, error)
}

// Stack provides various stack based operations
type Stack struct {
	// profile id the stack was loaded with
	profile string
	// builder is docker runtime
	crt *crt.Docker
	// config to use for this instance
	conf *config.ThrapConfig
	// there can only be one vcs provider
	vcs vcs.VCS
	// available registries
	regs map[string]registry.Registry
	// runtime orchestrator
	orch orchestrator.Orchestrator
	// packs
	packs *packs.Packs
	// stack store
	sst StackStorage

	log *log.Logger
}

// Assembler returns a new assembler for the stack
func (st *Stack) Assembler(cwd string, stack *thrapb.Stack) (*asm.StackAsm, error) {
	scopeVars := st.conf.VCS[st.vcs.ID()].ScopeVars("vcs.")
	return asm.NewStackAsm(stack, cwd, st.vcs, nil, scopeVars, st.packs)
}

// Register registers a new stack. It returns an error if the stack is
// already registered or fails to register
func (st *Stack) Register(stack *thrapb.Stack) (*thrapb.Stack, []*thrapb.ActionReport, error) {
	errs := stack.Validate()
	if len(errs) > 0 {
		return nil, nil, utils.FlattenErrors(errs)
	}

	stack, err := st.sst.Create(stack)
	if err != nil {
		if err == store.ErrStackExists {
			return nil, nil, ErrStackAlreadyRegistered
		}
		return nil, nil, err
	}

	reports := st.ensureStackResources(stack)

	// Temp
	for _, r := range reports {
		fmt.Printf("%v '%v'\n", r.Action, r.Error)
	}

	return stack, reports, err
}

// Validate validates the stack manifest
func (st *Stack) Validate(stack *thrapb.Stack) error {
	// stack.Version = vcs.GetRepoVersion(ctxDir).String()
	errs := stack.Validate()
	if len(errs) > 0 {
		return utils.FlattenErrors(errs)
	}
	return nil
}

// Commit updates a stack definition
func (st *Stack) Commit(stack *thrapb.Stack) (*thrapb.Stack, error) {
	errs := stack.Validate()
	if len(errs) > 0 {
		return nil, utils.FlattenErrors(errs)
	}

	return st.sst.Update(stack)
}

// Init initializes a basic stack with the configuration and options provided. This should only be
// used in the local cli case as the config is merged with the global.
func (st *Stack) Init(stconf *asm.BasicStackConfig, opt ConfigureOptions) (*thrapb.Stack, error) {

	_, err := ConfigureLocal(st.conf, opt)
	if err != nil {
		return nil, err
	}

	repo := opt.VCS.Repo
	vcsp, gitRepo, err := vcs.SetupLocalGitRepo(repo.Name, repo.Owner, opt.DataDir, opt.VCS.Addr)
	if err != nil {
		return nil, err
	}

	stack, err := asm.NewBasicStack(stconf, st.packs)
	if err != nil {
		return nil, err
	}
	if errs := stack.Validate(); len(errs) > 0 {
		return nil, utils.FlattenErrors(errs)
	}

	st.populateFromImageConf(stack)

	scopeVars := st.conf.VCS[st.vcs.ID()].ScopeVars("vcs.")
	stasm, err := asm.NewStackAsm(stack, opt.DataDir, vcsp, gitRepo, scopeVars, st.packs)
	if err != nil {
		return stack, err
	}

	err = stasm.AssembleMaterialize()
	if err == nil {
		err = stasm.WriteManifest()
	}

	return stack, err
}

func (st *Stack) scopeVars(stack *thrapb.Stack) scope.Variables {
	svars := stack.ScopeVars()
	for k, v := range stack.Components {

		ipvar := ast.Variable{
			Type:  ast.TypeString,
			Value: k + "." + stack.ID,
		}

		// Set container ip var
		svars[v.ScopeVarName(consts.CompVarPrefixKey+".", "container.ip")] = ipvar
		// Set container.addr var per port label
		for pl, p := range v.Ports {
			svars[v.ScopeVarName(consts.CompVarPrefixKey+".", "container.addr."+pl)] = ast.Variable{
				Type:  ast.TypeString,
				Value: fmt.Sprintf("%s:%d", ipvar.Value, p),
			}
		}

	}

	return svars
}

// startServices starts services needed to perform the build that themselves do not need
// to be built
func (st *Stack) startServices(ctx context.Context, stack *thrapb.Stack) error {
	var err error

	fmt.Printf("Services:\n\n")

	for _, comp := range stack.Components {
		if comp.IsBuildable() {
			continue
		}

		// Pull image if we do not locally have it
		imageID := comp.Name + ":" + comp.Version
		if !st.crt.HaveImage(ctx, imageID) {
			err = st.crt.ImagePull(ctx, imageID)
			if err != nil {
				break
			}
		}

		if err = st.startContainer(ctx, stack.ID, comp); err != nil {
			break
		}

		fmt.Println(" ", comp.ID)

	}

	return err
}

func (st *Stack) startContainer(ctx context.Context, sid string, comp *thrapb.Component) error {
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
	warnings, err := st.crt.Run(ctx, cfg)
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
	cstate, err = st.crt.Inspect(ctx, cfg.Name)
	if err == nil {
		if cstate.State.Dead {
			if cstate.State.ExitCode != 0 {
				err = errors.New(cstate.State.Error)
			}
		}
	}

	return err
}

// Log writes the log for a running component to the writers
func (st *Stack) Log(ctx context.Context, id string, stdout, stderr io.Writer) error {
	return st.crt.Logs(ctx, id, stdout, stderr)
}

// Logs writes all running component logs for the stack
func (st *Stack) Logs(ctx context.Context, stack *thrapb.Stack, stdout, stderr io.Writer) error {
	var err error
	for _, comp := range stack.Components {
		er := st.crt.Logs(ctx, comp.ID+"."+stack.ID, stdout, stderr)
		if er != nil {
			err = er
		}
	}

	return err
}

// Status returns a CompStatus slice containing the status of each component
// in the stack
func (st *Stack) Status(ctx context.Context, stack *thrapb.Stack) []*CompStatus {
	out := make([]*CompStatus, 0, len(stack.Components))
	for _, comp := range stack.Components {
		id := comp.ID + "." + stack.ID
		ss := st.getCompStatus(ctx, id)
		ss.ID = comp.ID

		out = append(out, ss)
	}

	return out
}

func (st *Stack) getCompStatus(ctx context.Context, id string) *CompStatus {

	ss := &CompStatus{}
	ss.Details, ss.Error = st.crt.Inspect(ctx, id)

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

// Artifacts returns all known artifacts for the stack
func (st *Stack) Artifacts(stack *thrapb.Stack) []*thrapb.Artifact {
	images := make([]*thrapb.Artifact, 0, len(stack.Components))

	ctx := context.Background()
	conts, err := st.crt.ListImagesWithLabel(ctx, "stack="+stack.ID)
	if err != nil {
		return images
	}

	for _, c := range conts {
		ci := thrapb.NewArtifact(c.ID, c.RepoTags)
		ci.Labels = c.Labels
		ci.Created = c.Created
		ci.DataSize = c.Size
		images = append(images, ci)
	}

	return images
}

// Get returns a stack from the store by id
func (st *Stack) Get(id string) (*thrapb.Stack, error) {
	return st.sst.Get(id)
}

// Iter iterates over each stack definition in the store.
func (st *Stack) Iter(prefix string, f func(*thrapb.Stack) error) error {
	return st.sst.Iter(prefix, f)
}

// Build starts all require services, then starts all the builds
func (st *Stack) Build(ctx context.Context, stack *thrapb.Stack, opt BuildOptions) error {
	var (
		totalTime          = (&metrics.Runtime{}).Start()
		svcTime            = &metrics.Runtime{}
		buildTime, pubTime = &metrics.Runtime{}, &metrics.Runtime{}
	)

	if errs := stack.Validate(); len(errs) > 0 {
		return utils.FlattenErrors(errs)
	}

	scopeVars := st.scopeVars(stack)
	printScopeVars(scopeVars)

	var err error
	// Eval variables
	for _, comp := range stack.Components {
		if err = st.evalComponent(comp, scopeVars); err != nil {
			return err
		}
	}

	err = st.crt.CreateNetwork(ctx, stack.ID)
	if err != nil {
		return err
	}

	defer st.Destroy(ctx, stack)

	// Start containers needed for build
	svcTime.Start()
	err = st.startServices(ctx, stack)
	if err != nil {
		return err
	}
	svcTime.End()

	// Start non-head builds
	buildTime.Start()
	nonHeadResults, err := st.startBuilds(ctx, stack, false)
	if err != nil {
		return err
	}

	// Start head builds
	headResults, err := st.startBuilds(ctx, stack, true)
	if err != nil {
		return err
	}
	buildTime.End()

	st.printArtifacts(stack)

	defer func() {
		totalTime.End()
		fmt.Printf("Timing:\n\n")
		fmt.Printf("  Service:\t%v\n", svcTime.Duration(time.Millisecond))
		fmt.Printf("  Build:\t%v\n", buildTime.Duration(time.Millisecond))
		for k, v := range nonHeadResults {
			fmt.Printf("    %s:\t%v\n", k, v.Runtime.Duration(time.Millisecond))
		}
		for k, v := range headResults {
			fmt.Printf("    %s:\t%v\n", k, v.Runtime.Duration(time.Millisecond))
		}
		fmt.Printf("  Publish:\t%v\n", pubTime.Duration(time.Millisecond))
		fmt.Printf("  Total:\t%v\n", totalTime.Duration(time.Millisecond))
		fmt.Println()
	}()

	// Check if we can publish artifacts after the build completes
	status, err := st.vcs.Status(vcs.Option{Path: opt.Workdir})
	if err != nil {
		return err
	}

	if !status.IsClean() {
		fmt.Printf("Uncommitted code:\n\n")
		fmt.Println(status)

		if !opt.Publish {
			fmt.Printf("Artifacts will not be published!\n\n")
			return nil
		}

		fmt.Println("** Explicit artifact publish requested (source code & artifacts may be out of sync) **")
	}

	pubTime.Start()
	st.publish(stack)
	pubTime.End()

	return nil
}

func (st *Stack) printArtifacts(stack *thrapb.Stack) {
	fmt.Printf("\nArtifacts:\n\n")
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		names := st.getBuildImageTags(stack.ID, comp)
		for _, name := range names {
			fmt.Println(" ", name)
		}
		fmt.Println()
	}
}

func (st *Stack) publish(stack *thrapb.Stack) {
	fmt.Printf("\n[TODO] Publishing artifacts:\n\n")
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		names := st.getBuildImageTags(stack.ID, comp)
		for _, name := range names {
			fmt.Println(" ", name)
		}
		fmt.Println()
	}
}

// Deploy deploys all components of the stack.
func (st *Stack) Deploy(stack *thrapb.Stack) error {
	if errs := stack.Validate(); len(errs) > 0 {
		return utils.FlattenErrors(errs)
	}

	// Evaluate variables
	svars := st.scopeVars(stack)
	for _, comp := range stack.Components {
		if err := st.evalComponent(comp, svars); err != nil {
			return err
		}
	}

	ctx := context.Background()

	opts := orchestrator.RequestOptions{}
	_, _, err := st.orch.Deploy(ctx, stack, opts)
	if err != nil {
		st.orch.Destroy(ctx, stack)
	}

	return err
}

// Destroy removes call components of the stack from the container runtime
func (st *Stack) Destroy(ctx context.Context, stack *thrapb.Stack) []*thrapb.ActionReport {
	ar := make([]*thrapb.ActionReport, 0, len(stack.Components))

	for _, c := range stack.Components {
		r := &thrapb.ActionReport{Action: thrapb.NewAction("destroy", "comp", c.ID)}
		r.Error = st.crt.Remove(ctx, c.ID+"."+stack.ID)
		ar = append(ar, r)
	}
	return ar
}

// Stop shutsdown any running containers in the stack.
func (st *Stack) Stop(ctx context.Context, stack *thrapb.Stack) []*thrapb.ActionReport {
	ar := make([]*thrapb.ActionReport, 0, len(stack.Components))

	for _, c := range stack.Components {
		r := &thrapb.ActionReport{Action: thrapb.NewAction("stop", "comp", c.ID)}
		r.Error = st.crt.Stop(ctx, c.ID+"."+stack.ID)
		ar = append(ar, r)
	}
	return ar
}

func (st *Stack) startBuilds(ctx context.Context, stack *thrapb.Stack, head bool) (map[string]*CompBuildResult, error) {
	var (
		results = make(map[string]*CompBuildResult, len(stack.Components))
		err     error
	)

	// Start build containers after
	for id, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		// Build based on whether head was requested
		if comp.Head != head {
			continue
		}

		// err = st.buildStages(ctx, stack.ID, comp)
		results[id], err = st.doBuild(ctx, stack.ID, stack.Version, comp)
		if err != nil {
			break
		}

		// Start container from image that was just built, if this component
		// is not the head
		if !comp.Head {
			err = st.startContainer(ctx, stack.ID, comp)
			if err != nil {
				break
			}
		}
	}

	return results, err
}

// getBuildImageTags returns tags that should be applied to a given image build
func (st *Stack) getBuildImageTags(stackID string, comp *thrapb.Component) []string {
	base := filepath.Join(stackID, comp.ID)
	out := []string{base}
	if len(comp.Version) > 0 {
		out = append(out, base+":"+comp.Version)
	}

	rconf := st.conf.GetDefaultRegistry()
	if rconf != nil && len(rconf.Addr) > 0 {
		rbase := filepath.Join(rconf.Addr, base)
		out = append(out, rbase)
		if len(comp.Version) > 0 {
			out = append(out, rbase+":"+comp.Version)
		}
	}
	return out
}

func (st *Stack) makeBuildRequest(sid, sver string, comp *thrapb.Component) *crt.BuildRequest {
	req := &crt.BuildRequest{
		Output:     crt.NewDockerBuildLog(os.Stdout),
		ContextDir: comp.Build.Context,
		BuildOpts: &types.ImageBuildOptions{
			Tags: st.getBuildImageTags(sid, comp),
			// ID to use in order to cancel the build
			BuildID:     comp.ID,
			Dockerfile:  comp.Build.Dockerfile,
			NetworkMode: sid,
			// Add labels to query later
			Labels: map[string]string{
				"stack":               sid,
				"component":           comp.ID,
				vars.ComponentVersion: comp.Version,
				vars.StackVersion:     sver,
			},
		},
	}

	if comp.HasEnvVars() {
		args := make(map[string]*string, len(comp.Env.Vars))

		fmt.Printf("\nBuild arguments:\n\n")
		for k := range comp.Env.Vars {
			fmt.Println(" ", k)

			v := comp.Env.Vars[k]
			args[k] = &v
		}
		fmt.Println()

		req.BuildOpts.BuildArgs = args
	}

	return req
}

func (st *Stack) doBuild(ctx context.Context, sid, sver string, comp *thrapb.Component) (*CompBuildResult, error) {
	result := &CompBuildResult{
		Runtime: (&metrics.Runtime{}).Start(),
	}

	req := st.makeBuildRequest(sid, sver, comp)

	fmt.Printf("Building %s:\n\n", comp.ID)

	// Blocking
	err := st.crt.Build(ctx, req)
	result.Runtime.End()
	result.Labels = req.BuildOpts.Labels

	return result, err
}

func (st *Stack) evalComponent(comp *thrapb.Component, scopeVars scope.Variables) error {

	var (
		vm  = pseudo.NewVM()
		err error
	)

	if comp.HasEnvVars() {
		err = st.evalCompEnv(comp, vm, scopeVars)
	}

	// TODO: In the future, eval other parts of the component

	return err
}

func (st *Stack) evalCompEnv(comp *thrapb.Component, vm *pseudo.VM, scopeVars scope.Variables) error {
	for k, v := range comp.Env.Vars {
		result, err := vm.ParseEval(v, scopeVars)
		if err != nil {
			return err
		}
		if result.Type != hil.TypeString {
			return fmt.Errorf("env value must be string key=%s value=%s", k, v)
		}
		comp.Env.Vars[k] = result.Value.(string)
	}

	return nil
}

func printScopeVars(scopeVars scope.Variables) {
	fmt.Printf("\nScope:\n\n")
	for _, name := range scopeVars.Names() {
		fmt.Println(" ", name)
	}
	fmt.Println()
}
