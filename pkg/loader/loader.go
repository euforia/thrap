package loader

import "github.com/euforia/thrap/pkg/provider"

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
