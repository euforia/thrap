package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/hclencoder"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
	"github.com/hashicorp/hcl"
)

// ThrapConfig holds configs for all providers
type ThrapConfig struct {
	VCS          map[string]*VCSConfig          `hcl:"vcs"`
	Orchestrator map[string]*OrchestratorConfig `hcl:"orchestrator"`
	Registry     map[string]*RegistryConfig     `hcl:"registry"`
	Secrets      map[string]*SecretsConfig      `hcl:"secrets"`
}

// Clone returns a copy of the config
func (conf *ThrapConfig) Clone() *ThrapConfig {
	if conf == nil {
		return nil
	}

	c := &ThrapConfig{
		VCS:          make(map[string]*VCSConfig, len(conf.VCS)),
		Orchestrator: make(map[string]*OrchestratorConfig, len(conf.Orchestrator)),
		Registry:     make(map[string]*RegistryConfig, len(conf.Registry)),
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
func (conf *ThrapConfig) Merge(other *ThrapConfig) {
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
func (conf *ThrapConfig) DefaultVCS() *VCSConfig {
	for _, v := range conf.VCS {
		return v
	}
	return nil
}

// DefaultOrchestrator returns the first available orchestrator
func (conf *ThrapConfig) DefaultOrchestrator() *OrchestratorConfig {
	for _, v := range conf.Orchestrator {
		return v
	}
	return nil
}

// DefaultRegistry returns the first registry from the map
func (conf *ThrapConfig) DefaultRegistry() *RegistryConfig {
	for _, v := range conf.Registry {
		return v
	}
	return nil
}

// DefaultSecrets returns the first secrets provider
func (conf *ThrapConfig) DefaultSecrets() *SecretsConfig {
	for _, v := range conf.Secrets {
		return v
	}
	return nil
}

// ScopeVars returns the scoped variables usable for interpolation
func (conf *ThrapConfig) ScopeVars() scope.Variables {
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

// DefaultThrapConfig returns a minimal config
func DefaultThrapConfig() *ThrapConfig {
	return &ThrapConfig{
		VCS: map[string]*VCSConfig{
			"github": &VCSConfig{
				ID:   "github",
				Addr: "github.com",
			},
		},
		Orchestrator: map[string]*OrchestratorConfig{
			"docker": &OrchestratorConfig{},
		},
		Registry: make(map[string]*RegistryConfig),
		Secrets: map[string]*SecretsConfig{
			"file": &SecretsConfig{},
		},
	}
}

// WriteThrapConfig write the given config to the specified file path
func WriteThrapConfig(conf *ThrapConfig, filename string) error {
	b, err := hclencoder.Encode(conf)
	if err == nil {
		err = ioutil.WriteFile(filename, b, 0644)
	}
	return err
}

// ReadThrapConfig reads an hcl thrap config from filename
func ReadThrapConfig(filename string) (*ThrapConfig, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf ThrapConfig
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
func ReadProjectConfig(projPath string) (*ThrapConfig, error) {
	filename := filepath.Join(projPath, consts.WorkDir, consts.ConfigFile)
	return ReadThrapConfig(filename)
}

// func ReadGlobalConfig() (*ThrapConfig, error) {
// 	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.ConfigFile)
// 	if err == nil {
// 		return ReadThrapConfig(filename)
// 	}
// 	return nil, err
// }

// func ReadGlobalCreds() (*CredsConfig, error) {
// 	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.CredsFile)
// 	if err == nil {
// 		return ReadCredsConfig(filename)
// 	}
// 	return nil, err
// }
