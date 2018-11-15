package core

import (
	"path/filepath"

	"github.com/euforia/thrap/pkg/provider/vcs"
	"github.com/euforia/thrap/thrapb"
	"github.com/pkg/errors"
)

// EnsureResources ensures that all stack resources exist or creates them as
// necessary
func (st *Stack) EnsureResources(stack *thrapb.Stack) thrapb.ActionsResults {
	results := make(map[string][]*thrapb.ActionResult, 3)
	// Source repo
	results["vcs"] = []*thrapb.ActionResult{st.createVcsRepo(stack)}
	// Registries
	results["registry"] = st.ensureRegistryRepos(stack)
	// Secrets
	results["secrets"] = st.ensureSecrets(stack)

	return results
}

func (st *Stack) ensureRegistryRepos(stack *thrapb.Stack) []*thrapb.ActionResult {
	reports := make([]*thrapb.ActionResult, 0, len(stack.Components))
	for id, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		repoName := stack.ArtifactName(id)

		report := &thrapb.ActionResult{
			Action:   "create",
			Resource: st.reg.ImageName(repoName),
		}

		// Exists
		if _, err := st.reg.GetRepo(repoName); err != nil {
			// Create registry repo
			_, report.Error = st.reg.CreateRepo(repoName)
			report.Data = "created"
		} else {
			report.Data = "exists"
		}

		reports = append(reports, report)
	}

	return reports
}

func (st *Stack) ensureSecrets(stack *thrapb.Stack) []*thrapb.ActionResult {
	reports := make([]*thrapb.ActionResult, 0, len(stack.Components))
	for id, comp := range stack.Components {
		if !comp.HasSecrets() {
			continue
		}

		report := &thrapb.ActionResult{
			Action:   "create",
			Resource: filepath.Join(stack.ID, id),
			Error:    errors.New("to be implemented"),
		}

		reports = append(reports, report)
	}

	return reports
}

func (st *Stack) createVcsRepo(stack *thrapb.Stack) *thrapb.ActionResult {

	er := &thrapb.ActionResult{
		Action:   "create",
		Resource: st.vcs.ID(),
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
	_, created, err := st.vcs.Create(repo, vcsOpt)
	if err == nil {
		if created {
			er.Data = "created"
		} else {
			er.Data = "exists"
		}
	}

	return er
}
