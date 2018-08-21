package core

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/pkg/config"
	"github.com/euforia/thrap/pkg/credentials"
)

// Config holds the core configuration
type Config struct {
	// This is the local project config merged with the global user config for the
	// instance
	*config.Config
	// Load creds
	Creds *credentials.Credentials
	// Overall logger
	Logger *log.Logger
	// Data directory. This must exist
	DataDir string
}

// Validate checks required fields and sets defaults where ever possible.  It
// returns an error if any fields are missing
func (conf *Config) Validate() error {
	if conf.DataDir == "" {
		return errDataDirMissing
	}

	if conf.Logger == nil {
		conf.Logger = DefaultLogger(ioutil.Discard)
	}

	return nil
}

// DefaultConfig returns a basic core config
func DefaultConfig() *Config {
	return &Config{DataDir: consts.DefaultDataDir}
}

// DefaultLogger returns a default logger with the underlying writer
func DefaultLogger(w io.Writer) *log.Logger {
	return log.New(w, "", log.LstdFlags|log.Lmicroseconds)
}
