package core

import (
	"crypto/ecdsa"
	"log"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/crt"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/pkg/config"
	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/provider/registry"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vcs"
)

var (
	errProviderNotConfigured = errors.New("provider not configured")
	errPacksDirMissing       = errors.New("packs directory missing")
	errDataDirMissing        = errors.New("data directory missing")
	errOrchNotLoaded         = errors.New("orchestrator not loaded")
	errRegNotLoaded          = errors.New("registry not loaded")
)

const (
	// Temporary default
	defaultPacksRepoURL = "https://github.com/euforia/thrap-packs.git"
)

// Core is the thrap core
type Core struct {
	conf  *config.Config
	creds *credentials.Credentials

	// Remote VCS github etc.
	vcs vcs.VCS

	// Loaded registries
	regs map[string]registry.Registry

	// Deployment orchestrator
	orchs map[string]orchestrator.Orchestrator

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
	c := &Core{}
	err := c.loadConfigs(conf)
	if err != nil {
		return nil, err
	}

	// Init CRT
	c.crt, err = crt.NewDocker()
	if err != nil {
		return nil, err
	}

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

// Config returns the currently loaded config.  This is the merged global and
// local config
func (core *Core) Config() *config.Config {
	return core.conf
}

// Packs returns a pack instance containing the currently loaded packs
func (core *Core) Packs() *packs.Packs {
	return core.packs
}

// StackStorage returns the underlying persistent store
func (core *Core) StackStorage() StackStorage {
	return core.sst
}

// Stack returns a Stack instance that can be used to perform operations
// against a stack.  It is loaded pased on the profile provided.  All
// stack fields or constructed based on the profile
func (core *Core) Stack(profile *thrapb.Profile) (*Stack, error) {
	orch, ok := core.orchs[profile.Orchestrator]
	if !ok {
		return nil, errors.Wrap(errOrchNotLoaded, profile.Orchestrator)
	}

	stack := &Stack{
		crt:   core.crt,
		orch:  orch,
		conf:  core.conf.Clone(),
		vcs:   core.vcs,
		packs: core.packs,
		sst:   core.sst,
		log:   core.log,
	}

	// The registry may be empty for local builds
	if profile.Registry != "" {
		reg, ok := core.regs[profile.Registry]
		if !ok {
			return nil, errors.Wrap(errRegNotLoaded, profile.Registry)
		}
		stack.reg = reg
	}

	return stack, nil
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
