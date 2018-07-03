package core

import (
	"crypto/sha256"
	"log"
	"math/rand"
	"os"

	"github.com/euforia/thrap/crt"

	"github.com/euforia/thrap/packs"

	"github.com/pkg/errors"

	"github.com/euforia/base58"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/registry"
	"github.com/euforia/thrap/secrets"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	errProviderNotConfigured = errors.New("provider not configured")
	errPacksDirMissing       = errors.New("packs directory missing")
)

const (
	// temporary default
	defaultPacksRepoURL = "https://github.com/euforia/thrap-packs.git"
)

// ???
// type CodeBuilder interface {
// 	Build(stackID string, comp *thrapb.Component, opts orchestrator.RequestOptions) error
// }

// Config holds the core configuration
type Config struct {
	// This is the local config merged with the global for the
	// instance
	*config.ThrapConfig
	// Load creds
	Creds *config.CredsConfig
	// Overall logger
	Logger *log.Logger
	// Directory path where packs are read and written
	PacksDir string
}

// Core is the thrap core
type Core struct {
	conf  *config.ThrapConfig
	creds *config.CredsConfig

	vcs  vcs.VCS                      // remote vcs
	regs map[string]registry.Registry // registries
	sec  secrets.Secrets              // secrets
	orch orchestrator.Orchestrator    // orchestrator

	packs *packs.Packs

	crt *crt.Docker

	sst *store.StackStore    // local store
	ist *store.IdentityStore // local store

	log *log.Logger
}

// NewCore loads the core engine with the global configs
func NewCore(conf *Config) (*Core, error) {

	dkr, err := crt.NewDocker()
	if err != nil {
		return nil, err
	}

	c := &Core{
		crt: dkr,
	}

	gconf, err := config.ReadGlobalConfig()
	if err != nil {
		return nil, err
	}
	// Merge supplied with global for this core instance
	if conf != nil {
		gconf.Merge(conf.ThrapConfig)
		c.log = conf.Logger
	}
	c.conf = gconf

	creds, err := config.ReadGlobalCreds()
	if err != nil {
		return nil, err
	}
	creds.Merge(creds)
	c.creds = creds

	if c.log == nil {
		c.log = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}

	err = c.initPacks(conf.PacksDir)
	if err != nil {
		return nil, err
	}

	err = c.initProviders()
	if err == nil {
		err = c.initStores()
	}

	return c, err
}

func (core *Core) Config() *config.ThrapConfig {
	return core.conf
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

func (core *Core) Packs() *packs.Packs {
	return core.packs
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

func (core *Core) initStores() error {
	dir := "~/" + consts.WorkDir + "/db"
	dbdir, _ := homedir.Expand(dir)

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

// RegisterIdentity registers a new identity. It returns an error if the identity exists
// or fails to register
func (core *Core) RegisterIdentity(ident *thrapb.Identity) (*thrapb.Identity, []*ActionReport, error) {
	err := ident.Validate()
	if err != nil {
		return nil, nil, err
	}
	ident.Nonce = rand.Uint64()

	er := &ActionReport{}
	er.Data, _, er.Error = core.ist.Create(ident)
	er.Action = NewAction("create", "identity", ident.ID)
	if err == nil {
		core.log.Printf("User registration request user=%s", ident.ID)
	}

	return ident, []*ActionReport{er}, er.Error
}

// ConfirmIdentity confirms a identity registration request and completes it.
// In this case the public key field is the signature from the client
func (core *Core) ConfirmIdentity(ident *thrapb.Identity) (*thrapb.Identity, error) {

	sident, _, err := core.ist.Get(ident.ID)
	if err != nil {
		return nil, err
	}

	shash := sident.SigHash(sha256.New())

	b58e := base58.Encode(shash)
	core.log.Printf("Verifying user registration code=%s", b58e)

	if !utils.VerifySignature(ident.PublicKey, shash, ident.Signature) {
		return nil, errors.New("signature verification failed")
	}

	sident.Signature = ident.Signature

	resp, _, err := core.ist.Update(sident)
	if err == nil {
		core.log.Printf("User registered user=%s", ident.ID)
	}
	return resp, err
}

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

// // RegisterStack registers a new stack. It returns an error if the stack is
// // already registered or fails to register
// func (core *Core) RegisterStack(stack *thrapb.Stack) (*thrapb.Stack, []*ActionReport, error) {
// 	errs := stack.Validate()
// 	if len(errs) > 0 {
// 		return nil, nil, utils.FlattenErrors(errs)
// 	}

// 	stack, _, err := core.sst.Create(stack)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	reports := core.ensureStackResources(stack)

// 	// Temp
// 	for _, r := range reports {
// 		fmt.Printf("%v '%v'\n", r.Action, r.Error)
// 	}

// 	return stack, reports, err
// }

// func (core *Core) ensureStackResources(stack *thrapb.Stack) []*ActionReport {
// 	report := core.createVcsRepo(stack)
// 	reports := core.ensureComponentResources(stack)

// 	// Concat results from both of the above
// 	return append([]*ActionReport{report}, reports...)
// }

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
