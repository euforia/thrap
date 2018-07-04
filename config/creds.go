package config

import (
	"io/ioutil"

	"github.com/euforia/hclencoder"
	"github.com/hashicorp/hcl"
)

// CredsConfig holds creds
type CredsConfig struct {
	Registry     map[string]map[string]string `hcl:"registry"`
	VCS          map[string]map[string]string `hcl:"vcs"`
	Secrets      map[string]map[string]string `hcl:"secrets"`
	Orchestrator map[string]map[string]string `hcl:"id"`
}

// GetRegistryCreds returns creds for the registry by the id
func (cc *CredsConfig) GetRegistryCreds(id string) map[string]string {
	return cc.Registry[id]
}

// GetVCSCreds returns creds for the vcs by the id
func (cc *CredsConfig) GetVCSCreds(id string) map[string]string {
	return cc.VCS[id]
}

// GetSecretsCreds returns creds for the secrets provider by the id
func (cc *CredsConfig) GetSecretsCreds(id string) map[string]string {
	return cc.Secrets[id]
}

// GetOrchestratorCreds returns creds for the orchestrator by the id
func (cc *CredsConfig) GetOrchestratorCreds(id string) map[string]string {
	return cc.Orchestrator[id]
}

// Merge merges other to this config.  Other takes precedence
func (cc *CredsConfig) Merge(other *CredsConfig) {
	if other == nil {
		return
	}
	cc.Registry = merge(cc.Registry, other.Registry)
	cc.VCS = merge(cc.VCS, other.VCS)
	cc.Secrets = merge(cc.Secrets, other.Secrets)
	cc.Orchestrator = merge(cc.Orchestrator, other.Orchestrator)
}

func merge(curr, newm map[string]map[string]string) map[string]map[string]string {
	if len(newm) == 0 {
		return curr
	}

	out := make(map[string]map[string]string, len(curr)+len(newm))
	for k, v := range curr {
		out[k] = v
	}

	for k, v := range newm {
		if _, ok := out[k]; !ok {
			out[k] = v
			continue
		}

		for k1, v1 := range v {
			if v1 != "" {
				out[k][k1] = v1
				continue
			}
		}

	}

	return out
}

// ReadCredsConfig reads a creds conf from the given file
func ReadCredsConfig(fpath string) (*CredsConfig, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	var cc CredsConfig
	err = hcl.Unmarshal(b, &cc)

	return &cc, err
}

// WriteCredsConfig writes the creds to the given file
func WriteCredsConfig(cc *CredsConfig, fpath string) error {
	b, err := hclencoder.Encode(cc)
	if err == nil {
		err = ioutil.WriteFile(fpath, b, 0644)
	}
	return err
}

// DefaultCredsConfig returns minimal credential config
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
