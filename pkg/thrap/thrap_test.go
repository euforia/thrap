package thrap

import (
	"io/ioutil"
	"os"

	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/utils"
)

func loadTestThrap() (*Thrap, *Config, error) {
	tdir, err := ioutil.TempDir("/tmp", "foo-")
	if err != nil {
		return nil, nil, err
	}

	confDir, err := utils.GetAbsPath("../../test-fixtures")
	if err != nil {
		return nil, nil, err
	}

	creds, err := credentials.ReadCredentials("../../.thrap/creds.hcl")
	if err != nil {
		return nil, nil, err
	}

	conf := &Config{
		// Providers:   pconf,
		Credentials: creds,
		DataDir:     tdir,
		ConfigDir:   confDir,
		Logger:      DefaultLogger(os.Stderr),
	}

	th, err := New(conf)
	return th, conf, err
}
