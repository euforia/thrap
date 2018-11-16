package thrap

import (
	"os"

	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/utils"
)

func loadTestThrap() (*Thrap, *Config, error) {
	confDir, err := utils.GetAbsPath("../../test-fixtures")
	if err != nil {
		return nil, nil, err
	}

	creds, err := credentials.ReadCredentials("../../.thrap/creds.hcl")
	if err != nil {
		return nil, nil, err
	}

	conf := &Config{
		Credentials: creds,
		ConfigDir:   confDir,
		Logger:      DefaultLogger(os.Stderr),
	}

	th, err := New(conf)
	return th, conf, err
}

func newProjects(t *Thrap) *Projects {
	return &Projects{
		t:     t,
		store: t.store.Project(),
	}
}
