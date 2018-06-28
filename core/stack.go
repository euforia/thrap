package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	// builder
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
func (st *Stack) Validate(stack *thrapb.Stack, ctxDir string) error {
	stack.Version = vcs.GetRepoVersion(ctxDir).String()
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
	stasm, err := asm.NewStackAsm(stack, vcsp, gitRepo, scopeVars, st.packs.Dev())
	if err != nil {
		return stack, err
	}

	err = stasm.Assemble()
	if err == nil {
		err = stasm.WriteManifest()
	}

	return stack, err
}

func (st *Stack) Build(ctx context.Context, stack *thrapb.Stack) error {
	if errs := stack.Validate(); len(errs) > 0 {
		return utils.FlattenErrors(errs)
	}

	opts := crt.RequestOptions{
		Output: os.Stdout,
	}

	err := st.crt.CreateNetwork(ctx, stack.ID)
	if err != nil {
		return err
	}

	defer st.Destroy(ctx, stack)

	for _, comp := range stack.Components {
		if comp.IsBuildable() {
			continue
		}

		cfg := thrapb.NewContainer(stack.ID, comp.ID)
		cfg.Container.Image = comp.Name + ":" + comp.Version

		// Non-blocking
		err := st.crt.Run(ctx, cfg, opts)
		if err != nil {
			return err
		}
	}

	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		// Blocking
		err := st.crt.Build(ctx, stack.ID, comp, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *Stack) Destroy(ctx context.Context, stack *thrapb.Stack) []*ActionReport {
	ar := make([]*ActionReport, 0, len(stack.Components))

	for _, c := range stack.Components {
		r := &ActionReport{Action: NewAction("component", c.ID, "remove")}
		r.Error = st.crt.Remove(ctx, c.ID+"."+stack.ID)
		ar = append(ar, r)
	}
	return ar
}

func (st *Stack) Deploy(stack *thrapb.Stack) error {

	var (
		ctx = context.Background()
		err = st.crt.CreateNetwork(ctx, stack.ID)
	)

	if err != nil {
		return err
	}

	opts := crt.RequestOptions{Output: os.Stdout}

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

		err = st.crt.Run(ctx, cfg, opts)
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
