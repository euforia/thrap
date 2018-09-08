package loader

import (
	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/pkg/provider/registry"
)

func LoadRegistries(rconfs map[string]*provider.Config, creds *credentials.Credentials) (map[string]registry.Registry, error) {

	regs := make(map[string]registry.Registry)

	for k, regConf := range rconfs {

		creds := creds.RegistryCreds(k)
		AddCredsToProviderConfig(regConf, creds)

		reg, err := registry.New(regConf)
		if err != nil {
			return nil, err
		}

		regs[k] = reg
	}

	return regs, nil

}

func AddCredsToProviderConfig(conf *provider.Config, creds map[string]string) bool {
	if len(creds) == 0 {
		return false
	}

	if conf.Config == nil {
		conf.Config = make(map[string]interface{})
	}

	for k, v := range creds {
		conf.Config[k] = v
	}

	return true
}
