package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"

	"github.com/euforia/pseudo"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/crt"
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

// Stack provides various stack based operations
type Stack struct {
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
	sst *store.StackStore // local store

	log *log.Logger
}

// Assembler returns a new assembler for the stack
func (st *Stack) Assembler(stack *thrapb.Stack) (*asm.StackAsm, error) {
	// fmt.Println(st.conf.VCS, st.conf.VCS[st.vcs.ID()])
	scopeVars := st.conf.VCS[st.vcs.ID()].ScopeVars("vcs.")

	//vconf := st.conf.VCS[st.vcs.ID()]
	// r, err := st.vcs.Open()
	// if err != nil {
	// 	return nil, err
	// }
	// gitRepo := r.(*git.Repository)
	return asm.NewStackAsm(stack, st.vcs, nil, scopeVars, st.packs)
}

// Register registers a new stack. It returns an error if the stack is
// already registered or fails to register
func (st *Stack) Register(stack *thrapb.Stack) (*thrapb.Stack, []*ActionReport, error) {
	errs := stack.Validate()
	if len(errs) > 0 {
		return nil, nil, utils.FlattenErrors(errs)
	}

	stack, _, err := st.sst.Create(stack)
	if err != nil {
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

// populatePorts populates ports into the stack from the container images if
// no ports have been defined
func (st *Stack) populatePorts(stack *thrapb.Stack) error {
	var err error
	for _, comp := range stack.Components {
		if len(comp.Ports) > 0 {
			continue
		}

		er := st.crt.ImagePull(context.Background(), comp.Name+":"+comp.Version)
		if er != nil {
			continue
		}

		ic, er := st.crt.ImageConfig(comp.Name, comp.Version)
		if er != nil {
			err = er
			continue
		}

		comp.Ports = make(map[string]int32, len(ic.ExposedPorts))
		if len(ic.ExposedPorts) == 1 {
			for k := range ic.ExposedPorts {
				comp.Ports["default"] = int32(k.Int())
				break
			}
		} else {
			for k := range ic.ExposedPorts {
				// HCL does not allow numbers as keys
				comp.Ports["port"+k.Port()] = int32(k.Int())
			}
		}
	}

	return err
}

// Init initializes a basic stack with the configuration and options provided. This should only be
// used in the local cli case as the config is merged with the global.
func (st *Stack) Init(stconf *asm.BasicStackConfig, opt ConfigureOptions) (*thrapb.Stack, error) {

	lconf, err := ConfigureLocal(opt)
	if err != nil {
		return nil, err
	}
	st.conf.Merge(lconf)

	repo := opt.VCS.Repo
	vcsp, gitRepo, err := vcs.SetupLocalGitRepo(repo.Name, repo.Owner, opt.DataDir, opt.VCS.Addr)
	if err != nil {
		return nil, err
	}

	stack, err := asm.NewBasicStack(stconf, st.packs)
	if err != nil {
		return nil, err
	}

	st.populatePorts(stack)

	if errs := stack.Validate(); len(errs) > 0 {
		return stack, utils.FlattenErrors(errs)
	}

	scopeVars := st.conf.VCS[st.vcs.ID()].ScopeVars("vcs.")
	stasm, err := asm.NewStackAsm(stack, vcsp, gitRepo, scopeVars, st.packs)
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
		svars[v.ScopeVarName("comps.", "container.ip")] = ipvar
		// Set container.addr var per port label
		for pl, p := range v.Ports {
			svars[v.ScopeVarName("comps.", "container.addr."+pl)] = ast.Variable{
				Type:  ast.TypeString,
				Value: fmt.Sprintf("%s:%d", ipvar.Value, p),
			}
		}

	}

	return svars
}

// startServices starts services needed to perform the build that themselves do not need
// to be built
func (st *Stack) startServices(ctx context.Context, stack *thrapb.Stack, scopeVars scope.Variables) error {
	var (
		err error
		//opt = crt.RequestOptions{Output: os.Stdout}
	)

	fmt.Printf("\nServices:\n\n")

	for _, comp := range stack.Components {
		if comp.IsBuildable() {
			continue
		}

		// eval hcl/hil
		if err = st.evalComponent(comp, scopeVars); err != nil {
			break
		}

		if err = st.startContainer(ctx, stack.ID, comp); err != nil {
			break
		}

		fmt.Println(comp.ID)

	}

	return err
}

func (st *Stack) startContainer(ctx context.Context, sid string, comp *thrapb.Component) error {
	cfg := thrapb.NewContainer(sid, comp.ID)
	cfg.Container.Image = comp.Name + ":" + comp.Version

	if comp.HasEnvVars() {
		cfg.Container.Env = make([]string, 0, len(comp.Env.Vars))
		for k, v := range comp.Env.Vars {
			cfg.Container.Env = append(cfg.Container.Env, k+"="+v)
		}
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
	return nil
}

// Build starts all require services, then starts all the builds
func (st *Stack) Build(ctx context.Context, stack *thrapb.Stack) error {
	if errs := stack.Validate(); len(errs) > 0 {
		return utils.FlattenErrors(errs)
	}

	err := st.crt.CreateNetwork(ctx, stack.ID)
	if err != nil {
		return err
	}

	scopeVars := st.scopeVars(stack)
	fmt.Printf("\nScope:\n\n")
	fmt.Println(strings.Join(scopeVars.Names(), "\n"))
	fmt.Println()
	// st.log.Printf("DEBUG Scope: %v", scopeVars.Names())

	defer st.Destroy(ctx, stack)

	// Start containers needed for build
	err = st.startServices(ctx, stack, scopeVars)
	if err != nil {
		return err
	}

	// Start non-head builds
	err = st.startBuilds(ctx, stack, scopeVars, false)
	if err != nil {
		return err
	}

	// Start head builds
	err = st.startBuilds(ctx, stack, scopeVars, true)

	return err
}

func (st *Stack) startBuilds(ctx context.Context, stack *thrapb.Stack, scopeVars scope.Variables, head bool) error {
	var err error

	// Start build containers after
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		// Build based on whether head was requested
		if comp.Head != head {
			continue
		}

		// eval hcl/hil
		if err = st.evalComponent(comp, scopeVars); err != nil {
			break
		}

		err = st.doBuild(ctx, stack.ID, comp)
		if err != nil {
			break
		}

		// Start container from image that was just built, if this component
		// is not the head
		if !comp.Head {
			fmt.Printf("%s\n", comp.ID)
			err = st.startContainer(ctx, stack.ID, comp)
			if err != nil {
				break
			}
		}
	}

	return err
}

func (st *Stack) doBuild(ctx context.Context, sid string, comp *thrapb.Component) error {
	tbase := filepath.Join(sid, comp.ID)
	req := &crt.BuildRequest{
		Output:     os.Stdout,
		ContextDir: comp.Build.Context,
		BuildOpts: &types.ImageBuildOptions{
			Tags:        []string{tbase},
			BuildID:     comp.ID, // todo add vcs repo version
			Dockerfile:  comp.Build.Dockerfile,
			NetworkMode: sid,
			// Remove:      true,
		},
	}
	if len(comp.Version) > 0 {
		req.BuildOpts.Tags = append(req.BuildOpts.Tags, tbase+":"+comp.Version)
	}

	if comp.HasEnvVars() {
		req.BuildOpts.BuildArgs = make(map[string]*string, len(comp.Env.Vars))
		fmt.Printf("\nBuild arguments:\n\n")
		for k := range comp.Env.Vars {
			v := comp.Env.Vars[k]
			fmt.Println(k)
			req.BuildOpts.BuildArgs[k] = &v
		}
		fmt.Println()
	}

	// Blocking
	fmt.Printf("Building %s:\n\n", comp.ID)
	return st.crt.Build(ctx, req)
}

func (st *Stack) evalComponent(comp *thrapb.Component, scopeVars scope.Variables) error {
	if !comp.HasEnvVars() {
		return nil
	}

	vm := pseudo.NewVM()

	for k, v := range comp.Env.Vars {
		result, err := vm.ParseEval(v, scopeVars)
		if err != nil {
			return err
		}
		if result.Type != hil.TypeString {
			return fmt.Errorf("env value must be string key=%s value=%s", k, v)
		}
		comp.Env.Vars[k] = result.Value.(string)
		// st.log.Printf("DEBUG env.%s='%s'", k, comp.Env.Vars[k])
	}

	return nil
}

// func (st *Stack) checkDockerfileArgs(comp *thrapb.Component) error {
// 	raw, err := dockerfile.Parse(comp.Build.Dockerfile)
// 	if err != nil {
// 		return err
// 	}
// 	df := dockerfile.ParseRaw(raw)

// 	args := make(map[string]string)
// 	for _, ins := range df.Stages[0] {
// 		if ins.Key() == dockerfile.KeyArg {
// 			args[ins.String()] = ""
// 		}
// 	}

// 	for k, _ := range comp.Env.Vars {
// 		if _, ok := args[k]; !ok {
// 			return fmt.Errorf("missing build arg: %s", k)
// 		}
// 	}

// 	return nil
// }

// Destroy removes call components of the stack from the container runtime
func (st *Stack) Destroy(ctx context.Context, stack *thrapb.Stack) []*ActionReport {
	ar := make([]*ActionReport, 0, len(stack.Components))

	for _, c := range stack.Components {
		r := &ActionReport{Action: NewAction("component", c.ID, "remove")}
		r.Error = st.crt.Remove(ctx, c.ID+"."+stack.ID)
		ar = append(ar, r)
	}
	return ar
}

// Deploy deploys all components of the stack.
func (st *Stack) Deploy(stack *thrapb.Stack) error {

	var (
		ctx = context.Background()
		err = st.crt.CreateNetwork(ctx, stack.ID)
	)

	if err != nil {
		return err
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

		_, err = st.crt.Run(ctx, cfg)
		if err != nil {
			break
		}
	}

	if err != nil {
		defer st.Destroy(ctx, stack)
	}

	// 	return nil, nil, err
	return err
}
