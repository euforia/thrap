package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

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

// type DeployOptions struct{
// 	Dryrun bool
// }

var (
	// ErrStackAlreadyRegistered is used when a stack is already registered
	ErrStackAlreadyRegistered = errors.New("stack already registered")
	errComponentNotBuildable  = errors.New("component not buildable")
	errArtifactsMissing       = errors.New("one or more artifacts missing")
)

// StackStore implements a thrap stack store
type StackStore interface {
	Get(id string) (*thrapb.Stack, error)
	Iter(prefix string, f func(*thrapb.Stack) error) error
	Register(stack *thrapb.Stack) (*thrapb.Stack, []*thrapb.ActionReport, error)
}

// Stack provides various stack based operations. It is initialized based
// on the supplied profile
type Stack struct {
	// builder is docker runtime
	crt *crt.Docker

	// config to use for this instance
	conf *config.ThrapConfig

	// code version control provider
	vcs vcs.VCS

	// registry loaded based on profile
	reg registry.Registry

	// orchestrator loaded based on profile
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
func (st *Stack) Status(ctx context.Context, stack *thrapb.Stack) []*thrapb.CompStatus {
	return st.orch.Status(ctx, stack)
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
		pubTime   = &metrics.Runtime{}
		scopeVars = st.scopeVars(stack)
		// rconf     = st.conf.DefaultRegistry()
		err error
	)

	printScopeVars(scopeVars)

	// Eval variables
	for _, comp := range stack.Components {
		if err = st.evalComponent(comp, scopeVars); err != nil {
			return err
		}
	}

	bldr := newStackBuilder(st.crt, st.reg, stack)
	err = bldr.Build(ctx)
	if err != nil {
		return err
	}

	results := bldr.Results()

	// Write timings at the end
	defer func() {
		totalTime.End()
		printBuildStats(bldr, totalTime, pubTime)

		var status string
		if bldr.Succeeded() {
			status = "SUCCEEDED"
		} else {
			status = "FAILED"
		}

		fmt.Printf("\n%s\n", status)
	}()

	defer func() {
		fmt.Printf("SUMMARY\n\n  Build:\n")
		printBuildResults(stack, results, os.Stdout)
	}()

	if !bldr.Succeeded() {
		return err
	}

	//
	// Publish
	//
	canPublish, err := st.checkWorktree(opt)
	if err != nil {
		return err
	}

	fmt.Printf("\nArtifacts:\n\n Generated:\n\n")
	st.printArtifacts(stack, true)

	if canPublish {
		pubTime = st.publishArtifacts(stack)
	}

	return nil
}

// Deploy deploys all components of the stack.
func (st *Stack) Deploy(stack *thrapb.Stack, opts orchestrator.RequestOptions) error {
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

	printScopeVarsWithVals(svars)

	// TODO: check artifact existence
	fmt.Printf("\nArtifacts:\n\n")
	err := st.checkArtifactsExist(stack)
	if err != nil {
		return err
	}

	ctx := context.Background()

	_, _, err = st.orch.Deploy(ctx, stack, opts)
	if err != nil {
		st.orch.Destroy(ctx, stack)
	}

	return err
}

func (st *Stack) checkArtifactsExist(stack *thrapb.Stack) error {

	var (
		reg     = st.reg
		reports = make(map[string]error, len(stack.Components))
		failed  bool
	)

	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		name := stack.ArtifactName(comp.ID)
		_, err := reg.GetManifest(name, comp.Version)
		if err != nil {
			failed = true
		}
		reports[st.reg.ImageName(name+":"+comp.Version)] = err
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, " \tArtifact\tStatus\n")
	fmt.Fprintf(tw, " \t--------\t------\n")
	for k, err := range reports {
		if err != nil {
			failed = true
			fmt.Fprintf(tw, " \t%s\t%v\n", k, err)
		} else {
			fmt.Fprintf(tw, " \t%s\t%v\n", k, "ok")
		}
	}
	tw.Flush()
	fmt.Println()

	if failed {
		return errArtifactsMissing
	}

	return nil
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
	// runtime := &metrics.Runtime{}

	// // Check if we can publish artifacts after the build completes
	// status, err := st.vcs.Status(vcs.Option{Path: opt.Workdir})
	// if err != nil {
	// 	return runtime, err
	// }
	// // We only auto-publish if the working tree is clean
	// if !status.IsClean() {
	// 	fmt.Printf("\nUncommitted code:\n\n")
	// 	fmt.Println(status)

	// 	if !opt.Publish {
	// 		fmt.Println("Artifacts will not be published!")
	// 		return runtime, nil
	// 	}

	// 	fmt.Println("** Explicit artifact publish requested (source code & artifacts may be out of sync) **")
	// }
	canPublish, err := st.checkWorktree(opt)
	if err != nil {
		return &metrics.Runtime{}, err
	}

	if canPublish {
		return st.publishArtifacts(stack), nil
	}

	return &metrics.Runtime{}, nil
}

// returns true if we can publish
func (st *Stack) checkWorktree(opt BuildOptions) (bool, error) {
	status, err := st.vcs.Status(vcs.Option{Path: opt.Workdir})
	if err != nil {
		return false, err
	}

	// We only auto-publish if the working tree is clean
	if !status.IsClean() {
		fmt.Printf("\nUncommitted code:\n\n")
		fmt.Println(status)

		if !opt.Publish {
			fmt.Println("Artifacts will not be published!")
			return false, nil
		}

		fmt.Println("** Explicit artifact publish requested (source code & artifacts may be out of sync) **")
	}

	return true, nil
}

func (st *Stack) publishArtifacts(stack *thrapb.Stack) *metrics.Runtime {
	runtime := &metrics.Runtime{}
	runtime.Start()
	fmt.Printf(" Publishing:\n\n")
	for id, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		fmt.Printf("  %s:\n\n", id)

		name := stack.ArtifactName(id)
		name = st.reg.ImageName(name)
		fmt.Printf("   %s\n", name)
		fmt.Printf("   %s:%s\n\n", name, comp.Version)
	}
	runtime.End()
	return runtime
}

func (st *Stack) printArtifacts(stack *thrapb.Stack, printBase bool) {
	// fmt.Printf("\nArtifacts:\n\n")
	for k, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		fmt.Printf("  %s:\n\n", comp.ID)
		name := stack.ArtifactName(k)
		name = st.reg.ImageName(name)

		if printBase {
			fmt.Printf("    %s\n", name)
		}
		fmt.Printf("    %s:%s\n\n", name, comp.Version)
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
