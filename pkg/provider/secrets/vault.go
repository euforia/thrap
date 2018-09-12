package secrets

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/provider"
	vault "github.com/hashicorp/vault/api"
)

const defaultPolicy = `{
	capabilities = ["create", "read", "update", "delete", "list"]
}`

const defaultTimeout = 3 * time.Second

var (
	vmap = map[uint8]map[string]string{
		1: map[string]string{
			"get":  "",
			"list": "",
		},
		2: map[string]string{
			"get":  "data",
			"list": "metadata",
		},
	}
)

// VaultSecrets is used to manage vault interactions. Currently it only support kv v1
type VaultSecrets struct {
	conf *provider.Config

	// prefix the kv is mounted on. this is to handle v1 vs v2
	prefix string

	//
	kvVersion uint8

	client *vault.Client
}

// Envionment Variables:
// VAULT_ADDR
// VAULT_TOKEN (required)
func (sec *VaultSecrets) Init(config *provider.Config) error {

	var (
		conf = vault.DefaultConfig()
		err  error
	)

	conf.Timeout = defaultTimeout

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

	// Default v1 prefix
	if sec.prefix == "" {
		sec.prefix = "secret"
	}

	if val, ok := c["version"]; ok {
		if s, ok := val.(int); ok {
			sec.kvVersion = uint8(s)
		} else {
			return errors.New("version not a integer")
		}
	}

	if sec.kvVersion == 0 {
		sec.kvVersion = 1
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
	conf.Timeout = defaultTimeout
	client, err := vault.NewClient(conf)
	if err != nil {
		return nil, err
	}

	if sec.conf.Addr != "" {
		client.SetAddress(sec.conf.Addr)
	}

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
	// path := sec.SecretsPath(startPath)
	// return sec.recursiveGet(path)
	return sec.recursiveGet(startPath)
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

// List returns a list of all keys under the prefix
func (sec *VaultSecrets) List(prefix string) ([]string, error) {
	// pre := sec.listPath(prefix)
	return sec.list(prefix)
}

func (sec *VaultSecrets) listPath(key string) string {
	return filepath.Join(sec.prefix, vmap[sec.kvVersion]["list"], key)
}

func (sec *VaultSecrets) list(prefix string) ([]string, error) {
	vlt := sec.client.Logical()
	resp, err := vlt.List(sec.listPath(prefix))
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

// Set a given key to the given value
func (sec *VaultSecrets) Set(key string, value map[string]interface{}) error {
	var (
		path = sec.SecretsPath(key)
		val  map[string]interface{}
	)

	if sec.kvVersion == 2 {
		val = map[string]interface{}{"data": value}
	} else {
		val = value
	}

	vlt := sec.client.Logical()
	_, err := vlt.Write(path, val)
	return err
}

func (sec *VaultSecrets) Get(key string) (map[string]interface{}, error) {
	// path := sec.SecretsPath(key)
	return sec.get(key)
}

func (sec *VaultSecrets) Delete(key string) error {
	path := sec.SecretsPath(key)
	vlt := sec.client.Logical()
	_, err := vlt.Delete(path)
	return err
}

func (sec *VaultSecrets) get(key string) (map[string]interface{}, error) {
	vlt := sec.client.Logical()
	resp, err := vlt.Read(sec.SecretsPath(key))
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	if sec.kvVersion == 1 {
		return resp.Data, nil
	}

	// Assume v2
	v := resp.Data["data"]
	d, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data not a map: %v", v)
	}

	return d, nil
}

// SecretsPath returns the data path for th key
func (sec *VaultSecrets) SecretsPath(key string) string {
	return filepath.Join(sec.prefix, vmap[sec.kvVersion]["get"], key)
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
