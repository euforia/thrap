package thrap

import (
	"fmt"
	"path/filepath"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/loader"
	"github.com/euforia/thrap/pkg/provider/iam"
	"github.com/euforia/thrap/pkg/provider/orchestrator"
	"github.com/euforia/thrap/pkg/provider/secrets"
	"github.com/euforia/thrap/pkg/storage"
)

func (t *Thrap) init() error {
	err := t.initIAM()
	if err != nil {
		return err
	}

	if err = t.initRegistries(); err != nil {
		return err
	}

	t.initOrchestrators()

	return nil
}

func (t *Thrap) initStorage(conf *Config) error {
	if conf.Profiles == nil {
		ppath := filepath.Join(conf.ConfigDir, profilesFile)
		t.log.Printf("Loading profiles from: %s", ppath)
		profiles, err := storage.ReadHCLFileProfileStorage(ppath)
		if err != nil {
			return err
		}
		conf.Profiles = profiles
	}

	t.profiles = conf.Profiles

	profs := t.profiles.List()
	for _, p := range profs {
		t.log.Printf("Profile loaded: %s", p.Name)
	}

	var err error
	t.ds, err = kvdb.NewBadgerDatastore(filepath.Join(conf.DataDir, dbDir), t.log)

	return err
}

func (t *Thrap) initSecretBackend(id string, creds map[string]string) (secrets.Secrets, error) {
	conf := t.conf.Secrets[id]
	if conf == nil {
		return nil, fmt.Errorf("secret backend not found: %s", id)
	}

	if !loader.AddCredsToProviderConfig(conf, creds) {
		t.log.Printf("No credentials for secrets: %s", id)
	}

	return secrets.New(conf)
}

func (t *Thrap) initIAM() error {
	iams := make(map[string]iam.IAM)

	for k, conf := range t.conf.IAM {
		creds := t.creds.SecretsCreds(k)
		if !loader.AddCredsToProviderConfig(conf, creds) {
			t.log.Printf("No credentials for IAM: %s", k)
		}

		ia, err := iam.New(conf)
		if err != nil {
			return err
		}
		iams[k] = ia
		t.log.Println("IAM loaded:", k)
	}

	t.iams = iams

	return nil
}

// load all configured orchestrators
func (t *Thrap) initOrchestrators() {
	orchs := make(map[string]orchestrator.Orchestrator)

	for k, conf := range t.conf.Orchestrator {
		creds := t.creds.OrchestratorCreds(k)
		if !loader.AddCredsToProviderConfig(conf, creds) {
			t.log.Printf("No credentials for orchestrator: %s", k)
		}

		orch, err := orchestrator.New(conf)
		if err == nil {
			orchs[k] = orch
			t.log.Println("Orchestrator loaded:", k)
			continue
		}

		t.log.Printf("Orchestrator '%s' not loaded: %v", k, err)

	}
	t.orchs = orchs
}

// initRegistries initialized all registries using the supplied configs
func (t *Thrap) initRegistries() error {
	regs, err := loader.LoadRegistries(t.conf.Registry, t.creds)
	if err != nil {
		return err
	}
	t.regs = regs

	for k := range t.regs {
		t.log.Println("Registry loaded:", k)
	}

	return nil
}
