package secrets

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/euforia/thrap/pkg/provider"
	vault "github.com/hashicorp/vault/api"
)

// func defaultPolicies() []string {
// 	return []string{"create", "read", "update", "delete"}
// }
const defaultPolicy = `{
	capabilities = ["create", "read", "update", "delete", "list"]
}`

// only support kv v1
type vaultSecrets struct {
	conf *provider.Config

	// prefix the kv is mounted on. this is to handle v1 vs v2
	prefix string

	client *vault.Client
}

// Envionment Variables:
// VAULT_ADDR
// VAULT_TOKEN (required)
func (sec *vaultSecrets) Init(config *provider.Config) error {
	// Default v1 prefix
	sec.prefix = "secret"

	var (
		conf = vault.DefaultConfig()
		err  error
	)

	sec.client, err = vault.NewClient(conf)
	if err != nil {
		return err
	}

	sec.conf = config

	if config.Addr != "" {
		sec.client.SetAddress(config.Addr)
	}

	c := config.Config
	if val, ok := c["prefix"]; ok {
		if s, ok := val.(string); ok {
			sec.prefix = s
		} else {
			return errors.New("prefix not a string")
		}
	}

	if val, ok := c["token"]; ok {
		if s, ok := val.(string); ok {
			sec.client.SetToken(s)
		}
	}

	return err
}

func (sec *vaultSecrets) RecursiveGet(startPath string) (map[string]map[string]interface{}, error) {
	out := make(map[string]map[string]interface{})
	kv, err := sec.Get(startPath)
	if err == nil {
		out[startPath] = kv
	}

	keys, err := sec.List(startPath + "/")
	if err != nil {
		fmt.Println(err)
		return out, nil
	}

	for _, key := range keys {
		k := filepath.Join(startPath, key)

		if key[len(key)-1] != '/' {
			kv, err := sec.Get(k)
			if err != nil {
				fmt.Println(err)
				continue
			}
			out[k] = kv

		} else {
			kvs, err := sec.RecursiveGet(k)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for k, v := range kvs {
				out[k] = v
			}
		}
	}

	return out, nil
}

func (sec *vaultSecrets) List(prefix string) ([]string, error) {
	pre := sec.SecretsPath(prefix)

	vlt := sec.client.Logical()
	// pre := filepath.Join("secret", prefix)
	resp, err := vlt.List(pre)
	if err != nil {
		return nil, err
	}

	list, _ := extractListData(resp)
	return list, nil
}

// Create creates an all access policy under the path of vaultKeyPrefix/name
func (sec *vaultSecrets) Create(path string) error {
	rules := `path "` + sec.SecretsPath(path) + `/*" ` + defaultPolicy

	vlt := sec.client.Sys()
	return vlt.PutPolicy(path, rules)
}

func (sec *vaultSecrets) Set(key string, value map[string]interface{}) error {
	path := sec.SecretsPath(key)
	// req := map[string]interface{}{
	// 	"data": value,
	// }

	vlt := sec.client.Logical()
	_, err := vlt.Write(path, value)
	return err
}

func (sec *vaultSecrets) Get(key string) (map[string]interface{}, error) {
	path := sec.SecretsPath(key)

	vlt := sec.client.Logical()
	resp, err := vlt.Read(path)
	if err == nil {
		if resp == nil {
			return nil, nil
		}
		// fmt.Println(path)
		return resp.Data, nil
	}

	return nil, err
}

func (sec *vaultSecrets) SecretsPath(key string) string {
	return filepath.Join(sec.prefix, key)
}

func extractListData(secret *vault.Secret) ([]string, bool) {
	if secret == nil || secret.Data == nil {
		return nil, false
	}

	k, ok := secret.Data["keys"]
	if !ok || k == nil {
		return nil, false
	}

	i, ok := k.([]interface{})
	if !ok {
		return nil, false
	}

	out := make([]string, len(i))
	for j, v := range i {
		out[j] = v.(string)
	}
	return out, true
}
