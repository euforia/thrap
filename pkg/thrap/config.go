package thrap

import (
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"io/ioutil"
	"log"

	"github.com/euforia/thrap/pkg/config"
	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/storage"
)

var (
	errDataDirMissing   = errors.New("data directory missing")
	errConfigDirMissing = errors.New("config directory missing")
)

// Config holds the configuration with which a Thrap instance is initialized
type Config struct {
	// Additional provider config on top of the one that is loaded
	// from the config dir.
	Providers *config.Config
	// Load creds
	Credentials *credentials.Credentials
	// Profile storage
	Profiles storage.ProfileStorage
	// Projects store
	Projects storage.ProjectStorage
	//
	Deployments storage.DeploymentStorage
	// Overall logger
	Logger *log.Logger
	// Data directory. This must exist
	DataDir string
	// Config directory
	ConfigDir string
	// Hash function for versioning
	HashFunc func() hash.Hash
}

// Validate checks required fields and sets defaults where ever possible.  It
// returns an error if any fields are missing
func (conf *Config) Validate() error {
	if conf.DataDir == "" {
		return errDataDirMissing
	}
	if conf.ConfigDir == "" {
		return errConfigDirMissing
	}

	if conf.Logger == nil {
		conf.Logger = DefaultLogger(ioutil.Discard)
	}
	if conf.HashFunc == nil {
		conf.HashFunc = sha256.New
	}

	return nil
}

// DefaultLogger returns a default logger with the underlying writer
func DefaultLogger(w io.Writer) *log.Logger {
	return log.New(w, "", log.LstdFlags|log.Lmicroseconds)
}
