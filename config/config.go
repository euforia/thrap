package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/hclencoder"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
	"github.com/hashicorp/hcl"
	homedir "github.com/mitchellh/go-homedir"
)

type OrchestratorConfig struct {
	ID   string `hcl:"id" hcle:"omit"`
	Addr string `hcl:"addr" hcle:"omitempty"`
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *OrchestratorConfig) Merge(other *OrchestratorConfig) {
	if other == nil {
		return
	}

	if other.ID != "" {
		conf.ID = other.ID
	}

	if other.Addr != "" {
		conf.Addr = other.Addr
	}
}

type SecretsConfig struct {
	ID   string `hcl:"id"   hcle:"omit"`
	Addr string `hcl:"addr" hcle:"omitempty"`
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *SecretsConfig) Merge(other *SecretsConfig) {
	if other == nil {
		return
	}

	if other.ID != "" {
		conf.ID = other.ID
	}

	if other.Addr != "" {
		conf.Addr = other.Addr
	}
}

type ThrapConfig struct {
	VCS          map[string]*VCSConfig          `hcl:"vcs"`
	Orchestrator map[string]*OrchestratorConfig `hcl:"orchestrator"`
	Registry     map[string]*RegistryConfig     `hcl:"registry"`
	Secrets      map[string]*SecretsConfig      `hcl:"secrets"`
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

func (conf *ThrapConfig) GetDefaultVCS() *VCSConfig {
	for _, v := range conf.VCS {
		return v
	}
	return nil
}
func (conf *ThrapConfig) GetDefaultOrchestrator() *OrchestratorConfig {
	for _, v := range conf.Orchestrator {
		return v
	}
	return nil
}
func (conf *ThrapConfig) GetDefaultRegistry() *RegistryConfig {
	for _, v := range conf.Registry {
		return v
	}
	return nil
}
func (conf *ThrapConfig) GetDefaultSecrets() *SecretsConfig {
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
			// svars["vcs."+k+"."+k1] = v1
			svars[k1] = v1
		}
	}

	for k, v := range conf.Registry {
		vars := v.ScopeVars("registry." + k + ".")
		for k1, v1 := range vars {
			// svars["registry."+k+"."+k1] = v1
			svars[k1] = v1
		}
	}

	return svars
}

func DefaultThrapConfig() *ThrapConfig {
	return &ThrapConfig{
		VCS: map[string]*VCSConfig{
			"github": &VCSConfig{
				Addr: "github.com",
			},
		},
		Orchestrator: map[string]*OrchestratorConfig{
			"docker": &OrchestratorConfig{},
		},
		Registry: make(map[string]*RegistryConfig),
		Secrets: map[string]*SecretsConfig{
			"file": &SecretsConfig{
				// For file this is the path
				//Addr: "secrets",
			},
		},
	}
}

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

func ReadGlobalConfig() (*ThrapConfig, error) {
	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.ConfigFile)
	if err == nil {
		return ReadThrapConfig(filename)
	}
	return nil, err
}

func ReadProjectConfig(projPath string) (*ThrapConfig, error) {
	filename := filepath.Join(projPath, consts.WorkDir, consts.ConfigFile)
	return ReadThrapConfig(filename)
}

func ReadGlobalCreds() (*CredsConfig, error) {
	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.CredsFile)
	if err == nil {
		return ReadCredsConfig(filename)
	}
	return nil, err
}
