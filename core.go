package thrap

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"

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
)

type CoreConfig struct {
	*config.ThrapConfig
	Creds  *config.CredsConfig
	Logger *log.Logger
}

type Core struct {
	conf  *config.ThrapConfig
	creds *config.CredsConfig

	vcs  vcs.VCS                   // remote vcs
	reg  registry.Registry         // remote registry
	sec  secrets.Secrets           // remote secrets
	orch orchestrator.Orchestrator // remote orchestrator

	sst *store.StackStore    // local store
	ist *store.IdentityStore // local store

	log *log.Logger
}

// NewCore loads the core engine with the global configs
func NewCore(conf *CoreConfig) (*Core, error) {
	gconf, err := ReadGlobalConfig()
	if err != nil {
		return nil, err
	}
	gconf.Merge(conf.ThrapConfig)

	creds, err := ReadGlobalCreds()
	if err != nil {
		return nil, err
	}
	creds.Merge(creds)

	c := &Core{
		conf:  gconf,
		creds: creds,
		log:   conf.Logger,
	}

	if c.log == nil {
		c.log = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}

	err = c.initProviders()
	if err == nil {
		return nil, err
	}

	err = c.initStores()

	return c, err
}

func (core *Core) initProviders() (err error) {
	if err = core.initVCS(); err != nil {
		return err
	}
	if err = core.initRegistry(); err != nil {
		return err
	}

	return core.initSecrets()
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
		core.log.Println("VCS loaded:", vc.ID)
	}
	return err
}

func (core *Core) initRegistry() (err error) {
	rc := core.conf.GetDefaultRegistry()
	if rc == nil || rc.ID == "" {
		return nil
	}

	rconf := registry.DefaultConfig()
	rconf.Provider = rc.ID
	for k, v := range rc.Config {
		rconf.Conf[k] = v
	}
	core.reg, err = registry.New(rconf)
	if err == nil {
		core.log.Println("Registry loaded:", rc.ID)
	}

	return err
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
		core.log.Println("Secrets loaded:", sc.ID)
	}
	return err
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

	if !verifySignature(ident.PublicKey, shash, ident.Signature) {
		return nil, errors.New("signature verification failed")
	}

	sident.Signature = ident.Signature

	resp, _, err := core.ist.Update(sident)
	if err == nil {
		core.log.Printf("User registered user=%s", ident.ID)
	}
	return resp, err
}

func (core *Core) RegisterStack(stack *thrapb.Stack) (*thrapb.Stack, []*ActionReport, error) {
	errs := stack.Validate()
	if len(errs) > 0 {
		return nil, nil, utils.FlattenErrors(errs)
	}

	stack, _, err := core.sst.Create(stack)
	if err != nil {
		return nil, nil, err
	}

	reports := core.ensureStackResources(stack)

	// Temp
	for _, r := range reports {
		fmt.Printf("%v '%v'\n", r.Action, r.Error)
	}

	return stack, reports, err
}

func (core *Core) ensureStackResources(stack *thrapb.Stack) []*ActionReport {
	report := core.createVcsRepo(stack)
	reports := core.ensureComponentResources(stack)

	// Concat results from both of the above
	return append([]*ActionReport{report}, reports...)
}

// returns the report and an error if any component failed
func (core *Core) ensureComponentResources(stack *thrapb.Stack) []*ActionReport {

	r := make([]*ActionReport, 0, len(stack.Components))
	for i, comp := range stack.Components {

		if !comp.IsBuildable() {
			continue
		}

		r = append(r, core.createRegistryRepo(stack.ID, i))

		if comp.HasSecrets() {
			r = append(r, core.createSecrets(i))
		}

	}

	return r
}

func (core *Core) createVcsRepo(stack *thrapb.Stack) *ActionReport {

	er := &ActionReport{Action: NewAction("create", "vcs.repo."+core.vcs.ID(), stack.ID)}

	if core.vcs == nil {
		er.Error = errProviderNotConfigured
		return er
	}

	var vcsOpt vcs.Option
	repo := &vcs.Repository{
		Name:        stack.ID,
		Description: stack.Description,
	}
	vc := core.conf.VCS[core.vcs.ID()]
	if vc.Repo.Owner != "" {
		repo.Owner = vc.Repo.Owner
	}

	er.Data, er.Error = core.vcs.Create(repo, vcsOpt)

	return er
}

// createSecrets performs actions for secret creation for the given component
// id
func (core *Core) createSecrets(compID string) *ActionReport {
	return &ActionReport{
		Action: NewAction("create", "secrets", compID),
		Error:  errors.New("to be implemented"),
	}
}

// createRegistryRepo performs actions for registry creation
func (core *Core) createRegistryRepo(stackID, compID string) *ActionReport {
	repoName := stackID + "/" + compID
	er := &ActionReport{
		Action: NewAction("create", "registry.repo", compID),
	}

	if core.reg == nil {
		er.Error = errProviderNotConfigured
		return er
	}

	// Create registry repo
	er.Data, er.Error = core.reg.Create(repoName)

	return er
}
