package config

import (
	"github.com/euforia/pseudo/scope"
	"github.com/hashicorp/hil/ast"
)

// RegistryRepoConfig is the configuration for a container registry repo
type RegistryRepoConfig struct {
	Name string `hcl:"name"`
}

// Clone returns a copy of the config
func (conf *RegistryRepoConfig) Clone() *RegistryRepoConfig {
	if conf == nil {
		return nil
	}
	return &RegistryRepoConfig{
		Name: conf.Name,
	}
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

// RegistryConfig holds configurations for a registry
type RegistryConfig struct {
	ID       string                 `hcl:"id"     hcle:"omit"`
	Provider string                 `hcl:"provider"`
	Addr     string                 `hcl:"addr"   hcle:"omitempty"`
	Repo     *RegistryRepoConfig    `hcl:"repo"   hcle:"omitempty"`
	Config   map[string]interface{} `hcl:"config" hcle:"omitempty"`
}

// Clone returns a copy of the config
func (conf *RegistryConfig) Clone() *RegistryConfig {
	if conf == nil {
		return nil
	}
	rc := &RegistryConfig{
		ID:       conf.ID,
		Provider: conf.Provider,
		Addr:     conf.Addr,
		Repo:     conf.Repo.Clone(),
		Config:   make(map[string]interface{}, len(conf.Config)),
	}
	for k, v := range conf.Config {
		rc.Config[k] = v
	}

	return rc
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

	if other.Provider != "" {
		conf.Provider = other.Provider
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
