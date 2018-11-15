package core

import (
	"os"
	"path/filepath"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/pkg/config"
	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/loader"
	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/pkg/provider/vcs"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/utils"
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
	creds, err := credentials.ReadCredentials(credsFile)
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

	err = core.initOrchestrators()

	return err
}

func (core *Core) initVCS() (err error) {
	vc := core.conf.DefaultVCS()
	vconf := &provider.Config{
		Provider: vc.ID,
		Config:   map[string]interface{}{"username": vc.Username},
	}

	vcreds := core.creds.VCSCreds(vc.ID)
	for k, v := range vcreds {
		vconf.Config[k] = v
	}

	core.vcs, err = vcs.New(vconf)
	core.log.Printf("VCS loaded: id=%s user=%s", core.vcs.ID(), vc.Username)
	return err
}

// DEPRECATE load all configured registries
func (core *Core) initRegistries() error {
	regs, err := loader.LoadRegistries(core.conf.Registry, core.creds)
	if err != nil {
		return err
	}
	core.regs = regs

	for k := range core.regs {
		core.log.Println("Registry loaded:", k)
	}
	return nil
}

// load all configured orchestrators
func (core *Core) initOrchestrators() error {
	orchs := make(map[string]orchestrator.Orchestrator, len(core.conf.Orchestrator))
	for k, conf := range core.conf.Orchestrator {
		// TODO: add config
		// conf := &provider.Config{Provider: k}
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
