package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/hclencoder"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/pkg/provider"
	"github.com/hashicorp/hcl"
)

// Config holds configs for all providers
type Config struct {
	VCS          map[string]*VCSConfig       `hcl:"vcs"`
	Orchestrator map[string]*provider.Config `hcl:"orchestrator"`
	Registry     map[string]*provider.Config `hcl:"registry"`
	Secrets      map[string]*SecretsConfig   `hcl:"secrets"`
}

// Clone returns a copy of the config
func (conf *Config) Clone() *Config {
	if conf == nil {
		return nil
	}

	c := &Config{
		VCS:          make(map[string]*VCSConfig, len(conf.VCS)),
		Orchestrator: make(map[string]*provider.Config, len(conf.Orchestrator)),
		Registry:     make(map[string]*provider.Config, len(conf.Registry)),
		Secrets:      make(map[string]*SecretsConfig, len(conf.Secrets)),
	}

	for k, v := range conf.VCS {
		c.VCS[k] = v.Clone()
	}
	for k, v := range conf.Orchestrator {
		c.Orchestrator[k] = v.Clone()
	}
	for k, v := range conf.Registry {
		c.Registry[k] = v.Clone()
	}
	for k, v := range conf.Secrets {
		c.Secrets[k] = v.Clone()
	}

	return conf
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	if other.VCS != nil {
		for k, v := range other.VCS {
			if cv, ok := conf.VCS[k]; ok {
				cv.Merge(v)
			} else {
				conf.VCS[k] = other.VCS[k]
			}
		}
	}

	if other.Orchestrator != nil {
		for k, v := range other.Orchestrator {
			if cv, ok := conf.Orchestrator[k]; ok {
				cv.Merge(v)
			} else {
				conf.Orchestrator[k] = other.Orchestrator[k]
			}
		}
	}

	if other.Registry != nil {
		for k, v := range other.Registry {
			if cv, ok := conf.Registry[k]; ok {
				cv.Merge(v)
			} else {
				conf.Registry[k] = other.Registry[k]
			}
		}
	}

	if other.Secrets != nil {
		for k, v := range other.Secrets {
			if cv, ok := conf.Secrets[k]; ok {
				cv.Merge(v)
			} else {
				conf.Secrets[k] = other.Secrets[k]
			}
		}
	}

}

// DefaultVCS returns the first available vcs
func (conf *Config) DefaultVCS() *VCSConfig {
	for _, v := range conf.VCS {
		return v
	}
	return nil
}

// DefaultOrchestrator returns the first available orchestrator
func (conf *Config) DefaultOrchestrator() *provider.Config {
	for _, v := range conf.Orchestrator {
		return v
	}
	return nil
}

// DefaultRegistry returns the first registry from the map
func (conf *Config) DefaultRegistry() *provider.Config {
	for _, v := range conf.Registry {
		return v
	}
	return nil
}

// DefaultSecrets returns the first secrets provider
func (conf *Config) DefaultSecrets() *SecretsConfig {
	for _, v := range conf.Secrets {
		return v
	}
	return nil
}

// ScopeVars returns the scoped variables usable for interpolation
func (conf *Config) ScopeVars() scope.Variables {
	svars := make(scope.Variables)
	for k, v := range conf.VCS {
		vars := v.ScopeVars("vcs." + k + ".")
		for k1, v1 := range vars {
			svars[k1] = v1
		}
	}

	for k, v := range conf.Registry {
		vars := v.ScopeVars("registry." + k + ".")
		for k1, v1 := range vars {
			svars[k1] = v1
		}
	}

	return svars
}

// DefaultConfig returns a minimal config
func DefaultConfig() *Config {
	return &Config{
		VCS: map[string]*VCSConfig{
			"github": &VCSConfig{
				ID:   "github",
				Addr: "github.com",
			},
		},
		Orchestrator: map[string]*provider.Config{
			"docker": &provider.Config{
				Provider: "docker",
			},
		},
		Registry: map[string]*provider.Config{
			"docker": &provider.Config{
				Provider: "docker",
			},
		},
		Secrets: map[string]*SecretsConfig{
			"file": &SecretsConfig{},
		},
	}
}

// WriteConfig write the given config to the specified file path
func WriteConfig(conf *Config, filename string) error {
	b, err := hclencoder.Encode(conf)
	if err == nil {
		err = ioutil.WriteFile(filename, b, 0644)
	}
	return err
}

// ReadConfig reads an hcl thrap config from filename
func ReadConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = hcl.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	for k, v := range conf.VCS {
		v.ID = k
	}
	for k, v := range conf.Orchestrator {
		v.ID = k
	}
	for k, v := range conf.Registry {
		v.ID = k
	}
	for k, v := range conf.Secrets {
		v.ID = k
	}

	return &conf, nil
}

// ReadProjectConfig reads the configs from the project config directory under a given
// path
func ReadProjectConfig(projPath string) (*Config, error) {
	filename := filepath.Join(projPath, consts.WorkDir, consts.ConfigFile)
	return ReadConfig(filename)
}
