package config

import (
	"github.com/euforia/pseudo/scope"
	"github.com/hashicorp/hil/ast"
)

// VCSConfig holds vcs configurations
type VCSConfig struct {
	ID       string      `hcl:"id" hcle:"omit"`
	Addr     string      `hcl:"addr" hcle:"omitempty"`
	Username string      `hcl:"username"`
	Repo     *RepoConfig `hcl:"repo" hcle:"omitempty"`
}

// Clone returns a copy of the config
func (conf *VCSConfig) Clone() *VCSConfig {
	if conf == nil {
		return nil
	}

	return &VCSConfig{
		ID:       conf.ID,
		Addr:     conf.Addr,
		Username: conf.Username,
		Repo:     conf.Repo.Clone(),
	}
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *VCSConfig) Merge(other *VCSConfig) {
	if other == nil {
		return
	}

	if other.ID != "" {
		conf.ID = other.ID
	}

	if other.Addr != "" {
		conf.Addr = other.Addr
	}
	if other.Username != "" {
		conf.Username = other.Username
	}

	if conf.Repo == nil {
		conf.Repo = &RepoConfig{}
	}
	conf.Repo.Merge(other.Repo)
}

// ScopeVars returns the scoped variables usable for interpolation
func (conf *VCSConfig) ScopeVars(prefix string) scope.Variables {
	svars := scope.Variables{
		prefix + "id": ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		prefix + "addr": ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
		prefix + "username": ast.Variable{
			Value: conf.Username,
			Type:  ast.TypeString,
		},
	}

	if conf.Repo != nil {
		svars[prefix+"repo.owner"] = ast.Variable{
			Value: conf.Repo.Owner,
			Type:  ast.TypeString,
		}
		svars[prefix+"repo.name"] = ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		}
	}

	return svars
}
