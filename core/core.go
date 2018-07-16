package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"log"
	"os"
	"path/filepath"

	"github.com/euforia/thrap/thrapb"

	"github.com/euforia/thrap/crt"

	"github.com/euforia/thrap/packs"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/registry"
	"github.com/euforia/thrap/secrets"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
)

var (
	errProviderNotConfigured = errors.New("provider not configured")
	errPacksDirMissing       = errors.New("packs directory missing")
	errDataDirMissing        = errors.New("data directory missing")
)

const (
	// temporary default
	defaultPacksRepoURL = "https://github.com/euforia/thrap-packs.git"
)

// Config holds the core configuration
type Config struct {
	// This is the local project config merged with the global user config for the
	// instance
	*config.ThrapConfig
	// Load creds
	Creds *config.CredsConfig
	// Overall logger
	Logger *log.Logger
	// Data directory. This must exist
	DataDir string
}

func (conf *Config) Validate() error {
	if conf.DataDir == "" {
		return errDataDirMissing
	}

	if conf.Logger == nil {
		conf.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}

	return nil
}

// DefaultConfig returns a basic core config
func DefaultConfig() *Config {
	return &Config{DataDir: consts.DefaultDataDir}
}

// StackStorage is a stack storage interface
type StackStorage interface {
	Get(string) (*thrapb.Stack, error)
	Create(*thrapb.Stack) (*thrapb.Stack, error)
}

// IdentityStorage is a identity storage interface
type IdentityStorage interface {
	Get(string) (*thrapb.Identity, error)
	Create(*thrapb.Identity) (*thrapb.Identity, error)
	Update(*thrapb.Identity) (*thrapb.Identity, error)
}

// Core is the thrap core
type Core struct {
	conf  *config.ThrapConfig
	creds *config.CredsConfig

	// Remote VCS github etc.
	vcs vcs.VCS

	// Loaded registries
	regs map[string]registry.Registry

	// Secrets engine
	sec secrets.Secrets

	// Deployment orchestrator
	orch orchestrator.Orchestrator

	// Loaded extension packs
	packs *packs.Packs

	// Container runtime. Currently docker
	crt *crt.Docker

	sst StackStorage
	ist IdentityStorage

	// Load keypair. Currently 1 per core
	kp *ecdsa.PrivateKey

	// Logger
	log *log.Logger
}

// NewCore loads the core engine with the global configs
func NewCore(conf *Config) (*Core, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	conf.DataDir, err = utils.GetAbsPath(conf.DataDir)
	if err != nil {
		return nil, err
	}

	if !utils.FileExists(conf.DataDir) {
		return nil, errDataDirMissing
	}

	dkr, err := crt.NewDocker()
	if err != nil {
		return nil, err
	}

	c := &Core{
		crt: dkr,
		log: conf.Logger,
	}

	cfile := filepath.Join(conf.DataDir, consts.ConfigFile)
	gconf, err := config.ReadThrapConfig(cfile)
	if err != nil {
		return nil, err
	}

	// Merge user supplied with global for this core instance
	gconf.Merge(conf.ThrapConfig)
	c.conf = gconf

	credsFile := filepath.Join(conf.DataDir, consts.CredsFile)
	creds, err := config.ReadCredsConfig(credsFile)
	if err != nil {
		return nil, err
	}
	creds.Merge(creds)
	c.creds = creds

	err = c.initKeyPair(conf.DataDir)
	if err != nil {
		return nil, err
	}

	err = c.initPacks(filepath.Join(conf.DataDir, consts.PacksDir))
	if err != nil {
		return nil, err
	}

	err = c.initProviders()
	if err == nil {
		err = c.initStores(conf.DataDir)
	}

	return c, err
}

// Config returns the currently loaded config
func (core *Core) Config() *config.ThrapConfig {
	return core.conf
}

// Packs returns a pack instance containing the currently loaded packs
func (core *Core) Packs() *packs.Packs {
	return core.packs
}

// Stack returns a Stack instance that can be used to perform operations
// against a stack
func (core *Core) Stack() *Stack {
	return &Stack{
		regs:  core.regs,
		crt:   core.crt,
		conf:  core.conf.Clone(),
		vcs:   core.vcs,
		packs: core.packs,
		sst:   core.sst,
		log:   core.log,
	}
}

// Identity returns an Identity instance to perform operations against
// identities
func (core *Core) Identity() *Identity {
	return &Identity{
		store: core.ist,
		log:   core.log,
	}
}

