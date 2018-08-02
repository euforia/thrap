package core

import (
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vcs"
	"github.com/pkg/errors"
)

func (st *Stack) ensureStackResources(stack *thrapb.Stack) []*thrapb.ActionReport {
	report := st.createVcsRepo(stack)

	// Registries
	reports := st.ensureRegistryRepos(stack)
	// Secrets
	reports = append(reports, st.ensureSecrets(stack)...)

	// Concat results from both of the above
	return append([]*thrapb.ActionReport{report}, reports...)
}

func (st *Stack) ensureRegistryRepos(stack *thrapb.Stack) []*thrapb.ActionReport {
	reports := make([]*thrapb.ActionReport, 0, len(stack.Components))
	for id, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		report := &thrapb.ActionReport{
			Action: thrapb.NewAction("create", "registry.repo", id),
		}

		// Create registry repo
		repoName := stack.ArtifactName(id)
		report.Data, report.Error = st.reg.Create(repoName)
		reports = append(reports, report)
	}

	return reports
}

func (st *Stack) ensureSecrets(stack *thrapb.Stack) []*thrapb.ActionReport {
	reports := make([]*thrapb.ActionReport, 0, len(stack.Components))
	for id, comp := range stack.Components {
		if !comp.HasSecrets() {
			continue
		}

		report := &thrapb.ActionReport{
			Action: thrapb.NewAction("create", "secrets", id),
			Error:  errors.New("to be implemented"),
		}

		reports = append(reports, report)
	}

	return reports
}

func (st *Stack) createVcsRepo(stack *thrapb.Stack) *thrapb.ActionReport {

	er := &thrapb.ActionReport{
		Action: thrapb.NewAction("create", "vcs.repo."+st.vcs.ID(), stack.ID),
	}

	if st.vcs == nil {
		er.Error = errProviderNotConfigured
		return er
	}

	repo := &vcs.Repository{
		Name:        stack.ID,
		Description: stack.Description,
	}

	vc := st.conf.VCS[st.vcs.ID()]
	if vc.Repo.Owner != "" {
		repo.Owner = vc.Repo.Owner
	}

	var vcsOpt vcs.Option
	er.Data, er.Error = st.vcs.Create(repo, vcsOpt)

	return er
}
