package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/consts"
	"github.com/hashicorp/hcl"
)

// Credentials holds creds
type Credentials struct {
	Registry     map[string]map[string]string `hcl:"registry"`
	VCS          map[string]map[string]string `hcl:"vcs"`
	Secrets      map[string]map[string]string `hcl:"secrets"`
	Orchestrator map[string]map[string]string `hcl:"orchestrator"`
}

// RegistryCreds returns creds for the registry by the id
func (cc *Credentials) RegistryCreds(id string) map[string]string {
	return cc.Registry[id]
}

// VCSCreds returns creds for the vcs by the id
func (cc *Credentials) VCSCreds(id string) map[string]string {
	return cc.VCS[id]
}

// SecretsCreds returns creds for the secrets provider by the id
func (cc *Credentials) SecretsCreds(id string) map[string]string {
	return cc.Secrets[id]
}

// OrchestratorCreds returns creds for the orchestrator by the id
func (cc *Credentials) OrchestratorCreds(id string) map[string]string {
	return cc.Orchestrator[id]
}

// Merge merges other to this config.  Other takes precedence
func (cc *Credentials) Merge(other *Credentials) {
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

	out := make(map[string]map[string]string)
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

// ReadProjectCredentials reads the creds file from the directory
func ReadProjectCredentials(dir string) (*Credentials, error) {
	filename := filepath.Join(dir, consts.WorkDir, consts.CredsFile)
	return ReadCredentials(filename)
}

// ReadCredentials reads a creds conf from the given file
func ReadCredentials(fpath string) (*Credentials, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	var cc Credentials
	err = hcl.Unmarshal(b, &cc)

	return &cc, err
}

// WriteCredentials writes the creds to the given file
func WriteCredentials(cc *Credentials, fpath string) error {
	b, err := hclencoder.Encode(cc)
	if err == nil {
		err = ioutil.WriteFile(fpath, b, 0644)
	}
	return err
}

// DefaultCredentials returns minimal credential config
func DefaultCredentials() *Credentials {
	return &Credentials{
		Registry: make(map[string]map[string]string),
		VCS: map[string]map[string]string{
			"github": map[string]string{"token": ""},
		},
		Secrets:      make(map[string]map[string]string),
		Orchestrator: make(map[string]map[string]string),
	}
}
