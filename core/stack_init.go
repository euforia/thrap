package core

import (
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vcs"
	"github.com/pkg/errors"
)

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
