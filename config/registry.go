package config

import (
	"github.com/euforia/pseudo/scope"
	"github.com/hashicorp/hil/ast"
)

type RegistryRepoConfig struct {
	Name string `hcl:"name"`
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *RegistryRepoConfig) Merge(other *RegistryRepoConfig) {
	if other != nil {
		if other.Name != "" {
			conf.Name = other.Name
		}
	}
}

type RegistryConfig struct {
	ID     string                 `hcl:"id"     hcle:"omit"`
	Addr   string                 `hcl:"addr"   hcle:"omitempty"`
	Repo   *RegistryRepoConfig    `hcl:"repo"   hcle:"omitempty"`
	Config map[string]interface{} `hcl:"config" hcle:"omitempty"`
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *RegistryConfig) Merge(other *RegistryConfig) {
	if other == nil {
		return
	}

	if other.ID != "" {
		conf.ID = other.ID
	}

	if other.Addr != "" {
		conf.Addr = other.Addr
	}

	conf.Repo.Merge(other.Repo)

	if other.Config != nil {
		if conf.Config == nil {
			conf.Config = make(map[string]interface{})
		}
		for k, v := range other.Config {
			conf.Config[k] = v
		}
	}
}

// ScopeVars returns the scoped variables usable for interpolation
func (conf *RegistryConfig) ScopeVars(prefix string) scope.Variables {
	return scope.Variables{
		prefix + "id": ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		prefix + "addr": ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
		prefix + "repo.name": ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		},
	}
}
