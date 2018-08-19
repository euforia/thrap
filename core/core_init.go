package core

import (
	"os"
	"path/filepath"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/registry"
	"github.com/euforia/thrap/secrets"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
)

// load all thrap configs from the core config.
func (core *Core) loadConfigs(conf *Config) error {
	err := conf.Validate()
	if err != nil {
		return err
	}

	core.log = conf.Logger

	conf.DataDir, err = utils.GetAbsPath(conf.DataDir)
	if err != nil {
		return err
	}

	if !utils.FileExists(conf.DataDir) {
		return errDataDirMissing
	}

	cfile := filepath.Join(conf.DataDir, consts.ConfigFile)
	gconf, err := config.ReadConfig(cfile)
	if err != nil {
		return err
	}
	// Merge user supplied with global for this core instance
	gconf.Merge(conf.Config)
	core.conf = gconf

	// Creds
	credsFile := filepath.Join(conf.DataDir, consts.CredsFile)
	creds, err := config.ReadCredsConfig(credsFile)
	if err != nil {
		return err
	}

	creds.Merge(conf.Creds)
	core.creds = creds

	return nil
}

func (core *Core) initKeyPair(dir string) (err error) {
	name := filepath.Join(dir, "ecdsa256")
	core.kp, err = utils.LoadECDSAKeyPair(name)
	return
}

func (core *Core) initPacks(dir string) error {
	pks, err := packs.New(dir)
	if err != nil {
		return err
	}

	core.packs = pks
	if !utils.FileExists(pks.Dir()) {
		err = core.packs.Load(defaultPacksRepoURL)
	}
	return err
}

func (core *Core) initProviders() (err error) {
	if err = core.initVCS(); err != nil {
		return err
	}
	if err = core.initRegistries(); err != nil {
		return err
	}
	if err = core.initSecrets(); err != nil {
		return err
	}
	err = core.initOrchestrators()
	return err
}

func (core *Core) initVCS() (err error) {
	vc := core.conf.DefaultVCS()
	vconf := &vcs.Config{
		Provider: vc.ID,
		Conf:     map[string]interface{}{"username": vc.Username},
	}

	vcreds := core.creds.GetVCSCreds(vc.ID)
	for k, v := range vcreds {
		vconf.Conf[k] = v
	}

	core.vcs, err = vcs.New(vconf)
	core.log.Printf("VCS loaded: id=%s user=%s", core.vcs.ID(), vc.Username)
	return err
}

func (core *Core) initSecrets() (err error) {
	sc := core.conf.DefaultSecrets()

	screds := core.creds.GetSecretsCreds(sc.ID)
	sconf := &secrets.Config{
		Provider: sc.ID,
		Conf:     make(map[string]interface{}),
	}
	sconf.Conf["addr"] = sc.Addr
	for k, v := range screds {
		sconf.Conf[k] = v
	}

	core.sec, err = secrets.New(sconf)

	return err
}

// load all configured registries
func (core *Core) initRegistries() error {
	core.regs = make(map[string]registry.Registry)

	for k, rc := range core.conf.Registry {
		conf := core.creds.GetRegistryCreds(k)
		if conf != nil {
			if rc.Config == nil {
				rc.Config = make(map[string]interface{})
			}

			for k, v := range conf {
				rc.Config[k] = v
			}
		} else {
			core.log.Printf("Credentials not found for registry: %s", k)
		}

		reg, err := registry.New(rc)
		if err != nil {
			return err
		}
		core.log.Println("Registry loaded:", k)
		core.regs[k] = reg
	}

	return nil
}

// load all configured orchestrators
func (core *Core) initOrchestrators() error {
	orchs := make(map[string]orchestrator.Orchestrator, len(core.conf.Orchestrator))
	for k := range core.conf.Orchestrator {
		// TODO: add config
		conf := &orchestrator.Config{Provider: k}
		orch, err := orchestrator.New(conf)
		if err != nil {
			return err
		}
		orchs[k] = orch
		core.log.Println("Orchestrator loaded:", k)
	}
	core.orchs = orchs

	return nil
}

func (core *Core) initStores(datadir string) error {
	dbdir := filepath.Join(datadir, "db")

	if !utils.FileExists(dbdir) {
		core.log.Println("Initializing new db:", dbdir)
		os.MkdirAll(dbdir, 0755)
	}

	db, err := store.NewBadgerDB(dbdir)
	if err != nil {
		return err
	}

	core.sst = store.NewBadgerStackStorage(db)
	core.ist = store.NewBadgerIdentityStorage(db)

	return nil
}
