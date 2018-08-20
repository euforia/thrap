package provider

import (
	"github.com/euforia/pseudo/scope"
	"github.com/hashicorp/hil/ast"
)

// Config holds configurations for a registry
type Config struct {
	ID       string                 `hcl:"id"     hcle:"omit"`
	Provider string                 `hcl:"provider"`
	Addr     string                 `hcl:"addr"   hcle:"omitempty"`
	Config   map[string]interface{} `hcl:"config" hcle:"omitempty"`
}

// Clone returns a copy of the config
func (conf *Config) Clone() *Config {
	if conf == nil {
		return nil
	}

	rc := &Config{
		ID:       conf.ID,
		Provider: conf.Provider,
		Addr:     conf.Addr,
	}

	if conf.Config != nil {
		conf.Config = make(map[string]interface{})
		for k, v := range conf.Config {
			rc.Config[k] = v
		}
	}

	return rc
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *Config) Merge(other *Config) {
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
func (conf *Config) ScopeVars(prefix string) scope.Variables {
	return scope.Variables{
		prefix + "id": ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		prefix + "addr": ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
	}
}
