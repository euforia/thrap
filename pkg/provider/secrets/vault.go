package secrets

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/provider"
	vault "github.com/hashicorp/vault/api"
)

const defaultPolicy = `{
	capabilities = ["create", "read", "update", "delete", "list"]
}`

// VaultSecrets is used to manage vault interactions. Currently it only support kv v1
type VaultSecrets struct {
	conf *provider.Config

	// prefix the kv is mounted on. this is to handle v1 vs v2
	prefix string

	client *vault.Client
}

// Envionment Variables:
// VAULT_ADDR
// VAULT_TOKEN (required)
func (sec *VaultSecrets) Init(config *provider.Config) error {
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

// Authenticate authenticates the token returning details about it.
func (sec *VaultSecrets) Authenticate(token string) (*vault.Secret, error) {
	// We use a complete new client so as to not mess with the current configured
	// client loaded via a profile
	conf := vault.DefaultConfig()
	client, err := vault.NewClient(conf)
	if err != nil {
		return nil, err
	}
	client.SetAddress(sec.conf.Addr)
	client.SetToken(token)

	// Validate token
	secret, err := client.Auth().Token().LookupSelf()
	if err != nil {
		return nil, errors.Wrap(err, "error looking up token")
	}

	if secret == nil {
		return nil, fmt.Errorf("empty response from lookup-self")
	}

	return secret, nil
}

func (sec *VaultSecrets) RecursiveGet(startPath string) (map[string]map[string]interface{}, error) {
	path := sec.SecretsPath(startPath)
	return sec.recursiveGet(path)
}

func (sec *VaultSecrets) recursiveGet(startPath string) (map[string]map[string]interface{}, error) {
	out := make(map[string]map[string]interface{})
	kv, err := sec.get(startPath)
	if err == nil {
		out[startPath] = kv
	}

	keys, err := sec.list(startPath + "/")
	if err != nil {
		fmt.Println(err)
		return out, nil
	}
	fmt.Println(startPath)
	fmt.Println(keys)

	for _, key := range keys {
		k := filepath.Join(startPath, key)

		if key[len(key)-1] != '/' {
			kv, err := sec.get(k)
			if err != nil {
				fmt.Println(err)
				continue
			}
			out[k] = kv

		} else {
			kvs, err := sec.recursiveGet(k)
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

func (sec *VaultSecrets) List(prefix string) ([]string, error) {
	pre := sec.SecretsPath(prefix)
	return sec.list(pre)
}

func (sec *VaultSecrets) list(prefix string) ([]string, error) {
	vlt := sec.client.Logical()
	resp, err := vlt.List(prefix)
	if err != nil {
		return nil, err
	}

	list, _ := extractListData(resp)
	return list, nil
}

// Create creates an all access policy under the path of vaultKeyPrefix/name
func (sec *VaultSecrets) Create(path string) error {
	rules := `path "` + sec.SecretsPath(path) + `/*" ` + defaultPolicy

	vlt := sec.client.Sys()
	return vlt.PutPolicy(path, rules)
}

func (sec *VaultSecrets) Set(key string, value map[string]interface{}) error {
	path := sec.SecretsPath(key)
	// req := map[string]interface{}{
	// 	"data": value,
	// }

	vlt := sec.client.Logical()
	_, err := vlt.Write(path, value)
	return err
}

func (sec *VaultSecrets) Get(key string) (map[string]interface{}, error) {
	path := sec.SecretsPath(key)
	return sec.get(path)
}

func (sec *VaultSecrets) get(key string) (map[string]interface{}, error) {
	vlt := sec.client.Logical()
	resp, err := vlt.Read(key)
	if err == nil {
		if resp == nil {
			return nil, nil
		}
		return resp.Data, nil
	}

	return nil, err
}

func (sec *VaultSecrets) SecretsPath(key string) string {
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
