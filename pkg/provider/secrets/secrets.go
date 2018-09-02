package secrets

import (
	"fmt"

	"github.com/euforia/thrap/pkg/provider"
)

// Secrets implements an interface to encrypt and store secrets
type Secrets interface {
	// This is called to initializes the secrets provider at startup
	Init(*provider.Config) error

	// Path to a given secrets
	SecretsPath(key string) string

	// Create should bootstrap secrets under the given name. This includes
	// creating any necessary access policies. It takes an optional initial
	// set for initialization
	Create(path string) error

	// List keys under the prefix
	List(prefix string) ([]string, error)

	// Get a single path
	Get(path string) (map[string]interface{}, error)

	// Recursively gets all key-value pairs under a path including the root
	RecursiveGet(path string) (map[string]map[string]interface{}, error)

	// Set key-value pairs to the path key
	Set(path string, data map[string]interface{}) error
}

// New initializes a new secrets engine based on the config
func New(conf *provider.Config) (Secrets, error) {
	var (
		sec Secrets
		err error
	)

	switch conf.Provider {
	case "vault":
		sec = &VaultSecrets{}

	// case "file":
	// 	sec = NewFileSecrets()

	default:
		err = fmt.Errorf("unsupported secrets provider: '%s'", conf.Provider)

	}

	if err != nil {
		return nil, err
	}

	err = sec.Init(conf)

	return sec, err
}
