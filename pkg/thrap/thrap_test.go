package thrap

import (
	"io/ioutil"
	"os"

	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/storage"
	"github.com/euforia/thrap/utils"
	"github.com/hashicorp/consul/api"
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

	consulConf := api.DefaultConfig()
	consulConf.Address = "http://127.0.0.1:8500"
	pstore, err := storage.NewConsulProjectStorage(consulConf, "thrap/project")
	conf := &Config{
		// Providers:   pconf,
		Projects:    pstore,
		Credentials: creds,
		DataDir:     tdir,
		ConfigDir:   confDir,
		Logger:      DefaultLogger(os.Stderr),
	}

	th, err := New(conf)
	return th, conf, err
}
