package config

import (
	"io/ioutil"

	"github.com/euforia/hclencoder"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/vars"
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
	ID       string        `hcl:"id"`
	Addr     string        `hcl:"addr"`
	Username string        `hcl:"username"`
	Repo     VCSRepoConfig `hcl:"repo"`
}

func (conf *VCSConfig) ScopeVars() scope.Variables {
	return scope.Variables{
		vars.VcsID: ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		vars.VcsAddr: ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
		vars.VcsRepoName: ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		},
		vars.VcsRepoOwner: ast.Variable{
			Value: conf.Repo.Owner,
			Type:  ast.TypeString,
		},
	}
}

type OrchestratorConfig struct {
	ID   string `hcl:"id"`
	Addr string `hcl:"addr"`
}

type RegistryRepoConfig struct {
	Name string `hcl:"name"`
}

type RegistryConfig struct {
	ID   string                 `hcl:"id"`
	Addr string                 `hcl:"addr"`
	Repo *RegistryRepoConfig    `hcl:"repo" hcle:"omitempty"`
	Conf map[string]interface{} `hcl:"conf" hcle:"omitempty"`
}

func (conf *RegistryConfig) ScopeVars() scope.Variables {
	return scope.Variables{
		vars.RegistryID: ast.Variable{
			Value: conf.ID,
			Type:  ast.TypeString,
		},
		vars.RegistryAddr: ast.Variable{
			Value: conf.Addr,
			Type:  ast.TypeString,
		},
		vars.RegistryRepoName: ast.Variable{
			Value: conf.Repo.Name,
			Type:  ast.TypeString,
		},
	}
}

type SecretsConfig struct {
	ID   string `hcl:"id"`
	Addr string `hcl:"addr"`
}

type ThrapConfig struct {
	VCS          *VCSConfig          `hcl:"vcs"`
	Orchestrator *OrchestratorConfig `hcl:"orchestrator"`
	Registry     *RegistryConfig     `hcl:"registry"`
	Secrets      *SecretsConfig      `hcl:"secrets"`
}

func (conf *ThrapConfig) ScopeVars() scope.Variables {
	svars := conf.VCS.ScopeVars()

	rvars := conf.Registry.ScopeVars()
	for k, v := range rvars {
		svars[k] = v
	}

	return svars
}

func DefaultThrapConfig() *ThrapConfig {
	return &ThrapConfig{
		VCS: &VCSConfig{
			ID:   "github",
			Addr: "github.com",
		},
		Orchestrator: &OrchestratorConfig{},
		Registry:     &RegistryConfig{},
		Secrets: &SecretsConfig{
			ID: "file",
			// For file this is the path
			Addr: "secrets",
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
