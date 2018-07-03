package secrets

import (
	"errors"

	vault "github.com/hashicorp/vault/api"
)

// func defaultPolicies() []string {
// 	return []string{"create", "read", "update", "delete"}
// }

type vaultSecrets struct {
	prefix string
	client *vault.Client
}

// Envionment Variables:
// VAULT_ADDR
// VAULT_TOKEN (required)
func (sec *vaultSecrets) Init(c map[string]interface{}) error {
	var (
		conf = vault.DefaultConfig()
		err  error
	)

	sec.client, err = vault.NewClient(conf)
	if err != nil {
		return err
	}

	if val, ok := c["addr"]; ok {
		if s, ok := val.(string); ok {
			sec.client.SetAddress(s)
		}
	}

	if val, ok := c["prefix"]; ok {
		if s, ok := val.(string); ok {
			sec.prefix = s
		} else {
			return errors.New("prefix not string")
		}
	} else {
		return errors.New("prefix missing")
	}

	if val, ok := c["token"]; ok {
		if s, ok := val.(string); ok {
			sec.client.SetToken(s)
		}
	}

	return err
}

func (sec *vaultSecrets) getOpaque() string {
	return "secret/data" + sec.prefix
}

func (sec *vaultSecrets) Set(value map[string]interface{}) error {

	path := sec.getOpaque()
	req := map[string]interface{}{
		"data": value,
	}

	vlt := sec.client.Logical()
	_, err := vlt.Write(path, req)

	return err
}

func (sec *vaultSecrets) Get() (map[string]interface{}, error) {

	path := sec.getOpaque()

	vlt := sec.client.Logical()
	resp, err := vlt.Read(path)
	if err == nil {
		return resp.Data["data"].(map[string]interface{}), nil
	}

	return nil, err
}
