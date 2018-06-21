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
	errProviderNotSpecified = errors.New("provider not specified")
)

type Core struct {
	conf *config.ThrapConfig

	vcs  vcs.VCS                   // remote vcs
	reg  registry.Registry         // remote registry
	sec  secrets.Secrets           // remote secrets
	orch orchestrator.Orchestrator // remote orchestrator

	sst *store.StackStore    // local store
	ist *store.IdentityStore // local store

	log *log.Logger
}

func NewCore(logger *log.Logger) (*Core, error) {
	conf, err := ReadGlobalConfig()
	if err != nil {
		return nil, err
	}
	creds, err := ReadGlobalCreds()
	if err != nil {
		return nil, err
	}

	c := &Core{conf: conf}
	if logger == nil {
		c.log = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}

	err = c.initResourceProviders(conf, creds)
	if err != nil {
		return nil, err
	}

	err = c.initStores()

	return c, err
}

func (core *Core) initResourceProviders(conf *config.ThrapConfig, creds *config.CredsConfig) (err error) {
	vconf := &vcs.Config{
		Provider: conf.VCS.ID,
		Conf:     map[string]interface{}{"username": conf.VCS.Username},
	}

	vcreds := creds.GetVCSCreds(conf.VCS.ID)
	for k, v := range vcreds {
		vconf.Conf[k] = v
	}

	core.vcs, err = vcs.New(vconf)
	if err != nil {
		return err
	}

	if conf.Registry.ID != "" {
		rconf := registry.DefaultConfig()
		rconf.Provider = conf.Registry.ID
		for k, v := range conf.Registry.Conf {
			rconf.Conf[k] = v
		}
		core.reg, err = registry.New(rconf)
		if err != nil {
			return err
		}
	}

	screds := creds.GetSecretsCreds(conf.Secrets.ID)
	sconf := &secrets.Config{
		Provider: conf.Secrets.ID,
		Conf:     make(map[string]interface{}),
	}
	sconf.Conf["addr"] = conf.Secrets.Addr
	for k, v := range screds {
		sconf.Conf[k] = v
	}

	core.sec, err = secrets.New(sconf)
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

	iobj := store.NewBadgerObjectStore(db, sha256.New, "/stack")
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

	//pubkey := makePubKeyFromBytes(elliptic.P256(), sident.PublicKey)

	shash := sident.SigHash(sha256.New())

	b58e := base58.Encode(shash)
	core.log.Printf("Verifying code: %s", b58e)

	core.log.Println(ident.Signature)
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
	reports := core.executeComponents(stack)

	// Concat results from both of the above
	return append([]*ActionReport{report}, reports...)
}

// returns the report and an error if any component failed
func (core *Core) executeComponents(stack *thrapb.Stack) []*ActionReport {

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
	var vcsOpt vcs.Option
	er := &ActionReport{Action: NewAction("create", "vcs.repo."+core.vcs.ID(), stack.ID)}
	if core.vcs != nil {
		repo := &vcs.Repository{
			Name:        stack.ID,
			Description: stack.Description,
		}
		if core.conf.VCS.Repo.Owner != "" {
			repo.Owner = core.conf.VCS.Repo.Owner
		}
		er.Data, er.Error = core.vcs.Create(repo, vcsOpt)
	} else {
		er.Error = errProviderNotSpecified
	}
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
	if core.reg != nil {
		// Create registry repo
		er.Data, er.Error = core.reg.Create(repoName)
	} else {
		er.Error = errProviderNotSpecified
	}

	return er
}
