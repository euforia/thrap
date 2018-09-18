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
