package thrap

import (
	"context"
	"fmt"
	"hash"
	"log"
	"path/filepath"

	"github.com/euforia/kvdb"
	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/config"
	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/provider/iam"
	"github.com/euforia/thrap/pkg/provider/orchestrator"
	"github.com/euforia/thrap/pkg/provider/registry"
	"github.com/euforia/thrap/pkg/provider/secrets"
	"github.com/euforia/thrap/pkg/storage"
)

// ContextKey is used for go context keys
type ContextKey string

const (
	// IAMContextKey represents the context used for IAM data
	IAMContextKey ContextKey = "iam"
	// CredsContextKey represents the context used for credentials
	CredsContextKey ContextKey = "creds"
)

const (
	descTableKey = "descriptor"
)

const (
	dbDir               = "kvdb"
	profilesFile        = "profiles.hcl"
	providersConfigFile = "config.hcl"
)

var (
	errIAMNotEnabled = errors.New("IAM not enabled")
)

// Thrap is the core engine
type Thrap struct {
	// Config of all providers
	conf *config.Config

	// Credentials for all providers
	creds *credentials.Credentials

	// Loaded registries based on config
	regs map[string]registry.Registry

	// Loaded orchestrators
	orchs map[string]orchestrator.Orchestrator

	// Loaded secrets
	secs map[string]secrets.Secrets

	// Load iam's
	iams map[string]iam.IAM

	// Profiles
	profiles storage.ProfileStorage

	// Project store
	projects storage.ProjectStorage

	// Deployment descriptors
	descs storage.DeployDescStorage

	// Deployments
	deploys storage.DeploymentStorage

	// Datastore
	ds kvdb.Datastore

	// hash function
	hashFunc func() hash.Hash

	// logger
	log *log.Logger
}

// New returns a new thrap instance
func New(conf *Config) (*Thrap, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	t := &Thrap{
		creds:    conf.Credentials,
		log:      conf.Logger,
		hashFunc: conf.HashFunc,
		projects: conf.Projects,
		deploys:  conf.Deployments,
	}

	t.conf, err = config.ReadConfig(filepath.Join(conf.ConfigDir, providersConfigFile))
	if err != nil {
		return nil, err
	}

	// Merge additional with persisted
	t.conf.Merge(conf.Providers)

	err = t.initStorage(conf)
	if err == nil {
		err = t.init()
	}

	return t, err
}

// Profiles returns the load profile store
func (t *Thrap) Profiles() storage.ProfileStorage {
	return t.profiles
}

// Engine returns a runnable engine with the profile specified.
func (t *Thrap) Engine(ctx context.Context, profileID string) (Engine, error) {
	profile, err := t.profiles.Get(profileID)
	if err != nil {
		return nil, err
	}

	eng := &engine{
		profile:   profile,
		o:         t.orchs[profile.Orchestrator],
		r:         t.regs[profile.Registry],
		dockerHub: t.regs["dockerhub"],
		log:       t.log,
	}

	// We load the secrets engines based on the credentials
	credCtx := ctx.Value(CredsContextKey)
	if credCtx != nil {
		if creds, ok := credCtx.(*credentials.Credentials); ok {
			eng.sec, err = t.initSecretBackend(profile.Secrets, creds.SecretsCreds(profile.Secrets))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("wrong credentials data type")
		}
	}

	return eng, eng.validateResourceProviders()
}

// IAMEnabled returns true if IAM is enabled
func (t *Thrap) IAMEnabled() bool {
	return len(t.iams) > 0
}

// Authenticate authenticates against the profile with the credentials supplied.
func (t *Thrap) Authenticate(profile, token string) (interface{}, error) {
	if !t.IAMEnabled() {
		return nil, errIAMNotEnabled
	}
	ia, ok := t.iams[profile]
	if !ok {
		return nil, fmt.Errorf("profile not found: %s", profile)
	}
	return ia.Authenticate(token)
}
