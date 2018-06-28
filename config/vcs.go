package config

import (
	"github.com/euforia/pseudo/scope"
	"github.com/hashicorp/hil/ast"
)

type VCSRepoConfig struct {
	Name  string `hcl:"name"`
	Owner string `hcl:"owner"`
}

func (conf *VCSRepoConfig) Clone() *VCSRepoConfig {
	if conf == nil {
		return nil
	}

	return &VCSRepoConfig{
		Name:  conf.Name,
		Owner: conf.Owner,
	}
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *VCSRepoConfig) Merge(other *VCSRepoConfig) {
	if other == nil {
		return
	}

	if other.Name != "" {
		conf.Name = other.Name
	}
	if other.Owner != "" {
		conf.Owner = other.Owner
	}

}

type VCSConfig struct {
	ID       string         `hcl:"id" hcle:"omit"`
	Addr     string         `hcl:"addr" hcle:"omitempty"`
	Username string         `hcl:"username"`
	Repo     *VCSRepoConfig `hcl:"repo" hcle:"omitempty"`
}

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
		conf.Repo = &VCSRepoConfig{}
	}
	conf.Repo.Merge(other.Repo)
}

// ScopeVars returns the scoped variables usable for interpolation
func (conf *VCSConfig) ScopeVars(prefix string) scope.Variables {
	return scope.Variables{
		prefix + "id": ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		prefix + "addr": ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
		prefix + "username": ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		},
		prefix + "repo.owner": ast.Variable{
			Value: conf.Repo.Owner,
			Type:  ast.TypeString,
		},
		prefix + "repo.name": ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		},
	}
}
