package iam

import (
	"sync"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/pkg/provider/secrets"
	"github.com/pkg/errors"

	vault "github.com/hashicorp/vault/api"
)

var (
	errIAMProviderMissing      = errors.New("IAM provider missing")
	errIAMProviderNotSupported = errors.New("IAM provider not supported")
)

// IAM implements an identity and access manager
type IAM interface {
	Authenticate(token string) (interface{}, error)
}

func New(conf *provider.Config) (IAM, error) {
	switch conf.Provider {
	case "vault":
		vlt := &secrets.VaultSecrets{}
		err := vlt.Init(conf)
		if err == nil {
			return newVaultIAM(vlt), nil
		}

		return nil, err

	case "":
		return nil, errors.Wrap(errIAMProviderMissing, conf.Provider)
	}

	return nil, errors.Wrap(errIAMProviderNotSupported, conf.Provider)
}

type vaultIAM struct {
	mu sync.RWMutex
	m  map[string]*vault.Secret

	// Vault auth backend
	vault *secrets.VaultSecrets
}

func newVaultIAM(v *secrets.VaultSecrets) *vaultIAM {
	return &vaultIAM{
		vault: v,
		m:     make(map[string]*vault.Secret),
	}
}

func (mgr *vaultIAM) Authenticate(token string) (interface{}, error) {

	mgr.mu.RLock()
	if val, ok := mgr.m[token]; ok {
		defer mgr.mu.RUnlock()
		return val, nil
	}
	mgr.mu.RUnlock()

	resp, err := mgr.vault.Authenticate(token)
	if err != nil {
		return nil, err
	}

	mgr.mu.Lock()
	mgr.m[token] = resp
	mgr.mu.Unlock()

	return resp, nil
}