// KeyPair returns the public-private key currently held by the core
func (core *Core) KeyPair() *ecdsa.PrivateKey {
	return core.kp
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
	if err = core.initRegistry(); err != nil {
		return err
	}

	if err = core.initSecrets(); err != nil {
		return err
	}

	err = core.initOrchestrator()
	return err
}

func (core *Core) initVCS() (err error) {
	vc := core.conf.GetDefaultVCS()
	vconf := &vcs.Config{
		Provider: vc.ID,
		Conf:     map[string]interface{}{"username": vc.Username},
	}

	vcreds := core.creds.GetVCSCreds(vc.ID)
	for k, v := range vcreds {
		vconf.Conf[k] = v
	}

	core.vcs, err = vcs.New(vconf)
	if err == nil {
		core.log.Println("DEBUG VCS:", vc.ID)
	}
	return err
}

func (core *Core) initRegistry() error {
	core.regs = make(map[string]registry.Registry)
	for k, rc := range core.conf.Registry {
		rconf := registry.DefaultConfig()
		rconf.Provider = rc.ID
		for k, v := range rc.Config {
			rconf.Conf[k] = v
		}

		reg, err := registry.New(rconf)
		if err != nil {
			return err
		}

		core.regs[k] = reg
		core.log.Println("DEBUG Registry:", rc.ID)
	}

	return nil
}

func (core *Core) initSecrets() (err error) {
	sc := core.conf.GetDefaultSecrets()

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
	if err == nil {
		core.log.Println("DEBUG Secrets:", sc.ID)
	}

	return err
}

func (core *Core) initOrchestrator() (err error) {
	c := core.conf.GetDefaultOrchestrator()
	conf := &orchestrator.Config{Provider: c.ID}
	core.orch, err = orchestrator.New(conf)
	if err == nil {
		core.log.Println("DEBUG Orchestrator:", core.orch.ID())
	}
	return
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

	sobj := store.NewBadgerObjectStore(db, sha256.New, "/stack")
	core.sst, err = store.NewStackStore(sobj)
	if err != nil {
		return err
	}

	iobj := store.NewBadgerObjectStore(db, sha256.New, "/identity")
	core.ist = store.NewIdentityStore(iobj)
	return nil
}

// // returns the report and an error if any component failed
// func (core *Core) ensureComponentResources(stack *thrapb.Stack) []*ActionReport {

// 	r := make([]*ActionReport, 0, len(stack.Components))
// 	for i, comp := range stack.Components {

// 		if !comp.IsBuildable() {
// 			continue
// 		}
// 		// TEMP use ecr
// 		r = append(r, core.createRegistryRepo("ecr", stack.ID, i))

// 		if comp.HasSecrets() {
// 			r = append(r, core.createSecrets(i))
// 		}

// 	}

// 	return r
// }

// func (core *Core) createVcsRepo(stack *thrapb.Stack) *ActionReport {

// 	er := &ActionReport{Action: NewAction("create", "vcs.repo."+core.vcs.ID(), stack.ID)}

// 	if core.vcs == nil {
// 		er.Error = errProviderNotConfigured
// 		return er
// 	}

// 	var vcsOpt vcs.Option
// 	repo := &vcs.Repository{
// 		Name:        stack.ID,
// 		Description: stack.Description,
// 	}
// 	vc := core.conf.VCS[core.vcs.ID()]
// 	if vc.Repo.Owner != "" {
// 		repo.Owner = vc.Repo.Owner
// 	}

// 	er.Data, er.Error = core.vcs.Create(repo, vcsOpt)

// 	return er
// }

// // createSecrets performs actions for secret creation for the given component
// // id
// func (core *Core) createSecrets(compID string) *ActionReport {
// 	return &ActionReport{
// 		Action: NewAction("create", "secrets", compID),
// 		Error:  errors.New("to be implemented"),
// 	}
// }

// // createRegistryRepo performs actions for registry creation
// func (core *Core) createRegistryRepo(regID, stackID, compID string) *ActionReport {
// 	repoName := stackID + "/" + compID
// 	er := &ActionReport{
// 		Action: NewAction("create", "registry.repo", compID),
// 	}

// 	reg, ok := core.regs[regID]
// 	if !ok {
// 		er.Error = errors.Wrap(errProviderNotConfigured, regID)
// 		return er
// 	}

// 	// Create registry repo
// 	er.Data, er.Error = reg.Create(repoName)

// 	return er
// }
