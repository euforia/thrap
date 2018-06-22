package config

import (
	"io/ioutil"

	"github.com/euforia/hclencoder"
	"github.com/euforia/pseudo/scope"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hil/ast"
)

// CredsConfig holds creds
type CredsConfig struct {
	Registry     map[string]map[string]string `hcl:"registry"`
	VCS          map[string]map[string]string `hcl:"vcs"`
	Secrets      map[string]map[string]string `hcl:"secrets"`
	Orchestrator map[string]map[string]string `hcl:"id"`
}

func (cc *CredsConfig) GetRegistryCreds(id string) map[string]string {
	return cc.Registry[id]
}
func (cc *CredsConfig) GetVCSCreds(id string) map[string]string {
	return cc.VCS[id]
}
func (cc *CredsConfig) GetSecretsCreds(id string) map[string]string {
	return cc.Secrets[id]
}
func (cc *CredsConfig) GetOrchestratorCreds(id string) map[string]string {
	return cc.Orchestrator[id]
}

func ReadCredsConfig(fpath string) (*CredsConfig, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	var cc CredsConfig
	err = hcl.Unmarshal(b, &cc)

	return &cc, err
}

func WriteCredsConfig(cc *CredsConfig, fpath string) error {
	b, err := hclencoder.Encode(cc)
	if err == nil {
		err = ioutil.WriteFile(fpath, b, 0644)
	}
	return err
}

func DefaultCredsConfig() *CredsConfig {
	return &CredsConfig{
		Registry: make(map[string]map[string]string),
		VCS: map[string]map[string]string{
			"github": map[string]string{"token": ""},
		},
		Secrets:      make(map[string]map[string]string),
		Orchestrator: make(map[string]map[string]string),
	}
}

type VCSRepoConfig struct {
	Name  string `hcl:"name"`
	Owner string `hcl:"owner"`
}

type VCSConfig struct {
	ID       string        `hcl:"id" hcle:"omit"`
	Addr     string        `hcl:"addr" hcle:"omitempty"`
	Username string        `hcl:"username"`
	Repo     VCSRepoConfig `hcl:"repo" hcle:"omitempty"`
}

func (conf *VCSConfig) ScopeVars() scope.Variables {
	return scope.Variables{
		"id": ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		"addr": ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
		"username": ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		},
		"repo.owner": ast.Variable{
			Value: conf.Repo.Owner,
			Type:  ast.TypeString,
		},
	}
}

type OrchestratorConfig struct {
	ID   string `hcl:"id" hcle:"omit"`
	Addr string `hcl:"addr" hcle:"omitempty"`
}

type RegistryRepoConfig struct {
	Name string `hcl:"name"`
}

type RegistryConfig struct {
	ID   string                 `hcl:"id" hcle:"omit"`
	Addr string                 `hcl:"addr" hcle:"omitempty"`
	Repo *RegistryRepoConfig    `hcl:"repo" hcle:"omitempty"`
	Conf map[string]interface{} `hcl:"conf" hcle:"omitempty"`
}

func (conf *RegistryConfig) ScopeVars() scope.Variables {
	return scope.Variables{
		"id": ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		"addr": ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
		"repo.name": ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		},
	}
}

type SecretsConfig struct {
	ID   string `hcl:"id" hcle:"omit"`
	Addr string `hcl:"addr" hcle:"omitempty"`
}

type ThrapConfig struct {
	VCS          map[string]*VCSConfig          `hcl:"vcs"`
	Orchestrator map[string]*OrchestratorConfig `hcl:"orchestrator"`
	Registry     map[string]*RegistryConfig     `hcl:"registry"`
	Secrets      map[string]*SecretsConfig      `hcl:"secrets"`
}

func (conf *ThrapConfig) GetVCS(id string) *VCSConfig {
	return conf.VCS["id"]
}
func (conf *ThrapConfig) GetOrchestrator(id string) *OrchestratorConfig {
	return conf.Orchestrator["id"]
}
func (conf *ThrapConfig) GetRegistry(id string) *RegistryConfig {
	return conf.Registry["id"]
}
func (conf *ThrapConfig) GetSecrets(id string) *SecretsConfig {
	return conf.Secrets["id"]
}

func (conf *ThrapConfig) ScopeVars() scope.Variables {
	svars := make(scope.Variables)
	for k, v := range conf.VCS {
		vars := v.ScopeVars()
		for k1, v1 := range vars {
			svars["vcs."+k+"."+k1] = v1
		}
	}

	for k, v := range conf.Registry {
		vars := v.ScopeVars()
		for k1, v1 := range vars {
			svars["registry."+k+"."+k1] = v1
		}
	}

	// svars := conf.VCS.ScopeVars()

	// rvars := conf.Registry.ScopeVars()
	// for k, v := range rvars {
	// 	svars[k] = v
	// }

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

// ReadThrapConfig reads an hcl thrap config from filename
func ReadThrapConfig(filename string) (*ThrapConfig, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf ThrapConfig
	err = hcl.Unmarshal(b, &conf)
	return &conf, err
}
