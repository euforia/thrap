package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

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
	// profile string

	// builder is docker runtime
	crt *crt.Docker
	// common functions
	// run *bdCommon

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

	if errs := stack.Validate(); len(errs) > 0 {
		return utils.FlattenErrors(errs)
	}

	var (
		totalTime = (&metrics.Runtime{}).Start()
		scopeVars = st.scopeVars(stack)
		rconf     = st.conf.DefaultRegistry()
		pubTime   *metrics.Runtime
		err       error
	)

	printScopeVars(scopeVars)

	// Eval variables
	for _, comp := range stack.Components {
		if err = st.evalComponent(comp, scopeVars); err != nil {
			return err
		}
	}

	bldr := newStackBuilder(st.crt, rconf, stack)
	err = bldr.Build(ctx)
	if err != nil {
		return err
	}

	results := bldr.Results()

	// Write timings at the end
	defer func() {
		totalTime.End()
		printBuildStats(bldr, totalTime, pubTime)
		fmt.Println()
	}()

	defer func() {
		fmt.Printf("\nBuild summary:\n")
		printBuildResults(stack, results, os.Stdout)
	}()

	pubTime, err = st.publish(rconf, stack, opt)

	return err
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
	return st.orch.Destroy(ctx, stack)
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

// publish checks if artifacts are publishable and only then will publish them
func (st *Stack) publish(rconf *config.RegistryConfig, stack *thrapb.Stack, opt BuildOptions) (*metrics.Runtime, error) {
	runtime := &metrics.Runtime{}

	// Check if we can publish artifacts after the build completes
	status, err := st.vcs.Status(vcs.Option{Path: opt.Workdir})
	if err != nil {
		return runtime, err
	}

	if !status.IsClean() {
		fmt.Printf("\nUncommitted code:\n\n")
		fmt.Println(status)

		if !opt.Publish {
			fmt.Println("Artifacts will not be published!")
			return runtime, nil
		}

		fmt.Println("** Explicit artifact publish requested (source code & artifacts may be out of sync) **")
	}

	printArtifacts(stack, rconf)

	runtime.Start()
	st.publishArtifacts(stack, rconf)
	runtime.End()

	return runtime, nil
}

func (st *Stack) publishArtifacts(stack *thrapb.Stack, rconf *config.RegistryConfig) {
	fmt.Printf("[TODO] Publishing artifacts:\n\n")
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		names := getBuildImageTags(stack.ID, comp, rconf)
		for _, name := range names {
			fmt.Println(" ", name)
		}
		fmt.Println()
	}
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
