package config

import (
	"io/ioutil"

	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/hcl"
)

type ImageVersionConfig struct {
	Image          string
	Versions       []string
	DefaultVersion string
}

// Config is the application configuration
type Config struct {
	Languages  map[string]*thrapb.Language
	WebServers map[string]*ImageVersionConfig
	DataStores map[string]*ImageVersionConfig
}

// ParseFile parses a config file to a Config struct
func ParseFile(fpath string) (*Config, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	var c Config
	err = hcl.Unmarshal(bs, &c)
	return &c, err
}
