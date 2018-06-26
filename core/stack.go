package core

import (
	"fmt"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/registry"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
	"github.com/pkg/errors"
)

// Stack provides various stack based operations
type Stack struct {
	st *thrapb.Stack

	dir string // project directory

	// config to use for this instance
	conf *config.ThrapConfig
	// there can only be one vcs provider
	vcs vcs.VCS
	// Available registries
	regs map[string]registry.Registry
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
func (st *Stack) Validate() error {
	st.st.Version = vcs.GetRepoVersion(st.dir).String()

	errs := st.st.Validate()
	if len(errs) > 0 {
		return utils.FlattenErrors(errs)
	}
	return nil
}

func (st *Stack) Init(opt ConfigureOptions) error {
	_, err := ConfigureLocal(opt)
	if err != nil {
		return err
	}

	repo := opt.VCS.Repo
	_, _, err = vcs.SetupLocalGitRepo(repo.Name, repo.Owner, opt.DataDir, opt.VCS.Addr)

	return err
}

func (st *Stack) Build() {}

func (st *Stack) Deploy() {}

func (st *Stack) ensureStackResources(stack *thrapb.Stack) []*ActionReport {
	report := st.createVcsRepo(stack)
	reports := st.ensureComponentResources(stack)

	// Concat results from both of the above
	return append([]*ActionReport{report}, reports...)
}

// returns the report and an error if any component failed
func (st *Stack) ensureComponentResources(stack *thrapb.Stack) []*ActionReport {

	r := make([]*ActionReport, 0, len(stack.Components))
	for i, comp := range stack.Components {

		if !comp.IsBuildable() {
			continue
		}
		// TEMP use ecr
		r = append(r, st.createRegistryRepo("ecr", stack.ID, i))

		if comp.HasSecrets() {
			r = append(r, st.createSecrets(i))
		}

	}

	return r
}

func (st *Stack) createVcsRepo(stack *thrapb.Stack) *ActionReport {

	er := &ActionReport{Action: NewAction("create", "vcs.repo."+st.vcs.ID(), stack.ID)}

	if st.vcs == nil {
		er.Error = errProviderNotConfigured
		return er
	}

	var vcsOpt vcs.Option
	repo := &vcs.Repository{
		Name:        stack.ID,
		Description: stack.Description,
	}
	vc := st.conf.VCS[st.vcs.ID()]
	if vc.Repo.Owner != "" {
		repo.Owner = vc.Repo.Owner
	}

	er.Data, er.Error = st.vcs.Create(repo, vcsOpt)

	return er
}

// createSecrets performs actions for secret creation for the given component
// id
func (st *Stack) createSecrets(compID string) *ActionReport {
	return &ActionReport{
		Action: NewAction("create", "secrets", compID),
		Error:  errors.New("to be implemented"),
	}
}

// createRegistryRepo performs actions for registry creation
func (st *Stack) createRegistryRepo(regID, stackID, compID string) *ActionReport {
	repoName := stackID + "/" + compID
	er := &ActionReport{
		Action: NewAction("create", "registry.repo", compID),
	}

	reg, ok := st.regs[regID]
	if !ok {
		er.Error = errors.Wrap(errProviderNotConfigured, regID)
		return er
	}

	// Create registry repo
	er.Data, er.Error = reg.Create(repoName)

	return er
}
